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

package powerflex

import (
	types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ProtectionDomainResourceSchema defines the schema for Protection Domain resource
var ProtectionDomainResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage protection domains on a PowerFlex array.",
	MarkdownDescription: "This resource can be used to manage protection domains on a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the protection domain instance.",
			MarkdownDescription: "Unique identifier of the protection domain instance.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description:         "Unique name of the protection domain instance.",
			MarkdownDescription: "Unique name of the protection domain instance.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"active": schema.BoolAttribute{
			Description:         "Whether the PD should be in 'Active' state. Default value is 'true'.",
			MarkdownDescription: "Whether the PD should be in `Active` state. Default value is `true`.",
			Computed:            true,
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolDefault(true),
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"state": schema.StringAttribute{
			Description:         "State of the PD.",
			MarkdownDescription: "State of the PD.",
			Computed:            true,
		},
		"rf_cache_accp_id": schema.StringAttribute{
			Description:         "ID of the RF Cache Acceleration Pool attached to the PD.",
			MarkdownDescription: "ID of the RF Cache Acceleration Pool attached to the PD.",
			Computed:            true,
		},
		"rf_cache_enabled": schema.BoolAttribute{
			Description:         "Whether SDS Rf Cache is enabled or not. Default value is 'true'.",
			MarkdownDescription: "Whether SDS Rf Cache is enabled or not. Default value is `true`.",
			Computed:            true,
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolDefault(true),
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"rf_cache_operational_mode": schema.StringAttribute{
			Description: "Operational Mode of the SDS RF Cache." +
				" Accepted values are 'Read', 'Write', 'ReadAndWrite' and 'WriteMiss'." +
				" Can be set only when 'rf_cache_enabled' is set to 'true'.",
			MarkdownDescription: "Operational Mode of the SDS RF Cache." +
				" Accepted values are `Read`, `Write`, `ReadAndWrite` and `WriteMiss`." +
				" Can be set only when `rf_cache_enabled` is set to `true`.",
			Computed: true,
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(types.PDRCModeRead),
					string(types.PDRCModeWrite),
					string(types.PDRCModeReadAndWrite),
					string(types.PDRCModeWriteMiss),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"rf_cache_page_size_kb": schema.Int64Attribute{
			Description: "Page size of the SDS RF Cache in KB." +
				" Can be set only when 'rf_cache_enabled' is set to 'true'.",
			MarkdownDescription: "Page size of the SDS RF Cache in KB." +
				" Can be set only when `rf_cache_enabled` is set to `true`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rf_cache_max_io_size_kb": schema.Int64Attribute{
			Description: "Maximum IO of the SDS RF Cache in KB." +
				" Can be set only when 'rf_cache_enabled' is set to 'true'.",
			MarkdownDescription: "Maximum IO of the SDS RF Cache in KB." +
				" Can be set only when `rf_cache_enabled` is set to `true`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"fgl_default_num_concurrent_writes": schema.Int64Attribute{
			Description:         "Fine Granularity default number of concurrent writes.",
			MarkdownDescription: "Fine Granularity default number of concurrent writes.",
			Computed:            true,
		},
		"fgl_metadata_cache_enabled": schema.BoolAttribute{
			Description:         "Whether Fine Granularity Metadata Cache is enabled or not.",
			MarkdownDescription: "Whether Fine Granularity Metadata Cache is enabled or not.",
			Computed:            true,
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"fgl_default_metadata_cache_size": schema.Int64Attribute{
			Description: "Fine Granularity Metadata Cache size." +
				" Can be set only when 'fgl_metadata_cache_enabled' is set to 'true'.",
			MarkdownDescription: "Fine Granularity Metadata Cache size." +
				" Can be set only when `fgl_metadata_cache_enabled` is set to `true`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"protected_maintenance_mode_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for protected maintenance mode in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for protected maintenance mode in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for rebuilding in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for rebuilding in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for rebalancing in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for rebalancing in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"vtree_migration_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for vtree migration in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for vtree migration in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for protected maintenance mode in KBps. Must be greater than any other network throttling parameter." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for protected maintenance mode in KBps. Must be greater than any other network throttling parameter." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"replication_capacity_max_ratio": schema.Int64Attribute{
			Description:         "Maximum Replication Capacity Ratio.",
			MarkdownDescription: "Maximum Replication Capacity Ratio.",
			Computed:            true,
		},
		"links": schema.ListNestedAttribute{
			Description:         "Underlying REST API links.",
			MarkdownDescription: "Underlying REST API links.",
			Computed:            true,
			// TODO: there is a bug in listplanmodifier which is fixed in tpf-v1.2.0
			// https://github.com/hashicorp/terraform-plugin-framework/issues/644
			// PlanModifiers: []planmodifier.List{
			// 	listplanmodifier.UseStateForUnknown(),
			// },
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description:         "Specifies the relationship with the Protection Domain.",
						MarkdownDescription: "Specifies the relationship with the Protection Domain.",
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
	},
}
