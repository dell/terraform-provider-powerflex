# cd ../ && make install && cd examples
# terraform init && terraform apply --auto-approve
terraform {
  required_providers {
    powerflex = {
      version = "0.1"
      source  = "registry.terraform.io/dell/powerflex"
    }
  }
}

provider "powerflex" {
    username = ""
    password = ""
    endpoint = ""
    #insecure = true
}

# # -----------------------------------------------------------------------------------
# # Read all volumes if completely blank, otherwise reads specific volume based on id or name
# # -----------------------------------------------------------------------------------
    # name is optional , if mentioned then will retrieve the specific volume with that name
    # id is optional , if mentioned then will retrieve the specific volume with that id
    # id and name both are empty then will return all sdc

data "powerflex_volume" "volume" {

    #name = "cosu-ce5b8a2c48"
    id = "4570761d00000024"
    #storage_pool_id= "7630a24800000002"
    #storage_pool_name= "pool2"
}

output "volumeResult" {
  value = data.powerflex_volume.volume.volumes
}

