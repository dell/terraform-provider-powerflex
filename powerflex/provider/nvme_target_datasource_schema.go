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

// NvmeTargetDataSourceSchema defines the schema for NvmeTarget datasource
var NvmeTargetDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing NVMe targets from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing NVMe targets from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the NVMe targets Datasource",
			MarkdownDescription: "ID of the NVMe targets Datasource",
			Computed:            true,
		},
		"nvme_target_details": schema.ListNestedAttribute{
			Description:         "List of NVMe targets",
			MarkdownDescription: "List of NVMe targets",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: NvmeTargetModelSchema,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.NvmeTargetFilter{})),
		},
	},
}

// NvmeTargetModelSchema defines the schema for NVMe target model
var NvmeTargetModelSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "ID of the NVMe target",
		MarkdownDescription: "ID of the NVMe target",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the NVMe target",
		MarkdownDescription: "Name of the NVMe target",
		Computed:            true,
	},
	"system_id": schema.StringAttribute{
		Description:         "The ID of the system.",
		MarkdownDescription: "The ID of the system.",
		Computed:            true,
	},
	"protection_domain_id": schema.StringAttribute{
		Description:         "Protection Domain ID of the replicatio of the NVMe target.",
		MarkdownDescription: "Protection Domain ID of the replicatio of the NVMe target.",
		Computed:            true,
	},
	"ip_list": schema.ListNestedAttribute{
		Description:         "List of IPs associated with the NVMe target.",
		MarkdownDescription: "List of IPs associated with the NVMe target.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"ip": schema.StringAttribute{
					Description:         "NVMe Target IP.",
					MarkdownDescription: "NVMe Target IP.",
					Computed:            true,
				},
				"role": schema.StringAttribute{
					Description:         "NVMe Target IP role.",
					MarkdownDescription: "NVMe Target IP role.",
					Computed:            true,
				},
			},
		},
	},
	"storage_port": schema.Int64Attribute{
		Description:         "The storage port of the NVMe target.",
		MarkdownDescription: "The storage port of the NVMe target.",
		Computed:            true,
	},
	"nvme_port": schema.Int64Attribute{
		Description:         "The NVMe port of the NVMe target.",
		MarkdownDescription: "The NVMe port of the NVMe target.",
		Computed:            true,
	},
	"discovery_port": schema.Int64Attribute{
		Description:         "The discovery port of the NVMe target.",
		MarkdownDescription: "The discovery port of the NVMe target.",
		Computed:            true,
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
		Description:         "The maintenance state of the NVMe target.",
		MarkdownDescription: "The maintenance state of the NVMe target.",
		Computed:            true,
	},
	"authentication_error": schema.StringAttribute{
		Description:         "The authentication error of the NVMe target.",
		MarkdownDescription: "The authentication error of the NVMe target.",
		Computed:            true,
	},
	"certificate_info": schema.SingleNestedAttribute{
		Description:         "Certificate Information.",
		MarkdownDescription: "Certificate Information.",
		Computed:            true,
		Attributes:          getCertificateInfoSchema(),
	},
	"links": schema.ListNestedAttribute{
		Description:         "Specifies the links associated with NVMe target.",
		MarkdownDescription: "Specifies the links associated with NVMe target.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"rel": schema.StringAttribute{
					Description:         "Specifies the relationship with the NVMe target.",
					MarkdownDescription: "Specifies the relationship with the NVMe target.",
					Computed:            true,
				},
				"href": schema.StringAttribute{
					Description:         "Specifies the exact path to fetch the details.",
					MarkdownDescription: "Specifies the exact path to fetch the details.",
					Computed:            true,
				},
			},
		},
	},
	"host_list": schema.ListNestedAttribute{
		Description:         "Hosts attached to the NVMe target.",
		MarkdownDescription: "Hosts attached to the NVMe target.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"host_ip": schema.StringAttribute{
					Description:         "Host IP address.",
					MarkdownDescription: "Host IP address.",
					Computed:            true,
				},
				"is_connected": schema.BoolAttribute{
					Description:         "Specifies whether the host is connected to the NVMe target.",
					MarkdownDescription: "Specifies whether the host is connected to the NVMe target.",
					Computed:            true,
				},
				"host_name": schema.StringAttribute{
					Description:         "Host name.",
					MarkdownDescription: "Host name.",
					Computed:            true,
				},
				"host_id": schema.StringAttribute{
					Description:         "Host ID.",
					MarkdownDescription: "Host ID.",
					Computed:            true,
				},
				"sys_port_ip": schema.StringAttribute{
					Description:         "Specifies the target IP address of the NVMe controller.",
					MarkdownDescription: "Specifies the target IP address of the NVMe controller.",
					Computed:            true,
				},
			},
		},
	},
}
