# /*
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import firmware respository by it's id
terraform import powerflex_firmware_repository.fr_import_by_id "<id>"

# After Import, username and password is not needed for approving the unsigned file in case of CIFS share. For approving the file in case of import, please refer the below config(change the value as per your use-case):
resource "powerflex_firmware_repository" "fr_import_by_id" {
  source_location = "https://10.10.10.1/artifactory/Denver/RCMs/SoftwareOnly/PowerFlex_Software_4.5.0.0_287_r1.zip"
  approve = true
  timeout = 45
}
