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

# Get OS Repository details by ID
data "powerflex_os_repository" "example2" {
  # this datasource supports filters like os repsoitory ids, names. 
  # Note: If both filters are used simultaneously, the results will include any records that match either of the filters.
  filter {
    os_repo_ids = ["1234","5678"]
  }
}

# Get OS Repository details by Name
data "powerflex_os_repository" "example3" {
  # this datasource supports filters like os repsoitory ids, names.
  # Note: If both filters are used simultaneously, the results will include any records that match either of the filters.
  filter {
    os_repo_names = ["test"]
  }
}
output "os_repository_result" {
  value = data.powerflex_os_repository.example1
}
