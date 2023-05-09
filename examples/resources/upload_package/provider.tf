terraform {
  required_providers {
    powerflex = {
      version = "0.0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}
provider "powerflex" {
  username = "admin"
  password = "Password123"
  endpoint = "https://10.247.103.159:443"
  insecure = true
  timeout  = 120
  gatewayinstaller = true
}