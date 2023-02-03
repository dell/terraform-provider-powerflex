# commands to run this tf file : terraform init && terraform apply --auto-approve
# To read SDS, either protection_domain_name or protection_domain_id must be provided
# This datasource reads SDS either by sds_names or sds_ids where user can provide a list of ids or names
# if both sds_names and sds_ids are not provided , then it will read all the sds under the protection domain
# Both sds_ids and sds_names can't be provided together .
# Both protection_domain_name and protection_domain_id can't be provided together



data "powerflex_sds" "example2" {
  # require field is either of protection_domain_name or protection_domain_id
	protection_domain_name = "domain1"
  # protection_domain_id = "4eeb304600000000"
  sds_names = [ "SDS_01_MOD" ]
  # sds_ids = ["6adfec1000000000"]
}

output "allsdcresult" {
  value = data.powerflex_sds.example2
}

