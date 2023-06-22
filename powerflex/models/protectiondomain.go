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

type ProtectionDomainResourceModel struct {
	ReplicationCapacityMaxRatio types.Int64 `tfsdk:"replication_capacity_max_ratio"`

	// Network throttling params
	RebuildNetworkThrottlingInKbps                  types.Int64 `tfsdk:"rebuild_network_throttling_in_kbps"`
	RebalanceNetworkThrottlingInKbps                types.Int64 `tfsdk:"rebalance_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingInKbps                types.Int64 `tfsdk:"overall_io_network_throttling_in_kbps"`
	VTreeMigrationNetworkThrottlingInKbps           types.Int64 `tfsdk:"vtree_migration_network_throttling_in_kbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps types.Int64 `tfsdk:"protected_maintenance_mode_network_throttling_in_kbps"`

	// Fine Granularity Params
	FglDefaultNumConcurrentWrites types.Int64 `tfsdk:"fgl_default_num_concurrent_writes"`
	FglMetadataCacheEnabled       types.Bool  `tfsdk:"fgl_metadata_cache_enabled"`
	FglDefaultMetadataCacheSize   types.Int64 `tfsdk:"fgl_default_metadata_cache_size"`

	// RfCache Params
	RfCacheEnabled         types.Bool   `tfsdk:"rf_cache_enabled"`
	RfCacheAccpID          types.String `tfsdk:"rf_cache_accp_id"`
	RfCacheOperationalMode types.String `tfsdk:"rf_cache_operational_mode"`
	RfCachePageSizeKb      types.Int64  `tfsdk:"rf_cache_page_size_kb"`
	RfCacheMaxIoSizeKb     types.Int64  `tfsdk:"rf_cache_max_io_size_kb"`

	Active types.Bool   `tfsdk:"active"`
	State  types.String `tfsdk:"state"`
	Name   types.String `tfsdk:"name"`
	ID     types.String `tfsdk:"id"`
	Links  types.List   `tfsdk:"links"`
}
