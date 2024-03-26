---
page_title: "Creating Multiple Volumes with For Each"
title: "Creating Multiple Volumes with For Each"
linkTitle: "Creating Multiple Volumes with For Each"
---

<!--
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
-->

You can use the for_each meta-argument to create multiple volumes.

## Example

To create 3 different volumes use the following configuration:

```terraform
resource "powerflex_volume" "volume_example" {
  for_each = toset( ["Volume1", "Volume2", "Volume3"] )
  name = each.key
  protection_domain_name = "domain1"
  storage_pool_name = "pool1"
  size          = 8
  capacity_unit = "GB"
}
```
