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
	"strconv"

	"terraform-provider-powerflex/clientgen"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &deviceActionV2Resource{}
	_ resource.ResourceWithConfigure = &deviceActionV2Resource{}
)

// NewDeviceActionV2Resource returns a new device action v2 resource instance.
func NewDeviceActionV2Resource() resource.Resource {
	return &deviceActionV2Resource{}
}

// deviceActionV2Resource is the resource implementation for Gen2 device actions.
type deviceActionV2Resource struct {
	client *clientgen.APIClient
}

// deviceActionV2ResourceModel maps the resource schema data.
type deviceActionV2ResourceModel struct {
	ID                types.String `tfsdk:"id"`
	DeviceID          types.String `tfsdk:"device_id"`
	Action            types.String `tfsdk:"action"`
	CapacityLimitInKb types.String `tfsdk:"capacity_limit_in_kb"`
}

func (r *deviceActionV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action"
}

func (r *deviceActionV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to perform Gen2 device actions on PowerFlex systems (5.0+). Supported actions: activate, clear_error, set_capacity_limit.",
		MarkdownDescription: "This resource is used to perform Gen2 device actions on PowerFlex systems (5.0+). Supported actions: `activate`, `clear_error`, `set_capacity_limit`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The identifier for this action resource.",
				MarkdownDescription: "The identifier for this action resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the device to perform the action on.",
				MarkdownDescription: "The ID of the device to perform the action on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action": schema.StringAttribute{
				Required:            true,
				Description:         "The action to perform. Valid values: activate, clear_error, set_capacity_limit.",
				MarkdownDescription: "The action to perform. Valid values: `activate`, `clear_error`, `set_capacity_limit`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("activate", "clear_error", "set_capacity_limit"),
				},
			},
			"capacity_limit_in_kb": schema.StringAttribute{
				Optional:            true,
				Description:         "Capacity limit in KB. Required when action is set_capacity_limit.",
				MarkdownDescription: "Capacity limit in KB. Required when action is `set_capacity_limit`.",
			},
		},
	}
}

func (r *deviceActionV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"The Gen2 API client is required for device actions. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	r.client = p.genClient
}

// Create performs the device action.
func (r *deviceActionV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deviceActionV2ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := plan.DeviceID.ValueString()
	action := plan.Action.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing device action '%s' on device '%s'", action, deviceID))

	switch action {
	case "activate":
		_, err := r.client.DeviceAPI.ActivateDevice(ctx, deviceID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating Device",
				fmt.Sprintf("Could not activate device %s: %s", deviceID, err.Error()),
			)
			return
		}

	case "clear_error":
		_, err := r.client.DeviceAPI.ClearDeviceError(ctx, deviceID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Clearing Device Error",
				fmt.Sprintf("Could not clear error on device %s: %s", deviceID, err.Error()),
			)
			return
		}

	case "set_capacity_limit":
		if plan.CapacityLimitInKb.IsNull() || plan.CapacityLimitInKb.IsUnknown() {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"capacity_limit_in_kb is required when action is set_capacity_limit.",
			)
			return
		}

		// Validate it's a valid number
		capStr := plan.CapacityLimitInKb.ValueString()
		if _, err := strconv.ParseInt(capStr, 10, 64); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Capacity Limit",
				fmt.Sprintf("capacity_limit_in_kb must be a valid integer: %s", err.Error()),
			)
			return
		}

		payload := clientgen.SetDeviceCapacityLimitGen2{
			CapacityLimitInKb: capStr,
		}
		_, err := r.client.DeviceAPI.SetDeviceCapacityLimitGen2(ctx, deviceID).SetDeviceCapacityLimitGen2(payload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Device Capacity Limit",
				fmt.Sprintf("Could not set capacity limit on device %s: %s", deviceID, err.Error()),
			)
			return
		}

	default:
		resp.Diagnostics.AddError(
			"Unknown Action",
			fmt.Sprintf("Unknown action '%s'", action),
		)
		return
	}

	// Set the state
	plan.ID = types.StringValue(fmt.Sprintf("%s-%s", deviceID, action))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Completed device action '%s' on device '%s'", action, deviceID))
}

// Read is a no-op for action resources.
func (r *deviceActionV2Resource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Action resources don't have a persistent state to read
}

// Update handles updates to the action resource.
func (r *deviceActionV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deviceActionV2ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For set_capacity_limit, re-execute with the new value
	if plan.Action.ValueString() == "set_capacity_limit" {
		if plan.CapacityLimitInKb.IsNull() || plan.CapacityLimitInKb.IsUnknown() {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"capacity_limit_in_kb is required when action is set_capacity_limit.",
			)
			return
		}

		capStr := plan.CapacityLimitInKb.ValueString()
		if _, err := strconv.ParseInt(capStr, 10, 64); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Capacity Limit",
				fmt.Sprintf("capacity_limit_in_kb must be a valid integer: %s", err.Error()),
			)
			return
		}

		payload := clientgen.SetDeviceCapacityLimitGen2{
			CapacityLimitInKb: capStr,
		}
		_, err := r.client.DeviceAPI.SetDeviceCapacityLimitGen2(ctx, plan.DeviceID.ValueString()).SetDeviceCapacityLimitGen2(payload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Device Capacity Limit",
				fmt.Sprintf("Could not set capacity limit: %s", err.Error()),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete is a no-op for action resources.
func (r *deviceActionV2Resource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Action resources don't need cleanup
}
