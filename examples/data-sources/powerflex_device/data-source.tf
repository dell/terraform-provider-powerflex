# commands to run this tf file : terraform init && terraform apply --auto-approve
# empty block of the powerflex_device datasource will give list of all device within the system

data "powerflex_device" "dev" {
}

output "deviceResult" {
  value = data.powerflex_device.dev.device_model
}
