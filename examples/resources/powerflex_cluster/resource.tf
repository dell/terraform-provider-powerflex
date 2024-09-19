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