package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DeviceDataSourceSchema defines the schema for device datasource
var DeviceDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource can be used to fetch information related to devices from a PowerFlex array.",
	MarkdownDescription: "This datasource can be used to fetch information related to devices from a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier of the device instance." +
				" Conflicts with 'name'.",
			MarkdownDescription: "Unique identifier of the device instance." +
				" Conflicts with `name`.",
			Optional: true,
		},
		"name": schema.StringAttribute{
			Description: "Unique name of the device instance." +
				" Conflicts with 'id'.",
			MarkdownDescription: "Unique name of the device instance." +
				" Conflicts with `id`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("id")),
			},
		},
		"storage_pool_id": schema.StringAttribute{
			Description:         "ID of the storage pool. Conflicts with 'storage_pool_name'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "ID of the storage pool. Conflicts with `storage_pool_name`.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("storage_pool_name")),
			},
		},
		"storage_pool_name": schema.StringAttribute{
			Description:         "Name of the storage pool. Conflicts with 'storage_pool_id'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Name of the storage pool. Conflicts with `storage_pool_id`.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "ID of the protection domain. Conflicts with 'protection_domain_name'.",
			MarkdownDescription: "ID of the protection domain. Conflicts with `protection_domain_name`.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description:         "Name of the protection domain. Conflicts with 'protection_domain_id'.",
			MarkdownDescription: "Name of the protection domain. Conflicts with `protection_domain_id`.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"sds_id": schema.StringAttribute{
			Description:         "ID of the SDS. Conflicts with 'sds_name'.",
			MarkdownDescription: "ID of the SDS. Conflicts with `sds_name`.",
			Computed:            true,
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("sds_name")),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"sds_name": schema.StringAttribute{
			Description:         "Name of the SDS. Conflicts with 'sds_id'.",
			MarkdownDescription: "Name of the SDS. Conflicts with `sds_id`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"current_path": schema.StringAttribute{
			Description:         "Path of device",
			MarkdownDescription: "Path of device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"device_model": schema.ListNestedAttribute{
			Description:         "List of devices fetched.",
			MarkdownDescription: "List of devices fetched.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: DeviceDataAttributes,
			},
		},
	},
}

var DeviceDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"fgl_nvdimm_metadata_amortization_x100": schema.Int64Attribute{
		Description:         "fgl_nvdimm_metadata_amortization_x100.",
		MarkdownDescription: "fgl_nvdimm_metadata_amortization_x100.",
		Computed:            true,
	},
	"logical_sector_size_in_bytes": schema.Int64Attribute{
		Description:         "logical_sector_size_in_bytes",
		MarkdownDescription: "logical_sector_size_in_bytes",
		Computed:            true,
	},
	"fgl_nvdimm_write_cache_size":  schema.Int64Attribute{
		Description:         "fgl_nvdimm_write_cache_size",
		MarkdownDescription: "fgl_nvdimm_write_cache_size",
		Computed:            true,
	},
	"acceleration_pool_id": schema.StringAttribute{
		Description:         "acceleration_pool_id",
		MarkdownDescription: "acceleration_pool_id",
		Computed:            true,
	},
	"rfcache_props": schema.SingleNestedAttribute{
		Description:         "rfcache_props",
		MarkdownDescription: "rfcache_props",
		Computed:            true,
		Attributes:          getRfcachePropsParamsSchema(),
	},
	"sds_id": schema.StringAttribute{
		Description:         "sds_id",
		MarkdownDescription: "sds_id",
		Computed:            true,
	},
	"storage_pool_id": schema.StringAttribute{
		Description:         "storage_pool_id",
		MarkdownDescription: "storage_pool_id",
		Computed:            true,
	},
	"capacity_limit_in_kb": schema.Int64Attribute{
		Description:         "capacity_limit_in_kb",
		MarkdownDescription: "capacity_limit_in_kb",
		Computed:            true,
	},
	"error_state": schema.StringAttribute{
		Description:         "error_state",
		MarkdownDescription: "error_state",
		Computed:            true,
	},
	"capacity": schema.Int64Attribute{
		Description:         "capacity",
		MarkdownDescription: "capacity",
		Computed:            true,
	},
	"device_type": schema.StringAttribute{
		Description:         "device_type",
		MarkdownDescription: "device_type",
		Computed:            true,
	},
	"persistent_checksum_state": schema.StringAttribute{
		Description:         "persistent_checksum_state",
		MarkdownDescription: "persistent_checksum_state",
		Computed:            true,
	},
	"device_state": schema.StringAttribute{
		Description:         "List of devices fetched.",
		MarkdownDescription: "List of devices fetched.",
		Computed:            true,
	},
	"led_setting": schema.StringAttribute{
		Description:         "led_setting",
		MarkdownDescription: "led_setting",
		Computed:            true,
	},
	"max_capacity_in_kb": schema.Int64Attribute{
		Description:         "max_capacity_in_kb",
		MarkdownDescription: "max_capacity_in_kb",
		Computed:            true,
	},
	"sp_sds_id": schema.StringAttribute{
		Description:         "sp_sds_id",
		MarkdownDescription: "sp_sds_id",
		Computed:            true,
	},
	"long_successful_ios": schema.SingleNestedAttribute{
		Description:         "long_successful_ios",
		MarkdownDescription: "long_successful_ios",
		Computed:            true,
		Attributes:          getLongSuccessfulIosPropsParamsSchema(),
	},
	"aggregated_state": schema.StringAttribute{
		Description:         "aggregated_state",
		MarkdownDescription: "aggregated_state",
		Computed:            true,
	},
	"temperature_state": schema.StringAttribute{
		Description:         "temperature_state",
		MarkdownDescription: "temperature_state",
		Computed:            true,
	},
	"ssd_end_of_life_state": schema.StringAttribute{
		Description:         "ssd_end_of_life_state",
		MarkdownDescription: "ssd_end_of_life_state",
		Computed:            true,
	},
	"model_name": schema.StringAttribute{
		Description:         "model_name",
		MarkdownDescription: "model_name",
		Computed:            true,
	},
	"vendor_name": schema.StringAttribute{
		Description:         "vendor_name",
		MarkdownDescription: "vendor_name",
		Computed:            true,
	},
	"raid_controller_serial_number": schema.StringAttribute{
		Description:         "raid_controller_serial_number",
		MarkdownDescription: "raid_controller_serial_number",
		Computed:            true,
	},
	"firmware_version": schema.StringAttribute{
		Description:         "firmware_version",
		MarkdownDescription: "firmware_version",
		Computed:            true,
	},
	"cache_look_ahead_active": schema.BoolAttribute{
		Description:         "cache_look_ahead_active",
		MarkdownDescription: "cache_look_ahead_active",
		Computed:            true,
	},
	"write_cache_active": schema.BoolAttribute{
		Description:         "write_cache_active",
		MarkdownDescription: "write_cache_active",
		Computed:            true,
	},
	"ata_security_active": schema.BoolAttribute{
		Description:         "ata_security_active",
		MarkdownDescription: "ata_security_active",
		Computed:            true,
	},
	"physical_sector_size_in_bytes": schema.Int64Attribute{
		Description:         "physical_sector_size_in_bytes",
		MarkdownDescription: "physical_sector_size_in_bytes",
		Computed:            true,
	},
	"media_failing": schema.BoolAttribute{
		Description:         "media_failing",
		MarkdownDescription: "media_failing",
		Computed:            true,
	},
	"slot_number": schema.StringAttribute{
		Description:         "slot_number",
		MarkdownDescription: "slot_number",
		Computed:            true,
	},
	"external_acceleration_type": schema.StringAttribute{
		Description:         "external_acceleration_type",
		MarkdownDescription: "external_acceleration_type",
		Computed:            true,
	},
	"auto_detect_media_type": schema.StringAttribute{
		Description:         "auto_detect_media_type",
		MarkdownDescription: "auto_detect_media_type",
		Computed:            true,
	},
	"storage_props": schema.SingleNestedAttribute{
		Description:         "storage_props",
		MarkdownDescription: "storage_props",
		Computed:            true,
		Attributes:          getStoragePropsParamsSchema(),
	},
	"acceleration_props": schema.SingleNestedAttribute{
		Description:         "acceleration_props",
		MarkdownDescription: "acceleration_props",
		Computed:            true,
		Attributes:          getAccelerationPropsParamsSchema(),
	},
	"device_current_path_name": schema.StringAttribute{
		Description:         "device_current_path_name",
		MarkdownDescription: "device_current_path_name",
		Computed:            true,
	},
	"device_original_path_name": schema.StringAttribute{
		Description:         "device_original_path_name",
		MarkdownDescription: "device_original_path_name",
		Computed:            true,
	},
	"rfcache_error_device_does_not_exist": schema.BoolAttribute{
		Description:         "rfcache_error_device_does_not_exist",
		MarkdownDescription: "rfcache_error_device_does_not_exist",
		Computed:            true,
	},
	"media_type": schema.StringAttribute{
		Description:         "media_type",
		MarkdownDescription: "media_type",
		Computed:            true,
	},
	"serial_number": schema.StringAttribute{
		Description:         "serial_number",
		MarkdownDescription: "serial_number",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "name",
		MarkdownDescription: "name",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		Description:         "id",
		MarkdownDescription: "id",
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

func getRfcachePropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"device_uuid": schema.StringAttribute{
			Description:         "device_uuid",
			MarkdownDescription: "device_uuid",
			Computed:            true,
		},
		"rfcache_error_stuck_io": schema.BoolAttribute{
			Description:         "rfcache_error_stuck_io",
			MarkdownDescription: "rfcache_error_stuck_io",
			Computed:            true,
		},
		"rfcache_error_heavy_load_cache_skip": schema.BoolAttribute{
			Description:         "rfcache_error_heavy_load_cache_skip",
			MarkdownDescription: "rfcache_error_heavy_load_cache_skip",
			Computed:            true,
		},
		"rfcache_error_card_io_error": schema.BoolAttribute{
			Description:         "rfcache_error_card_io_error",
			MarkdownDescription: "rfcache_error_card_io_error",
			Computed:            true,
		},
	}
}

func getLongSuccessfulIosPropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			Description:         "Short Window Parameters.",
			MarkdownDescription: "Short Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsDeviceSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			Description:         "Medium Window Parameters.",
			MarkdownDescription: "Medium Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsDeviceSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			Description:         "Long Window Parameters.",
			MarkdownDescription: "Long Window Parameters.",
			Computed:            true,
			Attributes:          getWindowParamsDeviceSchema(),
		},
	}
}

func getWindowParamsDeviceSchema() map[string]schema.Attribute {
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
		"last_oscillation_count": schema.Int64Attribute{
			Description:         "last_oscillation_count",
			MarkdownDescription: "last_oscillation_count",
			Computed:            true,
		},
		"last_oscillation_time": schema.Int64Attribute{
			Description:         "last_oscillation_time",
			MarkdownDescription: "last_oscillation_time",
			Computed:            true,
		},
		"max_failures_count": schema.Int64Attribute{
			MarkdownDescription: "max_failures_count",
			Description:         "max_failures_count",
			Computed:            true,
		},
	}
}

func getStoragePropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"fgl_acc_device_id": schema.StringAttribute{
			Description:         "fgl_acc_device_id",
			MarkdownDescription: "fgl_acc_device_id",
			Computed:            true,
		},
		"fgl_nvdimm_size_mb": schema.Int64Attribute{
			Description:         "fgl_nvdimm_size_mb",
			MarkdownDescription: "fgl_nvdimm_size_mb",
			Computed:            true,
		},
		"dest_fgl_nvdimm_size_mb": schema.Int64Attribute{
			Description:         "dest_fgl_nvdimm_size_mb",
			MarkdownDescription: "dest_fgl_nvdimm_size_mb",
			Computed:            true,
		},
		"dest_fgl_acc_device_id": schema.StringAttribute{
			Description:         "dest_fgl_acc_device_id",
			MarkdownDescription: "dest_fgl_acc_device_id",
			Computed:            true,
		},
		"checksum_mode": schema.StringAttribute{
			Description:         "checksum_mode",
			MarkdownDescription: "checksum_mode",
			Computed:            true,
		},
		"dest_checksum_mode": schema.StringAttribute{
			Description:         "dest_checksum_mode",
			MarkdownDescription: "dest_checksum_mode",
			Computed:            true,
		},
		"checksum_acc_device_id": schema.StringAttribute{
			Description:         "checksum_acc_device_id",
			MarkdownDescription: "checksum_acc_device_id",
			Computed:            true,
		},
		"dest_checksum_acc_device_id": schema.StringAttribute{
			Description:         "dest_checksum_acc_device_id",
			MarkdownDescription: "dest_checksum_acc_device_id",
			Computed:            true,
		},
		"checksum_size_mb": schema.Int64Attribute{
			Description:         "checksum_size_mb",
			MarkdownDescription: "checksum_size_mb",
			Computed:            true,
		},
		"is_checksum_fully_calculated": schema.BoolAttribute{
			Description:         "is_checksum_fully_calculated",
			MarkdownDescription: "is_checksum_fully_calculated",
			Computed:            true,
		},
		"checksum_changelog_acc_device_id": schema.StringAttribute{
			Description:         "checksum_changelog_acc_device_id",
			MarkdownDescription: "checksum_changelog_acc_device_id",
			Computed:            true,
		},
		"dest_checksum_changelog_acc_device_id": schema.StringAttribute{
			Description:         "dest_checksum_changelog_acc_device_id",
			MarkdownDescription: "dest_checksum_changelog_acc_device_id",
			Computed:            true,
		},
		"checksum_changelog_size_mb": schema.Int64Attribute{
			Description:         "checksum_changelog_size_mb",
			MarkdownDescription: "checksum_changelog_size_mb",
			Computed:            true,
		},
		"dest_checksum_changelog_size_mb": schema.Int64Attribute{
			Description:         "dest_checksum_changelog_size_mb",
			MarkdownDescription: "dest_checksum_changelog_size_mb",
			Computed:            true,
		},
	} 
}

func getAccelerationPropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"acc_used_capacity_in_kb": schema.StringAttribute{
			Description:         "acc_used_capacity_in_kb",
			MarkdownDescription: "acc_used_capacity_in_kb",
			Computed:            true,
		},
	}
}