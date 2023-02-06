# Below are the steps to import snapshot :
# Step 1 - To import a snapshot , we need the id of that snapshot 
# Step 2 - to get the id of the snapshot , use volume datasource and pass the name of the snapshot to get snapshot id
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_snapshot" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_snapshot.resource_block_name" "id_of_the_snapshot" (resource_block_name must be taken from step 3 and id can be taken from Step 2)
# Step 5 - After successful execution of the command , check the state file


