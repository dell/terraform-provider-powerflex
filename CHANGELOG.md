<!--
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# v1.4.0 (March 27, 2024)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_service` for managing Service in PowerFlex.
* `powerflex_snapshot_policy` for managing Snapshot Policy in PowerFlex.

### Data Sources

* `powerflex_node` for reading Node(Resource) details in PowerFlex.
* `powerflex_template` for reading Template details in PowerFlex.
* `powerflex_service` for reading Service details in PowerFlex.
  
### Enhancements

* All existing resources and datasources are qualified against PowerFlex v4.5 on Azure.
* Added support for including SDS in fault set.
* Added support virtual ip/interfcaes in SDC resource.
  
### Bug Fixes

* For SDC Resource, Fixed the change in value of is_sdc from No to Yes, which was previously giving an error "Provider produced inconsistent results after apply".

# v1.3.0 (December 22, 2023)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_fault_set` for managing Fault Set in PowerFlex.

### Data Sources

* `powerflex_fault_set` for reading Fault Set details in PowerFlex.

### Enhancements

* All existing resources and datasources are qualified against PowerFlex v4.5.
* Support installation for SDC on nodes running RockyLinux OS.

# v1.2.0 (September 27, 2023)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_cluster` for deploying PowerFlex cluster.
* `powerflex_mdm_cluster` for managing MDM cluster in PowerFlex.
* `powerflex_user` for managing users in PowerFlex.

### Data Sources

* `powerflex_vtree` for reading VTree details in PowerFlex.

### Notes

* `name` attribute is removed from sdc resource.

# v1.1.0 (June 28, 2023)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_device` for managing devices in PowerFlex.
* `powerflex_protection_domain` for managing protection domains in PowerFlex.
* `powerflex_package` for managing packages on the PowerFlex Gateway.

### Data Sources

* `powerflex_device` for reading devices in PowerFlex.

### Enhancements

* `powerflex_storage_pool` is enhanced to support additional attributes in PowerFlex.
* `powerflex_sdc` is enhanced to create/delete/import multiple SDCs in PowerFlex.

### Deprecations

* `name` attribute from SDC resource.

### Notes

* `sdc_list` attribute is removed from volume and snapshot resource.

# v1.0.1 (May 23, 2023)

## Release Summary

The release supports resources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_sdc_volumes_mapping` for managing map/unmap operations between SDC and volumes in PowerFlex.

### Deprecations

* sdc_list attribute in powerflex_volume and powerflex_snapshot resource.

***
<br>

# v1.0.0 (Feb 27, 2023)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Data Sources

* `powerflex_protection_domain` for reading protection domain in PowerFlex.
* `powerflex_sdc` for reading SDC in PowerFlex.
* `powerflex_sds` for reading SDS in PowerFlex.
* `powerflex_snapshot_policy` for reading snapshot policy in PowerFlex.
* `powerflex_storage_pool` for reading storage pool in PowerFlex.
* `powerflex_volume` for reading volume in PowerFlex.

### Resources

* `powerflex_sdc` for managing SDC in PowerFlex.
* `powerflex_sds` for managing SDS in PowerFlex.
* `powerflex_snapshot` for managing Snapshot in PowerFlex.
* `powerflex_storage_pool` for managing Storage Pool in PowerFlex.
* `powerflex_volume` for managing Volume in PowerFlex.

### Others

N/A

### Enhancements

N/A

### Bug Fixes

N/A
