# Below are the steps to import a user :
# Step 1 - To import a user , we either need id or  username of that user 
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_user" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_user.resource_block_name" "name:userName" (where userName is the name of the user which you wanna import)
# 2nd approach to import the user is using the id. The import command to import using the id is : terraform import "powerflex_user.resource_block_name" "id:userId" OR terraform import "powerflex_user.resource_block_name" "userId"
# Step 5 - After successful execution of the command , check the state file


