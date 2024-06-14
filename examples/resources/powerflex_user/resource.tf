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

# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# Import is supported.

# Example for creating user. After successful execution, user will be created with Monitor role. Below example works for PowerFlex version 3.6.
resource "powerflex_user" "newUser" {

  # name, role and password is the required parameter to create or update. Only role can be updated.
  name     = "terraform"
  role     = "Monitor" # Administrator/Monitor/Configure/Security/FrontendConfig/BackendConfig
  password = "Password"
}

# Example for creating user. After successful execution, user will be created with SystemAdmin role. Below example works for PowerFlex version 4.5.
resource "powerflex_user" "newUser" {

  # name, role and password is the required parameter to create or update. first_name and last_name are optional parameters in PowerFlex version 4.5.
  name     = "terraform"
  role     = "SystemAdmin" # Monitor/SuperUser/SystemAdmin/StorageAdmin/LifecycleAdmin/ReplicationManager/SnapshotManager/SecurityAdmin/DriveReplacer/Technician/Support
  password = "Password123@"
}

# Example for creating user. After successful execution, user will be created with LifecycleAdmin role. Below example works for PowerFlex version 4.6.
resource "powerflex_user" "test" {
  name       = "terraform"
  role       = "LifecycleAdmin" # Monitor/SuperUser/SystemAdmin/StorageAdmin/LifecycleAdmin/ReplicationManager/SnapshotManager/SecurityAdmin/DriveReplacer/Technician/Support
  password   = "Password123@"
  first_name = "terraform"
  last_name  = "terraform"
}
