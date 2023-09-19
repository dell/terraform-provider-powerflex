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
# Create, Update, Delete is supported for this resource.
# To import, check import.sh for more info.
# To create/update, either SDC ID or SDC name must be provided.
# volume_list attribute is optional. 
# To check which attributes of the sdc_volumes_mappping resource can be updated, please refer Product Guide in the documentation

resource "powerflex_sdc_volumes_mapping" "mapping-test" {
  # SDC id
  id = "e3ce1fb600000001"
  volume_list = [
    {
      # id of the volume which needs to be mapped. 
      # either volume_id or volume_name can be used.
      volume_id = "edb2059700000002"

      # Valid values are 0 or integers greater than 10
      limit_iops = 140

      # Default value is 0
      limit_bw_in_mbps = 19

      access_mode = "ReadOnly" # ReadOnly/ReadWrite/NoAccess
    },
    {
      volume_name      = "terraform-vol"
      access_mode      = "ReadWrite"
      limit_iops       = 120
      limit_bw_in_mbps = 25
    }
  ]
}

# To unmap all the volumes mapped to SDC, below config can be used. 

# resource "powerflex_sdc_volumes_mapping" "mapping-test" {
#   id          = "e3ce1fb600000001"
#   volume_list = []
# }
