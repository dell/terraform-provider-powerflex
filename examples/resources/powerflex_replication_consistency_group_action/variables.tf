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
variable "username" {
  type        = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type        = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
  type        = string
  description = "Stores the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}

variable "replication_consistency_group_name" {
  type        = string
  description = "The Replication Consistency Group name."
}

variable "action" {
  type        = string
  description = "The Replication Consistency Group action to be preformed. Options Failover, Restore, Sync, Reverse, Switchover and Snapshot"
}
