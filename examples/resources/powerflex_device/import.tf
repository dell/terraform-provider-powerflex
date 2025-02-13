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

//Gather all existing devices
data "powerflex_device" "all" {
}

//Import all devices
import {
    for_each = data.powerflex_device.all.device_model
    to = powerflex_device.import_test_device[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_device" "import_test_device" {
    count = length(data.powerflex_device.all.device_model)
    name = data.powerflex_device.all.device_model[count.index].name
    sds_id = data.powerflex_device.all.device_model[count.index].sds_id
    storage_pool_id = data.powerflex_device.all.device_model[count.index].storage_pool_id
    device_path                = data.powerflex_device.all.device_model[count.index].device_current_path_name
    media_type                 = data.powerflex_device.all.device_model[count.index].media_type
    external_acceleration_type = data.powerflex_device.all.device_model[count.index].external_acceleration_type
}