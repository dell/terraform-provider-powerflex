# terraform init && terraform apply --auto-approve
# Read all protection domains
# Filters by name or id if provided, but not both
data "powerflex_protection_domain" "pd" {
    name = "domain1"
    # id = "4eeb304600000000"
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

