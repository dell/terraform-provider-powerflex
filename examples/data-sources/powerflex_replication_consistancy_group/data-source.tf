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
# If both name and id is not provided , then it reads all the replication conistancy group
# id and name can't be given together to fetch the replication conistancy group .

// Empty filter block will return all the replication conistancy group
data "powerflex_replication_consistancy_group" "rcg" {}
output "rcgResult" {
  value = data.powerflex_replication_consistancy_group.rcg
}

data "powerflex_replication_consistancy_group" "rcg_filter" {
  filter {
    # id = ["id", "id2"]
    # name = ["name-1", "name-2"]
    # remote_id = ["remote_id-1", "remote_id-2"]
    # rpo_in_seconds = [15]
    # protection_domain_id = ["protection_domain_id-1", "protection_domain_id-2"]
    # remote_protection_domain_id = ["remote_protection_domain_id-1", "remote_protection_domain_id-2"]
    # destination_system_id = ["destination_system_id-1", "destination_system_id-2"]
    # peer_mdm_id = ["peer_mdm_id-1", "peer_mdm_id-2"]
    # remote_mdm_id = ["remote_mdm_id-1", "remote_mdm_id-2"]
    # replication_direction = ["replication_direction-1", "replication_direction-2"]
    # curr_consist_mode = ["curr_consist_mode-1", "curr_consist_mode-2"]
    # freeze_state = ["freeze_state-1", "freeze_state-2"]
    # pause_mode = ["pause_mode-1", "pause_mode-2"]
    # lifetime_state = ["lifetime_state-1", "lifetime_state-2"]
    # snap_creation_in_progress = false
    # last_snap_group_id = ["last_snap_group_id-1", "last_snap_group_id-2"]
    # type = ["type-1", "type-2"]
    # disaster_recovery_state = ["disaster_recovery_state-1", "disaster_recovery_state-2"]
    # remote_disaster_recovery_state = ["remote_disaster_recovery_state-1", "remote_disaster_recovery_state-2"]
    # target_volume_access_mode = ["target_volume_access_mode-1", "target_volume_access_mode-2"]
    # failover_type = ["failover_type-1", "failover_type-2"]
    # failover_state = ["failover_state-1", "failover_state-2"]
    # active_local = false
    # active_remote = false
    # abstract_state = ["abstract_state-1", "abstract_state-2"]
    # error = [403]
    # local_activity_state = ["local_activity_state-1", "local_activity_state-2"]
    # remote_activity_state = ["remote_activity_state-1", "remote_activity_state-2"]
    # inactive_reason = [1]
  }
}
output "rcgResultFiltered" {
  value = data.powerflex_replication_consistancy_group.rcg_filter
}