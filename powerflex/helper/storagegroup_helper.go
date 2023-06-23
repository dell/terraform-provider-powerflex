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
	"strconv"
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Function to update the State for Storagepool Resource
func UpdateStoragepoolState(storagepool *scaleiotypes.StoragePool, plan models.StoragepoolResourceModel) models.StoragepoolResourceModel {
	state := plan
	state.ProtectionDomainID = types.StringValue(storagepool.ProtectionDomainID)
	state.ID = types.StringValue(storagepool.ID)
	state.Name = types.StringValue(storagepool.Name)
	state.MediaType = types.StringValue(storagepool.MediaType)
	state.UseRmcache = types.BoolValue(storagepool.UseRmcache)
	state.UseRfcache = types.BoolValue(storagepool.UseRfcache)
	state.ZeroPaddingEnabled = types.BoolValue(storagepool.ZeroPaddingEnabled)
	state.ReplicationJournalCapacity = types.Int64Value(int64(storagepool.ReplicationCapacityMaxRatio))
	state.CapacityAlertHighThreshold = types.Int64Value(int64(storagepool.CapacityAlertHighThreshold))
	state.CapacityAlertCriticalThreshold = types.Int64Value(int64(storagepool.CapacityAlertCriticalThreshold))
	state.ProtectedMaintenanceModeIoPriorityPolicy = types.StringValue(storagepool.ProtectedMaintenanceModeIoPriorityPolicy)
	state.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice))
	state.ProtectedMaintenanceModeBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps))
	state.RebalanceEnabled = types.BoolValue(storagepool.RebalanceEnabled)
	state.RebalanceIoPriorityPolicy = types.StringValue(storagepool.RebalanceioPriorityPolicy)
	state.RebalanceNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.RebalanceioPriorityNumOfConcurrentIosPerDevice))
	state.RebalanceBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.RebalanceioPriorityBwLimitPerDeviceInKbps))
	state.VtreeMigrationIoPriorityPolicy = types.StringValue(storagepool.VtreeMigrationIoPriorityPolicy)
	state.VtreeMigrationNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice))
	state.VtreeMigrationBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps))
	state.SparePercentage = types.Int64Value(int64(storagepool.SparePercentage))
	state.RmCacheWriteHandlingMode = types.StringValue(storagepool.RmCacheWriteHandlingMode)
	state.RebuildEnabled = types.BoolValue(storagepool.RebuildEnabled)
	state.RebuildRebalanceParallelism = types.Int64Value(int64(storagepool.NumofParallelRebuildRebalanceJobsPerDevice))
	state.Fragmentation = types.BoolValue(storagepool.FragmentationEnabled)
	return state
}

func IsCritcalAlert(plan, state models.StoragepoolResourceModel) (*scaleiotypes.CapacityAlertThresholdParam, bool) {
	payload, ok := scaleiotypes.CapacityAlertThresholdParam{}, true
	if !plan.CapacityAlertHighThreshold.IsUnknown() && !state.CapacityAlertHighThreshold.Equal(plan.CapacityAlertHighThreshold) {
		ok = false
		payload.CapacityAlertHighThresholdPercent = strconv.FormatInt(plan.CapacityAlertHighThreshold.ValueInt64(), 10)
	}
	if !plan.CapacityAlertCriticalThreshold.IsUnknown() && !state.CapacityAlertCriticalThreshold.Equal(plan.CapacityAlertCriticalThreshold) {
		ok = false
		payload.CapacityAlertCriticalThresholdPercent = strconv.FormatInt(plan.CapacityAlertCriticalThreshold.ValueInt64(), 10)
	}
	return &payload, ok
}

