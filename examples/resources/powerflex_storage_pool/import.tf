/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.
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

//Gather all existing storage pools
data "powerflex_storage_pool" "all" {
}

//Gather all existing protection domains
data "powerflex_protection_domain" "all"{
}

//Import all storage pools
import{
    for_each = data.powerflex_storage_pool.all.storage_pools
    to = powerflex_storage_pool.import_test_storage_pool[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_storage_pool" "import_test_storage_pool" {
    count = length(data.powerflex_storage_pool.all.storage_pools)
    name                         = data.powerflex_storage_pool.all.storage_pools[count.index].name
    protection_domain_name       = data.powerflex_protection_domain.all.protection_domains[count.index].name
    media_type                   = data.powerflex_storage_pool.all.storage_pools[count.index].media_type
    use_rmcache                  = data.powerflex_storage_pool.all.storage_pools[count.index].use_rm_cache
    use_rfcache                  = data.powerflex_storage_pool.all.storage_pools[count.index].use_rf_cache
    zero_padding_enabled         = data.powerflex_storage_pool.all.storage_pools[count.index].zero_padding_enabled
    rebalance_enabled            = data.powerflex_storage_pool.all.storage_pools[count.index].rebalance_enabled
    spare_percentage  = data.powerflex_storage_pool.all.storage_pools[count.index].spare_percentage
    rebuild_enabled  = data.powerflex_storage_pool.all.storage_pools[count.index].rebuild_enabled
    fragmentation  = data.powerflex_storage_pool.all.storage_pools[count.index].fragmentation_enabled
}