package powerflex

import (
	"context"
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

// ModifyPlan modify resource plan attribute value
func (r *volumeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan VolumeResourceModel
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
	if !plan.ProtectionDomainName.IsUnknown() {
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
	} else {
		protectionDomain, err := sr.FindProtectionDomain(plan.ProtectionDomainID.ValueString(), "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting protection domain with id",
				"Could not get protection domain with id: "+plan.ProtectionDomainName.ValueString()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		pdr.ProtectionDomain = protectionDomain
		plan.ProtectionDomainName = types.StringValue(protectionDomain.Name)
	}
	if !plan.StoragePoolName.IsUnknown() {
		storagePool, err := pdr.FindStoragePool("", plan.StoragePoolName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting storage pool",
				"Could not get storage pool with name: "+plan.StoragePoolName.ValueString()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		plan.StoragePoolID = types.StringValue(storagePool.ID)
	} else {
		storagePool, err := pdr.FindStoragePool(plan.StoragePoolID.ValueString(), "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting storage pool with id",
				"Could not get storage pool with with id: "+plan.StoragePoolID.ValueString()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		plan.StoragePoolName = types.StringValue(storagePool.Name)
	}

	if !plan.Size.IsNull() {
		VSIKB := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		plan.SizeInKb = types.Int64Value(VSIKB)
	}
	sdcList := []SDCItemize{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcList, true)
	resp.Diagnostics.Append(diags...)
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
					"unexpected error: "+errA.Error(),
				)
				return
			}
			si.SdcID = foundsdc.Sdc.ID
		}
		if si.SdcName == "" {
			foundsdc, errA := sr.FindSdc("ID", si.SdcID)
			if errA != nil {
				resp.Diagnostics.AddError(
					"Error getting sdc name from sdc id: "+si.SdcID,
					"unexpected error: "+errA.Error(),
				)
				return
			}
			si.SdcName = foundsdc.Sdc.Name
		}
		if si.LimitIops <= 10 && si.LimitIops > 0 {
			resp.Diagnostics.AddError(
				"Error setting the limit iops",
				"sdc  "+si.SdcID+" "+si.SdcName+" limit iops must be a number larger than 10 or 0.",
			)
			return
		}
		obj := map[string]attr.Value{
			"sdc_id":           types.StringValue(si.SdcID),
			"limit_iops":       types.Int64Value(int64(si.LimitIops)),
			"limit_bw_in_mbps": types.Int64Value(int64(si.LimitBwInMbps)),
			"sdc_name":         types.StringValue(si.SdcName),
			"access_mode":      types.StringValue(si.AccessMode),
		}
		objVal, dgs := types.ObjectValue(SdcInfoAttrTypes, obj)
		diags = append(diags, dgs...)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	mappedSdcInfoVal, dgs := types.SetValue(sdcInfoElemType, objectSdcInfos)
	diags = append(diags, dgs...)
	resp.Diagnostics.Append(diags...)
	plan.SdcList = mappedSdcInfoVal
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
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	volumeCreate := &pftypes.VolumeParam{
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
		StoragePoolID:      plan.StoragePoolID.ValueString(),
		UseRmCache:         strconv.FormatBool(plan.UseRmCache.ValueBool()),
		CompressionMethod:  plan.CompressionMethod.ValueString(),
		VolumeType:         plan.VolumeType.ValueString(),
		VolumeSizeInKb:     strconv.FormatInt(plan.SizeInKb.ValueInt64(), 10),
		Name:               plan.Name.ValueString(),
	}
	spr, err0 := getStoragePoolInstance(r.client, volumeCreate.StoragePoolID, volumeCreate.ProtectionDomainID)
	if err0 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool with id: "+volumeCreate.StoragePoolID+" or protection pool with id: "+volumeCreate.ProtectionDomainID,
			"unexpected error: "+err0.Error(),
		)
		return
	}
	// platform fails silently for compression method "None".
	if (spr.StoragePool.DataLayout != "FineGranularity") && (plan.CompressionMethod.ValueString() != "") {
		resp.Diagnostics.AddError(
			"error setting the compression method",
			"compression may only be set on volumes with Fine Granularity layout on storage pool. This storage pool has "+spr.StoragePool.DataLayout+" layout.",
		)
		return
	}
	volCreateResponse, err1 := spr.CreateVolume(volumeCreate)
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			"unexpected error: "+err1.Error(),
		)
		return
	}
	volsResponse, err2 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after creation",
			"unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	vr := goscaleio.NewVolume(r.client)
	vr.Volume = vol
	if !plan.AccessMode.IsNull() {
		err3 := vr.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error setting access mode on volume",
				"unexpected error: "+err3.Error(),
			)
		}
	}
	sdcItems := []SDCItemize{}
	diags = plan.SdcList.ElementsAs(ctx, &sdcItems, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range sdcItems {
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 si.SdcID,
			AccessMode:            si.AccessMode,
			AllowMultipleMappings: "true",
		}
		// mapping the snapshot to sdc
		err4 := vr.MapVolumeSdc(&pfmvsp)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Error mapping sdc: "+si.SdcID+" "+si.SdcName,
				"unexpected error: "+err4.Error(),
			)
		} else {
			// setting limits on mapped sdc
			smslp := pftypes.SetMappedSdcLimitsParam{
				SdcID:                si.SdcID,
				BandwidthLimitInKbps: strconv.FormatInt(int64(si.LimitBwInMbps*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(si.LimitIops), 10),
			}
			err5 := vr.SetMappedSdcLimits(&smslp)
			if err4 != nil {
				resp.Diagnostics.AddError(
					"Error Setting Limits to mapped sdc: "+si.SdcID+" "+si.SdcName,
					"unexpected error: "+err5.Error(),
				)
			}
		}
	}
	volsResponse, err7 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err7 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after mapping the sdcs",
			"Could not get volume after mapping the sdcs, unexpected error: "+err2.Error(),
		)
		return
	}
	vol = volsResponse[0]
	dgs := refreshVolumeState(vol, &plan)
	resp.Diagnostics.Append(dgs...)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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
	volsResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	dgs := refreshVolumeState(vol, &state)
	resp.Diagnostics.Append(dgs...)
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
	volsplan, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"unexpected error: "+err2.Error(),
		)
		return
	}
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]
	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		err3 := volresource.SetVolumeName(plan.Name.ValueString())
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error renaming the volume",
				"unexpected error: "+err3.Error(),
			)
		}
	}

	// updating the size of the volume if there is change in plan
	if plan.SizeInKb.ValueInt64() != state.SizeInKb.ValueInt64() {
		sizeInGb := plan.SizeInKb.ValueInt64() / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		err4 := volresource.SetVolumeSize(sizeInGB)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Error setting the volume size",
				"unexpected error: "+err4.Error(),
			)
		}
	}

	// prompt error on change in volume type, as we can't update the volume type after the creation
	if !plan.VolumeType.Equal(state.VolumeType) {
		resp.Diagnostics.AddError(
			"volume type cannot be update after volume creation.",
			"unexpected error: volume type change is not supported",
		)
	}

	// updating the use rm cache if there is change in plan
	if !plan.UseRmCache.Equal(state.UseRmCache) {
		err5 := volresource.SetVolumeUseRmCache(plan.UseRmCache.ValueBool())
		if err5 != nil {
			resp.Diagnostics.AddError(
				"Error setting the use rm cache",
				"unexpected error: "+err5.Error(),
			)
		}
	}

	// updating the compression if there is change in plan
	if !plan.CompressionMethod.IsUnknown() && !plan.CompressionMethod.Equal(state.CompressionMethod) {
		err6 := volresource.SetCompressionMethod(plan.CompressionMethod.ValueString())
		if err6 != nil {
			resp.Diagnostics.AddError(
				"Error setting the compression method",
				"unexpected error: "+err6.Error(),
			)
		}
	}

	planSdcList := []SDCItemize{}
	stateSdcList := []SDCItemize{}

	//unmarshall the tfsdk sdclist type to go sdclist type
	diags = plan.SdcList.ElementsAs(ctx, &planSdcList, true)
	resp.Diagnostics.Append(diags...)
	diags = state.SdcList.ElementsAs(ctx, &stateSdcList, true)
	resp.Diagnostics.Append(diags...)

	// getting all the sdc_ids from plan and state sdc list and saving into respective list for finding sdc to map and unmap.
	planSdcIds := []string{}
	stateSdcIds := []string{}
	for _, psl := range planSdcList {
		planSdcIds = append(planSdcIds, psl.SdcID)
	}
	for _, ssl := range stateSdcList {
		stateSdcIds = append(stateSdcIds, ssl.SdcID)
	}

	// mapSdcIds will be storing the sdc id for which mapping action need to perform.
	mapSdcIds := Difference(planSdcIds, stateSdcIds)

	// unmapSdcIds will be storing the sdc id for which unmapping action need to perform.
	unmapSdcIds := Difference(stateSdcIds, planSdcIds)

	// nonchangeSdcIds will be storing the sdc id for which mapping parameter change action need to perform.
	nonchangeSdcIds := Difference(planSdcIds, mapSdcIds)

	// changing the access mode in case of change in access mode state from readonly to readwrite
	if (plan.AccessMode.ValueString() == READWRITE) && (state.AccessMode.ValueString() == READONLY) {
		err7 := volresource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err7 != nil {
			resp.Diagnostics.AddError(
				"Error setting the access mode",
				"unexpected error: "+err7.Error(),
			)
		}
	}

	// retervial of sdc id from mapping list and performing map operation.
	for _, ssl := range planSdcList {
		for _, msi := range mapSdcIds {
			if ssl.SdcID == msi {
				pfmvsp := pftypes.MapVolumeSdcParam{
					SdcID:                 ssl.SdcID,
					AccessMode:            ssl.AccessMode,
					AllowMultipleMappings: "true",
				}
				err8 := volresource.MapVolumeSdc(&pfmvsp)
				if err8 != nil {
					resp.Diagnostics.AddError(
						"Error mapping volume to sdc: "+msi,
						"unexpected error: "+err8.Error(),
					)
				} else {
					smslp := pftypes.SetMappedSdcLimitsParam{
						SdcID:                ssl.SdcID,
						BandwidthLimitInKbps: strconv.FormatInt(int64(ssl.LimitBwInMbps*1024), 10),
						IopsLimit:            strconv.FormatInt(int64(ssl.LimitIops), 10),
					}
					err9 := volresource.SetMappedSdcLimits(&smslp)
					if err9 != nil {
						resp.Diagnostics.AddError(
							"Error setting limits to sdc: "+ssl.SdcID+" "+ssl.SdcName,
							"unexpected error: "+err9.Error(),
						)
					}
				}
			}
		}
	}

	// reterival of sdc id from unmapping list and performing unmap operation.
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

	// reterival of sdc id from non changed sdc list and performing change operation.
	for _, ncsi := range nonchangeSdcIds {
		var planObj SDCItemize
		var stateObj SDCItemize

		// getting the plan sdc obj for comparision with state
		for _, psl := range planSdcList {
			if ncsi == psl.SdcID {
				planObj = psl
				break
			}
		}

		// getting the state sdc obj for comparision with plan
		for _, ssl := range stateSdcList {
			if ncsi == ssl.SdcID {
				stateObj = ssl
				break
			}
		}

		// updating the sdc mapping parameters: limit iops and bandwidth limits if there is change in plan and state
		if (planObj.LimitIops != stateObj.LimitIops) || (planObj.LimitBwInMbps != stateObj.LimitBwInMbps) {
			smslp := pftypes.SetMappedSdcLimitsParam{
				SdcID:                planObj.SdcID,
				BandwidthLimitInKbps: strconv.FormatInt(int64(planObj.LimitBwInMbps*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(planObj.LimitIops), 10),
			}
			err11 := volresource.SetMappedSdcLimits(&smslp)
			if err11 != nil {
				resp.Diagnostics.AddError(
					"Error setting access mode to sdc: "+planObj.SdcID+" "+planObj.SdcName,
					"unexpected error: "+err11.Error(),
				)
			}
		}

		// updating the access mode for sdc mapping if there is change in plan and state
		if planObj.AccessMode != stateObj.AccessMode {
			err12 := volresource.SetVolumeMappingAccessMode(planObj.AccessMode, planObj.SdcID)
			if err12 != nil {
				resp.Diagnostics.AddError(
					"Error setting access mode to sdc: "+planObj.SdcID+" "+planObj.SdcName,
					"unexpected error: "+err12.Error(),
				)
			}
		}

	}

	// changing the access mode in case of change in access mode state from readwrite to readonly
	if (plan.AccessMode.ValueString() == READONLY) && (state.AccessMode.ValueString() == READWRITE) {
		err13 := volresource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err13 != nil {
			resp.Diagnostics.AddError(
				"Error setting the access mode",
				"unexpected error: "+err13.Error(),
			)
		}
	}

	vols, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := vols[0]
	dgs := refreshVolumeState(vol, &plan)
	resp.Diagnostics.Append(dgs...)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

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
	volsplan, err1 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err1.Error(),
		)
		return
	}
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]

	sdcsToUnmap := []SDCItemize{}
	//unmarshall the tfsdk sdclist type to go sdclist type
	diags = state.SdcList.ElementsAs(ctx, &sdcsToUnmap, true)
	resp.Diagnostics.Append(diags...)

	// reterival of sdc id from sdcsToUnmap slice of struct SDCItemize and performing unmap operation.
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

	// finally removing the volume after unmap operation
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
