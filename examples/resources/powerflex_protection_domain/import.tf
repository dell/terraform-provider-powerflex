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

//Gather all exsisting protection domains
data "powerflex_protection_domain" "all" {
}

//import all protection domains
import {
    for_each = data.powerflex_protection_domain.all.protection_domains
    id = powerflex_protection_domain.import_test_protection_domain[each.key]
    to = each.value.id
}

//Add them to terraform state
resource "powerflex_protection_domain" "import_test_protection_domain" {
    count = length(data.powerflex_protection_domain.all.protection_domains)
    name = data.powerflex_protection_domain.all.protection_domains[count.index].name
    rebuild_network_throttling_in_kbps = data.powerflex_protection_domain.all.protection_domains[count.index].rebuild_network_throttling_in_kbps
    rebalance_network_throttling_in_kbps  = data.powerflex_protection_domain.all.protection_domains[count.index].rebalance_network_throttling_in_kbps
    vtree_migration_network_throttling_in_kbps = data.powerflex_protection_domain.all.protection_domains[count.index].vtree_migration_network_throttling_in_kbps
    overall_io_network_throttling_in_kbps = data.powerflex_protection_domain.all.protection_domains[count.index].overall_io_network_throttling_in_kbps
    fgl_metadata_cache_enabled  = data.powerflex_protection_domain.all.protection_domains[count.index].fgl_metadata_cache_enabled
    fgl_default_metadata_cache_size = data.powerflex_protection_domain.all.protection_domains[count.index].fgl_default_metadata_cache_size
    rf_cache_enabled = data.powerflex_protection_domain.all.protection_domains[count.index].rf_cache_enabled
    rf_cache_operational_mode = data.powerflex_protection_domain.all.protection_domains[count.index].rf_cache_operational_mode
    rf_cache_page_size_kb  = data.powerflex_protection_domain.all.protection_domains[count.index].rf_cache_page_size_kb
    rf_cache_max_io_size_kb  = data.powerflex_protection_domain.all.protection_domains[count.index].rf_cache_max_io_size_kb
}