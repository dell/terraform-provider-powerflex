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

title: "powerflex_device data source"
linkTitle: "powerflex_device"
page_title: "powerflex_device Data Source - powerflex"
subcategory: ""
description: |-
  This datasource is used to query the existing device from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.
---

# powerflex_device (Data Source)

This datasource is used to query the existing device from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.

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
# empty block of the powerflex_device datasource will give list of all device within the system

data "powerflex_device" "all" {
}

output "deviceResult" {
  value = data.powerflex_device.all.device_model
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_device" "device_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     temperature_state = ["^.*Failed$"]
#   }
# }

# output "faultSetFilterRegexResult"{
#  value = data.powerflex_device.device_filter_regex.device_model
# }

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_device" "filter" {
  filter {
    # fgl_nvdimm_metadata_amortization_x100 = [1,2]
    # logical_sector_size_in_bytes = [1,2]
    # fgl_nvdimm_write_cache_size = [1,2]
    # acceleration_pool_id = ["acceleration_pool_id1", "acceleration_pool_id2"]
    # sds_id = ["sds_id1", "sds_id2"]
    # storage_pool_id = ["storage_pool_id1", "storage_pool_id2"]
    # capacity_limit_in_kb = [1,2]
    # error_state = ["error_state1", "error_state2"]
    # capacity = [1,2]
    # device_type = ["device_type1", "device_type2"]
    # persistent_checksum_state = ["persistent_checksum_state1", "persistent_checksum_state2"]
    # device_state = ["device_state1", "device_state2"]
    # led_setting = ["led_setting1", "led_setting2"]
    # max_capacity_in_kb = [1,2]
    # sp_sds_id = ["sp_sds_id1", "sp_sds_id2"]
    # aggregated_state = ["aggregated_state1", "aggregated_state2"]
    # temperature_state = ["temperature_state1", "temperature_state2"]
    # ssd_end_of_life_state = ["ssd_end_of_life_state1", "ssd_end_of_life_state2"]
    # model_name = ["model_name1", "model_name2"]
    # vendor_name = ["vendor_name1", "vendor_name2"]
    # raid_controller_serial_number = ["raid_controller_serial_number1", "raid_controller_serial_number2"]
    # firmware_version = ["firmware_version1", "firmware_version2"]
    # cache_look_ahead_active = false
    # write_cache_active = false
    # ata_security_active = false
    # physical_sector_size_in_bytes = [1,2]
    # media_failing = false
    # slot_number = ["slot1", "slot2"]
    # external_acceleration_type = ["external_acceleration_type1", "external_acceleration_type2"]
    # auto_detect_media_type = ["auto_detect_media_type1", "auto_detect_media_type2"]
    # device_current_path_name = ["device_current_path_name1", "device_current_path_name2"]
    # device_original_path_name = ["device_original_path_name1", "device_original_path_name2"]
    # rfcache_error_device_does_not_exist = false
    # media_type = ["SSD", "HDD"]
    # serial_number = ["serial_number1", "serial_number2"]
    # name = ["name1", "name2"]
    # id = ["id1", "id2"]
  }
}

