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
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProtectionDomainResourceModel defines struct for protection domain resource
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

// ProtectionDomainDataSourceModel defines struct for protection domain data source
type ProtectionDomainDataSourceModel struct {
	ProtectionDomains      []ProtectionDomainModel `tfsdk:"protection_domains"`
	ID                     types.String            `tfsdk:"id"`
	ProtectionDomainFilter *ProtectionDomainFilter `tfsdk:"filter"`
}

// ProtectionDomainFilter defines struct for protection domain filter
type ProtectionDomainFilter struct {
	SystemID                                         []types.String `tfsdk:"system_id"`
	ReplicationCapacityMaxRatio                      []types.Int64  `tfsdk:"replication_capacity_max_ratio"`
	RebuildNetworkThrottlingInKbps                   []types.Int64  `tfsdk:"rebuild_network_throttling_in_kbps"`
	RebalanceNetworkThrottlingInKbps                 []types.Int64  `tfsdk:"rebalance_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingInKbps                 []types.Int64  `tfsdk:"overall_io_network_throttling_in_kbps"`
	VTreeMigrationNetworkThrottlingInKbps            []types.Int64  `tfsdk:"vtree_migration_network_throttling_in_kbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps  []types.Int64  `tfsdk:"protected_maintenance_mode_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingEnabled                types.Bool     `tfsdk:"overall_io_network_throttling_enabled"`
	RebuildNetworkThrottlingEnabled                  types.Bool     `tfsdk:"rebuild_network_throttling_enabled"`
	RebalanceNetworkThrottlingEnabled                types.Bool     `tfsdk:"rebalance_network_throttling_enabled"`
	VTreeMigrationNetworkThrottlingEnabled           types.Bool     `tfsdk:"vtree_migration_network_throttling_enabled"`
	ProtectedMaintenanceModeNetworkThrottlingEnabled types.Bool     `tfsdk:"protected_maintenance_mode_network_throttling_enabled"`
	FglDefaultNumConcurrentWrites                    []types.Int64  `tfsdk:"fgl_default_num_concurrent_writes"`
	FglMetadataCacheEnabled                          types.Bool     `tfsdk:"fgl_metadata_cache_enabled"`
	FglDefaultMetadataCacheSize                      []types.Int64  `tfsdk:"fgl_default_metadata_cache_size"`
	RfCacheEnabled                                   types.Bool     `tfsdk:"rf_cache_enabled"`
	RfCacheAccpID                                    []types.String `tfsdk:"rf_cache_accp_id"`
	RfCacheOperationalMode                           []types.String `tfsdk:"rf_cache_opertional_mode"`
	RfCachePageSizeKb                                []types.Int64  `tfsdk:"rf_cache_page_size_kb"`
	RfCacheMaxIoSizeKb                               []types.Int64  `tfsdk:"rf_cache_max_io_size_kb"`
	State                                            []types.String `tfsdk:"state"`
	Name                                             []types.String `tfsdk:"name"`
	ID                                               []types.String `tfsdk:"id"`
}

