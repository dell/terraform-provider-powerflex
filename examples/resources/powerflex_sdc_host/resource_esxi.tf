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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
# sdc_details is the required parameter for the SDC resource.

# Example for adding an ESXi host as SDC.
# In this example, we are using passwordless ssh authentication using private key and host key.

# load the private key
data "local_sensitive_file" "ssh_key" {
  filename = "/root/.ssh/esxi_rsa"
}

# load the host key
data "local_sensitive_file" "host_key" {
  filename = "esxi_host_ecdsa_key.pub"
}

# generate a random guid. This is required only for ESXi hosts.
resource "random_uuid" "sdc_guid" {
}

resource "powerflex_sdc_host" "sdc" {
  ip = "10.10.10.10"
  remote = {
    user = "root"
    # we are not using password auth here, but it can be used as well
    # password = "W0uldntUWannaKn0w!"
    private_key = data.local_sensitive_file.ssh_key.content_base64
    host_key    = data.local_sensitive_file.host_key.content_base64
  }
  os_family = "esxi"
  esxi = {
    guid         = random_uuid.sdc_guid.result
    drv_cfg_path = "/root/terraform-provider-powerflex/drv_cfg-3.6.500.106-esx7.x"
  }
  name         = "sdc-esxi"
  package_path = "/root/terraform-provider-powerflex/sdc-3.6.500.106-esx7.x.zip"
  mdm_ips      = ["10.10.10.5", "10.10.10.6"]
}

