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
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NvmeHostDataSourceSchema defines the schema for NvmeHost datasource
var NvmeHostDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing NVMe hosts from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing NVMe hosts from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the NVMe hosts Datasource",
			MarkdownDescription: "ID of the NVMe hosts Datasource",
			Computed:            true,
		},
		"nvme_host_details": schema.ListNestedAttribute{
			Description:         "List of NVMe hosts",
			MarkdownDescription: "List of NVMe hosts",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: NvmeHostModelSchema,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"ids": schema.SetAttribute{
					Description:         "List of NVMe host Ids.",
					MarkdownDescription: "List of NVMe host Ids.",
					ElementType:         types.StringType,
					Optional:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				"names": schema.SetAttribute{
					Description:         "List of NVMe host names.",
					MarkdownDescription: "List of NVMe host names.",
					ElementType:         types.StringType,
					Optional:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				"nqns": schema.SetAttribute{
					Description:         "List of NVMe host nqn.",
					MarkdownDescription: "List of NVMe host nqn.",
					ElementType:         types.StringType,
					Optional:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
			},
		},
	},
}

// NvmeHostModelSchema defines the schema for NVme host model
var NvmeHostModelSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "ID of the NVMe host",
		MarkdownDescription: "ID of the NVMe host",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the NVMe host",
		MarkdownDescription: "Name of the NVMe host",
		Computed:            true,
	},
	"system_id": schema.StringAttribute{
		Description:         "The ID of the system.",
		MarkdownDescription: "The ID of the system.",
		Computed:            true,
	},
	"nqn": schema.StringAttribute{
		Description:         "NQN of the NVMe host.",
		MarkdownDescription: "NQN of the NVMe host.",
		Computed:            true,
	},
	"max_num_paths": schema.Int64Attribute{
		Description:         "Number of Paths Per Volume.",
		MarkdownDescription: "Number of Paths Per Volume.",
		Computed:            true,
	},
	"max_num_sys_ports": schema.Int64Attribute{
		Description:         "Number of System Ports per Protection Domain.",
		MarkdownDescription: "Number of System Ports per Protection Domain.",
		Computed:            true,
	},
	"links": schema.ListNestedAttribute{
		Description:         "Specifies the links associated with NVMe host.",
		MarkdownDescription: "Specifies the links associated with NVMe host.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"rel": schema.StringAttribute{
					Description:         "Specifies the relationship with the NVMe host.",
					MarkdownDescription: "Specifies the relationship with the NVMe host.",
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
}
