/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# sdc_details is the required parameter for the SDC resource.

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

# Example demonstrating the use of the Package Resource and Virtual IP configuration 

resource "powerflex_package" "upload-test" {
  file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
  "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm"]
}

resource "powerflex_sdc" "test" {
  depends_on   = [powerflex_package.upload-test]
  mdm_password = "mdm_password"
  lia_password = "lia_password"
  sdc_details = [
    {
      ip               = "primary_mdm_ip"
      password         = "vm_password"
      operating_system = "linux"
      is_mdm_or_tb     = "Primary"
      is_sdc           = "Yes"
      virtual_ips      = "virtual_ip"
      virtual_ip_nics  = "virtual_nic"
      data_network_ip  = "data_network_ip_pmdm"
    },
    {
      ip               = "secondary_mdm_ip"
      password         = "vm_password"
      operating_system = "linux"
      is_mdm_or_tb     = "Secondary"
      is_sdc           = "Yes"
      virtual_ips      = "virtual_ip"
      virtual_ip_nics  = "virtual_nic"
      data_network_ip  = "data_network_ip_smdm"
    },
    {
      ip               = "tiebreaker_mdm"
      password         = "vm_password"
      operating_system = "linux"
      is_mdm_or_tb     = "TB"
      is_sdc           = "No"
      data_network_ip  = "data_network_ip_tmdm"
    }
  ]
}


# Example for deleting all MDMs installed as SDCs. After successful execution, SDCs will be removed from the cluster. 
resource "powerflex_sdc" "expansion" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details  = []
}

# Example for installing non-MDM node as SDC. After successful execution, one SDC will be added.
resource "powerflex_sdc" "sdc-example" {
  mdm_password = "Password"
  lia_password = "Password"
  sdc_details = [
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Primary"
      is_sdc           = "No"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Secondary"
      is_sdc           = "No"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "TB"
      is_sdc           = "No"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_sdc           = "Yes"
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
      name   = "rename_sdc"
    },
  ]
}


# To perform Multiple SDC Detail Update only
resource "powerflex_sdc" "sdc_update" {
  sdc_details = [
    {
      sdc_id              = "sdc_id"
      name                = "SDC_NAME"
      performance_profile = "HighPerformance"
    },
    {
      sdc_id              = "sdc_id"
      name                = "SDC_NAME"
      performance_profile = "HighPerformance"
    },
  ]
}
