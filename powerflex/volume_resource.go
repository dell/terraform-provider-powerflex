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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewVolumeResource is a helper function to simplify the provider implementation.
func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

// volumeResource is the resource implementation.
type volumeResource struct {
	client *goscaleio.Client
}

// Metadata returns the resource type name.
func (r *volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema defines the schema for the resource.
func (r *volumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VolumeResourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *volumeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

func (r *volumeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan VolumeResourceModel
	var state VolumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	sr, err := getFirstSystem(r.client)
	pdr := goscaleio.NewProtectionDomain(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"Could not get first system, unexpected error: "+err.Error(),
		)
		return
	}
	if !plan.ProtectionDomainName.IsNull() {
		protectionDomain, err := sr.FindProtectionDomain("", plan.ProtectionDomainName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting protection domain",
				"Could not get protection domain with name: "+plan.ProtectionDomainName.String()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		pdr.ProtectionDomain = protectionDomain
		plan.ProtectionDomainID = types.StringValue(protectionDomain.ID)
	}
	if !plan.StoragePoolName.IsNull() {
		storagePool, err := pdr.FindStoragePool("", plan.StoragePoolName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting storage pool",
				"Could not get storage pool with name: "+plan.StoragePoolName.String()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		plan.StoragePoolID = types.StringValue(storagePool.ID)
	}

	if !plan.Size.IsNull() {
		if plan.Size.ValueInt64()%8 != 0 {
			resp.Diagnostics.AddError(
				"Error: Size Must be in granularity of 8GB",
				"Could not assign volume with size. sizeInGb ("+strconv.FormatInt(plan.Size.ValueInt64(), 10)+") must be a positive number in granularity of 8 GB.",
			)
			return
		}
		VSIKB := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		plan.SizeInKb = types.Int64Value(int64(VSIKB))
	}

	sdcList := []SDCItemize{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcList, true)
	resp.Diagnostics.Append(diags...)

	if len(sdcList) > 0 && plan.AccessMode.ValueString() == "ReadOnly" {
		resp.Diagnostics.AddError(
			"Error: SDC can't be mapped, planned access_mode sets to 'ReadOnly`",
			"Could not map sdc to volume with ReadOnly access mode",
		)
		return
	}
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SdcInfoAttrTypes,
	}
	objectSdcInfos := []attr.Value{}
	for _, si := range sdcList {
		if si.SdcID == "" {
			foundsdc, errA := sr.FindSdc("Name", si.SdcName)
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error getting sdc with the name: "+si.SdcName,
					"Couldn't get the sdc, unexpected error: "+errA.Error(),
				)
				return
			}
			si.SdcID = foundsdc.Sdc.ID

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

	// update scenario
	if !req.State.Raw.IsNull() {
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if plan.SizeInKb.ValueInt64() < state.SizeInKb.ValueInt64() {
			resp.Diagnostics.AddError(
				"Error: Volume capacity can only be increased",
				"volume size with "+state.Size.String()+" size couldn't be decreased to "+plan.Size.String()+" size",
			)
			return
		}
		// volume type can't be updated after volume creation
		if plan.VolumeType.ValueString() != state.VolumeType.ValueString() {
			resp.Diagnostics.AddError(
				"Error: Volume Type can't be updated",
				state.VolumeType.ValueString()+" volume type couldn't be converted to "+plan.VolumeType.ValueString()+" volume type",
			)
			return
		}
		// volume rmcahce can't be update after volume creation
		if plan.UseRmCache.ValueBool() != state.UseRmCache.ValueBool() {
			resp.Diagnostics.AddError(
				"Error: Volume RMCache can't be updated",
				"volume RMCache sets to "+state.UseRmCache.String()+" can't be change to "+plan.UseRmCache.String(),
			)
			return
		}
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *volumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan VolumeResourceModel
	errMsg := make(map[string]string, 0)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	volumeCreate := &pftypes.VolumeParam{
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
		StoragePoolID:      plan.StoragePoolID.ValueString(),
		UseRmCache:         strconv.FormatBool(plan.UseRmCache.ValueBool()),
		VolumeType:         plan.VolumeType.ValueString(),
		VolumeSizeInKb:     strconv.FormatInt(plan.SizeInKb.ValueInt64(), 10),
		Name:               plan.Name.ValueString(),
	}
	spr, _ := getStoragePoolInstance(r.client, volumeCreate.StoragePoolID, volumeCreate.ProtectionDomainID)
	volCreateResponse, err1 := spr.CreateVolume(volumeCreate)
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err1.Error(),
		)
		return
	}
	volsResponse, err2 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after creation",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	vr := goscaleio.NewVolume(r.client)
	vr.Volume = vol
	if !plan.AccessMode.IsNull() {
		err := vr.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			errMsg["access_mode"] = err.Error()
		}
	}
	sdcItems := []SDCItemize{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcItems, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range sdcItems {
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 si.SdcID,
			AllowMultipleMappings: "true",
		}
		// mapping the snapshot to sdc
		err3 := vr.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			errMsg["sdc_map_err_"+si.SdcID] += "\n" + err3.Error()
		}
		// setting limits on mapped sdc
		smslp := pftypes.SetMappedSdcLimitsParam{
			SdcID:                si.SdcID,
			BandwidthLimitInKbps: strconv.FormatInt(int64(si.LimitBwInMbps*1024), 10),
			IopsLimit:            strconv.FormatInt(int64(si.LimitIops), 10),
		}
		err4 := vr.SetMappedSdcLimits(&smslp)
		if err4 != nil {
			errMsg["sdc_map_err_"+si.SdcID] += "\n" + err4.Error()
		}
		err5 := vr.SetVolumeMappingAccessMode(si.AccessMode, si.SdcID)
		if err5 != nil {
			errMsg["sdc_map_err_"+si.SdcID] += "\n" + err5.Error()
		}
	}
	volsResponse, err2 = spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after mapping the sdcs",
			"Could not get volume after mapping the sdcs, unexpected error: "+err2.Error(),
		)
		return
	}
	vol = volsResponse[0]
	refreshVolumeState(vol, &plan)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if len(errMsg) > 0 {
		failureAction := ""
		failureMessage := ""
		for key, value := range errMsg {
			failureAction += key + " "
			failureMessage += key + " : " + value + "\n"
		}
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Failed Actions [%v]", failureAction),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *volumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state VolumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsResponse, err2 := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	refreshVolumeState(vol, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan VolumeResourceModel
	errMsg := make(map[string]string, 0)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get current state
	var state VolumeResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsplan, _ := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]
	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		errRename := volresource.SetVolumeName(plan.Name.ValueString())
		if errRename != nil {
			errMsg["name"] = "Error renaming the volume -> " + plan.Name.ValueString() + " : " + state.Name.ValueString() + " \n Could not rename the volume, unexpected error:" + errRename.Error()
		}
	}

	// updating the size of the volume if there is change in plan
	if plan.SizeInKb.ValueInt64() != state.SizeInKb.ValueInt64() {
		sizeInGb := plan.SizeInKb.ValueInt64() / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		// sizeInGb = ((sizeInGb / 8) + 1) * 8
		// newSizeIn8Gb := strconv.FormatInt(int64(sizeInGb), 10)
		err3 := volresource.SetVolumeSize(sizeInGB)
		if err3 != nil {
			errMsg["size"] = "Error setting volume size -> " + plan.SizeInKb.String() + ":" + state.SizeInKb.String() + "\nCould not set new volume size -> " + sizeInGB + ", unexpected err: " + err3.Error()
		}
	}
	planSdcList := []SDCItemize{}
	stateSdcList := []SDCItemize{}
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

	if !plan.AccessMode.IsNull() {
		err := volresource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			errMsg["access_mode"] = err.Error()
		}
	}
	for _, msi := range mapSdcIds {
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msi,
			AllowMultipleMappings: "true",
		}
		err3 := volresource.MapVolumeSdc(&pfmvsp)
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
				err4 := volresource.SetMappedSdcLimits(&smslp)
				if err4 != nil {
					errMsg["sdc_map_err_"+msi] += "\n" + err4.Error()
				}
				err5 := volresource.SetVolumeMappingAccessMode(ssl.AccessMode, ssl.SdcID)
				if err5 != nil {
					errMsg["sdc_map_err_"+msi] += "\n" + err5.Error()
				}
			}
		}
	}

	for _, usi := range unmapSdcIds {
		err4 := volresource.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: usi,
			},
		)
		if err4 != nil {
			errMsg["sdc_unmap_err_"+usi] += "\n" + err4.Error()
		}
	}

	for _, ncsi := range nonchangeSdcIds {
		var planObj SDCItemize
		var stateObj SDCItemize

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
			err4 := volresource.SetMappedSdcLimits(&smslp)
			if err4 != nil {
				errMsg["sdc_map_err_"+planObj.SdcID] += "\n" + err4.Error()
			}
		}

		if planObj.AccessMode != stateObj.AccessMode {
			err5 := volresource.SetVolumeMappingAccessMode(planObj.AccessMode, planObj.SdcID)
			if err5 != nil {
				errMsg["sdc_map_err_"+planObj.SdcID] += "\n" + err5.Error()
			}
		}

	}
	vols, err2 := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := vols[0]
	refreshVolumeState(vol, &plan)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if len(errMsg) > 0 {
		failureAction := ""
		failureMessage := ""
		for key, value := range errMsg {
			failureAction += key + " "
			failureMessage += key + " : " + value + "\n"
		}
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Failed Actions [%v]", failureAction),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *volumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state VolumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsplan, _ := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]
	sdcsToUnmap := []SDCItemize{}
	diags = state.SdcList.ElementsAs(ctx, &sdcsToUnmap, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	for _, stu := range sdcsToUnmap {
		err := volresource.UnmapVolumeSdc(
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
	err := volresource.RemoveVolume(state.RemoveMode.ValueString())
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

func (r *volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
