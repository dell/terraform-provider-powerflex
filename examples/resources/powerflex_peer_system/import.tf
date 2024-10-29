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

// Get all of the exiting peer systems
data "powerflex_peer_system" "all_current_peer_systems" {
}

// Import all of the peers
import {
    for_each = data.powerflex_peer_system.all_current_peer_systems.peer_system_details
    to = powerflex_peer_system.imported_peer_systems[each.key]
    id = each.value.id
}

// Add them to the terraform state
resource "powerflex_peer_system" "imported_peer_systems" {
    count = length(data.powerflex_peer_system.all_current_peer_systems.peer_system_details)
    name = data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].name
    peer_system_id = data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].peer_system_id
    ip_list = [
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[0].ip,
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[1].ip,
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[2].ip,
    ]
}

