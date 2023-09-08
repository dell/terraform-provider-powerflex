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

# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check storagepool_resource_import.tf for more info
# To create / update, either protection_domain_id or protection_domain_name must be provided
# name and media_type is the required parameter to create or update
# other  atrributes like : use_rmcache, use_rfcache, replication_journal_capacity, capacity_alert_high_threshold, capacity_alert_critical_threshold etc. are optional 
# To check which attributes of the storage pool can be updated, please refer Product Guide in the documentation

resource "powerflex_storage_pool" "sp" {
  name                         = "newstoragepool"
  protection_domain_name       = "domain1"
  media_type                   = "HDD"
  use_rmcache                  = false
  use_rfcache                  = true
  replication_journal_capacity = 34
  zero_padding_enabled         = false
  rebalance_enabled            = false

  # Alert Thresholds
  # Critical threshold must be greater than high threshold
  capacity_alert_high_threshold     = 66
  capacity_alert_critical_threshold = 77

  # Protected Maintenance Mode Parameters
  # When the policy is set to "favorAppIos", then concurrent IOs and bandwidth limit can be set.
  # When the policy is set to "limitNumOfConcurrentIos", then only concurrent IOs can be set.
  # When the policy is set to "unlimited", then concurrent IOs and bandwidth limit can't be set.
  protected_maintenance_mode_io_priority_policy               = "favorAppIos"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps      = 1028

  # Rebalance Parameters
  # When the policy is set to "favorAppIos", then concurrent IOs and bandwidth limit can be set.
  # When the policy is set to "limitNumOfConcurrentIos", then only concurrent IOs can be set.
  # When the policy is set to "unlimited", then concurrent IOs and bandwidth limit can't be set.  
  rebalance_io_priority_policy               = "favorAppIos"
  rebalance_num_of_concurrent_ios_per_device = 7
  rebalance_bw_limit_per_device_in_kbps      = 1032

  #VTree Migration Parameters
  # When the policy is set to "favorAppIos", then concurrent IOs and bandwidth limit can be set.
  # When the policy is set to "limitNumOfConcurrentIos", then only concurrent IOs can be set.
  vtree_migration_io_priority_policy               = "favorAppIos"
  vtree_migration_num_of_concurrent_ios_per_device = 7
  vtree_migration_bw_limit_per_device_in_kbps      = 1030


  spare_percentage              = 66
  rm_cache_write_handling_mode  = "Passthrough"
  rebuild_enabled               = true
  rebuild_rebalance_parallelism = 5
  fragmentation                 = false
}

output "created_storagepool" {
  value = powerflex_storage_pool.sp
}
