package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceSchema is the schema for reading the storage pool data
var DataSourceSchema schema.Schema = schema.Schema{
	Description: "This data-source can be used to fetch information related to the storage pools from a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder identifier attribute.",
			MarkdownDescription: "Placeholder identifier attribute.",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "ID of the Protection Domain from which storage pools will be fetched." +
				" Must be provided if and only if 'protection_domain_name' is not provided.",
			MarkdownDescription: "ID of the Protection Domain from which storage pools will be fetched." +
				" Must be provided if and only if `protection_domain_name` is not provided.",
			Optional: true,
		},
		"protection_domain_name": schema.StringAttribute{
			Description: "Name of the Protection Domain from which storage pools will be fetched." +
				" Must be provided if and only if 'protection_domain_id' is not provided.",
			MarkdownDescription: "Name of the Protection Domain from which storage pools will be fetched." +
				" Must be provided if and only if `protection_domain_id` is not provided.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"storage_pool_ids": schema.ListAttribute{
			Description: "List of storage pool IDs." +
				" Can be provided only if 'storage_pool_names' is not provided.",
			MarkdownDescription: "List of storage pool IDs." +
				" Can be provided only if `storage_pool_names` is not provided.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"storage_pool_names": schema.ListAttribute{
			Description: "List of storage pool names." +
				" Can be provided only if 'storage_pool_ids' is not provided.",
			MarkdownDescription: "List of storage pool names." +
				" Can be provided only if `storage_pool_ids` is not provided.",
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRoot("storage_pool_ids")),
			},
		},
		"storage_pools": schema.ListNestedAttribute{
			Description:         "List of fetched storage pools.",
			MarkdownDescription: "List of fetched storage pools.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "Storage pool ID.",
						MarkdownDescription: "Storage pool ID.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Storage pool name.",
						MarkdownDescription: "Storage pool name.",
						Computed:            true,
					},
					"rebalance_io_priority_policy": schema.StringAttribute{
						Description:         "Rebalance IO Priority Policy.",
						MarkdownDescription: "Rebalance IO Priority Policy.",
						Computed:            true,
					},
					"rebuild_io_priority_policy": schema.StringAttribute{
						Description:         "Rebuild IO Priority Policy.",
						MarkdownDescription: "Rebuild IO Priority Policy.",
						Computed:            true,
					},
					"rebuild_io_priority_bw_limit_per_device_in_kbps": schema.Int64Attribute{
						Description:         "Rebuild Bandwidth Limit per Device.",
						MarkdownDescription: "Rebuild Bandwidth Limit per Device.",
						Computed:            true,
					},
					"rebuild_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
						Description:         "Number of Concurrent Rebuild IOPS per Device.",
						MarkdownDescription: "Number of Concurrent Rebuild IOPS per Device.",
						Computed:            true,
					},
					"rebalance_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
						Description:         "Number of Concurrent Rebalance IOPS per Device.",
						MarkdownDescription: "Number of Concurrent Rebalance IOPS per Device.",
						Computed:            true,
					},
					"rebalance_io_priority_bw_limit_per_device_kbps": schema.Int64Attribute{
						Description:         "Rebalance Bandwidth Limit per Device.",
						MarkdownDescription: "Rebalance Bandwidth Limit per Device.",
						Computed:            true,
					},
					"rebuild_io_priority_app_iops_per_device_threshold": schema.Int64Attribute{
						Description:         "Rebuild Application IOPS per Device Threshold.",
						MarkdownDescription: "Rebuild Application IOPS per Device Threshold.",
						Computed:            true,
					},
					"rebalance_io_priority_app_iops_per_device_threshold": schema.Int64Attribute{
						Description:         "Rebalance Application IOPS per Device Threshold.",
						MarkdownDescription: "Rebalance Application IOPS per Device Threshold.",
						Computed:            true,
					},
					"rebuild_io_priority_app_bw_per_device_threshold_kbps": schema.Int64Attribute{
						Description:         "Rebuild Application Bandwidth per Device Threshold.",
						MarkdownDescription: "Rebuild Application Bandwidth per Device Threshold.",
						Computed:            true,
					},
					"rebalance_io_priority_app_bw_per_device_threshold_kbps": schema.Int64Attribute{
						Description:         "Rebalance Application Bandwidth per Device Threshold.",
						MarkdownDescription: "Rebalance Application Bandwidth per Device Threshold.",
						Computed:            true,
					},
					"rebuild_io_priority_quiet_period_msec": schema.Int64Attribute{
						Description:         "Rebuild Quiet Period.",
						MarkdownDescription: "Rebuild Quiet Period.",
						Computed:            true,
					},
					"rebalance_io_priority_quiet_period_msec": schema.Int64Attribute{
						Description:         "Rebalance Quiet Period.",
						MarkdownDescription: "Rebalance Quiet Period.",
						Computed:            true,
					},
					"zero_padding_enabled": schema.BoolAttribute{
						Description:         "Zero Padding Enabled.",
						MarkdownDescription: "Zero Padding Enabled.",
						Computed:            true,
					},
					"use_rm_cache": schema.BoolAttribute{
						Description:         "Use RAM Read Cache.",
						MarkdownDescription: "Use RAM Read Cache.",
						Computed:            true,
					},
					"spare_percentage": schema.Int64Attribute{
						Description:         "Spare Percentage.",
						MarkdownDescription: "Spare Percentage.",
						Computed:            true,
					},
					"rm_cache_write_handling_mode": schema.StringAttribute{
						Description:         "RAM Read Cache Write Handling Mode.",
						MarkdownDescription: "RAM Read Cache Write Handling Mode.",
						Computed:            true,
					},
					"rebuild_enabled": schema.BoolAttribute{
						Description:         "Rebuild Enabled.",
						MarkdownDescription: "Rebuild Enabled.",
						Computed:            true,
					},
					"rebalance_enabled": schema.BoolAttribute{
						Description:         "Rebalance Enabled.",
						MarkdownDescription: "Rebalance Enabled.",
						Computed:            true,
					},
					"num_of_parallel_rebuild_rebalance_jobs_per_device": schema.Int64Attribute{
						Description:         "Number of Parallel Rebuild/Rebalance Jobs per Device.",
						MarkdownDescription: "Number of Parallel Rebuild/Rebalance Jobs per Device.",
						Computed:            true,
					},
					"background_scanner_bw_limit_kbps": schema.Int64Attribute{
						Description:         "Background Scanner Bandwidth Limit.",
						MarkdownDescription: "Background Scanner Bandwidth Limit.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
						Description:         "Number of Concurrent Protected Maintenance Mode IOPS per Device.",
						MarkdownDescription: "Number of Concurrent Protected Maintenance Mode IOPS per Device.",
						Computed:            true,
					},
					"data_layout": schema.StringAttribute{
						Description:         "Data Layout.",
						MarkdownDescription: "Data Layout.",
						Computed:            true,
					},
					"vtree_migration_io_priority_bw_limit_per_device_kbps": schema.Int64Attribute{
						Description:         "VTree Migration Bandwidth Limit per Device.",
						MarkdownDescription: "VTree Migration Bandwidth Limit per Device.",
						Computed:            true,
					},
					"vtree_migration_io_priority_policy": schema.StringAttribute{
						Description:         "VTree Migration IO Priority Policy.",
						MarkdownDescription: "VTree Migration IO Priority Policy.",
						Computed:            true,
					},
					"address_space_usage": schema.StringAttribute{
						Description:         "Address space usage.",
						MarkdownDescription: "Address space usage.",
						Computed:            true,
					},
					"external_acceleration_type": schema.StringAttribute{
						Description:         "External acceleration type.",
						MarkdownDescription: "External acceleration type.",
						Computed:            true,
					},
					"persistent_checksum_state": schema.StringAttribute{
						Description:         "Persistent Checksum State.",
						MarkdownDescription: "Persistent Checksum State.",
						Computed:            true,
					},
					"use_rf_cache": schema.BoolAttribute{
						Description:         "Use Read Flash Cache.",
						MarkdownDescription: "Use Read Flash Cache.",
						Computed:            true,
					},
					"checksum_enabled": schema.BoolAttribute{
						Description:         "Checksum Enabled.",
						MarkdownDescription: "Checksum Enabled.",
						Computed:            true,
					},
					"compression_method": schema.StringAttribute{
						Description:         "Compression method.",
						MarkdownDescription: "Compression method.",
						Computed:            true,
					},
					"fragmentation_enabled": schema.BoolAttribute{
						Description:         "Fragmentation Enabled.",
						MarkdownDescription: "Fragmentation Enabled.",
						Computed:            true,
					},
					"capacity_usage_state": schema.StringAttribute{
						Description:         "Capacity usage state (normal/high/critical/full).",
						MarkdownDescription: "Capacity usage state (normal/high/critical/full).",
						Computed:            true,
					},
					"capacity_usage_type": schema.StringAttribute{
						Description:         "Usage state reason.",
						MarkdownDescription: "Usage state reason.",
						Computed:            true,
					},
					"address_space_usage_type": schema.StringAttribute{
						Description:         "Address space usage reason.",
						MarkdownDescription: "Address space usage reason.",
						Computed:            true,
					},
					"bg_scanner_compare_error_action": schema.StringAttribute{
						Description:         "Scanner compare-error action.",
						MarkdownDescription: "Scanner compare-error action.",
						Computed:            true,
					},
					"bg_scanner_read_error_action": schema.StringAttribute{
						Description:         "Scanner read-error action.",
						MarkdownDescription: "Scanner read-error action.",
						Computed:            true,
					},
					"replication_capacity_max_ratio": schema.Int64Attribute{
						Description:         "Replication allowed capacity.",
						MarkdownDescription: "Replication allowed capacity.",
						Computed:            true,
					},
					"persistent_checksum_enabled": schema.BoolAttribute{
						Description:         "Persistent checksum enabled.",
						MarkdownDescription: "Persistent checksum enabled.",
						Computed:            true,
					},
					"persistent_checksum_builder_limit_kb": schema.Int64Attribute{
						Description:         "Persistent checksum builder limit.",
						MarkdownDescription: "Persistent checksum builder limit.",
						Computed:            true,
					},
					"persistent_checksum_validate_on_read": schema.BoolAttribute{
						Description:         "Persistent checksum validation on read.",
						MarkdownDescription: "Persistent checksum validation on read.",
						Computed:            true,
					},
					"vtree_migration_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
						Description:         "Number of concurrent VTree migration IOPS per device.",
						MarkdownDescription: "Number of concurrent VTree migration IOPS per device.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_policy": schema.StringAttribute{
						Description:         "Protected maintenance mode IO priority policy.",
						MarkdownDescription: "Protected maintenance mode IO priority policy.",
						Computed:            true,
					},
					"background_scanner_mode": schema.StringAttribute{
						Description:         "Scanner mode.",
						MarkdownDescription: "Scanner mode.",
						Computed:            true,
					},
					"media_type": schema.StringAttribute{
						Description:         "Media type.",
						MarkdownDescription: "Media type.",
						Computed:            true,
					},
					"capacity_alert_high_threshold": schema.Int64Attribute{
						Description:         "Capacity alert high threshold.",
						MarkdownDescription: "Capacity alert high threshold.",
						Computed:            true,
					},
					"capacity_alert_critical_threshold": schema.Int64Attribute{
						Description:         "Capacity alert critical threshold.",
						MarkdownDescription: "Capacity alert critical threshold.",
						Computed:            true,
					},
					"vtree_migration_io_priority_app_iops_per_device_threshold": schema.Int64Attribute{
						Description:         "VTree migration IO priority App IOPS per device threshold.",
						MarkdownDescription: "VTree migration IO priority App IOPS per device threshold.",
						Computed:            true,
					},
					"vtree_migration_io_priority_app_bw_per_device_threshold_kbps": schema.Int64Attribute{
						Description:         "VTree migration IO priority App bandwidth per device threshold in Kbps.",
						MarkdownDescription: "VTree migration IO priority App bandwidth per device threshold in Kbps.",
						Computed:            true,
					},
					"vtree_migration_io_priority_quiet_period_msec": schema.Int64Attribute{
						Description:         "VTree migration IO priority quiet period in Msec.",
						MarkdownDescription: "VTree migration IO priority quiet period in Msec.",
						Computed:            true,
					},
					"fgl_accp_id": schema.StringAttribute{
						Description:         "FGL ID.",
						MarkdownDescription: "FGL ID.",
						Computed:            true,
					},
					"fgl_extra_capacity": schema.Int64Attribute{
						Description:         "FGL extra capacity.",
						MarkdownDescription: "FGL extra capacity.",
						Computed:            true,
					},
					"fgl_overprovisioning_factor": schema.Int64Attribute{
						Description:         "FGL overprovisioning factor.",
						MarkdownDescription: "FGL overprovisioning factor.",
						Computed:            true,
					},
					"fgl_write_atomicity_size": schema.Int64Attribute{
						Description:         "FGL write atomicity size.",
						MarkdownDescription: "FGL write atomicity size.",
						Computed:            true,
					},
					"fgl_nvdimm_write_cache_size_mb": schema.Int64Attribute{
						Description:         "FGL NVDIMM write cache size in Mb.",
						MarkdownDescription: "FGL NVDIMM write cache size in Mb.",
						Computed:            true,
					},
					"fgl_nvdimm_metadata_amotization_x100": schema.Int64Attribute{
						Description:         "FGL NVDIMM metadata amortization.",
						MarkdownDescription: "FGL NVDIMM metadata amortization.",
						Computed:            true,
					},
					"fgl_perf_profile": schema.StringAttribute{
						Description:         "FGL performance profile.",
						MarkdownDescription: "FGL performance profile.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_bw_limit_per_device_kbps": schema.Int64Attribute{
						Description:         "Protected maintenance mode IO priority bandwidth limit per device in Kbps.",
						MarkdownDescription: "Protected maintenance mode IO priority bandwidth limit per device in Kbps.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_app_iops_per_device_threshold": schema.Int64Attribute{
						Description:         "Protected maintenance mode IO priority app IOPS per device threshold.",
						MarkdownDescription: "Protected maintenance mode IO priority app IOPS per device threshold.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps": schema.Int64Attribute{
						Description:         "Protected maintenance mode IO priority app bandwidth per device threshold in Kbps.",
						MarkdownDescription: "Protected maintenance mode IO priority app bandwidth per device threshold in Kbps.",
						Computed:            true,
					},
					"protected_maintenance_mode_io_priority_quiet_period_msec": schema.Int64Attribute{
						Description:         "Protected maintenance mode IO priority quiet period in Msec.",
						MarkdownDescription: "Protected maintenance mode IO priority quiet period in Msec.",
						Computed:            true,
					},
					"volumes": schema.ListNestedAttribute{
						Description:         "List of volumes associated with storage pool.",
						MarkdownDescription: "List of volumes associated with storage pool.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "Volume ID.",
									MarkdownDescription: "Volume ID.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Volume name.",
									MarkdownDescription: "Volume name.",
									Computed:            true,
								},
							},
						},
					},
					"sds": schema.ListNestedAttribute{
						Description:         "List of SDS associated with storage pool.",
						MarkdownDescription: "List of SDS associated with storage pool.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "SDS ID.",
									MarkdownDescription: "SDS ID.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "SDS name.",
									MarkdownDescription: "SDS name.",
									Computed:            true,
								},
							},
						},
					},
					"links": schema.ListNestedAttribute{
						Description:         "Specifies the links asscociated with storage pool.",
						MarkdownDescription: "Specifies the links asscociated with storage pool.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description:         "Specifies the relationship with the storage pool.",
									MarkdownDescription: "Specifies the relationship with the storage pool.",
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
			},
		},
	},
}
