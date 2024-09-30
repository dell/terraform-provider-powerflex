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

# Get NVMe host details using NVMe host Names
data "powerflex_nvme_host" "example2" {
  filter {
    names = ["name1", "name2"]
  }
}


# Get NVMe host details using NVMe host IDs
data "powerflex_nvme_host" "example3" {
  filter {
    ids = ["ID1", "ID2"]
  }
}

# Get NVMe host details using NVMe host Nqns
data "powerflex_nvme_host" "example4" {
  filter {
    nqns = ["NQN1", "NQN2"]
  }
}

output "nvme_host_result" {
  value = data.powerflex_nvme_host.example1
}
