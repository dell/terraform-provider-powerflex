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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateSnapshotPolicyState iterates over the snapshotpolicy list and update the state
func UpdateSnapshotPolicyState(sps []*scaleiotypes.SnapshotPolicy) (response []models.SnapshotPolicyModel) {
	for _, sp := range sps {
		spState := models.SnapshotPolicyModel{
			ID:                                    types.StringValue(sp.ID),
			Name:                                  types.StringValue(sp.Name),
			SnapshotPolicyState:                   types.StringValue(sp.SnapshotPolicyState),
			AutoSnapshotCreationCadenceInMin:      types.Int64Value((int64)(sp.AutoSnapshotCreationCadenceInMin)),
			MaxVTreeAutoSnapshots:                 types.Int64Value((int64)(sp.MaxVTreeAutoSnapshots)),
			NumOfSourceVolumes:                    types.Int64Value((int64)(sp.NumOfSourceVolumes)),
			NumOfExpiredButLockedSnapshots:        types.Int64Value((int64)(sp.NumOfExpiredButLockedSnapshots)),
			NumOfCreationFailures:                 types.Int64Value((int64)(sp.NumOfCreationFailures)),
			SnapshotAccessMode:                    types.StringValue(sp.SnapshotAccessMode),
			SecureSnapshots:                       types.BoolValue(sp.SecureSnapshots),
			TimeOfLastAutoSnapshot:                types.Int64Value((int64)(sp.TimeOfLastAutoSnapshot)),
			NextAutoSnapshotCreationTime:          types.Int64Value((int64)(sp.NextAutoSnapshotCreationTime)),
			TimeOfLastAutoSnapshotCreationFailure: types.Int64Value((int64)(sp.TimeOfLastAutoSnapshotCreationFailure)),
			LastAutoSnapshotCreationFailureReason: types.StringValue(sp.LastAutoSnapshotCreationFailureReason),
			LastAutoSnapshotFailureInFirstLevel:   types.BoolValue(sp.LastAutoSnapshotFailureInFirstLevel),
			NumOfAutoSnapshots:                    types.Int64Value((int64)(sp.NumOfAutoSnapshots)),
			NumOfLockedSnapshots:                  types.Int64Value((int64)(sp.NumOfLockedSnapshots)),
			SystemID:                              types.StringValue(sp.SystemID),
		}

		for _, link := range sp.Links {
			spState.Links = append(spState.Links, models.SnapshotPolicyLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		for _, rspl := range sp.NumOfRetainedSnapshotsPerLevel {
			spState.NumOfRetainedSnapshotsPerLevel = append(spState.NumOfRetainedSnapshotsPerLevel, types.Int64Value((int64)(rspl)))
		}
		response = append(response, spState)
	}
	return
}
