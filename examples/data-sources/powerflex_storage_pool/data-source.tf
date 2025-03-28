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

# commands to run this tf file : terraform init && terraform apply --auto-approve

data "powerflex_storage_pool" "all" {

}


output "storagePoolallresult" {
  value = data.powerflex_storage_pool.all.storage_pools
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_storage_pool" "storage_pool_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     rebuild_io_priority_policy = ["^limit.*$"]
#   }
# }

# output "storagePoolFilterRegexResult"{
#  value = data.powerflex_storage_pool.storage_pool_regex.storage_pools
# }


# Get Peer System details using filter with all values
# If there is no intersection between the filters then an empty datasource will be returned
# For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples
data "powerflex_storage_pool" "filtered" {
  filter {
    # id= ["id1","id2"]
    # name= ["name1","name2"]
    # rebalance_io_priority_policy= ["rebalanceIoPriorityPolicy1","rebalanceIoPriorityPolicy2"]
    # rebuild_io_priority_policy= ["rebuildIoPriorityPolicy1","rebuildIoPriorityPolicy2"]
    # rebuild_io_priority_bw_limit_per_device_in_kbps= [1,2]
    # rebuild_io_priority_num_of_concurrent_ios_per_device= [1,2]
    # rebalance_io_priority_num_of_concurrent_ios_per_device= [1,2]
    # rebalance_io_priority_bw_limit_per_device_kbps= [1,2]
    # rebuild_io_priority_app_iops_per_device_threshold= [1,2]
    # rebalance_io_priority_app_iops_per_device_threshold= [1,2]
    # rebuild_io_priority_app_bw_per_device_threshold_kbps= [1,2]
    # rebalance_io_priority_app_bw_per_device_threshold_kbps= [1,2]
    # rebuild_io_priority_quiet_period_msec= [1,2]
    # rebalance_io_priority_quiet_period_msec= [1,2]
    # zero_padding_enabled= true
    # use_rm_cache= true
    # spare_percentage= [1,2]
    # rm_cache_write_handling_mode= ["rmCacheWriteHandlingMode1","rmCacheWriteHandlingMode2"]
    # rebuild_enabled= true
    # rebalance_enabled= true
    # num_of_parallel_rebuild_rebalance_jobs_per_device= [1,2]
    # background_scanner_bw_limit_kbps= [1,2]
    # protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device= [1,2]
    # data_layout= ["dataLayout1","dataLayout2"]
    # vtree_migration_io_priority_bw_limit_per_device_kbps= [1,2]
    # vtree_migration_io_priority_policy= ["vtreeMigrationIoPriorityPolicy1","vtreeMigrationIoPriorityPolicy2"]
    # address_space_usage= ["addressSpaceUsage1","addressSpaceUsage2"]
    # external_acceleration_type= ["externalAccelerationType1","externalAccelerationType2"]
    # persistent_checksum_state= ["checksumState1","checksumState2"]
    # use_rf_cache= false
    # checksum_enabled= false
    # compression_method= ["compressionMethod1","compressionMethod2"]
    # fragmentation_enabled= true
    # capacity_usage_state= ["capacityUsageState1","capacityUsageState2"]
    # capacity_usage_type= ["capacityUsageType1","capacityUsageType2"]
    # address_space_usage_type= ["addressSpaceUsageType1","addressSpaceUsageType2"]
    # bg_scanner_compare_error_action= ["bgScannerCompareErrorAction1","bgScannerCompareErrorAction2"]
    # bg_scanner_read_error_action= ["bgScannerReadErrorAction1","bgScannerReadErrorAction2"]
    # replication_capacity_max_ratio= [1,2]
    # persistent_checksum_enabled= false
    # persistent_checksum_builder_limit_kb= [1,2]
    # persistent_checksum_validate_on_read= false
    # vtree_migration_io_priority_num_of_concurrent_ios_per_device= [1,2]
    # protected_maintenance_mode_io_priority_policy= ["protectedMaintenanceModeIoPriorityPolicy1","protectedMaintenanceModeIoPriorityPolicy2"]
    # background_scanner_mode= ["backgroundScannerMode1","backgroundScannerMode2"]
    # media_type= ["HDD","SSD"]
    # capacity_alert_high_threshold= [1,2]
    # capacity_alert_critical_threshold= [1,2]
    # vtree_migration_io_priority_app_iops_per_device_threshold= [1,2]
    # vtree_migration_io_priority_app_bw_per_device_threshold_kbps= [1,2]
    # vtree_migration_io_priority_quiet_period_msec= [1,2]
    # fgl_accp_id= ["fglAccpId1","fglAccpId2"]
    # fgl_extra_capacity= [1,2]
    # fgl_overprovisioning_factor= [1,2]
    # fgl_write_atomicity_size= [1,2]
    # fgl_nvdimm_write_cache_size_mb= [1,2]
    # fgl_nvdimm_metadata_amotization_x100= [1,2]
    # fgl_perf_profile= ["fglPerfProfile1","fglPerfProfile2"]
    # protected_maintenance_mode_io_priority_bw_limit_per_device_kbps= [1,2]
    # protected_maintenance_mode_io_priority_app_iops_per_device_threshold= [1,2]
    # protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps= [1,2]
    # protected_maintenance_mode_io_priority_quiet_period_msec= [1,2]
  }
}


# output "storagePoolallresult" {
#   value = data.powerflex_storage_pool.filtered.storage_pools
# }
