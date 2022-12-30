# cd ../../.. && make install && cd examples/sdc/resource
# terraform init && terraform import powerflex_sdc.sdc c423b09800000003
# terraform init && terraform apply --auto-approve
# terraform destroy
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}

module "base" {
  source = "../../config"
}

provider "powerflex" {
    username = "${module.base.username}"
    password = "${module.base.password}"
    endpoint = "${module.base.host}"
}



# # -----------------------------------------------------------------------------------
# # Change name of sdc and read that sdc
# # -----------------------------------------------------------------------------------
resource "powerflex_sdc" "sdc" {
  sdc_id = "c423b09800000003"
  name = "powerflex_sdc25"
}

output "changed_sdc" {
  value = powerflex_sdc.sdc
}
# # -----------------------------------------------------------------------------------
