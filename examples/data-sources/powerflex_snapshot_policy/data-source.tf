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

# commands to run this tf file : terraform init && terraform apply --auto-approve

data "powerflex_snapshot_policy" "sp" {
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/
# data "powerflex_snapshot_policy" "sp2" {
#   filter{
#     id = ["896a535700000000"]
#     name = ["snap-create-test"]
#     snapshot_policy_state = ["enabled"]
#     auto_snapshot_creation_cadence_in_min = [5]
#     max_vtree_auto_snapshots = [2]
#     system_id = ["1234567890abcdef"]
#     num_of_source_volumes = [7]
#     num_of_expired_but_locked_snapshots = [5]
#     num_of_creation_failures = [3]
#     num_of_retained_snapshots_per_level = [2]
#     snapshot_access_mode = ["read_write"]
#     secure_snapshots = true
#     time_of_last_auto_snapshot = [5]
#     time_of_last_auto_snapshot_creation_failure = [5]
#     last_auto_snapshot_creation_failure_reason = ["reason_for_failure"]
#     last_auto_snapshot_failure_in_first_level = true
#     num_of_auto_snapshots = [5]
#     num_of_locked_snapshots = [5]
#   }
# }

output "spResult" {
  value = data.powerflex_snapshot_policy.sp
}

