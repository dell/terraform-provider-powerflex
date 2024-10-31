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
# To import , check import.sh for more info.
# name and ip_list attributes are required, and either protection_domain_name or protection_domain_id must be specified.
# To check which attributes can be updated, please refer Product Guide in the documentation
# To avoid potential issues, it is recommended to operate NVMe targets using the default ports.

# Example for adding NVMe target.
resource "powerflex_nvme_target" "test-nvme-target" {
  name                   = "nvme_target_name"
  protection_domain_name = "demo-pd"
  ip_list = [
    {
      ip   = "10.10.10.13"
      role = "StorageAndHost" #StorageAndHost, StorageOnly, HostOnly
    },
    {
      ip   = "10.10.10.14"
      role = "StorageAndHost" #StorageAndHost, StorageOnly, HostOnly
    }
  ]
}