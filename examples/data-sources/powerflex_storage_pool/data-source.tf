# commands to run this tf file : terraform init && terraform apply --auto-approve
# To read storage pool, either protection_domain_name or protection_domain_id must be provided
# This datasource reads a list of storage pools either by storage_pool_ids or storage_pool_names where user can provide a list of ids or names
# if both storage_pool_ids and storage_pool_names are not provided , then it will read all the storage pool under the protection domain
# Both storage_pool_ids and storage_pool_names can't be provided together .
# Both protection_domain_name and protection_domain_id can't be provided together


data "powerflex_storage_pool" "example" {
  //protection_domain_name = "domain1"
  protection_domain_id = "4eeb304600000000"
  //storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
  storage_pool_names = ["pool2", "pool1"]
}


output "allsdcresult" {
  value = data.powerflex_storage_pool.example.storage_pools
}
