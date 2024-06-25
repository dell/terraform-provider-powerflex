---
page_title: "Deploying a PowerFlex cluster using multiple values for mdm_ips"
title: "Deploying a PowerFlex cluster using multiple values for mdm_ips"
linkTitle: "Deploying a PowerFlex cluster using multiple values for mdm_ips"
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

# Deploying a PowerFlex cluster having multiple values for mdm_ips

This guide explains how to deploy a cluster which has multiple values for mdm_ips.

### Example

```terraform

resource "powerflex_cluster" "test" {
  mdm_password = "Password"
  lia_password = "Password"

  # Advance Security Configuration
  allow_non_secure_communication_with_lia = false
  allow_non_secure_communication_with_mdm = false
  disable_non_mgmt_components_auth        = false

  # Cluster Configuration related fields
  cluster = [
    {
      # MDM Configuration Fields
      username                 = "user",
      password                 = "password",
      operating_system         = "linux",
      is_mdm_or_tb             = "Primary",
      mdm_ips                  = "10.10.10.1,10.10.10.2,10.10.10.3",
      mdm_mgmt_ip              = "10.10.10.1",
      mdm_name                 = "mdm1",
      perf_profile_for_mdm     = "HighPerformance",
      virtual_ips              = "10.10.10.4,10.10.10.5",
      virtual_ip_nics          = "eth1,eth2",
      is_sds                   = "Yes",
      sds_name                 = "sds1",
      sds_all_ips              = "10.10.10.2,10.10.10.3",
      protection_domain        = "domain_1",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc1",
      perf_profile_for_sdc     = "HighPerformance"

    },
    {
      username                 = "user",
      password                 = "password",
      operating_system         = "linux",
      is_mdm_or_tb             = "Secondary",
      mdm_name                 = "mdm2",
      mdm_ips                  = "10.10.10.6,10.10.10.7,10.10.10.8",
      mdm_mgmt_ip              = "10.10.10.6",
      perf_profile_for_mdm     = "HighPerformance",
      virtual_ips              = "10.10.10.4,10.10.10.5",
      virtual_ip_nics          = "eth1,eth2",
      is_sds                   = "Yes",
      sds_name                 = "sds2",
      sds_all_ips              = "10.10.10.7,10.10.10.8",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      protection_domain        = "domain_1",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc2",
      perf_profile_for_sdc     = "HighPerformance"
    },
    {
      username                 = "user",
      password                 = "password",
      operating_system         = "linux",
      is_mdm_or_tb             = "TB",
      mdm_name                 = "tb1",
      mdm_ips                  = "10.10.10.9,10.10.10.10,10.01.10.11",
      mdm_mgmt_ip              = "10.10.10.9",
      perf_profile_for_mdm     = "HighPerformance",
      is_sds                   = "Yes",
      sds_name                 = "sds3",
      sds_all_ips              = "10.10.10.10,10.01.10.11",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      protection_domain        = "domain_1",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc3",
      perf_profile_for_sdc     = "HighPerformance"
    }
  ]

    storage_pools = [
    {
      media_type        = "SSD"
      protection_domain = "domain_1"
      storage_pool      = "pool1"
      daya_layout       = "MG"
      zero_padding      = "true"
    }
  ]
}

```
This Terraform configuration sets up a PowerFlex cluster having multiple values.