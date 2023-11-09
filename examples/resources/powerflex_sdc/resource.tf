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
