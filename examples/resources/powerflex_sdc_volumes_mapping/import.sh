# Below are the steps to import sdc along with mapped volumes :
# Step 1 - To import a sdc , we need the id of that sdc
# Step 2 - To check the id of the sdc we can make use of sdc datasource . Please refer sdc_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_sdc_volumes_mapping" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_sdc_volumes_mapping.resource_block_name" "id_of_the_sdc" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file.

# /*
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import using SDC id
terraform import powerflex_sdc_volumes_mapping.sdc_mapping_import_by_id "<id>"



