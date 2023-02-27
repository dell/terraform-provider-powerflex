<!--
Copyright (c) 2022 Dell Inc., or its subsidiaries. All Rights Reserved.

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

The Terraform Provider can be used to manage SDCs, volumes, snapshots, snapshot-policies, storage pools, SDSs and protection domains.

## Table of contents

* [Support](https://github.com/dell/dell-terraform-providers/blob/main/docs/SUPPORT.md)
* [License](#license)
* [Prerequisites](#prerequisites)
* [List of DataSources in Terraform Provider for Dell PowerFlex](#list-of-datasources-in-terraform-provider-for-dell-powerflex)
* [List of Resources in Terraform Provider for Dell PowerFlex](#list-of-resources-in-terraform-provider-for-dell-powerflex)
* [Releasing, Maintenance and Deprecation](#releasing-maintenance-and-deprecation)

## License
The Terraform Provide for PowerFlex is released and licensed under the MPL-2.0 license. See [LICENSE](LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerFlex/VxFlex OS Version** | **OS** | **Terraform** | **Golang** |
|---------------------|-----------------------|-------|--------------------|--------------------------|
| v1.0.0 | 3.6 | ubuntu22.04 <br> rhel8.x <br> rhel7.x | 1.3.2 <br> 1.2.9 <br> | 1.19.x

## List of DataSources in Terraform Provider for Dell PowerFlex
  * [SDC](docs/data-sources/sdc.md)
  * [Storage pool](docs/data-sources/storagepool.md)
  * [Volume](docs/data-sources/volume.md)
  * [SDS](docs/data-sources/sds.md)
  * [Protection Domain](docs/data-sources/protection_domain.md)
  * [Snapshot Policy](docs/data-sources/protection_domain.md)

## List of Resources in Terraform Provider for Dell PowerFlex
  * [SDC](docs/resources/sdc.md)
  * [Storage pool](docs/resources/storagepool.md)
  * [Volume](docs/resources/volume.md)
  * [SDS](docs/resources/sds.md)
  * [Snapshot](docs/resources/snapshot.md)

## Installation and execution of Terraform Provider for Dell PowerFlex
The installation and execution steps of Terraform Provider for Dell PowerFlex can be found [here](about/INSTALLATION.md).

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technnologies PowerFlex follows [Semantic Versioning](https://semver.org/).

New version will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags with names of the form "vx.y.z" where x.y.z corresponds to the version number.
