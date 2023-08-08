# Below are the steps to import Cluster :
# Step 1 - To import a Cluster , we need the MDM IP, MDM Password, LIA Password  of that Cluster 
# Step 2 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_cluster" "resource_block_name" {
# }
# Step 4 - execute the command : terraform init && terraform import "powerflex_cluster.resource_block_name" "id_of_the_sdc" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file
