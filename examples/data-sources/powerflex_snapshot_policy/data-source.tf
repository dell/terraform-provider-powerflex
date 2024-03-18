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
# Reads snapshot policy either by name or by id , if provided
# If both name and id is not provided , then it reads all the snapshot policies
# id and name can't be given together to fetch the snapshot policy

data "powerflex_snapshot_policy" "sp" {

  #name = "sample_snap_policy_1"
  id = "896a535700000000"
}

output "spResult" {
  value = data.powerflex_snapshot_policy.sp.snapshotpolicies
}

