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

provider "powerflex" {
  username = var.username_destination
  password = var.password_destination
  endpoint = var.endpoint_destination
  insecure = true
  timeout  = 120
}

// Get all of the exiting replication pairs
data "powerflex_replication_pair" "existing" {
}

// Import all of the replication pairs
import {
    for_each = data.powerflex_replication_pair.existing.replication_pair_details
    to = powerflex_replication_pair.this[each.key]
    id = each.value.id
}

// Add them to the terraform state
resource "powerflex_replication_pair" "this" {
    count = length(data.powerflex_replication_pair.existing.replication_pair_details)
    name = data.powerflex_replication_pair.existing.replication_pair_details[count.index].name
    source_volume_id = data.powerflex_replication_pair.existing.replication_pair_details[count.index].local_volume_id
    destination_volume_id = data.powerflex_replication_pair.existing.replication_pair_details[count.index].remote_volume_id
    replication_consistency_group_id = data.powerflex_replication_pair.existing.replication_pair_details[count.index].replication_consistency_group_id
}