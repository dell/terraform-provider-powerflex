/*
Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# Create, Update, Delete is supported for this resource
# To import, check import.tf for more info
# protection_domain_id is required
# name and ip_list are required parameters to create
# port and fault_set_id are optional
# This resource is for PowerFlex Gen2 systems (5.0+)

# Example for creating a Storage Node (Gen2 replacement for SDS)
resource "powerflex_storage_node" "create" {
  name                 = "demo-storage-node-01"
  protection_domain_id = "202a046600000000"
  ip_list = [
    {
      ip   = "10.10.10.12"
      role = "all" # all/sdsOnly/sdcOnly
    },
    {
      ip   = "10.10.10.11"
      role = "sdcOnly" # all/sdsOnly/sdcOnly
    },
  ]
}

# Example for creating a Storage Node with fault set
resource "powerflex_fault_set" "test" {
  protection_domain_id = "202a046600000000"
  name                 = "demo_fault_set"
}

resource "powerflex_storage_node" "create_with_fs" {
  name                 = "demo-storage-node-02"
  protection_domain_id = "202a046600000000"
  fault_set_id         = powerflex_fault_set.test.id
  port                 = "6611"
  ip_list = [
    {
      ip   = "10.10.10.13"
      role = "all"
    },
  ]
}

output "storage_node" {
  value = powerflex_storage_node.create
}
