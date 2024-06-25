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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NodeDataSourceSchema defines the schema for node datasource
var NodeDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing nodes from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing nodes from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"node_ids": schema.SetAttribute{
			Description:         "List of node IDs",
			MarkdownDescription: "List of node IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("service_tags"),
					path.MatchRoot("ip_addresses"),
					path.MatchRoot("node_pool_ids"),
					path.MatchRoot("node_pool_names"),
				),
			},
		},
		"service_tags": schema.SetAttribute{
			Description:         "List of node service tags",
			MarkdownDescription: "List of node service tags",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("ip_addresses"),
					path.MatchRoot("node_pool_ids"),
					path.MatchRoot("node_pool_names"),
				),
			},
		},
		"ip_addresses": schema.SetAttribute{
			Description:         "List of node IP addresses",
			MarkdownDescription: "List of node IP addresses",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("node_pool_ids"),
					path.MatchRoot("node_pool_names"),
				),
			},
		},
		"node_pool_ids": schema.SetAttribute{
			Description:         "List of node pool IDs",
			MarkdownDescription: "List of node pool IDs",
			Optional:            true,
			ElementType:         types.Int64Type,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"node_pool_names": schema.SetAttribute{
			Description:         "List of node pool names",
			MarkdownDescription: "List of node pool names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("node_pool_ids"),
				),
			},
		},
		"node_details": schema.SetNestedAttribute{
			Description:         "Node details",
			MarkdownDescription: "Node details",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ref_id": schema.StringAttribute{
						Description:         "Reference ID of the node.",
						MarkdownDescription: "Reference ID of the node.",
						Computed:            true,
					},
					"ip_address": schema.StringAttribute{
						Description:         "IP address of the node.",
						MarkdownDescription: "IP address of the node.",
						Computed:            true,
					},
					"current_ip_address": schema.StringAttribute{
						Description:         "Current IP address of the node.",
						MarkdownDescription: "Current IP address of the node.",
						Computed:            true,
					},
					"service_tag": schema.StringAttribute{
						Description:         "Service tag of the node.",
						MarkdownDescription: "Service tag of the node.",
						Computed:            true,
					},
					"model": schema.StringAttribute{
						Description:         "Model of the node.",
						MarkdownDescription: "Model of the node.",
						Computed:            true,
					},
					"device_type": schema.StringAttribute{
						Description:         "Device type of the node.",
						MarkdownDescription: "Device type of the node.",
						Computed:            true,
					},
					"discover_device_type": schema.StringAttribute{
						Description:         "Discover device type of the node.",
						MarkdownDescription: "Discover device type of the node.",
						Computed:            true,
					},
					"display_name": schema.StringAttribute{
						Description:         "Display name of the node.",
						MarkdownDescription: "Display name of the node.",
						Computed:            true,
					},
					"managed_state": schema.StringAttribute{
						Description:         "Managed state of the node.",
						MarkdownDescription: "Managed state of the node.",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						Description:         "State of the node.",
						MarkdownDescription: "State of the node.",
						Computed:            true,
					},
					"in_use": schema.BoolAttribute{
						Description:         "Flag specifying if node is in use.",
						MarkdownDescription: "Flag specifying if node is in use.",
						Computed:            true,
					},
					"custom_firmware": schema.BoolAttribute{
						Description:         "Custom firmware of the node.",
						MarkdownDescription: "Custom firmware of the node.",
						Computed:            true,
					},
					"needs_attention": schema.BoolAttribute{
						Description:         "Flag specifying if node needs attention.",
						MarkdownDescription: "Flag specifying if node needs attention.",
						Computed:            true,
					},
					"manufacturer": schema.StringAttribute{
						Description:         "Manufacturer of the node.",
						MarkdownDescription: "Manufacturer of the node.",
						Computed:            true,
					},
					"system_id": schema.StringAttribute{
						Description:         "System ID.",
						MarkdownDescription: "System ID.",
						Computed:            true,
					},
					"health": schema.StringAttribute{
						Description:         "Health of the node.",
						MarkdownDescription: "Health of the node.",
						Computed:            true,
					},
					"health_message": schema.StringAttribute{
						Description:         "Health message.",
						MarkdownDescription: "Health message.",
						Computed:            true,
					},
					"operating_system": schema.StringAttribute{
						Description:         "Operating system of the node.",
						MarkdownDescription: "Operating system of the node.",
						Computed:            true,
					},
					"number_of_cpus": schema.Int64Attribute{
						Description:         "Number of CPUs of the node.",
						MarkdownDescription: "Number of CPUs of the node.",
						Computed:            true,
					},
					"nics": schema.Int64Attribute{
						Description:         "NICs of the node.",
						MarkdownDescription: "NICs of the node.",
						Computed:            true,
					},
					"memory_in_gb": schema.Int64Attribute{
						Description:         "Memory in GB.",
						MarkdownDescription: "Memory in GB.",
						Computed:            true,
					},
					"compliance_check_date": schema.StringAttribute{
						Description:         "Compliance check date.",
						MarkdownDescription: "Compliance check date.",
						Computed:            true,
					},
					"discovered_date": schema.StringAttribute{
						Description:         "Discovered date of the node.",
						MarkdownDescription: "Discovered date of the node.",
						Computed:            true,
					},
					"cred_id": schema.StringAttribute{
						Description:         "Cred ID.",
						MarkdownDescription: "Cred ID.",
						Computed:            true,
					},
					"compliance": schema.StringAttribute{
						Description:         "Node compliance.",
						MarkdownDescription: "Node compliance.",
						Computed:            true,
					},
					"failures_count": schema.Int64Attribute{
						Description:         "Failures count.",
						MarkdownDescription: "Failures count.",
						Computed:            true,
					},
					"facts": schema.StringAttribute{
						Description:         "Facts of the node.",
						MarkdownDescription: "Facts of the node.",
						Computed:            true,
					},
					"puppet_cert_name": schema.StringAttribute{
						Description:         "Puppet cert name of the node.",
						MarkdownDescription: "Puppet cert name of the node.",
						Computed:            true,
					},
					"flex_os_maint_mode": schema.Int64Attribute{
						Description:         "FLEX OS maintenance mode.",
						MarkdownDescription: "FLEX OS maintenance mode.",
						Computed:            true,
					},
					"esxi_maint_mode": schema.Int64Attribute{
						Description:         "ESXi maintenance mode.",
						MarkdownDescription: "ESXi maintenance mode.",
						Computed:            true,
					},
					"device_group_list": schema.SingleNestedAttribute{
						Description:         "Device group list.",
						MarkdownDescription: "Device group list.",
						Computed:            true,
						Attributes:          getDeviceGroupListSchema(),
					},
				},
			},
		},
	},
}

func getDeviceGroupListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"device_group": schema.ListNestedAttribute{
			Description:         "Device group information.",
			MarkdownDescription: "Device group information.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getDeviceGroupSchema(),
			},
		},
	}
}

func getDeviceGroupSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"group_seq_id": schema.Int64Attribute{
			Description:         "Group Sequence ID.",
			MarkdownDescription: "Group Sequence ID.",
			Computed:            true,
		},
		"group_name": schema.StringAttribute{
			Description:         "Group name.",
			MarkdownDescription: "Group name.",
			Computed:            true,
		},
		"group_description": schema.StringAttribute{
			Description:         "Group description.",
			MarkdownDescription: "Group description.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			Description:         "Creation date.",
			MarkdownDescription: "Creation date.",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			Description:         "User who created the group.",
			MarkdownDescription: "User who created the group.",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			Description:         "Updated date.",
			MarkdownDescription: "Updated date.",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			Description:         "User who updated the group.",
			MarkdownDescription: "User who updated the group.",
			Computed:            true,
		},
		"group_user_list": schema.SingleNestedAttribute{
			Description:         "Group user list.",
			MarkdownDescription: "Group user list.",
			Computed:            true,
			Attributes:          getGroupUserListSchema(),
		},
	}
}

func getGroupUserListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"total_records": schema.Int64Attribute{
			Description:         "Total number of records.",
			MarkdownDescription: "Total number of records.",
			Computed:            true,
		},
		"group_users": schema.ListNestedAttribute{
			Description:         "Group user information.",
			MarkdownDescription: "Group user information.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getGroupUserSchema(),
			},
		},
	}
}

func getGroupUserSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"user_seq_id": schema.Int64Attribute{
			Description:         "User sequence ID.",
			MarkdownDescription: "User sequence ID.",
			Computed:            true,
		},
		"user_name": schema.StringAttribute{
			Description:         "User name.",
			MarkdownDescription: "User name.",
			Computed:            true,
		},
		"first_name": schema.StringAttribute{
			Description:         "First name of the user.",
			MarkdownDescription: "First name of the user.",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			Description:         "Last name of the user.",
			MarkdownDescription: "Last name of the user.",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			Description:         "Enabled flag.",
			MarkdownDescription: "Enabled flag.",
			Computed:            true,
		},
	}
}
