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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ComplianceReportResourceGroupSchema defines the schema for resource group compliance report
var ComplianceReportResourceGroupSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the compliance report for Resource Group from PowerFlex array.",
	MarkdownDescription: "This datasource is used to query the compliance report for Resource Group from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Unique identifier Of The compliance report Datasource.",
			MarkdownDescription: "Unique identifier Of The compliance report Datasource.",
			Computed:            true,
		},
		"compliance_reports": schema.ListNestedAttribute{
			Description:         "List of compliance report.",
			MarkdownDescription: "List of compliance report.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: ComplianceReportAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.ComplianceReportFilterType{})),
		},
	},
}

// ComplianceReportAttributes defines the schema of compliance report attributes
var ComplianceReportAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"service_tag": schema.StringAttribute{
		Description: "The service tag of the resource.",
		Computed:    true,
	},
	"ip_address": schema.StringAttribute{
		Description: "The IP address of the resource.",
		Computed:    true,
	},
	"firmware_repository_name": schema.StringAttribute{
		Description: "The name of the firmware repository.",
		Computed:    true,
	},
	"compliant": schema.BoolAttribute{
		Description: "The compliance status of the resource.",
		Computed:    true,
	},
	"device_type": schema.StringAttribute{
		Description: "The type of the device.",
		Computed:    true,
	},
	"model": schema.StringAttribute{
		Description: "The model of the device.",
		Computed:    true,
	},
	"available": schema.BoolAttribute{
		Description: "The availability status of the device.",
		Computed:    true,
	},
	"managed_state": schema.StringAttribute{
		Description: "The managed state of the device.",
		Computed:    true,
	},
	"embedded_report": schema.BoolAttribute{
		Description: "The presence of an embedded report.",
		Computed:    true,
	},
	"device_state": schema.StringAttribute{
		Description: "The state of the device.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The unique identifier of the resource group.",
		Computed:    true,
	},
	"host_name": schema.StringAttribute{
		Description: "The hostname of the resource group.",
		Computed:    true,
	},
	"can_update": schema.BoolAttribute{
		Description: "The update capability of the resource group.",
		Computed:    true,
	},
	"firmware_compliance_report_components": schema.ListNestedAttribute{
		Description: "The list of firmware compliance report components.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: "The unique identifier of the component.",
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "The name of the component.",
					Computed:    true,
				},
				"current_version": schema.SingleNestedAttribute{
					Description: "The current version of the component.",
					Attributes:  ComplianceReportComponentVersionAttr,
					Computed:    true,
				},
				"target_version": schema.SingleNestedAttribute{
					Description: "The target version of the component.",
					Attributes:  ComplianceReportComponentVersionAttr,
					Computed:    true,
				},
				"vendor": schema.StringAttribute{
					Description: "The vendor of the component.",
					Computed:    true,
				},
				"operating_system": schema.StringAttribute{
					Description: "The operating system of the component.",
					Computed:    true,
				},
				"compliant": schema.BoolAttribute{
					Description: "The compliance status of the component.",
					Computed:    true,
				},
				"software": schema.BoolAttribute{
					Computed: true,
				},
				"rpm": schema.BoolAttribute{
					Computed: true,
				},
				"os_compatible": schema.BoolAttribute{
					Computed: true,
				},
			},
		},
	},
}

// ComplianceReportComponentVersionAttr defines the schema of compliance report component version
var ComplianceReportComponentVersionAttr = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique identifier of the version.",
		Computed:    true,
	},
	"firmware_name": schema.StringAttribute{
		Description: "The name of the firmware.",
		Computed:    true,
	},
	"firmware_type": schema.StringAttribute{
		Description: "The type of the firmware.",
		Computed:    true,
	},
	"firmware_version": schema.StringAttribute{
		Description: "The version of the firmware.",
		Computed:    true,
	},
	"firmware_last_update": schema.StringAttribute{
		Description: "The last update time of the firmware.",
		Computed:    true,
	},
	"firmware_level": schema.StringAttribute{
		Description: "The level of the firmware.",
		Computed:    true,
	},
}
