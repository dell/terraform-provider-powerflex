# Below are the steps to import volume :
# Step 1 - To import a volume , we need the id of that volume 
# Step 2 - To check the id of the volume we can make use of volume datasource . Please refer volume_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_volume" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_volume.resource_block_name" "id_of_the_volume" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file


