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
# To read SDS, either protection_domain_name or protection_domain_id must be provided
# This datasource reads a specific SDS either by sds_names or sds_ids where user can provide a list of sds ids or names
# if both sds_names and sds_ids are not provided , then it will read all the sds under the protection domain
# Both sds_ids and sds_names can't be provided together .
# Both protection_domain_name and protection_domain_id can't be provided together



data "powerflex_sds" "example2" {
  # require field is either of protection_domain_name or protection_domain_id
  protection_domain_name = "domain1"
  # protection_domain_id = "202a046600000000"
  sds_names = ["SDS_01_MOD", "sds_1", "node4"]
  # sds_ids = ["6adfec1000000000", "6ae14ba900000006", "6ad58bd200000002"]
}

output "allsdcresult" {
  value = data.powerflex_sds.example2
}

