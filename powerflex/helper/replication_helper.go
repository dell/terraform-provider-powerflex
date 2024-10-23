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
	"fmt"
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

// CreateReplicationPair POST replication pair
func CreateReplicationPair(client *goscaleio.Client, plan models.ReplicationPairResourceModel) (string, error) {
	rp := &scaleiotypes.QueryReplicationPair{
		Name:                          plan.Name.ValueString(),
		SourceVolumeID:                plan.SourceVolumeID.ValueString(),
		DestinationVolumeID:           plan.DestinationVolumeID.ValueString(),
		ReplicationConsistencyGroupID: plan.ReplicationConsistencyGroupID.ValueString(),
		// OnlineCopy is the only supported copy type for replication pair
		CopyType: "OnlineCopy",
	}
	res, err := client.CreateReplicationPair(rp)
	if err != nil {
		return "", err
	}
	return res.ID, err
}

// PauseReplicationPair Pause initial replication pair
func PauseReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.PausePairInitialCopy(id)
}

// ResumeReplicationPair Resume initial replication pair
func ResumeReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.ResumePairInitialCopy(id)
}

// GetSpecificReplicationPair GET a replication pair
func GetSpecificReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.GetReplicationPair(id)
}

// MapReplicationPairState map single replication pair state
func MapReplicationPairState(val scaleiotypes.ReplicationPair, state models.ReplicationPairResourceModel) models.ReplicationPairResourceModel {
	state.ID = types.StringValue(val.ID)
	state.Name = types.StringValue(val.Name)
	state.RemoteID = types.StringValue(val.RemoteID)
	state.UserRequestedPauseTransmitInitCopy = types.BoolValue(val.UserRequestedPauseTransmitInitCopy)
	state.RemoteCapacityInMB = types.Int64Value(int64(val.RemoteCapacityInMB))
	state.LocalVolumeID = types.StringValue(val.LocalVolumeID)
	state.RemoteVolumeID = types.StringValue(val.RemoteID)
	state.RemoteVolumeName = types.StringValue(val.RemoteVolumeName)
	state.ReplicationConsistencyGroupID = types.StringValue(val.ReplicationConsistencyGroupID)
	state.CopyType = types.StringValue(val.LifetimeState)
	state.LifetimeState = types.StringValue(val.CopyType)
	state.PeerSystemName = types.StringValue(val.LifetimeState)
	state.InitialCopyState = types.StringValue(val.InitialCopyState)
	state.InitialCopyPriority = types.Int64Value(int64(val.InitialCopyPriority))
	return state
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

// GetReplicationConsistancyGroups GET RCGs
func GetReplicationConsistancyGroups(client *goscaleio.Client) ([]scaleiotypes.ReplicationConsistencyGroup, error) {
	rps := []scaleiotypes.ReplicationConsistencyGroup{}

	// Get All RCGs
	rcgs, err := client.GetReplicationConsistencyGroups()
	if err != nil {
		return nil, err
	}
	for _, val := range rcgs {
		rps = append(rps, *val)
	}

	return rps, nil
}

// GetSpecificReplicationConsistencyGroup GET a specific RCG
func GetSpecificReplicationConsistencyGroup(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationConsistencyGroup, error) {
	return client.GetReplicationConsistencyGroupByID(id)
}

// CreateReplicationConsistencyGroup POST replication consistency group
func CreateReplicationConsistencyGroup(client *goscaleio.Client, plan models.ReplicationConsistancyGroupModel) (string, error) {
	rcg := scaleiotypes.ReplicationConsistencyGroupCreatePayload{
		Name:                     plan.Name.ValueString(),
		RpoInSeconds:             plan.RpoInSeconds.String(),
		ProtectionDomainID:       plan.ProtectionDomainID.ValueString(),
		RemoteProtectionDomainID: plan.RemoteProtectionDomainID.ValueString(),
		DestinationSystemID:      plan.DestinationSystemID.ValueString(),
	}
	res, err := client.CreateReplicationConsistencyGroup(&rcg)
	if err != nil {
		return "", err
	}
	return res.ID, err
}

// MapReplicationConsistancyGroupsState map Replication Consistancy Groups state
func MapReplicationConsistancyGroupsState(rcgs []scaleiotypes.ReplicationConsistencyGroup, state models.ReplicationConsistancyGroupDataSourceModel) models.ReplicationConsistancyGroupDataSourceModel {
	mappedRps := []models.ReplicationConsistancyGroupModel{}
	for _, val := range rcgs {
		temp := models.ReplicationConsistancyGroupModel{
			ID:                          types.StringValue(val.ID),
			Name:                        types.StringValue(val.Name),
			RemoteID:                    types.StringValue(val.RemoteID),
			RpoInSeconds:                types.Int64Value(int64(val.RpoInSeconds)),
			ProtectionDomainID:          types.StringValue(val.ProtectionDomainID),
			RemoteProtectionDomainID:    types.StringValue(val.RemoteProtectionDomainID),
			DestinationSystemID:         types.StringValue(val.DestinationSystemID),
			PeerMdmID:                   types.StringValue(val.PeerMdmID),
			RemoteMdmID:                 types.StringValue(val.RemoteMdmID),
			ReplicationDirection:        types.StringValue(val.ReplicationDirection),
			CurrConsistMode:             types.StringValue(val.CurrConsistMode),
			FreezeState:                 types.StringValue(val.FreezeState),
			PauseMode:                   types.StringValue(val.PauseMode),
			LifetimeState:               types.StringValue(val.LifetimeState),
			SnapCreationInProgress:      types.BoolValue(val.SnapCreationInProgress),
			LastSnapGroupID:             types.StringValue(val.LastSnapGroupID),
			Type:                        types.StringValue(val.Type),
			DisasterRecoveryState:       types.StringValue(val.DisasterRecoveryState),
			RemoteDisasterRecoveryState: types.StringValue(val.RemoteDisasterRecoveryState),
			TargetVolumeAccessMode:      types.StringValue(val.TargetVolumeAccessMode),
			FailoverType:                types.StringValue(val.FailoverType),
			FailoverState:               types.StringValue(val.FailoverState),
			ActiveLocal:                 types.BoolValue(val.ActiveLocal),
			ActiveRemote:                types.BoolValue(val.ActiveRemote),
			AbstractState:               types.StringValue(val.AbstractState),
			Error:                       types.Int64Value(int64(val.Error)),
			LocalActivityState:          types.StringValue(val.LocalActivityState),
			RemoteActivityState:         types.StringValue(val.RemoteActivityState),
			InactiveReason:              types.Int64Value(int64(val.InactiveReason)),
		}
		mappedRps = append(mappedRps, temp)
	}
	return models.ReplicationConsistancyGroupDataSourceModel{
		ID:                                 types.StringValue("replication_consistancy_group_id"),
		ReplicationConsistancyGroupFilter:  state.ReplicationConsistancyGroupFilter,
		ReplicationConsistancyGroupDetails: mappedRps,
	}
}

// MapReplicationConsistancyGroupsResourceState map Replication Consistancy Groups state
func MapReplicationConsistancyGroupsResourceState(rcg scaleiotypes.ReplicationConsistencyGroup, state models.ReplicationConsistancyGroupModel) models.ReplicationConsistancyGroupModel {
	rcgMap := models.ReplicationConsistancyGroupModel{
		ID:                          types.StringValue(rcg.ID),
		Name:                        types.StringValue(rcg.Name),
		RemoteID:                    types.StringValue(rcg.RemoteID),
		RpoInSeconds:                state.RpoInSeconds,
		ProtectionDomainID:          types.StringValue(rcg.ProtectionDomainID),
		RemoteProtectionDomainID:    types.StringValue(rcg.RemoteProtectionDomainID),
		DestinationSystemID:         state.DestinationSystemID,
		PeerMdmID:                   types.StringValue(rcg.PeerMdmID),
		RemoteMdmID:                 types.StringValue(rcg.RemoteMdmID),
		ReplicationDirection:        types.StringValue(rcg.ReplicationDirection),
		CurrConsistMode:             state.CurrConsistMode,
		FreezeState:                 state.FreezeState,
		PauseMode:                   state.PauseMode,
		LifetimeState:               types.StringValue(rcg.LifetimeState),
		SnapCreationInProgress:      types.BoolValue(rcg.SnapCreationInProgress),
		LastSnapGroupID:             types.StringValue(rcg.LastSnapGroupID),
		Type:                        types.StringValue(rcg.Type),
		DisasterRecoveryState:       types.StringValue(rcg.DisasterRecoveryState),
		RemoteDisasterRecoveryState: types.StringValue(rcg.RemoteDisasterRecoveryState),
		TargetVolumeAccessMode:      state.TargetVolumeAccessMode,
		FailoverType:                types.StringValue(rcg.FailoverType),
		FailoverState:               types.StringValue(rcg.FailoverState),
		ActiveLocal:                 types.BoolValue(rcg.ActiveLocal),
		ActiveRemote:                types.BoolValue(rcg.ActiveRemote),
		AbstractState:               types.StringValue(rcg.AbstractState),
		Error:                       types.Int64Value(int64(rcg.Error)),
		LocalActivityState:          state.LocalActivityState,
		RemoteActivityState:         types.StringValue(rcg.RemoteActivityState),
		InactiveReason:              types.Int64Value(int64(rcg.InactiveReason)),
	}
	return rcgMap
}

// RCGUpdates Update the RCG
func RCGUpdates(client *goscaleio.Client, state models.ReplicationConsistancyGroupModel, plan models.ReplicationConsistancyGroupModel) error {
	rcgClient := goscaleio.NewReplicationConsistencyGroup(client)
	rcgClient.ReplicationConsistencyGroup.ID = state.ID.ValueString()
	// Update RPO
	if state.RpoInSeconds.ValueInt64() != plan.RpoInSeconds.ValueInt64() {
		rpoErr := rcgClient.SetRPOOnReplicationGroup(scaleiotypes.SetRPOReplicationConsistencyGroup{
			RpoInSeconds: fmt.Sprint(plan.RpoInSeconds.ValueInt64()),
		})
		if rpoErr != nil {
			return rpoErr
		}
	}

	// Update Activity State
	if state.LocalActivityState.ValueString() != plan.LocalActivityState.ValueString() {
		if plan.LocalActivityState.ValueString() == "Active" {
			activateErr := rcgClient.ExecuteActivateOnReplicationGroup()
			if activateErr != nil {
				return activateErr
			}
		} else {
			inactivateErr := rcgClient.ExecuteTerminateOnReplicationGroup()
			if inactivateErr != nil {
				return inactivateErr
			}
		}

	}

	// Update Access Mode
	if state.TargetVolumeAccessMode.ValueString() != plan.TargetVolumeAccessMode.ValueString() {
		vamErr := rcgClient.SetTargetVolumeAccessModeOnReplicationGroup(scaleiotypes.SetTargetVolumeAccessModeOnReplicationGroup{
			TargetVolumeAccessMode: plan.TargetVolumeAccessMode.ValueString(),
		})
		if vamErr != nil {
			return vamErr
		}
	}

	// Update Pause Mode
	if state.PauseMode.ValueString() != plan.PauseMode.ValueString() {
		if plan.PauseMode.ValueString() == "Pause" {
			pauseModeErr := rcgClient.ExecutePauseOnReplicationGroup()
			if pauseModeErr != nil {
				return pauseModeErr
			}
		} else {
			resumeModeErr := rcgClient.ExecuteResumeOnReplicationGroup()
			if resumeModeErr != nil {
				return resumeModeErr
			}
		}
	}

	// Update Freeze State
	if state.FreezeState.ValueString() != plan.FreezeState.ValueString() {
		if plan.FreezeState.ValueString() == "Frozen" {
			freezeErr := rcgClient.FreezeReplicationConsistencyGroup(state.ID.ValueString())
			if freezeErr != nil {
				return freezeErr
			}
		} else {
			unfreezeErr := rcgClient.UnfreezeReplicationConsistencyGroup()
			if unfreezeErr != nil {
				return unfreezeErr
			}
		}
	}

	// Update Consistency Mode
	if plan.CurrConsistMode.ValueString() != state.CurrConsistMode.ValueString() {
		if plan.CurrConsistMode.ValueString() == "Consistent" {
			consistentErr := rcgClient.ExecuteConsistentOnReplicationGroup()
			if consistentErr != nil {
				return consistentErr
			}
		} else {
			inconsistentErr := rcgClient.ExecuteInconsistentOnReplicationGroup()
			if inconsistentErr != nil {
				return inconsistentErr
			}
		}
	}

	// Update Name
	if state.Name.ValueString() != plan.Name.ValueString() {
		nameErr := rcgClient.SetNewNameOnReplicationGroup(scaleiotypes.SetNewNameOnReplicationGroup{
			NewName: plan.Name.ValueString(),
		})
		if nameErr != nil {
			return nameErr
		}
	}
	return nil
}
