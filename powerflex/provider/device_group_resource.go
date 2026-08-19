/*
Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"terraform-provider-powerflex/clientgen"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &deviceGroupResource{}
	_ resource.ResourceWithConfigure   = &deviceGroupResource{}
	_ resource.ResourceWithImportState = &deviceGroupResource{}
)

// NewDeviceGroupResource returns a new device group resource instance.
func NewDeviceGroupResource() resource.Resource {
	return &deviceGroupResource{}
}

// deviceGroupResource is the resource implementation.
type deviceGroupResource struct {
	client *clientgen.APIClient
}

// deviceGroupResourceModel maps the resource schema data.
type deviceGroupResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	NumOfDevices       types.Int64  `tfsdk:"num_of_devices"`
	UsableCapacityInKb types.Int64  `tfsdk:"usable_capacity_in_kb"`
	UsedCapacityInKb   types.Int64  `tfsdk:"used_capacity_in_kb"`
	Status             types.String `tfsdk:"status"`
}

func (r *deviceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

func (r *deviceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage Device Groups on PowerFlex Gen2 systems (5.0+). Device Groups allow logical grouping of devices for management.",
		MarkdownDescription: "This resource is used to manage Device Groups on PowerFlex Gen2 systems (5.0+). Device Groups allow logical grouping of devices for management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The unique identifier of the device group.",
				MarkdownDescription: "The unique identifier of the device group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name of the device group.",
				MarkdownDescription: "The name of the device group.",
			},
			"protection_domain_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the protection domain.",
				MarkdownDescription: "The ID of the protection domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"num_of_devices": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of devices in the device group.",
				MarkdownDescription: "The number of devices in the device group.",
			},
			"usable_capacity_in_kb": schema.Int64Attribute{
				Computed:            true,
				Description:         "The usable capacity in KB.",
				MarkdownDescription: "The usable capacity in KB.",
			},
			"used_capacity_in_kb": schema.Int64Attribute{
				Computed:            true,
				Description:         "The used capacity in KB.",
				MarkdownDescription: "The used capacity in KB.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the device group.",
				MarkdownDescription: "The status of the device group.",
			},
		},
	}
}

func (r *deviceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*powerflexProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *powerflexProvider, got: %T.", req.ProviderData),
		)
		return
	}

	if p.genClient == nil {
		resp.Diagnostics.AddError(
			"Gen2 API Client Not Available",
			"The Gen2 API client is required for device group management. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	r.client = p.genClient
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deviceGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := clientgen.DeviceGroupParam{
		ProtectionDomainId: plan.ProtectionDomainID.ValueString(),
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		params.Name = clientgen.PtrString(plan.Name.ValueString())
	}

	tflog.Debug(ctx, "Creating Device Group")

	result, httpResponse, err := r.client.DeviceGroupAPI.CreateDeviceGroup(ctx).DeviceGroupParam(params).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("CreateDeviceGroup error: %v, status: %d", err, httpResponse.StatusCode))
		resp.Diagnostics.AddError(
			"Error Creating Device Group",
			"Could not create device group: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("CreateDeviceGroup success, status: %d", httpResponse.StatusCode))

	// Check for non-2xx status codes - this should not happen if err is nil, but check anyway
	if httpResponse != nil && (httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300) {
		tflog.Debug(ctx, fmt.Sprintf("CreateDeviceGroup non-2xx status: %d", httpResponse.StatusCode))
		resp.Diagnostics.AddError(
			"Error Creating Device Group",
			fmt.Sprintf("API returned status %d", httpResponse.StatusCode),
		)
		return
	}

	// Check if result is nil (shouldn't happen if no error, but defensive)
	if result == nil {
		resp.Diagnostics.AddError(
			"Error Creating Device Group",
			"API returned nil result",
		)
		return
	}

	// Read back the created device group
	groupID := result.GetId()
	tflog.Debug(ctx, fmt.Sprintf("Reading device group with ID: %s", groupID))
	group, _, err := r.client.DeviceGroupAPI.GetDeviceGroup(ctx, groupID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device Group After Creation",
			"Could not read device group: "+err.Error(),
		)
		return
	}

	r.mapDeviceGroupToState(group, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Created Device Group with ID: %s", groupID))
}

// Read refreshes the Terraform state with the latest data.
func (r *deviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, _, err := r.client.DeviceGroupAPI.GetDeviceGroup(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device Group",
			fmt.Sprintf("Could not read device group ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	r.mapDeviceGroupToState(group, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state deviceGroupResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle name change
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != state.Name.ValueString() {
		namePayload := clientgen.SetDeviceGroupName{
			Name: plan.Name.ValueString(),
		}
		_, err := r.client.DeviceGroupAPI.SetDeviceGroupName(ctx, state.ID.ValueString()).SetDeviceGroupName(namePayload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Renaming Device Group",
				fmt.Sprintf("Could not rename device group: %s", err.Error()),
			)
			return
		}
	}

	// Read back the updated device group
	group, _, err := r.client.DeviceGroupAPI.GetDeviceGroup(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device Group After Update",
			err.Error(),
		)
		return
	}

	r.mapDeviceGroupToState(group, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeviceGroupAPI.DeleteDeviceGroup(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Device Group",
			fmt.Sprintf("Could not delete device group ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Deleted Device Group ID: %s", state.ID.ValueString()))
}

func (r *deviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapDeviceGroupToState maps the API response to the Terraform state model.
func (r *deviceGroupResource) mapDeviceGroupToState(group *clientgen.DeviceGroup, model *deviceGroupResourceModel) {
	if group == nil {
		return
	}

	if group.Id != nil {
		model.ID = types.StringValue(*group.Id)
	}
	if group.Name != nil {
		model.Name = types.StringValue(*group.Name)
	}
	if group.ProtectionDomainId != nil {
		model.ProtectionDomainID = types.StringValue(*group.ProtectionDomainId)
	}
	if group.NumOfDevices != nil {
		model.NumOfDevices = types.Int64Value(int64(*group.NumOfDevices))
	}
	if group.UsableCapacityInKb != nil {
		model.UsableCapacityInKb = types.Int64Value(int64(*group.UsableCapacityInKb))
	}
	if group.UsedCapacityInKb != nil {
		model.UsedCapacityInKb = types.Int64Value(int64(*group.UsedCapacityInKb))
	}
	if group.Status != nil {
		model.Status = types.StringValue(*group.Status)
	}
}
