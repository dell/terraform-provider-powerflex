---
page_title: "Deploying a PowerFlex cluster using a CSV Topology File"
title: "Deploying a PowerFlex cluster using a CSV Topology File"
linkTitle: "Deploying a PowerFlex cluster using a CSV Topology File"
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

# Deploying a PowerFlex cluster using a CSV Topology File

This guide explains how to use CSV file(s) for managing PowerFlex cluster via Terraform's Cluster Resource.

### Example

~> **Note:** It is recommended to use 2 separate csv files for cluster node related information and storage pool related information. Refer the following examples for more information. The first row will be treated as header.


**cluster_node.csv**
```
IPs,Username,Password,Operating System,Is MDM/TB,Is SDS,SDS Storage Device List,Is SDC
10.76.60.1,root,Password1,linux,Primary,Yes,/dev/sdb,Yes
10.76.60.2,root,Password1,linux,Secondary,Yes,/dev/sdb,Yes
10.76.60.3,root,Password1,linux,TB,Yes,/dev/sdb,Yes
10.76.60.4,root,Password1,linux,,Yes,/dev/sdb,Yes
```

**sp_data.csv**
```
ProtectionDomain,StoragePool,Media Type,Replication journal capacity percentage
domain_1,pool1,HDD,50
```

To perform SDC operations with a CSV file, use the following configuration:

```terraform

locals {
  cluster_node_data = csvdecode(file("cluster_node.csv"))
  sp_data = csvdecode(file("sp_data.csv"))
}

resource "powerflex_cluster" "test" {
  # Security Related Fields
  mdm_password = "Password"
  lia_password = "Password"

  # Advance Security Configuration
  allow_non_secure_communication_with_lia = false
  allow_non_secure_communication_with_mdm = false
  disable_non_mgmt_components_auth = false

  # Cluster Configuration related fields
  cluster = [
    for row in local.cluster_node_data : {
      ips                  = row.IPs
      username             = row.Username
      password             = row.Password
      operating_system     = row["Operating System"]
      is_mdm_or_tb         = row["Is MDM/TB"]
      is_sds               = row["Is SDS"]
      sdc_name             = row["SDS Name"]
      perf_profile_for_sdc = row.perfProfileForSDC
      is_rfcache           = row["RFcache"]
      is_sdr               = row["Is SDR"]
    }
    if row["Is MDM/TB"] != ""
  ]

  storage_pools = [
    for row in local.sp_data : {
      media_type                              = row["Media Type"]
      protection_domain                       = row.ProtectionDomain
      storage_pool                            = row.StoragePool
      replication_journal_capacity_percentage = tonumber(row["Replication journal capacity percentage"])
    }
  ]
}
```
This Terraform configuration sets up a PowerFlex cluster by defining its configuration, including security settings, node information, and storage pool details. The data for this configuration is read from two CSV files using the csvdecode function and organized into the required format for the resource.

Proposed: Two Terraform local variables are getting declared here (cluster_node_data and sp_data) by using Terraform's built-in 'csvdecode' function.

`locals`: The scope of the variable is within the configuration file. These variables are used to store data or perform calculations within this configuration.

`cluster_node_data`: Reads and decodes the CSV file “cluster_node.csv” into a list of maps. This file contains information about cluster nodes.

`sp_data`: Reads and decodes the CSV file “sp_data.csv” into a list of maps. This file contains information about storage pools.

`cluster`: A list variable which is getting created using a for loop which iterates over each element in the 'cluster_node_data' (local list variable).

`storage_pools`: A list variable which is getting created using a for loop which iterates over each element in the 'sp_data' (local list variable defined earlier).

The fields inside the object (e.g., ip, password, operating_system, etc.) are populated with values from the CSV file. You can also apply some conditions in the if statement to filter out rows where IPs and Password are not empty and not equal to "null."