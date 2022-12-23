# cd ../../.. && make install && cd examples/sdc/resource
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "dell.com/dev/powerflex"
    }
  }
}

module "base" {
  source = "../../config"
}

provider "powerflex" {
    username = "${module.base.username}"
    password = "${module.base.password}"
    host = "${module.base.host}"
    insecure = ""
    usecerts = ""
    powerflex_version = ""
}



# # -----------------------------------------------------------------------------------
# # Change name of sdc and read that sdc
# # -----------------------------------------------------------------------------------
resource "powerflex_sdc" "sdc" {
  sdcid = "c423b09800000003"
  systemid = "0e7a082862fedf0f"
  name = "powerflex_sdc11"
}

output "changed_sdc" {
  value = powerflex_sdc.sdc
}
# # -----------------------------------------------------------------------------------
