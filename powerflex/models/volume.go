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

// VolumeResourceModel maps the resource schema data.
type VolumeResourceModel struct {
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	StoragePoolName      types.String `tfsdk:"storage_pool_name"`
	StoragePoolID        types.String `tfsdk:"storage_pool_id"`
	VolumeType           types.String `tfsdk:"volume_type"`
	UseRmCache           types.Bool   `tfsdk:"use_rm_cache"`
	CompressionMethod    types.String `tfsdk:"compression_method"`
	Size                 types.Int64  `tfsdk:"size"`
	CapacityUnit         types.String `tfsdk:"capacity_unit"`
	Name                 types.String `tfsdk:"name"`
	SizeInKb             types.Int64  `tfsdk:"size_in_kb"`
	ID                   types.String `tfsdk:"id"`
	AccessMode           types.String `tfsdk:"access_mode"`
	RemoveMode           types.String `tfsdk:"remove_mode"`
}

// SDCItemize maps the sdc_list schema data
type SDCItemize struct {
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName       types.String `tfsdk:"sdc_name"`
	AccessMode    types.String `tfsdk:"access_mode"`
}

// VolumeDataSourceModel defines struct for volume data source
type VolumeDataSourceModel struct {
	Volumes      []VolumeModel `tfsdk:"volumes"`
	ID           types.String  `tfsdk:"id"`
	VolumeFilter *VolumeFilter `tfsdk:"filter"`
}

// VolumeModel define struct for volume model
type VolumeModel struct {
	ID                                 types.String         `tfsdk:"id"`
	Name                               types.String         `tfsdk:"name"`
	CreationTime                       types.Int64          `tfsdk:"creation_time"`
	SizeInKb                           types.Int64          `tfsdk:"size_in_kb"`
	AncestorVolumeID                   types.String         `tfsdk:"ancestor_volume_id"`
	VTreeID                            types.String         `tfsdk:"vtree_id"`
	ConsistencyGroupID                 types.String         `tfsdk:"consistency_group_id"`
	VolumeType                         types.String         `tfsdk:"volume_type"`
	UseRmCache                         types.Bool           `tfsdk:"use_rm_cache"`
	StoragePoolID                      types.String         `tfsdk:"storage_pool_id"`
	DataLayout                         types.String         `tfsdk:"data_layout"`
	NotGenuineSnapshot                 types.Bool           `tfsdk:"not_genuine_snapshot"`
	AccessModeLimit                    types.String         `tfsdk:"access_mode_limit"`
	SecureSnapshotExpTime              types.Int64          `tfsdk:"secure_snapshot_exp_time"`
	ManagedBy                          types.String         `tfsdk:"managed_by"`
	LockedAutoSnapshot                 types.Bool           `tfsdk:"locked_auto_snapshot"`
	LockedAutoSnapshotMarkedForRemoval types.Bool           `tfsdk:"locked_auto_snapshot_marked_for_removal"`
	CompressionMethod                  types.String         `tfsdk:"compression_method"`
	TimeStampIsAccurate                types.Bool           `tfsdk:"time_stamp_is_accurate"`
	OriginalExpiryTime                 types.Int64          `tfsdk:"original_expiry_time"`
	VolumeReplicationState             types.String         `tfsdk:"volume_replication_state"`
	ReplicationJournalVolume           types.Bool           `tfsdk:"replication_journal_volume"`
	ReplicationTimeStamp               types.Int64          `tfsdk:"replication_time_stamp"`
	Links                              []VolumeLinkModel    `tfsdk:"links"`
	MappedSdcInfo                      []MappedSdcInfoModel `tfsdk:"mapped_sdc_info"`
}

// VolumeFilter define struct for volume filter model
type VolumeFilter struct {
	ID                                 []types.String `tfsdk:"id"`
	Name                               []types.String `tfsdk:"name"`
	CreationTime                       []types.Int64  `tfsdk:"creation_time"`
	SizeInKb                           []types.Int64  `tfsdk:"size_in_kb"`
	AncestorVolumeID                   []types.String `tfsdk:"ancestor_volume_id"`
	VTreeID                            []types.String `tfsdk:"vtree_id"`
	ConsistencyGroupID                 []types.String `tfsdk:"consistency_group_id"`
	VolumeType                         []types.String `tfsdk:"volume_type"`
	UseRmCache                         types.Bool     `tfsdk:"use_rm_cache"`
	StoragePoolID                      []types.String `tfsdk:"storage_pool_id"`
	DataLayout                         []types.String `tfsdk:"data_layout"`
	NotGenuineSnapshot                 types.Bool     `tfsdk:"not_genuine_snapshot"`
	AccessModeLimit                    []types.String `tfsdk:"access_mode_limit"`
	SecureSnapshotExpTime              []types.Int64  `tfsdk:"secure_snapshot_exp_time"`
	ManagedBy                          []types.String `tfsdk:"managed_by"`
	LockedAutoSnapshot                 types.Bool     `tfsdk:"locked_auto_snapshot"`
	LockedAutoSnapshotMarkedForRemoval types.Bool     `tfsdk:"locked_auto_snapshot_marked_for_removal"`
	CompressionMethod                  []types.String `tfsdk:"compression_method"`
	TimeStampIsAccurate                types.Bool     `tfsdk:"time_stamp_is_accurate"`
	OriginalExpiryTime                 []types.Int64  `tfsdk:"original_expiry_time"`
	VolumeReplicationState             []types.String `tfsdk:"volume_replication_state"`
	ReplicationJournalVolume           types.Bool     `tfsdk:"replication_journal_volume"`
	ReplicationTimeStamp               []types.Int64  `tfsdk:"replication_time_stamp"`
}

// VolumeLinkModel defines struct for volume links
type VolumeLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// MappedSdcInfoModel defines struct for mapped SDC info
type MappedSdcInfoModel struct {
	SdcID                 types.String `tfsdk:"sdc_id"`
	SdcIP                 types.String `tfsdk:"sdc_ip"`
	LimitIops             types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps         types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName               types.String `tfsdk:"sdc_name"`
	AccessMode            types.String `tfsdk:"access_mode"`
	IsDirectBufferMapping types.Bool   `tfsdk:"is_direct_buffer_mapping"`
}
