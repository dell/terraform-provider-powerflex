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

package helper

import (
	"strconv"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateSnapshotPolicyState iterates over the snapshotpolicy list and update the state
func UpdateSnapshotPolicyState(sps []scaleiotypes.SnapshotPolicy) (response []models.SnapshotPolicyModel) {
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

// CreateSnapshotPolicy creates snapshot policy
func CreateSnapshotPolicy(system *goscaleio.System, payload *scaleiotypes.SnapshotPolicyCreateParam) (string, error) {
	return system.CreateSnapshotPolicy(payload)
}

// GetSnapshotPolicy gets snapshot policy
func GetSnapshotPolicy(client *goscaleio.Client, id string) ([]*scaleiotypes.SnapshotPolicy, error) {
	return client.GetSnapshotPolicy("", id)
}

// AssignVolumeToSnapshotPolicy assigns volume to snapshot policy
func AssignVolumeToSnapshotPolicy(system *goscaleio.System, payload *scaleiotypes.AssignVolumeToSnapshotPolicyParam, id string) error {
	return system.AssignVolumeToSnapshotPolicy(payload, id)
}

// ModifySnapshotPolicy modifies snapshot policy
func ModifySnapshotPolicy(system *goscaleio.System, payload *scaleiotypes.SnapshotPolicyModifyParam, id string) error {
	return system.ModifySnapshotPolicy(payload, id)
}

// RenameSnapshotPolicy modifies snapshot policy name
func RenameSnapshotPolicy(system *goscaleio.System, id string, name string) error {
	return system.RenameSnapshotPolicy(id, name)
}

// UpdateSnapshotPolicyResourceState updates the state file attributes for snapshot policy resource.
func UpdateSnapshotPolicyResourceState(sps []*scaleiotypes.SnapshotPolicy, volumes []*scaleiotypes.Volume, state *models.SnapshotPolicyResourceModel) (response models.SnapshotPolicyResourceModel) {
	for _, sp := range sps {
		response = models.SnapshotPolicyResourceModel{
			ID:                               types.StringValue(string(sp.ID)),
			Name:                             types.StringValue(sp.Name),
			AutoSnapshotCreationCadenceInMin: types.Int64Value((int64)(sp.AutoSnapshotCreationCadenceInMin)),
			SnapshotAccessMode:               types.StringValue(string(sp.SnapshotAccessMode)),
			SecureSnapshots:                  types.BoolValue(sp.SecureSnapshots),
		}
		for _, rspl := range sp.NumOfRetainedSnapshotsPerLevel {
			response.NumOfRetainedSnapshotsPerLevel = append(response.NumOfRetainedSnapshotsPerLevel, (int64)(rspl))
		}
		if sp.SnapshotPolicyState == "Active" {
			response.Paused = types.BoolValue(false)
		} else {
			response.Paused = types.BoolValue(true)
		}
	}
	for _, v := range volumes {
		response.VolumeIds = append(response.VolumeIds, types.StringValue(v.ID))
	}
	response.RemoveMode = state.RemoveMode
	return response
}

// DifferenceArray function to find the state difference b/w volumes
func DifferenceArray(a, b []string) ([]string, []string) {
	var addedItems, removedItems []string
	//Find added items
	for _, item := range b {
		found := false
		for _, val := range a {
			if item == val {
				found = true
				break
			}
		}
		if !found {
			addedItems = append(addedItems, item)
		}
	}
	// find removed items
	for _, item := range a {
		found := false
		for _, val := range b {
			if item == val {
				found = true
				break
			}
		}
		if !found {
			removedItems = append(removedItems, item)
		}
	}
	return addedItems, removedItems
}

// ListToSlice converts the list to slice for num of retained snapshots per level
func ListToSlice(snap models.SnapshotPolicyResourceModel) []string {
	stringList := make([]string, len(snap.NumOfRetainedSnapshotsPerLevel))
	for i, v := range snap.NumOfRetainedSnapshotsPerLevel {
		stringList[i] = strconv.Itoa((int)(v))
	}
	return stringList
}

// ListToSliceVol converts the list to slice for volumes
func ListToSliceVol(snap models.SnapshotPolicyResourceModel) []string {
	stringList := make([]string, len(snap.VolumeIds))
	for i, v := range snap.VolumeIds {
		stringList[i] = v.ValueString()
	}
	return stringList
}

func GetAllSnapshotPolicies(client *goscaleio.Client) ([]scaleiotypes.SnapshotPolicy, error) {
	sps := []scaleiotypes.SnapshotPolicy{}
	resp, err := client.GetSnapshotPolicy("", "")
	if err != nil {
		return nil, err
	}
	for _, val := range resp {
		sps = append(sps, *val)
	}
	return sps, nil
}
