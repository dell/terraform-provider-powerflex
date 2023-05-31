# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
resource "powerflex_sdc_expansion" "test" {
  mdm_password = "Password"
  lia_password = "Password"
  cluster_details = [
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Primary"
      is_sdc           = "Yes"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Secondary"
      is_sdc           = "Yes"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "TB"
      is_sdc           = "Yes"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Standby"
      is_sdc           = "Yes"
    },
  ]
}
