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
	_ resource.Resource                = &storageNodeResource{}
	_ resource.ResourceWithConfigure   = &storageNodeResource{}
	_ resource.ResourceWithImportState = &storageNodeResource{}
)

// NewStorageNodeResource returns a new storage node resource instance.
func NewStorageNodeResource() resource.Resource {
	return &storageNodeResource{}
}

// storageNodeResource is the resource implementation.
type storageNodeResource struct {
	client *clientgen.APIClient
}

// storageNodeResourceModel maps the resource schema data.
type storageNodeResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	IPList             types.List   `tfsdk:"ip_list"`
	Port               types.String `tfsdk:"port"`
	FaultSetID         types.String `tfsdk:"fault_set_id"`
	Status             types.String `tfsdk:"status"`
	State              types.String `tfsdk:"state"`
	MembershipState    types.String `tfsdk:"membership_state"`
	MdmConnectionState types.String `tfsdk:"mdm_connection_state"`
	SoftwareVersion    types.String `tfsdk:"software_version"`
	NumOfDevices       types.Int64  `tfsdk:"num_of_devices"`
	TotalCapacityInKb  types.Int64  `tfsdk:"total_capacity_in_kb"`
	UsedCapacityInKb   types.Int64  `tfsdk:"used_capacity_in_kb"`
	MaintenanceState   types.String `tfsdk:"maintenance_state"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
}

// storageNodeIPModel maps IP list items.
type storageNodeIPModel struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

func (r *storageNodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_node"
}

func (r *storageNodeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage Storage Nodes on PowerFlex Gen2 systems (5.0+). Storage Nodes are the Gen2 replacement for SDS.",
		MarkdownDescription: "This resource is used to manage Storage Nodes on PowerFlex Gen2 systems (5.0+). Storage Nodes are the Gen2 replacement for SDS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The unique identifier of the storage node.",
				MarkdownDescription: "The unique identifier of the storage node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name of the storage node.",
				MarkdownDescription: "The name of the storage node.",
			},
			"protection_domain_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the protection domain.",
				MarkdownDescription: "The ID of the protection domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The port number for the storage node.",
				MarkdownDescription: "The port number for the storage node.",
			},
			"fault_set_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the fault set.",
				MarkdownDescription: "The ID of the fault set.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the storage node.",
				MarkdownDescription: "The status of the storage node.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "The state of the storage node.",
				MarkdownDescription: "The state of the storage node.",
			},
			"membership_state": schema.StringAttribute{
				Computed:            true,
				Description:         "The membership state of the storage node.",
				MarkdownDescription: "The membership state of the storage node.",
			},
			"mdm_connection_state": schema.StringAttribute{
				Computed:            true,
				Description:         "The MDM connection state of the storage node.",
				MarkdownDescription: "The MDM connection state of the storage node.",
			},
			"software_version": schema.StringAttribute{
				Computed:            true,
				Description:         "The software version of the storage node.",
				MarkdownDescription: "The software version of the storage node.",
			},
			"num_of_devices": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of devices attached to the storage node.",
				MarkdownDescription: "The number of devices attached to the storage node.",
			},
			"total_capacity_in_kb": schema.Int64Attribute{
				Computed:            true,
				Description:         "The total capacity of the storage node in KB.",
				MarkdownDescription: "The total capacity of the storage node in KB.",
			},
			"used_capacity_in_kb": schema.Int64Attribute{
				Computed:            true,
				Description:         "The used capacity of the storage node in KB.",
				MarkdownDescription: "The used capacity of the storage node in KB.",
			},
			"maintenance_state": schema.StringAttribute{
				Computed:            true,
				Description:         "The maintenance state of the storage node.",
				MarkdownDescription: "The maintenance state of the storage node.",
			},
			"performance_profile": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The performance profile of the storage node.",
				MarkdownDescription: "The performance profile of the storage node.",
			},
		},
		Blocks: map[string]schema.Block{
			"ip_list": schema.ListNestedBlock{
				Description:         "List of IPs associated with the storage node.",
				MarkdownDescription: "List of IPs associated with the storage node.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required:            true,
							Description:         "IP address.",
							MarkdownDescription: "IP address.",
						},
						"role": schema.StringAttribute{
							Required:            true,
							Description:         "Role of the IP. Valid values: SdsOnly, SdcOnly, All.",
							MarkdownDescription: "Role of the IP. Valid values: `SdsOnly`, `SdcOnly`, `All`.",
						},
					},
				},
			},
		},
	}
}

func (r *storageNodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"The Gen2 API client is required for storage node management. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	r.client = p.genClient
}

// Create creates the resource and sets the initial Terraform state.
func (r *storageNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan storageNodeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build IP list
	var ipListItems []storageNodeIPModel
	diags = plan.IPList.ElementsAs(ctx, &ipListItems, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sdsIPList []clientgen.SdsIPList
	for _, item := range ipListItems {
		sdsIPList = append(sdsIPList, clientgen.SdsIPList{
			SdsIp: &clientgen.SdsIP{
				Ip:   clientgen.PtrString(item.IP.ValueString()),
				Role: clientgen.PtrString(item.Role.ValueString()),
			},
		})
	}

	// Build create params
	params := clientgen.StorageNodeParam{
		ProtectionDomainId: plan.ProtectionDomainID.ValueString(),
		IpList:             sdsIPList,
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		params.Name = clientgen.PtrString(plan.Name.ValueString())
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		params.Port = clientgen.PtrString(plan.Port.ValueString())
	}
	if !plan.FaultSetID.IsNull() && !plan.FaultSetID.IsUnknown() {
		params.FaultSetId = clientgen.PtrString(plan.FaultSetID.ValueString())
	}

	tflog.Debug(ctx, "Creating Storage Node")

	// Create the storage node
	result, _, err := r.client.StorageNodeAPI.CreateStorageNode(ctx).StorageNodeParam(params).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Storage Node",
			"Could not create storage node: "+err.Error(),
		)
		return
	}

	// Read back the created storage node to get full details
	nodeID := result.GetId()
	node, _, err := r.client.StorageNodeAPI.GetStorageNode(ctx, nodeID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storage Node After Creation",
			"Could not read storage node: "+err.Error(),
		)
		return
	}

	// Map response to state
	r.mapStorageNodeToState(ctx, node, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Created Storage Node with ID: %s", nodeID))
}

// Read refreshes the Terraform state with the latest data.
func (r *storageNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state storageNodeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	node, _, err := r.client.StorageNodeAPI.GetStorageNode(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storage Node",
			fmt.Sprintf("Could not read storage node ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	r.mapStorageNodeToState(ctx, node, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *storageNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state storageNodeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle name change
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != state.Name.ValueString() {
		namePayload := clientgen.SetStorageNodeName{
			Name: plan.Name.ValueString(),
		}
		_, err := r.client.StorageNodeAPI.SetStorageNodeName(ctx, state.ID.ValueString()).SetStorageNodeName(namePayload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Renaming Storage Node",
				fmt.Sprintf("Could not rename storage node: %s", err.Error()),
			)
			return
		}
	}

	// Read back the updated storage node
	node, _, err := r.client.StorageNodeAPI.GetStorageNode(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storage Node After Update",
			err.Error(),
		)
		return
	}

	r.mapStorageNodeToState(ctx, node, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *storageNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state storageNodeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.StorageNodeAPI.DeleteStorageNode(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Storage Node",
			fmt.Sprintf("Could not delete storage node ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Deleted Storage Node ID: %s", state.ID.ValueString()))
}

func (r *storageNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapStorageNodeToState maps the API response to the Terraform state model.
func (r *storageNodeResource) mapStorageNodeToState(_ context.Context, node *clientgen.StorageNode, model *storageNodeResourceModel) {
	if node == nil {
		return
	}

	if node.Id != nil {
		model.ID = types.StringValue(*node.Id)
	}
	if node.Name != nil {
		model.Name = types.StringValue(*node.Name)
	}
	if node.ProtectionDomainId != nil {
		model.ProtectionDomainID = types.StringValue(*node.ProtectionDomainId)
	}
	if node.FaultSetId != nil {
		model.FaultSetID = types.StringValue(*node.FaultSetId)
	}
	if node.Status != nil {
		model.Status = types.StringValue(*node.Status)
	}
	if node.State != nil {
		model.State = types.StringValue(*node.State)
	}
	if node.MembershipState != nil {
		model.MembershipState = types.StringValue(*node.MembershipState)
	}
	if node.MdmConnectionState != nil {
		model.MdmConnectionState = types.StringValue(*node.MdmConnectionState)
	}
	if node.SoftwareVersion != nil {
		model.SoftwareVersion = types.StringValue(*node.SoftwareVersion)
	}
	if node.Port != nil {
		model.Port = types.StringValue(fmt.Sprintf("%d", *node.Port))
	}
	if node.NumOfDevices != nil {
		model.NumOfDevices = types.Int64Value(int64(*node.NumOfDevices))
	}
	if node.TotalCapacityInKb != nil {
		model.TotalCapacityInKb = types.Int64Value(int64(*node.TotalCapacityInKb))
	}
	if node.UsedCapacityInKb != nil {
		model.UsedCapacityInKb = types.Int64Value(int64(*node.UsedCapacityInKb))
	}
	if node.MaintenanceState != nil {
		model.MaintenanceState = types.StringValue(*node.MaintenanceState)
	}
	if node.PerformanceProfile != nil {
		model.PerformanceProfile = types.StringValue(*node.PerformanceProfile)
	}
}
