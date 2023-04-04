# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check storagepool_resource_import.tf for more info
# To create / update, either protection_domain_id or protection_domain_name must be provided
# name and media_type is the required parameter to create or update
# other  atrributes like : use_rmcache, use_rfcache are optional 
# To check which attributes of the storage pool can be updated, please refer Product Guide in the documentation

resource "powerflex_storage_pool" "sp" {
  name                 = "storagepool3"
  protection_domain_id = "4eeb304600000000"
  # protection_domain_name = "domain1"
  media_type  = "HDD"
  use_rmcache = true
  use_rfcache = false
}

output "created_storagepool" {
  value = powerflex_storage_pool.sp
}
