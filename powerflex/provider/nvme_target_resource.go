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
	"fmt"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewNvmeTargetResource is a helper function to simplify the provider implementation.
func NewNvmeTargetResource() resource.Resource {
	return &NvmeTargetResource{}
}

// NvmeTargetResource is the resource implementation.
type NvmeTargetResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

// Metadata describes the resource arguments.
func (r *NvmeTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvme_target"
}

// Schema describes the resource arguments.
func (r *NvmeTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource is used to manager NVMe target from the PowerFlex array. We can Create, Update and Delete the PowerFlex NVMe target using this resource. We can also import an existing NVMe target from PowerFlex array.  \n" +
			" **Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",
		MarkdownDescription: "This resource is used to manager NVMe target from the PowerFlex array. We can Create, Update and Delete the PowerFlex NVMe target using this resource. We can also import an existing NVMe target from PowerFlex array.  \n" +
			" **Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the NVMe target",
				MarkdownDescription: "ID of the NVMe target",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the NVMe target",
				MarkdownDescription: "Name of the NVMe target",
				Required:            true,
			},
			"system_id": schema.StringAttribute{
				Description:         "The ID of the system.",
				MarkdownDescription: "The ID of the system.",
				Computed:            true,
			},
			"protection_domain_id": schema.StringAttribute{
				Description: "ID of the Protection Domain under which the NVMe target will be created." +
					" Conflicts with 'protection_domain_name'." +
					" Cannot be updated." +
					"**Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",
				MarkdownDescription: "ID of the Protection Domain under which the NVMe target will be created." +
					" Conflicts with 'protection_domain_name'." +
					" Cannot be updated." +
					"**Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name"), path.MatchRoot("protection_domain_id")),
				},
			},
			"protection_domain_name": schema.StringAttribute{
				Description: "Name of the Protection Domain under which the NVMe target will be created." +
					" Conflicts with 'protection_domain_id'." +
					" Cannot be updated." +
					"**Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",
				MarkdownDescription: "Name of the Protection Domain under which the NVMe target will be created." +
					" Conflicts with 'protection_domain_id'." +
					" Cannot be updated." +
					"**Note:** Either `protection_domain_id` or `protection_domain_name` must be specified.",
				Optional: true,
			},
			"ip_list": schema.SetNestedAttribute{
				Description:         "List of IPs associated with the NVMe target.",
				MarkdownDescription: "List of IPs associated with the NVMe target.",
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description:         "NVMe Target IP.",
							MarkdownDescription: "NVMe Target IP.",
							Required:            true,
						},
						"role": schema.StringAttribute{
							Description:         "NVMe Target IP role.",
							MarkdownDescription: "NVMe Target IP role.",
							Required:            true,
							Validators: []validator.String{stringvalidator.OneOf(
								"StorageAndHost",
								"StorageOnly",
								"HostOnly",
							)},
						},
					},
				},
			},
			"storage_port": schema.Int64Attribute{
				Description:         "The storage port of the NVMe target (default: 12200).",
				MarkdownDescription: "The storage port of the NVMe target (default: 12200).",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(12200),
				Validators: []validator.Int64{
					int64validator.Between(1025, 65535),
				},
			},
			"nvme_port": schema.Int64Attribute{
				Description:         "The NVMe port of the NVMe target (default: 4420).",
				MarkdownDescription: "The NVMe port of the NVMe target (default: 4420).",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(4420),
				Validators: []validator.Int64{
					int64validator.Between(1025, 65535),
				},
			},
			"discovery_port": schema.Int64Attribute{
				Description:         "The discovery port of the NVMe target (default: 8009).",
				MarkdownDescription: "The discovery port of the NVMe target (default: 8009).",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(8009),
				Validators: []validator.Int64{
					int64validator.Between(1025, 65535),
				},
			},
			"sdt_state": schema.StringAttribute{
				Description:         "The state of the NVMe target.",
				MarkdownDescription: "The state of the NVMe target.",
				Computed:            true,
			},
			"mdm_connection_state": schema.StringAttribute{
				Description:         "The MDM connection state of the NVMe target.",
				MarkdownDescription: "The MDM connection state of the NVMe target.",
				Computed:            true,
			},
			"membership_state": schema.StringAttribute{
				Description:         "The membership state of the NVMe target.",
				MarkdownDescription: "The membership state of the NVMe target.",
				Computed:            true,
			},
			"fault_set_id": schema.StringAttribute{
				Description:         "The fault set ID of the NVMe target.",
				MarkdownDescription: "The fault set ID of the NVMe target.",
				Computed:            true,
			},
			"software_version_info": schema.StringAttribute{
				Description:         "The software version information of the NVMe target.",
				MarkdownDescription: "The software version information of the NVMe target.",
				Computed:            true,
			},
			"maintenance_state": schema.StringAttribute{
				Description:         "The maintenance state of the NVMe target. Active to trigger enter maintenance mode; Inactive to exit maintenance mode.",
				MarkdownDescription: "The maintenance state of the NVMe target. Active to trigger enter maintenance mode; Inactive to exit maintenance mode.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					"Active",
					"Inactive",
				)},
			},
			"authentication_error": schema.StringAttribute{
				Description:         "The authentication error of the NVMe target.",
				MarkdownDescription: "The authentication error of the NVMe target.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NvmeTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create used to Create NVMe target Resource
func (r *NvmeTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create NVMe target")

	var plan models.NvmeTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := scaleiotypes.SdtParam{}
	err := helper.ReadFromState(ctx, plan, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing plan model",
			fmt.Sprintf("Could not read NVMe target param with error: %s", err.Error()),
		)
		return
	}

	if !plan.Name.IsNull() && plan.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Name cannot be empty",
			"Name cannot be empty",
		)
		return
	}

	if len(plan.IPList) >= 0 {
		for _, ip := range plan.IPList {
			params.IPList = append(params.IPList, &scaleiotypes.SdtIP{IP: ip.IP.ValueString(), Role: ip.Role.ValueString()})
		}
	}

	pdm, err := helper.GetNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	params.ProtectionDomainID = pdm.ProtectionDomain.ID
	nvmeTargetResp, err := pdm.CreateSdt(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create NVMe target with name %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	if !plan.MaintenanceState.IsUnknown() && !plan.MaintenanceState.IsNull() {
		err = helper.ToggleSdtMaintenanceMode(r.system, nvmeTargetResp.ID, plan.MaintenanceState.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set maintenance mode to %s", plan.MaintenanceState.ValueString()),
				err.Error(),
			)
			return
		}
	}

	nvmeTarget, err := helper.GetNvmeTargetByID(r.system, nvmeTargetResp.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not read NVMe target with ID %s", nvmeTargetResp.ID),
			err.Error(),
		)
		return
	}

	err = helper.CopyFields(ctx, nvmeTarget, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating NVMe target",
			fmt.Sprintf("Could not parse NVMe target struct with error: %s", err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read NVMe target Resource
func (r *NvmeTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read NVMe target")
	// Get current state
	var state models.NvmeTargetResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nvmeTarget, err := helper.GetNvmeTargetByID(r.system, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the NVMe target Details",
			err.Error(),
		)
		return
	}

	err = helper.CopyFields(ctx, nvmeTarget, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading NVMe target",
			fmt.Sprintf("Could not parse NVMe target struct with error: %s", err.Error()),
		)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update used to NVMe target Resource
func (r *NvmeTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.NvmeTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.NvmeTargetResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.IsNull() && plan.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Name cannot be empty",
			"Name cannot be empty",
		)
		return
	}

	if !plan.ProtectionDomainID.IsUnknown() && plan.ProtectionDomainID.ValueString() != state.ProtectionDomainID.ValueString() {
		resp.Diagnostics.AddError(
			"Protection domain ID cannot be updated",
			"Protection domain ID cannot be updated")
		return
	}

	if !plan.ProtectionDomainName.IsNull() && plan.ProtectionDomainName.ValueString() != state.ProtectionDomainName.ValueString() {
		protectionDomain, err := r.system.FindProtectionDomain("", plan.ProtectionDomainName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Read Powerflex Protection domain by name %v", plan.ProtectionDomainName.ValueString()),
				err.Error(),
			)
			return
		}
		if protectionDomain.ID != state.ProtectionDomainID.ValueString() {
			resp.Diagnostics.AddError(
				"Protection domain name does not match the original Protection domain",
				"Protection domain name does not match the original Protection domain")
			return
		}
		state.ProtectionDomainName = types.StringValue(protectionDomain.Name)
	}

	err := helper.NvmeTargetUpdate(r.system, state, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update NVMe target",
			err.Error(),
		)
	}

	if !plan.MaintenanceState.IsUnknown() && !plan.MaintenanceState.IsNull() &&
		plan.MaintenanceState.ValueString() != state.MaintenanceState.ValueString() {
		err = helper.ToggleSdtMaintenanceMode(r.system, state.ID.ValueString(), plan.MaintenanceState.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set maintenance mode to %s", plan.MaintenanceState.ValueString()),
				err.Error(),
			)
		}
	}

	nvmeTarget, err := helper.GetNvmeTargetByID(r.system, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the NVMe target Details",
			err.Error(),
		)
		return
	}

	err = helper.CopyFields(ctx, nvmeTarget, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating NVMe target",
			fmt.Sprintf("Could not parse NVMe target struct with error: %s", err.Error()),
		)
		return
	}
	// in case of switching back and forth between PD ID and name
	state.ProtectionDomainName = plan.ProtectionDomainName

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *NvmeTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.NvmeTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.system.DeleteSdt(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete NVMe target",
			err.Error(),
		)

		return
	}
	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource state.
func (r *NvmeTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
