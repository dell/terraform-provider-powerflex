# commands to run this tf file : terraform init && terraform apply --auto-approve
# Reads protection domain either by name or by id , if provided
# If both name and id is not provided , then it reads all the protection domain
# id and name can't be given together to fetch the protection domain .

data "powerflex_protection_domain" "pd" {
  name = "domain1"
  # id = "202a046600000000"
}

output "inputPdID" {
  value = data.powerflex_protection_domain.pd.id
}

output "inputPdName" {
  value = data.powerflex_protection_domain.pd.name
}

output "pdResult" {
  value = data.powerflex_protection_domain.pd.protection_domains
}

