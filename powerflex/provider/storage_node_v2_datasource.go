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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &storageNodeV2DataSource{}
	_ datasource.DataSourceWithConfigure = &storageNodeV2DataSource{}
)

// StorageNodeV2DataSource returns a new storage node v2 datasource.
func StorageNodeV2DataSource() datasource.DataSource {
	return &storageNodeV2DataSource{}
}

type storageNodeV2DataSource struct {
	client *clientgen.APIClient
}

// storageNodeV2DataSourceModel describes the datasource data model.
type storageNodeV2DataSourceModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	ProtectionDomainID types.String             `tfsdk:"protection_domain_id"`
	StorageNodes       []storageNodeV2ItemModel `tfsdk:"storage_nodes"`
}

// storageNodeV2ItemModel describes each storage node item in the list.
type storageNodeV2ItemModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	FaultSetID         types.String `tfsdk:"fault_set_id"`
	Status             types.String `tfsdk:"status"`
	State              types.String `tfsdk:"state"`
	MembershipState    types.String `tfsdk:"membership_state"`
	MdmConnectionState types.String `tfsdk:"mdm_connection_state"`
	SoftwareVersion    types.String `tfsdk:"software_version"`
	Port               types.Int64  `tfsdk:"port"`
	NumOfDevices       types.Int64  `tfsdk:"num_of_devices"`
	TotalCapacityInKb  types.Int64  `tfsdk:"total_capacity_in_kb"`
	UsedCapacityInKb   types.Int64  `tfsdk:"used_capacity_in_kb"`
	MaintenanceState   types.String `tfsdk:"maintenance_state"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
}

func (d *storageNodeV2DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_node_v2"
}

func (d *storageNodeV2DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This datasource is used to query Storage Nodes from a PowerFlex Gen2 system (5.0+). Storage Nodes are the Gen2 replacement for SDS.",
		MarkdownDescription: "This datasource is used to query Storage Nodes from a PowerFlex Gen2 system (5.0+). Storage Nodes are the Gen2 replacement for SDS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Storage Node ID to query. Mutually exclusive with name.",
				MarkdownDescription: "Storage Node ID to query. Mutually exclusive with `name`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Storage Node name to query. Mutually exclusive with id.",
				MarkdownDescription: "Storage Node name to query. Mutually exclusive with `id`.",
			},
			"protection_domain_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter storage nodes by protection domain ID.",
				MarkdownDescription: "Filter storage nodes by protection domain ID.",
			},
			"storage_nodes": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of storage nodes.",
				MarkdownDescription: "List of storage nodes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Storage Node ID.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Storage Node name.",
						},
						"protection_domain_id": schema.StringAttribute{
							Computed:    true,
							Description: "Protection Domain ID.",
						},
						"fault_set_id": schema.StringAttribute{
							Computed:    true,
							Description: "Fault Set ID.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the storage node.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "State of the storage node.",
						},
						"membership_state": schema.StringAttribute{
							Computed:    true,
							Description: "Membership state.",
						},
						"mdm_connection_state": schema.StringAttribute{
							Computed:    true,
							Description: "MDM connection state.",
						},
						"software_version": schema.StringAttribute{
							Computed:    true,
							Description: "Software version.",
						},
						"port": schema.Int64Attribute{
							Computed:    true,
							Description: "Port number.",
						},
						"num_of_devices": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices.",
						},
						"total_capacity_in_kb": schema.Int64Attribute{
							Computed:    true,
							Description: "Total capacity in KB.",
						},
						"used_capacity_in_kb": schema.Int64Attribute{
							Computed:    true,
							Description: "Used capacity in KB.",
						},
						"maintenance_state": schema.StringAttribute{
							Computed:    true,
							Description: "Maintenance state.",
						},
						"performance_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Performance profile.",
						},
					},
				},
			},
		},
	}
}

func (d *storageNodeV2DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*powerflexProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *powerflexProvider, got: %T.", req.ProviderData),
		)
		return
	}

	if p.genClient == nil {
		resp.Diagnostics.AddError(
			"Gen2 API Client Not Available",
			"The Gen2 API client is required for storage node queries. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	d.client = p.genClient
}

func (d *storageNodeV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config storageNodeV2DataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If ID is specified, query single node
	if !config.ID.IsNull() && config.ID.ValueString() != "" {
		node, _, err := d.client.StorageNodeAPI.GetStorageNode(ctx, config.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Storage Node",
				fmt.Sprintf("Could not read storage node ID %s: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}

		item := mapStorageNodeToItem(node)
		config.StorageNodes = []storageNodeV2ItemModel{item}
		config.ID = types.StringValue(config.ID.ValueString())

		diags = resp.State.Set(ctx, &config)
		resp.Diagnostics.Append(diags...)
		return
	}

	// List all nodes
	nodes, _, err := d.client.StorageNodeAPI.ListStorageNodes(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Storage Nodes",
			fmt.Sprintf("Could not list storage nodes: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d storage nodes", len(nodes)))

	var filteredNodes []storageNodeV2ItemModel
	for _, node := range nodes {
		item := mapStorageNodeToItem(&node)

		// Filter by name if specified
		if !config.Name.IsNull() && config.Name.ValueString() != "" {
			if node.Name == nil || *node.Name != config.Name.ValueString() {
				continue
			}
		}

		// Filter by protection domain ID if specified
		if !config.ProtectionDomainID.IsNull() && config.ProtectionDomainID.ValueString() != "" {
			if node.ProtectionDomainId == nil || *node.ProtectionDomainId != config.ProtectionDomainID.ValueString() {
				continue
			}
		}

		filteredNodes = append(filteredNodes, item)
	}

	config.StorageNodes = filteredNodes
	if config.ID.IsNull() || config.ID.ValueString() == "" {
		config.ID = types.StringValue("storage-node-datasource")
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func mapStorageNodeToItem(node *clientgen.StorageNode) storageNodeV2ItemModel {
	item := storageNodeV2ItemModel{}
	if node.Id != nil {
		item.ID = types.StringValue(*node.Id)
	}
	if node.Name != nil {
		item.Name = types.StringValue(*node.Name)
	}
	if node.ProtectionDomainId != nil {
		item.ProtectionDomainID = types.StringValue(*node.ProtectionDomainId)
	}
	if node.FaultSetId != nil {
		item.FaultSetID = types.StringValue(*node.FaultSetId)
	} else {
		item.FaultSetID = types.StringNull()
	}
	if node.Status != nil {
		item.Status = types.StringValue(*node.Status)
	}
	if node.State != nil {
		item.State = types.StringValue(*node.State)
	}
	if node.MembershipState != nil {
		item.MembershipState = types.StringValue(*node.MembershipState)
	}
	if node.MdmConnectionState != nil {
		item.MdmConnectionState = types.StringValue(*node.MdmConnectionState)
	}
	if node.SoftwareVersion != nil {
		item.SoftwareVersion = types.StringValue(*node.SoftwareVersion)
	}
	if node.Port != nil {
		item.Port = types.Int64Value(int64(*node.Port))
	}
	if node.NumOfDevices != nil {
		item.NumOfDevices = types.Int64Value(int64(*node.NumOfDevices))
	}
	if node.TotalCapacityInKb != nil {
		item.TotalCapacityInKb = types.Int64Value(int64(*node.TotalCapacityInKb))
	}
	if node.UsedCapacityInKb != nil {
		item.UsedCapacityInKb = types.Int64Value(int64(*node.UsedCapacityInKb))
	}
	if node.MaintenanceState != nil {
		item.MaintenanceState = types.StringValue(*node.MaintenanceState)
	}
	if node.PerformanceProfile != nil {
		item.PerformanceProfile = types.StringValue(*node.PerformanceProfile)
	}
	return item
}
