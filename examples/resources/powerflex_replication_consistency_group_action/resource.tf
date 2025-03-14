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
