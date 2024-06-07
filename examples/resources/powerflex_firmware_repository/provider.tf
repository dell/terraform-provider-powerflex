terraform {
  required_providers {
    powerflex = {
      version = "1.5.0"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}
provider "powerflex" {
  username = "admin"
  password = "Password123@"
  endpoint = "https://pflex4env12.pie.lab.emc.com"
  insecure = true
  timeout  = 120
}