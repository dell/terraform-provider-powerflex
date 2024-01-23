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

# Get all node details present on the cluster
data "powerflex_node" "example1" {
}

# Get node details using node IDs
data "powerflex_node" "example2" {
  node_ids = ["Node_ID1", "Node_ID2"]
}

# Get node details using IP addresses
data "powerflex_node" "example3" {
  ip_addresses = ["IP1", "IP2"]
}

# Get node details using service tags
data "powerflex_node" "example4" {
  service_tags = ["Service_Tag1", "Service_Tag2"]
}

# Get node details using nodepool IDs
data "powerflex_node" "example5" {
  node_pool_ids = ["NodePool_ID1", "NodePool_ID2"]
}

# Get node details using nodepool names
data "powerflex_node" "example6" {
  node_pool_names = ["NodePool_Name1", "NodePool_Name2"]
}


output "node_result" {
  value = data.powerflex_node.example1.node_details
}
