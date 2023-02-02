data "powerflex_sds" "example2" {
  # require field is either of protection_domain_name or protection_domain_id
	protection_domain_name = "domain1"
  # protection_domain_id = "XXXYIO123456"
  sds_names = [ "SDS_01_MOD" ]
  # sds_ids = ["XXXY01230456"]
}

output "allsdcresult" {
  value = data.powerflex_sds.example2
}

