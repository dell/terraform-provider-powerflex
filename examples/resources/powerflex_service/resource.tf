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


# Commands to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To create / update, either template_id or firmware_id must be provided
# deployment_name and deployment_description are the required parameters to create or update
# other  atrributes like : nodes, port, deployment_timeout are optional 
# To check which attributes can be updated, please refer Product Guide in the documentation

resource "powerflex_service" "service" {
  deployment_name        = "Test-Create-U1"
  deployment_description = "Test Service-U1"
  template_id            = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
  firmware_id            = "8aaaee208c8c467e018cd37813250614"
  nodes                  = 5
  clone_from_host        = "pfmc-k8s-20230809-160"
  deployment_timeout     = 50
}
