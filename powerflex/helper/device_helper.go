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

	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateDeviceState function to update the state of device resource in state file
func UpdateDeviceState(deviceResponse *goscaleio_types.Device, plan models.DeviceModel) (models.DeviceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	state.ID = types.StringValue(deviceResponse.ID)
	if deviceResponse.Name == "" {
		state.Name = types.StringNull()
	} else {
		state.Name = types.StringValue(deviceResponse.Name)
	}

	state.DevicePath = types.StringValue(deviceResponse.DeviceCurrentPathName)
	state.DeviceOriginalPath = types.StringValue(deviceResponse.DeviceOriginalPathName)
	state.MediaType = types.StringValue(deviceResponse.MediaType)
	state.ExternalAccelerationType = types.StringValue(deviceResponse.ExternalAccelerationType)
	state.DeviceCapacityInKB = types.Int64Value(int64(deviceResponse.CapacityLimitInKb))
	state.DeviceState = types.StringValue(deviceResponse.DeviceState)
	state.SdsID = types.StringValue(deviceResponse.SdsID)
	state.StoragePoolID = types.StringValue(deviceResponse.StoragePoolID)
	return state, diags
}

// GetAllDeviceState saves the state of device datasource
func GetAllDeviceState(devices []goscaleio_types.Device) (response []models.DeviceModelData) {
	for _, device := range devices {
		deviceState := models.DeviceModelData{
			FglNvdimmMetadataAmortizationX100: types.Int64Value(int64(device.FglNvdimmMetadataAmortizationX100)),
			LogicalSectorSizeInBytes:          types.Int64Value(int64(device.LogicalSectorSizeInBytes)),
			FglNvdimmWriteCacheSize:           types.Int64Value(int64(device.FglNvdimmWriteCacheSize)),
			AccelerationPoolID:                types.StringValue(device.AccelerationPoolID),
			SdsID:                             types.StringValue(device.SdsID),
			StoragePoolID:                     types.StringValue(device.StoragePoolID),
			CapacityLimitInKb:                 types.Int64Value(int64(device.CapacityLimitInKb)),
			ErrorState:                        types.StringValue(device.ErrorState),
			Capacity:                          types.Int64Value(int64(device.Capacity)),
			DeviceType:                        types.StringValue(device.DeviceType),
			PersistentChecksumState:           types.StringValue(device.PersistentChecksumState),
			DeviceState:                       types.StringValue(device.DeviceState),
			LedSetting:                        types.StringValue(device.LedSetting),
			MaxCapacityInKb:                   types.Int64Value(int64(device.MaxCapacityInKb)),
			SpSdsID:                           types.StringValue(device.SpSdsID),
			AggregatedState:                   types.StringValue(device.AggregatedState),
			TemperatureState:                  types.StringValue(device.TemperatureState),
			SsdEndOfLifeState:                 types.StringValue(device.SsdEndOfLifeState),
			ModelName:                         types.StringValue(device.ModelName),
			VendorName:                        types.StringValue(device.VendorName),
			RaidControllerSerialNumber:        types.StringValue(device.RaidControllerSerialNumber),
			FirmwareVersion:                   types.StringValue(device.FirmwareVersion),
			CacheLookAheadActive:              types.BoolValue(device.CacheLookAheadActive),
			WriteCacheActive:                  types.BoolValue(device.WriteCacheActive),
			AtaSecurityActive:                 types.BoolValue(device.AtaSecurityActive),
			PhysicalSectorSizeInBytes:         types.Int64Value(int64(device.PhysicalSectorSizeInBytes)),
			MediaFailing:                      types.BoolValue(device.MediaFailing),
			SlotNumber:                        types.StringValue(device.SlotNumber),
			ExternalAccelerationType:          types.StringValue(device.ExternalAccelerationType),
			AutoDetectMediaType:               types.StringValue(device.AutoDetectMediaType),
			DeviceCurrentPathName:             types.StringValue(device.DeviceCurrentPathName),
			DeviceOriginalPathName:            types.StringValue(device.DeviceOriginalPathName),
			RfcacheErrorDeviceDoesNotExist:    types.BoolValue(device.RfcacheErrorDeviceDoesNotExist),
			MediaType:                         types.StringValue(device.MediaType),
			SerialNumber:                      types.StringValue(device.SerialNumber),
			Name:                              types.StringValue(device.Name),
			ID:                                types.StringValue(device.ID),
		}
		deviceState.RfcacheProps = models.RfcachePropsModel{
			DeviceUUID:                     types.StringValue(device.RfcacheProps.DeviceUUID),
			RfcacheErrorStuckIO:            types.BoolValue(device.RfcacheProps.RfcacheErrorStuckIO),
			RfcacheErrorHeavyLoadCacheSkip: types.BoolValue(device.RfcacheProps.RfcacheErrorHeavyLoadCacheSkip),
			RfcacheErrorCardIoError:        types.BoolValue(device.RfcacheProps.RfcacheErrorCardIoError),
		}
		deviceState.LongSuccessfulIos = models.LongSuccessfulIosModel{
			ShortWindow: models.DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.MaxFailuresCount)),
			},
			MediumWindow: models.DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.MaxFailuresCount)),
			},
			LongWindow: models.DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.MaxFailuresCount)),
			},
		}
		deviceState.StorageProps = models.StoragePropsModel{
			FglAccDeviceID:                   types.StringValue(device.StorageProps.FglAccDeviceID),
			FglNvdimmSizeMb:                  types.Int64Value(int64(device.StorageProps.FglNvdimmSizeMb)),
			DestFglNvdimmSizeMb:              types.Int64Value(int64(device.StorageProps.DestFglNvdimmSizeMb)),
			DestFglAccDeviceID:               types.StringValue(device.StorageProps.DestFglAccDeviceID),
			ChecksumMode:                     types.StringValue(device.StorageProps.ChecksumMode),
			DestChecksumMode:                 types.StringValue(device.StorageProps.DestChecksumMode),
			ChecksumAccDeviceID:              types.StringValue(device.StorageProps.ChecksumAccDeviceID),
			DestChecksumAccDeviceID:          types.StringValue(device.StorageProps.DestChecksumAccDeviceID),
			ChecksumSizeMb:                   types.Int64Value(int64(device.StorageProps.ChecksumSizeMb)),
			IsChecksumFullyCalculated:        types.BoolValue(device.StorageProps.IsChecksumFullyCalculated),
			ChecksumChangelogAccDeviceID:     types.StringValue(device.StorageProps.ChecksumChangelogAccDeviceID),
			DestChecksumChangelogAccDeviceID: types.StringValue(device.StorageProps.DestChecksumChangelogAccDeviceID),
			ChecksumChangelogSizeMb:          types.Int64Value(int64(device.StorageProps.ChecksumChangelogSizeMb)),
			DestChecksumChangelogSizeMb:      types.Int64Value(int64(device.StorageProps.DestChecksumChangelogSizeMb)),
		}
		deviceState.AccelerationProps = models.AccelerationPropsModel{
			AccUsedCapacityInKb: types.StringValue(device.AccelerationProps.AccUsedCapacityInKb),
		}
		for _, link := range device.Links {
			deviceState.Links = append(deviceState.Links, models.DeviceLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		response = append(response, deviceState)
	}
	return
}
