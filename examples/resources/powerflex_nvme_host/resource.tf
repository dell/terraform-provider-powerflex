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
# Create, Update, Read, and Delete operations are supported for this resource.
# To import , check import.sh for more info
# nqn is the required parameter to create
# To check which attributes can be updated, please refer Product Guide in the documentation
# Please note that 
# 1. NVMe over TCP is supported in PowerFlex 4.0 and later versions, therefore this resource is not supported in PowerFlex 3.x.
# 2. Due to certain limitations, updating the NVMe host in PowerFlex versions earlier than 4.6 is not supported

# Example for adding NVMe host.
resource "powerflex_nvme_host" "test-nvme-host" {
  name              = "nvme_host_name"
  nqn               = "nqn.2014-08.org.nvmexpress:uuid:a10e4d56-a2c0-4cab-9a0a-9a7a4ebb8c0e"
  max_num_paths     = 4
  max_num_sys_ports = 10
}