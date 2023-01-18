resource "powerflex_volume" "avengers" {
  name = "<volume-name>"
  storage_pool_id = "<storage-pool-id>"
  protection_domain_id = "<protection-domain-id>"
  capacity_unit = "<capacity-unit> - GB/TB "
  size = 10
  map_sdcs_id = ["<sdc-id-1>","<sdc-id-2>"]
}
