# https://github.com/hashicorp/terraform-provider-hashicups/blob/main/hashicups/provider.go
# https://developer.hashicorp.com/terraform/tutorials/providers/provider-debug
# https://developer.hashicorp.com/terraform/tutorials/providers/provider-complex-read
# https://github.com/hashicorp/terraform-provider-hashicups/blob/main/examples/main.tf

# cd .. && make install && cd example
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
    insecure = "true"
    usecerts = "true"
    pfxm_version = ""
}

# -----------------------------------------------------------------------------------
# Read all sdcs if id is blank, otherwise reads all sdcs
# -----------------------------------------------------------------------------------
data "powerflex_sdcs" "selected" {
    id = "595a0bb600000006"
    systemid = ""
}

# Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdcs.selected.sdcs
}
# -----------------------------------------------------------------------------------



# -----------------------------------------------------------------------------------
# Change name of sdc and read that sdc
# -----------------------------------------------------------------------------------
# data "powerflex_sdc_name_change" "selected" {
#     id = "595a0bb600000006"
#     name = "boomboomshakalaka"
#     systemid = "bae9b21d5a53850f"
# }

# # Returns changed sdcs
# output "allsdcresult" {
#   value = data.powerflex_sdc_name_change.selected.sdcs
# }
# -----------------------------------------------------------------------------------
