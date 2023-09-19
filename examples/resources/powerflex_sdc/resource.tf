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

# To perform Multiple SDC Detail Update and SDC Expansion
resource "powerflex_sdc" "expansion" {
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
      name                = "SDC_NAME"
      performance_profile = "HighPerformance"
      sdc_id              = "sdc_id"
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
      sdc_id              = "sdc_id"
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
      sdc_id              = "sdc_id"
    },
    {
      ip                  = "IP"
      username            = "Username"
      password            = "Password"
      operating_system    = "linux"
      is_mdm_or_tb        = "Standby"
      is_sdc              = "Yes"
      name                = "SDC_NAME"
      performance_profile = "Compact"
      sdc_id              = "sdc_id"
    },
  ]
}
