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

# Example for uploading compliance file. After successful execution, compliance file will be uploaded to the manager.
resource "powerflex_firmware_repository" "upload-test" {
  source_location = "https://10.10.10.1/artifactory/Denver/RCMs/SoftwareOnly/PowerFlex_Software_4.5.0.0_287_r1.zip"
  username = "user" # To be provided in case of CIFS share
  password = "password" # To be provided in case of CIFS share
  approve = true # To be used to approve the unsigned file
  timeout = 45 #controls that till what time the upload compliance will run
}
