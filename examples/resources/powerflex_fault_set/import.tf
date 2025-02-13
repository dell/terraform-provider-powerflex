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

//Gather all existing fault sets
data "powerflex_fault_set" "all" {
}

//Import all fault sets
import {
    for_each = data.powerflex_fault_set.all.fault_set_details
    to = powerflex_fault_set.import_test_fault_set[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_fault_set" "import_test_fault_set" {
    count = length(data.powerflex_fault_set.all.fault_set_details)
    name = data.powerflex_fault_set.all.fault_set_details[count.index].name
    protection_domain_id = data.powerflex_fault_set.all.fault_set_details[count.index].protection_domain_id
}