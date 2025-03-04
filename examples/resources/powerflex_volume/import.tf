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

//Gather all existing volumes
data "powerflex_volume" "all"{
}

//Gather all existing protection domains
data "powerflex_protection_domain" "all" {
}

//Import all volumes
import {
    for_each = data.powerflex_volume.all.volumes
    to = powerflex_volume.import_test_volume[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_volume" "import_test_volume" {
    count = length(data.powerflex_volume.all.volumes)
    name = data.powerflex_volume.all.volumes[count.index].name
    protection_domain_name = data.powerflex_protection_domain.all.protection_domains[count.index].name # Value is not gathered from volume datasource
    storage_pool_id = data.powerflex_volume.all.volumes[count.index].storage_pool_id
    size          = floor(data.powerflex_volume.all.volumes[count.index].size_in_kb / 1000000)
    capacity_unit = "GB" # GB/TB
    use_rm_cache = data.powerflex_volume.all.volumes[count.index].use_rm_cache
    volume_type  = data.powerflex_volume.all.volumes[count.index].volume_type
}