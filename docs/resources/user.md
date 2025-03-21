---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerflex_user resource"
linkTitle: "powerflex_user"
page_title: "powerflex_user Resource - powerflex"
subcategory: "User Management"
description: |-
  This resource is used to manage the User entity of the PowerFlex Array. We can Create, Update and Delete the user using this resource. We can also import an existing user from the PowerFlex array.
---

# powerflex_user (Resource)

This resource is used to manage the User entity of the PowerFlex Array. We can Create, Update and Delete the user using this resource. We can also import an existing user from the PowerFlex array.

## Example Usage

```terraform
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
```

After the execution of above resource block, new user would have been created on the PowerFlex array. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the user. For PowerFlex version 3.6, cannot be updated.
- `password` (String) Password of the user. For PowerFlex version 3.6, cannot be updated.
- `role` (String) The role of the user. Accepted values for PowerFlex version 3.6 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'. Accepted values for PowerFlex version 4.5 are 'Monitor', 'SuperUser', 'SystemAdmin', 'StorageAdmin', 'LifecycleAdmin', 'ReplicationManager', 'SnapshotManager', 'SecurityAdmin', 'DriveReplacer', 'Technician', 'Support'.

### Optional

- `first_name` (String) First name of the user. PowerFlex version 3.6 does not support the first_name attribute. It is mandatory for PowerFlex version 4.6.
- `last_name` (String) Last name of the user. PowerFlex version 3.6 does not support the last_name attribute. It is mandatory for PowerFlex version 4.6.

### Read-Only

- `id` (String) The ID of the user.
- `system_id` (String) The ID of the system.

## Import

Import is supported using the following syntax:

```shell
# /*
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

# import user by it's id
terraform import powerflex_user.user_import_by_id "<id>"

# import user by it's id - alternative approach by prefixing it with "id:"
terraform import powerflex_user.user_import_by_id "<id:id_of_the_user>"

# import user by it's name
terraform import powerflex_user.user_import_by_name "<name:name_of_the_user>"
```

1. This will import the User instance with specified ID into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.