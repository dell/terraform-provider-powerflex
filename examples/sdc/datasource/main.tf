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



module "base" {
  source = "../../config"
}

provider "powerflex" {
    username = "${module.base.username}"
    password = "${module.base.password}"
    endpoint = "${module.base.host}"
}

# # -----------------------------------------------------------------------------------
# # Read all sdcs if id is blank, otherwise reads all sdcs
# # -----------------------------------------------------------------------------------
    # name is optional if empty then will return all sdc
    # sdcid is optional if empty then will return all sdc
    # sdcid and name both are empty then will return all sdc
data "powerflex_sdc" "selected" {
    # id = "c423b09800000003"
    # name = ""
}

# # Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected
}
# # -----------------------------------------------------------------------------------

