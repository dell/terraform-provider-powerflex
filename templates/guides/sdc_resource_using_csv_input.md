---
page_title: "Deploying a PowerFlex SDC using a CSV Topology File"
title: "Deploying a PowerFlex SDC using a CSV Topology File"
linkTitle: "Deploying a PowerFlex SDC using a CSV Topology File"
---

<!--
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Deploying a PowerFlex SDC using a CSV Topology File

This guide explains how to use a CSV file as input for an SDC Resource.

### Example

~> **Note:** In the CSV file, provide only Cluster Node-related details as shown in the example below:

**[sdc_input_minimal.csv](https://dell-my.sharepoint.com/:x:/p/krunal_thakkar/Eb6I5vEbfANCpsf2BCC1AC4BfZzMsNzkyI1Vfv6CPUypQw?e=UX8WPg)**  &&  **[sdc_complete_config.csv](https://dell-my.sharepoint.com/:x:/p/krunal_thakkar/EYYDXGjdIRNCtPxqAPk0rNcBMCvbf3ogB2m-bPOg4DkLJw?e=Nwp1TU)**

```
IPs,Password,Operating System,Is MDM/TB,Is SDC,perfProfileForSDC,SDC Name
10.10.10.10,password,linux,Primary,Yes,HighPerformance,SDC1
10.10.10.11,password,linux,Secondary,Yes,HighPerformance,SDC2
10.10.10.12,password,linux,TB,Yes,Compact,SDC3
10.10.10.13,password,linux,Standby,Yes,Compact,SDC4
```

To perform SDC operations via CSV file(s), use the following configuration:

```terraform

locals {
  csv_data = csvdecode(file("data/sdc_input.csv"))
}

resource "powerflex_sdc" "sdc_rs" {
  mdm_password = "MDM_Password"
  lia_password = "LIA_Password"

  sdc_details = [
    for row in local.csv_data : {
      ip                  = row.IPs # if using complere config use row["MDM IPs"]
      password            = row.Password
      operating_system    = row["Operating System"]
      is_mdm_or_tb        = row["Is MDM/TB"]
      is_sdc              = row["Is SDC"]
      name                = row["SDC Name"]
      performance_profile = row.perfProfileForSDC
    }
    if row.IPs != "" && row.IPs != "null" && row.Password != "null"
  ]
}
```
This code demonstrates how you can use data from the CSV file to dynamically configure the powerflex_sdc resource.

`locals`: This section defines local variables within the Terraform configuration. These variables are used to store data or perform calculations within your configuration.

`csv_data`: This is a local variable that is being assigned the result of csvdecode(file("sdc_input.csv")). It decodes the content of the "sdc_input.csv" file into a data structure that can be used in the Terraform configuration.

`file("sdc_input.csv")`: This part of the code specifies the path to the CSV file. In this example, it assumes that "sdc_input.csv" is located in the same directory as your Terraform configuration file. If the CSV file is in a different directory, you can provide the relative or absolute path to that file. For example, if the CSV file is in a subdirectory named "data," you can specify the path as file("data/sdc_input.csv").

`sdc_details`: This is an attribute of the powerflex_sdc resource. You are using a list comprehension (the for loop) to create a list of objects based on the data from the CSV file (local.csv_data). Each object represents a configuration for an SDC (Storage Data Client).

The fields inside the object (e.g., ip, password, operating_system, etc.) are populated with values from the CSV file. You are also applying some conditions in the if statement to filter out rows where IPs and Password are not empty and not equal to "null."