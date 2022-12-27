# cd ../../.. && make install && cd examples
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell/powerflex"
    }
  }
}

provider "powerflex" {
    username = var.username
    password = var.password
    endpoint = var.endpoint
    insecure = true
}

# # -----------------------------------------------------------------------------------
# # Create and Modify Storage Pool
# # -----------------------------------------------------------------------------------

resource "powerflex_storagepool" "storagepool" {
  name = "Storage_91"
  protection_domain_id = "4eeb304600000000"
  media_type = "HDD"
}

output "created_storagepool" {
  value = powerflex_storagepool.storagepool
}
# # -----------------------------------------------------------------------------------
