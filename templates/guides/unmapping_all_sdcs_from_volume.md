---
page_title: "Unmapping all SDCs from volumes"
title: "Unmapping all SDCs from volumes"
linkTitle: "Unmapping all SDCs from volumes"
---

<!--
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
-->

In order to unmap all SDCs from volumes or snapshots, there is a small trick.

Suppose you created a volume as

```terraform
resource "powerflex_volume" "volume" {
  name                   = "security-footage"
  protection_domain_name = "domain1"
  storage_pool_name      = "pool1"
  size                   = 8
  use_rm_cache           = true
  volume_type            = "ThinProvisioned"
  access_mode            = "ReadWrite"
  sdc_list               = [
    {
        sdc_name = "sdc1"
        limit_iops = 14o
        limit_bw_in_mbps = 20
        access_mode = "ReadWrite"
    }
  ]
}
```

Now you change your config to

```terraform
resource "powerflex_volume" "volume" {
  name                   = "security-footage"
  protection_domain_name = "domain1"
  storage_pool_name      = "pool1"
  size                   = 8
  use_rm_cache           = true
  volume_type            = "ThinProvisioned"
  access_mode            = "ReadWrite"
}
```

Then you would simply get a message saying
```
No changes. Your infrastructure matches the configuration.
```

You have to change your configuartion to

```terraform
resource "powerflex_volume" "volume" {
  name                   = "security-footage"
  protection_domain_name = "domain1"
  storage_pool_name      = "pool1"
  size                   = 8
  use_rm_cache           = true
  volume_type            = "ThinProvisioned"
  access_mode            = "ReadWrite"
  sdc_list               = []
}
```
Then your last SDC will get unmapped. Happy unmapping!
