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

# commands to run this tf file : terraform init && terraform apply --auto-approve
# Get all Resource Group details present in the PowerFlex
data "powerflex_compliance_report_resource_group" "example1" {
}

# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_compliance_report_resource_group" "compliance_report_resource_group_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     model = ["^Powerflex.*$"]
#   }
# }

# output "complianceReportResourceGroupFilterRegexResult"{
#  value = data.powerflex_compliance_report_resource_group.compliance_report_resource_group_filter_regex.compliance_reports
# }

# this datasource supports multiple filters like ip_address, host_name, service_tag, compliant,etc.
# Note: If both filters are used simultaneously, the results will include any records that match either of the filters.
# data "powerflex_compliance_report_resource_group" "complianceReport" {
#   filter {
#      firmware_repository_name = ["name1", "name2"]
#      device_type  = ["devicetype1", "devicetype2"]
#      model = ["model1", "model2"]
#      service_tag = ["servicetag1", "servicetag2"]
#      compliant = true
#      embedded_report = true
#      can_update = true
#      available = true
#      device_state  = ["devicestate1", "devicestate2"]
#      managed_state = ["managedstate1", "managedstate2"]
#      host_name = ["hostname1", "hostname2"]
#      id = ["id1", "id2"]
#      ip_address = ["ip1", "ip2"]
#   }
# }

output "result" {
  value = data.powerflex_compliance_report_resource_group.example1.compliance_reports
}