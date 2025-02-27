/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# commands to run this tf file : terraform init && terraform apply
# This feature is only supported for PowerFlex 4.5 and above.

// Empty filter block will return all the replication pairs
data "powerflex_resource_credential" "rc" {}
output "rpResult" {
  value = data.powerflex_resource_credential.rc
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_resource_credential" "rc_filter" {
  filter {
    # id = ["id", "id2"]
    # label = ["label-1", "label-2"]
    # created_by = ["created_by-1", "created_by-2"]
    # domain = ["domain-1", "domain-2"]
    # type = ["type-1", "type-2"]
    # updated_by = ["updated_by-1", "updated_by-2"]
    # updated_date = ["updated_date-1", "updated_date-2"]
    # username = ["username-1", "username-2"]
  }
}
output "rcResultFiltered" {
  value = data.powerflex_resource_credential.rc_filter
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_resource_credential" "resource_credentialp_filter_regex" {
#   filter{
#     label = ["^System_.*$"]
#     id = ["^Powerflex.*$"]
#   }
# }