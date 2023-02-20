# commands to run this tf file : terraform init && terraform apply --auto-approve
# Reads SDC either by name or by id , if provided
# If both name and id is not provided , then it reads all the SDC
# id and name can't be given together to fetch the SDC .
# id can't be empty

data "powerflex_sdc" "selected" {
  #id = "c423b09800000003"
  name = "sdc_01"
}

# # Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected
}
# # -----------------------------------------------------------------------------------

