---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerflex_firmware_repository data source"
linkTitle: "powerflex_firmware_repository"
page_title: "powerflex_firmware_repository Data Source - powerflex"
subcategory: ""
description: |-
  This datasource is used to query the existing firmware repository from the PowerFlex array. The information fetched from this datasource can be used for getting the necessary details regarding the bundles and their components in that firmware repository.
---

# powerflex_firmware_repository (Data Source)

This datasource is used to query the existing firmware repository from the PowerFlex array. The information fetched from this datasource can be used for getting the necessary details regarding the bundles and their components in that firmware repository.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `firmware_repository_ids` (Set of String) List of firmware repository IDs
- `firmware_repository_names` (Set of String) List of firmware repository names

### Read-Only

- `firmware_repository_details` (Attributes List) Firmware Repository details (see [below for nested schema](#nestedatt--firmware_repository_details))
- `id` (String) Placeholder attribute.

<a id="nestedatt--firmware_repository_details"></a>
### Nested Schema for `firmware_repository_details`

Read-Only:

- `bundle_count` (Number) Bundle Count
- `component_count` (Number) Component Count
- `created_by` (String) Created By
- `created_date` (String) Created Date
- `custom` (Boolean) Custom
- `default_catalog` (Boolean) Default Catalog
- `disk_location` (String) Disk Location
- `download_progress` (Number) Download Progress
- `download_status` (String) Download Status
- `embedded` (Boolean) Embedded
- `extract_progress` (Number) Extract Progress
- `file_size_in_gigabytes` (Number) File Size In Gigabytes
- `filename` (String) Filename
- `id` (String) ID of the Firmware Repository
- `job_id` (String) Job ID
- `minimal` (Boolean) Minimal
- `name` (String) Firmware Repository name
- `needs_attention` (Boolean) Needs Attention
- `password` (String) Password
- `rcmapproved` (Boolean) Rcmapproved
- `signature` (String) Signature
- `software_bundles` (Attributes List) Software Bundles (see [below for nested schema](#nestedatt--firmware_repository_details--software_bundles))
- `software_components` (Attributes List) Software Components (see [below for nested schema](#nestedatt--firmware_repository_details--software_components))
- `source_location` (String) Source Location
- `source_type` (String) Source Type
- `state` (String) State
- `updated_by` (String) Updated By
- `updated_date` (String) Updated Date
- `user_bundle_count` (Number) User Bundle Count
- `username` (String) Username

<a id="nestedatt--firmware_repository_details--software_bundles"></a>
### Nested Schema for `firmware_repository_details.software_bundles`

Read-Only:

- `bundle_date` (String) Bundle Date
- `bundle_type` (String) Bundle Type
- `created_by` (String) Created By
- `created_date` (String) Created Date
- `custom` (Boolean) Custom
- `description` (String) Description
- `device_model` (String) Device Model
- `device_type` (String) Device Type
- `fw_repository_id` (String) Fw Repository ID
- `id` (String) ID
- `name` (String) Name
- `needs_attention` (Boolean) Needs Attention
- `software_components` (Attributes List) Software Components (see [below for nested schema](#nestedatt--firmware_repository_details--software_bundles--software_components))
- `updated_by` (String) Updated By
- `updated_date` (String) Updated Date
- `user_bundle` (Boolean) User Bundle
- `user_bundle_path` (String) User Bundle Path
- `version` (String) Version

<a id="nestedatt--firmware_repository_details--software_bundles--software_components"></a>
### Nested Schema for `firmware_repository_details.software_bundles.software_components`

Read-Only:

- `category` (String) Category
- `component_id` (String) Component ID
- `component_type` (String) Component Type
- `created_by` (String) Created By
- `created_date` (String) Created Date
- `custom` (Boolean) Custom
- `dell_version` (String) Dell Version
- `device_id` (String) Device ID
- `firmware_repo_name` (String) Firmware Repo Name
- `hash_md5` (String) Hash Md5
- `id` (String) ID
- `ignore` (Boolean) Ignore
- `name` (String) Name
- `needs_attention` (Boolean) Needs Attention
- `operating_system` (String) Operating System
- `original_component_id` (String) Original Component ID
- `package_id` (String) Package ID
- `path` (String) Path
- `sub_device_id` (String) Sub Device ID
- `sub_vendor_id` (String) Sub Vendor ID
- `system_ids` (List of String) System IDs
- `updated_by` (String) Updated By
- `updated_date` (String) Updated Date
- `vendor_id` (String) Vendor ID
- `vendor_version` (String) Vendor Version



<a id="nestedatt--firmware_repository_details--software_components"></a>
### Nested Schema for `firmware_repository_details.software_components`

Read-Only:

- `category` (String) Category
- `component_id` (String) Component ID
- `component_type` (String) Component Type
- `created_by` (String) Created By
- `created_date` (String) Created Date
- `custom` (Boolean) Custom
- `dell_version` (String) Dell Version
- `device_id` (String) Device ID
- `firmware_repo_name` (String) Firmware Repo Name
- `hash_md5` (String) Hash Md5
- `id` (String) ID
- `ignore` (Boolean) Ignore
- `name` (String) Name
- `needs_attention` (Boolean) Needs Attention
- `operating_system` (String) Operating System
- `original_component_id` (String) Original Component ID
- `package_id` (String) Package ID
- `path` (String) Path
- `sub_device_id` (String) Sub Device ID
- `sub_vendor_id` (String) Sub Vendor ID
- `system_ids` (List of String) System IDs
- `updated_by` (String) Updated By
- `updated_date` (String) Updated Date
- `vendor_id` (String) Vendor ID
- `vendor_version` (String) Vendor Version


