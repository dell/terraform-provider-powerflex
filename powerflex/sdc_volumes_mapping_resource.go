package powerflex

import (
	"context"
	"strconv"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sdcVolumeMappingResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	VolumeList types.Set    `tfsdk:"volume_list"`
}

type sdcVolumeModel struct {
	VolumeID   types.String `tfsdk:"volume_id"`
	VolumeName types.String `tfsdk:"volume_name"`
	IOPSLimit  types.Int64  `tfsdk:"limit_iops"`
	BWLimit    types.Int64  `tfsdk:"limit_bw_in_mbps"`
	AccessMode types.String `tfsdk:"access_mode"`
}

// NewSDSResource is a helper function to simplify the provider implementation.
func NewSDCVolumesMappingResource() resource.Resource {
	return &sdcVolumeMappingResource{}
}

// sdsVolumeMappingResource is the resource implementation.
type sdcVolumeMappingResource struct {
	client *goscaleio.Client
}

func (r *sdcVolumeMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_volumes_mapping"
}

func (r *sdcVolumeMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage SDC-Volume mapping resource.",
		MarkdownDescription: "Manage SDC-Volume mapping resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the SDC.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the SDC.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRoot("id")),
				},
			},
			"volume_list": schema.SetNestedAttribute{
				Description:         "List of volumes mapped to SDC.",
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "List of volumes mapped to SDC.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"volume_id": schema.StringAttribute{
							Description:         "The ID of the volume.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The ID of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"volume_name": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("volume_id")),
							},
						},
						"limit_iops": schema.Int64Attribute{
							Description:         "IOPS limit",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "IOPS limit",
						},
						"limit_bw_in_mbps": schema.Int64Attribute{
							Description:         "Bandwidth limit in MBPS",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Bandwidth limit in MBPS",
						},
						"access_mode": schema.StringAttribute{
							Description:         "The Access Mode of the SDC",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The Access Mode of the SDC",
							Validators: []validator.String{stringvalidator.OneOf(
								"ReadOnly",
								"ReadWrite",
								"NoAccess",
							)},
							PlanModifiers: []planmodifier.String{
								stringDefault("ReadOnly"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *sdcVolumeMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// ModifyPlan modify resource plan attribute value
func (r *sdcVolumeMappingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan sdcVolumeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Get the system on the PowerFlex cluster
	system, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	var sdc *goscaleio.Sdc

	// Populate SDC name in the plan if ID is provided in the config
	if !plan.ID.IsUnknown() {
		sdc, err = system.GetSdcByID(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting SDC with ID",
				"Could not get SDC with ID: "+plan.ID.ValueString()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		plan.Name = types.StringValue(sdc.Sdc.Name)
	} else if !plan.Name.IsUnknown() {
		sdc, err = system.FindSdc("Name", plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting SDC with name",
				"Could not get SDC with name: "+plan.ID.ValueString()+", \n unexpected error: "+err.Error(),
			)
			return
		}
		plan.ID = types.StringValue(sdc.Sdc.ID)
	}

	volList := []sdcVolumeModel{}
	diags = plan.VolumeList.ElementsAs(ctx, &volList, true)
	resp.Diagnostics.Append(diags...)

	for index, vol := range volList {
		// Populate volume name in the plan if volume ID is provided in the config
		if !vol.VolumeID.IsUnknown() {
			volume, err := r.client.GetVolume("", vol.VolumeID.ValueString(), "", "", false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume with ID: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			volList[index].VolumeName = types.StringValue(volume[0].Name)
		} else if !vol.VolumeName.IsUnknown() {
			volume, err := r.client.GetVolume("", "", "", vol.VolumeName.ValueString(), false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume with name: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			volList[index].VolumeID = types.StringValue(volume[0].ID)
		}
	}

	// Modify the plan to populate volume list
	mappedVolumeList, dgs := GetVolSetValueFromItems(volList)
	resp.Diagnostics.Append(dgs...)
	plan.VolumeList = mappedVolumeList
	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Return the type for volume list
func GetVolSetValueFromItems(volumes []sdcVolumeModel) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	volInfoElemType := types.ObjectType{
		AttrTypes: GetVolType(),
	}

	if len(volumes) == 0 {
		return types.SetNull(volInfoElemType), diags
	}

	objectVolInfos := []attr.Value{}
	for _, vol := range volumes {
		obj := map[string]attr.Value{
			"volume_id":        vol.VolumeID,
			"volume_name":      vol.VolumeName,
			"limit_iops":       vol.IOPSLimit,
			"limit_bw_in_mbps": vol.BWLimit,
			"access_mode":      vol.AccessMode,
		}
		objVal, dgs := types.ObjectValue(GetVolType(), obj)
		diags = append(diags, dgs...)
		objectVolInfos = append(objectVolInfos, objVal)
	}
	mappedSdcInfoVal, dgs := types.SetValue(volInfoElemType, objectVolInfos)
	diags = append(diags, dgs...)
	return mappedSdcInfoVal, diags
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdcVolumeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sdcVolumeMappingResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volList := []sdcVolumeModel{}
	diags = plan.VolumeList.ElementsAs(ctx, &volList, true)
	resp.Diagnostics.Append(diags...)

	for _, vol := range volList {
		volType, err := getVolumeType(r.client, vol.VolumeID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting volume",
				"unexpected error: "+err.Error(),
			)
			return
		}

		// Add mapped SDC
		mapType := goscaleio_types.MapVolumeSdcParam{
			SdcID:                 plan.ID.ValueString(),
			AccessMode:            vol.AccessMode.ValueString(),
			AllowMultipleMappings: "true",
		}

		err = volType.MapVolumeSdc(&mapType)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error mapping sdc: "+plan.ID.ValueString(),
				"unexpected error: "+err.Error(),
			)
		} else {
			// setting limits on mapped sdc
			limitType := goscaleio_types.SetMappedSdcLimitsParam{
				SdcID:                plan.ID.ValueString(),
				BandwidthLimitInKbps: strconv.FormatInt(int64(vol.BWLimit.ValueInt64()*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(vol.IOPSLimit.ValueInt64()), 10),
			}
			err = volType.SetMappedSdcLimits(&limitType)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error setting limits to sdc: "+plan.ID.String(),
					"unexpected error: "+err.Error(),
				)
			}
		}
	}

	sdcType, err1 := getSdcType(r.client, plan.ID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error Getting SDC type: "+plan.ID.String(),
			"unexpected error: "+err1.Error(),
		)
	}

	// Get the volumes mapped to SDC
	mappedVolumes, err2 := sdcType.GetVolume()
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volumes mapped to SDC : "+plan.ID.String(),
			"unexpected error: "+err2.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := updateSDCVolMapState(mappedVolumes, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func GetVolType() map[string]attr.Type {
	return map[string]attr.Type{
		"volume_id":        types.StringType,
		"volume_name":      types.StringType,
		"limit_iops":       types.Int64Type,
		"limit_bw_in_mbps": types.Int64Type,
		"access_mode":      types.StringType,
	}
}

func GetVolValue(vol *goscaleio_types.Volume) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetVolType(), map[string]attr.Value{
		"volume_id":        types.StringValue(vol.ID),
		"volume_name":      types.StringValue(vol.Name),
		"limit_iops":       types.Int64Value(int64(vol.MappedSdcInfo[0].LimitIops)),
		"limit_bw_in_mbps": types.Int64Value(int64(vol.MappedSdcInfo[0].LimitBwInMbps)),
		"access_mode":      types.StringValue(vol.MappedSdcInfo[0].AccessMode),
	})
}

