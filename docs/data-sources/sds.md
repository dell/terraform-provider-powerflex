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

title: "powerflex_sds data source"
linkTitle: "powerflex_sds"
page_title: "powerflex_sds Data Source - powerflex"
subcategory: ""
description: |-
  This datasource is used to query the existing Storage Data Servers from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.
---

# powerflex_sds (Data Source)

This datasource is used to query the existing Storage Data Servers from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.

> **Note:** Exactly one of `protection_domain_name` and `protection_domain_id` is required.

> **Note:** Only one of `sds_names` and `sds_ids` can be provided at a time.

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

data "powerflex_sds" "example1" {
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_sds" "sds_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     maintenance_type = ["^.*Maintenance$"]
#   }
# }

# output "sdsFilterRegexResult"{
#  value = data.powerflex_sds.sds_filter_regex.sds_details
# }

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
  value = data.powerflex_sds.example1.sds_details
}
```

After the successful execution of above said block, We can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_sds.example2.attribute_name` where attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) Placeholder identifier attribute.
- `sds_details` (Attributes List) List of fetched SDS. (see [below for nested schema](#nestedatt--sds_details))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `authentication_error` (Set of String) List of authentication_error
- `configured_drl_mode` (Set of String) List of configured_drl_mode
- `drl_mode` (Set of String) List of drl_mode
- `fault_set_id` (Set of String) List of fault_set_id
- `fgl_metadata_cache_size` (Set of Number) List of fgl_metadata_cache_size
- `fgl_metadata_cache_state` (Set of String) List of fgl_metadata_cache_state
- `fgl_num_concurrent_writes` (Set of Number) List of fgl_num_concurrent_writes
- `id` (Set of String) List of id
- `last_upgrade_time` (Set of Number) List of last_upgrade_time
- `maintenance_state` (Set of String) List of maintenance_state
- `maintenance_type` (Set of String) List of maintenance_type
- `mdm_connection_state` (Set of String) List of mdm_connection_state
- `membership_state` (Set of String) List of membership_state
- `name` (Set of String) List of name
- `num_io_buffers` (Set of Number) List of num_io_buffers
- `num_restarts` (Set of Number) List of num_restarts
- `on_vmware` (Boolean) Value for on_vmware
- `performance_profile` (Set of String) List of performance_profile
- `port` (Set of Number) List of port
- `rfcache_enabled` (Boolean) Value for rfcache_enabled
- `rfcache_error_api_version_mismatch` (Boolean) Value for rfcache_error_api_version_mismatch
- `rfcache_error_device_does_not_exist` (Boolean) Value for rfcache_error_device_does_not_exist
- `rfcache_error_inconsistent_cache_configuration` (Boolean) Value for rfcache_error_inconsistent_cache_configuration
- `rfcache_error_inconsistent_source_configuration` (Boolean) Value for rfcache_error_inconsistent_source_configuration
- `rfcache_error_invalid_driver_path` (Boolean) Value for rfcache_error_invalid_driver_path
- `rfcache_error_low_resources` (Boolean) Value for rfcache_error_low_resources
- `rmcache_enabled` (Boolean) Value for rmcache_enabled
- `rmcache_frozen` (Boolean) Value for rmcache_frozen
- `rmcache_memory_allocation_state` (Set of String) List of rmcache_memory_allocation_state
- `rmcache_size` (Set of Number) List of rmcache_size
- `sds_state` (Set of String) List of sds_state
- `software_version_info` (Set of String) List of software_version_info


<a id="nestedatt--sds_details"></a>
### Nested Schema for `sds_details`

Read-Only:

- `authentication_error` (String) Authentication error.
- `certificate_info` (Attributes) Certificate Information. (see [below for nested schema](#nestedatt--sds_details--certificate_info))
- `configured_drl_mode` (String) Configured DRL mode.
- `drl_mode` (String) DRL mode.
- `fault_set_id` (String) Fault set ID.
- `fgl_metadata_cache_size` (Number) FGL metadata cache size.
- `fgl_metadata_cache_state` (String) FGL metadata cache state.
- `fgl_num_concurrent_writes` (Number) FGL concurrent writes.
- `id` (String) SDS ID.
- `ip_list` (Attributes List) List of IPs associated with SDS. (see [below for nested schema](#nestedatt--sds_details--ip_list))
- `last_upgrade_time` (Number) Last time SDS was upgraded.
- `links` (Attributes List) Specifies the links associated with SDS. (see [below for nested schema](#nestedatt--sds_details--links))
- `maintenance_state` (String) Maintenance state.
- `maintenance_type` (String) Maintenance type.
- `mdm_connection_state` (String) MDM connection state.
- `membership_state` (String) Membership state.
- `name` (String) SDS name.
- `num_io_buffers` (Number) Number of IO buffers.
- `num_restarts` (Number) Number of restarts.
- `on_vmware` (Boolean) Presence on VMware.
- `performance_profile` (String) Performance profile.
- `port` (Number) SDS port.
- `raid_controllers` (Attributes List) RAID controllers information. (see [below for nested schema](#nestedatt--sds_details--raid_controllers))
- `rfcache_enabled` (Boolean) Whether RF cache is enabled or not.
- `rfcache_error_api_version_mismatch` (Boolean) RF cache error for API version mismatch.
- `rfcache_error_device_does_not_exist` (Boolean) RF cache error for device does not exist.
- `rfcache_error_inconsistent_cache_configuration` (Boolean) RF cache error for inconsistent cache configuration.
- `rfcache_error_inconsistent_source_configuration` (Boolean) RF cache error for inconsistent source configuration.
- `rfcache_error_invalid_driver_path` (Boolean) RF cache error for invalid driver path.
- `rfcache_error_low_resources` (Boolean) RF cache error for low resources.
- `rmcache_enabled` (Boolean) Whether RM cache is enabled or not.
- `rmcache_frozen` (Boolean) Indicates whether the Read RAM Cache is currently temporarily not in use.
- `rmcache_memory_allocation_state` (String) Indicates the state of the memory allocation process. Can be one of `in progress` and `done`.
- `rmcache_size` (Number) Indicates the size of Read RAM Cache on the specified SDS in KB.
- `sds_configuration_failure` (Attributes) SDS configuration failure windows. (see [below for nested schema](#nestedatt--sds_details--sds_configuration_failure))
- `sds_decoupled` (Attributes) SDS decoupled windows. (see [below for nested schema](#nestedatt--sds_details--sds_decoupled))
- `sds_receive_buffer_allocation_failures` (Attributes) SDS receive buffer allocation failure windows. (see [below for nested schema](#nestedatt--sds_details--sds_receive_buffer_allocation_failures))
- `sds_state` (String) SDS state.
- `software_version_info` (String) Software version information.

<a id="nestedatt--sds_details--certificate_info"></a>
### Nested Schema for `sds_details.certificate_info`

Read-Only:

- `issuer` (String) Certificate issuer.
- `subject` (String) Certificate subject.
- `thumbprint` (String) Certificate thumbprint.
- `valid_from` (String) The start date of the certificate validity.
- `valid_from_asn1_format` (String) The start date of the Asn1 format.
- `valid_to` (String) The end date of the certificate validity.
- `valid_to_asn1_format` (String) The end date of the Asn1 format.


<a id="nestedatt--sds_details--ip_list"></a>
### Nested Schema for `sds_details.ip_list`

Read-Only:

- `ip` (String) SDS IP.
- `role` (String) SDS IP role.


<a id="nestedatt--sds_details--links"></a>
### Nested Schema for `sds_details.links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details.
- `rel` (String) Specifies the relationship with the SDS.


<a id="nestedatt--sds_details--raid_controllers"></a>
### Nested Schema for `sds_details.raid_controllers`

Read-Only:

- `battery_status` (String) Battery status
- `driver_name` (String) Driver name.
- `driver_version` (String) Driver version.
- `firmware_version` (String) Firmware version.
- `model_name` (String) Model name.
- `pci_address` (String) PCI address.
- `serial_number` (String) Serial number.
- `status` (String) RAID status.
- `vendor_name` (String) Vendor name.


<a id="nestedatt--sds_details--sds_configuration_failure"></a>
### Nested Schema for `sds_details.sds_configuration_failure`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_configuration_failure--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_configuration_failure--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_configuration_failure--short_window))

<a id="nestedatt--sds_details--sds_configuration_failure--long_window"></a>
### Nested Schema for `sds_details.sds_configuration_failure.long_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_configuration_failure--medium_window"></a>
### Nested Schema for `sds_details.sds_configuration_failure.medium_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_configuration_failure--short_window"></a>
### Nested Schema for `sds_details.sds_configuration_failure.short_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--sds_details--sds_decoupled"></a>
### Nested Schema for `sds_details.sds_decoupled`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_decoupled--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_decoupled--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_decoupled--short_window))

<a id="nestedatt--sds_details--sds_decoupled--long_window"></a>
### Nested Schema for `sds_details.sds_decoupled.long_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_decoupled--medium_window"></a>
### Nested Schema for `sds_details.sds_decoupled.medium_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_decoupled--short_window"></a>
### Nested Schema for `sds_details.sds_decoupled.short_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--sds_details--sds_receive_buffer_allocation_failures"></a>
### Nested Schema for `sds_details.sds_receive_buffer_allocation_failures`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_receive_buffer_allocation_failures--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_receive_buffer_allocation_failures--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--sds_details--sds_receive_buffer_allocation_failures--short_window))

<a id="nestedatt--sds_details--sds_receive_buffer_allocation_failures--long_window"></a>
### Nested Schema for `sds_details.sds_receive_buffer_allocation_failures.long_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_receive_buffer_allocation_failures--medium_window"></a>
### Nested Schema for `sds_details.sds_receive_buffer_allocation_failures.medium_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--sds_details--sds_receive_buffer_allocation_failures--short_window"></a>
### Nested Schema for `sds_details.sds_receive_buffer_allocation_failures.short_window`

Read-Only:

- `last_oscillation_count` (Number) Last oscillation count.
- `last_oscillation_time` (Number) Last oscillation time.
- `max_failures_count` (Number) Maximum failures count.
- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


