/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"strconv"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewSDCVolumesMappingResource is a helper function to simplify the provider implementation.
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
		Description:         "This resource can be used to manage mapping of volumes to an SDC on a PowerFlex array.",
		MarkdownDescription: "This resource can be used to manage mapping of volumes to an SDC on a PowerFlex array.",
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
				Description:         "List of volumes mapped to SDC. At least one of 'volume_id' and 'volume_name' is required.",
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "List of volumes mapped to SDC. At least one of `volume_id` and `volume_name` is required.",
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
							Description:         "IOPS limit. Valid values are 0 or integers greater than 10. '0' represents unlimited IOPS. Default value is '0'.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "IOPS limit. Valid values are 0 or integers greater than 10. `0` represents unlimited IOPS. Default value is `0`.",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"limit_bw_in_mbps": schema.Int64Attribute{
							Description:         "Bandwidth limit in MBPS. '0' represents unlimited bandwith. Default value is '0'.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Bandwidth limit in MBPS. `0` represents unlimited bandwith. Default value is `0`.",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"access_mode": schema.StringAttribute{
							Description:         "The Access Mode of the SDC. Valid values are 'ReadOnly', 'ReadWrite' and 'NoAccess'. Default value is 'ReadOnly'.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The Access Mode of the SDC. Valid values are `ReadOnly`, `ReadWrite` and `NoAccess`. Default value is `ReadOnly`.",
							Validators: []validator.String{stringvalidator.OneOf(
								"ReadOnly",
								"ReadWrite",
								"NoAccess",
							)},
							PlanModifiers: []planmodifier.String{
								helper.StringDefault("ReadOnly"),
								stringplanmodifier.UseStateForUnknown(),
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
	var plan models.SdcVolumeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Get the system on the PowerFlex cluster
	system, err := helper.GetFirstSystem(r.client)
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

	_ = r.VerifyVolumes(ctx, &plan)

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdcVolumeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.SdcVolumeMappingResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.VerifyVolumes(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volList := []models.SdcVolumeModel{}
	diags = plan.VolumeList.ElementsAs(ctx, &volList, true)
	resp.Diagnostics.Append(diags...)

	for _, vol := range volList {
		volType, err := helper.GetVolumeType(r.client, vol.VolumeID.ValueString())
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

	sdcType, err1 := helper.GetSdcType(r.client, plan.ID.ValueString())
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
	state, dgs := helper.UpdateSDCVolMapState(mappedVolumes, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sdcVolumeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state models.SdcVolumeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdcType, err1 := helper.GetSdcType(r.client, state.ID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error Getting SDC type: "+state.ID.String(),
			"unexpected error: "+err1.Error(),
		)
		return
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
	state, dgs := helper.UpdateSDCVolMapState(mappedVolumes, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdcVolumeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan models.SdcVolumeMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.VerifyVolumes(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state models.SdcVolumeMappingResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate if there is change in plan and state w.r.t SDC ID
	if !plan.ID.IsUnknown() && plan.ID.ValueString() != state.ID.ValueString() {
		resp.Diagnostics.AddError(
			"SDC ID cannot be updated",
			"SDC ID cannot be updated")
		return
	}

	planVolList := []models.SdcVolumeModel{}
	stateVolList := []models.SdcVolumeModel{}

	// Populate planVolList with volumes defined in plan
	diags = plan.VolumeList.ElementsAs(ctx, &planVolList, true)
	resp.Diagnostics.Append(diags...)

	// Populate stateVolList with volumes stored in state
	diags = state.VolumeList.ElementsAs(ctx, &stateVolList, true)
	resp.Diagnostics.Append(diags...)

	planVolIds := make(map[string]string)
	stateVolIds := make(map[string]string)

	// Populate planVolIds with the volume IDs defined in plan
	for _, vol := range planVolList {
		planVolIds[vol.VolumeID.ValueString()] = vol.VolumeID.ValueString()
	}

	// Populate stateVolIds with the volume IDs stored in state
	for _, vol := range stateVolList {
		stateVolIds[vol.VolumeID.ValueString()] = vol.VolumeID.ValueString()
	}

	// mapVolIds will be storing the volume ids for which mapping needs to be performed
	mapVolIds := helper.DifferenceMap(planVolIds, stateVolIds)

	// unmapVolIds will be storing the volume id for which unmapping needs to be performed
	unmapVolIds := helper.DifferenceMap(stateVolIds, planVolIds)

	// nonchangeVolIds will be storing the volume id for which limits/access mode needs to be performed
	nonchangeVolIds := helper.DifferenceMap(planVolIds, mapVolIds)

	// Perform mapping and setting limits operation
	for _, planVol := range planVolList {
		if _, ok := mapVolIds[planVol.VolumeID.ValueString()]; ok {
			volType, err := helper.GetVolumeType(r.client, planVol.VolumeID.ValueString())
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
					"Error mapping volume to sdc: "+planVol.VolumeID.ValueString(),
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

	// Perform unmap operation
	for volID := range unmapVolIds {
		volType, err := helper.GetVolumeType(r.client, volID)
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
	for volID := range nonchangeVolIds {
		var planObj models.SdcVolumeModel
		var stateObj models.SdcVolumeModel

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
		if (!planObj.IOPSLimit.IsUnknown() && planObj.IOPSLimit != stateObj.IOPSLimit) || (!planObj.BWLimit.IsUnknown() && planObj.BWLimit != stateObj.BWLimit) {
			smslp := goscaleio_types.SetMappedSdcLimitsParam{
				SdcID:                plan.ID.ValueString(),
				BandwidthLimitInKbps: strconv.FormatInt(int64(planObj.BWLimit.ValueInt64()*1024), 10),
				IopsLimit:            strconv.FormatInt(int64(planObj.IOPSLimit.ValueInt64()), 10),
			}

			volType, err := helper.GetVolumeType(r.client, planObj.VolumeID.ValueString())
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
		if !planObj.AccessMode.IsUnknown() && planObj.AccessMode != stateObj.AccessMode {
			volType, err := helper.GetVolumeType(r.client, planObj.VolumeID.ValueString())
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

	sdcType, err1 := helper.GetSdcType(r.client, state.ID.ValueString())
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
	state, dgs := helper.UpdateSDCVolMapState(mappedVolumes, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdcVolumeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.SdcVolumeMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove the mapped volumes to SDC
	if len(state.VolumeList.Elements()) > 0 {
		volList := []models.SdcVolumeModel{}
		diags = state.VolumeList.ElementsAs(ctx, &volList, true)
		resp.Diagnostics.Append(diags...)

		for _, vol := range volList {
			volType, err := helper.GetVolumeType(r.client, vol.VolumeID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting volume",
					"unexpected error: "+err.Error(),
				)
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
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource
func (r *sdcVolumeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// VerifyVolumes checks if volume exists or not before mapping
func (r *sdcVolumeMappingResource) VerifyVolumes(ctx context.Context, plan *models.SdcVolumeMappingResourceModel) (diags diag.Diagnostics) {
	volList := []models.SdcVolumeModel{}
	diag := plan.VolumeList.ElementsAs(ctx, &volList, true)
	diags.Append(diag...)

	for index, vol := range volList {
		// Populate volume name in the plan if volume ID is provided in the config
		if !vol.VolumeID.IsUnknown() {
			volume, err := r.client.GetVolume("", vol.VolumeID.ValueString(), "", "", false)
			if err != nil {
				diags.AddError(
					"Error getting volume with ID: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			volList[index].VolumeName = types.StringValue(volume[0].Name)
		} else if !vol.VolumeName.IsUnknown() {
			volume, err := r.client.GetVolume("", "", "", vol.VolumeName.ValueString(), false)
			if err != nil {
				diags.AddError(
					"Error getting volume with name: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			volList[index].VolumeID = types.StringValue(volume[0].ID)
		}
	}

	// Modify the plan to populate volume list
	if len(volList) > 0 {
		mappedVolumeList, dgs := helper.GetVolSetValueFromItems(volList)
		diags.Append(dgs...)
		plan.VolumeList = mappedVolumeList
	}
	return diags
}
