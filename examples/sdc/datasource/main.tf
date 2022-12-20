# cd ../../.. && make install && cd examples/sdc/datasource
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell.com/dev/powerflex"
    }
  }
}

module "base" {
  source = "../../config"
}


provider "powerflex" {
    insecure = ""
    usecerts = ""
    powerflex_version = ""
}

# # -----------------------------------------------------------------------------------
# # Read all sdcs if id is blank, otherwise reads all sdcs
# # -----------------------------------------------------------------------------------
    # systemid is required field
    # name is optional if empty then will return all sdc
    # sdcid is optional if empty then will return all sdc
    # sdcid and name both are empty then will return all sdc
data "powerflex_sdc" "selected" {
    sdcid = "c423b09800000003"
    systemid = "0e7a082862fedf0f"
    # name = "LGLW6090"
}

# # Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected
}
# # -----------------------------------------------------------------------------------

