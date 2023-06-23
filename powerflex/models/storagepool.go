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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StoragepoolResourceModel defines struct storage pool resource
type StoragepoolResourceModel struct {
	ID                                                  types.String `tfsdk:"id"`
	ProtectionDomainID                                  types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName                                types.String `tfsdk:"protection_domain_name"`
	Name                                                types.String `tfsdk:"name"`
	MediaType                                           types.String `tfsdk:"media_type"`
	UseRmcache                                          types.Bool   `tfsdk:"use_rmcache"`
	UseRfcache                                          types.Bool   `tfsdk:"use_rfcache"`
	ZeroPaddingEnabled                                  types.Bool   `tfsdk:"zero_padding_enabled"`
	ReplicationJournalCapacity                          types.Int64  `tfsdk:"replication_journal_capacity"`
	CapacityAlertHighThreshold                          types.Int64  `tfsdk:"capacity_alert_high_threshold"`
	CapacityAlertCriticalThreshold                      types.Int64  `tfsdk:"capacity_alert_critical_threshold"`
	ProtectedMaintenanceModeIoPriorityPolicy            types.String `tfsdk:"protected_maintenance_mode_io_priority_policy"`
	ProtectedMaintenanceModeNumOfConcurrentIosPerDevice types.Int64  `tfsdk:"protected_maintenance_mode_num_of_concurrent_ios_per_device"`
	ProtectedMaintenanceModeBwLimitPerDeviceInKbps      types.Int64  `tfsdk:"protected_maintenance_mode_bw_limit_per_device_in_kbps"`
	RebalanceEnabled                                    types.Bool   `tfsdk:"rebalance_enabled"`
	RebalanceIoPriorityPolicy                           types.String `tfsdk:"rebalance_io_priority_policy"`
	RebalanceNumOfConcurrentIosPerDevice                types.Int64  `tfsdk:"rebalance_num_of_concurrent_ios_per_device"`
	RebalanceBwLimitPerDeviceInKbps                     types.Int64  `tfsdk:"rebalance_bw_limit_per_device_in_kbps"`
	VtreeMigrationIoPriorityPolicy                      types.String `tfsdk:"vtree_migration_io_priority_policy"`
	VtreeMigrationNumOfConcurrentIosPerDevice           types.Int64  `tfsdk:"vtree_migration_num_of_concurrent_ios_per_device"`
	VtreeMigrationBwLimitPerDeviceInKbps                types.Int64  `tfsdk:"vtree_migration_bw_limit_per_device_in_kbps"`
	SparePercentage                                     types.Int64  `tfsdk:"spare_percentage"`
	RmCacheWriteHandlingMode                            types.String `tfsdk:"rm_cache_write_handling_mode"`
	RebuildEnabled                                      types.Bool   `tfsdk:"rebuild_enabled"`
	RebuildRebalanceParallelism                         types.Int64  `tfsdk:"rebuild_rebalance_parallelism"`
	Fragmentation                                       types.Bool   `tfsdk:"fragmentation"`
}

