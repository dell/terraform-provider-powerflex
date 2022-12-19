terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell.com/dev/powerflex"
    }
  }
}

provider "powerflex" {
    username = ""
    password = ""
    host = ""
    insecure = ""
    usecerts = ""
    powerflex_version = ""
}

resource "powerflex_volume" "avengers" {
  name = "avengers-powerstone"
  storage_pool_id = ""
  protection_domain_id = ""
  volume_size_in_kb = ""
}