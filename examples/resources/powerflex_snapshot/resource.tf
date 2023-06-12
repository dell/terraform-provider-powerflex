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
# To import , check snapshot_resource_import.tf for more info
# To create / update, either volume_id or volume_name must be provided
# name is the required parameter to create or update
# other  atrributes like : access_mode, size, capacity_unit, lock_auto_snapshot, desired_retention, retention_unit, remove_mode are optional 
# To check which attributes of the snapshot can be updated, please refer Product Guide in the documentation

resource "powerflex_snapshot" "snapshots-create" {
  name      = "snapshots-create"
  volume_id = "4577c84000000120"
}

resource "powerflex_snapshot" "snapshots-create-01" {
  name          = "snapshots-create-epsilon"
  volume_id     = "4577c84000000120"
  access_mode   = "ReadWrite"
  size          = 16
  capacity_unit = "GB"
  remove_mode   = "INCLUDING_DESCENDANTS"
}


# General guidlines for furnishing this resource block 
# resource "powerflex_snapshot" "snapshots-create-1" {
# 	name = "<snapshot name>"
# 	volume_name = "<volume name>"
# 	access_mode = "<access mode options are ReadOnly/ReadWrite, default value ReadOnly>"
# 	size = "<size[int] associated with capacity unit>"
# 	capacity_unit =  "<capacity unit options are gb/tb, default value gb>"
# 	lock_auto_snapshot = "<lock auto snapshot, snapshot which are created by snapshot policy can be locked.>"
# 	desired_retention = "<desired retention[int] associated with retention unit>"
# 	retention_unit = "<retention unit options are hours/days, default value hours>"
# 	remove_mode = "<remove mode options are ONLY_ME/INCLUDING_DESCENDANTS, default value ONLY_ME>"
# }