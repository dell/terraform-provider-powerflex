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

# Source Vars
variable "username_source" {
  type        = string
  description = "Stores the source username of PowerFlex host."
}

variable "password_source" {
  type        = string
  description = "Stores source the password of PowerFlex host."
}

variable "endpoint_source" {
  type        = string
  description = "Stores source the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}

variable "name" {
  type        = string
  description = "The replication consistancy group name."
}

variable "rpo_in_seconds" {
  type        = number
  description = "The replication consistancy group rpo_in_seconds."
}

variable "source_protection_domain_name" {
  type        = string
  description = "The source protection domain name."
}

# Destination Vars
variable "username_destination" {
  type        = string
  description = "Stores the destination username of PowerFlex host."
}

variable "password_destination" {
  type        = string
  description = "Stores destination the password of PowerFlex host."
}

variable "endpoint_destination" {
  type        = string
  description = "Stores destination the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}

variable "destination_protection_domain_name" {
  type        = string
  description = "The destination protection domain name."
}