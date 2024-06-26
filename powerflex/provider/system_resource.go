/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sdcResource{}
	_ resource.ResourceWithConfigure   = &sdcResource{}
	_ resource.ResourceWithImportState = &sdcResource{}
)

// SystemResource - function to return resource interface
func SystemResource() resource.Resource {
	return &systemResource{}
}

// systemResource - struct to define system resource
type systemResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

// Metadata - function to return metadata for system resource.
func (r *systemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

// Schema - function to return Schema for system resource.
func (r *systemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SystemResourceSchema
}

// Configure - function to return Configuration for system resource.
func (r *systemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	r.client = req.ProviderData.(*powerflexProvider).client
	system, err := helper.GetFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	r.system = system
}

// SystemResourceSchema defines the schema for system resource
var SystemResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the cluster level operations of the PowerFlex Array. This resource supports Create, Update and Delete operations.",
	MarkdownDescription: "This resource is used to manage the cluster level operations of the PowerFlex Array. This resource supports Create, Update and Delete operations.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "System ID",
			MarkdownDescription: "System ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"restricted_mode": schema.StringAttribute{
			Required:            true,
			Description:         "Restricted mode of the cluster. Accepted values are `None`, `Guid`, `ApprovedIp`.",
			MarkdownDescription: "Restricted mode of the cluster. Accepted values are `None`, `Guid`, `ApprovedIp`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"None",
				"Guid",
				"ApprovedIp",
			),
			},
		},
		"sdc_guids": schema.ListAttribute{
			Optional:            true,
			Computed:            true,
			Description:         "Specifies list of SDC GUIDs.",
			MarkdownDescription: "Specifies list of SDC GUIDs.",
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				listvalidator.ConflictsWith(
					path.MatchRoot("sdc_ids"),
					path.MatchRoot("sdc_names"),
				),
			},
		},
		"sdc_ids": schema.ListAttribute{
			Optional:            true,
			Description:         "Specifies list of SDC IDs.",
			MarkdownDescription: "Specifies list of SDC IDs.",
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				listvalidator.ConflictsWith(
					path.MatchRoot("sdc_names"),
				),
			},
		},
		"sdc_names": schema.ListAttribute{
			Optional:            true,
			Description:         "Specifies list of SDC names.",
			MarkdownDescription: "Specifies list of SDC names.",
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"sdc_approved_ips": schema.ListNestedAttribute{
			Optional:            true,
			Description:         "Specifies list of SDC IPs.",
			MarkdownDescription: "Specifies list of SDC IPs.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "SDC ID.",
						MarkdownDescription: "SDC ID.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"ips": schema.SetAttribute{
						Description:         "SDC IPs.",
						MarkdownDescription: "SDC IPs.",
						Required:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
		},
	},
}

