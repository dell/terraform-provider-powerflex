title: "powerflex_peer_system resource"
linkTitle: "powerflex_peer_system"
page_title: "powerflex_peer_system Resource - powerflex"
subcategory: ""
description: |-
  This resource is used to manage the Peer System entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Peer Systems using this resource. We can also Import an existing Peer Systems from the PowerFlex array. Peer system refers to the setup where multiple MDM nodes work together as peers to provide redundancy and high availability. This means that if one MDM node fails, other peer MDM nodes can take over its responsibilities, ensuring continuous operation without disruptions
---

# powerflex_peer_system (Resource)

This resource is used to manage the Peer System entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Peer Systems using this resource. We can also Import an existing Peer Systems from the PowerFlex array. Peer system refers to the setup where multiple MDM nodes work together as peers to provide redundancy and high availability. This means that if one MDM node fails, other peer MDM nodes can take over its responsibilities, ensuring continuous operation without disruptions

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

terraform {
  required_providers {
    powerflex = {
      source                = "registry.terraform.io/dell/powerflex"
      configuration_aliases = [powerflex.system_1, powerflex.system_2]
    }
  }
}

provider "powerflex" {
  alias    = "system_1"
  username = var.username_system_1
  password = var.password_system_1
  endpoint = "https://${var.endpoint_system_1}"
  insecure = true
  timeout  = 120
}

provider "powerflex" {
  alias    = "system_2"
  username = var.username_system_2
  password = var.password_system_2
  endpoint = "https://${var.endpoint_system_2}"
  insecure = true
  timeout  = 120
}

data "powerflex_protection_domain" "protection_domain_system_2" {
  provider = powerflex.system_2
  name     = var.protection_domain_name_system_2
}

data "powerflex_protection_domain" "protection_domain_system_1" {
  provider = powerflex.system_1
  name     = var.protection_domain_name_system_1
}


resource "powerflex_peer_system" "system_1" {
  provider = powerflex.system_1
  
  // This should be done in order to avoid a confict while sshing
  depends_on = [ resource.powerflex_peer_system.system_1 ]
  ### Required Values

  # New name of the Peer System
  name = var.name
  # Peer System (System 2) ID
  peer_system_id = data.powerflex_protection_domain.protection_domain_system_2.protection_domains[0].system_id
  # List of Peer MDM Ips at the destination
  ip_list = var.mdm_ips_system_2

  ### Optional with defaults if unset


  # Add certificate flag, default: false. 
  # If true source_primary_mdm_information and destination_primary_mdm_information must be filled out in order to get and set the certificate
  #add_certificate = true

  # source_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_1
  #   management_username = var.username_system_1
  #   management_password = var.password_system_1
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # destination_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_2
  #   management_username = var.username_system_2
  #   management_password = var.password_system_2
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # Port of the Peer System Default: 7611
  #port = 7611
  # Sets the Performance Profile, Options (Compact, HighPerformance) Default: HighPerformance
  #perf_profile = "HighPerformance"
}

resource "powerflex_peer_system" "system_2" {
  provider = powerflex.system_2
  ### Required Values

  # New name of the Peer System
  name = var.name
  # Peer System (System 1) ID
  peer_system_id = data.powerflex_protection_domain.protection_domain_system_1.protection_domains[0].system_id
  # List of Peer MDM Ips at the destination
  ip_list = var.mdm_ips_system_1

  ### Optional with defaults if unset


  # Add certificate flag, default: false. 
  # If true source_primary_mdm_information and destination_primary_mdm_information must be filled out in order to get and set the certificate
  # add_certificate = true

  # source_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_2
  #   management_username = var.username_system_2
  #   management_password = var.password_system_2
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # destination_primary_mdm_information = {
  #   # Required fields
  #   ip = "1.2.3.4"
  #   ssh_username = "user"
  #   ssh_password = "pass"
  #   management_ip = var.endpoint_system_1
  #   management_username = var.username_system_1
  #   management_password = var.password_system_1
  #   # Optional field defaults to 22
  #   #ssh_port = "22"
  # }

  # Port of the Peer System Default: 7611
  #port = 7611
  # Sets the Performance Profile, Options (Compact, HighPerformance) Default: HighPerformance
  #perf_profile = "HighPerformance"
}
```

After the execution of above resource block, the two PowerFlex Peer Systems will be connected. The user can then setup Replication Groups and Pairs in order to replicate volumes across powerflex systems.
For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ip_list` (Set of String) List of ips for the peer mdm system.
- `name` (String) Name of the peer mdm instance.
- `peer_system_id` (String) Unique identifier of the peer mdm system.

### Optional

- `add_certificate` (Boolean) Flag that if set to true will attempt to add certificate of the peer mdm destination to source. This flag is only used during create.
- `destination_primary_mdm_information` (Attributes) Only used if add_certificate is set to true during create. The destination primary mdm information to get the root certificate. (see [below for nested schema](#nestedatt--destination_primary_mdm_information))
- `perf_profile` (String) Performance profile of the peer mdm instance.
- `port` (Number) Port of the peer mdm instance.
- `source_primary_mdm_information` (Attributes) Only used if add_certificate is set to true during create. The source primary mdm information to get the root certificate. (see [below for nested schema](#nestedatt--source_primary_mdm_information))

### Read-Only

- `coupling_rc` (String) Coupling return code number of the peer mdm system.
- `id` (String) Unique identifier of the peer mdm instance.
- `membership_state` (String) membership state of the peer mdm instance.
- `network_type` (String) Network type of the peer mdm system.
- `software_version_info` (String) Software version details of the peer mdm instance.
- `system_id` (String) Unique identifier of the peer mdm system.

<a id="nestedatt--destination_primary_mdm_information"></a>
### Nested Schema for `destination_primary_mdm_information`

Optional:

- `ip` (String) ip of the primary destination mdm instance.
- `management_ip` (String) ip of the destination management instance.
- `management_password` (String, Sensitive) password of the management instance.
- `management_username` (String) ssh username of the destination management instance.
- `ssh_password` (String, Sensitive) ssh password of the primary destination mdm instance.
- `ssh_port` (String) port of the primary destination mdm instance.
- `ssh_username` (String) ssh username of the destination primary mdm instance.


<a id="nestedatt--source_primary_mdm_information"></a>
### Nested Schema for `source_primary_mdm_information`

Optional:

- `ip` (String) ip of the primary source mdm instance.
- `management_ip` (String) ip of the source management instance.
- `management_password` (String, Sensitive) password of the source instance.
- `management_username` (String) ssh username of the source management instance.
- `ssh_password` (String, Sensitive) ssh password of the source primary mdm instance.
- `ssh_port` (String) port of the primary source mdm instance.
- `ssh_username` (String) ssh username of the source primary mdm instance.

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

# import Peer System by its id
terraform import powerflex_peer_system.example "<id>"
```

1. This will import the resource instance with specified ID into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.
