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

# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check import.sh for more info
# name, size is the required parameter to create or update
# other  atrributes like : capacity_unit, volume_type, use_rm_cache, compression_method, access_mode, remove_mode are optional 
# To check which attributes of the volume can be updated, please refer Product Guide in the documentation


resource "powerflex_volume" "avengers-volume-create" {
  name = "avengers-volume-create"

  # To create / update, either protection_domain_id or protection_domain_name must be provided
  protection_domain_name = "domain1"

  # To create / update, either storage_pool_id or storage_pool_name must be provided
  storage_pool_name = "pool1"

  # The unit of size of the volume is defined by capacity_unit whose default value is "GB".
  size          = 8
  capacity_unit = "GB"

  use_rm_cache = true
  volume_type  = "ThickProvisioned"      # ThickProvisioned/ThinProvisioned volume type
  access_mode  = "ReadWrite"             # ReadWrite/ReadOnly volume access mode
  remove_mode  = "INCLUDING_DESCENDANTS" # INCLUDING_DESCENDANTS/ONLY_ME remove mode
}