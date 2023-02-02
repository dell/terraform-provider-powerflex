
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
      {
        ip = "10.10.10.12"
        role = "sdcOnly"
      }
    ]
}

output "changed_sds" {
  value = powerflex_sds.create
}
