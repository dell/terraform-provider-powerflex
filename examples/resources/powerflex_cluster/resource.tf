terraform {
  required_providers {
    powerflex = {
      version = "1.5.0"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}
provider "powerflex" {
  username = "admin"
  password = "Password123@"
  endpoint = "https://pflex4env12.pie.lab.emc.com"
  insecure = true
  timeout  = 120
}


resource "powerflex_package" "upload-test" {
  file_path = ["/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-lia-4.5-0.287.sles15.3.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-mdm-4.5-0.287.sles15.3.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-sdc-4.5-0.287.sles15.3.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-sds-4.5-0.287.sles15.3.x86_64.rpm",
    "/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-sdr-4.5-0.287.sles15.3.x86_64.rpm",
  "/root/powerflex_packages/PowerFlex_4.5.0.287_SLES15.3/EMC-ScaleIO-activemq-5.16.4-62.noarch.rpm"]
}

resource "powerflex_cluster" "test" {
  depends_on   = [powerflex_package.upload-test]
  mdm_password = "Password123@"
  lia_password = "Password123@"
  # Advance Security Configuration
  allow_non_secure_communication_with_lia = false
  allow_non_secure_communication_with_mdm = false
  disable_non_mgmt_components_auth        = false
  # Cluster Configuration related fields
  cluster = [
    {
      # MDM Configuration Fields
      username         = "root",
      password         = "admin",
      operating_system = "linux",
      is_mdm_or_tb     = "Primary",
      ips          = "10.247.103.160" #,10.247.39.124,10.247.78.174", 
      # mdm_mgmt_ip              = "10.247.103.160",
      mdm_name             = "mdm1",
      perf_profile_for_mdm = "HighPerformance",
      # virtual_ips              = "10.247.39.126,10.247.78.177",
      # virtual_ip_nics          = "eth1,eth2",
      is_sds = "Yes",
      #   sds_name                 = "sds1",
      #   sds_all_ips              = "10.247.39.124,10.247.78.174",
      #   protection_domain        = "domain_1",
      #  sds_storage_device_list  = "/dev/sdb",
      #  sds_storage_device_names = "sdb",
      #  storage_pool_list        = "pool1",
      #  perf_profile_for_sds     = "HighPerformance",
      is_sdc = "No"
    },
    {
      username         = "root",
      password         = "admin",
      operating_system = "linux",
      is_mdm_or_tb     = "Secondary",
      mdm_name         = "mdm2",
      ips          = "10.247.103.161" #,10.247.39.122,10.247.78.175", 
      # mdm_mgmt_ip              = "10.247.103.161",
      perf_profile_for_mdm = "HighPerformance",
      # virtual_ips              = "10.247.39.126,10.247.78.177",
      # virtual_ip_nics          = "eth1,eth2",
      is_sds = "Yes",
      #  sds_name                 = "sds2",
      # sds_all_ips              = "10.247.39.122,10.247.78.175",
      # sds_storage_device_list  = "/dev/sdb",
      # sds_storage_device_names = "sdb",
      # protection_domain        = "domain_1",
      # storage_pool_list        = "pool1",
      # perf_profile_for_sds     = "HighPerformance",
      is_sdc = "No"
    },
    {
      username         = "root",
      password         = "admin",
      operating_system = "linux",
      is_mdm_or_tb     = "TB",
      mdm_name         = "tb1",
      ips              = "10.247.103.162" #,10.247.39.130,10.247.78.176",
      # mdm_mgmt_ip              = "10.247.103.162",
      perf_profile_for_mdm = "HighPerformance",
      is_sds               = "Yes",
      # sds_name                 = "sds3",
      # sds_all_ips              = "10.247.39.130,10.247.78.176",
      # sds_storage_device_list  = "/dev/sdb",
      # sds_storage_device_names = "sdb",
      # protection_domain        = "domain_1",
      # storage_pool_list        = "pool1",
      # perf_profile_for_sds     = "HighPerformance",
      is_sdc = "No"
    }
  ]
  storage_pools = [
    {
      media_type = "HDD"
      # protection_domain = "domain_1"
      # storage_pool      = "pool1"
      # daya_layout       = "MG"
      # zero_padding      = "true"
    }
  ]
}
