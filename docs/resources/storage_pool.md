---
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerflex_storage_pool resource"
linkTitle: "powerflex_storage_pool"
page_title: "powerflex_storage_pool Resource - powerflex"
subcategory: ""
description: |-
  This resource is used to manage the Storage Pool entity of PowerFlex Array. We can Create, Update and Delete the storage pool using this resource. We can also import an existing storage pool from PowerFlex array.
---

# powerflex_storage_pool (Resource)

This resource is used to manage the Storage Pool entity of PowerFlex Array. We can Create, Update and Delete the storage pool using this resource. We can also import an existing storage pool from PowerFlex array.

> **Caution:** <span style='color: red;' >Storage Pool creation or update is not atomic. In case of partially completed create operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
If the taint is not removed, terraform will destroy and recreate the resource.</span>

> **Note:** Either `protection_domain_name` or `protection_domain_id` is required. But not both. 

## Example Usage

```terraform
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

# Example for creating storage pool. After successful execution, storage pool will be created under a specified protection domain.
resource "powerflex_storage_pool" "sp" {
  name                         = "newstoragepool"
  protection_domain_name       = "domain1"
  media_type                   = "HDD" # HDD/SSD/Transitional
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
  protected_maintenance_mode_io_priority_policy               = "favorAppIos" # favorAppIos/limitNumOfConcurrentIos/unlimited
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps      = 1028

  # Rebalance Parameters
  # When the policy is set to "favorAppIos", then concurrent IOs and bandwidth limit can be set.
  # When the policy is set to "limitNumOfConcurrentIos", then only concurrent IOs can be set.
  # When the policy is set to "unlimited", then concurrent IOs and bandwidth limit can't be set.  
  rebalance_io_priority_policy               = "favorAppIos" # favorAppIos/limitNumOfConcurrentIos/unlimited
  rebalance_num_of_concurrent_ios_per_device = 7
  rebalance_bw_limit_per_device_in_kbps      = 1032

  #VTree Migration Parameters
  # When the policy is set to "favorAppIos", then concurrent IOs and bandwidth limit can be set.
  # When the policy is set to "limitNumOfConcurrentIos", then only concurrent IOs can be set.
  vtree_migration_io_priority_policy               = "favorAppIos" # favorAppIos/limitNumOfConcurrentIos
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
```

After the execution of above resource block, storage pool would have been created on the PowerFlex array. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `media_type` (String) Media Type of the storage pool. Valid values are `HDD`, `SSD` and `Transitional`
- `name` (String) Name of the Storage pool

### Optional

- `capacity_alert_critical_threshold` (Number) Set the threshold for triggering capacity usage critical-priority alert.
- `capacity_alert_high_threshold` (Number) Set the threshold for triggering capacity usage high-priority alert.
- `fragmentation` (Boolean) Enable or disable fragmentation in the Storage Pool
- `protected_maintenance_mode_bw_limit_per_device_in_kbps` (Number) The maximum bandwidth of protected maintenance mode migration I/Os, in KB per second, per device
- `protected_maintenance_mode_io_priority_policy` (String) Set the I/O priority policy for protected maintenance mode for a specific Storage Pool. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`
- `protected_maintenance_mode_num_of_concurrent_ios_per_device` (Number) The maximum number of concurrent protected maintenance mode migration I/Os per device
- `protection_domain_id` (String) ID of the Protection Domain under which the storage pool will be created. Conflicts with `protection_domain_name`. Cannot be updated.
- `protection_domain_name` (String) Name of the Protection Domain under which the storage pool will be created. Conflicts with `protection_domain_id`. Cannot be updated.
- `rebalance_bw_limit_per_device_in_kbps` (Number) The maximum bandwidth of rebalance I/Os, in KB/s, per device
- `rebalance_enabled` (Boolean) Enable or disable rebalancing in the specified Storage Pool
- `rebalance_io_priority_policy` (String) Policy to use for rebalance I/O priority. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`
- `rebalance_num_of_concurrent_ios_per_device` (Number) The maximum number of concurrent rebalance I/Os per device
- `rebuild_enabled` (Boolean) Enable or disable rebuilds in the specified Storage Pool
- `rebuild_rebalance_parallelism` (Number) Maximum number of concurrent rebuild and rebalance activities on SDSs in the Storage Pool
- `replication_journal_capacity` (Number) This defines the maximum percentage of Storage Pool capacity that can be used by replication for the journal. Before deleting the storage pool, this has to be set to 0.
- `rm_cache_write_handling_mode` (String) Sets the Read RAM Cache write handling mode of the specified Storage Pool
- `spare_percentage` (Number) Sets the spare capacity reservation policy
- `use_rfcache` (Boolean) Enable/Disable RFcache on a specific storage pool
- `use_rmcache` (Boolean) Enable/Disable RMcache on a specific storage pool
- `vtree_migration_bw_limit_per_device_in_kbps` (Number) The maximum bandwidth of V-Tree migration IOs, in KB per second, per device
- `vtree_migration_io_priority_policy` (String) Set the I/O priority policy for V-Tree migration for a specific Storage Pool. Valid values are `limitNumOfConcurrentIos` and `favorAppIos`
- `vtree_migration_num_of_concurrent_ios_per_device` (Number) The maximum number of concurrent V-Tree migration I/Os per device
- `zero_padding_enabled` (Boolean) Enable/Disable padding policy on a specific storage pool

### Read-Only

- `id` (String) ID of the Storage pool

## Import

Import is supported using the following syntax:

```shell
# /*
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import storage pool by it's id
terraform import powerflex_storage_pool.storage_pool_import_by_id "<id>"
```

1. This will import the storage pool instance with specified ID into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.