output "deviceFilterResult" {
  value = data.powerflex_device.filter.device_model
}
```

After the successful execution of above said block, We can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_device.dev.attribute_name` where attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `device_model` (Attributes List) List of devices fetched. (see [below for nested schema](#nestedatt--device_model))
- `id` (String) Placeholder id of device datasource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `acceleration_pool_id` (Set of String) List of acceleration_pool_id
- `aggregated_state` (Set of String) List of aggregated_state
- `ata_security_active` (Boolean) Value for ata_security_active
- `auto_detect_media_type` (Set of String) List of auto_detect_media_type
- `cache_look_ahead_active` (Boolean) Value for cache_look_ahead_active
- `capacity` (Set of Number) List of capacity
- `capacity_limit_in_kb` (Set of Number) List of capacity_limit_in_kb
- `device_current_path_name` (Set of String) List of device_current_path_name
- `device_original_path_name` (Set of String) List of device_original_path_name
- `device_state` (Set of String) List of device_state
- `device_type` (Set of String) List of device_type
- `error_state` (Set of String) List of error_state
- `external_acceleration_type` (Set of String) List of external_acceleration_type
- `fgl_nvdimm_metadata_amortization_x100` (Set of Number) List of fgl_nvdimm_metadata_amortization_x100
- `fgl_nvdimm_write_cache_size` (Set of Number) List of fgl_nvdimm_write_cache_size
- `firmware_version` (Set of String) List of firmware_version
- `id` (Set of String) List of id
- `led_setting` (Set of String) List of led_setting
- `logical_sector_size_in_bytes` (Set of Number) List of logical_sector_size_in_bytes
- `max_capacity_in_kb` (Set of Number) List of max_capacity_in_kb
- `media_failing` (Boolean) Value for media_failing
- `media_type` (Set of String) List of media_type
- `model_name` (Set of String) List of model_name
- `name` (Set of String) List of name
- `persistent_checksum_state` (Set of String) List of persistent_checksum_state
- `physical_sector_size_in_bytes` (Set of Number) List of physical_sector_size_in_bytes
- `raid_controller_serial_number` (Set of String) List of raid_controller_serial_number
- `rfcache_error_device_does_not_exist` (Boolean) Value for rfcache_error_device_does_not_exist
- `sds_id` (Set of String) List of sds_id
- `serial_number` (Set of String) List of serial_number
- `slot_number` (Set of String) List of slot_number
- `sp_sds_id` (Set of String) List of sp_sds_id
- `ssd_end_of_life_state` (Set of String) List of ssd_end_of_life_state
- `storage_pool_id` (Set of String) List of storage_pool_id
- `temperature_state` (Set of String) List of temperature_state
- `vendor_name` (Set of String) List of vendor_name
- `write_cache_active` (Boolean) Value for write_cache_active


<a id="nestedatt--device_model"></a>
### Nested Schema for `device_model`

Read-Only:

- `acceleration_pool_id` (String) Acceleration Pool_id Of The Device Instance.
- `acceleration_props` (Attributes) Acceleration Props Of The Device Instance. (see [below for nested schema](#nestedatt--device_model--acceleration_props))
- `aggregated_state` (String) Aggregated State Of The Device Instance.
- `ata_security_active` (Boolean) Ata Security Active Of The Device Instance.
- `auto_detect_media_type` (String) Auto Detect Media Type Of The Device Instance.
- `cache_look_ahead_active` (Boolean) Cache Look Ahead Active Of The Device Instance.
- `capacity` (Number) Capacity Of The Device Instance.
- `capacity_limit_in_kb` (Number) Capacity Limit In Kb Of The Device Instance.
- `device_current_path_name` (String) Device Current Path Name Of The Device Instance.
- `device_original_path_name` (String) Device Original Path Name Of The Device Instance.
- `device_state` (String) State Of The Device Instance.
- `device_type` (String) Device Type Of The Device Instance.
- `error_state` (String) Error State Of The Device Instance.
- `external_acceleration_type` (String) External Acceleration Type Of The Device Instance.
- `fgl_nvdimm_metadata_amortization_x100` (Number) Fgl Nvdimm Metadata Amortization X100 Of The Device Instance.
- `fgl_nvdimm_write_cache_size` (Number) Fgl Nvdimm Write Cache Size Of The Device Instance.
- `firmware_version` (String) Firmware Version Of The Device Instance.
- `id` (String) Unique ID Of The Device Instance.
- `led_setting` (String) LED Setting Of The Device Instance.
- `links` (Attributes List) Underlying REST API links. (see [below for nested schema](#nestedatt--device_model--links))
- `logical_sector_size_in_bytes` (Number) Logical Sector Size In Bytes Of The Device Instance.
- `long_successful_ios` (Attributes) Long Successful Ios Of The Device Instance. (see [below for nested schema](#nestedatt--device_model--long_successful_ios))
- `max_capacity_in_kb` (Number) Max Capacity In Kb Of The Device Instance.
- `media_failing` (Boolean) Media Failing Of The Device Instance.
- `media_type` (String) Media Type Of The Device Instance.
- `model_name` (String) Model Name Of The Device Instance.
- `name` (String) Name Of The Device Instance.
- `persistent_checksum_state` (String) Persistent Checksum State Of The Device Instance.
- `physical_sector_size_in_bytes` (Number) Physical Sector Size In Bytes Of The Device Instance.
- `raid_controller_serial_number` (String) Raid Controller Serial Number Of The Device Instance.
- `rfcache_error_device_does_not_exist` (Boolean) Rfcache Error Device Does Not Exist Of The Device Instance.
- `rfcache_props` (Attributes) Rfcache Props Of The Device Instance. (see [below for nested schema](#nestedatt--device_model--rfcache_props))
- `sds_id` (String) Sds ID Of The Device Instance.
- `serial_number` (String) Serial Number Of The Device Instance.
- `slot_number` (String) Slot Number Of The Device Instance.
- `sp_sds_id` (String) Sp Sds Id Of The Device Instance.
- `ssd_end_of_life_state` (String) Ssd End Of Life State Of The Device Instance.
- `storage_pool_id` (String) Storage Pool ID Of The Device Instance.
- `storage_props` (Attributes) Storage Props Of The Device Instance. (see [below for nested schema](#nestedatt--device_model--storage_props))
- `temperature_state` (String) Temperature State Of The Device Instance.
- `vendor_name` (String) Vendor Name Of The Device Instance.
- `write_cache_active` (Boolean) Write Cache Active Of The Device Instance.

<a id="nestedatt--device_model--acceleration_props"></a>
### Nested Schema for `device_model.acceleration_props`

Read-Only:

- `acc_used_capacity_in_kb` (String) Accelerator(ACC) Used Capacity In KB Acceleration Properties Parameters Of The Device Instance.


<a id="nestedatt--device_model--links"></a>
### Nested Schema for `device_model.links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details.
- `rel` (String) Specifies the relationship with the Protection Domain.


<a id="nestedatt--device_model--long_successful_ios"></a>
### Nested Schema for `device_model.long_successful_ios`

Read-Only:

- `long_window` (Attributes) Long Window Parameters. (see [below for nested schema](#nestedatt--device_model--long_successful_ios--long_window))
- `medium_window` (Attributes) Medium Window Parameters. (see [below for nested schema](#nestedatt--device_model--long_successful_ios--medium_window))
- `short_window` (Attributes) Short Window Parameters. (see [below for nested schema](#nestedatt--device_model--long_successful_ios--short_window))

<a id="nestedatt--device_model--long_successful_ios--long_window"></a>
### Nested Schema for `device_model.long_successful_ios.long_window`

Read-Only:

- `last_oscillation_count` (Number) Last Oscillation Count Window Parameters Of The Device Instance.
- `last_oscillation_time` (Number) Last Oscillation Time Window Parameters Of The Device Instance.
- `max_failures_count` (Number) Max Failures Count Window Parameters Of The Device Instance.
- `threshold` (Number) Threshold Window Parameters Of The Device Instance.
- `window_size_in_sec` (Number) Window Size in seconds Window Parameters Of The Device Instance.


<a id="nestedatt--device_model--long_successful_ios--medium_window"></a>
### Nested Schema for `device_model.long_successful_ios.medium_window`

Read-Only:

- `last_oscillation_count` (Number) Last Oscillation Count Window Parameters Of The Device Instance.
- `last_oscillation_time` (Number) Last Oscillation Time Window Parameters Of The Device Instance.
- `max_failures_count` (Number) Max Failures Count Window Parameters Of The Device Instance.
- `threshold` (Number) Threshold Window Parameters Of The Device Instance.
- `window_size_in_sec` (Number) Window Size in seconds Window Parameters Of The Device Instance.


<a id="nestedatt--device_model--long_successful_ios--short_window"></a>
### Nested Schema for `device_model.long_successful_ios.short_window`

Read-Only:

- `last_oscillation_count` (Number) Last Oscillation Count Window Parameters Of The Device Instance.
- `last_oscillation_time` (Number) Last Oscillation Time Window Parameters Of The Device Instance.
- `max_failures_count` (Number) Max Failures Count Window Parameters Of The Device Instance.
- `threshold` (Number) Threshold Window Parameters Of The Device Instance.
- `window_size_in_sec` (Number) Window Size in seconds Window Parameters Of The Device Instance.



<a id="nestedatt--device_model--rfcache_props"></a>
### Nested Schema for `device_model.rfcache_props`

Read-Only:

- `device_uuid` (String) Device UUID RfCache Parameters Of The Device Instance.
- `rfcache_error_card_io_error` (Boolean) Rfcache Error Card Io error RfCache Parameters Of The Device Instance.
- `rfcache_error_heavy_load_cache_skip` (Boolean) rfcache_error_heavy_load_cache_skip RfCache Parameters Of The Device Instance.
- `rfcache_error_stuck_io` (Boolean) Rfcache Error Stuck Io RfCache Parameters Of The Device Instance.


<a id="nestedatt--device_model--storage_props"></a>
### Nested Schema for `device_model.storage_props`

Read-Only:

- `checksum_acc_device_id` (String) Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.
- `checksum_changelog_acc_device_id` (String) Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.
- `checksum_changelog_size_mb` (Number) Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.
- `checksum_mode` (String) Checksum Mode Storage Properties Parameters Of The Device Instance.
- `checksum_size_mb` (Number) Checksum Size MB Storage Properties Parameters Of The Device Instance.
- `dest_checksum_acc_device_id` (String) Destination Checksum Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.
- `dest_checksum_changelog_acc_device_id` (String) Destination Checksum Changelog Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.
- `dest_checksum_changelog_size_mb` (Number) Destination Checksum Changelog Size MB Storage Properties Parameters Of The Device Instance.
- `dest_checksum_mode` (String) Destination Checksum Mode Storage Properties Parameters Of The Device Instance.
- `dest_fgl_acc_device_id` (String) Destination FGL(Fujitsu General Limited) Accelerator(ACC) Device ID Storage Properties Parameters Of The Device Instance.
- `dest_fgl_nvdimm_size_mb` (Number) Destination FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.
- `fgl_acc_device_id` (String) FGL(Fujitsu General Limited) Accelerator(ACC) Device Id Storage Properties Parameters Of The Device Instance.
- `fgl_nvdimm_size_mb` (Number) FGL(Fujitsu General Limited) Non-Volatile Dual In-line Memory Module(NVDIMM) Size In MB Storage Properties Parameters Of The Device Instance.
- `is_checksum_fully_calculated` (Boolean) Is Checksum Fully Calculated Storage Properties Parameters Of The Device Instance.


