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
