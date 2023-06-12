/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# Commands to run this tf file : terraform init && terraform plan && terraform apply
# Only importing the sdc resource or renaming of sdc resource is supported
# For Renaming , id and name are required fields
# For importing , please check sdc_resource_import.tf file for more details
# name can't be empty


resource "powerflex_sdc" "sdc" {
  id   = "e3ce1fb500000000"
  name = "terraform_sdc"
}


#output "changed_sdc" {
# value = powerflex_sdc.sdc
#}
# # -----------------------------------------------------------------------------------


