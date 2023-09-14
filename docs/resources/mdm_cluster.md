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

title: "powerflex_mdm_cluster resource"
linkTitle: "powerflex_mdm_cluster"
page_title: "powerflex_mdm_cluster Resource - powerflex"
subcategory: ""
description: |-
  This resource can be used to manage MDM cluster on a PowerFlex array. Supports adding or removing standby MDMs, migrate from 3-node to 5-node cluster or vice-versa, changing MDM ownership, changing performance profile, and renaming MDMs.
---

# powerflex_mdm_cluster (Resource)

This resource can be used to manage MDM cluster on a PowerFlex array. Supports adding or removing standby MDMs, migrate from 3-node to 5-node cluster or vice-versa, changing MDM ownership, changing performance profile, and renaming MDMs.

!> **Caution:** MDM cluster creation or update is not atomic. In case of partially completed create operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
If the taint is not removed, terraform will destroy and recreate the resource.

~> **Note:** Use of MDM cluster resource requires the presence of the MDM cluster. The purpose of this resource is to update MDM cluster, not create or delete. Import operation is not supported for MDM cluster.

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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, and Delete operations are supported for this resource.
# For this resource, primary_mdm, secondary_mdm, tiebreaker_mdm and cluster_mode are mandatory parameters.
# While specifying primary_mdm, secondary_mdm and tiebreaker_mdm, id or ips is mandatory.
# For adding standby mdm, ips and role parameters are mandatory.
# To check which attributes of the MDM cluster resource can be updated, please refer Product Guide in the documentation

# Example for adding standby MDMs. Before adding standby MDMs, MDM package must be installed on VM with respective role. 
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "7f328d0b71711801"
  }
  secondary_mdm = [{
    id = "10912d8a5d412800"
  }]
  tiebreaker_mdm = [{
    id = "0e4f0a2f5978ae02"
  }]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for changing MDM ownership. In this case, the id of the primary and secondary MDM will be swapped.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [{
    id = "7f328d0b71711801"
  }]
  tiebreaker_mdm = [{
    id = "0e4f0a2f5978ae02"
  }]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for switching cluster mode to 5 nodes from 3 nodes. The cluster mode will be FiveNodes. Previously added standby MDMs will be added as Secondary/TieBreaker MDM.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "FiveNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    },
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
  standby_mdm = [
  ]
}

# Example for switching cluster mode to 3 nodes from 5 nodes. The cluster mode will be ThreeNodes. One of the active Secondary/TieBreaker MDM will be moved to standby MDMs.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    }
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    }
  ]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for removing standby MDMs. In this case, standby_mdm will be empty list.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    }
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    }
  ]
  standby_mdm = []
}
```

After the execution of the above resource block, the MDM cluster would become a 3-node cluster. For more information, please check the state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_mode` (String) Mode of the MDM cluster. Accepted values are `ThreeNodes` and `FiveNodes`.
- `primary_mdm` (Attributes) Primary MDM details. (see [below for nested schema](#nestedatt--primary_mdm))
- `secondary_mdm` (Attributes List) Secondary MDM details. (see [below for nested schema](#nestedatt--secondary_mdm))
- `tiebreaker_mdm` (Attributes List) TieBreaker MDM details. (see [below for nested schema](#nestedatt--tiebreaker_mdm))

### Optional

- `performance_profile` (String) Performance profile of the MDM cluster. Accepted values are `Compact` and `HighPerformance`.
- `standby_mdm` (Attributes List) StandBy MDM details. StandBy MDM can be added/removed/promoted to manager/tiebreaker role. (see [below for nested schema](#nestedatt--standby_mdm))

### Read-Only

- `id` (String) Unique identifier of the MDM cluster.

<a id="nestedatt--primary_mdm"></a>
### Nested Schema for `primary_mdm`

Optional:

- `id` (String) ID of the primary MDM.
- `ips` (Set of String) The Ips of the primary MDM.
- `name` (String) Name of the the primary MDM.

Read-Only:

- `management_ips` (Set of String) The management ips of the primary MDM.
- `port` (Number) Port of the primary MDM.


<a id="nestedatt--secondary_mdm"></a>
### Nested Schema for `secondary_mdm`

Optional:

- `id` (String) ID of the secondary MDM.
- `ips` (Set of String) The Ips of the secondary MDM.
- `name` (String) Name of the the secondary MDM.

Read-Only:

- `management_ips` (Set of String) The management ips of the secondary MDM.
- `port` (Number) Port of the secondary MDM.


<a id="nestedatt--tiebreaker_mdm"></a>
### Nested Schema for `tiebreaker_mdm`

Optional:

- `id` (String) ID of the tiebreaker MDM.
- `ips` (Set of String) The Ips of the tiebreaker MDM.
- `name` (String) Name of the the tiebreaker MDM.

Read-Only:

- `management_ips` (Set of String) The management ips of the tiebreaker MDM.
- `port` (Number) Port of the tiebreaker MDM.


<a id="nestedatt--standby_mdm"></a>
### Nested Schema for `standby_mdm`

Required:

- `ips` (Set of String) The Ips of the standby MDM. Cannot be updated.
- `role` (String) Role of the standby mdm. Accepted values are `Manager` and `TieBreaker`. Cannot be updated.

Optional:

- `allow_asymmetric_ips` (Boolean) Allow the added MDM to have a different number of IPs from the primary MDM. Cannot be updated.
- `management_ips` (Set of String) The management ips of the standby MDM. Cannot be updated.
- `name` (String) Name of the the standby MDM.
- `port` (Number) Port of the standby MDM. Cannot be updated.

Read-Only:

- `id` (String) ID of the standby MDM.

