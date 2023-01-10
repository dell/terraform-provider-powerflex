data "powerflex_storagepool" "example" {
  //protection_domain_name = "domain1"
  protection_domain_id = "4eeb304600000000"
  //storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
  storage_pool_names = ["pool2", "pool1"]
}
