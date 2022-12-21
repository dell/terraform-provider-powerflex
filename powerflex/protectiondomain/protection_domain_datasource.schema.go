package protectiondomain

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// DataSourceSchema defines the schema for Protection Domain datasource
var DataSourceSchema schema.Schema = schema.Schema{
	Description: ".",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"protection_domains": schema.ListNestedAttribute{
			Description: "List of volumes",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"state": schema.StringAttribute{
						Description: "A PD can be Active, ActivePending, Inactive or InactivePending",
						Computed:    true,
					},
					"system_id": schema.StringAttribute{
						Description: "System ID of the PD",
						Computed:    true,
					},
					"rf_cache_accp_id": schema.StringAttribute{
						Description: "Rf Cache  Acceleration Pool ID",
						Computed:    true,
					},
					"rf_cache_enabled": schema.BoolAttribute{
						Description: "Whether SDS Rf Cache is enabled or not",
						Computed:    true,
					},
					"rf_cache_opertional_mode": schema.StringAttribute{
						Description: "Operational Mode of the SDS RF Cache",
						Computed:    true,
					},
					"rf_cache_page_size_kb": schema.Int64Attribute{
						Description: "Page size of the SDS RF Cache in KB",
						Computed:    true,
					},
					"rf_cache_max_io_size_kb": schema.Int64Attribute{
						Description: "Maximum io of the SDS RF Cache in KB",
						Computed:    true,
					},
					"fgl_default_num_concurrent_writes": schema.Int64Attribute{
						Description: "Fine Granularity default number of concurrent writes",
						Computed:    true,
					},
					"fgl_metadata_cache_enabled": schema.BoolAttribute{
						Description: "Whether Fine Granularity Metadata Cache is enabled or not",
						Computed:    true,
					},
					"fgl_default_metadata_cache_size": schema.Int64Attribute{
						Description: "Fine Granularity Metadata Cache size",
						Computed:    true,
					},
					"protected_maintenance_mode_network_throttling_enabled": schema.BoolAttribute{
						Description: "Whether network throttling is enabled for protected maintenance mode",
						Computed:    true,
					},
					"protected_maintenance_mode_network_throttling_in_kbps": schema.Int64Attribute{
						Description: "Maximum allowed io for protected maintenance mode in KBps",
						Computed:    true,
					},
					"rebuild_network_throttling_enabled": schema.BoolAttribute{
						Description: "Whether network throttling is enabled for rebuilding",
						Computed:    true,
					},
					"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
						Description: "Maximum allowed io for rebuilding in KBps",
						Computed:    true,
					},
					"rebalance_network_throttling_enabled": schema.BoolAttribute{
						Description: "Whether network throttling is enabled for rebalancing",
						Computed:    true,
					},
					"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
						Description: "Maximum allowed io for rebalancing in KBps",
						Computed:    true,
					},
					"vtree_migration_network_throttling_enabled": schema.BoolAttribute{
						Description: "Whether network throttling is enabled for vtree migration",
						Computed:    true,
					},
					"vtree_migration_network_throttling_in_kbps": schema.Int64Attribute{
						Description: "Maximum allowed io for vtree migration in KBps",
						Computed:    true,
					},
					"overall_io_network_throttling_enabled": schema.BoolAttribute{
						Description: "Whether network throttling is enabled for overall io",
						Computed:    true,
					},
					"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
						Description: "Maximum allowed io for protected maintenance mode in KBps. Must be greater than any other network throttling parameter.",
						Computed:    true,
					},
					"sdr_sds_connectivity": schema.SingleNestedAttribute{
						Description: "SDR-SDS Connectivity information",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"client_server_conn_status": schema.StringAttribute{
								Description: "Connectivity Status",
								Computed:    true,
							},
							"disconnected_client_id": schema.StringAttribute{
								Description: "ID of the disconnected client",
								Computed:    true,
							},
							"disconnected_client_name": schema.StringAttribute{
								Description: "Name of the disconnected client",
								Computed:    true,
							},
							"disconnected_server_id": schema.StringAttribute{
								Description: "ID of the disconnected server",
								Computed:    true,
							},
							"disconnected_server_name": schema.StringAttribute{
								Description: "Name of the disconnected server",
								Computed:    true,
							},
							"disconnected_server_ip": schema.StringAttribute{
								Description: "IP address of the disconnected server",
								Computed:    true,
							},
						},
					},
					"sds_decoupled_counter": schema.SingleNestedAttribute{
						Description: "SDS Decoupled Counter windows",
						Computed:    true,
						Attributes:  getAllWindowParamsSchema(),
					},
					"sds_configuration_failure_counter": schema.SingleNestedAttribute{
						Description: "SDS Configuration Failure Counter windows",
						Computed:    true,
						Attributes:  getAllWindowParamsSchema(),
					},
					"mdm_sds_network_disconnections_counter": schema.SingleNestedAttribute{
						Description: "MDM-SDS Network Disconnection Counter windows",
						Computed:    true,
						Attributes:  getAllWindowParamsSchema(),
					},
					"sds_sds_network_disconnections_counter": schema.SingleNestedAttribute{
						Description: "SDS-SDS Network Disconnection Counter windows",
						Computed:    true,
						Attributes:  getAllWindowParamsSchema(),
					},
					"sds_receive_buffer_allocation_failures_counter": schema.SingleNestedAttribute{
						Description: "SDS receive Buffer Allocation Failure Counter windows",
						Computed:    true,
						Attributes:  getAllWindowParamsSchema(),
					},
					"replication_capacity_max_ratio": schema.Int64Attribute{
						Description: "Maximum Replication Capacity Ratio",
						Computed:    true,
					},
					"links": schema.ListNestedAttribute{
						Description: "Underlying REST API links",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description: "Numeric identifier of the coffee ingredient.",
									Computed:    true,
								},
								"href": schema.StringAttribute{
									Description: "Numeric identifier of the coffee ingredient.",
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

func getAllWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			Description: "Short Window Params",
			Computed:    true,
			Attributes:  getWindowParamsSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			Description: "Medium Window Params",
			Computed:    true,
			Attributes:  getWindowParamsSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			Description: "Long Window Params",
			Computed:    true,
			Attributes:  getWindowParamsSchema(),
		},
	}
}

func getWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			Description: "Threshhold",
			Computed:    true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			Description: "Window Size in seconds",
			Computed:    true,
		},
	}
}
