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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetReplicationPairs GET replication pairs
func GetReplicationPairs(client *goscaleio.Client) ([]scaleiotypes.ReplicationPair, error) {
	rps := []scaleiotypes.ReplicationPair{}

	// Get All Replication Pairs
	pairs, err := client.GetAllReplicationPairs()
	if err != nil {
		return nil, err
	}
	for _, val := range pairs {
		rps = append(rps, *val)
	}

	return rps, nil
}

// MapReplicationPairsState map replication pairs state
func MapReplicationPairsState(pairs []scaleiotypes.ReplicationPair, state models.ReplicationPairDataSourceModel) models.ReplicationPairDataSourceModel {
	mappedRps := []models.ReplicationPairModel{}
	for _, val := range pairs {
		temp := models.ReplicationPairModel{
			ID:                                 types.StringValue(val.ID),
			Name:                               types.StringValue(val.Name),
			RemoteID:                           types.StringValue(val.RemoteID),
			UserRequestedPauseTransmitInitCopy: types.BoolValue(val.UserRequestedPauseTransmitInitCopy),
			RemoteCapacityInMB:                 types.Int64Value(int64(val.RemoteCapacityInMB)),
			LocalVolumeID:                      types.StringValue(val.LocalVolumeID),
			RemoteVolumeID:                     types.StringValue(val.RemoteID),
			RemoteVolumeName:                   types.StringValue(val.RemoteVolumeName),
			ReplicationConsistencyGroupID:      types.StringValue(val.ReplicationConsistencyGroupID),
			CopyType:                           types.StringValue(val.LifetimeState),
			LifetimeState:                      types.StringValue(val.CopyType),
			PeerSystemName:                     types.StringValue(val.LifetimeState),
			InitialCopyState:                   types.StringValue(val.InitialCopyState),
			InitialCopyPriority:                types.Int64Value(int64(val.InitialCopyPriority)),
		}
		mappedRps = append(mappedRps, temp)
	}
	return models.ReplicationPairDataSourceModel{
		ID:                     types.StringValue("replication_pair_id"),
		ReplicationPairFilter:  state.ReplicationPairFilter,
		ReplicationPairDetails: mappedRps,
	}
}
