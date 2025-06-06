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

//Gather all existing firmware repositories
data "powerflex_firmware_repository" "all" {
}

//Import all firmware repositories
import {
    for_each = data.powerflex_firmware_repository.all.firmware_repository_details
    id = powerflex_firmware_repository.import_test_firmware_repository[each.key]
    to = each.value.id
}

//Add them to terraform state
resource "powerflex_firmware_repository" "import_test_firmware_repository" {
    count = length(data.powerflex_firmware_repository.all.firmware_repository_details)
    source_location = data.powerflex_firmware_repository.all.firmware_repository_details[count.index].source_location
    username        = data.powerflex_firmware_repository.all.firmware_repository_details[count.index].username
}