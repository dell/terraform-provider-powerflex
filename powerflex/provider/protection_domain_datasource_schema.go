/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ProtectionDomainDataSourceSchema defines the schema for Protection Domain datasource
var ProtectionDomainDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing protection domain from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing protection domain from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier of the protection domain instance." +
				" Conflicts with 'name'.",
			MarkdownDescription: "Unique identifier of the protection domain instance." +
				" Conflicts with `name`.",
			Optional: true,
		},
		"name": schema.StringAttribute{
			Description: "Unique name of the protection domain instance." +
				" Conflicts with 'id'.",
			MarkdownDescription: "Unique name of the protection domain instance." +
				" Conflicts with `id`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("id")),
			},
		},
		"protection_domains": schema.ListNestedAttribute{
			Description:         "List of protection domains fetched.",
			MarkdownDescription: "List of protection domains fetched.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: protectionDomainDataAttributes,
			},
		},
	},
}

var protectionDomainDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "Unique identifier of the protection domain instance.",
		MarkdownDescription: "Unique identifier of the protection domain instance.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Unique name of the protection domain instance.",
		MarkdownDescription: "Unique name of the protection domain instance.",
		Computed:            true,
	},
	"state": schema.StringAttribute{
		Description:         "State of a PD. Valid values are Active, ActivePending, Inactive or InactivePending.",
		MarkdownDescription: "State of a PD. Valid values are `Active`, `ActivePending`, `Inactive` or `InactivePending`.",
		Computed:            true,
	},
	"system_id": schema.StringAttribute{
		Description:         "System ID of the PD.",
		MarkdownDescription: "System ID of the PD.",
		Computed:            true,
	},
	"rf_cache_accp_id": schema.StringAttribute{
		Description:         "ID of the Rf Cache Acceleration Pool attached to the PD.",
		MarkdownDescription: "ID of the Rf Cache Acceleration Pool attached to the PD.",
		Computed:            true,
	},
	"rf_cache_enabled": schema.BoolAttribute{
		Description:         "Whether SDS Rf Cache is enabled or not.",
		MarkdownDescription: "Whether SDS Rf Cache is enabled or not.",
		Computed:            true,
	},
	"rf_cache_opertional_mode": schema.StringAttribute{
		Description:         "Operational Mode of the SDS RF Cache.",
		MarkdownDescription: "Operational Mode of the SDS RF Cache.",
		Computed:            true,
	},
	"rf_cache_page_size_kb": schema.Int64Attribute{
		Description:         "Page size of the SDS RF Cache in KB.",
		MarkdownDescription: "Page size of the SDS RF Cache in KB.",
		Computed:            true,
	},
	"rf_cache_max_io_size_kb": schema.Int64Attribute{
		Description:         "Maximum io of the SDS RF Cache in KB.",
		MarkdownDescription: "Maximum io of the SDS RF Cache in KB.",
		Computed:            true,
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
	},
	"fgl_default_metadata_cache_size": schema.Int64Attribute{
		Description:         "Fine Granularity Metadata Cache size.",
		MarkdownDescription: "Fine Granularity Metadata Cache size.",
		Computed:            true,
	},
	"protected_maintenance_mode_network_throttling_enabled": schema.BoolAttribute{
		Description:         "Whether network throttling is enabled for protected maintenance mode.",
		MarkdownDescription: "Whether network throttling is enabled for protected maintenance mode.",
		Computed:            true,
	},
	"protected_maintenance_mode_network_throttling_in_kbps": schema.Int64Attribute{
		Description:         "Maximum allowed io for protected maintenance mode in KBps.",
		MarkdownDescription: "Maximum allowed io for protected maintenance mode in KBps.",
		Computed:            true,
	},
	"rebuild_network_throttling_enabled": schema.BoolAttribute{
		Description:         "Whether network throttling is enabled for rebuilding.",
		MarkdownDescription: "Whether network throttling is enabled for rebuilding.",
		Computed:            true,
	},
	"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
		Description:         "Maximum allowed io for rebuilding in KBps.",
		MarkdownDescription: "Maximum allowed io for rebuilding in KBps.",
		Computed:            true,
	},
	"rebalance_network_throttling_enabled": schema.BoolAttribute{
		Description:         "Whether network throttling is enabled for rebalancing.",
		MarkdownDescription: "Whether network throttling is enabled for rebalancing.",
		Computed:            true,
	},
	"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
		Description:         "Maximum allowed io for rebalancing in KBps.",
		MarkdownDescription: "Maximum allowed io for rebalancing in KBps.",
		Computed:            true,
	},
	"vtree_migration_network_throttling_enabled": schema.BoolAttribute{
		Description:         "Whether network throttling is enabled for vtree migration.",
		MarkdownDescription: "Whether network throttling is enabled for vtree migration.",
		Computed:            true,
	},
	"vtree_migration_network_throttling_in_kbps": schema.Int64Attribute{
		Description:         "Maximum allowed io for vtree migration in KBps.",
		MarkdownDescription: "Maximum allowed io for vtree migration in KBps.",
		Computed:            true,
	},
	"overall_io_network_throttling_enabled": schema.BoolAttribute{
		Description:         "Whether network throttling is enabled for overall io.",
		MarkdownDescription: "Whether network throttling is enabled for overall io.",
		Computed:            true,
	},
	"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
		Description:         "Maximum allowed io for protected maintenance mode in KBps. Must be greater than any other network throttling parameter.",
		MarkdownDescription: "Maximum allowed io for protected maintenance mode in KBps. Must be greater than any other network throttling parameter.",
		Computed:            true,
	},
	"sdr_sds_connectivity": schema.SingleNestedAttribute{
		Description:         "SDR-SDS Connectivity information.",
		MarkdownDescription: "SDR-SDS Connectivity information.",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"client_server_conn_status": schema.StringAttribute{
				Description:         "Connectivity Status.",
				MarkdownDescription: "Connectivity Status.",
				Computed:            true,
			},
			"disconnected_client_id": schema.StringAttribute{
				Description:         "ID of the disconnected client.",
				MarkdownDescription: "ID of the disconnected client.",
				Computed:            true,
			},
			"disconnected_client_name": schema.StringAttribute{
				Description:         "Name of the disconnected client.",
				MarkdownDescription: "Name of the disconnected client.",
				Computed:            true,
			},
			"disconnected_server_id": schema.StringAttribute{
				Description:         "ID of the disconnected server.",
				MarkdownDescription: "ID of the disconnected server.",
				Computed:            true,
			},
			"disconnected_server_name": schema.StringAttribute{
				Description:         "Name of the disconnected server.",
				MarkdownDescription: "Name of the disconnected server.",
				Computed:            true,
			},
			"disconnected_server_ip": schema.StringAttribute{
				Description:         "IP address of the disconnected server.",
				MarkdownDescription: "IP address of the disconnected server.",
				Computed:            true,
			},
		},
	},
	"sds_decoupled_counter": schema.SingleNestedAttribute{
		Description:         "SDS Decoupled Counter windows.",
		MarkdownDescription: "SDS Decoupled Counter windows.",
		Computed:            true,
		Attributes:          getAllWindowParamsSchema(),
	},
	"sds_configuration_failure_counter": schema.SingleNestedAttribute{
		Description:         "SDS Configuration Failure Counter windows.",
		MarkdownDescription: "SDS Configuration Failure Counter windows.",
		Computed:            true,
		Attributes:          getAllWindowParamsSchema(),
	},
	"mdm_sds_network_disconnections_counter": schema.SingleNestedAttribute{
		Description:         "MDM-SDS Network Disconnection Counter windows.",
		MarkdownDescription: "MDM-SDS Network Disconnection Counter windows.",
		Computed:            true,
		Attributes:          getAllWindowParamsSchema(),
	},
	"sds_sds_network_disconnections_counter": schema.SingleNestedAttribute{
		Description:         "SDS-SDS Network Disconnection Counter windows.",
		MarkdownDescription: "SDS-SDS Network Disconnection Counter windows.",
		Computed:            true,
		Attributes:          getAllWindowParamsSchema(),
	},
	"sds_receive_buffer_allocation_failures_counter": schema.SingleNestedAttribute{
		Description:         "SDS receive Buffer Allocation Failure Counter windows.",
		MarkdownDescription: "SDS receive Buffer Allocation Failure Counter windows.",
		Computed:            true,
		Attributes:          getAllWindowParamsSchema(),
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
}

func getAllWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			Description:         "Short Window Parameters.",
			MarkdownDescription: "Short Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			Description:         "Medium Window Parameters.",
			MarkdownDescription: "Medium Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			Description:         "Long Window Parameters.",
			MarkdownDescription: "Long Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsSchema(),
		},
	}
}

func getWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			Description:         "Threshold.",
			MarkdownDescription: "Threshold.",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			Description:         "Window Size in seconds.",
			MarkdownDescription: "Window Size in seconds.",
			Computed:            true,
		},
	}
}
