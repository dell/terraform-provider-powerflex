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
# Terraform Provider for Dell Technologies PowerFlex

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](LICENSE)

The Terraform Provider for Dell Technologies (Dell) PowerFlex allows Data Center and IT administrators to use Hashicorp Terraform to automate and orchestrate the provisioning and management of Dell PowerFlex storage systems.

The Terraform Provider can be used to manage SDCs, volumes, snapshots, snapshot-policies, storage pools, SDSs, protection domains, devices, users, MDM cluster, fault sets, firmware repository, peer systems, replication consistency groups, replication pairs, NVMe hosts and NVMe targets.

## Table of contents

* [Code of Conduct](https://github.com/dell/dell-terraform-providers/blob/main/docs/CODE_OF_CONDUCT.md)
* [Maintainer Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/MAINTAINER_GUIDE.md)
* [Committer Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/COMMITTER_GUIDE.md)
* [Contributing Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/CONTRIBUTING.md)
* [List of Adopters](https://github.com/dell/dell-terraform-providers/blob/main/docs/ADOPTERS.md)
* [Support](#support)
* [Security](https://github.com/dell/dell-terraform-providers/blob/main/docs/SECURITY.md)
* [License](#license)
* [Prerequisites](#prerequisites)
* [List of DataSources in Terraform Provider for Dell PowerFlex](#list-of-datasources-in-terraform-provider-for-dell-powerflex)
* [List of Resources in Terraform Provider for Dell PowerFlex](#list-of-resources-in-terraform-provider-for-dell-powerflex)
* [Releasing, Maintenance and Deprecation](#releasing-maintenance-and-deprecation)
* [Documentation](#documentation)
* [New to Terraform?](#new-to-terraform)

## Support
For any Terraform Provider for Dell PowerFlex issues, questions or feedback, please follow our [support process](https://github.com/dell/dell-terraform-providers/blob/main/docs/SUPPORT.md)

## License
The Terraform Provider for Dell PowerFlex is released and licensed under the MPL-2.0 license. See [LICENSE](https://github.com/dell/terraform-provider-powerflex/blob/main/LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerFlex/VxFlex OS Version** | **OS** | **Terraform** | **Golang** |
|---------------------|-----------------------|-------|--------------------|--------------------------|
| v1.7.0 | 3.6 <br> 4.5 <br> 4.6 | ubuntu22.04 <br> rhel9.x | 1.8.x <br> 1.9.x <br>| 1.23.x

## List of DataSources in Terraform Provider for Dell PowerFlex
  * [SDC](docs/data-sources/sdc.md)
  * [Storage pool](docs/data-sources/storage_pool.md)
  * [Volume](docs/data-sources/volume.md)
  * [SDS](docs/data-sources/sds.md)
  * [Protection Domain](docs/data-sources/protection_domain.md)
  * [Snapshot Policy](docs/data-sources/snapshot_policy.md)
  * [Device](docs/data-sources/device.md)
  * [VTree](docs/data-sources/vtree.md)
  * [Fault Set](docs/data-sources/fault_set.md)
  * [Node](docs/data-sources/node.md)
  * [Template](docs/data-sources/template.md)
  * [Service](docs/data-sources/service.md)
  * [Firmware Repository](docs/data-sources/firmware_repository.md)
  * [Peer System](docs/data-sources/peer_system.md)
  * [Replication Consistency Group](docs/data-sources/replication_consistency_group.md)
  * [Replication Pair](docs/data-sources/replication_pair.md)
  * [NVMe Host](docs/data-sources/nvme_host.md)
  * [NVMe Target](docs/data-sources/nvme_target.md)

## List of Resources in Terraform Provider for Dell PowerFlex
  * [SDC](docs/resources/sdc.md)
  * [Storage pool](docs/resources/storage_pool.md)
  * [Volume](docs/resources/volume.md)
  * [SDS](docs/resources/sds.md)
  * [Snapshot](docs/resources/snapshot.md)
  * [Protection Domain](docs/resources/protection_domain.md)
  * [SDC Volume Mapping](docs/resources/sdc_volumes_mapping.md)
  * [Device](docs/resources/device.md)
  * [Package](docs/resources/package.md)
  * [Cluster](docs/resources/cluster.md)
  * [MDM Cluster](docs/resources/mdm_cluster.md)
  * [User](docs/resources/user.md)
  * [Fault Set](docs/resources/fault_set.md)
  * [Service](docs/resources/service.md)
  * [Snapshot Policy](docs/resources/snapshot_policy.md)
  * [Firmware Repository](docs/resources/firmware_repository.md)
  * [SDC Host](docs/resources/sdc_host.md)
  * [System](docs/resources/system.md)
  * [Peer System](docs/resources/peer_system.md)
  * [Replication Consistency Group](docs/resources/replication_consistency_group.md)
  * [Replication Consistency Group Action](docs/resources/replication_consistency_group_action.md)
  * [Replication Pair](docs/resources/replication_pair.md)
  * [NVMe Host](docs/resources/nvme_host.md)
  * [NVMe Target](docs/resources/nvme_target.md)

## List of Modules in Terraform Provider for Dell PowerFlex
  * [User](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/user)
  * [SDC EXSi](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/sdc_host_esxi) 
  * [SDC Linux](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/sdc_host_linux) 
  * [SDC Windows](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/sdc_host_win)
  * [vSphere OVA Deployment](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/vsphere-ova-vm-deployment)
  * [Azure Deployment](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/azure_pfmp)
  * [AWS Block Storage Deployment](https://registry.terraform.io/modules/dell/modules/powerflex/latest/submodules/aws_install)

## Installation and execution of Terraform Provider for Dell PowerFlex
The installation and execution steps of Terraform Provider for Dell PowerFlex can be found [here](https://github.com/dell/terraform-provider-powerflex/blob/main/about/INSTALLATION.md).

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technologies PowerFlex follows [Semantic Versioning](https://semver.org/).

New versions will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags in the form of "vx.y.z" where x.y.z corresponds to the version number.

## Documentation

For more detailed information, please refer to [Dell Terraform Providers Documentation](https://dell.github.io/terraform-docs/).

## Terraform PowerFlex Modules

Check the following links for the terraform-modules repository and registry.
  * [Terraform PowerFlex Modules GitHub](https://github.com/dell/terraform-powerflex-modules)
  * [Terraform PowerFlex Modules Registry](https://registry.terraform.io/modules/dell/modules/powerflex/latest)

## New to Terraform?
**Here are some helpful links to get you started if you are new to terraform before using our provider:**

- Intro to Terraform: https://developer.hashicorp.com/terraform/intro 
- Providers: https://developer.hashicorp.com/terraform/language/providers 
- Resources: https://developer.hashicorp.com/terraform/language/resources
- Datasources: https://developer.hashicorp.com/terraform/language/data-sources
- Import: https://developer.hashicorp.com/terraform/language/import
- Variables: https://developer.hashicorp.com/terraform/language/values/variables
- Modules: https://developer.hashicorp.com/terraform/language/modules
- State: https://developer.hashicorp.com/terraform/language/state