---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerflex_compliance_report_resource_group data source"
linkTitle: "powerflex_compliance_report_resource_group"
page_title: "powerflex_compliance_report_resource_group Data Source - powerflex"
subcategory: ""
description: |-
  This datasource is used to query the compliance report for Resource Group from PowerFlex array.
---

# powerflex_compliance_report_resource_group (Data Source)

This datasource is used to query the compliance report for Resource Group from PowerFlex array.

## Example Usage

```terraform
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
```

After the successful execution of above said block, we can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_compliance_report_resource_group.datasource_block_name.attribute_name` where datasource_block_name is the name of the data source block and attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `resource_group_id` (String) Unique id Of resource group for which you want to get the compliance report. Conflicts with resource_group_name

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `compliance_reports` (Attributes List) List of compliance report. (see [below for nested schema](#nestedatt--compliance_reports))
- `id` (String) Unique identifier Of The compliance report Datasource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `compliant` (Boolean) Compliant status for resources.
- `host_names` (Set of String) List of host names for resources.
- `ip_addresses` (Set of String) List of Ip Address for resources.
- `resource_ids` (Set of String) List of resource ids.
- `service_tags` (Set of String) List of service tags for resources.


<a id="nestedatt--compliance_reports"></a>
### Nested Schema for `compliance_reports`

Read-Only:

- `available` (Boolean) The availability status of the device.
- `can_update` (Boolean) The update capability of the resource group.
- `compliant` (Boolean) The compliance status of the resource.
- `device_state` (String) The state of the device.
- `device_type` (String) The type of the device.
- `embedded_report` (Boolean) The presence of an embedded report.
- `firmware_compliance_report_components` (Attributes List) The list of firmware compliance report components. (see [below for nested schema](#nestedatt--compliance_reports--firmware_compliance_report_components))
- `firmware_repository_name` (String) The name of the firmware repository.
- `host_name` (String) The hostname of the resource group.
- `id` (String) The unique identifier of the resource group.
- `ip_address` (String) The IP address of the resource.
- `managed_state` (String) The managed state of the device.
- `model` (String) The model of the device.
- `service_tag` (String) The service tag of the resource.

<a id="nestedatt--compliance_reports--firmware_compliance_report_components"></a>
### Nested Schema for `compliance_reports.firmware_compliance_report_components`

Read-Only:

- `compliant` (Boolean) The compliance status of the component.
- `current_version` (Attributes) The current version of the component. (see [below for nested schema](#nestedatt--compliance_reports--firmware_compliance_report_components--current_version))
- `id` (String) The unique identifier of the component.
- `name` (String) The name of the component.
- `operating_system` (String) The operating system of the component.
- `os_compatible` (Boolean)
- `rpm` (Boolean)
- `software` (Boolean)
- `target_version` (Attributes) The target version of the component. (see [below for nested schema](#nestedatt--compliance_reports--firmware_compliance_report_components--target_version))
- `vendor` (String) The vendor of the component.

<a id="nestedatt--compliance_reports--firmware_compliance_report_components--current_version"></a>
### Nested Schema for `compliance_reports.firmware_compliance_report_components.current_version`

Read-Only:

- `firmware_last_update` (String) The last update time of the firmware.
- `firmware_level` (String) The level of the firmware.
- `firmware_name` (String) The name of the firmware.
- `firmware_type` (String) The type of the firmware.
- `firmware_version` (String) The version of the firmware.
- `id` (String) The unique identifier of the version.


<a id="nestedatt--compliance_reports--firmware_compliance_report_components--target_version"></a>
### Nested Schema for `compliance_reports.firmware_compliance_report_components.target_version`

Read-Only:

- `firmware_last_update` (String) The last update time of the firmware.
- `firmware_level` (String) The level of the firmware.
- `firmware_name` (String) The name of the firmware.
- `firmware_type` (String) The type of the firmware.
- `firmware_version` (String) The version of the firmware.
- `id` (String) The unique identifier of the version.

