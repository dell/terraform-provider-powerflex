/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

import "github.com/hashicorp/terraform-plugin-framework/types"

// ReplicationPairModel model for the replication pair
type ReplicationPairModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	RemoteID                           types.String `tfsdk:"remote_id"`
	UserRequestedPauseTransmitInitCopy types.Bool   `tfsdk:"user_requested_pause_transmit_init_copy"`
	RemoteCapacityInMB                 types.Int64  `tfsdk:"remote_capacity_in_mb"`
	LocalVolumeID                      types.String `tfsdk:"local_volume_id"`
	RemoteVolumeID                     types.String `tfsdk:"remote_volume_id"`
	RemoteVolumeName                   types.String `tfsdk:"remote_volume_name"`
	ReplicationConsistencyGroupID      types.String `tfsdk:"replication_consistency_group_id"`
	CopyType                           types.String `tfsdk:"copy_type"`
	LifetimeState                      types.String `tfsdk:"lifetime_state"`
	PeerSystemName                     types.String `tfsdk:"peer_system_name"`
	InitialCopyState                   types.String `tfsdk:"initial_copy_state"`
	InitialCopyPriority                types.Int64  `tfsdk:"initial_copy_priority"`
}

// ReplicationPairDataSourceModel is the tfsdk model of ReplicationPair data source schema
type ReplicationPairDataSourceModel struct {
	ReplicationPairDetails []ReplicationPairModel `tfsdk:"replication_pair_details"`
	ID                     types.String           `tfsdk:"id"`
	ReplicationPairFilter  *ReplicationPairFilter `tfsdk:"filter"`
}

// ReplicationPairFilter defines the model for filters used for ReplicationPairsDataSource
type ReplicationPairFilter struct {
	ID                                 []types.String `tfsdk:"id"`
	Name                               []types.String `tfsdk:"name"`
	RemoteID                           []types.String `tfsdk:"remote_id"`
	UserRequestedPauseTransmitInitCopy types.Bool     `tfsdk:"user_requested_pause_transmit_init_copy"`
	RemoteCapacityInMB                 []types.Int64  `tfsdk:"remote_capacity_in_mb"`
	LocalVolumeID                      []types.String `tfsdk:"local_volume_id"`
	RemoteVolumeID                     []types.String `tfsdk:"remote_volume_id"`
	RemoteVolumeName                   []types.String `tfsdk:"remote_volume_name"`
	ReplicationConsistencyGroupID      []types.String `tfsdk:"replication_consistency_group_id"`
	CopyType                           []types.String `tfsdk:"copy_type"`
	LifetimeState                      []types.String `tfsdk:"lifetime_state"`
	PeerSystemName                     []types.String `tfsdk:"peer_system_name"`
	InitialCopyState                   []types.String `tfsdk:"initial_copy_state"`
	InitialCopyPriority                []types.Int64  `tfsdk:"initial_copy_priority"`
}

// ReplicationConsistancyGroupModel model for the RCG
type ReplicationConsistancyGroupModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	RpoInSeconds                types.Int64  `tfsdk:"rpo_in_seconds"`
	ProtectionDomainID          types.String `tfsdk:"protection_domain_id"`
	RemoteProtectionDomainID    types.String `tfsdk:"remote_protection_domain_id"`
	DestinationSystemID         types.String `tfsdk:"destination_system_id"`
	PeerMdmID                   types.String `tfsdk:"peer_mdm_id"`
	RemoteID                    types.String `tfsdk:"remote_id"`
	RemoteMdmID                 types.String `tfsdk:"remote_mdm_id"`
	ReplicationDirection        types.String `tfsdk:"replication_direction"`
	CurrConsistMode             types.String `tfsdk:"curr_consist_mode"`
	FreezeState                 types.String `tfsdk:"freeze_state"`
	PauseMode                   types.String `tfsdk:"pause_mode"`
	LifetimeState               types.String `tfsdk:"lifetime_state"`
	SnapCreationInProgress      types.Bool   `tfsdk:"snap_creation_in_progress"`
	LastSnapGroupID             types.String `tfsdk:"last_snap_group_id"`
	Type                        types.String `tfsdk:"type"`
	DisasterRecoveryState       types.String `tfsdk:"disaster_recovery_state"`
	RemoteDisasterRecoveryState types.String `tfsdk:"remote_disaster_recovery_state"`
	TargetVolumeAccessMode      types.String `tfsdk:"target_volume_access_mode"`
	FailoverType                types.String `tfsdk:"failover_type"`
	FailoverState               types.String `tfsdk:"failover_state"`
	ActiveLocal                 types.Bool   `tfsdk:"active_local"`
	ActiveRemote                types.Bool   `tfsdk:"active_remote"`
	AbstractState               types.String `tfsdk:"abstract_state"`
	Error                       types.Int64  `tfsdk:"error"`
	LocalActivityState          types.String `tfsdk:"local_activity_state"`
	RemoteActivityState         types.String `tfsdk:"remote_activity_state"`
	InactiveReason              types.Int64  `tfsdk:"inactive_reason"`
}

