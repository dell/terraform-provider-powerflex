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
    password = ""
    host = ""
    insecure = ""
    usecerts = ""
    powerflex_version = ""
}

# # -----------------------------------------------------------------------------------
# # Read all sdcs if id is blank, otherwise reads all sdcs
# # -----------------------------------------------------------------------------------
data "powerflex_sdc" "selected" {
    sdcid = "595a0bb600000006"
    systemid = "bae9b21d5a53850f"
}

# Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected.sdcs
}
# # -----------------------------------------------------------------------------------



# # -----------------------------------------------------------------------------------
# # Change name of sdc and read that sdc
# # -----------------------------------------------------------------------------------
# # resource "powerflex_sdc" "sdc" {
# #   sdcid = "595a0bb600000006"
# #   name = "goodestname"
# #   systemid = "bae9b21d5a53850f"
# # }

# # output "changed_sdc" {
# #   value = powerflex_sdc.sdc
# #   sensitive   = true
# # }
# # -----------------------------------------------------------------------------------
