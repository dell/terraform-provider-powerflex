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

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetNodeState returns the state for node data source
func GetNodeState(node scaleiotypes.NodeDetails) (response models.NodeModel) {
	response = models.NodeModel{
		RefID:               types.StringValue(node.RefID),
		IPAddress:           types.StringValue(node.IPAddress),
		CurrentIPAddress:    types.StringValue(node.CurrentIPAddress),
		ServiceTag:          types.StringValue(node.ServiceTag),
		Model:               types.StringValue(node.Model),
		DeviceType:          types.StringValue(node.DeviceType),
		DiscoverDeviceType:  types.StringValue(node.DiscoverDeviceType),
		DisplayName:         types.StringValue(node.DisplayName),
		ManagedState:        types.StringValue(node.ManagedState),
		State:               types.StringValue(node.State),
		InUse:               types.BoolValue(node.InUse),
		CustomFirmware:      types.BoolValue(node.CustomFirmware),
		NeedsAttention:      types.BoolValue(node.NeedsAttention),
		Manufacturer:        types.StringValue(node.Manufacturer),
		SystemID:            types.StringValue(node.SystemID),
		Health:              types.StringValue(node.Health),
		HealthMessage:       types.StringValue(node.HealthMessage),
		OperatingSystem:     types.StringValue(node.OperatingSystem),
		NumberOfCPUs:        types.Int64Value(int64(node.NumberOfCPUs)),
		Nics:                types.Int64Value(int64(node.Nics)),
		MemoryInGB:          types.Int64Value(int64(node.MemoryInGB)),
		ComplianceCheckDate: types.StringValue(node.ComplianceCheckDate),
		DiscoveredDate:      types.StringValue(node.DiscoveredDate),
		CredID:              types.StringValue(node.CredID),
		Compliance:          types.StringValue(node.Compliance),
		FailuresCount:       types.Int64Value(int64(node.FailuresCount)),
		Facts:               types.StringValue(node.Facts),
		PuppetCertName:      types.StringValue(node.PuppetCertName),
		FlexosMaintMode:     types.Int64Value(int64(node.FlexosMaintMode)),
		EsxiMaintMode:       types.Int64Value(int64(node.EsxiMaintMode)),
	}

	var deviceList []models.DeviceGroup
	for _, device := range node.DeviceGroupList.DeviceGroup {
		var groupUserList models.GroupUserList
		groupUserList.TotalRecords = device.GroupUserList.TotalRecords
		var groupUser []models.GroupUsers
		for _, user := range device.GroupUserList.GroupUsers {
			groupUser = append(groupUser, models.GroupUsers{
				UserSeqID: user.UserSeqID,
				UserName:  user.UserName,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				Enabled:   user.Enabled,
			})
		}
		groupUserList.GroupUsers = groupUser

		deviceList = append(deviceList, models.DeviceGroup{
			GroupSeqID:       device.GroupSeqID,
			GroupName:        device.GroupName,
			GroupDescription: device.GroupDescription,
			CreatedDate:      device.CreatedDate,
			CreatedBy:        device.CreatedBy,
			UpdatedDate:      device.UpdatedDate,
			UpdatedBy:        device.UpdatedBy,
			GroupUserList:    groupUserList,
		})
	}
	response.DeviceGroupList.DeviceGroup = deviceList
	return
}
