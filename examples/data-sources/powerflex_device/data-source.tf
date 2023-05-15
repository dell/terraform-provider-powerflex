# commands to run this tf file : terraform init && terraform apply --auto-approve
# empty block of the powerflex_device datasource will give list of all device within the system

data "powerflex_device" "dev" {
}

data "powerflex_device" "dev1" {
  storage_pool_id = "c98e26e500000000"
}

data "powerflex_device" "dev2" {
  sds_name = "SDS_2"
}

output "deviceResult" {
  value = data.powerflex_device.dev.device_model
}
