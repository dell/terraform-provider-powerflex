terraform {
  required_providers {
    powerflex = {
      version = "0.0.1"
      source  = "registry.terraform.io/dell/powerflex"
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
      {
        ip = "10.247.100.232"
        role = "sdsOnly"
      },
      {
        ip = "10.10.10.1"
        role = "sdcOnly"
      },
      {
        ip = "10.10.10.2"
        role = "sdcOnly"
      }
    ]
  protection_domain_id = "4eeb304600000000"
}

output "changed_sds" {
  value = powerflex_sds.create
}
