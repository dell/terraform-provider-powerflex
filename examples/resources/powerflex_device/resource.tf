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

# Example for adding device. After successful execution, device will be added to the specified storage pool
resource "powerflex_device" "test-device" {
  device_path            = "/dev/sdc"
  storage_pool_name      = "pool1"
  protection_domain_name = "domain1"
  sds_name               = "SDS_2"
  media_type             = "HDD" # HDD/SSD
}
