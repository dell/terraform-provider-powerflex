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

# Commands to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check sds_resource_import.tf for more info
# To create / update, either protection_domain_name or protection_domain_id must be provided
# name and ip_list are the required parameters to create or update
# other  atrributes like : performance_profile, port, drl_mode, rmcache_enabled, rfcache_enabled, rmcache_size_in_mb are optional 
# To check which attributes can be updated, please refer Product Guide in the documentation

# Example for adding SDS. After successful execution, SDS will be added to the protection domain.
resource "powerflex_sds" "create" {
  name                   = "demo-sds-test-01"
  protection_domain_name = "demo-sds-pd"
  ip_list = [
    {
      ip   = "10.10.10.12"
      role = "sdsOnly" # all/sdsOnly/sdcOnly
    },
    {
      ip   = "10.10.10.11"
      role = "sdcOnly" # all/sdsOnly/sdcOnly
    },
  ]
}

output "changed_sds" {
  value = powerflex_sds.create
}
