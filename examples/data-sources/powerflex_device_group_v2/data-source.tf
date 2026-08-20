# Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# List all device groups
data "powerflex_device_group_v2" "all" {
}

output "all_device_groups" {
  value = data.powerflex_device_group_v2.all.device_groups
}

# Get device group by ID
data "powerflex_device_group_v2" "by_id" {
  id = "abc123def456"
}

# Filter by name
data "powerflex_device_group_v2" "by_name" {
  name = "device-group-1"
}

# Filter by protection domain
data "powerflex_device_group_v2" "by_pd" {
  protection_domain_id = "pd-123456"
}
