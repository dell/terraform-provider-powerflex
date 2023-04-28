# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Read, Delete and Import operations are supported for this resource.
# To add device, device_path is mandatory along with storage pool and SDS identifier.
# Along with storage_pool_name, we have to specify protection_domain_id or protection_domain_name.
# To check which attributes of the device resource can be updated, please refer Product Guide in the documentation

resource "powerflex_device" "test-device" {
  device_path = "/dev/sdc"
  storage_pool_name = "pool1"
  protection_domain_name = "domain1"
  sds_name = "SDS_2"
  media_type = "HDD"
}
