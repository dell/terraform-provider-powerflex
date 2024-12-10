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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceModel defines the struct for device resource
type DeviceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	DevicePath               types.String `tfsdk:"device_path"`
	DeviceOriginalPath       types.String `tfsdk:"device_original_path"`
	ProtectionDomainName     types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID       types.String `tfsdk:"protection_domain_id"`
	StoragePoolName          types.String `tfsdk:"storage_pool_name"`
	StoragePoolID            types.String `tfsdk:"storage_pool_id"`
	SdsID                    types.String `tfsdk:"sds_id"`
	SdsName                  types.String `tfsdk:"sds_name"`
	MediaType                types.String `tfsdk:"media_type"`
	ExternalAccelerationType types.String `tfsdk:"external_acceleration_type"`
	DeviceCapacity           types.Int64  `tfsdk:"device_capacity"`
	DeviceCapacityInKB       types.Int64  `tfsdk:"device_capacity_in_kb"`
	DeviceState              types.String `tfsdk:"device_state"`
}

// DeviceDataSourceModel defines struct for device datasource
type DeviceDataSourceModel struct {
	ID           types.String      `tfsdk:"id"`
	DeviceModel  []DeviceModelData `tfsdk:"device_model"`
	DeviceFilter *DeviceFilter     `tfsdk:"filter"`
}

// DeviceFilter defines struct for device filter
type DeviceFilter struct {
	FglNvdimmMetadataAmortizationX100 []types.Int64  `tfsdk:"fgl_nvdimm_metadata_amortization_x100"`
	LogicalSectorSizeInBytes          []types.Int64  `tfsdk:"logical_sector_size_in_bytes"`
	FglNvdimmWriteCacheSize           []types.Int64  `tfsdk:"fgl_nvdimm_write_cache_size"`
	AccelerationPoolID                []types.String `tfsdk:"acceleration_pool_id"`
	SdsID                             []types.String `tfsdk:"sds_id"`
	StoragePoolID                     []types.String `tfsdk:"storage_pool_id"`
	CapacityLimitInKb                 []types.Int64  `tfsdk:"capacity_limit_in_kb"`
	ErrorState                        []types.String `tfsdk:"error_state"`
	Capacity                          []types.Int64  `tfsdk:"capacity"`
	DeviceType                        []types.String `tfsdk:"device_type"`
	PersistentChecksumState           []types.String `tfsdk:"persistent_checksum_state"`
	DeviceState                       []types.String `tfsdk:"device_state"`
	LedSetting                        []types.String `tfsdk:"led_setting"`
	MaxCapacityInKb                   []types.Int64  `tfsdk:"max_capacity_in_kb"`
	SpSdsID                           []types.String `tfsdk:"sp_sds_id"`
	AggregatedState                   []types.String `tfsdk:"aggregated_state"`
	TemperatureState                  []types.String `tfsdk:"temperature_state"`
	SsdEndOfLifeState                 []types.String `tfsdk:"ssd_end_of_life_state"`
	ModelName                         []types.String `tfsdk:"model_name"`
	VendorName                        []types.String `tfsdk:"vendor_name"`
	RaidControllerSerialNumber        []types.String `tfsdk:"raid_controller_serial_number"`
	FirmwareVersion                   []types.String `tfsdk:"firmware_version"`
	CacheLookAheadActive              types.Bool     `tfsdk:"cache_look_ahead_active"`
	WriteCacheActive                  types.Bool     `tfsdk:"write_cache_active"`
	AtaSecurityActive                 types.Bool     `tfsdk:"ata_security_active"`
	PhysicalSectorSizeInBytes         []types.Int64  `tfsdk:"physical_sector_size_in_bytes"`
	MediaFailing                      types.Bool     `tfsdk:"media_failing"`
	SlotNumber                        []types.String `tfsdk:"slot_number"`
	ExternalAccelerationType          []types.String `tfsdk:"external_acceleration_type"`
	AutoDetectMediaType               []types.String `tfsdk:"auto_detect_media_type"`
	DeviceCurrentPathName             []types.String `tfsdk:"device_current_path_name"`
	DeviceOriginalPathName            []types.String `tfsdk:"device_original_path_name"`
	RfcacheErrorDeviceDoesNotExist    types.Bool     `tfsdk:"rfcache_error_device_does_not_exist"`
	MediaType                         []types.String `tfsdk:"media_type"`
	SerialNumber                      []types.String `tfsdk:"serial_number"`
	Name                              []types.String `tfsdk:"name"`
	ID                                []types.String `tfsdk:"id"`
}

