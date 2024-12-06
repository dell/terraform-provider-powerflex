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

# Get all fault set details present on the cluster
data "powerflex_fault_set" "all" {
}

output "fault_set_result_all" {
  value = data.powerflex_fault_set.all.fault_set_details
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/
data "powerflex_fault_set" "filtered" {
  filter {
    # protection_domain_id = ["protection_domain_id", "protection_domain_id2"]
    # name = ["name", "name2"]
    # id = ["id", "id2"]
  }
}

output "fault_set_result_filtered" {
  value = data.powerflex_fault_set.filtered.fault_set_details
}
