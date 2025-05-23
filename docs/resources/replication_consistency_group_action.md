---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerflex_replication_consistency_group_action resource"
linkTitle: "powerflex_replication_consistency_group_action"
page_title: "powerflex_replication_consistency_group_action Resource - powerflex"
subcategory: "Data Protection"
description: |-
  This resource is used to execute actions on the Replication Consistency Group entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.
---

# powerflex_replication_consistency_group_action (Resource)

This resource is used to execute actions on the Replication Consistency Group entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.

## Example Usage

```terraform
/*
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
*/

terraform {
  required_providers {
    powerflex = {
      source = "registry.terraform.io/dell/powerflex"
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

// Resource to manage lifecycle
resource "terraform_data" "always_run" {
  input = timestamp()
}

data "powerflex_replication_consistency_group" "rcg_filter" {
  filter {
    name = [var.replication_consistency_group_name]
  }
}

resource "powerflex_replication_consistency_group_action" "example" {
  # Required

  # Id of the replication consistency group
  id = data.powerflex_replication_consistency_group.rcg_filter.replication_consistency_group_details[0].id

  # Action to be performed on the replication consistency group.
  # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
  action = var.action

  // This will allow terraform create process to trigger each time we run terraform apply.
  lifecycle {
    replace_triggered_by = [
      terraform_data.always_run
    ]
  }
}
```

After the execution of above resource block, The selected Action will be preformed on the selected RCG.
Each time the RCG Actions resource is run it will preform the action again. 
For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Replication Consistency Group ID

### Optional

- `action` (String) Replication Consistency Group Action
