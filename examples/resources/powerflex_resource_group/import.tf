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

//Gather all exsisting resource groups
data "powerflex_resource_group" "all" {
}

//Gather all exsisting templates
data "powerflex_template" "all" {
}

//Import all resource groups
import {
    for_each = data.powerflex_resource_group.all.resource_group_details
    to = powerflex_resource_group.import_test_resource_group[each.key]
    id = each.value.id
}

//Add them to terraform state
resource "powerflex_resource_group" "import_test_resource_group" {
    count = length(data.powerflex_resource_group.all.resource_group_details)
    deployment_name        = data.powerflex_resource_group.all.resource_group_details[count.index].deployment_name
    deployment_description = data.powerflex_resource_group.all.resource_group_details[count.index].deployment_description
    template_id            = data.powerflex_template.all.template_details[count.index].id // Not gathered in Resource group datasource
    firmware_id            = data.powerflex_resource_group.all.resource_group_details[count.index].firmware_repository_id
}