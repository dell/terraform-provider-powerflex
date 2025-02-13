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

//Gather all existing sds
data "powerflex_sds" "all" {
}

//Import all sds
import {
    for_each = data.powerflex_sds.all.sds_details
    to = powerflex_sds.import_test_sds[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_sds" "import_test_sds" {
    count  = length(data.powerflex_sds.all.sds_details)
    name = data.powerflex_sds.all.sds_details[count.index].name
    protection_domain_id = data.powerflex_sds.all.sds_details[count.index].protection_domain_id
    ip_list = data.powerflex_sds.all.sds_details[count.index].ip_list
}