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
      configuration_aliases = [powerflex.system_1, powerflex.system_2]
    }
  }
}

provider "powerflex" {
  alias    = "system_1"
  username = var.username_system_1
  password = var.password_system_1
  endpoint = "https://${var.endpoint_system_1}"
  insecure = true
  timeout  = 120
}

provider "powerflex" {
  alias    = "system_2"
  username = var.username_system_2
  password = var.password_system_2
  endpoint = "https://${var.endpoint_system_2}"
  insecure = true
  timeout  = 120
}

data "powerflex_protection_domain" "protection_domain_system_2" {
  provider = powerflex.system_2
  name     = var.protection_domain_name_system_2
}

data "powerflex_protection_domain" "protection_domain_system_1" {
  provider = powerflex.system_1
  name     = var.protection_domain_name_system_1
}


resource "powerflex_peer_system" "system_1" {
  provider = powerflex.system_1

  // This should be done in order to avoid a confict while sshing
  depends_on = [resource.powerflex_peer_system.system_1]
  ### Required Values

  # New name of the Peer System
  name = var.name
  # Peer System (System 2) ID
  peer_system_id = data.powerflex_protection_domain.protection_domain_system_2.protection_domains[0].system_id
  # List of Peer MDM Ips at the destination
  ip_list = var.mdm_ips_system_2

  ### Optional with defaults if unset


  # Add certificate flag, default: false. 
  # If true source_primary_mdm_information and destination_primary_mdm_information must be filled out in order to get and set the certificate
  #add_certificate = true

  # source_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_1
  #   management_username = var.username_system_1
  #   management_password = var.password_system_1
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # destination_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_2
  #   management_username = var.username_system_2
  #   management_password = var.password_system_2
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # Port of the Peer System Default: 7611
  #port = 7611
  # Sets the Performance Profile, Options (Compact, HighPerformance) Default: HighPerformance
  #perf_profile = "HighPerformance"
}

resource "powerflex_peer_system" "system_2" {
  provider = powerflex.system_2
  ### Required Values

  # New name of the Peer System
  name = var.name
  # Peer System (System 1) ID
  peer_system_id = data.powerflex_protection_domain.protection_domain_system_1.protection_domains[0].system_id
  # List of Peer MDM Ips at the destination
  ip_list = var.mdm_ips_system_1

  ### Optional with defaults if unset


  # Add certificate flag, default: false. 
  # If true source_primary_mdm_information and destination_primary_mdm_information must be filled out in order to get and set the certificate
  # add_certificate = true

  # source_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_2
  #   management_username = var.username_system_2
  #   management_password = var.password_system_2
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # destination_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_1
  #   management_username = var.username_system_1
  #   management_password = var.password_system_1
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # Port of the Peer System Default: 7611
  #port = 7611
  # Sets the Performance Profile, Options (Compact, HighPerformance) Default: HighPerformance
  #perf_profile = "HighPerformance"
}

