terraform {
  required_providers {
    powerflex = {
      version = "0.0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}

provider "powerflex" {
    username = var.username
    password = var.password
    endpoint = var.endpoint
    insecure = true
}

data "powerflex_storagepool" "example" {
  //protection_domain_name = "domain1"
  protection_domain_id = "4eeb304600000000"
  //storage_pool_id = ["7630a24600000000", "7630a24800000002"]
  storage_pool_name = ["pool2", "pool1"]
}
