# terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check sds_resource_import.tf for more info
# To create / update, either protection_domain_name or protection_domain_id must be provided
# name and ip_list are the required parameters to create or update
# other  atrributes like : performance_profile, port, drl_mode, rmcache_enabled, rfcache_enabled, rmcache_size_in_mb are optional 
# To check which attributes can be updated, please refer Product Guide in the documentation

resource "powerflex_sds" "create" {
  name = "demo-sds-test-01"
  protection_domain_name = "demo-sds-pd"
  ip_list = [
      {
        ip = "10.247.100.231"
        role = "sdsOnly"
      },
      {
        ip = "10.10.10.11"
        role = "sdcOnly"
      },
    ]
}

output "changed_sds" {
  value = powerflex_sds.create
}
