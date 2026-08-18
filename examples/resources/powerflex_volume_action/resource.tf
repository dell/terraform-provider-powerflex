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

# Refresh volume metadata
resource "powerflex_volume_action" "refresh" {
  volume_id = "vol-abc123"
  action    = "refresh"
}

# Restore volume from a snapshot
resource "powerflex_volume_action" "restore" {
  volume_id          = "vol-abc123"
  action             = "restore"
  source_snapshot_id = "snap-def456"
}

# Map a volume to a host with read-write access
resource "powerflex_volume_action" "map" {
  volume_id   = "vol-abc123"
  action      = "map_to_host"
  host_id     = "host-ghi789"
  access_mode = "READ_WRITE"
}

# Unmap a volume from a host
resource "powerflex_volume_action" "unmap" {
  volume_id = "vol-abc123"
  action    = "unmap_from_host"
  host_id   = "host-ghi789"
}