// ReplicationConsistancyGroupDataSourceModel is the tfsdk model of ReplicationConsistancyGroup data source schema
type ReplicationConsistancyGroupDataSourceModel struct {
	ReplicationConsistancyGroupDetails []ReplicationConsistancyGroupModel `tfsdk:"replication_consistency_group_details"`
	ID                                 types.String                       `tfsdk:"id"`
	ReplicationConsistancyGroupFilter  *ReplicationConsistancyGroupFilter `tfsdk:"filter"`
}

// ReplicationConsistancyGroupFilter defines the model for filters used for ReplicationConsistancyGroupsDataSource
type ReplicationConsistancyGroupFilter struct {
	ID                          []types.String `tfsdk:"id"`
	Name                        []types.String `tfsdk:"name"`
	RpoInSeconds                []types.Int64  `tfsdk:"rpo_in_seconds"`
	ProtectionDomainID          []types.String `tfsdk:"protection_domain_id"`
	RemoteProtectionDomainID    []types.String `tfsdk:"remote_protection_domain_id"`
	DestinationSystemID         []types.String `tfsdk:"destination_system_id"`
	PeerMdmID                   []types.String `tfsdk:"peer_mdm_id"`
	RemoteID                    []types.String `tfsdk:"remote_id"`
	RemoteMdmID                 []types.String `tfsdk:"remote_mdm_id"`
	ReplicationDirection        []types.String `tfsdk:"replication_direction"`
	CurrConsistMode             []types.String `tfsdk:"curr_consist_mode"`
	FreezeState                 []types.String `tfsdk:"freeze_state"`
	PauseMode                   []types.String `tfsdk:"pause_mode"`
	LifetimeState               []types.String `tfsdk:"lifetime_state"`
	SnapCreationInProgress      types.Bool     `tfsdk:"snap_creation_in_progress"`
	LastSnapGroupID             []types.String `tfsdk:"last_snap_group_id"`
	Type                        []types.String `tfsdk:"type"`
	DisasterRecoveryState       []types.String `tfsdk:"disaster_recovery_state"`
	RemoteDisasterRecoveryState []types.String `tfsdk:"remote_disaster_recovery_state"`
	TargetVolumeAccessMode      []types.String `tfsdk:"target_volume_access_mode"`
	FailoverType                []types.String `tfsdk:"failover_type"`
	FailoverState               []types.String `tfsdk:"failover_state"`
	ActiveLocal                 types.Bool     `tfsdk:"active_local"`
	ActiveRemote                types.Bool     `tfsdk:"active_remote"`
	AbstractState               []types.String `tfsdk:"abstract_state"`
	Error                       []types.Int64  `tfsdk:"error"`
	LocalActivityState          []types.String `tfsdk:"local_activity_state"`
	RemoteActivityState         []types.String `tfsdk:"remote_activity_state"`
	InactiveReason              []types.Int64  `tfsdk:"inactive_reason"`
}
