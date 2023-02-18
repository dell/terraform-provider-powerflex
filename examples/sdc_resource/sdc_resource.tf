# Commands to run this tf file : terraform init && terraform plan && terraform apply
# Only importing the sdc resource or renaming of sdc resource is supported
# For Renaming , id and name are required fields
# For importing , please check sdc_resource_import.tf file for more details
# name can't be empty


resource "powerflex_sdc" "sdc" {
  id   = "c423b09800000003"
  name = "powerflex_sdc26"
}


#output "changed_sdc" {
# value = powerflex_sdc.sdc
#}
# # -----------------------------------------------------------------------------------


