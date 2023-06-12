# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
resource "powerflex_sdc" "test" {
  id   = "SDC_ID"
  name = "Name"
  mdm_password = "Password"
  lia_password = "Password"
  cluster_details = [
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Primary"
      is_sdc           = "No"
      name             = "SDC_NAME"
      performance_profile ="HighPerformance"
      sdc_id           = "sdc_id"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Secondary"
      is_sdc           = "Yes"
      name             = "SDC_NAME"
      performance_profile ="Compact"
      sdc_id           = "sdc_id"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "TB"
      is_sdc           = "Yes"
      name             = "SDC_NAME"
      performance_profile ="Compact"
      sdc_id           = "sdc_id"
    },
    {
      ip               = "IP"
      username         = "Username"
      password         = "Password"
      operating_system = "linux"
      is_mdm_or_tb     = "Standby"
      is_sdc           = "Yes"
      name             = "SDC_NAME"
      performance_profile ="Compact"
      sdc_id           = "sdc_id"
    },
  ]
}