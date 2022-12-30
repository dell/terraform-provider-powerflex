# cd ../../.. && make install && cd examples/sdc/datasource
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}



provider "powerflex" {
    username = ""
password = ""
endpoint = ""
insecure = true
}

# # -----------------------------------------------------------------------------------
# # Read all sdcs if id is blank, otherwise reads all sdcs
# # -----------------------------------------------------------------------------------
    # name is optional if empty then will return all sdc
    # sdcid is optional if empty then will return all sdc
    # sdcid and name both are empty then will return all sdc
data "powerflex_sdc" "selected" {
    # sdc_id = "c423b09800000003"
    # name = ""
}

# # Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected
}
# # -----------------------------------------------------------------------------------

