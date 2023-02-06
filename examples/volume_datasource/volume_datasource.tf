# commands to run this tf file : terraform init && terraform apply --auto-approve
# This datasource reads volumes either by id or name or storage_pool_id or storage_pool_name where user can provide a value to any one of them
# If it is a empty datsource block , then it will read all the volumes
# If id or name is provided then it reads a particular volume with that id or name
# If storage_pool_id or storage_pool_name is provided then it will return the volumes under that storage pool
# Only one of the attribute can be provided among id, name, storage_pool_id, storage_pool_name 

data "powerflex_volume" "volume" {

    #name = "cosu-ce5b8a2c48"
    id = "4570761d00000024"
    #storage_pool_id= "7630a24800000002"
    #storage_pool_name= "pool2"
}

output "volumeResult" {
  value = data.powerflex_volume.volume.volumes
}

