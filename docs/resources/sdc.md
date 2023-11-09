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

title: "powerflex_sdc resource"
linkTitle: "powerflex_sdc"
page_title: "powerflex_sdc Resource - powerflex"
subcategory: ""
description: |-
  This resource is used to manage the SDC entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex SDC using this resource. We can also Import an existing SDC from PowerFlex array.
---

# powerflex_sdc (Resource)

This resource is used to manage the SDC entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex SDC using this resource. We can also Import an existing SDC from PowerFlex array.


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
# Create, Update, Read, Delete and Import operations are supported for this resource.

# Example for adding MDMs as SDCs. After successful execution, three SDCs will be added.
resource "powerflex_sdc" "sdc-example" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details = [
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Primary"
      is_sdc              = "Yes"
      name                = "SDC_NAME"
      performance_profile = "HighPerformance"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Secondary"
      is_sdc              = "Yes"
      name                = "SDC_NAME"
      performance_profile = "Compact"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "TB"
      is_sdc              = "Yes"
      name                = "SDC_NAME"
      performance_profile = "Compact"
    }
  ]
}

# Example for deleting all MDMs installed as SDCs. After successful execution, SDCs will be removed from the cluster. 
resource "powerflex_sdc" "expansion" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details = []
}

# Example for installing non-MDM node as SDC. After successful execution, one SDC will be added.
resource "powerflex_sdc" "sdc-example" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details = [
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Primary"
      is_sdc              = "No"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Secondary"
      is_sdc              = "No"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "TB"
      is_sdc              = "No"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_sdc              = "Yes"
    }
  ]
}

# Example for renaming existing SDC using ID. After successful execution, SDC will be renamed.
data "powerflex_sdc" "all" {
}

locals {
  matching_sdc = [for sdc in data.powerflex_sdc.all.sdcs : sdc if sdc.sdc_ip == "IP address of the SDC node to get SDC ID"]
}

resource "powerflex_sdc" "test" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details = [
    {
      sdc_id = local.matching_sdc[0].id
      name = "rename_sdc"
    },
  ]
}
```

After the execution of above resource block, sdc would have been renamed on the PowerFlex array. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `lia_password` (String, Sensitive) LIA Password to connect MDM Server.
- `mdm_password` (String, Sensitive) MDM Password to connect MDM Server.
- `sdc_details` (Attributes List) List of SDC Expansion Server Details. (see [below for nested schema](#nestedatt--sdc_details))

### Read-Only

- `id` (String) Placeholder

<a id="nestedatt--sdc_details"></a>
### Nested Schema for `sdc_details`

Optional:

- `ip` (String, Sensitive) IP of the node. Conflict with `sdc_id`
- `is_mdm_or_tb` (String) Whether this works as MDM or Tie Breaker,The acceptable value are `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank
- `is_sdc` (String) Whether this node is to operate as an SDC or not. The acceptable values are `Yes` and `No`. Default value is `Yes`.
- `name` (String) Name of the SDC to manage.
- `operating_system` (String) Operating System on the node
- `password` (String, Sensitive) Password of the node
- `performance_profile` (String) Performance Profile of SDC, The acceptable value are `HighPerformance` or `Compact`.
- `sdc_id` (String) ID of the SDC to manage. This can be retrieved from the Datasource and PowerFlex Server. Cannot be updated. Conflict with `ip`
- `username` (String) Username of the node

Read-Only:

- `mdm_connection_state` (String) The MDM connection status of the fetched SDC.
- `on_vmware` (Boolean) If the fetched SDC is on vmware.
- `sdc_approved` (Boolean) If the fetched SDC is approved.
- `sdc_guid` (String) The GUID of the fetched SDC.
- `system_id` (String) The System ID of the fetched SDC.

## Import

Import is supported using the following syntax:

```shell
# /*
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import SDC by it's id
terraform import powerflex_sdc.sdc_data "<id>"

# import SDC by multiple id
terraform import powerflex_sdc.sdc_data "<id1>,<id2>,<id3>"
```

1. This will import the SDC instance with specified ID into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.