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

# Get all NVMe host details present on the cluster
data "powerflex_nvme_host" "example1" {
}

# Get NVMe host details by NVMe host Name
data "powerflex_nvme_host" "example2" {
  # this datasource supports filters like NVMe host name, id, nqn, max_num_paths, max_num_sys_ports, system_id.
  # Note: If both filters are used simultaneously, the results will include any records that match either of the filters.

  filter {
    name = ["name1", "name2"]
    # id = ["ID1", "ID2"]
    # nqn = ["NQN1", "NQN2"]
    # max_num_paths = [2]
    # max_num_sys_ports = [10]
    # system_id = ["systemID"]
  }
}

output "nvme_host_result" {
  value = data.powerflex_nvme_host.example1
}
