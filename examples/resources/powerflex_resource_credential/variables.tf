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

variable "username" {
  type        = string
  description = "username of the credential resource."
  default = ""
}

variable "password" {
  type        = string
  description = "password of the credential resource."
  default = ""
}

variable "name" {
  type        = string
  description = "name of the credential resource."
  default = ""
}

variable "type" {
  type        = string
  description = "type of the credential resource."
  default = ""
}

variable "os_username" {
  type        = string
  description = "os_username of the credential resource."
  default = ""
}

variable "os_password" {
  type        = string
  description = "os_password of the credential resource."
  default = ""
}

variable "snmpv2_community_string" {
  type        = string
  description = "snmpv2_community_string of the credential resource."
  default = ""
}

variable "snmpv3_security_name" {
  type        = string
  description = "snmpv3_security_name of the credential resource."
  default = ""
}

variable "snmpv3_security_level" {
  type        = string
  description = "snmpv3_security_level of the credential resource."
  default = ""
}

variable "snmpv3_md5_auth_password" {
  type        = string
  description = "snmpv3_md5_auth_password of the credential resource."
  default = ""
}

variable "snmpv3_des_private_password" {
  type        = string
  description = "snmpv3_des_private_password of the credential resource."
  default = ""  
}

variable "ssh_private_key_path" {
  type        = string
  description = "ssh_private_key_path of the credential resource."
  default = ""
}

variable "key_pair_name" {
  type        = string
  description = "key_pair_name of the credential resource."
  default = ""
}
