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

# commands to run this tf file : terraform init && terraform applye
# If both name and id is not provided , then it reads all the replication pairs
# id and name can't be given together to fetch the replication pairs .

// Empty filter block will return all the replication pairs
data "powerflex_replication_pair" "rp" {}
output "rpResult" {
  value = data.powerflex_replication_pair.rp
}

data "powerflex_replication_pair" "rp_filter" {
  filter {
    # id = ["id", "id2"]
    # name = ["name-1", "name-2"]
    # remote_id = ["remote_id-1", "remote_id-2"]
    # user_requested_pause_transmit_init_copy = false
    # remote_capacity_in_mb = [8192]
    # local_volume_id = ["local_volume_id-1", "local_volume_id-2"]
    # remote_volume_id = ["remote_volume_id-1", "remote_volume_id-2"]
    # remote_volume_name = ["remote_volume_name-1", "remote_volume_name-2"]
    # replication_consistency_group_id = ["replication_consistency_group_id-1", "replication_consistency_group_id-2"]
    # copy_type = ["copy_type-1", "copy_type-2"]
    # lifetime_state = ["lifetime_state-1", "lifetime_state-2"]
    # peer_system_name = ["peer_system_name-1", "peer_system_name-2"]
    # initial_copy_state = ["initial_copy_state-1", "initial_copy_state-2"]
    # initial_copy_priority = [0, 1]
  }
}
output "rpResultFiltered" {
  value = data.powerflex_replication_pair.rp_filter
}