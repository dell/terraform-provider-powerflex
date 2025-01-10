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

# Get all VTrees details present on the cluster
data "powerflex_vtree" "all" {
}

output "powerflex_vtree_all_result" {
  value = data.powerflex_vtree.all.vtree_details
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_vtree" "vtree_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     data_layout = ["^.*Granularity$"]
#   }
# }

# output "vtreeFilterRegexResult"{
#  value = data.powerflex_vtree.vtree_filter_regex.vtree_details
# }


# Get Peer System details using filter with all values
# If there is no intersection between the filters then an empty datasource will be returned
# For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples
data "powerflex_vtree" "filtered" {
  filter {
    # storage_pool_id = ["storage_pool_id", "storage_pool_id2"]
    # data_layout = ["data_layout", "data_layout2"]
    # compression_method = ["compression_method", "compression_method2"]
    # in_deletion = false
    # name = ["name", "name2"]
    # id = ["id", "id2"]
  }
}

output "powerflex_vtree_filtered_result" {
  value = data.powerflex_vtree.filtered.vtree_details
}