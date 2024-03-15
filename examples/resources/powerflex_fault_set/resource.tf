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

# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check import.sh for more info
# name and protection_domain_id is the required parameter to create or update 
# To check which attributes of the fault set can be updated, please refer Product Guide in the documentation


resource "powerflex_fault_set" "avengers-fs-create" {
  # Name of the fault set
  name = "avengers-fs-create2"

  # To create / update, protection_domain_id is required
  protection_domain_id = "202a046600000000"
}