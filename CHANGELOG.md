<!--
Copyright (c) 2022-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# v1.8.0 (March 2025)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

* Added Support for Powerflex Appliance 4.6.1

### Resources

* `powerflex_resource_credential` for managing resource credential details in PowerFlex.
* `powerflex_template_clone` for cloning pre-existing sample resource group templates

### Data Sources

* `powerflex_resource_credential` for reading resource credential details in PowerFlex.
  
### Enhancements

* Added support for SDT deployment using Cluster resource.

# v1.7.0 (December 2024)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_peer_system` for managing Peer Systems in PowerFlex.
* `powerflex_replication_consistency_group` for managing Replication Consistency Groups in PowerFlex.
* `powerflex_replication_pair` for managing Replication Pairs in PowerFlex.
* `powerflex_replication_consistency_group_action` for performing action on Replication Consistency Group in PowerFlex.
* `powerflex_nvme_host` for managing NVMe Hosts in PowerFlex.
* `powerflex_nvme_target` for managing NVMe Targets in PowerFlex.

### Data Sources

* `powerflex_peer_system` for reading Peer System details in PowerFlex.
* `powerflex_replication_consistency_group` for reading Replication Consistency Group details in PowerFlex.
* `powerflex_replication_pair` for reading Replication Pair details in PowerFlex.
* `powerflex_nvme_host` for reading NVMe Host details in PowerFlex.
* `powerflex_nvme_target` for reading NVMe Target details in PowerFlex.
  
### Enhancements

* Added support for SDT deployment using Cluster resource.

### Limitations

* NVMe over TCP is supported in PowerFlex 4.0 and later versions, therefore `powerflex_nvme_host`, `powerflex_nvme_target` and SDT deployment using Cluster resource are not supported in PowerFlex 3.x.
* Due to certain limitations, updating the NVMe host in PowerFlex versions earlier than 4.6 is not supported

# v1.6.0 (Aug 2024)

### Removed Deprecated Resources
* `powerflex_service`
* `powerflex_sdc`

### New Resources and Datasources
* `powerflex_compatibility_management`
* `powerflex_compliance_report`
* `powerflex_os_repositiory`

# v1.5.0 (June 28, 2024)

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.

## Features

### Resources

* `powerflex_firmware_repository` for managing firmware repository in PowerFlex.
* `powerflex_sdc_host` for managing SDCs in PowerFlex.
* `powerflex_system` for managing system level configuration in PowerFlex

### Data Sources

* `powerflex_firmware_repository` for reading firmware repository details in PowerFlex.
  
### Enhancements

* All existing resources and datasources are qualified against PowerFlex v4.5 on AWS.
* All existing resources and datasources are qualified against PowerFlex v4.6.
* Added support for multiple mdm ips in Cluster/SDC resource.
* Added support for passwordless authentication in Cluster resource.

### Deprecations

* Service Resource/Datasource is deprecated.
* SDC resource is deprecated.

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
