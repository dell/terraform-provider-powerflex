---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerflex_sds Data Source - powerflex"
subcategory: ""
description: |-
  Fetches the list of sds.
---

# powerflex_sds (Data Source)

Fetches the list of sds.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `protection_domain_id` (String) Protection Domain ID.
- `protection_domain_name` (String) Protection Domain Name.
- `sds_ids` (List of String) List of SDS IDs.
- `sds_names` (List of String) List of SDS names.

### Read-Only

- `id` (String) Placeholder identifier attribute.
- `sds_details` (Attributes List) List of SDS. (see [below for nested schema](#nestedatt--sds_details))

<a id="nestedatt--sds_details"></a>
### Nested Schema for `sds_details`

Read-Only:

- `authentication_error` (String) Authentication error.
- `certificate_info` (Attributes) Certificate Information. (see [below for nested schema](#nestedatt--sds_details--certificate_info))
- `configured_drl_mode` (String) Configured DRL mode.
- `drl_mode` (String) DRL mode.
- `faultset_id` (String) Fault set ID.
- `fgl_metadata_cache_size` (Number) FGL metadata cache size.
- `fgl_metadata_cache_state` (String) FGL metadata cache state.
- `fgl_num_concurrent_writes` (Number) FGL concurrent writes.
- `id` (String) SDS ID.
- `ip_list` (Attributes List) List of IPs associated with SDS. (see [below for nested schema](#nestedatt--sds_details--ip_list))
- `last_upgrade_time` (Number) Last time SDS was upgraded.
- `links` (Attributes List) Specifies the links asscociated with SDS. (see [below for nested schema](#nestedatt--sds_details--links))
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
- `rmcache_memory_allocation_state` (String) Indicates the state of the memory allocation process (in progress/done).
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

