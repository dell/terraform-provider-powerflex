---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerflex_os_repository resource"
linkTitle: "powerflex_os_repository"
page_title: "powerflex_os_repository Resource - powerflex"
subcategory: "Firmware and OS Management"
description: |-
  This resource is used to manage the OS Repository entity of the PowerFlex Array. We can Create ,Read and Delete the os image repository using this resource.
---

# powerflex_os_repository (Resource)

This resource is used to manage the OS Repository entity of the PowerFlex Array. We can Create ,Read and Delete the os image repository using this resource.

## Example Usage

```terraform
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

# Command to run this tf file : terraform plan && terraform apply.

// Resource to manage lifecycle
resource "terraform_data" "always_run" {
  input = timestamp()
}

# Example for creating OS repository. After successful execution, OS Repository will be created.
resource "powerflex_os_repository" "test" {
  # Required Fields

  # Name of the OS repository
  name = "TestOsRepo"
  # Source path of the OS image
  source_path = "https://pathtoimage.iso"
  # Supported image types are redhat7, vmware_esxi, sles, windows2016, windows2019
  image_type = "vmware_esxi"

  // This will allow terraform create process to trigger each time we run terraform apply.
  lifecycle {
    replace_triggered_by = [
      terraform_data.always_run
    ]
  }
}
```


After the execution of above resource block, OS Repository would have been created. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `image_type` (String) Type of the OS image. Supported types are redhat7, vmware_esxi, sles, windows2016, windows2019
- `name` (String) Name of the OS repository
- `source_path` (String) Source path of the OS image

### Optional

- `password` (String) Password of the OS repository
- `repo_type` (String) Type of the OS repository. Default is ISO
- `timeout` (Number) Describes the time in minutes to timeout the job.
- `username` (String) Username of the OS repository

### Read-Only

- `base_url` (String) Base URL of the OS repository
- `created_by` (String) User who created the OS repository
- `created_date` (String) Date of creation of the OS Repository
- `from_web` (Boolean) Whether the OS repository is from the web or not
- `id` (String) ID of the OS Repository
- `in_use` (Boolean) Whether the OS repository is in use or not
- `razor_name` (String) Name of the Razor
- `rcm_path` (String) Path of the RCM
- `state` (String) State of the OS repository
