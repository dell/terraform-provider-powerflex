---
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerflex_vtree data source"
linkTitle: "powerflex_vtree"
page_title: "powerflex_vtree Data Source - powerflex"
subcategory: ""
description: |-
  This datasource can be used to fetch information related to VTrees from a PowerFlex array.
---

# powerflex_vtree (Data Source)

This datasource can be used to fetch information related to VTrees from a PowerFlex array.

## Example Usage

```terraform
/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# commands to run this tf file : terraform init && terraform apply --auto-approve
# empty block of the powerflex_vtree datasource will give list of all VTrees within the system

data "powerflex_vtree" "example1" {
}

# Get VTree details using VTree IDs
data "powerflex_vtree" "example2" {
  vtree_ids = ["VTree_ID1", "VTree_ID2"]
}

# Get VTree details using Volume IDs
data "powerflex_vtree" "example3" {
  volume_ids = ["Volume_ID1", "Volume_ID2"]
}

# Get VTree details using Volume Names
data "powerflex_vtree" "example4" {
  volume_names = ["Volume_Name1", "Volume_Name2"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `volume_ids` (Set of String) List of volume IDs
- `volume_names` (Set of String) List of volume names
- `vtree_ids` (Set of String) List of VTree IDs

### Read-Only

- `id` (String) Placeholder identifier attribute.
- `vtree_details` (Attributes Set) VTree details (see [below for nested schema](#nestedatt--vtree_details))

<a id="nestedatt--vtree_details"></a>
### Nested Schema for `vtree_details`

Read-Only:

- `compression_method` (String) Compression method
- `data_layout` (String) Data layout
- `id` (String) VTree ID
- `in_deletion` (Boolean) In deletion
- `links` (Attributes List) Specifies the links asscociated with VTree (see [below for nested schema](#nestedatt--vtree_details--links))
- `name` (String) VTree name
- `root_volumes` (Set of String) Root volumes
- `storage_pool_id` (String) Storage pool ID
- `vtree_migration_info` (Attributes) Vtree migration information (see [below for nested schema](#nestedatt--vtree_details--vtree_migration_info))

<a id="nestedatt--vtree_details--links"></a>
### Nested Schema for `vtree_details.links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details
- `rel` (String) Specifies the relationship with the VTree


<a id="nestedatt--vtree_details--vtree_migration_info"></a>
### Nested Schema for `vtree_details.vtree_migration_info`

Read-Only:

- `destination_storage_pool_id` (String) Destination storage pool ID
- `migration_pause_reason` (String) Migration pause reason
- `migration_queue_position` (Number) Migration queue position
- `migration_status` (String) Migration status
- `source_storage_pool_id` (String) Source storage pool ID
- `thickness_conversion_type` (String) Thickness conversion type

