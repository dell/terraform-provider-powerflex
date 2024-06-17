---
page_title: "Import volumes mapped to SDC"
title: "Import volumes mapped to SDC"
linkTitle: "Import volumes mapped to SDC"
---

<!--
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
-->

# Importing volumes mapped to SDC using sdc_volumes_mapping resource

This guide explains how to import volumes mapped to SDC using sdc_volumes_mapping_resource. Below steps are taken from [this article](https://developer.hashicorp.com/terraform/language/import/generating-configuration).

### Step 1: Add the import block
```
    import {
        to = powerflex_sdc_volumes_mapping.test
        id = "sdc_id"
    }
```

### Step 2: Plan and generate configuration
```
    terraform plan -generate-config-out=test.tf
```

### Step 3: Review generated configuration in test.tf

### Step 4: Run terraform apply to import the volumes mapped to SDC. Choose only one argument from mutually exclusive arguments.<br>
For example, choose either volume_id or volume_name from sdc_volumes_mapping resource.