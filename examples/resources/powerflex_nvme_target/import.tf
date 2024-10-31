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

// Get all of the existing NVMe targets
data "powerflex_nvme_target" "existing" {
}

// Import all of the NVMe targets
import {
  for_each = data.powerflex_nvme_target.existing.nvme_target_details
  to       = powerflex_nvme_target.this[each.key]
  id       = each.value.id
}

// Add them to the terraform state
resource "powerflex_nvme_target" "this" {
  count                = length(data.powerflex_nvme_target.existing.nvme_target_details)
  name                 = data.powerflex_nvme_target.existing.nvme_target_details[count.index].name
  protection_domain_id = data.powerflex_nvme_target.existing.nvme_target_details[count.index].protection_domain_id
  ip_list              = data.powerflex_nvme_target.existing.nvme_target_details[count.index].ip_list
}