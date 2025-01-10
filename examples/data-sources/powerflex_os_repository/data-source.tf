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

# commands to run this tf file : terraform init && terraform apply --auto-approve

# Get all OS Repositories
data "powerflex_os_repository" "example1" {
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_os_repository" "os_repository_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     source_path = ["^c://.*$"]
#   }
# }

# output "osRepositoryFilterRegexResult"{
#  value = data.powerflex_os_repository.os_repository_filter_regex.os_repositories
# }

# # Get OS Repository details by ID
# data "powerflex_os_repository" "example2" {
#   # this datasource supports filters like os repository id, name, source path, etc.
#   # Note: If both filters are used simultaneously, the results will include any records that match either of the filters.
#   filter {
#     id = ["1234","5678"]
#   }
# }

# # Get OS Repository details by all fields
# data "powerflex_os_repository" "example3" {
#   # this datasource supports filters like os repository id, name, source path, etc.
#   # Note: If both filters are used simultaneously, the results will include any records that match either of the filters.
#   filter {
#     base_url = ["http://195.0.0.0:8080"]
#     created_by = ["system"]
#     created_date = ["20XX-XX-XXT13:58:13.978+00:00"]
#     from_web = false
#     id = ["1234","5678"]
#     image_type = ["vmwar_esxi"]
#     in_use = false
#     name   = ["Dell EMC PowerFlex Embedded OS"]
#     razor_name = ["DellEMCPowerFlexEmbeddedOS"]
#     rcm_path = ["c:\\..."]
#     repo_type  = ["ISO"]
#     source_path = ["z:\\ESXi\\ESXi-x.x.x-xxxxxx-x_Dell.zip"]
#     state = ["available"]
#     username = ["User"]
#   }
# }

output "os_repository_result" {
  value = data.powerflex_os_repository.example1.os_repositories
}
