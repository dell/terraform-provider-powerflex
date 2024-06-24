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

// SnapshotPolicyDataSourceModel defines struct for snapshot policy data source
type SnapshotPolicyDataSourceModel struct {
	SnapshotPolicies []SnapshotPolicyModel `tfsdk:"snapshotpolicies"`
	ID               types.String          `tfsdk:"id"`
	Name             types.String          `tfsdk:"name"`
}

// SnapshotPolicyModel defines struct for snapshot policy model
type SnapshotPolicyModel struct {
	ID                                    types.String              `tfsdk:"id"`
	Name                                  types.String              `tfsdk:"name"`
	SnapshotPolicyState                   types.String              `tfsdk:"snapshot_policy_state"`
	AutoSnapshotCreationCadenceInMin      types.Int64               `tfsdk:"auto_snapshot_creation_cadence_in_min"`
	MaxVTreeAutoSnapshots                 types.Int64               `tfsdk:"max_vtree_auto_snapshots"`
	NumOfSourceVolumes                    types.Int64               `tfsdk:"num_of_source_volumes"`
	NumOfExpiredButLockedSnapshots        types.Int64               `tfsdk:"num_of_expired_but_locked_snapshots"`
	NumOfCreationFailures                 types.Int64               `tfsdk:"num_of_creation_failures"`
	NumOfRetainedSnapshotsPerLevel        []types.Int64             `tfsdk:"num_of_retained_snapshots_per_level"`
	SnapshotAccessMode                    types.String              `tfsdk:"snapshot_access_mode"`
	SecureSnapshots                       types.Bool                `tfsdk:"secure_snapshots"`
	TimeOfLastAutoSnapshot                types.Int64               `tfsdk:"time_of_last_auto_snapshot"`
	NextAutoSnapshotCreationTime          types.Int64               `tfsdk:"next_auto_snapshot_creation_time"`
	TimeOfLastAutoSnapshotCreationFailure types.Int64               `tfsdk:"time_of_last_auto_snapshot_creation_failure"`
	LastAutoSnapshotCreationFailureReason types.String              `tfsdk:"last_auto_snapshot_creation_failure_reason"`
	LastAutoSnapshotFailureInFirstLevel   types.Bool                `tfsdk:"last_auto_snapshot_failure_in_first_level"`
	NumOfAutoSnapshots                    types.Int64               `tfsdk:"num_of_auto_snapshots"`
	NumOfLockedSnapshots                  types.Int64               `tfsdk:"num_of_locked_snapshots"`
	SystemID                              types.String              `tfsdk:"system_id"`
	Links                                 []SnapshotPolicyLinkModel `tfsdk:"links"`
}

// SnapshotPolicyLinkModel defines struct for snapshot policy links
type SnapshotPolicyLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// SnapshotPolicyResourceModel defines struct for snapshot policy resource
type SnapshotPolicyResourceModel struct {
	ID                               types.String   `tfsdk:"id"`
	Name                             types.String   `tfsdk:"name"`
	NumOfRetainedSnapshotsPerLevel   []int64        `tfsdk:"num_of_retained_snapshots_per_level"`
	AutoSnapshotCreationCadenceInMin types.Int64    `tfsdk:"auto_snapshot_creation_cadence_in_min"`
	Paused                           types.Bool     `tfsdk:"paused"`
	VolumeIds                        []types.String `tfsdk:"volume_ids"`
	RemoveMode                       types.String   `tfsdk:"remove_mode"`
	SecureSnapshots                  types.Bool     `tfsdk:"secure_snapshots"`
	SnapshotAccessMode               types.String   `tfsdk:"snapshot_access_mode"`
}
