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
      source                = "registry.terraform.io/dell/powerflex"
      configuration_aliases = [powerflex.source, powerflex.destination]
    }
  }
}

provider "powerflex" {
  alias    = "source"
  username = var.username_source
  password = var.password_source
  endpoint = var.endpoint_source
  insecure = true
  timeout  = 120
}

provider "powerflex" {
  alias    = "destination"
  username = var.username_destination
  password = var.password_destination
  endpoint = var.endpoint_destination
  insecure = true
  timeout  = 120
}

data "powerflex_protection_domain" "source_protection_domain" {
  provider = powerflex.source
  name     = var.source_protection_domain_name
}

data "powerflex_protection_domain" "destination_protection_domain" {
  provider = powerflex.destination
  name     = var.destination_protection_domain_name
}

resource "powerflex_replication_consistency_group" "example" {
  provider = powerflex.source
  ### Required Values

  # New name of the Replication Consistency Group
  name = var.name
  # Protection domain on the source machine
  protection_domain_id = data.powerflex_protection_domain.source_protection_domain.protection_domains[0].id
  # Protection domain on the destination machine
  remote_protection_domain_id = data.powerflex_protection_domain.destination_protection_domain.protection_domains[0].id
  # Destination System ID
  destination_system_id = data.powerflex_protection_domain.destination_protection_domain.protection_domains[0].system_id

  ### Optional with defaults if unset

  # Must be greater the 15 or less then 3600 (seconds) Default: 15
  #rpo_in_seconds = 15
  # Sets the Active state, Options (Active, Terminated) Default: Active
  #local_activity_state = "Active"
  # Sets the Volume Access Mode, Options (NoAccess, ReadOnly) Default: NoAccess
  #target_volume_access_mode = "NoAccess"
  # Pause Mode, Options (Resume, Pause) Default: Resume
  #pause_mode = "Resume"
  # Freeze State, Options (Unfrozen, Frozen) Default: Unfrozen
  #freeze_state = "Unfrozen"
  # Current Consistency Mode, Options (Consistent, Inconsistent) Default: Consistent
  #curr_consist_mode = "Consistent"
}
