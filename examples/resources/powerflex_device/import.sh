# Below are the steps to import device :
# Step 1 - To import a device , we need the id of that device
# Step 2 - To check the id of the device we can make use of device datasource . Please refer device_datasource.tf for more info.
# Step 3 - create a tf file with empty resource block . Refer the example below.
# Example :
# resource "powerflex_device" "resource_block_name" {
# }
# Step 4 - execute the command: terraform import "powerflex_device.resource_block_name" "id_of_the_device" (resource_block_name must be taken from step 3 and id must be taken from step 2)
# Step 5 - After successful execution of the command , check the state file