func IsProtectedMaintenance(plan, state models.StoragepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && !state.ProtectedMaintenanceModeIoPriorityPolicy.Equal(plan.ProtectedMaintenanceModeIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.ProtectedMaintenanceModeIoPriorityPolicy.ValueString()
	}
	if !plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && !state.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.Equal(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown() && !state.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.Equal(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}

func IsRebalance(plan, state models.StoragepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.RebalanceIoPriorityPolicy.IsUnknown() && !state.RebalanceIoPriorityPolicy.Equal(plan.RebalanceIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.RebalanceIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.RebalanceIoPriorityPolicy.ValueString()
	}
	if !plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && !state.RebalanceNumOfConcurrentIosPerDevice.Equal(plan.RebalanceNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.RebalanceNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown() && !state.RebalanceBwLimitPerDeviceInKbps.Equal(plan.RebalanceBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.RebalanceBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}

func IsVtreeMigration(plan, state models.StoragepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && !state.VtreeMigrationIoPriorityPolicy.Equal(plan.VtreeMigrationIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.VtreeMigrationIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.VtreeMigrationIoPriorityPolicy.ValueString()
	}
	if !plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && !state.VtreeMigrationNumOfConcurrentIosPerDevice.Equal(plan.VtreeMigrationNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.VtreeMigrationNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown() && !state.VtreeMigrationBwLimitPerDeviceInKbps.Equal(plan.VtreeMigrationBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.VtreeMigrationBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}

func GetStoragePoolState(volList []*scaleiotypes.Volume, sdsList []scaleiotypes.Sds, s1 *scaleiotypes.StoragePool) (storagePool models.StoragePoolModel) {
	storagePool = models.StoragePoolModel{
		ID:   types.StringValue(s1.ID),
		Name: types.StringValue(s1.Name),
	}

	// Iterate through volume list
	for _, vol := range volList {
		storagePool.Volumes = append(storagePool.Volumes, models.Volume{
			ID:   types.StringValue(vol.ID),
			Name: types.StringValue(vol.Name),
		})
	}

	// Iterate through SDS list
	for _, sds := range sdsList {
		storagePool.SDS = append(storagePool.SDS, models.SdsData{
			ID:   types.StringValue(sds.ID),
			Name: types.StringValue(sds.Name),
		})
	}

	// Iterate through the Links
	for _, link := range s1.Links {
		storagePool.Links = append(storagePool.Links, models.LinkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}

	storagePool.RebalanceioPriorityPolicy = types.StringValue(s1.RebalanceioPriorityPolicy)
	storagePool.RebalanceioPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.RebalanceioPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.RebalanceioPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.RebalanceioPriorityAppIopsPerDeviceThreshold))
	storagePool.RebalanceioPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.RebalanceioPriorityBwLimitPerDeviceInKbps))
	storagePool.RebalanceioPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.RebalanceioPriorityQuietPeriodInMsec))
	storagePool.RebalanceioPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.RebalanceioPriorityNumOfConcurrentIosPerDevice))
	storagePool.RebuildioPriorityPolicy = types.StringValue(s1.RebuildioPriorityPolicy)
	storagePool.RebuildioPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.RebuildioPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.RebuildioPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.RebuildioPriorityAppIopsPerDeviceThreshold))
	storagePool.RebuildioPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.RebalanceioPriorityBwLimitPerDeviceInKbps))
	storagePool.RebuildioPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.RebuildioPriorityQuietPeriodInMsec))
	storagePool.RebuildioPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.RebuildioPriorityNumOfConcurrentIosPerDevice))
	storagePool.ZeroPaddingEnabled = types.BoolValue(s1.ZeroPaddingEnabled)
	storagePool.UseRmcache = types.BoolValue(s1.UseRmcache)
	storagePool.SparePercentage = types.Int64Value(int64(s1.SparePercentage))
	storagePool.RmCacheWriteHandlingMode = types.StringValue(s1.RmCacheWriteHandlingMode)
	storagePool.RebalanceEnabled = types.BoolValue(s1.RebalanceEnabled)
	storagePool.RebuildEnabled = types.BoolValue(s1.RebuildEnabled)
	storagePool.NumofParallelRebuildRebalanceJobsPerDevice = types.Int64Value(int64(s1.NumofParallelRebuildRebalanceJobsPerDevice))
	storagePool.BackgroundScannerBWLimitKBps = types.Int64Value(int64(s1.BackgroundScannerBWLimitKBps))
	storagePool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice))
	storagePool.DataLayout = types.StringValue(s1.DataLayout)
	storagePool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps))
	storagePool.VtreeMigrationIoPriorityPolicy = types.StringValue(s1.VtreeMigrationIoPriorityPolicy)
	storagePool.AddressSpaceUsage = types.StringValue(s1.AddressSpaceUsage)
	storagePool.ExternalAccelerationType = types.StringValue(s1.ExternalAccelerationType)
	storagePool.PersistentChecksumState = types.StringValue(s1.PersistentChecksumState)
	storagePool.UseRfcache = types.BoolValue(s1.UseRfcache)
	storagePool.ChecksumEnabled = types.BoolValue(s1.ChecksumEnabled)
	storagePool.CompressionMethod = types.StringValue(s1.CompressionMethod)
	storagePool.FragmentationEnabled = types.BoolValue(s1.FragmentationEnabled)
	storagePool.CapacityUsageState = types.StringValue(s1.CapacityUsageState)
	storagePool.CapacityUsageType = types.StringValue(s1.CapacityUsageType)
	storagePool.AddressSpaceUsageType = types.StringValue(s1.AddressSpaceUsageType)
	storagePool.BgScannerCompareErrorAction = types.StringValue(s1.BgScannerCompareErrorAction)
	storagePool.BgScannerReadErrorAction = types.StringValue(s1.BgScannerReadErrorAction)
	storagePool.ReplicationCapacityMaxRatio = types.Int64Value(int64(s1.ReplicationCapacityMaxRatio))
	storagePool.PersistentChecksumEnabled = types.BoolValue(s1.PersistentChecksumEnabled)
	storagePool.PersistentChecksumBuilderLimitKb = types.Int64Value(int64(s1.PersistentChecksumBuilderLimitKb))
	storagePool.PersistentChecksumValidateOnRead = types.BoolValue(s1.PersistentChecksumValidateOnRead)
	storagePool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice))
	storagePool.ProtectedMaintenanceModeIoPriorityPolicy = types.StringValue(s1.ProtectedMaintenanceModeIoPriorityPolicy)
	storagePool.BackgroundScannerMode = types.StringValue(s1.BackgroundScannerMode)
	storagePool.MediaType = types.StringValue(s1.MediaType)
	storagePool.CapacityAlertHighThreshold = types.Int64Value(int64(s1.CapacityAlertHighThreshold))
	storagePool.CapacityAlertHighThreshold = types.Int64Value(int64(s1.CapacityAlertCriticalThreshold))
	storagePool.VtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.VtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.VtreeMigrationIoPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.VtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.VtreeMigrationIoPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.VtreeMigrationIoPriorityQuietPeriodInMsec))
	storagePool.FglAccpID = types.StringValue(s1.FglAccpID)
	storagePool.FglExtraCapacity = types.Int64Value(int64(s1.FglExtraCapacity))
	storagePool.FglOverProvisioningFactor = types.Int64Value(int64(s1.FglOverProvisioningFactor))
	storagePool.FglWriteAtomicitySize = types.Int64Value(int64(s1.FglWriteAtomicitySize))
	storagePool.FglNvdimmWriteCacheSizeInMb = types.Int64Value(int64(s1.FglNvdimmWriteCacheSizeInMb))
	storagePool.FglNvdimmMetadataAmortizationX100 = types.Int64Value(int64(s1.FglNvdimmMetadataAmortizationX100))
	storagePool.FglPerfProfile = types.StringValue(s1.FglPerfProfile)
	storagePool.ProtectedMaintenanceModeIoPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.ProtectedMaintenanceModeIoPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityAppIopsPerDeviceThreshold))
	storagePool.ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps))
	storagePool.ProtectedMaintenanceModeIoPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityQuietPeriodInMsec))
	return
}
