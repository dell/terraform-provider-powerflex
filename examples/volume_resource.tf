terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell.com/dev/powerflex"
    }
  }
}
provider "powerflex" {
    username = "celestials"
    password = ""
    host = ""
    insecure = ""
    usecerts = ""
    powerflex_version = ""
}
resource "powerflex_volume" "avengers" {
  name = "avengers-ironman"
  storage_pool_id = "76"
  protection_domain_id = "4e"
  capacity_unit = "GB"
  size = 10
  map_sdcs_id = ["c4","c3"]
}