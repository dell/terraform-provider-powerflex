---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

listIgnoreTitle: true
weight: 1
title: "powerflex provider"
linkTitle: "Provider"
page_title: "powerflex Provider"
subcategory: ""
description: |-
  The Terraform provider for Dell PowerFlex can be used to interact with a Dell PowerFlex array in order to manage the array resources.
---

# powerflex Provider

The Terraform provider for Dell PowerFlex can be used to interact with a Dell PowerFlex array in order to manage the array resources.

## Example Usage

provider.tf
```terraform
/*
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
*/

terraform {
  required_providers {
    powerflex = {
      version = "1.5.0"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}
provider "powerflex" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = true
  timeout  = 120
}

resource "powerflex_protection_domain" "pd" {
  name = "domain_1"
}

resource "powerflex_sds" "sds1" {
  name                 = "sds_1"
  protection_domain_id = powerflex_protection_domain.pd.id
  ip_list = [
    {
      ip   = "10.10.10.1"
      role = "all"
    }
  ]
  depends_on = [powerflex_protection_domain.pd]
}

resource "powerflex_sds" "sds2" {
  name                 = "sds_2"
  protection_domain_id = powerflex_protection_domain.pd.id
  ip_list = [
    {
      ip   = "10.10.10.2"
      role = "all"
    }
  ]
  depends_on = [powerflex_protection_domain.pd]
}

resource "powerflex_sds" "sds3" {
  name                 = "sds_3"
  protection_domain_id = powerflex_protection_domain.pd.id
  ip_list = [
    {
      ip   = "10.10.10.3"
      role = "all"
    }
  ]
  depends_on = [powerflex_protection_domain.pd]
}

resource "powerflex_storage_pool" "sp" {
  name                 = "SP"
  protection_domain_id = powerflex_protection_domain.pd.id
  media_type           = "HDD"
  use_rmcache          = true
  use_rfcache          = true
}

resource "powerflex_device" "device1" {
  name                       = "device1"
  device_path                = "/dev/sdb"
  sds_id                     = powerflex_sds.sds1.id
  storage_pool_id            = powerflex_storage_pool.sp.id
  media_type                 = "HDD"
  external_acceleration_type = "ReadAndWrite"
  depends_on                 = [powerflex_storage_pool.sp]
}

resource "powerflex_device" "device2" {
  name                       = "device2"
  device_path                = "/dev/sdb"
  sds_id                     = powerflex_sds.sds2.id
  storage_pool_id            = powerflex_storage_pool.sp.id
  media_type                 = "HDD"
  external_acceleration_type = "ReadAndWrite"
  depends_on                 = [powerflex_storage_pool.sp]
}

resource "powerflex_device" "device3" {
  name                       = "device3"
  device_path                = "/dev/sdb"
  sds_id                     = powerflex_sds.sds3.id
  storage_pool_id            = powerflex_storage_pool.sp.id
  media_type                 = "HDD"
  external_acceleration_type = "ReadAndWrite"
  depends_on                 = [powerflex_storage_pool.sp]
}

resource "powerflex_volume" "volume" {
  name                 = "volume1"
  protection_domain_id = powerflex_protection_domain.pd.id
  storage_pool_id      = powerflex_storage_pool.sp.id
  size                 = 16
  volume_type          = "ThinProvisioned"
  depends_on           = [powerflex_device.device1, powerflex_device.device2, powerflex_device.device3]
}

resource "powerflex_sdc_volumes_mapping" "map" {
  id = "e3d105e900000005"
  volume_list = [
    {
      volume_id        = powerflex_volume.volume.id
      limit_iops       = 140
      limit_bw_in_mbps = 19
      access_mode      = "ReadOnly"
    },
  ]
  depends_on = [powerflex_volume.volume]
}

resource "powerflex_package" "upload-test" {
  file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
  "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm"]
}

resource "powerflex_user" "user" {
  name     = "NewUser"
  role     = "Monitor"
  password = "Password123"
}
```

variables.tf
```terraform
/*
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
*/

variable "username" {
  type        = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type        = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
  type        = string
  description = "Stores the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `endpoint` (String) The PowerFlex Gateway server URL (inclusive of the port).
- `password` (String, Sensitive) The password required for the authentication.
- `username` (String) The username required for authentication.

### Optional

- `insecure` (Boolean) Specifies if the user wants to skip SSL verification.
- `timeout` (Number) HTTPS timeout.
