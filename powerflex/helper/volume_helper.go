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
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// MiKB to convert size in megabytes
	MiKB = 1024
	// GiKB to convert size in gigabytes
	GiKB = 1024 * MiKB
	// TiKB to convert size in terabytes
	TiKB = 1024 * GiKB
)

// ConvertToKB fucntion to convert size into kb
func ConvertToKB(capacityUnit string, size int64) int64 {
	var valInKiB int64
	switch capacityUnit {
	case "MB":
		valInKiB = size * MiKB
	case "TB":
		valInKiB = size * TiKB
	case "GB":
		valInKiB = size * GiKB
	}
	return valInKiB
}

// RefreshVolumeState function to update the state of volume resource in terraform.tfstate file
func RefreshVolumeState(vol *pftypes.Volume, state *models.VolumeResourceModel) (diags diag.Diagnostics) {
	state.StoragePoolID = types.StringValue(vol.StoragePoolID)
	state.UseRmCache = types.BoolValue(vol.UseRmCache)
	state.VolumeType = types.StringValue(vol.VolumeType)
	state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	state.Name = types.StringValue(vol.Name)
	state.ID = types.StringValue(vol.ID)
	state.AccessMode = types.StringValue(vol.AccessModeLimit)
	state.CompressionMethod = types.StringValue(vol.CompressionMethod)
	return diags
}

// GetStoragePoolInstance function to get storage pool from storage pool id and protection domain id
func GetStoragePoolInstance(c *goscaleio.Client, spID string, pdID string) (*goscaleio.StoragePool, error) {
	getSystems, _ := c.GetSystems()
	sr := goscaleio.NewSystem(c)
	sr.System = getSystems[0]
	pdr := goscaleio.NewProtectionDomain(c)
	protectionDomain, err := sr.FindProtectionDomain(pdID, "", "")
	if err != nil {
		return nil, err
	}
	pdr.ProtectionDomain = protectionDomain
	spr := goscaleio.NewStoragePool(c)
	storagePool, err := pdr.FindStoragePool(spID, "", "")
	if err != nil {
		return nil, err
	}
	spr.StoragePool = storagePool
	return spr, nil
}

// Difference function to find the state difference b/w sdcs
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// DifferenceMap function to find the state difference b/w sdcs
func DifferenceMap(a, b map[string]string) map[string]string {
	m := make(map[string]bool)
	diff := make(map[string]string)

	for item := range b {
		m[item] = true
	}

	for item := range a {
		if _, ok := m[item]; !ok {
			diff[item] = item
		}
	}
	return diff
}

// UpdateVolumeState iterates over the volume list and update the state
func UpdateVolumeState(volumes []*scaleiotypes.Volume) (response []models.VolumeModel) {
	for _, volumeValue := range volumes {
		volumeState := models.VolumeModel{
			ID:                                 types.StringValue(volumeValue.ID),
			Name:                               types.StringValue(volumeValue.Name),
			CreationTime:                       types.Int64Value((int64)(volumeValue.CreationTime)),
			SizeInKb:                           types.Int64Value((int64)(volumeValue.SizeInKb)),
			AncestorVolumeID:                   types.StringValue(volumeValue.AncestorVolumeID),
			VTreeID:                            types.StringValue(volumeValue.VTreeID),
			ConsistencyGroupID:                 types.StringValue(volumeValue.ConsistencyGroupID),
			VolumeType:                         types.StringValue(volumeValue.VolumeType),
			UseRmCache:                         types.BoolValue(volumeValue.UseRmCache),
			StoragePoolID:                      types.StringValue(volumeValue.StoragePoolID),
			DataLayout:                         types.StringValue(volumeValue.DataLayout),
			NotGenuineSnapshot:                 types.BoolValue(volumeValue.NotGenuineSnapshot),
			AccessModeLimit:                    types.StringValue(volumeValue.AccessModeLimit),
			SecureSnapshotExpTime:              types.Int64Value((int64)(volumeValue.SecureSnapshotExpTime)),
			ManagedBy:                          types.StringValue(volumeValue.ManagedBy),
			LockedAutoSnapshot:                 types.BoolValue(volumeValue.LockedAutoSnapshot),
			LockedAutoSnapshotMarkedForRemoval: types.BoolValue(volumeValue.LockedAutoSnapshotMarkedForRemoval),
			CompressionMethod:                  types.StringValue(volumeValue.CompressionMethod),
			TimeStampIsAccurate:                types.BoolValue(volumeValue.TimeStampIsAccurate),
			OriginalExpiryTime:                 types.Int64Value((int64)(volumeValue.OriginalExpiryTime)),
			VolumeReplicationState:             types.StringValue(volumeValue.VolumeReplicationState),
			ReplicationJournalVolume:           types.BoolValue(volumeValue.ReplicationJournalVolume),
			ReplicationTimeStamp:               types.Int64Value((int64)(volumeValue.ReplicationTimeStamp)),
		}

		for _, link := range volumeValue.Links {
			volumeState.Links = append(volumeState.Links, models.VolumeLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		for _, sdc := range volumeValue.MappedSdcInfo {
			volumeState.MappedSdcInfo = append(volumeState.MappedSdcInfo, models.MappedSdcInfoModel{
				SdcID:                 types.StringValue(sdc.SdcID),
				SdcIP:                 types.StringValue(sdc.SdcIP),
				LimitIops:             types.Int64Value((int64)(sdc.LimitIops)),
				LimitBwInMbps:         types.Int64Value((int64)(sdc.LimitBwInMbps)),
				SdcName:               types.StringValue(sdc.SdcName),
				AccessMode:            types.StringValue(sdc.AccessMode),
				IsDirectBufferMapping: types.BoolValue(sdc.IsDirectBufferMapping),
			})
		}
		response = append(response, volumeState)
	}
	return
}

// SetMappedSdcLimits setsMapped SDC limits
func SetMappedSdcLimits(volType *goscaleio.Volume, limitType scaleiotypes.SetMappedSdcLimitsParam) error {
	return volType.SetMappedSdcLimits(&limitType)
}

// MapVolumeSdc Map Volume to SDC
func MapVolumeSdc(volType *goscaleio.Volume, mapType scaleiotypes.MapVolumeSdcParam) error {
	return volType.MapVolumeSdc(&mapType)
}
