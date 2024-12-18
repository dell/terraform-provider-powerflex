/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

# commands to run this tf file : terraform init && terraform apply --auto-approve

# Get all Resource Group details present in the PowerFlex
data "powerflex_resource_group" "example1" {
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/
data "powerflex_resource_group" "resource_group" {
  # filter{
	# id = ["id1","id2"]
	# deployment_name = ["deployment_name1","deployment_name2"]
	# deployment_description = ["deployment_description1","deployment_description2"]
	# retry = true
	# teardown = false
	# teardown_after_cancel = true
	# remove_service = true
	# created_date = ["created_date1","created_date2"]
	# created_by = ["created_by1","created_by2"]
	# updated_date = ["updated_date1","updated_date2"]
	# updated_by = ["updated_by1","updated_by2"]
	# deployment_scheduled_date = ["deployment_scheduled_date1","deployment_scheduled_date2"]
	# deployment_started_date = ["deployment_started_date1","deployment_started_date2"]
    # deployment_finished_date = ["deployment_finished_date1","deployment_finished_date2"]
	# schedule_date = ["scheduled_date1","scheduled_date2"]
	# status = ["status1","status2"]
	# compliant = true
	# update_server_firmware = false
	# use_default_catalog = false
    # firmware_repository_id = ["firmware_repository_id1","firmware_repository_id2"]
	# license_repository_id = ["license_repository_id1","license_repository_id2"]
	# individual_teardown = true
	# deployment_health_status_type = ["deployment_health_status_type1","deployment_health_status_type2"]
	# all_users_allowed = true
	# owner = ["owner1","owner2"]
	# no_op = false
	# firmware_init = true
	# disruptive_firmware = false
	# preconfigure_svm = true
	# preconfigure_svm_and_update = true
	# services_deployed = ["services_deployed1","services_deployed2"]
	# precalculated_device_health = ["precalculated_device_health1","precalculated_device_health2"]
	# number_of_deployments = [0]
	# operation_type = ["operation_type1","operation_type2"]
	# operation_status = ["operation_status1","operation_status2"]
	# operation_data = ["operation_data1","operation_data2"]
	# current_step_count = ["current_step_count1", "current_step_count2"]
	# total_num_of_steps = ["total_num_of_steps1", "total_num_of_steps2"]
	# current_step_message = ["current_step_message1", "current_step_message2"]
	# custom_image = ["custom_image1", "custom_image2"]
	# original_deployment_id = ["original_deployment_id1", "original_deployment_id2"]
	# current_batch_count = ["current_batch_count1", "current_batch_count2"]
	# total_batch_count = ["total_batch_count1", "total_batch_count2"]
	# brownfield = true
	# overall_device_health = ["overall_device_health1","overall_device_health2"]
	# vds = false
	# scale_up = true
	# lifecycle_mode = false
	# can_migratev_clsv_ms = false
	# template_valid = true
	# configuration_change = true
	# detail_message = ["detail_message1", "detail_message2"]
	# timestamp = ["timestamp1","timestamp2"]
	# error = ["error1","error2"]
	# path = ["path1, "path2"]
  # }
}

output "resource_group_result" {
  value = data.powerflex_resource_group.example1.resource_group_details
}
