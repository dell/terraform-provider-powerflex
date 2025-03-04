/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.
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

//Gather all existing snapshot policies
data "powerflex_snapshot_policy" "all" {
}

//Import all snapshot policies
import{
    for_each =  data.powerflex_snapshot_policy.all.snapshot_policy
    to = powerflex_snapshot_policy.import_test_snapshot_policy[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_snapshot_policy" "import_test_snapshot_policy" {
    name                                  = data.powerflex_snapshot_policy.all.snapshot_policy[count.index].name
    num_of_retained_snapshots_per_level   = data.powerflex_snapshot_policy.all.snapshot_policy[count.index].num_of_retained_snapshots_per_level
    auto_snapshot_creation_cadence_in_min = data.powerflex_snapshot_policy.all.snapshot_policy[count.index].auto_snapshot_creation_cadence_in_min
    snapshot_access_mode                  = data.powerflex_snapshot_policy.all.snapshot_policy[count.index].snapshot_access_mode
    secure_snapshots                      = data.powerflex_snapshot_policy.all.snapshot_policy[count.index].secure_snapshots
}