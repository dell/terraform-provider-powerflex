# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
resource "powerflex_sdc_expansion" "test-csv2" {
	mdm_ip = "10.247.103.160"
	mdm_password = "ABCD"
	lia_password="ABCD"
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