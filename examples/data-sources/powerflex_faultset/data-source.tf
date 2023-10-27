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

# Get all faultset details present on the cluster
data "powerflex_faultset" "example1" {
}

# Get faultset details using faultset IDs
data "powerflex_faultset" "example2" {
  faultset_ids = ["FaultSet_ID1", "FaultSet_ID2"]
}

# Get faultset details using faultset names
data "powerflex_faultset" "example3" {
  faultset_names = ["FaultSet_Name1", "FaultSet_Name2"]
}