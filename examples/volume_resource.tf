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
    host = var.host
}
resource "powerflex_volume" "avengers" {
  name = var.volume_resource_name
  storage_pool_id = var.volume_resource_storage_pool_id
  protection_domain_id = var.volume_resource_storage_pool_id
  capacity_unit = var.volume_resource_capacity_unit
  size = var.volume_resource_size
  map_sdcs_id = var.volume_resource_map_sdcs_id
}
