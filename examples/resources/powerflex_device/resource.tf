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
# To add device, device_path is mandatory along with storage_pool_name/storage_pool_id and sds_name/sds_id.
# Along with storage_pool_name, we have to specify protection_domain_id or protection_domain_name.
# To check which attributes of the device resource can be updated, please refer Product Guide in the documentation

resource "powerflex_device" "test-device" {
  device_path            = "/dev/sdc"
  storage_pool_name      = "pool1"
  protection_domain_name = "domain1"
  sds_name               = "SDS_2"
  media_type             = "HDD"
}
