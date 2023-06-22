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
