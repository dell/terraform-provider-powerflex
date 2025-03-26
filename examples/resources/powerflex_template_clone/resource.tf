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

# terraform init && terraform plan && terraform apply
# Create, and read is supported for this resource

//gets original template id from sample templates
data "powerflex_template" "template" {
  filter{
    category = ["Sample Templates"]
  }
}

resource "powerflex_template_clone" "example" {
  template_name = "Template Clone"
  original_template_id=data.powerflex_template.template.template_details[0].original_template_id
}