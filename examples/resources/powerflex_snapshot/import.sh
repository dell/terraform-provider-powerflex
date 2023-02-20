# Below are the steps to import snapshot :
# Step 1 - To import a snapshot , we need the id of that snapshot 
# Step 2 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_snapshot" "resource_block_name" {
# }
# Step 3 - execute the command: terraform import "powerflex_snapshot.resource_block_name" "id_of_the_snapshot" (resource_block_name must be taken from step 2)
# Step 4 - After successful execution of the command , check the state file


