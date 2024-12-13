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

# Get all template details present on the cluster
data "powerflex_template" "example1" {
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_template" "template" {
  # filter{
  #   template_name = ["template_name"]
  #   id = ["template_id"]
  #   template_type = ["template_type"]
  #   template_version = ["template_version"]
  #   original_template_id = ["original_template_id"]
  #   template_locked =  false
  #   in_configuration = false
  #   created_date = ["created_date"]
  #   created_by = ["created_by"]
  #   updated_date = ["updated_date"]
  #   last_deployed_date = ["last_deployed_date"]
  #   updated_by = ["updated_by"]
  #   manage_firmware = true
  #   use_default_catalog = true
  #   all_users_allowed = false
  #   category = ["category"]
  #   server_count = [3]
  #   storage_count = [0]
  #   cluster_count = [1]
  #   service_count = [0]
  #   switch_count = [2]
  #   vm_count = [0]
  #   sdnas_count = [3]
  #   brownfield_template_type = ["brownfield_template_type"]
  #   Draft = true
  # }
}

output "template_result" {
  value = data.powerflex_template.example1.template_details
}
