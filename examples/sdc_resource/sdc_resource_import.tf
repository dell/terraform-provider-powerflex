# Below are the steps to import SDC :
# Step 1 - To import a sdc , we need the id of that sdc 
# Step 2 - To check the id of the sdc we can make use of sdc datasource . Please refer sdc_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_sdc" "resource_block_name" {
# }
# Step 4 - execute the command : terraform init && terraform import "powerflex_sdc.resource_block_name" "id_of_the_sdc" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file


