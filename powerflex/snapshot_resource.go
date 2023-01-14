package powerflex

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snapshotResource{}
	_ resource.ResourceWithConfigure   = &snapshotResource{}
	_ resource.ResourceWithImportState = &snapshotResource{}
)

// NewSnapshotResource is a helper function to simplify the provider implementation.
func NewSnapshotResource() resource.Resource {
	return &snapshotResource{}
}

// snapshotResource is the resource implementation.
type snapshotResource struct {
	client *goscaleio.Client
}

// Metadata returns the resource type name.
func (r *snapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema defines the schema for the resource.
func (r *snapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotResourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *snapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SnapshotResourceModel
	attrState := &attrState{
		name:             false,
		volumeID:         false,
		volumeName:       false,
		desiredRetention: false,
		retentionUnit:    false,
		accessMode:       false,
		capacityUnit:     false,
		size:             false,
		sizeInKb:         false,
		sdcList:          make(map[string]*sdcAttr),
		errMsg:           make(map[string]string),
	}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	sr, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
	snapshotReqs := make([]*pftypes.SnapshotDef, 0)

	snapReq := &pftypes.SnapshotDef{
		VolumeID:     plan.VolumeID.ValueString(),
		SnapshotName: plan.Name.ValueString(),
	}
	snapshotReqs = append(snapshotReqs, snapReq)
	snapParam := &pftypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotReqs,
		// RetentionPeriodInMin: convertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString()),
		AccessMode: plan.AccessMode.ValueString(),
	}
	// create a snapshot of one volume, this requires a volume id and snapshot as parameter
	snapResps, err := sr.CreateSnapshotConsistencyGroup(snapParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			"Could not create snapshot, unexpected err: "+err.Error(),
		)
		return
	}
	// setting success update for name, volume_id, desired_retention, retention_unit and volume_name.
	attrState.name = true
	attrState.volumeID = true
	attrState.volumeName = true
	attrState.accessMode = true
	snapID := snapResps.VolumeIDList[0]
	snapResponse, err2 := r.client.GetVolume("", snapID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		// add state update here before so we don't lose the creation of snapshot
		return
	}
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap

	// setting snapshot size. default value will be equal to volume size
	if !plan.Size.IsNull() {
		vikb, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		tflog.Info(ctx, "vikb"+strconv.FormatInt(vikb, 10))
		if int64(snapResource.Volume.SizeInKb) != vikb {
			err3 := snapResource.SetVolumeSize(strconv.FormatInt(plan.Size.ValueInt64(), 10))
			if err3 != nil {
				attrState.errMsg["size/capacity_unit"] = err3.Error()
			} else {
				attrState.capacityUnit = true
				attrState.size = true
			}
		}
	}
	attrState.sizeInKb = true
	// locking the auto snapshot on finding LockedAutoSnapshot parameter as true
	if plan.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			attrState.errMsg["lock_auto_snapshot"] = err.Error()
		} else {
			attrState.lockAutoSnapshot = true
		}
	}
	// unmarshalling the plan sdclist data into sdcList struct for iterative mapping of snapshot to sdc
	sdcList := []SdcList{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcList, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range sdcList {
		if si.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", si.SdcName)
			if errA != nil {
				attrState.errMsg["sdc_map_err_"+si.SdcName] = errA.Error()
				continue
			} else {
				si.SdcID = foundsdc.Sdc.ID
			}
		}
		attrState.sdcList[si.SdcID] = &sdcAttr{
			isSdcMapped:     false,
			isSdcLimitSet:   false,
			isAccessModeSet: false,
		}
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 si.SdcID,
			AllowMultipleMappings: "true",
		}
		// mapping the snapshot to sdc
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			attrState.errMsg["sdc_map_err_"+si.SdcID] += "\n" + err3.Error()
		} else {
			attrState.sdcList[si.SdcID].isSdcMapped = true
		}
		// setting limits on mapped sdc
		smslp := pftypes.SetMappedSdcLimitsParam{
			SdcID:                si.SdcID,
			BandwidthLimitInKbps: strconv.FormatInt(int64(si.LimitBwInMbps*1024), 10),
			IopsLimit:            strconv.FormatInt(int64(si.LimitIops), 10),
		}
		err4 := snapResource.SetMappedSdcLimits(&smslp)
		if err4 != nil {
			attrState.errMsg["sdc_map_err_"+si.SdcID] += "\n" + err4.Error()
		} else {
			attrState.sdcList[si.SdcID].isSdcLimitSet = true
		}
		err5 := snapResource.SetVolumeMappingAccessMode(si.AccessMode, si.SdcID)
		if err5 != nil {
			attrState.errMsg["sdc_map_err_"+si.SdcID] += "\n" + err5.Error()
		} else {
			attrState.sdcList[si.SdcID].isAccessModeSet = true
		}
	}
	if !plan.DesiredRetention.IsNull() {
		retentionInMin := convertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString())
		errRetention := snapResource.SetSnapshotSecurity(retentionInMin)
		if errRetention != nil {
			attrState.errMsg["desired_retention/retention_unit"] = errRetention.Error()
		} else {
			attrState.desiredRetention = true
		}
	}
	attrState.retentionUnit = true
	snapResponse, err2 = r.client.GetVolume("", snapID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap = snapResponse[0]
	state := SnapshotTerraformState(snap, plan, attrState)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if len(attrState.errMsg) > 0 {
		failureAction := ""
		failureMessage := ""
		for key, value := range attrState.errMsg {
			failureAction += key + ", "
			failureMessage += key + " : " + value + "\n"
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed Actions [%v]", failureAction),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	errMsg := make(map[string]string, 0)
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap := snapResponse[0]
	refreshState(snap, &state)
	// checking for volume from which snapshot is created
	vol, errVol := r.client.GetVolume("", state.VolumeID.ValueString(), "", "", false)
	if errVol != nil {
		errMsg["volume_name"] = errVol.Error()
	} else {
		state.VolumeName = types.StringValue(vol[0].Name)
	}
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	var state SnapshotResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	errMsg := make(map[string]string, 0)
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap

	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		err := snapResource.SetVolumeName(plan.Name.ValueString())
		if err != nil {
			errMsg["name"] = err.Error()
		}
	}
	VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	plan.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
	// updating the size of the volume if there is change in plan
	if plan.VolumeSizeInKb.ValueString() != state.VolumeSizeInKb.ValueString() {
		sizeInGb, _ := strconv.Atoi(strconv.FormatInt(VSIKB, 10))
		sizeInGb = sizeInGb / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		err := snapResource.SetVolumeSize(sizeInGB)
		if err != nil {
			errMsg["size/capacity_unit"] = err.Error()
		}
	}
	// locking the snapshot in case of change in LockedAutoSnapshot state to true
	if plan.LockAutoSnapshot.ValueBool() && !state.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			errMsg["lock_auto_snapshot"] = err.Error()
		}
	}
	// unlocking the snapshot in case of change in LockedAutoSnapshot state to false
	if !plan.LockAutoSnapshot.ValueBool() && state.LockAutoSnapshot.ValueBool() {
		err := snapResource.UnlockAutoSnapshot()
		if err != nil {
			errMsg["lock_auto_snapshot"] = err.Error()
		}
	}
	planSdcList := []SdcList{}
	stateSdcList := []SdcList{}
	diags = plan.SdcList.ElementsAs(ctx, &planSdcList, true)
	resp.Diagnostics.Append(diags...)
	diags = state.SdcList.ElementsAs(ctx, &stateSdcList, true)
	resp.Diagnostics.Append(diags...)
	planSdcIds := []string{}
	stateSdcIds := []string{}
	for _, psl := range planSdcList {
		planSdcIds = append(planSdcIds, psl.SdcID)
	}
	for _, ssl := range stateSdcList {
		stateSdcIds = append(stateSdcIds, ssl.SdcID)
	}
	mapSdcIds := Difference(planSdcIds, stateSdcIds)
	unmapSdcIds := Difference(stateSdcIds, planSdcIds)
	nonchangeSdcIds := Difference(planSdcIds, mapSdcIds)

	// changing the access mode in case of change in access mode state
	if (plan.AccessMode.ValueString() == "ReadWrite") && (state.AccessMode.ValueString() == "ReadOnly") {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			errMsg["access-readwrite"] = err.Error()
		}
	}

	for _, msi := range mapSdcIds {
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msi,
			AllowMultipleMappings: "true",
		}
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			errMsg["sdc_map_err_"+msi] += "\n" + err3.Error()
		}
		for _, ssl := range planSdcList {
			if ssl.SdcID == msi {
				smslp := pftypes.SetMappedSdcLimitsParam{
					SdcID:                ssl.SdcID,
					BandwidthLimitInKbps: strconv.FormatInt(int64(ssl.LimitBwInMbps*1024), 10),
					IopsLimit:            strconv.FormatInt(int64(ssl.LimitIops), 10),
				}
				err4 := snapResource.SetMappedSdcLimits(&smslp)
				if err4 != nil {
					errMsg["sdc_map_err_"+msi] += "\n" + err4.Error()
				}
				err5 := snapResource.SetVolumeMappingAccessMode(ssl.AccessMode, ssl.SdcID)
				if err5 != nil {
					errMsg["sdc_map_err_"+msi] += "\n" + err5.Error()
				}
			}
		}
	}

	for _, usi := range unmapSdcIds {
		err4 := snapResource.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: usi,
			},
		)
		if err4 != nil {
			errMsg["sdc_unmap_err_"+usi] += "\n" + err4.Error()
		}
	}

	for _, ncsi := range nonchangeSdcIds {
		var planObj SdcList
		var stateObj SdcList

		for _, psl := range planSdcList {
			if ncsi == psl.SdcID {
				planObj = psl
				break
			}
		}
		for _, ssl := range stateSdcList {
			if ncsi == ssl.SdcID {
				stateObj = ssl
				break
			}
		}
		if (planObj.LimitIops != stateObj.LimitIops) || (planObj.LimitBwInMbps != stateObj.LimitBwInMbps) {
			smslp := pftypes.SetMappedSdcLimitsParam{
				SdcID:                planObj.SdcID,
				BandwidthLimitInKbps: strconv.FormatInt(int64(planObj.LimitBwInMbps*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(planObj.LimitIops), 10),
			}
			err4 := snapResource.SetMappedSdcLimits(&smslp)
			if err4 != nil {
				errMsg["sdc_map_err_"+planObj.SdcID] += "\n" + err4.Error()
			}
		}

		if planObj.AccessMode != stateObj.AccessMode {
			err5 := snapResource.SetVolumeMappingAccessMode(planObj.AccessMode, planObj.SdcID)
			if err5 != nil {
				errMsg["sdc_map_err_"+planObj.SdcID] += "\n" + err5.Error()
			}
		}

	}

	// changing the access mode in case of change in access mode state
	if (plan.AccessMode.ValueString() == "ReadOnly") && (state.AccessMode.ValueString() == "ReadWrite") {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			errMsg["access-readwrite"] = err.Error()
		}
	}
	if !plan.DesiredRetention.IsNull() && ((plan.RetentionUnit.ValueString() != state.RetentionUnit.String()) || (plan.DesiredRetention.ValueInt64() != state.DesiredRetention.ValueInt64())) {
		retentionInMin := convertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString())
		err := snapResource.SetSnapshotSecurity(retentionInMin)
		if err != nil {
			errMsg["desired_retention/retention_unit"] = err.Error()
		}
	}
	snapResponse, err2 = r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap = snapResponse[0]
	snapResource.Volume = snap
	refreshState(snap, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if len(errMsg) > 0 {
		failureAction := ""
		failureMessage := ""
		for key, value := range errMsg {
			failureAction += key + ", "
			failureMessage += key + " : " + value + "\n"
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed Actions [%v]", failureAction),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snapshot := goscaleio.NewVolume(r.client)
	snapshot.Volume = snapResponse[0]
	sdcsToUnmap := []SdcList{}
	diags = state.SdcList.ElementsAs(ctx, &sdcsToUnmap, true)
	resp.Diagnostics.Append(diags...)
	for _, stu := range sdcsToUnmap {
		err := snapshot.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: stu.SdcID,
			},
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Volume to SDCs",
				"Couldn't unmap volume to scs with id: "+stu.SdcID+", unexpected error: "+err.Error(),
			)
			return
		}
	}
	err := snapshot.RemoveVolume(state.RemoveMode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing Volume",
			"Couldn't remove volume "+err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *snapshotResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return 
	}
	var plan SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	errMsg := make(map[string]string, 0)
	sr, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
	if !plan.VolumeName.IsNull() {
		snapResponse, err2 := r.client.GetVolume("", "", "", plan.VolumeName.ValueString(), false)
		if err2 != nil {
			resp.Diagnostics.AddError(
				"Error getting volume",
				"Could not get volume, unexpected error: "+err2.Error(),
			)
			return
		}
		plan.VolumeID = types.StringValue(snapResponse[0].ID)
	}
	sdcList := []SdcList{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcList, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range sdcList {
		if si.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", si.SdcName)
			if errA != nil {
				errMsg["sdc_map_err_"+si.SdcName] = errA.Error()
				continue
			} else {
				si.SdcID = foundsdc.Sdc.ID
			}
		}
	}
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SdcInfoAttrTypes,
	}
	objectSdcInfos := []attr.Value{}
	for _, si := range sdcList {
		// refreshing state for drift outside terraform
		if si.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", si.SdcName)
			if errA != nil {
				errMsg["sdc_map_err_"+si.SdcName] = errA.Error()
			} else {
				si.SdcID = foundsdc.Sdc.ID
			}
		}
		obj := map[string]attr.Value{
			"sdc_id":           types.StringValue(si.SdcID),
			"limit_iops":       types.Int64Value(int64(si.LimitIops)),
			"limit_bw_in_mbps": types.Int64Value(int64(si.LimitBwInMbps)),
			"sdc_name":         types.StringValue(si.SdcName),
			"access_mode":      types.StringValue(si.AccessMode),
		}
		objVal, _ := types.ObjectValue(SdcInfoAttrTypes, obj)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	mappedSdcInfoVal, _ := types.SetValue(sdcInfoElemType, objectSdcInfos)
	plan.SdcList = mappedSdcInfoVal
	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
