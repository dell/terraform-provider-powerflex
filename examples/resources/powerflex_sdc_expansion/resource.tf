# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
resource "powerflex_sdc_expansion" "test-csv2" {
	mdm_ip = "MDM_IPS"
	mdm_password = "ABCD"
	lia_password="ABCD"
	csv_detail = [
		{
			ip = "IP"
			password = "Password"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "IP"
			password = "Password"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "IP"
			password = "Password"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "IP"
			password = "Password"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "Yes"
			performance_profile = "Compact"
			sdc_name = "SDC"
   		},
	]
}