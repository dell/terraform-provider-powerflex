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
terraform {
  required_providers {
    powerflex = {
      source = "registry.terraform.io/dell/powerflex"
    }
  }
}

provider "powerflex" {
  username = var.username_destination
  password = var.password_destination
  endpoint = var.endpoint_destination
  insecure = true
  timeout  = 120
}

// Get all of the exiting replication consistency groups
data "powerflex_replication_consistency_group" "all_current_rcgs" {
}

// Import all of the replication consistency groups
import {
    for_each = data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details
    to = powerflex_replication_consistency_group.imported_rcgs[each.key]
    id = each.value.id
}

// Add them to the terraform state
resource "powerflex_replication_consistency_group" "imported_rcgs" {
    count = length(data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details)
    name = data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details[count.index].name
    protection_domain_id = data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details[count.index].protection_domain_id
    remote_protection_domain_id = data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details[count.index].remote_protection_domain_id
    destination_system_id = data.powerflex_replication_consistency_group.all_current_rcgs.replication_consistency_group_details[count.index].destination_system_id
}
