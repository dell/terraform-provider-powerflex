/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// DeviceDataSourceSchema defines the schema for device datasource
var DeviceDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing device from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing device from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder id of device datasource.",
			MarkdownDescription: "Placeholder id of device datasource.",
			Computed:            true,
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
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.DeviceFilter{})),
		},
	},
}

// DeviceDataAttributes define the schema of devices
var DeviceDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"fgl_nvdimm_metadata_amortization_x100": schema.Int64Attribute{
		Description:         "Fgl Nvdimm Metadata Amortization X100 Of The Device Instance.",
		MarkdownDescription: "Fgl Nvdimm Metadata Amortization X100 Of The Device Instance.",
		Computed:            true,
	},
	"logical_sector_size_in_bytes": schema.Int64Attribute{
		Description:         "Logical Sector Size In Bytes Of The Device Instance.",
		MarkdownDescription: "Logical Sector Size In Bytes Of The Device Instance.",
		Computed:            true,
	},
	"fgl_nvdimm_write_cache_size": schema.Int64Attribute{
		Description:         "Fgl Nvdimm Write Cache Size Of The Device Instance.",
		MarkdownDescription: "Fgl Nvdimm Write Cache Size Of The Device Instance.",
		Computed:            true,
	},
	"acceleration_pool_id": schema.StringAttribute{
		Description:         "Acceleration Pool_id Of The Device Instance.",
		MarkdownDescription: "Acceleration Pool_id Of The Device Instance.",
		Computed:            true,
	},
	"rfcache_props": schema.SingleNestedAttribute{
		Description:         "Rfcache Props Of The Device Instance.",
		MarkdownDescription: "Rfcache Props Of The Device Instance.",
		Computed:            true,
		Attributes:          getRfcachePropsParamsSchema(),
	},
	"sds_id": schema.StringAttribute{
		Description:         "Sds ID Of The Device Instance.",
		MarkdownDescription: "Sds ID Of The Device Instance.",
		Computed:            true,
	},
	"storage_pool_id": schema.StringAttribute{
		Description:         "Storage Pool ID Of The Device Instance.",
		MarkdownDescription: "Storage Pool ID Of The Device Instance.",
		Computed:            true,
	},
	"capacity_limit_in_kb": schema.Int64Attribute{
		Description:         "Capacity Limit In Kb Of The Device Instance.",
		MarkdownDescription: "Capacity Limit In Kb Of The Device Instance.",
		Computed:            true,
	},
	"error_state": schema.StringAttribute{
		Description:         "Error State Of The Device Instance.",
		MarkdownDescription: "Error State Of The Device Instance.",
		Computed:            true,
	},
	"capacity": schema.Int64Attribute{
		Description:         "Capacity Of The Device Instance.",
		MarkdownDescription: "Capacity Of The Device Instance.",
		Computed:            true,
	},
	"device_type": schema.StringAttribute{
		Description:         "Device Type Of The Device Instance.",
		MarkdownDescription: "Device Type Of The Device Instance.",
		Computed:            true,
	},
	"persistent_checksum_state": schema.StringAttribute{
		Description:         "Persistent Checksum State Of The Device Instance.",
		MarkdownDescription: "Persistent Checksum State Of The Device Instance.",
		Computed:            true,
	},
	"device_state": schema.StringAttribute{
		Description:         "State Of The Device Instance.",
		MarkdownDescription: "State Of The Device Instance.",
		Computed:            true,
	},
	"led_setting": schema.StringAttribute{
		Description:         "LED Setting Of The Device Instance.",
		MarkdownDescription: "LED Setting Of The Device Instance.",
		Computed:            true,
	},
	"max_capacity_in_kb": schema.Int64Attribute{
		Description:         "Max Capacity In Kb Of The Device Instance.",
		MarkdownDescription: "Max Capacity In Kb Of The Device Instance.",
		Computed:            true,
	},
	"sp_sds_id": schema.StringAttribute{
		Description:         "Sp Sds Id Of The Device Instance.",
		MarkdownDescription: "Sp Sds Id Of The Device Instance.",
		Computed:            true,
	},
	"long_successful_ios": schema.SingleNestedAttribute{
		Description:         "Long Successful Ios Of The Device Instance.",
		MarkdownDescription: "Long Successful Ios Of The Device Instance.",
		Computed:            true,
		Attributes:          getLongSuccessfulIosPropsParamsSchema(),
	},
	"aggregated_state": schema.StringAttribute{
		Description:         "Aggregated State Of The Device Instance.",
		MarkdownDescription: "Aggregated State Of The Device Instance.",
		Computed:            true,
	},
	"temperature_state": schema.StringAttribute{
		Description:         "Temperature State Of The Device Instance.",
		MarkdownDescription: "Temperature State Of The Device Instance.",
		Computed:            true,
	},
	"ssd_end_of_life_state": schema.StringAttribute{
		Description:         "Ssd End Of Life State Of The Device Instance.",
		MarkdownDescription: "Ssd End Of Life State Of The Device Instance.",
		Computed:            true,
	},
	"model_name": schema.StringAttribute{
		Description:         "Model Name Of The Device Instance.",
		MarkdownDescription: "Model Name Of The Device Instance.",
		Computed:            true,
	},
	"vendor_name": schema.StringAttribute{
		Description:         "Vendor Name Of The Device Instance.",
		MarkdownDescription: "Vendor Name Of The Device Instance.",
		Computed:            true,
	},
	"raid_controller_serial_number": schema.StringAttribute{
		Description:         "Raid Controller Serial Number Of The Device Instance.",
		MarkdownDescription: "Raid Controller Serial Number Of The Device Instance.",
		Computed:            true,
	},
	"firmware_version": schema.StringAttribute{
		Description:         "Firmware Version Of The Device Instance.",
		MarkdownDescription: "Firmware Version Of The Device Instance.",
		Computed:            true,
	},
	"cache_look_ahead_active": schema.BoolAttribute{
		Description:         "Cache Look Ahead Active Of The Device Instance.",
		MarkdownDescription: "Cache Look Ahead Active Of The Device Instance.",
		Computed:            true,
	},
	"write_cache_active": schema.BoolAttribute{
		Description:         "Write Cache Active Of The Device Instance.",
		MarkdownDescription: "Write Cache Active Of The Device Instance.",
		Computed:            true,
	},
	"ata_security_active": schema.BoolAttribute{
		Description:         "Ata Security Active Of The Device Instance.",
		MarkdownDescription: "Ata Security Active Of The Device Instance.",
		Computed:            true,
	},
	"physical_sector_size_in_bytes": schema.Int64Attribute{
		Description:         "Physical Sector Size In Bytes Of The Device Instance.",
		MarkdownDescription: "Physical Sector Size In Bytes Of The Device Instance.",
		Computed:            true,
	},
	"media_failing": schema.BoolAttribute{
		Description:         "Media Failing Of The Device Instance.",
		MarkdownDescription: "Media Failing Of The Device Instance.",
		Computed:            true,
	},
	"slot_number": schema.StringAttribute{
		Description:         "Slot Number Of The Device Instance.",
		MarkdownDescription: "Slot Number Of The Device Instance.",
		Computed:            true,
	},
	"external_acceleration_type": schema.StringAttribute{
		Description:         "External Acceleration Type Of The Device Instance.",
		MarkdownDescription: "External Acceleration Type Of The Device Instance.",
		Computed:            true,
	},
	"auto_detect_media_type": schema.StringAttribute{
		Description:         "Auto Detect Media Type Of The Device Instance.",
		MarkdownDescription: "Auto Detect Media Type Of The Device Instance.",
		Computed:            true,
	},
	"storage_props": schema.SingleNestedAttribute{
		Description:         "Storage Props Of The Device Instance.",
		MarkdownDescription: "Storage Props Of The Device Instance.",
		Computed:            true,
		Attributes:          getStoragePropsParamsSchema(),
	},
	"acceleration_props": schema.SingleNestedAttribute{
		Description:         "Acceleration Props Of The Device Instance.",
		MarkdownDescription: "Acceleration Props Of The Device Instance.",
		Computed:            true,
		Attributes:          getAccelerationPropsParamsSchema(),
	},
	"device_current_path_name": schema.StringAttribute{
		Description:         "Device Current Path Name Of The Device Instance.",
		MarkdownDescription: "Device Current Path Name Of The Device Instance.",
		Computed:            true,
	},
	"device_original_path_name": schema.StringAttribute{
		Description:         "Device Original Path Name Of The Device Instance.",
		MarkdownDescription: "Device Original Path Name Of The Device Instance.",
		Computed:            true,
	},
	"rfcache_error_device_does_not_exist": schema.BoolAttribute{
		Description:         "Rfcache Error Device Does Not Exist Of The Device Instance.",
		MarkdownDescription: "Rfcache Error Device Does Not Exist Of The Device Instance.",
		Computed:            true,
	},
	"media_type": schema.StringAttribute{
		Description:         "Media Type Of The Device Instance.",
		MarkdownDescription: "Media Type Of The Device Instance.",
		Computed:            true,
	},
	"serial_number": schema.StringAttribute{
		Description:         "Serial Number Of The Device Instance.",
		MarkdownDescription: "Serial Number Of The Device Instance.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name Of The Device Instance.",
		MarkdownDescription: "Name Of The Device Instance.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		Description:         "Unique ID Of The Device Instance.",
		MarkdownDescription: "Unique ID Of The Device Instance.",
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
			Description:         "Device UUID RfCache Parameters Of The Device Instance.",
			MarkdownDescription: "Device UUID RfCache Parameters Of The Device Instance.",
			Computed:            true,
		},
		"rfcache_error_stuck_io": schema.BoolAttribute{
			Description:         "Rfcache Error Stuck Io RfCache Parameters Of The Device Instance.",
			MarkdownDescription: "Rfcache Error Stuck Io RfCache Parameters Of The Device Instance.",
			Computed:            true,
		},
		"rfcache_error_heavy_load_cache_skip": schema.BoolAttribute{
			Description:         "Rfcache Error Heavy Load Cache Skip RfCache Parameters Of The Device Instance.",
			MarkdownDescription: "rfcache_error_heavy_load_cache_skip RfCache Parameters Of The Device Instance.",
			Computed:            true,
		},
		"rfcache_error_card_io_error": schema.BoolAttribute{
			Description:         "Rfcache Error Card Io error RfCache Parameters Of The Device Instance.",
			MarkdownDescription: "Rfcache Error Card Io error RfCache Parameters Of The Device Instance.",
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
			Description:         "Threshold Window Parameters Of The Device Instance.",
			MarkdownDescription: "Threshold Window Parameters Of The Device Instance.",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			Description:         "Window Size in seconds Window Parameters Of The Device Instance.",
			MarkdownDescription: "Window Size in seconds Window Parameters Of The Device Instance.",
			Computed:            true,
		},
		"last_oscillation_count": schema.Int64Attribute{
			Description:         "Last Oscillation Count Window Parameters Of The Device Instance.",
			MarkdownDescription: "Last Oscillation Count Window Parameters Of The Device Instance.",
			Computed:            true,
		},
		"last_oscillation_time": schema.Int64Attribute{
			Description:         "Last Oscillation Time Window Parameters Of The Device Instance.",
			MarkdownDescription: "Last Oscillation Time Window Parameters Of The Device Instance.",
			Computed:            true,
		},
		"max_failures_count": schema.Int64Attribute{
			MarkdownDescription: "Max Failures Count Window Parameters Of The Device Instance.",
			Description:         "Max Failures Count Window Parameters Of The Device Instance.",
			Computed:            true,
		},
	}
}

func getStoragePropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"fgl_acc_device_id": schema.StringAttribute{
			Description:         "FGL(Fujitsu General Limited) Accelerator(ACC) Device Id Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "FGL(Fujitsu General Limited) Accelerator(ACC) Device Id Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"fgl_nvdimm_size_mb": schema.Int64Attribute{
			Description:         "FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_fgl_nvdimm_size_mb": schema.Int64Attribute{
			Description:         "Destination FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_fgl_acc_device_id": schema.StringAttribute{
			Description:         "Destination FGL(Fujitsu General Limited) Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination FGL(Fujitsu General Limited) Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"checksum_mode": schema.StringAttribute{
			Description:         "Checksum Mode Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Checksum Mode Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_checksum_mode": schema.StringAttribute{
			Description:         "Destination Checksum Mode Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination Checksum Mode Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"checksum_acc_device_id": schema.StringAttribute{
			Description:         "Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_checksum_acc_device_id": schema.StringAttribute{
			Description:         "Destination Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"checksum_size_mb": schema.Int64Attribute{
			Description:         "Checksum Size MB Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Checksum Size MB Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"is_checksum_fully_calculated": schema.BoolAttribute{
			Description:         "Is Checksum Fully Calculated Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Is Checksum Fully Calculated Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"checksum_changelog_acc_device_id": schema.StringAttribute{
			Description:         "Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_checksum_changelog_acc_device_id": schema.StringAttribute{
			Description:         "Destination Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"checksum_changelog_size_mb": schema.Int64Attribute{
			Description:         "Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
		"dest_checksum_changelog_size_mb": schema.Int64Attribute{
			Description:         "Destination Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Destination Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
	}
}

func getAccelerationPropsParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"acc_used_capacity_in_kb": schema.StringAttribute{
			Description:         "Accelerator(ACC) Used Capacity In KB Acceleration Properties Parameters Of The Device Instance.",
			MarkdownDescription: "Accelerator(ACC) Used Capacity In KB Acceleration Properties Parameters Of The Device Instance.",
			Computed:            true,
		},
	}
}
