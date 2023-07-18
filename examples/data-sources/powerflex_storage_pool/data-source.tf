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

# commands to run this tf file : terraform init && terraform apply --auto-approve
# To read storage pool, either protection_domain_name or protection_domain_id must be provided
# This datasource reads a list of storage pools either by storage_pool_ids or storage_pool_names where user can provide a list of ids or names
# if both storage_pool_ids and storage_pool_names are not provided , then it will read all the storage pool under the protection domain
# Both storage_pool_ids and storage_pool_names can't be provided together .
# Both protection_domain_name and protection_domain_id can't be provided together


data "powerflex_storage_pool" "example" {
  //protection_domain_name = "domain1"
  protection_domain_id = "202a046600000000"
  //storage_pool_ids = ["c98ec35000000002", "c98e26e500000000"]
  storage_pool_names = ["pool2", "pool1"]
}


output "allsdcresult" {
  value = data.powerflex_storage_pool.example.storage_pools
}
