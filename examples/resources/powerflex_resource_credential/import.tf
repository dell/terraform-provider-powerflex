/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// Get the data of a specific resource credential
data "powerflex_resource_credential" "all_current_resource_credentials" {
    filter {
        label = [var.name]
    }
}

import {
  to = powerflex_resource_credential.imported_resource_credentials
  id = data.powerflex_resource_credential.all_current_resource_credentials.resource_credential_details[0].id
}

resource "powerflex_resource_credential" "imported_resource_credentials" {
  name = data.powerflex_resource_credential.all_current_resource_credentials.resource_credential_details[0].label
  username = data.powerflex_resource_credential.all_current_resource_credentials.resource_credential_details[0].username
  password = var.password
  type = var.type
}
