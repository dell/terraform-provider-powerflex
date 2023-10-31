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

# Get all fault set details present on the cluster
data "powerflex_fault_set" "example1" {
}

# Get fault set details using fault set IDs
data "powerflex_fault_set" "example2" {
  fault_set_ids = ["FaultSet_ID1", "FaultSet_ID2"]
}

# Get fault set details using fault set names
data "powerflex_fault_set" "example3" {
  fault_set_names = ["FaultSet_Name1", "FaultSet_Name2"]
}

output "fault_set_result" {
  value = data.powerflex_fault_set.example1.fault_set_details
}
