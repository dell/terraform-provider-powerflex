---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerflex_cluster resource"
linkTitle: "powerflex_cluster"
page_title: "powerflex_cluster Resource - powerflex"
subcategory: "Cluster and System"
description: |-
  This terraform resource is used to deploy the PowerFlex Cluster. We can Create and Delete the PowerFlex Cluster using this resource. We can also Import an existing Cluster of the PowerFlex.
---

# powerflex_cluster (Resource)

This terraform resource is used to deploy the PowerFlex Cluster. We can Create and Delete the PowerFlex Cluster using this resource. We can also Import an existing Cluster of the PowerFlex.

**Please consider the following points before using cluster resource.**

1. For PowerFlex 4.x, the PowerFlex Manager must be installed as a prerequisite. The required packages should be uploaded to the PowerFlex Manager.

2. For PowerFlex 3.x, a Gateway server is a prerequisite. The required packages should be uploaded to the gateway. The Package resource can be used for uploading packages to the gateway.

3. Support is provided for creating, importing, and deleting operations for this resource.

4. In multi-node cluster deployments, when some of the component installations fail, the partial deployment will not be rolled back.

5. If you've separately installed any SDR, SDS, or SDC and connected it to the cluster and if you face any security certificate issues during the destroy process, you'll have to manually accept the security certificate to resolve them.

6. During the destroy process, the entire cluster will be destroyed, not just specific individual resources. After destroy need to follow cleanup process.

7. `ips` attribute is used in minimal csv configuration whereas `mdm_ips` attribute is used in complete csv configuration.

8. For PowerFlex 4.x, there's no need to mention `allow_non_secure_communication_with_lia`, `allow_non_secure_communication_with_mdm`, and `disable_non_mgmt_components_auth`. And, **Rfcache** is not supported.

