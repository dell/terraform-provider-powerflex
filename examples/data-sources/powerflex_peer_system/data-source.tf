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
# This feature is only supported for PowerFlex 4.5 and above.

# Get all Peer System details present on the cluster
data "powerflex_peer_system" "example1" {
}

# Get Peer System details using filter with all values
# If there is no intersection between the filters then an empty datasource will be returned
# For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples
# data "powerflex_peer_system" "example2" {
#   filter {
#     id = ["b27112300000000"]
#     coupling_rc = ["SUCCESS"]
#     membership_state = ["Joined"]
#     name = ["10654"]
#     peer_system_id = ["89980fe50ff2243f"]
#     perf_profile = ["HighPerformance"]
#     network_type = ["External"]
#     port = [9091]
#     software_version_info = ["R4_6.2134.0"]
#     system_id = ["7d2f6023117d93f0"]
#   }
# }

output "peer_system_result" {
  value = data.powerflex_peer_system.example1
}