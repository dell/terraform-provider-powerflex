# cd ../../ && make install && cd examples/volume
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
    username = var.powerflex_username
    password = var.powerflex_password
    host = var.powerflex_host
}

# Read all protection domains
# Filters by name or id if provided, but not both
data "powerflex_protection_domain" "pd" {
    name = "domain1"
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

