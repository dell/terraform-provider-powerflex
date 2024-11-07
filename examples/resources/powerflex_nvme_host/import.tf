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

// Get all of the existing NVMe hosts
data "powerflex_nvme_host" "existing" {
}

// Import all of the NVMe hosts
import {
  for_each = data.powerflex_nvme_host.existing.nvme_host_details
  to       = powerflex_nvme_host.this[each.key]
  id       = each.value.id
}

// Add them to the terraform state
resource "powerflex_nvme_host" "this" {
  count             = length(data.powerflex_nvme_host.existing.nvme_host_details)
  name              = data.powerflex_nvme_host.existing.nvme_host_details[count.index].name
  nqn               = data.powerflex_nvme_host.existing.nvme_host_details[count.index].nqn
  max_num_paths     = data.powerflex_nvme_host.existing.nvme_host_details[count.index].max_num_paths
  max_num_sys_ports = data.powerflex_nvme_host.existing.nvme_host_details[count.index].max_num_sys_ports
}