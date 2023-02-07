# cd ../../.. && make install && cd examples/sdc/resource
# terraform init && terraform import powerflex_sdc.sdc c423b09800000003
# terraform init && terraform apply --auto-approve
# terraform destroy
terraform {
  required_providers {
    powerflex = {
      version = "0.0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}

# module "base" {
#   source = "../../config"
# }

provider "powerflex" {
    username = "admin"
    password = "Password123"
    endpoint = "https://10.247.101.69:443"
}



# # -----------------------------------------------------------------------------------
# # Change name of sdc and read that sdc
# # -----------------------------------------------------------------------------------
resource "powerflex_sdc" "sdc" {
  id = "c42a193500000096"
  name = " "
}

output "changed_sdc" {
  value = powerflex_sdc.sdc
}
# # -----------------------------------------------------------------------------------