// Volume maps the volume schema data.
type Volume struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// SdsData maps the SDS schema data
type SdsData struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// LinkModel maps the link schema data
type LinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// StoragePoolModel maps the storagepool schema data
type StoragePoolModel struct {
	ID                                                              types.String `tfsdk:"id"`
	Name                                                            types.String `tfsdk:"name"`
	RebalanceioPriorityPolicy                                       types.String `tfsdk:"rebalance_io_priority_policy"`
	RebuildioPriorityPolicy                                         types.String `tfsdk:"rebuild_io_priority_policy"`
	RebuildioPriorityBwLimitPerDeviceInKbps                         types.Int64  `tfsdk:"rebuild_io_priority_bw_limit_per_device_in_kbps"`
	RebuildioPriorityNumOfConcurrentIosPerDevice                    types.Int64  `tfsdk:"rebuild_io_priority_num_of_concurrent_ios_per_device"`
	RebalanceioPriorityNumOfConcurrentIosPerDevice                  types.Int64  `tfsdk:"rebalance_io_priority_num_of_concurrent_ios_per_device"`
	RebalanceioPriorityBwLimitPerDeviceInKbps                       types.Int64  `tfsdk:"rebalance_io_priority_bw_limit_per_device_kbps"`
	RebuildioPriorityAppIopsPerDeviceThreshold                      types.Int64  `tfsdk:"rebuild_io_priority_app_iops_per_device_threshold"`
	RebalanceioPriorityAppIopsPerDeviceThreshold                    types.Int64  `tfsdk:"rebalance_io_priority_app_iops_per_device_threshold"`
	RebuildioPriorityAppBwPerDeviceThresholdInKbps                  types.Int64  `tfsdk:"rebuild_io_priority_app_bw_per_device_threshold_kbps"`
	RebalanceioPriorityAppBwPerDeviceThresholdInKbps                types.Int64  `tfsdk:"rebalance_io_priority_app_bw_per_device_threshold_kbps"`
	RebuildioPriorityQuietPeriodInMsec                              types.Int64  `tfsdk:"rebuild_io_priority_quiet_period_msec"`
	RebalanceioPriorityQuietPeriodInMsec                            types.Int64  `tfsdk:"rebalance_io_priority_quiet_period_msec"`
	ZeroPaddingEnabled                                              types.Bool   `tfsdk:"zero_padding_enabled"`
	UseRmcache                                                      types.Bool   `tfsdk:"use_rm_cache"`
	SparePercentage                                                 types.Int64  `tfsdk:"spare_percentage"`
	RmCacheWriteHandlingMode                                        types.String `tfsdk:"rm_cache_write_handling_mode"`
	RebuildEnabled                                                  types.Bool   `tfsdk:"rebuild_enabled"`
	RebalanceEnabled                                                types.Bool   `tfsdk:"rebalance_enabled"`
	NumofParallelRebuildRebalanceJobsPerDevice                      types.Int64  `tfsdk:"num_of_parallel_rebuild_rebalance_jobs_per_device"`
	BackgroundScannerBWLimitKBps                                    types.Int64  `tfsdk:"background_scanner_bw_limit_kbps"`
	ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice   types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device"`
	DataLayout                                                      types.String `tfsdk:"data_layout"`
	VtreeMigrationIoPriorityBwLimitPerDeviceInKbps                  types.Int64  `tfsdk:"vtree_migration_io_priority_bw_limit_per_device_kbps"`
	VtreeMigrationIoPriorityPolicy                                  types.String `tfsdk:"vtree_migration_io_priority_policy"`
	AddressSpaceUsage                                               types.String `tfsdk:"address_space_usage"`
	ExternalAccelerationType                                        types.String `tfsdk:"external_acceleration_type"`
	PersistentChecksumState                                         types.String `tfsdk:"persistent_checksum_state"`
	UseRfcache                                                      types.Bool   `tfsdk:"use_rf_cache"`
	ChecksumEnabled                                                 types.Bool   `tfsdk:"checksum_enabled"`
	CompressionMethod                                               types.String `tfsdk:"compression_method"`
	FragmentationEnabled                                            types.Bool   `tfsdk:"fragmentation_enabled"`
	CapacityUsageState                                              types.String `tfsdk:"capacity_usage_state"`
	CapacityUsageType                                               types.String `tfsdk:"capacity_usage_type"`
	AddressSpaceUsageType                                           types.String `tfsdk:"address_space_usage_type"`
	BgScannerCompareErrorAction                                     types.String `tfsdk:"bg_scanner_compare_error_action"`
	BgScannerReadErrorAction                                        types.String `tfsdk:"bg_scanner_read_error_action"`
	ReplicationCapacityMaxRatio                                     types.Int64  `tfsdk:"replication_capacity_max_ratio"`
	PersistentChecksumEnabled                                       types.Bool   `tfsdk:"persistent_checksum_enabled"`
	PersistentChecksumBuilderLimitKb                                types.Int64  `tfsdk:"persistent_checksum_builder_limit_kb"`
	PersistentChecksumValidateOnRead                                types.Bool   `tfsdk:"persistent_checksum_validate_on_read"`
	VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice             types.Int64  `tfsdk:"vtree_migration_io_priority_num_of_concurrent_ios_per_device"`
	ProtectedMaintenanceModeIoPriorityPolicy                        types.String `tfsdk:"protected_maintenance_mode_io_priority_policy"`
	BackgroundScannerMode                                           types.String `tfsdk:"background_scanner_mode"`
	MediaType                                                       types.String `tfsdk:"media_type"`
	CapacityAlertHighThreshold                                      types.Int64  `tfsdk:"capacity_alert_high_threshold"`
	CapacityAlertCriticalThreshold                                  types.Int64  `tfsdk:"capacity_alert_critical_threshold"`
	VtreeMigrationIoPriorityAppIopsPerDeviceThreshold               types.Int64  `tfsdk:"vtree_migration_io_priority_app_iops_per_device_threshold"`
	VtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps           types.Int64  `tfsdk:"vtree_migration_io_priority_app_bw_per_device_threshold_kbps"`
	VtreeMigrationIoPriorityQuietPeriodInMsec                       types.Int64  `tfsdk:"vtree_migration_io_priority_quiet_period_msec"`
	FglAccpID                                                       types.String `tfsdk:"fgl_accp_id"`
	FglExtraCapacity                                                types.Int64  `tfsdk:"fgl_extra_capacity"`
	FglOverProvisioningFactor                                       types.Int64  `tfsdk:"fgl_overprovisioning_factor"`
	FglWriteAtomicitySize                                           types.Int64  `tfsdk:"fgl_write_atomicity_size"`
	FglNvdimmWriteCacheSizeInMb                                     types.Int64  `tfsdk:"fgl_nvdimm_write_cache_size_mb"`
	FglNvdimmMetadataAmortizationX100                               types.Int64  `tfsdk:"fgl_nvdimm_metadata_amotization_x100"`
	FglPerfProfile                                                  types.String `tfsdk:"fgl_perf_profile"`
	ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps        types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_bw_limit_per_device_kbps"`
	ProtectedMaintenanceModeIoPriorityAppIopsPerDeviceThreshold     types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_app_iops_per_device_threshold"`
	ProtectedMaintenanceModeIoPriorityAppBwPerDeviceThresholdInKbps types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps"`
	ProtectedMaintenanceModeIoPriorityQuietPeriodInMsec             types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_quiet_period_msec"`
	Volumes                                                         []Volume     `tfsdk:"volumes"`
	SDS                                                             []SdsData    `tfsdk:"sds"`
	Links                                                           []LinkModel  `tfsdk:"links"`
}

// StoragepoolDataSourceModel maps the storage pool data source schema data
type StoragepoolDataSourceModel struct {
	StoragePoolIDs       types.List         `tfsdk:"storage_pool_ids"`
	StoragePoolNames     types.List         `tfsdk:"storage_pool_names"`
	ProtectionDomainID   types.String       `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String       `tfsdk:"protection_domain_name"`
	StoragePools         []StoragePoolModel `tfsdk:"storage_pools"`
	ID                   types.String       `tfsdk:"id"`
}
