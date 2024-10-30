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

# System 1 Vars
variable "username_system_1" {
  type        = string
  description = "Stores the system 1 username of PowerFlex host."
}

variable "password_system_1" {
  type        = string
  description = "Stores system 1 the password of PowerFlex host."
}

variable "endpoint_system_1" {
  type        = string
  description = "Stores system 1 the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}

variable "name" {
  type        = string
  description = "The peer system name."
}

variable "mdm_ips_system_1" {
  type        = set(string)
  description = "The system 1 list of mdm ips"
}

variable "protection_domain_name_system_1" {
  type        = string
  description = "The system 1 protection domain name."
}

# System 2 Vars
variable "username_system_2" {
  type        = string
  description = "Stores the system 2 username of PowerFlex host."
}

variable "password_system_2" {
  type        = string
  description = "Stores system 2 the password of PowerFlex host."
}

variable "endpoint_system_2" {
  type        = string
  description = "Stores system 2 the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}

variable "mdm_ips_system_2" {
  type        = set(string)
  description = "The system 2 list of mdm ips"
}

variable "protection_domain_name_system_2" {
  type        = string
  description = "The system 2 protection domain name."
}
