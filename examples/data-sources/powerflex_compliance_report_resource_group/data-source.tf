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

# To get the resource_group_id you can use the powerflex_resource_group data source as shown below:

# Get all Resource Group details present in the PowerFlex
data "powerflex_resource_group" "example1" {
}

# Get Resource Group details using the Name of the Resource Group
data "powerflex_resource_group" "example3" {
  resource_group_names = ["Name_1", "Name_2"]
}

# Get all compliance report details for the given resource group
data "powerflex_compliance_report_resource_group" "complianceReport" {
    resource_group_id = "ID"
}

# Get compliance report details for the given resource group filtered by given ipaddresses
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    ip_addresses = ["10.xxx.xxx.xx","10.xxx.xxx.xx"]
  }
}

# Get compliance report details for the given resource group filtered by resource ids and compliant status
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    resource_ids = ["resourceid1","resourceid2"]
    compliant = true
  }
}

# Get compliance report details for the given resource group filtered by compliant resources
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    compliant = true
  }
}

# Get compliance report details for the given resource group filtered by hostnames
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    host_names = ["hostname1","hostname2"]
  }
}

# Get compliance report details for the given resource group filtered by service tags
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    service_tags = ["servicetag1","servicetag2"]
  }
}

# Get compliance report details for the given resource group filtered by resource ids
data "powerflex_compliance_report_resource_group" "complianceReport" {
  resource_group_id = "ID"
  # this datasource supports multiple filters like ip_addresses, host_names, service_tags, resource_ids, compliant
  # and gives an intersection of the results
  filter {
    resource_ids = ["resourceid1","resourceid2"]
  }
}

output "result" {
  value = data.powerflex_compliance_report_resource_group.complianceReport
}