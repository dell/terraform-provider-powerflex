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

// ReplicationPairResourceModel model for the replication pair
type ReplicationPairResourceModel struct {
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
	SourceVolumeID                     types.String `tfsdk:"source_volume_id"`
	DestinationVolumeID                types.String `tfsdk:"destination_volume_id"`
	PauseCopy                          types.Bool   `tfsdk:"pause_initial_copy"`
}

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

// ReplicationConsistencyGroupModel model for the RCG
type ReplicationConsistencyGroupModel struct {
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

// ReplicationConsistencyGroupDataSourceModel is the tfsdk model of ReplicationConsistencyGroup data source schema
type ReplicationConsistencyGroupDataSourceModel struct {
	ReplicationConsistencyGroupDetails []ReplicationConsistencyGroupModel `tfsdk:"replication_consistency_group_details"`
	ID                                 types.String                       `tfsdk:"id"`
	ReplicationConsistencyGroupFilter  *ReplicationConsistencyGroupFilter `tfsdk:"filter"`
}

// ReplicationConsistencyGroupFilter defines the model for filters used for ReplicationConsistencyGroupsDataSource
type ReplicationConsistencyGroupFilter struct {
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

// ReplicationConsistencyGroupAction defines the model for ReplicationConsistencyGroupAction
type ReplicationConsistencyGroupAction struct {
	ID     types.String `tfsdk:"id"`
	Action types.String `tfsdk:"action"`
}
