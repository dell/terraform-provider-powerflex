# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check volume_resource_import.tf for more info
# To create / update, either storage_pool_id or storage_pool_name must be provided
# Also , to create / update, either protection_domain_id or protection_domain_name must be provided
# name, size is the required parameter to create or update
# other  atrributes like : capacity_unit, volume_type, use_rm_cache, compression_method, access_mode, remove_mode, sdc_list are optional 
# To check which attributes of the snapshot can be updated, please refer Product Guide in the documentation

resource "powerflex_sdc_expansion" "test-csv2" {
	mdm_ip = "10.247.103.160"
	mdm_password = "Password123"
	lia_password="Password123"
	csv_detail = [
		{
			ip = "10.247.103.160"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.161"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.162"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "10.247.103.163"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "Yes"
   		},
	]
}