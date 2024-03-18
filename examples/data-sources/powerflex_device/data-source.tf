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

# commands to run this tf file : terraform init && terraform apply --auto-approve
# empty block of the powerflex_device datasource will give list of all device within the system

data "powerflex_device" "dev" {
}

data "powerflex_device" "dev1" {
  storage_pool_id = "c98e26e500000000"
}

data "powerflex_device" "dev2" {
  sds_name = "SDS_2"
}

output "deviceResult" {
  value = data.powerflex_device.dev.device_model
}
