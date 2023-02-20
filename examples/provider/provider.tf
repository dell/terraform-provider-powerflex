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
  timeout  = 120
}