// ProtectionDomainModel defines struct for protection domain data source
type ProtectionDomainModel struct {
	SystemID                    types.String    `tfsdk:"system_id"`
	SdrSdsConnectivityInfo      PdConnInfoModel `tfsdk:"sdr_sds_connectivity"`
	ReplicationCapacityMaxRatio types.Int64     `tfsdk:"replication_capacity_max_ratio"`

	// Network throttling params
	RebuildNetworkThrottlingInKbps                   types.Int64 `tfsdk:"rebuild_network_throttling_in_kbps"`
	RebalanceNetworkThrottlingInKbps                 types.Int64 `tfsdk:"rebalance_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingInKbps                 types.Int64 `tfsdk:"overall_io_network_throttling_in_kbps"`
	VTreeMigrationNetworkThrottlingInKbps            types.Int64 `tfsdk:"vtree_migration_network_throttling_in_kbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps  types.Int64 `tfsdk:"protected_maintenance_mode_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingEnabled                types.Bool  `tfsdk:"overall_io_network_throttling_enabled"`
	RebuildNetworkThrottlingEnabled                  types.Bool  `tfsdk:"rebuild_network_throttling_enabled"`
	RebalanceNetworkThrottlingEnabled                types.Bool  `tfsdk:"rebalance_network_throttling_enabled"`
	VTreeMigrationNetworkThrottlingEnabled           types.Bool  `tfsdk:"vtree_migration_network_throttling_enabled"`
	ProtectedMaintenanceModeNetworkThrottlingEnabled types.Bool  `tfsdk:"protected_maintenance_mode_network_throttling_enabled"`

	// Fine Granularity Params
	FglDefaultNumConcurrentWrites types.Int64 `tfsdk:"fgl_default_num_concurrent_writes"`
	FglMetadataCacheEnabled       types.Bool  `tfsdk:"fgl_metadata_cache_enabled"`
	FglDefaultMetadataCacheSize   types.Int64 `tfsdk:"fgl_default_metadata_cache_size"`

	// RfCache Params
	RfCacheEnabled         types.Bool   `tfsdk:"rf_cache_enabled"`
	RfCacheAccpID          types.String `tfsdk:"rf_cache_accp_id"`
	RfCacheOperationalMode types.String `tfsdk:"rf_cache_opertional_mode"`
	RfCachePageSizeKb      types.Int64  `tfsdk:"rf_cache_page_size_kb"`
	RfCacheMaxIoSizeKb     types.Int64  `tfsdk:"rf_cache_max_io_size_kb"`

	// Counter Params
	SdsConfigurationFailureCP            PdCounterModel `tfsdk:"sds_configuration_failure_counter"`
	SdsDecoupledCP                       PdCounterModel `tfsdk:"sds_decoupled_counter"`
	MdmSdsNetworkDisconnectionsCP        PdCounterModel `tfsdk:"mdm_sds_network_disconnections_counter"`
	SdsSdsNetworkDisconnectionsCP        PdCounterModel `tfsdk:"sds_sds_network_disconnections_counter"`
	SdsReceiveBufferAllocationFailuresCP PdCounterModel `tfsdk:"sds_receive_buffer_allocation_failures_counter"`

	State types.String                `tfsdk:"state"`
	Name  types.String                `tfsdk:"name"`
	ID    types.String                `tfsdk:"id"`
	Links []ProtectionDomainLinkModel `tfsdk:"links"`
}

// WindowModel defines struct for protection domain window model
type WindowModel struct {
	Threshold       types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec types.Int64 `tfsdk:"window_size_in_sec"`
}

// PdCounterModel defines struct for protection domain counter models
type PdCounterModel struct {
	ShortWindow  WindowModel `tfsdk:"short_window"`
	MediumWindow WindowModel `tfsdk:"medium_window"`
	LongWindow   WindowModel `tfsdk:"long_window"`
}

// PdCounterModelValue defines struct for protection domain model values
func PdCounterModelValue(p scaleiotypes.PDCounterParams) PdCounterModel {
	return PdCounterModel{
		ShortWindow: WindowModel{
			Threshold:       types.Int64Value(int64(p.ShortWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.ShortWindow.WindowSizeInSec)),
		},
		MediumWindow: WindowModel{
			Threshold:       types.Int64Value(int64(p.MediumWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.MediumWindow.WindowSizeInSec)),
		},
		LongWindow: WindowModel{
			Threshold:       types.Int64Value(int64(p.LongWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.LongWindow.WindowSizeInSec)),
		},
	}
}

// PdConnInfoModel defines struct for protection domain connection information
type PdConnInfoModel struct {
	ClientServerConnStatus types.String `tfsdk:"client_server_conn_status"`
	DisconnectedClientID   types.String `tfsdk:"disconnected_client_id"`
	DisconnectedClientName types.String `tfsdk:"disconnected_client_name"`
	DisconnectedServerID   types.String `tfsdk:"disconnected_server_id"`
	DisconnectedServerName types.String `tfsdk:"disconnected_server_name"`
	DisconnectedServerIP   types.String `tfsdk:"disconnected_server_ip"`
}

// ProtectionDomainLinkModel defines struct for protection domain links
type ProtectionDomainLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}
