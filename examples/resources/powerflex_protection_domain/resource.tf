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

resource "powerflex_protection_domain" "pd" {
  # required paramaters  ======

  name = "domain_1"

  # optional parameters  ======

  active = true

  # SDS IOPS throttling
  # overall_io_network_throttling_in_kbps must be greater than the rest of the parameters
  # 0 indicates unlimited IOPS
  protected_maintenance_mode_network_throttling_in_kbps = 10 * 1024
  rebuild_network_throttling_in_kbps                    = 10 * 1024
  rebalance_network_throttling_in_kbps                  = 10 * 1024
  vtree_migration_network_throttling_in_kbps            = 10 * 1024
  overall_io_network_throttling_in_kbps                 = 20 * 1024

  # Fine granularity metadata caching
  fgl_metadata_cache_enabled      = true
  fgl_default_metadata_cache_size = 1024

  # Read Flash cache
  rf_cache_enabled          = true
  rf_cache_operational_mode = "ReadAndWrite"
  rf_cache_page_size_kb     = 16
  rf_cache_max_io_size_kb   = 32
}
