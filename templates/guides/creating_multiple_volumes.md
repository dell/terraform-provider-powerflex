---
page_title: "Creating Multiple Volumes with Count"
title: "Creating Multiple Volumes with Count"
linkTitle: "Creating Multiple Volumes with Count"
---

U can use the count meta-argument to create multiple volumes.

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
