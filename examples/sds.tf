terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell/powerflex"
    }
  }
}

provider "powerflex" {
    username = var.powerflex_username
    password = var.powerflex_password
    endpoint = var.powerflex_endpoint
}

resource "powerflex_sds" "create" {
  name = "SDS_01"
  ip_list = [
      "10.247.101.60"
    ]
  protection_domain_id = "4eeb304600000000"
}

output "changed_sds" {
  value = powerflex_sds.create
}
