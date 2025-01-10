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