package powerflex

import (
	"context"
	"strconv"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
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
	if plan.VolumeName.ValueString() != "" {
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
	snapReq := &pftypes.SnapshotDef{
		VolumeID:     plan.VolumeID.ValueString(),
		SnapshotName: plan.Name.ValueString(),
	}
	snapshotReqs = append(snapshotReqs, snapReq)
	snapParam := &pftypes.SnapshotVolumesParam{
		SnapshotDefs:         snapshotReqs,
		RetentionPeriodInMin: convertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString()),
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
	snapID := snapResps.VolumeIDList[0]
	snapResponse, err2 := r.client.GetVolume("", snapID, "", "", false)
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
	// setting access mode limit on snapshot. default value - ReadOnly
	err = snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting Volume Access Mode",
			"Could not set snapshots, unexpected err: "+err.Error(),
		)
	}

	// if size of volume is defined, then checking if size of snapshot is defined greater than volume then expand the snapshot size otherwise throw error.
	if !plan.Size.IsNull() {
		vikb, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		tflog.Info(ctx, "vikb"+strconv.FormatInt(vikb, 10))
		if int64(snapResource.Volume.SizeInKb) != vikb {
			err3 := snapResource.SetVolumeSize(strconv.FormatInt(plan.Size.ValueInt64(), 10))
			if err3 != nil {
				resp.Diagnostics.AddError(
					"Error setting snapshot size",
					"Could not set snapshot size, unexpected err: "+err3.Error(),
				)
				snapResource.RemoveVolume("")
				return
			}
		}
	}

	// locking the auto snapshot on finding LockedAutoSnapshot parameter as true
	if plan.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
			snapResource.RemoveVolume("")
			return
		}
	}
	successMapped := make([]string, 0)
	sdcList := []SdcList{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcList, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range sdcList {
		if si.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", si.SdcName)
			si.SdcID = foundsdc.Sdc.ID
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error Finding SDC with name",
					"Could not get sdc with name"+si.SdcName+",unexpected error: "+errA.Error(),
				)
				for _, usi := range successMapped {
					snapResource.UnmapVolumeSdc(
						&pftypes.UnmapVolumeSdcParam{
							SdcID: usi,
						},
					)
				}
				snapResource.RemoveVolume("")
				return
			}
		}
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 si.SdcID,
			AllowMultipleMappings: "true",
		}
		// mapping the snapshot to sdc
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Snapshot to SDCs",
				"Could not map Snapshot to scs with id: "+si.SdcID+", unexpected error: "+err3.Error(),
			)
			for _, usi := range successMapped {
				snapResource.UnmapVolumeSdc(
					&pftypes.UnmapVolumeSdcParam{
						SdcID: usi,
					},
				)
			}
			snapResource.RemoveVolume("")
			return
		}
		// setting limits on mapped sdc
		smslp := pftypes.SetMappedSdcLimitsParam{
			SdcID:                si.SdcID,
			BandwidthLimitInKbps: strconv.FormatInt(int64(si.LimitBwInMbps*1024), 10),
			IopsLimit:            strconv.FormatInt(int64(si.LimitIops), 10),
		}
		err4 := snapResource.SetMappedSdcLimits(&smslp)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Error Setting Mapped Sdc Limits",
				"Could not set mapped sdc limit, unexpected error: "+err4.Error(),
			)
			for _, usi := range successMapped {
				snapResource.UnmapVolumeSdc(
					&pftypes.UnmapVolumeSdcParam{
						SdcID: usi,
					},
				)
			}
			snapResource.RemoveVolume("")
			return
		}
		err5 := snapResource.SetVolumeMappingAccessMode(si.AccessMode, si.SdcID)
		if err5 != nil {
			resp.Diagnostics.AddError(
				"Error Setting Access Mode On Mapped SDC To Snapshot",
				"Could not set access mode on mapped sdc, unexpected error: "+err5.Error(),
			)
			for _, usi := range successMapped {
				snapResource.UnmapVolumeSdc(
					&pftypes.UnmapVolumeSdcParam{
						SdcID: usi,
					},
				)
			}
			snapResource.RemoveVolume("")
			return
		}
		successMapped = append(successMapped, si.SdcID)
	}
	snapResponse, err2 = r.client.GetVolume("", snapID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap = snapResponse[0]
	state := SnapshotTerraformState(snap, plan, sdcList)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	sr, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	sdcList := []SdcList{}
	diags = state.SdcList.ElementsAs(ctx, &sdcList, true)
	for _, sl := range sdcList {
		if sl.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", sl.SdcName)
			sl.SdcID = foundsdc.Sdc.ID
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error Finding SDC with name",
					"Could not get sdc with name"+sl.SdcName+",unexpected error: "+errA.Error(),
				)
				return
			}
		}
	}
	resp.Diagnostics.Append(diags...)
	snap := snapResponse[0]
	state = SnapshotTerraformState(snap, state, sdcList)
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
	sr, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
	VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	plan.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
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
		errRename := snapResource.SetVolumeName(plan.Name.ValueString())
		if errRename != nil {
			resp.Diagnostics.AddError(
				"Error renaming the snapshot -> "+plan.Name.ValueString()+" : "+state.Name.ValueString(),
				"Could not rename the snapshot, unexpected error:"+errRename.Error(),
			)
			return
		}
	}
	// updating the size of the volume if there is change in plan
	if plan.VolumeSizeInKb.ValueString() != state.VolumeSizeInKb.ValueString() {
		sizeInGb, _ := strconv.Atoi(strconv.FormatInt(VSIKB, 10))
		sizeInGb = sizeInGb / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		err3 := snapResource.SetVolumeSize(sizeInGB)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error setting snapshot size -> "+plan.VolumeSizeInKb.ValueString()+":"+state.VolumeSizeInKb.ValueString(),
				"Could not set new snapshot size -> "+sizeInGB+", unexpected err: "+err3.Error(),
			)
			return
		}
	}

	// locking the snapshot in case of change in LockedAutoSnapshot state to true
	if plan.LockAutoSnapshot.ValueBool() && !state.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// unlocking the snapshot in case of change in LockedAutoSnapshot state to false
	if !plan.LockAutoSnapshot.ValueBool() && state.LockAutoSnapshot.ValueBool() {
		err := snapResource.UnlockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unlocking Auto Snapshots",
				"Could not unlock auto snapshots, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if !plan.DesiredRetention.IsNull() && (plan.RetentionUnit.ValueString() != state.RetentionUnit.String()) || (plan.DesiredRetention.ValueInt64() != state.DesiredRetention.ValueInt64()) {
		retentionInMin := convertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString())
		err := snapResource.SetSnapshotSecurity(retentionInMin)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Snapshots Security",
				"Could not set snapshot security, unexpected error: "+err.Error(),
			)
			return
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
	for idx, psl := range planSdcList {
		if psl.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", psl.SdcName)
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error Finding SDC with name",
					"Could not get sdc with name"+psl.SdcName+",unexpected error: "+errA.Error(),
				)
				return
			}
			planSdcList[idx].SdcID = foundsdc.Sdc.ID
			planSdcIds = append(planSdcIds, planSdcList[idx].SdcID)
		} else {
			planSdcIds = append(planSdcIds, psl.SdcID)
		}

	}
	for idx, ssl := range stateSdcList {
		if ssl.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", ssl.SdcName)
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error Finding SDC with name",
					"Could not get sdc with name"+ssl.SdcName+",unexpected error: "+errA.Error(),
				)
				return
			}
			stateSdcList[idx].SdcID = foundsdc.Sdc.ID
			stateSdcIds = append(stateSdcIds, stateSdcList[idx].SdcID)
		} else {
			stateSdcIds = append(stateSdcIds, ssl.SdcID)
		}

	}
	mapSdcIds := Difference(planSdcIds, stateSdcIds)
	unmapSdcIds := Difference(stateSdcIds, planSdcIds)
	nonchangeSdcIds := Difference(planSdcIds, mapSdcIds)

	// changing the access mode in case of change in access mode state
	if (plan.AccessMode.ValueString() == "ReadWrite") && (state.AccessMode.ValueString() == "ReadOnly") {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Snapshot Access Mode",
				"Could not set the Snapshot Access Mode, unexpected err: "+err.Error(),
			)
			return
		}
	}

	for _, msi := range mapSdcIds {

		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msi,
			AllowMultipleMappings: "true",
		}
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Snapshots to SDCs",
				"Could not map snapshot to scs with id: "+msi+", unexpected error: "+err3.Error(),
			)
			return
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
					resp.Diagnostics.AddError(
						"Error Setting Mapped Sdc Limits",
						"Could not set mapped sdc limit, unexpected error: "+err4.Error(),
					)
					return
				}
				err5 := snapResource.SetVolumeMappingAccessMode(ssl.AccessMode, ssl.SdcID)
				if err5 != nil {
					resp.Diagnostics.AddError(
						"Error Setting Access Mode On Mapped SDC To Snapshot",
						"Could not set access mode on mapped sdc, unexpected error: "+err5.Error(),
					)
					return
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
			resp.Diagnostics.AddError(
				"Error Unmapping Snapshot to SDCs",
				"Could not Unmap snapshot to scs with id: "+usi+", unexpected error: "+err4.Error(),
			)
			return
		}
	}

	// changing the access mode in case of change in access mode state
	if (plan.AccessMode.ValueString() == "ReadOnly") && (state.AccessMode.ValueString() == "ReadWrite") {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Snapshot Access Mode",
				"Could not set the Snapshot Access Mode, unexpected err: "+err.Error(),
			)
			return
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
				resp.Diagnostics.AddError(
					"Error Setting Mapped Sdc Limits",
					"Could not set mapped sdc limit, unexpected error: "+err4.Error(),
				)
				return
			}
		}

		if planObj.AccessMode != stateObj.AccessMode {
			err5 := snapResource.SetVolumeMappingAccessMode(planObj.AccessMode, planObj.SdcID)
			if err5 != nil {
				resp.Diagnostics.AddError(
					"Error Setting Access Mode On Mapped SDC To Snapshot",
					"Could not set access mode on mapped sdc, unexpected error: "+err5.Error(),
				)
				return
			}
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
	state = SnapshotTerraformState(snap, plan, planSdcList)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	sr, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
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
	for _, ssl := range sdcsToUnmap {
		if ssl.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", ssl.SdcName)
			ssl.SdcID = foundsdc.Sdc.ID
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error Finding SDC with name",
					"Could not get sdc with name"+ssl.SdcName+",unexpected error: "+errA.Error(),
				)
				return
			}
		}
	}
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
	err = snapshot.RemoveVolume(state.RemoveMode.ValueString())
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

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
