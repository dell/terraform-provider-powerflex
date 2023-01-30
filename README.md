# Terraform Provider for Dell Technologies PowerFlex

The Terraform Provider for Dell Technologies (Dell) PowerFlex allows Data Center and IT administrators to use Hashicorp Terraform to automate and orchestrate the provisioning and management of Dell PowerFlex storage systems.

The Terraform Provider can be used to manage SDCs, volumes, snapshots, snapshot-policies, storage pools, SDSs and protection domains.

## Table of contents

* [Support](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/SUPPORT.md)
* [License](#license)
* [Prerequisites](#prerequisites)
* [List of DataSources in Terraform Provider for Dell PowerFlex](#list-of-datasources-in-terraform-provider-for-dell-powerflex)
* [List of Resources in Terraform Provider for Dell PowerFlex](#list-of-resources-in-terraform-provider-for-dell-powerflex)
* [Releasing, Maintenance and Deprecation](#releasing-maintenance-and-deprecation)

## License
The Terraform Provide for PowerFlex is released and licensed under the MPL-2.0 license. See [LICENSE](https://github.com/dell/terraform-provider-powerflex/blob/main/LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerFlex/VxFlex OS Version** | **OS** | **Terraform** | **Golang** | **Terraform Plugin Framework version**              |
|---------------------|-----------------------|-------|--------------------|--------------------------|--------------------|
| v1.0.0 | 3.6 | ubuntu22.04 <br> <br> rhel8.x <br> rhel7.x | 1.3.2 <br> 1.2.9 <br> 1.3.2 <br> 1.2.9 <br> | 1.19.x | 1.0.1

## List of DataSources in Terraform Provider for Dell PowerFlex
  * [SDC](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/sdc.md)
  * [Storage pool](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/storagepool.md)
  * [Volume](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/volume.md)
  * [SDS](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/sds.md)
  * [Protection Domain](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/protection_domain.md)
  * [Snapshot Policy](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/data-sources/protection_domain.md)

## List of Resources in Terraform Provider for Dell PowerFlex
  * [SDC](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/resources/sdc.md)
  * [Storage pool](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/resources/storagepool.md)
  * [Volume](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/resources/volume.md)
  * [SDS](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/resources/sds.md)
  * [Snapshot](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/resources/snapshot.md)

## Installation and execution of Terraform Provider for Dell PowerFlex
The installation and execution steps of Terraform Provider for Dell PowerFlex can be found [here](https://github.com/dell/terraform-provider-powerflex/blob/main/docs/INSTALLATION.md).

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technnologies PowerFlex follows [Semantic Versioning](https://semver.org/).

New version will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags with names of the form "vx.y.z" where x.y.z corresponds to the version number.
