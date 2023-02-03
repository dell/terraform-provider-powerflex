# Below are the steps to import storage pool :
# Step 1 - To import a storage pool , we need the id of that storage pool 
# Step 2 - To check the id of the storage pool we can make use of storage pool datasource . Please refer storagepool_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_storagepool" "var_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_storagepool.var_name" "id_of_the_storagepool" (var_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file


