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

title: "powerflex_protection_domain data source"
linkTitle: "powerflex_protection_domain"
page_title: "powerflex_protection_domain Data Source - powerflex"
subcategory: "Storage Management"
description: |-
  This datasource is used to query the existing protection domain from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.
---

# powerflex_protection_domain (Data Source)

This datasource is used to query the existing protection domain from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.

> **Note:** Only one of `name` and `id` can be provided at a time.

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
```

After the successful execution of above said block, We can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_protection_domain.pd.attribute_name` where attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))
- `id` (String) Unique identifier of the protection domain instance. Conflicts with `name`.

### Read-Only

- `protection_domains` (Attributes List) List of protection domains fetched. (see [below for nested schema](#nestedatt--protection_domains))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `fgl_default_metadata_cache_size` (Set of Number) List of fgl_default_metadata_cache_size
- `fgl_default_num_concurrent_writes` (Set of Number) List of fgl_default_num_concurrent_writes
- `fgl_metadata_cache_enabled` (Boolean) Value for fgl_metadata_cache_enabled
- `id` (Set of String) List of id
- `name` (Set of String) List of name
- `overall_io_network_throttling_enabled` (Boolean) Value for overall_io_network_throttling_enabled
- `overall_io_network_throttling_in_kbps` (Set of Number) List of overall_io_network_throttling_in_kbps
- `protected_maintenance_mode_network_throttling_enabled` (Boolean) Value for protected_maintenance_mode_network_throttling_enabled
- `protected_maintenance_mode_network_throttling_in_kbps` (Set of Number) List of protected_maintenance_mode_network_throttling_in_kbps
- `rebalance_network_throttling_enabled` (Boolean) Value for rebalance_network_throttling_enabled
- `rebalance_network_throttling_in_kbps` (Set of Number) List of rebalance_network_throttling_in_kbps
- `rebuild_network_throttling_enabled` (Boolean) Value for rebuild_network_throttling_enabled
- `rebuild_network_throttling_in_kbps` (Set of Number) List of rebuild_network_throttling_in_kbps
- `replication_capacity_max_ratio` (Set of Number) List of replication_capacity_max_ratio
- `rf_cache_accp_id` (Set of String) List of rf_cache_accp_id
- `rf_cache_enabled` (Boolean) Value for rf_cache_enabled
- `rf_cache_max_io_size_kb` (Set of Number) List of rf_cache_max_io_size_kb
- `rf_cache_opertional_mode` (Set of String) List of rf_cache_opertional_mode
- `rf_cache_page_size_kb` (Set of Number) List of rf_cache_page_size_kb
- `state` (Set of String) List of state
- `system_id` (Set of String) List of system_id
- `vtree_migration_network_throttling_enabled` (Boolean) Value for vtree_migration_network_throttling_enabled
- `vtree_migration_network_throttling_in_kbps` (Set of Number) List of vtree_migration_network_throttling_in_kbps


<a id="nestedatt--protection_domains"></a>
### Nested Schema for `protection_domains`

Read-Only:

- `fgl_default_metadata_cache_size` (Number) Fine Granularity Metadata Cache size.
- `fgl_default_num_concurrent_writes` (Number) Fine Granularity default number of concurrent writes.
- `fgl_metadata_cache_enabled` (Boolean) Whether Fine Granularity Metadata Cache is enabled or not.
- `id` (String) Unique identifier of the protection domain instance.
- `links` (Attributes List) Underlying REST API links. (see [below for nested schema](#nestedatt--protection_domains--links))
- `mdm_sds_network_disconnections_counter` (Attributes) MDM-SDS Network Disconnection Counter windows. (see [below for nested schema](#nestedatt--protection_domains--mdm_sds_network_disconnections_counter))
- `name` (String) Unique name of the protection domain instance.
- `overall_io_network_throttling_enabled` (Boolean) Whether network throttling is enabled for overall io.
- `overall_io_network_throttling_in_kbps` (Number) Maximum allowed io for protected maintenance mode in KBps. Must be greater than any other network throttling parameter.
- `protected_maintenance_mode_network_throttling_enabled` (Boolean) Whether network throttling is enabled for protected maintenance mode.
- `protected_maintenance_mode_network_throttling_in_kbps` (Number) Maximum allowed io for protected maintenance mode in KBps.
- `rebalance_network_throttling_enabled` (Boolean) Whether network throttling is enabled for rebalancing.
- `rebalance_network_throttling_in_kbps` (Number) Maximum allowed io for rebalancing in KBps.
- `rebuild_network_throttling_enabled` (Boolean) Whether network throttling is enabled for rebuilding.
- `rebuild_network_throttling_in_kbps` (Number) Maximum allowed io for rebuilding in KBps.
- `replication_capacity_max_ratio` (Number) Maximum Replication Capacity Ratio.
- `rf_cache_accp_id` (String) ID of the Rf Cache Acceleration Pool attached to the PD.
- `rf_cache_enabled` (Boolean) Whether SDS Rf Cache is enabled or not.
- `rf_cache_max_io_size_kb` (Number) Maximum io of the SDS RF Cache in KB.
- `rf_cache_opertional_mode` (String) Operational Mode of the SDS RF Cache.
- `rf_cache_page_size_kb` (Number) Page size of the SDS RF Cache in KB.
- `sdr_sds_connectivity` (Attributes) SDR-SDS Connectivity information. (see [below for nested schema](#nestedatt--protection_domains--sdr_sds_connectivity))
- `sds_configuration_failure_counter` (Attributes) SDS Configuration Failure Counter windows. (see [below for nested schema](#nestedatt--protection_domains--sds_configuration_failure_counter))
- `sds_decoupled_counter` (Attributes) SDS Decoupled Counter windows. (see [below for nested schema](#nestedatt--protection_domains--sds_decoupled_counter))
- `sds_receive_buffer_allocation_failures_counter` (Attributes) SDS receive Buffer Allocation Failure Counter windows. (see [below for nested schema](#nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter))
- `sds_sds_network_disconnections_counter` (Attributes) SDS-SDS Network Disconnection Counter windows. (see [below for nested schema](#nestedatt--protection_domains--sds_sds_network_disconnections_counter))
- `state` (String) State of a PD. Valid values are `Active`, `ActivePending`, `Inactive` or `InactivePending`.
- `system_id` (String) System ID of the PD.
- `vtree_migration_network_throttling_enabled` (Boolean) Whether network throttling is enabled for vtree migration.
- `vtree_migration_network_throttling_in_kbps` (Number) Maximum allowed io for vtree migration in KBps.

<a id="nestedatt--protection_domains--links"></a>
### Nested Schema for `protection_domains.links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details.
- `rel` (String) Specifies the relationship with the Protection Domain.


<a id="nestedatt--protection_domains--mdm_sds_network_disconnections_counter"></a>
### Nested Schema for `protection_domains.mdm_sds_network_disconnections_counter`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--mdm_sds_network_disconnections_counter--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--mdm_sds_network_disconnections_counter--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--mdm_sds_network_disconnections_counter--short_window))

<a id="nestedatt--protection_domains--mdm_sds_network_disconnections_counter--long_window"></a>
### Nested Schema for `protection_domains.mdm_sds_network_disconnections_counter.long_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--mdm_sds_network_disconnections_counter--medium_window"></a>
### Nested Schema for `protection_domains.mdm_sds_network_disconnections_counter.medium_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--mdm_sds_network_disconnections_counter--short_window"></a>
### Nested Schema for `protection_domains.mdm_sds_network_disconnections_counter.short_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--protection_domains--sdr_sds_connectivity"></a>
### Nested Schema for `protection_domains.sdr_sds_connectivity`

Read-Only:

- `client_server_conn_status` (String) Connectivity Status.
- `disconnected_client_id` (String) ID of the disconnected client.
- `disconnected_client_name` (String) Name of the disconnected client.
- `disconnected_server_id` (String) ID of the disconnected server.
- `disconnected_server_ip` (String) IP address of the disconnected server.
- `disconnected_server_name` (String) Name of the disconnected server.


<a id="nestedatt--protection_domains--sds_configuration_failure_counter"></a>
### Nested Schema for `protection_domains.sds_configuration_failure_counter`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_configuration_failure_counter--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_configuration_failure_counter--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_configuration_failure_counter--short_window))

<a id="nestedatt--protection_domains--sds_configuration_failure_counter--long_window"></a>
### Nested Schema for `protection_domains.sds_configuration_failure_counter.long_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_configuration_failure_counter--medium_window"></a>
### Nested Schema for `protection_domains.sds_configuration_failure_counter.medium_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_configuration_failure_counter--short_window"></a>
### Nested Schema for `protection_domains.sds_configuration_failure_counter.short_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--protection_domains--sds_decoupled_counter"></a>
### Nested Schema for `protection_domains.sds_decoupled_counter`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_decoupled_counter--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_decoupled_counter--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_decoupled_counter--short_window))

<a id="nestedatt--protection_domains--sds_decoupled_counter--long_window"></a>
### Nested Schema for `protection_domains.sds_decoupled_counter.long_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_decoupled_counter--medium_window"></a>
### Nested Schema for `protection_domains.sds_decoupled_counter.medium_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_decoupled_counter--short_window"></a>
### Nested Schema for `protection_domains.sds_decoupled_counter.short_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter"></a>
### Nested Schema for `protection_domains.sds_receive_buffer_allocation_failures_counter`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--short_window))

<a id="nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--long_window"></a>
### Nested Schema for `protection_domains.sds_receive_buffer_allocation_failures_counter.long_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--medium_window"></a>
### Nested Schema for `protection_domains.sds_receive_buffer_allocation_failures_counter.medium_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_receive_buffer_allocation_failures_counter--short_window"></a>
### Nested Schema for `protection_domains.sds_receive_buffer_allocation_failures_counter.short_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.



<a id="nestedatt--protection_domains--sds_sds_network_disconnections_counter"></a>
### Nested Schema for `protection_domains.sds_sds_network_disconnections_counter`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_sds_network_disconnections_counter--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_sds_network_disconnections_counter--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--protection_domains--sds_sds_network_disconnections_counter--short_window))

<a id="nestedatt--protection_domains--sds_sds_network_disconnections_counter--long_window"></a>
### Nested Schema for `protection_domains.sds_sds_network_disconnections_counter.long_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_sds_network_disconnections_counter--medium_window"></a>
### Nested Schema for `protection_domains.sds_sds_network_disconnections_counter.medium_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


<a id="nestedatt--protection_domains--sds_sds_network_disconnections_counter--short_window"></a>
### Nested Schema for `protection_domains.sds_sds_network_disconnections_counter.short_window`

Read-Only:

- `threshold` (Number) Threshold.
- `window_size_in_sec` (Number) Window Size in seconds.


