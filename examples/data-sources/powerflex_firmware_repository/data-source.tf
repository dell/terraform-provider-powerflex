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

# Example for fetching details of the firmware repository using names
data "powerflex_firmware_repository" "all" {
  
}

output "powerflex_firmware_repository_all_result" {
  value = data.powerflex_firmware_repository.all.firmware_repository_details
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_firmware_repository" "firmware_repository_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     disk_location = ["^https://powerflex.*$"]
#   }
# }

# output "firmwareRepositoryFilterRegexResult"{
#  value = data.powerflex_firmware_repository.firmware_repository_filter_regex.firmware_repository_details
# }

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_firmware_repository" "filter" {
  filter {
      #id = ["ID1", "ID2"]
      # name = ["Name1", "Name2"]
      # source_location = ["SourceLocation1", "SourceLocation2"]
      # source_type = ["SourceType1", "SourceType2"]
      # disk_location = ["DiskLocation1", "DiskLocation2"]
      # filename = ["Filename1", "Filename2"]
      # username = ["Username1", "Username2"]
      # download_status = ["DownloadStatus1", "DownloadStatus2"]
      # created_date = ["CreatedDate1", "CreatedDate2"]
      # created_by = ["CreatedBy1", "CreatedBy2"]
      # updated_date = ["UpdatedDate1", "UpdatedDate2"]
      # updated_by = ["UpdatedBy1", "UpdatedBy2"]
      # default_catalog = false
      # embedded = false
      # state = ["state1", "state2"]
      # bundle_count = [10, 11]
      # component_count = [1, 2]
      # user_bundle_count = [1, 2]
      # minimal = false
      # download_progress = [1, 2]
      # extract_progress = [1, 2]
      # signature = ["Signature1", "Signature2"]
      # custom = false
      # needs_attention = false
      # job_id = ["JobID1", "JobID2"]
      # rcmapproved = false
  }
}

output "powerflex_firmware_repository_filtered_result" {
  value = data.powerflex_firmware_repository.filter.firmware_repository_details
}