// DeviceModelData defines struct for device model
type DeviceModelData struct {
	FglNvdimmMetadataAmortizationX100 types.Int64            `tfsdk:"fgl_nvdimm_metadata_amortization_x100"`
	LogicalSectorSizeInBytes          types.Int64            `tfsdk:"logical_sector_size_in_bytes"`
	FglNvdimmWriteCacheSize           types.Int64            `tfsdk:"fgl_nvdimm_write_cache_size"`
	AccelerationPoolID                types.String           `tfsdk:"acceleration_pool_id"`
	RfcacheProps                      RfcachePropsModel      `tfsdk:"rfcache_props"`
	SdsID                             types.String           `tfsdk:"sds_id"`
	StoragePoolID                     types.String           `tfsdk:"storage_pool_id"`
	CapacityLimitInKb                 types.Int64            `tfsdk:"capacity_limit_in_kb"`
	ErrorState                        types.String           `tfsdk:"error_state"`
	Capacity                          types.Int64            `tfsdk:"capacity"`
	DeviceType                        types.String           `tfsdk:"device_type"`
	PersistentChecksumState           types.String           `tfsdk:"persistent_checksum_state"`
	DeviceState                       types.String           `tfsdk:"device_state"`
	LedSetting                        types.String           `tfsdk:"led_setting"`
	MaxCapacityInKb                   types.Int64            `tfsdk:"max_capacity_in_kb"`
	SpSdsID                           types.String           `tfsdk:"sp_sds_id"`
	LongSuccessfulIos                 LongSuccessfulIosModel `tfsdk:"long_successful_ios"`
	AggregatedState                   types.String           `tfsdk:"aggregated_state"`
	TemperatureState                  types.String           `tfsdk:"temperature_state"`
	SsdEndOfLifeState                 types.String           `tfsdk:"ssd_end_of_life_state"`
	ModelName                         types.String           `tfsdk:"model_name"`
	VendorName                        types.String           `tfsdk:"vendor_name"`
	RaidControllerSerialNumber        types.String           `tfsdk:"raid_controller_serial_number"`
	FirmwareVersion                   types.String           `tfsdk:"firmware_version"`
	CacheLookAheadActive              types.Bool             `tfsdk:"cache_look_ahead_active"`
	WriteCacheActive                  types.Bool             `tfsdk:"write_cache_active"`
	AtaSecurityActive                 types.Bool             `tfsdk:"ata_security_active"`
	PhysicalSectorSizeInBytes         types.Int64            `tfsdk:"physical_sector_size_in_bytes"`
	MediaFailing                      types.Bool             `tfsdk:"media_failing"`
	SlotNumber                        types.String           `tfsdk:"slot_number"`
	ExternalAccelerationType          types.String           `tfsdk:"external_acceleration_type"`
	AutoDetectMediaType               types.String           `tfsdk:"auto_detect_media_type"`
	StorageProps                      StoragePropsModel      `tfsdk:"storage_props"`
	AccelerationProps                 AccelerationPropsModel `tfsdk:"acceleration_props"`
	DeviceCurrentPathName             types.String           `tfsdk:"device_current_path_name"`
	DeviceOriginalPathName            types.String           `tfsdk:"device_original_path_name"`
	RfcacheErrorDeviceDoesNotExist    types.Bool             `tfsdk:"rfcache_error_device_does_not_exist"`
	MediaType                         types.String           `tfsdk:"media_type"`
	SerialNumber                      types.String           `tfsdk:"serial_number"`
	Name                              types.String           `tfsdk:"name"`
	ID                                types.String           `tfsdk:"id"`
	Links                             []DeviceLinkModel      `tfsdk:"links"`
}

// RfcachePropsModel defines struct for Device
type RfcachePropsModel struct {
	DeviceUUID                     types.String `tfsdk:"device_uuid"`
	RfcacheErrorStuckIO            types.Bool   `tfsdk:"rfcache_error_stuck_io"`
	RfcacheErrorHeavyLoadCacheSkip types.Bool   `tfsdk:"rfcache_error_heavy_load_cache_skip"`
	RfcacheErrorCardIoError        types.Bool   `tfsdk:"rfcache_error_card_io_error"`
}

// LongSuccessfulIosModel defines struct for Device
type LongSuccessfulIosModel struct {
	ShortWindow  DeviceWindowTypeModel `tfsdk:"short_window"`
	MediumWindow DeviceWindowTypeModel `tfsdk:"medium_window"`
	LongWindow   DeviceWindowTypeModel `tfsdk:"long_window"`
}

// DeviceWindowTypeModel defines struct for LongSuccessfulIosModel
type DeviceWindowTypeModel struct {
	Threshold            types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec      types.Int64 `tfsdk:"window_size_in_sec"`
	LastOscillationCount types.Int64 `tfsdk:"last_oscillation_count"`
	LastOscillationTime  types.Int64 `tfsdk:"last_oscillation_time"`
	MaxFailuresCount     types.Int64 `tfsdk:"max_failures_count"`
}

// AccelerationPropsModel defines struct for Device
type AccelerationPropsModel struct {
	AccUsedCapacityInKb types.String `tfsdk:"acc_used_capacity_in_kb"`
}

// StoragePropsModel defines struct for Device
type StoragePropsModel struct {
	FglAccDeviceID                   types.String `tfsdk:"fgl_acc_device_id"`
	FglNvdimmSizeMb                  types.Int64  `tfsdk:"fgl_nvdimm_size_mb"`
	DestFglNvdimmSizeMb              types.Int64  `tfsdk:"dest_fgl_nvdimm_size_mb"`
	DestFglAccDeviceID               types.String `tfsdk:"dest_fgl_acc_device_id"`
	ChecksumMode                     types.String `tfsdk:"checksum_mode"`
	DestChecksumMode                 types.String `tfsdk:"dest_checksum_mode"`
	ChecksumAccDeviceID              types.String `tfsdk:"checksum_acc_device_id"`
	DestChecksumAccDeviceID          types.String `tfsdk:"dest_checksum_acc_device_id"`
	ChecksumSizeMb                   types.Int64  `tfsdk:"checksum_size_mb"`
	IsChecksumFullyCalculated        types.Bool   `tfsdk:"is_checksum_fully_calculated"`
	ChecksumChangelogAccDeviceID     types.String `tfsdk:"checksum_changelog_acc_device_id"`
	DestChecksumChangelogAccDeviceID types.String `tfsdk:"dest_checksum_changelog_acc_device_id"`
	ChecksumChangelogSizeMb          types.Int64  `tfsdk:"checksum_changelog_size_mb"`
	DestChecksumChangelogSizeMb      types.Int64  `tfsdk:"dest_checksum_changelog_size_mb"`
}

// DeviceLinkModel defines struct for device links
type DeviceLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}
