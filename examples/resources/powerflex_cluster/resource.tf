/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# To perform Cluster Installation
resource "powerflex_cluster" "test" {
	mdm_password =  "Password"
	lia_password= "Password"
	allow_non_secure_communication_with_lia= false
	allow_non_secure_communication_with_mdm= false
	disable_non_mgmt_components_auth= false
	cluster = [
	{
		ips= "10.10.10.1",
		username= "root",
		password= "Password",
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		protection_domain = "domain_1"
		sds_storage_device_list = "/dev/sdb"
		storage_pool_list = "pool1"
		perf_profile_for_sdc= "High",
		ia_rfcache= "No",
		is_sdr= "Yes",
		sdr_all_ips = "10.10.20.1"
	 },
	 {
		ips= "10.10.10.2",
		username= "root",
		password= "Password",
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		protection_domain = "domain_1"
		sds_storage_device_list = "/dev/sdb"
		storage_pool_list = "pool1"
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "10.10.10.3",
		username= "root",
		password= "Password",
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		is_sds= "No",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
			protection_domain = "domain_1"
			storage_pool = "pool1"
			replication_journal_capacity_percentage = "50"
		}	
	]
}

# To perform Cluster Import MDM_IP, MDM_Password, LIA_Password
# terraform import powerflex_cluster.test "10.10.10.1,Password,Password"
# To Delete whole cluster
# terraform destroy
resource "powerflex_cluster" "test" {
}