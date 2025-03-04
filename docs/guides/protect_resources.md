---
page_title: "Protecting Terraform Resources"
title: "Protecting Terraform Resources"
linkTitle: "Protecting Terraform Resources"
---

<!--
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
-->
# Protecting Terraform Resources by Managing Resource Lifecycle
Terraform provides a lifecycle block within resource configurations to control the creation, update, and deletion of resources. This block is essential for managing resource behavior and ensuring critical resources are not accidentally destroyed.

Key Arguments in the lifecycle Block
## prevent_destroy:
Prevents Terraform from destroying the resource. If a plan includes the destruction of this resource, Terraform will produce an error instead.
```
resource "powerflex_example" "example" {
  # ... other configuration ...

  lifecycle {
    prevent_destroy = true
  }
}
```

## ignore_changes:
Ignores changes to specified attributes of the resource. This is useful when certain attributes are managed outside of Terraform.
```
resource "eg_instance" "example" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}
```

## create_before_destroy:
Ensures that a new resource is created before the existing one is destroyed. This is useful in certain cases and for minimizing downtime. 
```
resource "powerflex_example" "example" {
 # ... some other configuration ...

  lifecycle {
    create_before_destroy = true
  }
}
```

### Example: Preventing Accidental Deletion
To ensure a critical resource, such as a volume instance, is not accidentally deleted, you can use the prevent_destroy argument:
```
resource "powerflex_volume" "example-volume-create" {
  name = "example-volume-create"
  protection_domain_name = "domain1"
  storage_pool_name = "pool1"

  # The unit of size of the volume is defined by capacity_unit whose default value is "GB".
  size          = 8
  capacity_unit = "GB" # GB/TB

  use_rm_cache = true
  volume_type  = "ThickProvisioned"      # ThickProvisioned/ThinProvisioned volume type
  access_mode  = "ReadWrite"             # ReadWrite/ReadOnly volume access mode
  remove_mode  = "INCLUDING_DESCENDANTS" # INCLUDING_DESCENDANTS/ONLY_ME remove mode
  lifecycle {
    prevent_destroy = true
  }
}
```
In this example, the prevent_destroy argument ensures that any attempt to destroy the powerflex_volume resource will result in an error, protecting it from accidental deletion.

### Additional Resources
For more detailed information on the lifecycle block and its arguments, refer to the Terraform documentation on [lifecycle meta-arguments](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle)
