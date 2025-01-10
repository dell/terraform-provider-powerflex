/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# Get all NVMe target details present on the cluster
data "powerflex_nvme_target" "example1" {
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_nvme_target" "nvme_target_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     software_version_info = ["^R4_5.*$"]
#   }
# }

# output "nvmeTargetFilterRegexResult"{
#  value = data.powerflex_nvme_target.nvme_target_filter_regex.nvme_target_details
# }

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
// Please note that NVMe over TCP is supported in PowerFlex 4.0 and later versions, therefore this datasource is not supported in PowerFlex 3.x.

data "powerflex_nvme_target" "example" {
  filter {
    name = ["name1", "name2"]
    # id                                   = ["ID1", "ID2"]
    # protection_domain_id                 = ["PD1", "PD2"]
    # storage_port                         = [12200]
    # nvme_port                            = [4420]
    # discovery_port                       = [8009]
    # sdt_state                            = ["Normal"]
    # mdm_connection_state                 = ["Connected"]
    # membership_state                     = ["Joined"]
    # fault_set_id                         = ["FS1"]
    # software_version_info                = ["Version"]
    # maintenance_state                    = ["NoMaintenance"]
    # authentication_error                 = ["None"]
    # persistent_discovery_controllers_num = [0]
    # system_id                            = ["systemID"]
  }
}

output "nvme_target_result" {
  value = data.powerflex_nvme_target.example.nvme_target_details
}
