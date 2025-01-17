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
data "powerflex_protection_domain" "all" {
}

output "inputAll" {
  value = data.powerflex_protection_domain.all.protection_domains
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_protection_domain" "protection_domain_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     rf_cache_opertional_mode = ["^.*Write.*$"]
#   }
# }

# output "protectionDomainFilterRegexResult"{
#  value = data.powerflex_protection_domain.protection_domain_filter_regex.protection_domains
# }


# Get Peer System details using filter with all values
# If there is no intersection between the filters then an empty datasource will be returned
# For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples
data "powerflex_protection_domain" "filter" {
  filter {
    # system_id = ["systemID1", "systemID2"]
    # replication_capacity_max_ratio = [1,2]
    # rebuild_network_throttling_in_kbps = [1,2]
    # rebalance_network_throttling_in_kbps = [1,2]
    # overall_io_network_throttling_in_kbps = [1,2]
    # vtree_migration_network_throttling_in_kbps = [1,2]
    # protected_maintenance_mode_network_throttling_in_kbps = [1,2]
    # overall_io_network_throttling_enabled = false
    # rebuild_network_throttling_enabled = false
    # rebalance_network_throttling_enabled = false
    # vtree_migration_network_throttling_enabled = false
    # protected_maintenance_mode_network_throttling_enabled = false
    # fgl_default_num_concurrent_writes = [1,2]
    # fgl_metadata_cache_enabled = false
    # fgl_default_metadata_cache_size = [1,2]
    # rf_cache_enabled = false
    # rf_cache_accp_id = ["rfcache_accp_id1", "rfcache_accp_id2"]
    # rf_cache_opertional_mode = ["rfcache_opertional_mode1", "rfcache_opertional_mode2"]
    # rf_cache_page_size_kb = [1,2]
    # rf_cache_max_io_size_kb = [1,2]
    # state = ["state1", "state2"]
    # name = ["name1", "name2"]
    # id = ["id1", "id2"]
  }
}

output "inputAll" {
  value = data.powerflex_protection_domain.filter.protection_domains
}
