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

# Get all Resource Group details present on the Resource Group
data "powerflex_resource_group" "example1" {
}

# Get Resource Group details using the ID of the Resource Group
data "powerflex_resource_group" "example2" {
  resource_group_ids = ["ID1", "ID2"]
}

# Get Resource Group details using the Name of the Resource Group
data "powerflex_resource_group" "example3" {
  resource_group_names = ["Name_1", "Name_2"]
}

output "resource_group_result" {
  value = data.powerflex_resource_group.example1.resource_group_details
}