func updateSDCVolMapState(mappedVolumes []*goscaleio_types.Volume, plan sdcVolumeMappingResourceModel) (sdcVolumeMappingResourceModel, diag.Diagnostics) {
	var state sdcVolumeMappingResourceModel
	state = plan
	SDCAttrTypes := GetVolType()

	SDCElemType := types.ObjectType{
		AttrTypes: SDCAttrTypes,
	}

	objectSDCs := []attr.Value{}
	var diags diag.Diagnostics
	for _, vol := range mappedVolumes {
		objVal, dgs := GetVolValue(vol)
		diags = append(diags, dgs...)
		objectSDCs = append(objectSDCs, objVal)
		state.Name = types.StringValue(vol.MappedSdcInfo[0].SdcName)
		state.ID = types.StringValue(vol.MappedSdcInfo[0].SdcID)
	}
	setVal, dgs := types.SetValue(SDCElemType, objectSDCs)
	diags = append(diags, dgs...)
	state.VolumeList = setVal

	return state, diags
}

// Read refreshes the Terraform state with the latest data.
func (r *sdcVolumeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state sdcVolumeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdcType, err1 := getSdcType(r.client, state.ID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error Getting SDC type: "+state.ID.String(),
			"unexpected error: "+err1.Error(),
		)
	}

	mappedVolumes, err2 := sdcType.GetVolume()
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volumes mapped to SDC : "+state.ID.String(),
			"unexpected error: "+err2.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := updateSDCVolMapState(mappedVolumes, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdcVolumeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan sdcVolumeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state sdcVolumeMappingResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planVolList := []sdcVolumeModel{}
	stateVolList := []sdcVolumeModel{}

	// Populate planVolList with volumes defined in plan
	diags = plan.VolumeList.ElementsAs(ctx, &planVolList, true)
	resp.Diagnostics.Append(diags...)

	// Populate stateVolList with volumes stored in state
	diags = state.VolumeList.ElementsAs(ctx, &stateVolList, true)
	resp.Diagnostics.Append(diags...)

	planVolIds := []string{}
	stateVolIds := []string{}

	// Populate planVolIds with the volume IDs defined in plan
	for _, vol := range planVolList {
		planVolIds = append(planVolIds, vol.VolumeID.ValueString())
	}

	// Populate stateVolIds with the volume IDs stored in plan
	for _, vol := range stateVolList {
		stateVolIds = append(stateVolIds, vol.VolumeID.ValueString())
	}

	// mapVolIds will be storing the volume ids for which mapping needs to be performed
	mapVolIds := Difference(planVolIds, stateVolIds)

	// unmapVolIds will be storing the volume id for which unmapping needs to be performed
	unmapVolIds := Difference(stateVolIds, planVolIds)

	// nonchangeVolIds will be storing the volume id for which limits/access mode needs to be performed
	nonchangeVolIds := Difference(planVolIds, mapVolIds)

	// Perform mapping and setting limits operation
	for _, planVol := range planVolList {
		for _, volID := range mapVolIds {
			if planVol.VolumeID.ValueString() == volID {
				volType, err := getVolumeType(r.client, planVol.VolumeID.ValueString())
				if err != nil {
					resp.Diagnostics.AddError(
						"Error getting volume",
						"unexpected error: "+err.Error(),
					)
					return
				}

				pfmvsp := goscaleio_types.MapVolumeSdcParam{
					SdcID:                 plan.ID.ValueString(),
					AccessMode:            planVol.AccessMode.ValueString(),
					AllowMultipleMappings: "true",
				}

				err = volType.MapVolumeSdc(&pfmvsp)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error mapping volume to sdc: "+volID,
						"unexpected error: "+err.Error(),
					)
				} else {
					smslp := goscaleio_types.SetMappedSdcLimitsParam{
						SdcID:                plan.ID.ValueString(),
						BandwidthLimitInKbps: strconv.FormatInt(planVol.BWLimit.ValueInt64()*1024, 10),
						IopsLimit:            strconv.FormatInt(planVol.IOPSLimit.ValueInt64(), 10),
					}
					err := volType.SetMappedSdcLimits(&smslp)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error setting limits to sdc: "+plan.ID.ValueString(),
							"unexpected error: "+err.Error(),
						)
					}
				}
			}
		}
	}

	// Perform unmap operation
	for _, volID := range unmapVolIds {
		volType, err := getVolumeType(r.client, volID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting volume",
				"unexpected error: "+err.Error(),
			)
			return
		}

		err = volType.UnmapVolumeSdc(
			&goscaleio_types.UnmapVolumeSdcParam{
				SdcID: plan.ID.ValueString(),
			},
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error unmapping sdc: "+plan.ID.ValueString(),
				"unexpected error: "+err.Error(),
			)
		}
	}

	// retrieval of volume id from nonchangeVolIds and perform change operation
	for _, volID := range nonchangeVolIds {
		var planObj sdcVolumeModel
		var stateObj sdcVolumeModel

		// getting the plan volume object for comparison with state
		for _, planVol := range planVolList {
			if volID == planVol.VolumeID.ValueString() {
				planObj = planVol
				break
			}
		}

		// getting the state volume object for comparison with plan
		for _, stateVol := range stateVolList {
			if volID == stateVol.VolumeID.ValueString() {
				stateObj = stateVol
				break
			}
		}

		// update the volume mapping parameters: limit iops and bandwidth limits if plan and state differs
		if (planObj.IOPSLimit != stateObj.IOPSLimit) || (planObj.BWLimit != stateObj.BWLimit) {
			smslp := goscaleio_types.SetMappedSdcLimitsParam{
				SdcID:                plan.ID.ValueString(),
				BandwidthLimitInKbps: strconv.FormatInt(int64(planObj.BWLimit.ValueInt64()*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(planObj.IOPSLimit.ValueInt64()), 10),
			}

			volType, err := getVolumeType(r.client, planObj.VolumeID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume",
					"unexpected error: "+err.Error(),
				)
				return
			}

			err11 := volType.SetMappedSdcLimits(&smslp)
			if err11 != nil {
				resp.Diagnostics.AddError(
					"Error setting limits to sdc: "+plan.ID.ValueString(),
					"unexpected error: "+err11.Error(),
				)
			}
		}

		// update the access mode for volume mapping if plan and state differs
		if planObj.AccessMode != stateObj.AccessMode {
			volType, err := getVolumeType(r.client, planObj.VolumeID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume",
					"unexpected error: "+err.Error(),
				)
				return
			}

			err12 := volType.SetVolumeMappingAccessMode(planObj.AccessMode.ValueString(), plan.ID.ValueString())
			if err12 != nil {
				resp.Diagnostics.AddError(
					"Error setting access mode to sdc: "+plan.ID.ValueString(),
					"unexpected error: "+err12.Error(),
				)
			}
		}

	}

	sdcType, err1 := getSdcType(r.client, state.ID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error Getting SDC type: "+state.ID.String(),
			"unexpected error: "+err1.Error(),
		)
	}

	mappedVolumes, err2 := sdcType.GetVolume()
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volumes mapped to SDC : "+state.ID.String(),
			"unexpected error: "+err2.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := updateSDCVolMapState(mappedVolumes, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdcVolumeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sdcVolumeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove the mapped volumes to SDC
	if len(state.VolumeList.Elements()) > 0 {
		volList := []sdcVolumeModel{}
		diags = state.VolumeList.ElementsAs(ctx, &volList, true)
		resp.Diagnostics.Append(diags...)

		for _, vol := range volList {
			volType, err := getVolumeType(r.client, vol.VolumeID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume",
					"unexpected error: "+err.Error(),
				)
				return
			}

			err = volType.UnmapVolumeSdc(
				&goscaleio_types.UnmapVolumeSdcParam{
					SdcID: state.ID.ValueString(),
				},
			)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error Unmapping Volume to SDCs",
					"Couldn't unmap volume to SDC with id: "+state.ID.ValueString()+", unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *sdcVolumeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}