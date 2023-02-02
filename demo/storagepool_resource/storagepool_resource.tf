
# # -----------------------------------------------------------------------------------
# # Create and Modify Storage Pool
# # -----------------------------------------------------------------------------------

resource "powerflex_storagepool" "storagepool" {
  name = "storagepool3"
  protection_domain_id = "4eeb304600000000"
  media_type = "HDD"
  use_rmcache = true
  use_rfcache = false
}

output "created_storagepool" {
  value = powerflex_storagepool.storagepool
}
# # -----------------------------------------------------------------------------------
