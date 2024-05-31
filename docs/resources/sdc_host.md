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

title: "powerflex_sdc_host resource"
linkTitle: "powerflex_sdc_host"
page_title: "powerflex_sdc_host Resource - powerflex"
subcategory: ""
description: |-
  This resource is used to manage the SDC entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.
---

# powerflex_sdc_host (Resource)

This resource is used to manage the SDC entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.

## Example Usage

### With ESXi

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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
# sdc_details is the required parameter for the SDC resource.

# Example for adding an ESXi host as SDC.
# In this example, we are using passwordless ssh authentication using private key and host key.

# load the private key
data "local_sensitive_file" "ssh_key" {
  filename = "/root/.ssh/esxi_rsa"
}

# load the host key
data "local_sensitive_file" "host_key" {
  filename = "esxi_host_ecdsa_key.pub"
}

# generate a random guid. This is required only for ESXi hosts.
resource "random_uuid" "sdc_guid" {
}

resource "powerflex_sdc_host" "sdc" {
  ip = "10.10.10.10"
  remote = {
    user = "root"
    # we are not using password auth here, but it can be used as well
    # password = "W0uldntUWannaKn0w!"
    private_key = data.local_sensitive_file.ssh_key.content_base64
    host_key    = data.local_sensitive_file.host_key.content_base64
  }
  os_family = "esxi"
  esxi = {
    guid         = random_uuid.sdc_guid.result
    drv_cfg_path = "/root/terraform-provider-powerflex/drv_cfg-3.6.500.106-esx7.x"
  }
  name         = "sdc-esxi"
  package_path = "/root/terraform-provider-powerflex/sdc-3.6.500.106-esx7.x.zip"
  mdm_ips      = ["10.10.10.5", "10.10.10.6"]
}
```

After the execution of above resource block, the ESXi host would have been addes as an SDC to the PowerFlex array. For more information, please check the terraform state file.

### With Linux

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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
# sdc_details is the required parameter for the SDC resource.

# Example for adding an Linux host as SDC.
# In this example, we are using passwordless ssh authentication using private key and host key.

# load the private key
data "local_sensitive_file" "ssh_key" {
  filename = "/root/.ssh/linux_rsa"
}

# load the host key
data "local_sensitive_file" "host_key" {
  filename = "linux_host_ecdsa_key.pub"
}

# Example for adding an Linux host as SDC.
resource "powerflex_sdc_host" "sdc_linux" {
  ip = "10.10.10.10"
  remote = {
    user = "root"
    # we are not using password auth here, but it can be used as well
    # password = "password"
    private_key = data.local_sensitive_file.ssh_key.content_base64
    host_key    = data.local_sensitive_file.host_key.content_base64
  }
  os_family    = "linux"
  name         = "sdc-linux"
  package_path = "/root/terraform-provider-powerflex/EMC-ScaleIO-sdc-3.6-700.103.Ubuntu.22.04.x86_64.tar" # For Ubuntu
  # package_path = "/root/terraform-provider-powerflex/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm" # For RHEL
  # mdm_ips = ["10.10.10.5", "10.10.10.6"]   # Optional 
}
```

After the execution of above resource block, the Linux host would have been addes as an SDC to the PowerFlex array. For more information, please check the terraform state file.

### With Windows

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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
# sdc_details is the required parameter for the SDC resource.

# Example for adding an Windows host as SDC.
# In this example, we are using passwordless ssh authentication using private key and host key.


# Example for adding an Windows host as SDC.
resource "powerflex_sdc_host" "sdc_windows" {
  ip = "10.10.10.10"
  remote = {
    user     = "username"
    password = "password"
  }
  os_family    = "windows"
  name         = "sdc-windows"
  package_path = "/root/terraform-provider-powerflex/EMC-ScaleIO-sdc-3.6-200.105.msi"
  # mdm_ips = ["10.10.10.5", "10.10.10.6"]   # Optional 
}
```

After the execution of above resource block, the Windows Server host would have been addes as an SDC to the PowerFlex array. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ip` (String) IP address of the server to be configured as SDC.
- `os_family` (String) Operating System family of the SDC. Accepted values are 'linux', 'windows' and 'esxi'. Cannot be changed once set.
- `package_path` (String) Full path (on local machine) of the package to be installed on the SDC.
- `remote` (Attributes) Remote login details of the SDC. (see [below for nested schema](#nestedatt--remote))

### Optional

- `esxi` (Attributes) Details of the SDC host if the `os_family` is `esxi`. (see [below for nested schema](#nestedatt--esxi))
- `mdm_ips` (List of String) List of MDM IPs to be assigned to the SDC.
- `name` (String) Name of SDC.
- `performance_profile` (String) Performance profile of the SDC. Accepted values are 'HighPerformance' and 'Compact'.

### Read-Only

- `guid` (String) GUID of the HOST
- `id` (String) The id of the SDC
- `is_approved` (Boolean) Is Host Approved
- `mdm_connection_state` (String) MDM Connection State
- `on_vmware` (Boolean) Is Host on VMware
- `system_id` (String) System ID of the Host

<a id="nestedatt--remote"></a>
### Nested Schema for `remote`

Required:

- `user` (String) Remote Login username of the SDC server.

Optional:

- `certificate` (String) Remote Login certificate issued by a CA to the remote login user. Must be used with `private_key` and the private key must match the certificate.
- `dir` (String) Directory on the SDC server to upload packages to for Unix. Defaults to `/tmp` on Unix.
- `host_key` (String) Remote Login host key of the SDC server. Corresponds to the UserKnownHostsFile field of OpenSSH.
- `password` (String, Sensitive) Remote Login password of the SDC server.
- `port` (String) Remote Login port of the SDC server. Defaults to `22`.
- `private_key` (String) Remote Login private key of the SDC server. Corresponds to the IdentityFile field of OpenSSH.


<a id="nestedatt--esxi"></a>
### Nested Schema for `esxi`

Required:

- `guid` (String) GUID of the SDC.

Optional:

- `verify_vib_signature` (Boolean) Whether to verify the VIB signature or not. Defaults to `true`.

## Import

Import is supported using the following syntax:

```shell
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

# import SDC by it's IP address
terraform import powerflex_sdc_host.sdc "<ip>"
```

1. This will import the SDC instance with specified IP into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.