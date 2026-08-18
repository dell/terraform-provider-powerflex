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
	_ datasource.DataSource              = &deviceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceGroupDataSource{}
)

// DeviceGroupDataSource returns a new device group datasource.
func DeviceGroupDataSource() datasource.DataSource {
	return &deviceGroupDataSource{}
}

type deviceGroupDataSource struct {
	client *clientgen.APIClient
}

type deviceGroupDataSourceModel struct {
	ID                 types.String            `tfsdk:"id"`
	Name               types.String            `tfsdk:"name"`
	ProtectionDomainID types.String            `tfsdk:"protection_domain_id"`
	DeviceGroups       []deviceGroupItemModel  `tfsdk:"device_groups"`
}

type deviceGroupItemModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	NumOfDevices       types.Int64  `tfsdk:"num_of_devices"`
	UsableCapacityInKb types.Int64  `tfsdk:"usable_capacity_in_kb"`
	UsedCapacityInKb   types.Int64  `tfsdk:"used_capacity_in_kb"`
	Status             types.String `tfsdk:"status"`
}

func (d *deviceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

func (d *deviceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This datasource is used to query Device Groups from a PowerFlex Gen2 system (5.0+).",
		MarkdownDescription: "This datasource is used to query Device Groups from a PowerFlex Gen2 system (5.0+).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Device Group ID to query. Mutually exclusive with name.",
				MarkdownDescription: "Device Group ID to query. Mutually exclusive with `name`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Device Group name to query. Mutually exclusive with id.",
				MarkdownDescription: "Device Group name to query. Mutually exclusive with `id`.",
			},
			"protection_domain_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter device groups by protection domain ID.",
				MarkdownDescription: "Filter device groups by protection domain ID.",
			},
			"device_groups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of device groups.",
				MarkdownDescription: "List of device groups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Device Group ID.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Device Group name.",
						},
						"protection_domain_id": schema.StringAttribute{
							Computed:    true,
							Description: "Protection Domain ID.",
						},
						"num_of_devices": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices in the group.",
						},
						"usable_capacity_in_kb": schema.Int64Attribute{
							Computed:    true,
							Description: "Usable capacity in KB.",
						},
						"used_capacity_in_kb": schema.Int64Attribute{
							Computed:    true,
							Description: "Used capacity in KB.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the device group.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
			"The Gen2 API client is required for device group queries. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	d.client = p.genClient
}

func (d *deviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deviceGroupDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If ID is specified, query single group
	if !config.ID.IsNull() && config.ID.ValueString() != "" {
		group, _, err := d.client.DeviceGroupAPI.GetDeviceGroup(ctx, config.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Device Group",
				fmt.Sprintf("Could not read device group ID %s: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}

		item := mapDeviceGroupToItem(group)
		config.DeviceGroups = []deviceGroupItemModel{item}
		config.ID = types.StringValue(config.ID.ValueString())

		diags = resp.State.Set(ctx, &config)
		resp.Diagnostics.Append(diags...)
		return
	}

	// List all groups
	groups, _, err := d.client.DeviceGroupAPI.ListDeviceGroups(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Device Groups",
			fmt.Sprintf("Could not list device groups: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d device groups", len(groups)))

	var filteredGroups []deviceGroupItemModel
	for _, group := range groups {
		item := mapDeviceGroupToItem(&group)

		// Filter by name if specified
		if !config.Name.IsNull() && config.Name.ValueString() != "" {
			if group.Name == nil || *group.Name != config.Name.ValueString() {
				continue
			}
		}

		// Filter by protection domain ID if specified
		if !config.ProtectionDomainID.IsNull() && config.ProtectionDomainID.ValueString() != "" {
			if group.ProtectionDomainId == nil || *group.ProtectionDomainId != config.ProtectionDomainID.ValueString() {
				continue
			}
		}

		filteredGroups = append(filteredGroups, item)
	}

	config.DeviceGroups = filteredGroups
	if config.ID.IsNull() || config.ID.ValueString() == "" {
		config.ID = types.StringValue("device-group-datasource")
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func mapDeviceGroupToItem(group *clientgen.DeviceGroup) deviceGroupItemModel {
	item := deviceGroupItemModel{}
	if group.Id != nil {
		item.ID = types.StringValue(*group.Id)
	}
	if group.Name != nil {
		item.Name = types.StringValue(*group.Name)
	} else {
		item.Name = types.StringNull()
	}
	if group.ProtectionDomainId != nil {
		item.ProtectionDomainID = types.StringValue(*group.ProtectionDomainId)
	}
	if group.NumOfDevices != nil {
		item.NumOfDevices = types.Int64Value(int64(*group.NumOfDevices))
	}
	if group.UsableCapacityInKb != nil {
		item.UsableCapacityInKb = types.Int64Value(int64(*group.UsableCapacityInKb))
	}
	if group.UsedCapacityInKb != nil {
		item.UsedCapacityInKb = types.Int64Value(int64(*group.UsedCapacityInKb))
	}
	if group.Status != nil {
		item.Status = types.StringValue(*group.Status)
	}
	return item
}
