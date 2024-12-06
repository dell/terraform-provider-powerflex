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

data "powerflex_sds" "example1" {
}

# Get Sds details using filter with all values
# If there is no intersection between the filters then an empty datasource will be returned
# For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples
# data "powerflex_sds" "sds" {
#   filter{
    # authentication_error                            = ["None"]
    # configured_drl_mode                             = ["Volatile"]
    # drl_mode                                        = ["Volatile"]
    # fault_set_id                                    = ["fault_set_id1", "fault_set_id2"]
    # fgl_metadata_cache_size                         = [0]
    # fgl_metadata_cache_state                        = ["Enabled"]
    # fgl_num_concurrent_writes                       = [10000]
    # id                                              = ["id1", "id2"]
    # last_upgrade_time                               = [0]
    # maintenance_state                               = ["NoMaintenance"]
    # maintenance_type                                = ["NoMaintenance"]
    # mdm_connection_state                            = ["Connected"]
    # membership_state                                = ["Joined"]
    # name                                            = ["sds-env-name1", "sds-env-name2"]
    # num_io_buffers                                  = [0]
    # num_restarts                                    = [6]
    # on_vmware                                       = false
    # performance_profile                             = ["HighPerformance"]
    # port                                            = [8080]
    # rfcache_enabled                                 = true
    # rfcache_error_api_version_mismatch              = false
    # rfcache_error_device_does_not_exist             = false
    # rfcache_error_inconsistent_cache_configuration  = false
    # rfcache_error_inconsistent_source_configuration = false
    # rfcache_error_invalid_driver_path               = false
    # rfcache_error_low_resources                     = false
    # rmcache_enabled                                 = false
    # rmcache_frozen                                  = false
    # rmcache_memory_allocation_state                 = ["RmcacheDisabled"]
    # rmcache_size                                    = [1400]
    # sds_state                                       = ["Normal"]
    # software_version_info                           = ["R4_6.1234.0"]
#   }
# }

output "allsdcresult" {
  value = data.powerflex_sds.example1
}

