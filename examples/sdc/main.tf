# cd ../../ && make install && cd examples/sdc
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell.com/dev/powerflex"
    }
  }
}

provider "powerflex" {
    username = "admin"
    password = "Scaleio123!"
    host = "https://10.234.221.217"
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
# data "powerflex_sdc" "selected" {

    
#     # sdcid = "c42844760000005f"
#     systemid = "bae9b21d5a53850f"

#     # name = "LGLW6091" // LGLW6091
# }

# # Returns all sdcs matching criteria
# output "allsdcresult" {
#   value = data.powerflex_sdc.selected.sdcs
# }
# # -----------------------------------------------------------------------------------



# # -----------------------------------------------------------------------------------
# # Change name of sdc and read that sdc
# # -----------------------------------------------------------------------------------
resource "powerflex_sdc" "sdc" {
  id = "595a0bb300000003"
  name = "ffror046"
  systemid = "bae9b21d5a53850f"
}

output "changed_sdc" {
  value = powerflex_sdc.sdc
}
# # -----------------------------------------------------------------------------------
