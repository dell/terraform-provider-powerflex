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

# Example for adding an Windows host as SDC.
# In this example, we are using passwordless ssh authentication using private key and host key.


# Example for adding an Windows host as SDC.
resource "powerflex_sdc_host" "sdc_windows" {
  ip = "10.10.10.10"
  remote = {
    user     = "username"
    password = "password"
  }
  os_family    = "windows"
  name         = "sdc-windows"
  package_path = "/root/terraform-provider-powerflex/EMC-ScaleIO-sdc-3.6-200.105.msi"
  # mdm_ips = ["10.10.10.5", "10.10.10.6"]   # Optional 
}
