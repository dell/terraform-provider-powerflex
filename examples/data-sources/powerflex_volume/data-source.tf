/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// Empty filter block will return all the volumes
data "powerflex_volume" "volume" {

}

output "volumeResult" {
  value = data.powerflex_volume.volume.volumes
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_volume" "volume_filter" {
  filter {
  #   id                                     = ["id1", "id2"]
  #   name                                   = ["name1", "name2"]
  #   creation_time                          = [1, 2]
  #   size_in_kb                             = [123, 456]
  #   ancestor_volume_id                     = ["ancestor_volume_id1", "ancestor_volume_id2"]
  #   vtree_id                               = ["vtree_id1", "vtree_id2"]
  #   consistency_group_id                   = ["consistency_group_id1", "consistency_group_id2"]
  #   volume_type                            = ["volume_type1", "volume_type2"]
  #   use_rm_cache                           = false
  #   storage_pool_id                        = ["storage_pool_id1", "storage_pool_id2"]
  #   data_layout                            = ["data_layout1", "data_layout2"]
  #   not_genuine_snapshot                   = false
  #   access_mode_limit                      = ["access_mode_limit1", "access_mode_limit2"]
  #   secure_snapshot_exp_time               = [789, 123]
  #   managed_by                             = ["managed_by1", "managed_by2"]
  #   locked_auto_snapshot                    = false
  #   locked_auto_snapshot_marked_for_removal = false
  #   compression_method                     = ["compression_method1", "compression_method2"]
  #   time_stamp_is_accurate                 = true
  #   original_expiry_time                   = [43, 71]
  #   volume_replication_state               = ["volume_replication_state1", "volume_replication_state2"]
  #   replication_journal_volume             = false
  #   replication_time_stamp                 = [1,2]
  }
}


output "volumeFilterResult" {
  value = data.powerflex_volume.volume_filter.volumes
}
