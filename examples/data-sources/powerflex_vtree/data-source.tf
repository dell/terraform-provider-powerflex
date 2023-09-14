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

# Get all VTrees details present on the cluster
data "powerflex_vtree" "example1" {
}

# Get VTree details using VTree IDs
data "powerflex_vtree" "example2" {
  vtree_ids = ["VTree_ID1", "VTree_ID2"]
}

# Get VTree details using Volume IDs
data "powerflex_vtree" "example3" {
  volume_ids = ["Volume_ID1", "Volume_ID2"]
}

# Get VTree details using Volume Names
data "powerflex_vtree" "example4" {
  volume_names = ["Volume_Name1", "Volume_Name2"]
}