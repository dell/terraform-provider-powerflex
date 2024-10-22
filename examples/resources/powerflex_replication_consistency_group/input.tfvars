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
username_source               = "example_source_user"
password_source               = "example_source_password"
endpoint_source               = "example_source_endpoint"
name                          = "example_replication_consistancy_group_name"
rpo_in_seconds                = 15
source_protection_domain_name = "example_source_protection_domain"

# Destination Vars
username_destination               = "example_destination_user"
password_destination               = "example_destination_password"
endpoint_destination               = "example_destination_endpoint"
destination_protection_domain_name = "example_datasource_protection_domain"
