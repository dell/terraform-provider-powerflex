

---
page_title: "Installing a SDC on a PowerFlex cluster using multiple values for IP"
title: "Installing a SDC on a PowerFlex cluster using multiple values for IP"
linkTitle: "Installing a SDC on a PowerFlex cluster using multiple values for IP"
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

# Steps for installing a SDC on a PowerFlex cluster using multiple values for IP

This guide explains how to install a SDC in cas of multiple values for IP.

### Example

```terraform
resource "powerflex_sdc" "sdc-example" {
  mdm_password = "Password
  lia_password = "Password"
  sdc_details = [
    {
      ip                  = "10.10.10.1,10.10.10.2,10.10.10.3" 
      username            = "user"
      password            = "password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Primary"
      is_sdc              = "Yes"
      name                = "SDC_PRIMARY"
      performance_profile = "HighPerformance"
      virtual_ips         = "10.10.10.4,10.10.10.5"
      virtual_ip_nics     = "eth1,eth2"
      data_network_ip     = "10.10.10.2,10.10.10.3"
    },
    {
      ip                  = "10.10.10.6,10.10.10.7,10.10.10.8" 
      username            = "user"
      password            = "password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Secondary"
      is_sdc              = "Yes"
      name                = "SDC_SECONDARY"
      performance_profile = "HighPerformance"
      virtual_ips         = "10.10.10.9,10.10.10.10"
      virtual_ip_nics     = "eth1,eth2"
      data_network_ip     = "10.10.10.7,10.10.10.8"
    },
    {
      ip                  = "10.10.10.11,10.10.10.12,10.10.10.13" 
      username            = "user"
      password            = "password"
      operating_system    = "linux"
      is_mdm_or_tb        = "TB"
      is_sdc              = "Yes"
      name                = "SDC_TB"
      performance_profile = "Compact"
      data_network_ip     = "10.10.10.12,10.10.10.13"
    },
    {
      ip               = "10.10.10.14"
      username         = "user"
      password         = "password"
      operating_system = "linux"
      is_mdm_or_tb     = "Standby"
      is_sdc           = "Yes"
    }
  ]
}

```


