# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check storagepool_resource_import.tf for more info
# To create / update, either protection_domain_id or protection_domain_name must be provided
# name and media_type is the required parameter to create or update
# other  atrributes like : use_rmcache, use_rfcache are optional 
# To check which attributes of the storage pool can be updated, please refer Product Guide in the documentation

resource "powerflex_storage_pool" "sp" {
  name                 = "storagepool3"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  #replication_journal_capacity = 34
  capacity_alert_high_threshold = 66
  capacity_alert_critical_threshold = 77
  zero_padding_enabled = false
  protected_maintenance_mode_io_priority_policy = "favorAppIos"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  rebalance_enabled = false
  use_rmcache = false
  use_rfcache = true
  rebalance_io_priority_policy = "favorAppIos"
  rebalance_num_of_concurrent_ios_per_device = 7
  rebalance_bw_limit_per_device_in_kbps = 1032
  vtree_migration_io_priority_policy = "favorAppIos"
  vtree_migration_num_of_concurrent_ios_per_device = 7
  vtree_migration_bw_limit_per_device_in_kbps = 1030
  spare_percentage = 66
  rm_cache_write_handling_mode = "Passthrough"
  rebuild_enabled = true
  rebuild_rebalance_parallelism = 5
  fragmentation = false
}

resource "powerflex_storage_pool" "sp2" {
  name                 = "storagepool4"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  zero_padding_enabled = true
  capacity_alert_high_threshold = 68
  protected_maintenance_mode_io_priority_policy = "limitNumOfConcurrentIos"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 9
  rebalance_enabled = true
  use_rmcache = true
  use_rfcache = true
  rebalance_io_priority_policy = "limitNumOfConcurrentIos"
  rebalance_num_of_concurrent_ios_per_device = 8
  vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
  vtree_migration_num_of_concurrent_ios_per_device = 10
  spare_percentage = 67
  rm_cache_write_handling_mode = "Passthrough"
  rebuild_enabled = false
  rebuild_rebalance_parallelism = 6
  fragmentation = true
}

resource "powerflex_storage_pool" "sp3" {
  name                 = "storagepool5"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  zero_padding_enabled = false
  #replication_journal_capacity = 34
  capacity_alert_critical_threshold = 88
  protected_maintenance_mode_io_priority_policy = "unlimited"
  use_rmcache = true
  use_rfcache = false
  rebalance_io_priority_policy = "unlimited"
  vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
  spare_percentage = 68
  rm_cache_write_handling_mode = "Cached"
  rebuild_enabled = true
  rebuild_rebalance_parallelism = 7

}

resource "powerflex_storage_pool" "sp4" {
  name                 = "storagepool6"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
}

resource "powerflex_storage_pool" "sp5" {
  name                 = "storagepool7"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  protected_maintenance_mode_io_priority_policy = "favorAppIos"
  protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  rebalance_io_priority_policy = "favorAppIos"
  rebalance_bw_limit_per_device_in_kbps = 1032
  vtree_migration_io_priority_policy = "favorAppIos"
  vtree_migration_bw_limit_per_device_in_kbps = 1030
  
}

resource "powerflex_storage_pool" "sp6" {
  name                 = "storagepool8"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  rebalance_num_of_concurrent_ios_per_device = 7
  rebalance_bw_limit_per_device_in_kbps = 1032
  vtree_migration_num_of_concurrent_ios_per_device = 7
  vtree_migration_bw_limit_per_device_in_kbps = 1030
}

resource "powerflex_storage_pool" "sp7" {
  name                 = "storagepool9"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  rebalance_num_of_concurrent_ios_per_device = 7
  vtree_migration_num_of_concurrent_ios_per_device = 7
}

resource "powerflex_storage_pool" "sp8" {
  name                 = "storagepool10"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  rebalance_bw_limit_per_device_in_kbps = 1032
  vtree_migration_bw_limit_per_device_in_kbps = 1030
}


resource "powerflex_storage_pool" "sp9" {
  name                 = "storagepool11"
  #protection_domain_id = "4eeb304600000000"
  protection_domain_name = "domain1"
  media_type  = "HDD"
  protected_maintenance_mode_io_priority_policy = "limitNumOfConcurrentIos"
  protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
  protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  use_rmcache = false
  rebalance_io_priority_policy = "limitNumOfConcurrentIos"
  rebalance_num_of_concurrent_ios_per_device = 7
  rebalance_bw_limit_per_device_in_kbps = 1032
  vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
  vtree_migration_num_of_concurrent_ios_per_device = 7
  vtree_migration_bw_limit_per_device_in_kbps = 1030
  rm_cache_write_handling_mode = "Passthrough"
}

output "created_storagepool" {
  value = powerflex_storage_pool.sp
}
