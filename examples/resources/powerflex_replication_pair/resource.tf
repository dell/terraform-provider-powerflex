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
      source  = "registry.terraform.io/dell/powerflex"
      configuration_aliases = [ powerflex.source, powerflex.destination ]
    }
  }
}

provider "powerflex" {
  alias = "source"
  username = var.username_source
  password = var.password_source
  endpoint = var.endpoint_source
  insecure = true
  timeout  = 120
}

provider "powerflex" {
  alias = "destination"
  username = var.username_destination
  password = var.password_destination
  endpoint = var.endpoint_destination
  insecure = true
  timeout  = 120
}

# Used to get the id of the source volume
data "powerflex_volume" "source_volume" {
  provider = powerflex.source
  name = var.volume_name_source
}

# Used to get the id of the destination volume
data "powerflex_volume" "destination_volume" {
  provider = powerflex.destination
  name = var.volume_name_destination
}

# Used to get the id of the replication consistancy group
data "powerflex_replication_consistancy_group" "source_replication_consistancy_group" {
  provider = powerflex.source
  filter {
    name = [var.replication_consistancy_group_name_source]
  }
}

# Commands to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# name, source_volume_id, destination_volume_id, replication_consistency_group_id are the required parameters to create or update
# Only pause_initial_copy can be updated to true or false

# Create Replication Pair
resource "powerflex_replication_pair" "example" {
  provider = powerflex.source
  # Required values
  name                              = "example-replication-pair"
  source_volume_id                  = data.powerflex_volume.source_volume.volumes[0].id
  destination_volume_id             = data.powerflex_volume.destination_volume.volumes[0].id
  replication_consistency_group_id  = data.powerflex_replication_consistancy_group.source_replication_consistancy_group.replication_consistency_group_details[0].id
  # Optional values
  # pause_initial_copy = true # Pauses the replication pair (This will only work during the initial copy process), defaults to false
}
