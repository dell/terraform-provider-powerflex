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

# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check import.sh for more info
# name, num_of_retained_snapshots_per_level, auto_snapshot_creation_cadence_in_min is the required parameter to create or update 


resource "powerflex_snapshot_policy" "sp" {
  name                                  = "snap-create"
  num_of_retained_snapshots_per_level   = [2, 6, 7]
  auto_snapshot_creation_cadence_in_min = 6
  volume_ids                            = ["edd2fb3100000007", "edd322270000000a"] # assigning or unassigning volumes to snapshot policy
  snapshot_access_mode                  = "ReadWrite"                              # Cannot be updated after creation. It only supports two values : ReadOnly / ReadWrite
  secure_snapshots                      = false                                    # Cannot be updated after creation
  #remove_mode = "Remove" #remove_mode is applicable while unassingning the volumes from the snapshot policy. It only supports two values : Remove / Detach
}