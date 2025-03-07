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

title: "powerflex_package resource"
linkTitle: "powerflex_package"
page_title: "powerflex_package Resource - powerflex"
subcategory: "Firmware and OS Management"
description: |-
  This resource can be used to upload packages on a PowerFlex Gateway. We can add or remove packages from PowerFlex Gateway. Import functionality is not supported.
---

# powerflex_package (Resource)

This resource can be used to upload packages on a PowerFlex Gateway. We can add or remove packages from PowerFlex Gateway. Import functionality is not supported.

>**Note:** This resource can be used to upload packages to a PowerFlex Manager (4.x) or Gateway Server (3.x). We can add or remove packages from PowerFlex Manager or Gateway. Import functionality is not supported.

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

# Example for uploading packages. After successful execution, packages will be uploaded to the gateway.
resource "powerflex_package" "upload-test" {
  file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm",
  "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdt-3.6-700.103.el7.x86_64.rpm"]
}
```

After the execution of above resource block, package would have been uploaded on the PowerFlex Gateway. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `file_path` (List of String) The list of path of packages

### Read-Only

- `id` (String) The ID of the package.
- `package_details` (Attributes Set) Uploaded Packages details. (see [below for nested schema](#nestedatt--package_details))

<a id="nestedatt--package_details"></a>
### Nested Schema for `package_details`

Read-Only:

- `file_name` (String) The Name of package.
- `label` (String) Uploaded Package Minor Version with OS Combination.
- `latest` (Boolean) Package Version is latest or not
- `linux_flavour` (String) Type of Linux OS
- `operating_system` (String) Supported OS.
- `sio_patch_number` (Number) Package Patch Number.
- `size` (Number) Size of Package.
- `type` (String) Type of Package. Like. MDM, LIA, SDS, SDC, etc.
- `version` (String) Uploaded Package Version.

