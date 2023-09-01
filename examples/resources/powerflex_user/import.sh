# Below are the steps to import a user :
# Step 1 - To import a user , we need the unsername of that user 
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_user" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_user.resource_block_name" "name:user_name" (where user_name is the name of the user which you wanna import)
# Step 5 - After successful execution of the command , check the state file


