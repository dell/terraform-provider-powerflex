
resource "powerflex_protection_domain" "pd" {
  # required paramaters  ======

  name = "domain_1"

  # optional parameters  ======

  active = true

  # SDS IOPS throttling
  # overall_io_network_throttling_in_kbps must be greater than the rest of the parameters
  # 0 indicates unlimited IOPS
  protected_maintenance_mode_network_throttling_in_kbps = 10 * 1024
  rebuild_network_throttling_in_kbps                    = 10 * 1024
  rebalance_network_throttling_in_kbps                  = 10 * 1024
  vtree_migration_network_throttling_in_kbps            = 10 * 1024
  overall_io_network_throttling_in_kbps                 = 20 * 1024

  # Fine granularity metadata caching
  fgl_metadata_cache_enabled      = true
  fgl_default_metadata_cache_size = 1024

  # Read Flash cache
  rf_cache_enabled          = true
  rf_cache_operational_mode = "ReadAndWrite"
  rf_cache_page_size_kb     = 16
  rf_cache_max_io_size_kb   = 32
}
