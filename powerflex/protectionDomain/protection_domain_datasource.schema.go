package protectiondomaindatasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var ProtectionDomainDataSourceSchema schema.Schema = schema.Schema{
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
		"state": schema.StringAttribute{
			Description: "A PD can be Active, ActivePending, Inactive or InactivePending",
			Computed:    true,
		},
		"system_id": schema.StringAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rf_cache_accp_id": schema.StringAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rf_cache_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rf_cache_opertional_mode": schema.StringAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rf_cache_page_size_kb": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rf_cache_max_io_size_kb": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"fgl_default_num_concurrent_writes": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"fgl_metadata_cache_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"fgl_default_metadata_cache_size": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"protected_maintenance_mode_network_throttling_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"protected_maintenance_mode_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rebuild_network_throttling_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rebalance_network_throttling_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"vtree_migration_network_throttling_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"vtree_migration_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"overall_io_network_throttling_enabled": schema.BoolAttribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
		},
		"sdr_sds_connectivity": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"client_server_conn_status": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
				"disconnected_client_id	str": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
				"disconnected_client_name": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
				"disconnected_server_id": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
				"disconnected_server_name": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
				"disconnected_server_ip": schema.StringAttribute{
					Description: "System ID of the PD",
					Computed:    true,
				},
			},
		},
		"sds_decoupled_counter": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes:  getAllWindowParamsSchema(),
		},
		"sds_configuration_failure_counter": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes:  getAllWindowParamsSchema(),
		},
		"mdm_sds_network_disconnections_counter": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes:  getAllWindowParamsSchema(),
		},
		"sds_sds_network_disconnections_counter": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes:  getAllWindowParamsSchema(),
		},
		"sds_receive_buffer_allocation_failures_counter": schema.SingleNestedAttribute{
			Description: "System ID of the PD",
			Computed:    true,
			Attributes:  getAllWindowParamsSchema(),
		},
		"replicationCapacityMaxRatio": schema.Int64Attribute{
			Description: "System ID of the PD",
			Computed:    true,
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
		"longWindow": schema.SingleNestedAttribute{
			Description: "Long Window Params",
			Computed:    true,
			Attributes:  getWindowParamsSchema(),
		},
	}
}

func getWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			Computed: true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			Computed: true,
		},
	}
}
