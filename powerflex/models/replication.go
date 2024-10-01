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