9. To follow the installation process, you can refer to the [Deployment Guide 3.x](https://www.dell.com/support/manuals/en-us/scaleio/pfx_deploy_guide_3.6.x/deploy-powerflex?guid=guid-e9f70972-baac-42c9-9ff9-a3d2b0722f54&lang=en-us) & [Deployment Guide 4.x](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/introduction?guid=guid-e798f431-7df4-450c-8f86-60ee7f3d1e3e&lang=en-us)

10. NVMe over TCP is supported in PowerFlex 4.0 and later versions, therefore SDT deployment is not supported in PowerFlex 3.x.

## Example Usage

```terraform
/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Read, Delete and Import operations are supported for this resource.

# Example for deploying cluster. After successful execution, 3 node MDM cluster will be deployed with 3 SDCs and 2 SDS.
resource "powerflex_package" "upload-test" {
  file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm",
  "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdt-3.6-700.103.el7.x86_64.rpm"]
}

resource "powerflex_cluster" "test" {

  depends_on = [powerflex_package.upload-test]

  # Security Related Field
  mdm_password = "Password"
  lia_password = "Password"

  # Advance Security Configuration
  allow_non_secure_communication_with_lia = false
  allow_non_secure_communication_with_mdm = false
  disable_non_mgmt_components_auth        = false

  # Cluster Configuration related fields 
  cluster = [
    {
      # MDM Configuration Fields 
      ips                  = "10.10.10.1",
      username             = "root",
      password             = "Password",
      operating_system     = "linux",
      is_mdm_or_tb         = "primary",
      mdm_ips              = "10.10.10.1",
      mdm_mgmt_ip          = "10.10.10.1",
      mdm_name             = "MDM_1",
      perf_profile_for_mdm = "HighPerformance",
      virtual_ips          = "10.30.30.1",
      virtual_ip_nics      = "ens192",

      # SDS Configuration Fields
      is_sds      = "yes",
      sds_name    = "sds1",
      sds_all_ips = "10.20.20.3", # conflict with sds_to_sds_only_ips,sds_to_sdc_only_ips
      # sds_to_sdc_only_ips      = "10.20.20.2", 
      # sds_to_sds_only_ips      = "10.20.20.1",
      fault_set                = "fs1",
      protection_domain        = "domain_1"
      sds_storage_device_list  = "/dev/sdb"
      sds_storage_device_names = "device1"
      storage_pool_list        = "pool1"
      perf_profile_for_sds     = "HighPerformance"

      # SDC Configuration Fields
      is_sdc               = "yes",
      sdc_name             = "sdc1",
      perf_profile_for_sdc = "HighPerformance",

      # Rfcache Configuration Fields
      is_rfcache               = "No",
      rf_cache_ssd_device_list = "/dev/sdd"

      # SDR Configuration Fields
      is_sdr   = "Yes",
      sdr_name = "SDR_1"
      sdr_port = "2000"
      # sdr_application_ips  = "10.20.30.1"
      # sdr_storage_ips      = "10.20.30.2"
      # sdr_external_ips     = "10.20.30.3" 
      sdr_all_ips          = "10.10.20.1" # conflict with sdr_application_ips, sdr_storage_ips, sdr_external_ips
      perf_profile_for_sdr = "Compact"

      # SDT Configuration Fields
      is_sdt      = "Yes"
      sdt_name    = "SDT_1"
      sdt_all_ips = "10.20.40.1"
    },
    {
      ips                     = "10.10.10.2",
      username                = "root",
      password                = "Password",
      operating_system        = "linux",
      is_mdm_or_tb            = "Secondary",
      protection_domain       = "domain_1"
      sds_storage_device_list = "/dev/sdb"
      storage_pool_list       = "pool1"
      is_sds                  = "yes",
      sds_name                = "sds2",
      is_sdc                  = "yes",
      sdc_name                = "sdc2",
      perf_profile_for_sdc    = "compact",
      is_rfcache              = "No",
      is_sdr                  = "No",
      is_sdt                  = "Yes"
      sdt_name                = "SDT_2"
      sdt_all_ips             = "10.20.40.2"
    },
    {
      ips                  = "10.10.10.3",
      username             = "root",
      password             = "Password",
      operating_system     = "linux",
      is_mdm_or_tb         = "TB",
      is_sds               = "No",
      is_sdc               = "yes",
      sdc_name             = "sdc3",
      perf_profile_for_sdc = "compact",
      is_rfcache           = "No",
      is_sdr               = "No",
      is_sdt               = "No"
    },
  ]
  # Storage Pool Configuration Fields
  storage_pools = [
    {
      media_type                              = "HDD"
      protection_domain                       = "domain_1"
      storage_pool                            = "pool1"
      replication_journal_capacity_percentage = "50"
    }
  ]
}
```

After the execution of above resource block, Cluster would have been created on the PowerFlex array. For more information, please check the terraform state file.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (Attributes List) Cluster Installation Details (see [below for nested schema](#nestedatt--cluster))
- `lia_password` (String) Lia Password
- `mdm_password` (String) MDM Password

### Optional

- `allow_non_secure_communication_with_lia` (Boolean) Allow Non Secure Communication With lia
- `allow_non_secure_communication_with_mdm` (Boolean) Allow Non Secure Communication With MDM
- `disable_non_mgmt_components_auth` (Boolean) Disable Non Mgmt Components Auth
- `id` (String) ID
- `storage_pools` (Attributes List) Storage Pool Details (see [below for nested schema](#nestedatt--storage_pools))

### Read-Only

- `mdm_list` (Attributes Set) Cluster MDM Details (see [below for nested schema](#nestedatt--mdm_list))
- `protection_domains` (Attributes List) Cluster Protection Domain Details (see [below for nested schema](#nestedatt--protection_domains))
- `sdc_list` (Attributes Set) Cluster SDC Details (see [below for nested schema](#nestedatt--sdc_list))
- `sdr_list` (Attributes Set) Cluster SDR Details (see [below for nested schema](#nestedatt--sdr_list))
- `sds_list` (Attributes Set) Cluster SDS Details (see [below for nested schema](#nestedatt--sds_list))
- `sdt_list` (Attributes Set) Cluster SDT Details (see [below for nested schema](#nestedatt--sdt_list))

<a id="nestedatt--cluster"></a>
### Nested Schema for `cluster`

Required:

- `is_mdm_or_tb` (String) Is Mdm Or Tb
- `operating_system` (String) Operating System

Optional:

- `fault_set` (String) Fault Set
- `ips` (String) Use this field to assign a single IP address for all the MDM IP, MDM Mgmt IP, and SDS All IP. This option is useful when separate networks for data and management are not required.
- `is_rfcache` (String) Is RFCache. The acceptable values are `Yes` and `No`. Default value is `No`.
- `is_sdc` (String) Is Sdc. The acceptable values are `Yes` and `No`. Default value is `No`.
- `is_sdr` (String) Is SDR. The acceptable values are `Yes` and `No`. Default value is `No`.
- `is_sds` (String) Is Sds. The acceptable values are `Yes` and `No`. Default value is `No`.
- `is_sdt` (String) Is Sdt. The acceptable values are `Yes` and `No`. Default value is `No`.
- `mdm_ips` (String) MDM IP addresses used to communicate with other PowerFlex components in the storage network. This is required for all MDMs, Tiebreakers and Standbys.Leave this field blank for hosts that are not part of the MDM cluster.
- `mdm_mgmt_ip` (String) This IP address is for the management-only network. The management ip is not required for Tiebreaker MDM, Standby Tiebreaker MDM and any host that is not an MDM.
- `mdm_name` (String) MDMName
- `password` (String, Sensitive) Password used to log in to the node.
- `perf_profile_for_mdm` (String) Performance Profile For MDM
- `perf_profile_for_sdc` (String) Performance Profile For SDC
- `perf_profile_for_sdr` (String) Performance Profile For SDR
- `perf_profile_for_sds` (String) Performance Profile For SDS
- `protection_domain` (String) Protection Domain
- `rf_cache_ssd_device_list` (String) List of SSD devices to provide RFcache acceleration for Medium Granularity data layout Storage Pools.
- `sdc_name` (String) SDC Name
- `sdr_all_ips` (String) SDR IP addresses to be used for communication among all nodes (including all three roles)
- `sdr_application_ips` (String) The IP addresses through which the SDC communicates with the SDR.
- `sdr_external_ips` (String) The IP addresses through which the SDR communicates with peer systems SDRs
- `sdr_name` (String) SDR Name
- `sdr_port` (String) SDR Port
- `sdr_storage_ips` (String) The IP addresses through which the SDR communicates with the MDM for server side control communications.
- `sds_all_ips` (String) SDS IP addresses to be used for communication among all nodes.
- `sds_name` (String) SDS Name
- `sds_storage_device_list` (String) Storage devices to be added to an SDS. For more than one device, use a comma separated list, with no spaces.
- `sds_storage_device_names` (String) Sets names for devices.
- `sds_to_sdc_only_ips` (String) SDS IP addresses to be used for communication among SDS and SDC nodes only.
- `sds_to_sds_only_ips` (String) SDS IP addresses to be used for communication among SDS nodes. When the replication feature is used, these addresses are also used for SDS-SDR communication.
- `sdt_all_ips` (String) SDT IP addresses used for both hosts communication and MDM communication (including both roles).
- `sdt_name` (String) SDT Name
- `storage_pool_list` (String) Sets Storage Pool names
- `username` (String) The value can be either `root` or any non-root user name with appropriate permissions.
- `virtual_ip_nics` (String) The NIC to which the virtual IP addresses are mapped.
- `virtual_ips` (String) Virtual IPs


<a id="nestedatt--storage_pools"></a>
### Nested Schema for `storage_pools`

Required:

- `media_type` (String) Media Type

Optional:

- `compression_method` (String) Compression Method
- `data_layout` (String) Data Layout
- `external_acceleration` (String) External Acceleration
- `protection_domain` (String) Protection Domain
- `replication_journal_capacity_percentage` (String) Replication Journal Capacity Percentage
- `storage_pool` (String) Storage Pool
- `zero_padding` (String) Zero Padding


<a id="nestedatt--mdm_list"></a>
### Nested Schema for `mdm_list`

Read-Only:

- `id` (String) ID
- `ip` (String) MDM Node IP
- `mdm_ip` (String) MDM IP
- `mgmt_ip` (String) MGMTIP
- `mode` (String) Mode
- `name` (String) Name
- `role` (String) Role
- `virtual_ip` (String) Virtual IP
- `virtual_ip_nic` (String) Virtual IPNIC


<a id="nestedatt--protection_domains"></a>
### Nested Schema for `protection_domains`

Read-Only:

- `name` (String) Name
- `storage_pool_list` (Attributes List) Storage Pools (see [below for nested schema](#nestedatt--protection_domains--storage_pool_list))

<a id="nestedatt--protection_domains--storage_pool_list"></a>
### Nested Schema for `protection_domains.storage_pool_list`

Read-Only:

- `compression_method` (String) Compression Method
- `data_layout` (String) Data Layout
- `external_acceleration` (String) External Acceleration
- `media_type` (String) Media Type
- `name` (String) Name
- `replication_journal_capacity_percentage` (Number) Replication Journal Capacity Percentage
- `zero_padding` (String) Zero Padding



<a id="nestedatt--sdc_list"></a>
### Nested Schema for `sdc_list`

Read-Only:

- `guid` (String) GUID
- `id` (String) ID
- `ip` (String) SDC Node IP
- `name` (String) Name


<a id="nestedatt--sdr_list"></a>
### Nested Schema for `sdr_list`

Read-Only:

- `all_ips` (String) All IP
- `application_ips` (String) Application IP
- `external_ips` (String) External IP
- `id` (String) ID
- `ip` (String) SDR Node IP
- `name` (String) Name
- `port` (Number) Port
- `storage_ips` (String) Storage IP


<a id="nestedatt--sds_list"></a>
### Nested Schema for `sds_list`

Read-Only:

- `all_ips` (String) All IP
- `devices` (Attributes Set) Devices (see [below for nested schema](#nestedatt--sds_list--devices))
- `fault_set` (String) Fault Set
- `id` (String) ID
- `ip` (String) SDS Node IP
- `name` (String) Name
- `protection_domain_id` (String) Protection Domain Name
- `protection_domain_name` (String) Protection Domain Name
- `sds_only_ips` (String) SDSOnly IP
- `sds_sdc_ips` (String) SDSSDCIP

<a id="nestedatt--sds_list--devices"></a>
### Nested Schema for `sds_list.devices`

Read-Only:

- `max_capacity_in_kb` (Number) Max Capacity
- `name` (String) Name
- `path` (String) Path
- `storage_pool` (String) Storage Pool Name



<a id="nestedatt--sdt_list"></a>
### Nested Schema for `sdt_list`

Read-Only:

- `all_ips` (String) All IP
- `discovery_port` (Number) Discovery Port
- `host_only_ips` (String) Host Only IP
- `id` (String) ID
- `ip` (String) SDT Node IP
- `name` (String) Name
- `nvme_port` (Number) NVMe Port
- `protection_domain_id` (String) Protection Domain Name
- `protection_domain_name` (String) Protection Domain Name
- `storage_only_ips` (String) Storage Only IP
- `storage_port` (Number) Storage Port

## Import

Import is supported using the following syntax:

```shell
# /*
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import existing cluster
terraform import powerflex_cluster.cluster_data "<MDM_IP>,<MDM_Password>,<LIA_Password>"
```

1. This will import the cluster instance using specified details into your Terraform state.
2. After successful import, you can run terraform state list to ensure the resource has been imported successfully.
3. Now, you can fill in the resource block with the appropriate arguments and settings that match the imported resource's real-world configuration.
4. Execute terraform plan to see if your configuration and the imported resource are in sync. Make adjustments if needed.
5. Finally, execute terraform apply to bring the resource fully under Terraform's management.
6. Now, the resource which was not part of terraform became part of Terraform managed infrastructure.