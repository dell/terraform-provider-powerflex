---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerflex_storage_pool data source"
linkTitle: "powerflex_storage_pool"
page_title: "powerflex_storage_pool Data Source - powerflex"
subcategory: ""
description: |-
  This datasource is used to query the existing storage pools from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.
---

# powerflex_storage_pool (Data Source)

## Example Usage

```terraform
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


output "storagePoolallresult" {
  value = data.powerflex_storage_pool.filtered.storage_pools
}
```

After the successful execution of above said block, We can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_storage_pool.example.attribute_name` where attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) Placeholder identifier attribute.
- `storage_pools` (Attributes List) List of fetched storage pools. (see [below for nested schema](#nestedatt--storage_pools))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `address_space_usage` (Set of String) List of address_space_usage
- `address_space_usage_type` (Set of String) List of address_space_usage_type
- `background_scanner_bw_limit_kbps` (Set of Number) List of background_scanner_bw_limit_kbps
- `background_scanner_mode` (Set of String) List of background_scanner_mode
- `bg_scanner_compare_error_action` (Set of String) List of bg_scanner_compare_error_action
- `bg_scanner_read_error_action` (Set of String) List of bg_scanner_read_error_action
- `capacity_alert_critical_threshold` (Set of Number) List of capacity_alert_critical_threshold
- `capacity_alert_high_threshold` (Set of Number) List of capacity_alert_high_threshold
- `capacity_usage_state` (Set of String) List of capacity_usage_state
- `capacity_usage_type` (Set of String) List of capacity_usage_type
- `checksum_enabled` (Boolean) Value for checksum_enabled
- `compression_method` (Set of String) List of compression_method
- `data_layout` (Set of String) List of data_layout
- `external_acceleration_type` (Set of String) List of external_acceleration_type
- `fgl_accp_id` (Set of String) List of fgl_accp_id
- `fgl_extra_capacity` (Set of Number) List of fgl_extra_capacity
- `fgl_nvdimm_metadata_amotization_x100` (Set of Number) List of fgl_nvdimm_metadata_amotization_x100
- `fgl_nvdimm_write_cache_size_mb` (Set of Number) List of fgl_nvdimm_write_cache_size_mb
- `fgl_overprovisioning_factor` (Set of Number) List of fgl_overprovisioning_factor
- `fgl_perf_profile` (Set of String) List of fgl_perf_profile
- `fgl_write_atomicity_size` (Set of Number) List of fgl_write_atomicity_size
- `fragmentation_enabled` (Boolean) Value for fragmentation_enabled
- `id` (Set of String) List of id
- `media_type` (Set of String) List of media_type
- `name` (Set of String) List of name
- `num_of_parallel_rebuild_rebalance_jobs_per_device` (Set of Number) List of num_of_parallel_rebuild_rebalance_jobs_per_device
- `persistent_checksum_builder_limit_kb` (Set of Number) List of persistent_checksum_builder_limit_kb
- `persistent_checksum_enabled` (Boolean) Value for persistent_checksum_enabled
- `persistent_checksum_state` (Set of String) List of persistent_checksum_state
- `persistent_checksum_validate_on_read` (Boolean) Value for persistent_checksum_validate_on_read
- `protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps` (Set of Number) List of protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps
- `protected_maintenance_mode_io_priority_app_iops_per_device_threshold` (Set of Number) List of protected_maintenance_mode_io_priority_app_iops_per_device_threshold
- `protected_maintenance_mode_io_priority_bw_limit_per_device_kbps` (Set of Number) List of protected_maintenance_mode_io_priority_bw_limit_per_device_kbps
- `protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device` (Set of Number) List of protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device
- `protected_maintenance_mode_io_priority_policy` (Set of String) List of protected_maintenance_mode_io_priority_policy
- `protected_maintenance_mode_io_priority_quiet_period_msec` (Set of Number) List of protected_maintenance_mode_io_priority_quiet_period_msec
- `rebalance_enabled` (Boolean) Value for rebalance_enabled
- `rebalance_io_priority_app_bw_per_device_threshold_kbps` (Set of Number) List of rebalance_io_priority_app_bw_per_device_threshold_kbps
- `rebalance_io_priority_app_iops_per_device_threshold` (Set of Number) List of rebalance_io_priority_app_iops_per_device_threshold
- `rebalance_io_priority_bw_limit_per_device_kbps` (Set of Number) List of rebalance_io_priority_bw_limit_per_device_kbps
- `rebalance_io_priority_num_of_concurrent_ios_per_device` (Set of Number) List of rebalance_io_priority_num_of_concurrent_ios_per_device
- `rebalance_io_priority_policy` (Set of String) List of rebalance_io_priority_policy
- `rebalance_io_priority_quiet_period_msec` (Set of Number) List of rebalance_io_priority_quiet_period_msec
- `rebuild_enabled` (Boolean) Value for rebuild_enabled
- `rebuild_io_priority_app_bw_per_device_threshold_kbps` (Set of Number) List of rebuild_io_priority_app_bw_per_device_threshold_kbps
- `rebuild_io_priority_app_iops_per_device_threshold` (Set of Number) List of rebuild_io_priority_app_iops_per_device_threshold
- `rebuild_io_priority_bw_limit_per_device_in_kbps` (Set of Number) List of rebuild_io_priority_bw_limit_per_device_in_kbps
- `rebuild_io_priority_num_of_concurrent_ios_per_device` (Set of Number) List of rebuild_io_priority_num_of_concurrent_ios_per_device
- `rebuild_io_priority_policy` (Set of String) List of rebuild_io_priority_policy
- `rebuild_io_priority_quiet_period_msec` (Set of Number) List of rebuild_io_priority_quiet_period_msec
- `replication_capacity_max_ratio` (Set of Number) List of replication_capacity_max_ratio
- `rm_cache_write_handling_mode` (Set of String) List of rm_cache_write_handling_mode
- `spare_percentage` (Set of Number) List of spare_percentage
- `use_rf_cache` (Boolean) Value for use_rf_cache
- `use_rm_cache` (Boolean) Value for use_rm_cache
- `vtree_migration_io_priority_app_bw_per_device_threshold_kbps` (Set of Number) List of vtree_migration_io_priority_app_bw_per_device_threshold_kbps
- `vtree_migration_io_priority_app_iops_per_device_threshold` (Set of Number) List of vtree_migration_io_priority_app_iops_per_device_threshold
- `vtree_migration_io_priority_bw_limit_per_device_kbps` (Set of Number) List of vtree_migration_io_priority_bw_limit_per_device_kbps
- `vtree_migration_io_priority_num_of_concurrent_ios_per_device` (Set of Number) List of vtree_migration_io_priority_num_of_concurrent_ios_per_device
- `vtree_migration_io_priority_policy` (Set of String) List of vtree_migration_io_priority_policy
- `vtree_migration_io_priority_quiet_period_msec` (Set of Number) List of vtree_migration_io_priority_quiet_period_msec
- `zero_padding_enabled` (Boolean) Value for zero_padding_enabled


<a id="nestedatt--storage_pools"></a>
### Nested Schema for `storage_pools`

Read-Only:

- `address_space_usage` (String) Address space usage.
- `address_space_usage_type` (String) Address space usage reason.
- `background_scanner_bw_limit_kbps` (Number) Background Scanner Bandwidth Limit.
- `background_scanner_mode` (String) Scanner mode.
- `bg_scanner_compare_error_action` (String) Scanner compare-error action.
- `bg_scanner_read_error_action` (String) Scanner read-error action.
- `capacity_alert_critical_threshold` (Number) Capacity alert critical threshold.
- `capacity_alert_high_threshold` (Number) Capacity alert high threshold.
- `capacity_usage_state` (String) Capacity usage state (normal/high/critical/full).
- `capacity_usage_type` (String) Usage state reason.
- `checksum_enabled` (Boolean) Checksum Enabled.
- `compression_method` (String) Compression method.
- `data_layout` (String) Data Layout.
- `external_acceleration_type` (String) External acceleration type.
- `fgl_accp_id` (String) FGL ID.
- `fgl_extra_capacity` (Number) FGL extra capacity.
- `fgl_nvdimm_metadata_amotization_x100` (Number) FGL NVDIMM metadata amortization.
- `fgl_nvdimm_write_cache_size_mb` (Number) FGL NVDIMM write cache size in Mb.
- `fgl_overprovisioning_factor` (Number) FGL overprovisioning factor.
- `fgl_perf_profile` (String) FGL performance profile.
- `fgl_write_atomicity_size` (Number) FGL write atomicity size.
- `fragmentation_enabled` (Boolean) Fragmentation Enabled.
- `id` (String) Storage pool ID.
- `links` (Attributes List) Specifies the links associated with storage pool. (see [below for nested schema](#nestedatt--storage_pools--links))
- `media_type` (String) Media type.
- `name` (String) Storage pool name.
- `num_of_parallel_rebuild_rebalance_jobs_per_device` (Number) Number of Parallel Rebuild/Rebalance Jobs per Device.
- `persistent_checksum_builder_limit_kb` (Number) Persistent checksum builder limit.
- `persistent_checksum_enabled` (Boolean) Persistent checksum enabled.
- `persistent_checksum_state` (String) Persistent Checksum State.
- `persistent_checksum_validate_on_read` (Boolean) Persistent checksum validation on read.
- `protected_maintenance_mode_io_priority_app_bw_per_device_threshold_kbps` (Number) Protected maintenance mode IO priority app bandwidth per device threshold in Kbps.
- `protected_maintenance_mode_io_priority_app_iops_per_device_threshold` (Number) Protected maintenance mode IO priority app IOPS per device threshold.
- `protected_maintenance_mode_io_priority_bw_limit_per_device_kbps` (Number) Protected maintenance mode IO priority bandwidth limit per device in Kbps.
- `protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device` (Number) Number of Concurrent Protected Maintenance Mode IOPS per Device.
- `protected_maintenance_mode_io_priority_policy` (String) Protected maintenance mode IO priority policy.
- `protected_maintenance_mode_io_priority_quiet_period_msec` (Number) Protected maintenance mode IO priority quiet period in Msec.
- `rebalance_enabled` (Boolean) Rebalance Enabled.
- `rebalance_io_priority_app_bw_per_device_threshold_kbps` (Number) Rebalance Application Bandwidth per Device Threshold.
- `rebalance_io_priority_app_iops_per_device_threshold` (Number) Rebalance Application IOPS per Device Threshold.
- `rebalance_io_priority_bw_limit_per_device_kbps` (Number) Rebalance Bandwidth Limit per Device.
- `rebalance_io_priority_num_of_concurrent_ios_per_device` (Number) Number of Concurrent Rebalance IOPS per Device.
- `rebalance_io_priority_policy` (String) Rebalance IO Priority Policy.
- `rebalance_io_priority_quiet_period_msec` (Number) Rebalance Quiet Period.
- `rebuild_enabled` (Boolean) Rebuild Enabled.
- `rebuild_io_priority_app_bw_per_device_threshold_kbps` (Number) Rebuild Application Bandwidth per Device Threshold.
- `rebuild_io_priority_app_iops_per_device_threshold` (Number) Rebuild Application IOPS per Device Threshold.
- `rebuild_io_priority_bw_limit_per_device_in_kbps` (Number) Rebuild Bandwidth Limit per Device.
- `rebuild_io_priority_num_of_concurrent_ios_per_device` (Number) Number of Concurrent Rebuild IOPS per Device.
- `rebuild_io_priority_policy` (String) Rebuild IO Priority Policy.
- `rebuild_io_priority_quiet_period_msec` (Number) Rebuild Quiet Period.
- `replication_capacity_max_ratio` (Number) Replication allowed capacity.
- `rm_cache_write_handling_mode` (String) RAM Read Cache Write Handling Mode.
- `sds` (Attributes List) List of SDS associated with storage pool. (see [below for nested schema](#nestedatt--storage_pools--sds))
- `spare_percentage` (Number) Spare Percentage.
- `use_rf_cache` (Boolean) Use Read Flash Cache.
- `use_rm_cache` (Boolean) Use RAM Read Cache.
- `volumes` (Attributes List) List of volumes associated with storage pool. (see [below for nested schema](#nestedatt--storage_pools--volumes))
- `vtree_migration_io_priority_app_bw_per_device_threshold_kbps` (Number) VTree migration IO priority App bandwidth per device threshold in Kbps.
- `vtree_migration_io_priority_app_iops_per_device_threshold` (Number) VTree migration IO priority App IOPS per device threshold.
- `vtree_migration_io_priority_bw_limit_per_device_kbps` (Number) VTree Migration Bandwidth Limit per Device.
- `vtree_migration_io_priority_num_of_concurrent_ios_per_device` (Number) Number of concurrent VTree migration IOPS per device.
- `vtree_migration_io_priority_policy` (String) VTree Migration IO Priority Policy.
- `vtree_migration_io_priority_quiet_period_msec` (Number) VTree migration IO priority quiet period in Msec.
- `zero_padding_enabled` (Boolean) Zero Padding Enabled.

<a id="nestedatt--storage_pools--links"></a>
### Nested Schema for `storage_pools.links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details.
- `rel` (String) Specifies the relationship with the storage pool.


<a id="nestedatt--storage_pools--sds"></a>
### Nested Schema for `storage_pools.sds`

Read-Only:

- `id` (String) SDS ID.
- `name` (String) SDS name.


<a id="nestedatt--storage_pools--volumes"></a>
### Nested Schema for `storage_pools.volumes`

Read-Only:

- `id` (String) Volume ID.
- `name` (String) Volume name.


