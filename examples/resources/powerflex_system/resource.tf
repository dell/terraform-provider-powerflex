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

# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import, check import.sh for more info
# restricted_mode is the required parameter

resource "powerflex_system" "test" {
  restricted_mode = "Guid"
  sdc_approved_ips = [
    {
      id = "sdc_id1"
      ips = ["sdc_ip1", "sdc_ip2"]
    },
    {
      id = "sdc_id2"
      ips = ["sdc_ip3"]
    },
  ]
}