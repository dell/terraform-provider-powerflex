# cd ../../ && make install && cd examples/volume
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
    username = var.powerflex_username
    password = var.powerflex_password
    host = var.powerflex_host
    insecure = "true"
    usecerts = "false"
    # powerflex_version = ""
}

# # -----------------------------------------------------------------------------------
# # Read all volumes if completely blank, otherwise reads specific volume based on id or name
# # -----------------------------------------------------------------------------------
    # name is optional , if mentioned then will retrieve the specific volume with that name
    # id is optional , if mentioned then will retrieve the specific volume with that id
    # id and name both are empty then will return all sdc

data "powerflex_protection_domain" "pd" {

    # name = "domain1"
    # id = "4eeb304600000000"
}

output "inputPdID" {
  value = data.powerflex_protection_domain.pd.id
}

output "inputPdName" {
  value = data.powerflex_protection_domain.pd.name
}

output "pdResult" {
  value = data.powerflex_protection_domain.pd.protection_domains
}

