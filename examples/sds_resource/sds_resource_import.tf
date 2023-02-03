# Below are the steps to import SDS :
# Step 1 - To import a sds , we need the id of that sds 
# Step 2 - To check the id of the sds we can make use of sds datasource . Please refer sds_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_sds" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_sds.resource_block_name" "id_of_the_sds" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file


