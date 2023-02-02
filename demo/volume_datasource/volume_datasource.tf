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

