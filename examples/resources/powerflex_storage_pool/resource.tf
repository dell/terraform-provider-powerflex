# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check storagepool_resource_import.tf for more info
# To create / update, either protection_domain_id or protection_domain_name must be provided
# name and media_type is the required parameter to create or update
# other  atrributes like : use_rmcache, use_rfcache, replication_journal_capacity, capacity_alert_high_threshold, capacity_alert_critical_threshold etc. are optional 
# To check which attributes of the storage pool can be updated, please refer Product Guide in the documentation

resource "powerflex_storage_pool" "sp" {
  name = "storagepool3"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type             = "HDD"
  use_rmcache            = false
  use_rfcache            = true
  #replication_journal_capacity = 34
  capacity_alert_high_threshold                               = 66
  capacity_alert_critical_threshold                           = 77
  zero_padding_enabled                                        = false
  protected_maintenance_mode_io_priority_policy               = "favorAppIos"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps      = 1028
  rebalance_enabled                                           = false
  rebalance_io_priority_policy                                = "favorAppIos"
  rebalance_num_of_concurrent_ios_per_device                  = 7
  rebalance_bw_limit_per_device_in_kbps                       = 1032
  vtree_migration_io_priority_policy                          = "favorAppIos"
  vtree_migration_num_of_concurrent_ios_per_device            = 7
  vtree_migration_bw_limit_per_device_in_kbps                 = 1030
  spare_percentage                                            = 66
  rm_cache_write_handling_mode                                = "Passthrough"
  rebuild_enabled                                             = true
  rebuild_rebalance_parallelism                               = 5
  fragmentation                                               = false
}

output "created_storagepool" {
  value = powerflex_storage_pool.sp
}