func (r *systemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "In create operation")
	// Retrieve values from plan
	var (
		plan    models.SystemModel
		systems []*scaleiotypes.System
		err     error
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.PopulateSDCDetails(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Change restricted mode
	err = r.system.SetRestrictedMode(plan.RestrictedMode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error changing restricted mode",
			"Could not change restricted mode, unexpected err: "+err.Error(),
		)
		return
	}

	SdcGuids := make([]string, 0)
	diags.Append(plan.SdcGuids.ElementsAs(ctx, &SdcGuids, true)...)

	// Approve SDCs if restricted mode is set to Guid
	if plan.RestrictedMode.ValueString() == "Guid" {
		for _, SdcGUID := range SdcGuids {
			diags.Append(r.ApproveSdcGUID(SdcGUID)...)
		}
	}

	// Approve SDC IPs
	if !plan.SdcApprovedIPs.IsNull() {
		planApprovedIPs := make([]models.SdcApprovedIPsModel, 0)
		diags.Append(plan.SdcApprovedIPs.ElementsAs(ctx, &planApprovedIPs, true)...)

		for _, sdcPlan := range planApprovedIPs {
			// Check if SDC exists with the given ID
			sdc, err := r.system.FindSdc("ID", sdcPlan.ID.ValueString())
			if err != nil {
				diags.AddError(
					"Error getting SDC with ID: "+sdcPlan.ID.ValueString(),
					"unexpected error: "+err.Error(),
				)
				continue
			}

			// Approve SDC if its not approved already
			if !sdc.Sdc.SdcApproved {
				diags.Append(r.ApproveSdcIP(sdc.Sdc.SdcIPs)...)
			}

			sdcIps := make([]string, 0)
			diags.Append(sdcPlan.IPs.ElementsAs(ctx, &sdcIps, true)...)
			err = r.system.SetApprovedIps(sdc.Sdc.ID, sdcIps)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in setting approved IPs",
					"Error in setting approved IPs for SDC "+sdc.Sdc.ID+", unexpected err: "+err.Error(),
				)
			}
		}
	}

	systems, err = r.client.GetInstance("")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			"Could not get system instance, unexpected err: "+err.Error(),
		)
	}

	state, diags := helper.UpdateSystemState(&plan, systems[0], r.system)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *systemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "In read operation")
	var (
		state   *models.SystemModel
		systems []*scaleiotypes.System
		err     error
	)

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	systems, err = r.client.GetInstance("")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			"Could not get system instance, unexpected err: "+err.Error(),
		)
	}

	state, diags = helper.UpdateSystemState(state, systems[0], r.system)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *systemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "In update operation")
	var (
		plan  models.SystemModel
		state *models.SystemModel
		err   error
	)

	// Get the plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Get Current State
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.PopulateSDCDetails(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Change restricted mode
	if plan.RestrictedMode.ValueString() != state.RestrictedMode.ValueString() {
		err = r.system.SetRestrictedMode(plan.RestrictedMode.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error changing restricted mode",
				"Could not change restricted mode, unexpected err: "+err.Error(),
			)
		}
	}

	systems, err := r.client.GetInstance("")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			"Could not get system instance, unexpected err: "+err.Error(),
		)
	}

	planSdcGuids := make([]string, 0)
	diags.Append(plan.SdcGuids.ElementsAs(ctx, &planSdcGuids, true)...)

	stateSdcGuids := make([]string, 0)
	diags.Append(state.SdcGuids.ElementsAs(ctx, &stateSdcGuids, true)...)

	approveSdcs, _ := helper.DifferenceArray(stateSdcGuids, planSdcGuids)

	for _, SdcGUID := range approveSdcs {
		diags.Append(r.ApproveSdcGUID(SdcGUID)...)
	}

	if !plan.SdcApprovedIPs.IsNull() {
		planSdcApprovedIPs := make([]models.SdcApprovedIPsModel, 0)
		stateSdcApprovedIPs := make([]models.SdcApprovedIPsModel, 0)
		planSdcMap := make(map[string][]string)
		stateSdcMap := make(map[string][]string)

		diags.Append(plan.SdcApprovedIPs.ElementsAs(ctx, &planSdcApprovedIPs, true)...)
		diags.Append(state.SdcApprovedIPs.ElementsAs(ctx, &stateSdcApprovedIPs, true)...)

		planSdcIds := make([]string, 0)
		stateSdcIds := make([]string, 0)

		for _, sdc := range planSdcApprovedIPs {
			planSdcIPs := make([]string, 0)
			planSdcIds = append(planSdcIds, sdc.ID.ValueString())
			diags.Append(sdc.IPs.ElementsAs(ctx, &planSdcIPs, true)...)
			planSdcMap[sdc.ID.ValueString()] = planSdcIPs
		}

		for _, sdc := range stateSdcApprovedIPs {
			stateSdcIPs := make([]string, 0)
			stateSdcIds = append(stateSdcIds, sdc.ID.ValueString())
			diags.Append(sdc.IPs.ElementsAs(ctx, &stateSdcIPs, true)...)
			stateSdcMap[sdc.ID.ValueString()] = stateSdcIPs
		}

		// Set approved IPs for new SDCs
		diffSdcs, _ := helper.DifferenceArray(stateSdcIds, planSdcIds)
		if len(diffSdcs) > 0 {
			for _, sdcID := range diffSdcs {
				sdc, err := r.system.FindSdc("ID", sdcID)
				if err != nil {
					diags.AddError(
						"Error getting SDC with ID: "+sdcID,
						"unexpected error: "+err.Error(),
					)
					continue
				}

				if !sdc.Sdc.SdcApproved {
					diags.Append(r.ApproveSdcIP(sdc.Sdc.SdcIPs)...)
				}

				err = r.system.SetApprovedIps(sdcID, planSdcMap[sdcID])
				if err != nil {
					resp.Diagnostics.AddError(
						"Error in setting approved IPs",
						"Error in setting approved IPs for SDC "+sdcID+", unexpected err: "+err.Error(),
					)
				}
			}
		}

		// Set approved IPs for already approved SDCs
		for key := range planSdcMap {
			if !helper.CompareStringSlice(planSdcMap[key], stateSdcMap[key]) {
				err = r.system.SetApprovedIps(key, planSdcMap[key])
				if err != nil {
					resp.Diagnostics.AddError(
						"Error in setting approved IPs",
						"Error in setting approved IPs for SDC "+key+", unexpected err: "+err.Error(),
					)
				}
			}
		}

	}

	state, diags = helper.UpdateSystemState(&plan, systems[0], r.system)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *systemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "In delete operation")
	// Retrieve values from state
	var state models.SystemModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *systemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	if req.ID != "" {
		systems, err := r.client.GetInstance("")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting system instance on the PowerFlex cluster",
				"Could not get system instance, unexpected err: "+err.Error(),
			)
			return
		}

		if req.ID != systems[0].ID {
			resp.Diagnostics.AddError(
				"Error in importing system",
				"Could not import system with ID: "+req.ID,
			)
			return
		}
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// PopulateSDCDetails populates the SDC Guid based on SDC ID/Name.
func (r *systemResource) PopulateSDCDetails(ctx context.Context, plan *models.SystemModel) (diags diag.Diagnostics) {
	SdcGuids := make([]string, 0)
	SdcIDs := make([]string, 0)
	SdcNames := make([]string, 0)

	if plan.SdcGuids.IsUnknown() && plan.SdcIDs.IsNull() && plan.SdcNames.IsNull() {
		plan.SdcGuids = types.ListNull(types.StringType)
	}

	if !plan.SdcIDs.IsNull() {
		diags.Append(plan.SdcIDs.ElementsAs(ctx, &SdcIDs, true)...)
		for _, id := range SdcIDs {
			sdc, err := r.system.FindSdc("ID", id)
			if err != nil {
				diags.AddError(
					"Error getting SDC with ID: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			SdcGuids = append(SdcGuids, sdc.Sdc.SdcGUID)
		}
		plan.SdcGuids, _ = types.ListValueFrom(ctx, types.StringType, SdcGuids)
	} else if !plan.SdcNames.IsNull() {
		diags.Append(plan.SdcNames.ElementsAs(ctx, &SdcNames, true)...)
		for _, name := range SdcNames {
			sdc, err := r.system.FindSdc("Name", name)
			if err != nil {
				diags.AddError(
					"Error getting SDC with Name: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
			SdcGuids = append(SdcGuids, sdc.Sdc.SdcGUID)
		}
		plan.SdcGuids, _ = types.ListValueFrom(ctx, types.StringType, SdcGuids)
	} else {
		diags.Append(plan.SdcGuids.ElementsAs(ctx, &SdcGuids, true)...)
		for _, guid := range SdcGuids {
			_, err := r.system.FindSdc("SdcGUID", guid)
			if err != nil {
				diags.AddError(
					"Error getting SDC with GUID: ",
					"unexpected error: "+err.Error(),
				)
				return
			}
		}
	}
	return
}

// ApproveSdcGUID approves the SDC based on given Guid
func (r *systemResource) ApproveSdcGUID(sdcGUID string) diag.Diagnostics {
	var diags diag.Diagnostics

	sdc, err := r.system.FindSdc("SdcGUID", sdcGUID)
	if err != nil {
		diags.AddError(
			"Error getting SDC with GUID: ",
			"unexpected error: "+err.Error(),
		)
		return diags
	}

	payload := scaleiotypes.ApproveSdcParam{
		SdcGUID: sdcGUID,
	}

	if !sdc.Sdc.SdcApproved {
		_, err := r.system.ApproveSdc(&payload)
		if err != nil {
			diags.AddError(
				"Error in approving SDC with GUID",
				"Error in approving SDC with GUID "+sdcGUID+"+, unexpected err: "+err.Error(),
			)
		}
	}
	return diags
}

// ApproveSdcIP approves the SDC based on given IPs
func (r *systemResource) ApproveSdcIP(sdcIPs []string) diag.Diagnostics {
	var diags diag.Diagnostics
	payload := scaleiotypes.ApproveSdcParam{
		SdcIps: sdcIPs,
	}

	_, err := r.system.ApproveSdc(&payload)
	if err != nil {
		diags.AddError(
			"Error in approving SDC with IP",
			"Error in approving SDC with IP: "+err.Error(),
		)
	}
	return diags
}
