---
page_title: "Creating Multiple Volumes with Count"
title: "Creating Multiple Volumes with Count"
linkTitle: "Creating Multiple Volumes with Count"
---

<!--
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

You can use the count meta-argument to create multiple volumes.

## Example

To create 7 different volumes using the following configuration:

```terraform
resource "powerflex_volume" "volumes" {
  count                  = 7
  name                   = "security-footage-${count.index}"
  protection_domain_name = "domain1"
  storage_pool_name      = "pool1"
  size                   = 8
  use_rm_cache           = true
  volume_type            = "ThinProvisioned"
  access_mode            = "ReadWrite"
}
```
