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
    port     = 5985
  }
  os_family    = "windows"
  name         = "sdc-windows"
  package_path = "/root/terraform-provider-powerflex/EMC-ScaleIO-sdc-3.6-200.105.msi"
  # Optional all the mdms(either primary,secondary or virtual ips) in a comma separated list by cluster if unset will use the mdms of the cluster set in the provider block
  # Removal of mdms is not supported for windows, if you wish to remove a cluster from the sdc please follow steps here: https://www.dell.com/support/kbdoc/en-us/000167031/how-do-i-remove-the-mdm-entry-from-the-sdc-as-displayed-in-the-output-of-drv-cfg-binary-in-query-mdms-on-the-sdc-on-windows-or-linux-os#:~:text=Resolution%201%20For%20Linux%20SDC%20host%2C%20open%20%2Fbin%2Femc%2Fscaleio%2Fdrv_cfg.txt,4%20Reboot%20Linux%20SDC%20host%20to%20apply%20changes.?msockid=0ee30a4c8e9f67f610c21ecc8f89664a
  # clusters_mdm_ips = ["10.10.10.5,10.10.10.6", "10.10.10.7,10.10.10.8"]   
}
