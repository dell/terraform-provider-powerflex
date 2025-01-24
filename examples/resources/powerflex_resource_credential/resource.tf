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

terraform {
  required_providers {
    powerflex = {
      source  = "registry.terraform.io/dell/powerflex"
    }
    local = {
      source = "hashicorp/local"
      version = "2.5.2"
    }
  }
}
provider "powerflex" {
  username =  "user"
  password = "password"
  endpoint = "https://example.com"
  insecure = true
  timeout  = 120
}

# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Read, Delete and Import operations are supported for this resource

## Optional only needed if using ssh_private_key field
# data "local_file" "ssh_key" {
#   filename = var.ssh_private_key_path
# }

# Create Resource Credential
resource "powerflex_resource_credential" "example" {
 ## Required values for all credential types
 name = var.name
 type = var.type // Options: Node, Switch, vCenter, ElementManager, PowerflexGateway, PresentationServer, OSAdmin, OSUser
 password = var.password
 username = var.username

 ## Required value for vCenter, ElementManager, OSUser
 #domain = var.domain

 ## Required value for PowerflexGateway
 #os_username = var.os_username
 #os_password = var.os_password

 ## Optional values for Node, Switch, ElementManager
 #snmp_v2_community_string = var.snmpv2_community_string

 ## Optional for Node
 #snmp_v3_security_level = var.snmpv3_security_level // Options "Minimal", "Moderate", or "Maximal"
 #snmp_v3_security_name = var.snmpv3_security_name
 #snmp_v3_md5_authentication_password = var.snmpv3_md5_auth_password // required for level "Moderate" and "Maximal"
 #snmp_v3_des_authentication_password = var.snmpv3_des_private_password // required for level "Maximal"

 ## Optional for Node, Switch, OSAdmin, OSUser
 #ssh_private_key = data.local_file.ssh_key.content
 #key_pair_name = var.key_pair_name
}