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

# Get all template details present on the cluster
data "powerflex_template" "example1" {
}

# Get template details using the ID of the template.
data "powerflex_template" "example2" {
  template_ids = ["ID1", "ID2"]
}

# Get template details using the Name of the template.
data "powerflex_template" "example3" {
  template_names = ["Name_1", "Name_2"]
}

output "template_result" {
  value = data.powerflex_template.example1.template_details
}
