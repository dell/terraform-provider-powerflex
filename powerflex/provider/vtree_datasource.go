/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &vtreeDataSource{}
	_ datasource.DataSourceWithConfigure = &vtreeDataSource{}
)

// VTreeDataSource returns the VTree data source
func VTreeDataSource() datasource.DataSource {
	return &vtreeDataSource{}
}

type vtreeDataSource struct {
	client *goscaleio.Client
}

func (d *vtreeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vtree"
}

func (d *vtreeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VTreeDataSourceSchema
}

func (d *vtreeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

// Read refreshes the Terraform state with the latest data.
func (d *vtreeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started vtree data source read method")
	var (
		state models.VTreeDataSourceModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch VTree details for all the VTrees
	vTrees, err := helper.GetAllVTrees(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting vTree details", err.Error(),
		)
		return
	}

	if state.VTreeFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.VTreeFilter, vTrees)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering vtrees: %v please validate the filter", state.VTreeFilter), err.Error(),
			)
			return
		}
		filteredvTrees := []scaleiotypes.VTreeDetails{}
		for _, val := range filtered {
			filteredvTrees = append(filteredvTrees, val.(scaleiotypes.VTreeDetails))
		}
		vTrees = filteredvTrees
	}

	state.VTrees = helper.GetAllVTreeState(vTrees)
	state.ID = types.StringValue("vtree_datasource")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// VTreeDataSourceSchema defines the schema for VTree data source
var VTreeDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing vtrees from the PowerFlex array. The information fetched from this datasource can be used for getting the details.",
	MarkdownDescription: "This datasource is used to query the existing vtrees from the PowerFlex array. The information fetched from this datasource can be used for getting the details.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder identifier attribute.",
			MarkdownDescription: "Placeholder identifier attribute.",
			Computed:            true,
		},
		"vtree_details": schema.SetNestedAttribute{
			Description:         "VTree details",
			MarkdownDescription: "VTree details",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"storage_pool_id": schema.StringAttribute{
						MarkdownDescription: "Storage pool ID",
						Description:         "Storage pool ID",
						Computed:            true,
					},
					"data_layout": schema.StringAttribute{
						MarkdownDescription: "Data layout",
						Description:         "Data layout",
						Computed:            true,
					},
					"compression_method": schema.StringAttribute{
						MarkdownDescription: "Compression method",
						Description:         "Compression method",
						Computed:            true,
					},
					"root_volumes": schema.SetAttribute{
						MarkdownDescription: "Root volumes",
						Description:         "Root volumes",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"vtree_migration_info": schema.SingleNestedAttribute{
						MarkdownDescription: "Vtree migration information",
						Description:         "Vtree migration information",
						Computed:            true,
						Attributes:          VtreeMigrationInfoSchema(),
					},
					"in_deletion": schema.BoolAttribute{
						MarkdownDescription: "In deletion",
						Description:         "In deletion",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "VTree name",
						Description:         "VTree name",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "VTree ID",
						Description:         "VTree ID",
						Computed:            true,
					},
					"links": schema.ListNestedAttribute{
						MarkdownDescription: "Specifies the links associated with VTree",
						Description:         "Specifies the links associated with VTree",
						Computed:            true,
						NestedObject:        schema.NestedAttributeObject{Attributes: LinksSchema()},
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.VTreeFilter{})),
		},
	},
}

// VtreeMigrationInfoSchema specifies the schema for VTree migration
func VtreeMigrationInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"migration_queue_position": schema.Int64Attribute{
			MarkdownDescription: "Migration queue position",
			Description:         "Migration queue position",
			Computed:            true,
		},
		"migration_pause_reason": schema.StringAttribute{
			MarkdownDescription: "Migration pause reason",
			Description:         "Migration pause reason",
			Computed:            true,
		},
		"migration_status": schema.StringAttribute{
			MarkdownDescription: "Migration status",
			Description:         "Migration status",
			Computed:            true,
		},
		"source_storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "Source storage pool ID",
			Description:         "Source storage pool ID",
			Computed:            true,
		},
		"destination_storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "Destination storage pool ID",
			Description:         "Destination storage pool ID",
			Computed:            true,
		},
		"thickness_conversion_type": schema.StringAttribute{
			MarkdownDescription: "Thickness conversion type",
			Description:         "Thickness conversion type",
			Computed:            true,
		},
	}
}

// LinksSchema specifies the schema for VTree links
func LinksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"rel": schema.StringAttribute{
			MarkdownDescription: "Specifies the relationship with the VTree",
			Description:         "Specifies the relationship with the VTree",
			Computed:            true,
		},
		"href": schema.StringAttribute{
			MarkdownDescription: "Specifies the exact path to fetch the details",
			Description:         "Specifies the exact path to fetch the details",
			Computed:            true,
		},
	}
}
