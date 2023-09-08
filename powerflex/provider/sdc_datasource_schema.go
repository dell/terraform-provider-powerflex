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
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing Storage Data Clients from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing Storage Data Clients from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "ID of the SDC to fetch." +
				" Conflicts with 'name'",
			MarkdownDescription: "ID of the SDC to fetch." +
				" Conflicts with `name`",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("name")),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"name": schema.StringAttribute{
			Description: "Name of the SDC to fetch." +
				" Conflicts with 'id'",
			MarkdownDescription: "Name of the SDC to fetch." +
				" Conflicts with `id`",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"sdcs": schema.ListNestedAttribute{
			Description: "List of fetched SDCs.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.ID,
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.Name,
						Computed:    true,
					},
					"sdc_guid": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.SdcGUID,
						Computed:    true,
					},
					"on_vmware": schema.BoolAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.OnVMWare,
						Computed:    true,
					},
					"sdc_approved": schema.BoolAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.SdcApproved,
						Computed:    true,
					},
					"system_id": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.SystemID,
						Computed:    true,
					},
					"sdc_ip": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.SdcIP,
						Computed:    true,
					},
					"mdm_connection_state": schema.StringAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.MdmConnectionState,
						Computed:    true,
					},
					"links": schema.ListNestedAttribute{
						Description: models.SdcDatasourceSchemaDescriptions.Links,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description: models.SdcDatasourceSchemaDescriptions.LinksRel,
									Computed:    true,
								},
								"href": schema.StringAttribute{
									Description: models.SdcDatasourceSchemaDescriptions.LinksHref,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	},
}
