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

## Required fields for all credential type
username = "example-username"
password = "example-password"
name = "example-resource-credential-name"
type = "Node" // Options: Node, Switch, vCenter, ElementManager, PowerflexGateway, PresentationServer, OSAdmin, OSUser

## Required value for vCenter, ElementManager, OSUser
#domain = "1.1.1.1"

## Required value for PowerflexGateway
#os_username = "example_os_username"
#os_password = "example_os_password"

## Optional for Node, Switch, ElementManager
#snmpv2_community_string = "public-test"

## Optional for Node 
#snmpv3_security_name = "example-security_name"
#snmpv3_security_level = "Moderate" // Options "Minimal", "Moderate", or "Maximal"
#snmpv3_md5_auth_password = "example_md5_auth" // required for level "Moderate" and "Maximal" 
#snmpv3_des_private_password = "example_des_private_password" // required for level "Maximal"

## Optional for Node, Switch, OSAdmin, OSUser
#ssh_private_key_path = "../../../../terraform-demos/resource-cred-test.pem" // Note: if either of these values are set they both need to be set
#key_pair_name = "example_key_pair_name" // Note: if either of these values are set they both need to be set  
