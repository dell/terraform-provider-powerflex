---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerflex_resource_group data source"
linkTitle: "powerflex_resource_group"
page_title: "powerflex_resource_group Data Source - powerflex"
subcategory: "Resource Group Management"
description: |-
  This datasource is used to query the existing ResourceGroup from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.
---

# powerflex_resource_group (Data Source)

This datasource is used to query the existing ResourceGroup from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.

## Example Usage

```terraform
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


# if a filter is of type string it has the ability to allow regular expressions
# data "powerflex_resource_group" "resource_group_filter_regex" {
#   filter{
#     name = ["^System_.*$"]
#     deployment_finished_date = ["^2024-01-10.*$"]
#   }
# }

# output "resourceGroupFilterRegexResult"{
#  value = data.powerflex_resource_group.resource_group_filter_regex.resource_group_details
# }

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
```

After the successful execution of above said block, We can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_resource_group.example1.attribute_name` where attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) Placeholder attribute.
- `resource_group_details` (Attributes Set) Resource Group details (see [below for nested schema](#nestedatt--resource_group_details))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `all_users_allowed` (Boolean) Value for all_users_allowed
- `brownfield` (Boolean) Value for brownfield
- `can_migratev_clsv_ms` (Boolean) Value for can_migratev_clsv_ms
- `compliant` (Boolean) Value for compliant
- `configuration_change` (Boolean) Value for configuration_change
- `created_by` (Set of String) List of created_by
- `created_date` (Set of String) List of created_date
- `current_batch_count` (Set of String) List of current_batch_count
- `current_step_count` (Set of String) List of current_step_count
- `current_step_message` (Set of String) List of current_step_message
- `custom_image` (Set of String) List of custom_image
- `deployment_description` (Set of String) List of deployment_description
- `deployment_finished_date` (Set of String) List of deployment_finished_date
- `deployment_health_status_type` (Set of String) List of deployment_health_status_type
- `deployment_name` (Set of String) List of deployment_name
- `deployment_scheduled_date` (Set of String) List of deployment_scheduled_date
- `deployment_started_date` (Set of String) List of deployment_started_date
- `detail_message` (Set of String) List of detail_message
- `disruptive_firmware` (Boolean) Value for disruptive_firmware
- `error` (Set of String) List of error
- `firmware_init` (Boolean) Value for firmware_init
- `firmware_repository_id` (Set of String) List of firmware_repository_id
- `id` (Set of String) List of id
- `individual_teardown` (Boolean) Value for individual_teardown
- `license_repository_id` (Set of String) List of license_repository_id
- `lifecycle_mode` (Boolean) Value for lifecycle_mode
- `no_op` (Boolean) Value for no_op
- `number_of_deployments` (Set of Number) List of number_of_deployments
- `operation_data` (Set of String) List of operation_data
- `operation_status` (Set of String) List of operation_status
- `operation_type` (Set of String) List of operation_type
- `original_deployment_id` (Set of String) List of original_deployment_id
- `overall_device_health` (Set of String) List of overall_device_health
- `owner` (Set of String) List of owner
- `path` (Set of String) List of path
- `precalculated_device_health` (Set of String) List of precalculated_device_health
- `preconfigure_svm` (Boolean) Value for preconfigure_svm
- `preconfigure_svm_and_update` (Boolean) Value for preconfigure_svm_and_update
- `remove_service` (Boolean) Value for remove_service
- `retry` (Boolean) Value for retry
- `scale_up` (Boolean) Value for scale_up
- `schedule_date` (Set of String) List of schedule_date
- `services_deployed` (Set of String) List of services_deployed
- `status` (Set of String) List of status
- `teardown` (Boolean) Value for teardown
- `teardown_after_cancel` (Boolean) Value for teardown_after_cancel
- `template_valid` (Boolean) Value for template_valid
- `timestamp` (Set of String) List of timestamp
- `total_batch_count` (Set of String) List of total_batch_count
- `total_num_of_steps` (Set of String) List of total_num_of_steps
- `update_server_firmware` (Boolean) Value for update_server_firmware
- `updated_by` (Set of String) List of updated_by
- `updated_date` (Set of String) List of updated_date
- `use_default_catalog` (Boolean) Value for use_default_catalog
- `vds` (Boolean) Value for vds


<a id="nestedatt--resource_group_details"></a>
### Nested Schema for `resource_group_details`

Read-Only:

- `all_users_allowed` (Boolean) Indicates whether all users are allowed for the deployment.
- `assigned_users` (Attributes List) List of users assigned to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--assigned_users))
- `brownfield` (Boolean) Indicates whether the deployment involves brownfield operations.
- `can_migratev_clsv_ms` (Boolean) Indicates whether migration of cluster virtual machines is allowed.
- `compliant` (Boolean) Indicates whether the deployment is compliant.
- `configuration_change` (Boolean) Specifies whether there has been a change in the deployment configuration.
- `created_by` (String) The user who created the deployment.
- `created_date` (String) The date when the deployment was created.
- `current_batch_count` (String) The current batch count during deployment.
- `current_step_count` (String) The current step count during deployment.
- `current_step_message` (String) The message associated with the current step during deployment.
- `custom_image` (String) The custom image used for deployment.
- `deployment_description` (String) The description of the ResourceGroup
- `deployment_device` (Attributes List) List of devices associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--deployment_device))
- `deployment_finished_date` (String) The date when the deployment finished.
- `deployment_health_status_type` (String) The type of health status associated with the deployment.
- `deployment_name` (String) The name of the ResourceGroup
- `deployment_scheduled_date` (String) The scheduled date for the deployment.
- `deployment_started_date` (String) The date when the deployment started.
- `deployment_valid` (Attributes) Details about the validity of the Resource Group (see [below for nested schema](#nestedatt--resource_group_details--deployment_valid))
- `deployment_validation_response` (Attributes) Details about the validation response for the deployment. (see [below for nested schema](#nestedatt--resource_group_details--deployment_validation_response))
- `detail_message` (String) Detailed Message
- `disruptive_firmware` (Boolean) Indicates whether disruptive firmware actions are allowed.
- `error` (String) Error
- `firmware_init` (Boolean) Indicates whether firmware initialization is performed during deployment.
- `firmware_repository` (Attributes) Details about the firmware repository used by the template. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository))
- `firmware_repository_id` (String) The ID of the firmware repository associated with the deployment.
- `id` (String) The unique identifier for the Resource Group
- `individual_teardown` (Boolean) Indicates whether to perform individual teardown for the deployment.
- `job_details` (Attributes List) List of job details associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--job_details))
- `license_repository` (Attributes) Details about the license repository associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--license_repository))
- `license_repository_id` (String) The ID of the license repository associated with the deployment.
- `lifecycle_mode` (Boolean) Indicates whether the deployment is in lifecycle mode.
- `lifecycle_mode_reasons` (List of String) List of reasons for the lifecycle mode of the deployment.
- `messages` (Attributes List) Messages (see [below for nested schema](#nestedatt--resource_group_details--messages))
- `no_op` (Boolean) Indicates whether the deployment is a no-op (no operation).
- `number_of_deployments` (Number) The total number of deployments.
- `operation_data` (String) Additional data associated with the operation.
- `operation_status` (String) The status of the operation associated with the deployment.
- `operation_type` (String) The type of operation associated with the deployment.
- `original_deployment_id` (String) The ID of the original deployment.
- `overall_device_health` (String) The overall health status of the device in the deployment.
- `owner` (String) The owner of the deployment.
- `path` (String) Path
- `precalculated_device_health` (String) The precalculated health of devices associated with the deployment.
- `preconfigure_svm` (Boolean) Indicates whether to preconfigure SVM (Storage Virtual Machine).
- `preconfigure_svm_and_update` (Boolean) Indicates whether to preconfigure SVM and perform an update.
- `remove_service` (Boolean) Indicates whether the associated service should be removed.
- `retry` (Boolean) Indicates whether the deployment should be retried.
- `scale_up` (Boolean) Indicates whether the deployment involves scaling up.
- `schedule_date` (String) The date when the deployment is scheduled.
- `service_template` (Attributes) Template details (see [below for nested schema](#nestedatt--resource_group_details--service_template))
- `services_deployed` (String) Details about the services deployed during the deployment.
- `status` (String) The status of the deployment.
- `teardown` (Boolean) teardown
- `teardown_after_cancel` (Boolean) Indicates whether teardown should occur after canceling the deployment.
- `template_valid` (Boolean) Details about the validity of the template.
- `timestamp` (String) The timestamp indicating when the message was generated.
- `total_batch_count` (String) The total number of batches involved in the deployment.
- `total_num_of_steps` (String) The total number of steps involved in the deployment.
- `update_server_firmware` (Boolean) Indicates whether to update server firmware during the deployment.
- `updated_by` (String) The user who last updated the deployment.
- `updated_date` (String) The date when the deployment was last updated.
- `use_default_catalog` (Boolean) Indicates whether to use the default catalog for the deployment.
- `vds` (Boolean) Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.
- `vms` (Attributes List) List of virtual machines associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--vms))

<a id="nestedatt--resource_group_details--assigned_users"></a>
### Nested Schema for `resource_group_details.assigned_users`

Read-Only:

- `created_by` (String) The user who created the assigned user.
- `created_date` (String) The date when the assigned user was created.
- `domain_name` (String) The domain name of the assigned user.
- `email` (String) The email address of the assigned user.
- `enabled` (Boolean) Indicates whether the assigned user is enabled.
- `first_name` (String) The first name of the assigned user.
- `group_dn` (String) The distinguished name (DN) of the group associated with the assigned user.
- `group_name` (String) The name of the group associated with the assigned user.
- `id` (String) The unique identifier of the assigned user.
- `last_name` (String) The last name of the assigned user.
- `password` (String) The password associated with the assigned user.
- `phone_number` (String) The phone number of the assigned user.
- `role` (String) The role associated with the assigned user.
- `roles` (List of String) The roles associated with the assigned user.
- `system_user` (Boolean) Indicates whether the assigned user is a system user.
- `update_password` (Boolean) Indicates whether the user password needs to be updated.
- `updated_by` (String) The user who last updated the assigned user.
- `updated_date` (String) The date when the assigned user was last updated.
- `user_name` (String) The username of the assigned user.
- `user_preference` (String) The preferences of the assigned user.
- `user_seq_id` (Number) The sequential ID of the assigned user.


<a id="nestedatt--resource_group_details--deployment_device"></a>
### Nested Schema for `resource_group_details.deployment_device`

Read-Only:

- `brownfield` (Boolean) Indicates whether the deployment device is associated with a brownfield deployment.
- `brownfield_status` (String) The brownfield status of the deployment device.
- `cloud_link` (Boolean) Indicates whether the deployment device has a cloud link.
- `compliant_state` (String) The compliant state of the deployment device.
- `component_id` (String) The component ID associated with the deployment device.
- `current_ip_address` (String) The current IP address of the deployment device.
- `das_cache` (Boolean) Indicates whether the deployment device has Direct-Attached Storage (DAS) cache.
- `device_group_name` (String) The name of the device group associated with the deployment device.
- `device_health` (String) The health status of the deployment device.
- `device_state` (String) The state of the deployment device.
- `device_type` (String) The type of device associated with the deployment device.
- `health_message` (String) The health message associated with the deployment device.
- `ip_address` (String) The IP address of the deployment device.
- `log_dump` (String) The log dump information associated with the deployment device.
- `model` (String) The model of the deployment device.
- `puppet_cert_name` (String) The Puppet certificate name associated with the deployment device.
- `ref_id` (String) The reference ID associated with the deployment device.
- `ref_type` (String) The reference type associated with the deployment device.
- `service_tag` (String) The service tag associated with the deployment device.
- `status` (String) The status of the deployment device.
- `status_end_time` (String) The end time of the status for the deployment device.
- `status_message` (String) The status message associated with the deployment device.
- `status_start_time` (String) The start time of the status for the deployment device.


<a id="nestedatt--resource_group_details--deployment_valid"></a>
### Nested Schema for `resource_group_details.deployment_valid`

Read-Only:

- `messages` (Attributes List) List of messages related to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--deployment_valid--messages))
- `valid` (Boolean) Indicates whether the deployment is valid.

<a id="nestedatt--resource_group_details--deployment_valid--messages"></a>
### Nested Schema for `resource_group_details.deployment_valid.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.



<a id="nestedatt--resource_group_details--deployment_validation_response"></a>
### Nested Schema for `resource_group_details.deployment_validation_response`

Read-Only:

- `can_deploy` (Boolean) Indicates whether the deployment can be executed.
- `disk_type_mismatch` (Boolean) Indicates whether there is a disk type mismatch in the deployment.
- `drives_per_storage_pool` (Number) The number of drives per storage pool in the deployment.
- `hostnames` (List of String) A list of hostnames associated with the deployment.
- `max_scalability` (Number) The maximum scalability of the deployment.
- `new_node_disk_types` (List of String) The disk types associated with new nodes in the deployment.
- `no_of_fault_sets` (Number) The number of fault sets in the deployment.
- `nodes` (Number) The number of nodes in the deployment.
- `nodes_per_fault_set` (Number) The number of nodes per fault set in the deployment.
- `number_of_service_volumes` (Number) The number of service volumes in the deployment.
- `protection_domain` (String) The protection domain associated with the deployment.
- `storage_pool_disk_type` (List of String) The disk types associated with each storage pool in the deployment.
- `storage_pools` (Number) The number of storage pools in the deployment.
- `virtual_machines` (Number) The number of virtual machines in the deployment.
- `warning_messages` (List of String) A list of warning messages associated with the deployment validation.


<a id="nestedatt--resource_group_details--firmware_repository"></a>
### Nested Schema for `resource_group_details.firmware_repository`

Read-Only:

- `bundle_count` (Number) The count of software bundles in the firmware repository.
- `component_count` (Number) The count of software components in the firmware repository.
- `created_by` (String) The user who created the firmware repository.
- `created_date` (String) The date when the firmware repository was created.
- `custom` (Boolean) Indicates whether the firmware repository is custom.
- `default_catalog` (Boolean) Indicates whether the firmware repository is the default catalog.
- `deployments` (Attributes List) List of deployments associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments))
- `disk_location` (String) The location on disk where the firmware repository is stored.
- `download_progress` (Number) The progress of the download for the firmware repository.
- `download_status` (String) The download status of the firmware repository.
- `embedded` (Boolean) Indicates whether the firmware repository is embedded.
- `extract_progress` (Number) The progress of the extraction for the firmware repository.
- `file_size_in_gigabytes` (Number) The size of the firmware repository file in gigabytes.
- `filename` (String) The filename of the firmware repository.
- `id` (String) The unique identifier of the firmware repository.
- `job_id` (String) The unique identifier of the job associated with the firmware repository.
- `md_5_hash` (String) The MD5 hash of the firmware repository.
- `minimal` (Boolean) Indicates whether the firmware repository is minimal.
- `name` (String) The name of the firmware repository.
- `needs_attention` (Boolean) Indicates whether the firmware repository needs attention.
- `password` (String) The password associated with the firmware repository.
- `rcmapproved` (Boolean) Indicates whether the firmware repository is RCM approved.
- `signature` (String) The signature of the firmware repository.
- `signed_key_source_location` (String) The source location of the signed key associated with the firmware repository.
- `software_bundles` (Attributes List) List of software bundles associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--software_bundles))
- `software_components` (Attributes List) List of software components associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--software_components))
- `source_location` (String) The location of the source for the firmware repository.
- `source_type` (String) The type of the source for the firmware repository.
- `state` (String) The state of the firmware repository.
- `updated_by` (String) The user who last updated the firmware repository.
- `updated_date` (String) The date when the firmware repository was last updated.
- `user_bundle_count` (Number) The count of user-specific software bundles in the firmware repository.
- `username` (String) The username associated with the firmware repository.

<a id="nestedatt--resource_group_details--firmware_repository--deployments"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments`

Read-Only:

- `all_users_allowed` (Boolean) Indicates whether all users are allowed for the deployment.
- `assigned_users` (Attributes List) List of users assigned to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--assigned_users))
- `brownfield` (Boolean) Indicates whether the deployment involves brownfield operations.
- `can_migratev_clsv_ms` (Boolean) Indicates whether migration of cluster virtual machines is allowed.
- `compliant` (Boolean) Indicates whether the deployment is compliant.
- `configuration_change` (Boolean) Specifies whether there has been a change in the deployment configuration.
- `created_by` (String) The user who created the deployment.
- `created_date` (String) The date when the deployment was created.
- `current_batch_count` (String) The current batch count during deployment.
- `current_step_count` (String) The current step count during deployment.
- `current_step_message` (String) The message associated with the current step during deployment.
- `custom_image` (String) The custom image used for deployment.
- `deployment_description` (String) The description of the deployment.
- `deployment_device` (Attributes List) List of devices associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--deployment_device))
- `deployment_finished_date` (String) The date when the deployment finished.
- `deployment_health_status_type` (String) The type of health status associated with the deployment.
- `deployment_name` (String) The name of the deployment.
- `deployment_scheduled_date` (String) The scheduled date for the deployment.
- `deployment_started_date` (String) The date when the deployment started.
- `deployment_valid` (Attributes) Details about the validity of the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--deployment_valid))
- `deployment_validation_response` (Attributes) Details about the validation response for the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--deployment_validation_response))
- `disruptive_firmware` (Boolean) Indicates whether disruptive firmware actions are allowed.
- `firmware_init` (Boolean) Indicates whether firmware initialization is performed during deployment.
- `firmware_repository_id` (String) The ID of the firmware repository associated with the deployment.
- `id` (String) The unique identifier of the deployment.
- `individual_teardown` (Boolean) Indicates whether to perform individual teardown for the deployment.
- `job_details` (Attributes List) List of job details associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--job_details))
- `license_repository` (Attributes) Details about the license repository associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--license_repository))
- `license_repository_id` (String) The ID of the license repository associated with the deployment.
- `lifecycle_mode` (Boolean) Indicates whether the deployment is in lifecycle mode.
- `lifecycle_mode_reasons` (List of String) List of reasons for the lifecycle mode of the deployment.
- `no_op` (Boolean) Indicates whether the deployment is a no-op (no operation).
- `number_of_deployments` (Number) The total number of deployments.
- `operation_data` (String) Additional data associated with the operation.
- `operation_status` (String) The status of the operation associated with the deployment.
- `operation_type` (String) The type of operation associated with the deployment.
- `original_deployment_id` (String) The ID of the original deployment.
- `overall_device_health` (String) The overall health status of the device in the deployment.
- `owner` (String) The owner of the deployment.
- `precalculated_device_health` (String) The precalculated health of devices associated with the deployment.
- `preconfigure_svm` (Boolean) Indicates whether to preconfigure SVM (Storage Virtual Machine).
- `preconfigure_svm_and_update` (Boolean) Indicates whether to preconfigure SVM and perform an update.
- `remove_service` (Boolean) Indicates whether the associated service should be removed.
- `retry` (Boolean) Indicates whether the deployment should be retried.
- `scale_up` (Boolean) Indicates whether the deployment involves scaling up.
- `schedule_date` (String) The date when the deployment is scheduled.
- `services_deployed` (String) Details about the services deployed during the deployment.
- `status` (String) The status of the deployment.
- `teardown` (Boolean) Indicates whether the deployment should be torn down.
- `teardown_after_cancel` (Boolean) Indicates whether teardown should occur after canceling the deployment.
- `template_valid` (Boolean) Indicates if the deployment template is valid.
- `total_batch_count` (String) The total number of batches involved in the deployment.
- `total_num_of_steps` (String) The total number of steps involved in the deployment.
- `update_server_firmware` (Boolean) Indicates whether to update server firmware during the deployment.
- `updated_by` (String) The user who last updated the deployment.
- `updated_date` (String) The date when the deployment was last updated.
- `use_default_catalog` (Boolean) Indicates whether to use the default catalog for the deployment.
- `vds` (Boolean) Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.
- `vms` (Attributes List) List of virtual machines associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--vms))

<a id="nestedatt--resource_group_details--firmware_repository--deployments--assigned_users"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.assigned_users`

Read-Only:

- `created_by` (String) The user who created the assigned user.
- `created_date` (String) The date when the assigned user was created.
- `domain_name` (String) The domain name of the assigned user.
- `email` (String) The email address of the assigned user.
- `enabled` (Boolean) Indicates whether the assigned user is enabled.
- `first_name` (String) The first name of the assigned user.
- `group_dn` (String) The distinguished name (DN) of the group associated with the assigned user.
- `group_name` (String) The name of the group associated with the assigned user.
- `id` (String) The unique identifier of the assigned user.
- `last_name` (String) The last name of the assigned user.
- `password` (String) The password associated with the assigned user.
- `phone_number` (String) The phone number of the assigned user.
- `role` (String) The role associated with the assigned user.
- `roles` (List of String) The roles associated with the assigned user.
- `system_user` (Boolean) Indicates whether the assigned user is a system user.
- `update_password` (Boolean) Indicates whether the user password needs to be updated.
- `updated_by` (String) The user who last updated the assigned user.
- `updated_date` (String) The date when the assigned user was last updated.
- `user_name` (String) The username of the assigned user.
- `user_preference` (String) The preferences of the assigned user.
- `user_seq_id` (Number) The sequential ID of the assigned user.


<a id="nestedatt--resource_group_details--firmware_repository--deployments--deployment_device"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.deployment_device`

Read-Only:

- `brownfield` (Boolean) Indicates whether the deployment device is associated with a brownfield deployment.
- `brownfield_status` (String) The brownfield status of the deployment device.
- `cloud_link` (Boolean) Indicates whether the deployment device has a cloud link.
- `compliant_state` (String) The compliant state of the deployment device.
- `component_id` (String) The component ID associated with the deployment device.
- `current_ip_address` (String) The current IP address of the deployment device.
- `das_cache` (Boolean) Indicates whether the deployment device has Direct-Attached Storage (DAS) cache.
- `device_group_name` (String) The name of the device group associated with the deployment device.
- `device_health` (String) The health status of the deployment device.
- `device_state` (String) The state of the deployment device.
- `device_type` (String) The type of device associated with the deployment device.
- `health_message` (String) The health message associated with the deployment device.
- `ip_address` (String) The IP address of the deployment device.
- `log_dump` (String) The log dump information associated with the deployment device.
- `model` (String) The model of the deployment device.
- `puppet_cert_name` (String) The Puppet certificate name associated with the deployment device.
- `ref_id` (String) The reference ID associated with the deployment device.
- `ref_type` (String) The reference type associated with the deployment device.
- `service_tag` (String) The service tag associated with the deployment device.
- `status` (String) The status of the deployment device.
- `status_end_time` (String) The end time of the status for the deployment device.
- `status_message` (String) The status message associated with the deployment device.
- `status_start_time` (String) The start time of the status for the deployment device.


<a id="nestedatt--resource_group_details--firmware_repository--deployments--deployment_valid"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.deployment_valid`

Read-Only:

- `messages` (Attributes List) List of messages related to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--deployments--deployment_valid--messages))
- `valid` (Boolean) Indicates whether the deployment is valid.

<a id="nestedatt--resource_group_details--firmware_repository--deployments--deployment_valid--messages"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.deployment_valid.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.



<a id="nestedatt--resource_group_details--firmware_repository--deployments--deployment_validation_response"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.deployment_validation_response`

Read-Only:

- `can_deploy` (Boolean) Indicates whether the deployment can be executed.
- `disk_type_mismatch` (Boolean) Indicates whether there is a disk type mismatch in the deployment.
- `drives_per_storage_pool` (Number) The number of drives per storage pool in the deployment.
- `hostnames` (List of String) A list of hostnames associated with the deployment.
- `max_scalability` (Number) The maximum scalability of the deployment.
- `new_node_disk_types` (List of String) The disk types associated with new nodes in the deployment.
- `no_of_fault_sets` (Number) The number of fault sets in the deployment.
- `nodes` (Number) The number of nodes in the deployment.
- `nodes_per_fault_set` (Number) The number of nodes per fault set in the deployment.
- `number_of_service_volumes` (Number) The number of service volumes in the deployment.
- `protection_domain` (String) The protection domain associated with the deployment.
- `storage_pool_disk_type` (List of String) The disk types associated with each storage pool in the deployment.
- `storage_pools` (Number) The number of storage pools in the deployment.
- `virtual_machines` (Number) The number of virtual machines in the deployment.
- `warning_messages` (List of String) A list of warning messages associated with the deployment validation.


<a id="nestedatt--resource_group_details--firmware_repository--deployments--job_details"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.job_details`

Read-Only:

- `component_id` (String) The unique identifier of the component associated with the job.
- `execution_id` (String) The unique identifier of the job execution.
- `level` (String) The log level of the job.
- `message` (String) The log message of the job.
- `timestamp` (String) The timestamp of the job execution.


<a id="nestedatt--resource_group_details--firmware_repository--deployments--license_repository"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.license_repository`

Read-Only:

- `created_by` (String) The user who created the license repository.
- `created_date` (String) The date when the license repository was created.
- `disk_location` (String) The disk location of the license repository.
- `filename` (String) The filename associated with the license repository.
- `id` (String) The unique identifier of the license repository.
- `license_data` (String) The license data associated with the license repository.
- `name` (String) The name of the license repository.
- `state` (String) The state of the license repository.
- `type` (String) The type of the license repository.
- `updated_by` (String) The user who last updated the license repository.
- `updated_date` (String) The date when the license repository was last updated.


<a id="nestedatt--resource_group_details--firmware_repository--deployments--vms"></a>
### Nested Schema for `resource_group_details.firmware_repository.deployments.vms`

Read-Only:

- `certificate_name` (String) The certificate name associated with the virtual machine (VM).
- `vm_ipaddress` (String) The IP address of the virtual machine (VM).
- `vm_manufacturer` (String) The manufacturer of the virtual machine (VM).
- `vm_model` (String) The model of the virtual machine (VM).
- `vm_service_tag` (String) The service tag associated with the virtual machine (VM).



<a id="nestedatt--resource_group_details--firmware_repository--software_bundles"></a>
### Nested Schema for `resource_group_details.firmware_repository.software_bundles`

Read-Only:

- `bundle_date` (String) The date when the software bundle was created.
- `bundle_type` (String) The type of the software bundle.
- `created_by` (String) The user who initially created the software bundle.
- `created_date` (String) The date when the software bundle was initially created.
- `criticality` (String) The criticality level of the software bundle.
- `custom` (Boolean) Indicates whether the software bundle is custom.
- `description` (String) A brief description of the software bundle.
- `device_model` (String) The model of the device associated with the software bundle.
- `device_type` (String) The type of device associated with the software bundle.
- `fw_repository_id` (String) The identifier of the firmware repository associated with the software bundle.
- `id` (String) The unique identifier for the software bundle.
- `name` (String) The name of the software bundle.
- `needs_attention` (Boolean) Indicates whether the software bundle needs attention.
- `software_components` (Attributes List) List of software components associated with the software bundle. (see [below for nested schema](#nestedatt--resource_group_details--firmware_repository--software_bundles--software_components))
- `updated_by` (String) The user who last updated the software bundle.
- `updated_date` (String) The date when the software bundle was last updated.
- `user_bundle` (Boolean) Indicates whether the software bundle is a user-specific bundle.
- `user_bundle_hash_md_5` (String) The MD5 hash value of the user-specific software bundle.
- `user_bundle_path` (String) The path associated with the user-specific software bundle.
- `version` (String) The version of the software bundle.

<a id="nestedatt--resource_group_details--firmware_repository--software_bundles--software_components"></a>
### Nested Schema for `resource_group_details.firmware_repository.software_bundles.software_components`

Read-Only:

- `category` (String) The category to which the component belongs.
- `component_id` (String) The identifier of the component.
- `component_type` (String) The type of the component.
- `created_by` (String) The user who created the component.
- `created_date` (String) The date when the component was created.
- `custom` (Boolean) Indicates whether the component is custom or not.
- `dell_version` (String) The version of the component according to Dell standards.
- `device_id` (String) The identifier of the device associated with the component.
- `firmware_repo_name` (String) The name of the firmware repository associated with the component.
- `hash_md_5` (String) The MD5 hash value of the component.
- `id` (String) The unique identifier for the software component.
- `ignore` (Boolean) Indicates whether the component should be ignored.
- `name` (String) The name of the software component.
- `needs_attention` (Boolean) Indicates whether the component needs attention.
- `operating_system` (String) The operating system associated with the component.
- `original_component_id` (String) The identifier of the original component.
- `original_version` (String) The original version of the component.
- `package_id` (String) The identifier of the package to which the component belongs.
- `path` (String) The path where the component is stored.
- `sub_device_id` (String) The sub-identifier of the device associated with the component.
- `sub_vendor_id` (String) The sub-identifier of the vendor associated with the component.
- `system_ids` (List of String) List of system IDs associated with the component.
- `updated_by` (String) The user who last updated the component.
- `updated_date` (String) The date when the component was last updated.
- `vendor_id` (String) The identifier of the vendor associated with the component.
- `vendor_version` (String) The version of the component according to the vendor's standards.



<a id="nestedatt--resource_group_details--firmware_repository--software_components"></a>
### Nested Schema for `resource_group_details.firmware_repository.software_components`

Read-Only:

- `category` (String) The category to which the component belongs.
- `component_id` (String) The identifier of the component.
- `component_type` (String) The type of the component.
- `created_by` (String) The user who created the component.
- `created_date` (String) The date when the component was created.
- `custom` (Boolean) Indicates whether the component is custom or not.
- `dell_version` (String) The version of the component according to Dell standards.
- `device_id` (String) The identifier of the device associated with the component.
- `firmware_repo_name` (String) The name of the firmware repository associated with the component.
- `hash_md_5` (String) The MD5 hash value of the component.
- `id` (String) The unique identifier for the software component.
- `ignore` (Boolean) Indicates whether the component should be ignored.
- `name` (String) The name of the software component.
- `needs_attention` (Boolean) Indicates whether the component needs attention.
- `operating_system` (String) The operating system associated with the component.
- `original_component_id` (String) The identifier of the original component.
- `original_version` (String) The original version of the component.
- `package_id` (String) The identifier of the package to which the component belongs.
- `path` (String) The path where the component is stored.
- `sub_device_id` (String) The sub-identifier of the device associated with the component.
- `sub_vendor_id` (String) The sub-identifier of the vendor associated with the component.
- `system_ids` (List of String) List of system IDs associated with the component.
- `updated_by` (String) The user who last updated the component.
- `updated_date` (String) The date when the component was last updated.
- `vendor_id` (String) The identifier of the vendor associated with the component.
- `vendor_version` (String) The version of the component according to the vendor's standards.



<a id="nestedatt--resource_group_details--job_details"></a>
### Nested Schema for `resource_group_details.job_details`

Read-Only:

- `component_id` (String) The unique identifier of the component associated with the job.
- `execution_id` (String) The unique identifier of the job execution.
- `level` (String) The log level of the job.
- `message` (String) The log message of the job.
- `timestamp` (String) The timestamp of the job execution.


<a id="nestedatt--resource_group_details--license_repository"></a>
### Nested Schema for `resource_group_details.license_repository`

Read-Only:

- `created_by` (String) The user who created the license repository.
- `created_date` (String) The date when the license repository was created.
- `disk_location` (String) The disk location of the license repository.
- `filename` (String) The filename associated with the license repository.
- `id` (String) The unique identifier of the license repository.
- `license_data` (String) The license data associated with the license repository.
- `name` (String) The name of the license repository.
- `state` (String) The state of the license repository.
- `type` (String) The type of the license repository.
- `updated_by` (String) The user who last updated the license repository.
- `updated_date` (String) The date when the license repository was last updated.


<a id="nestedatt--resource_group_details--messages"></a>
### Nested Schema for `resource_group_details.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.


<a id="nestedatt--resource_group_details--service_template"></a>
### Nested Schema for `resource_group_details.service_template`

Read-Only:

- `all_users_allowed` (Boolean) Indicates whether all users are allowed for the template.
- `assigned_users` (Attributes List) List of users assigned to the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--assigned_users))
- `brownfield_template_type` (String) The type of template for brownfield deployments.
- `category` (String) The category to which the template belongs.
- `cluster_count` (Number) The count of clusters associated with the template.
- `components` (Attributes List) List of components included in the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--components))
- `configuration` (Attributes) Details about the configuration settings of the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--configuration))
- `created_by` (String) The user who created the template.
- `created_date` (String) The date when the template was created.
- `draft` (Boolean) Indicates whether the template is in draft mode.
- `firmware_repository` (Attributes) Details about the firmware repository used by the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository))
- `id` (String) The unique identifier for the template.
- `in_configuration` (Boolean) Indicates whether the template is part of the current configuration.
- `last_deployed_date` (String) The date when the template was last deployed.
- `license_repository` (Attributes) Details about the license repository used by the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--license_repository))
- `manage_firmware` (Boolean) Indicates whether firmware is managed by the template.
- `networks` (Attributes List) List of networks associated with the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--networks))
- `original_template_id` (String) The ID of the original template if this is a derived template.
- `sdnas_count` (Number) The count of software-defined network appliances associated with the template.
- `server_count` (Number) The count of servers associated with the template.
- `service_count` (Number) The count of services associated with the template.
- `storage_count` (Number) The count of storage devices associated with the template.
- `switch_count` (Number) The count of switches associated with the template.
- `template_description` (String) The description of the template.
- `template_locked` (Boolean) Indicates whether the template is locked or not.
- `template_name` (String) The name of the template.
- `template_type` (String) The type/category of the template.
- `template_valid` (Attributes) Details about the validity of the template. (see [below for nested schema](#nestedatt--resource_group_details--service_template--template_valid))
- `template_version` (String) The version of the template.
- `updated_by` (String) The user who last updated the template.
- `updated_date` (String) The date when the template was last updated.
- `use_default_catalog` (Boolean) Indicates whether the default catalog is used for the template.
- `vm_count` (Number) The count of virtual machines associated with the template.

<a id="nestedatt--resource_group_details--service_template--assigned_users"></a>
### Nested Schema for `resource_group_details.service_template.assigned_users`

Read-Only:

- `created_by` (String) The user who created the assigned user.
- `created_date` (String) The date when the assigned user was created.
- `domain_name` (String) The domain name of the assigned user.
- `email` (String) The email address of the assigned user.
- `enabled` (Boolean) Indicates whether the assigned user is enabled.
- `first_name` (String) The first name of the assigned user.
- `group_dn` (String) The distinguished name (DN) of the group associated with the assigned user.
- `group_name` (String) The name of the group associated with the assigned user.
- `id` (String) The unique identifier of the assigned user.
- `last_name` (String) The last name of the assigned user.
- `password` (String) The password associated with the assigned user.
- `phone_number` (String) The phone number of the assigned user.
- `role` (String) The role associated with the assigned user.
- `roles` (List of String) The roles associated with the assigned user.
- `system_user` (Boolean) Indicates whether the assigned user is a system user.
- `update_password` (Boolean) Indicates whether the user password needs to be updated.
- `updated_by` (String) The user who last updated the assigned user.
- `updated_date` (String) The date when the assigned user was last updated.
- `user_name` (String) The username of the assigned user.
- `user_preference` (String) The preferences of the assigned user.
- `user_seq_id` (Number) The sequential ID of the assigned user.


<a id="nestedatt--resource_group_details--service_template--components"></a>
### Nested Schema for `resource_group_details.service_template.components`

Read-Only:

- `asm_guid` (String) The ASM GUID (Global Unique Identifier) associated with the component.
- `brownfield` (Boolean) Indicates whether the component is brownfield.
- `changed` (Boolean) Indicates whether the component has changed.
- `cloned` (Boolean) Indicates whether the component is cloned.
- `cloned_from_asm_guid` (String) The ASM GUID from which the component is cloned.
- `cloned_from_id` (String) The identifier of the component from which this component is cloned.
- `component_id` (String) The identifier for the component.
- `component_valid` (Attributes) Information about the validity of the component. (see [below for nested schema](#nestedatt--resource_group_details--service_template--components--component_valid))
- `config_file` (String) The configuration file associated with the component.
- `help_text` (String) Help text associated with the component.
- `id` (String) The unique identifier for the component.
- `identifier` (String) The identifier for the component.
- `instances` (Number) The number of instances of the component.
- `ip` (String) The IP address associated with the component.
- `manage_firmware` (Boolean) Indicates whether firmware is managed for the component.
- `management_ip_address` (String) The management IP address of the component.
- `name` (String) The name of the component.
- `os_puppet_cert_name` (String) The OS Puppet certificate name associated with the component.
- `puppet_cert_name` (String) The Puppet certificate name associated with the component.
- `ref_id` (String) The reference identifier associated with the component.
- `related_components` (Map of String) Related components associated with this component.
- `resources` (Attributes List) List of resources associated with the component. (see [below for nested schema](#nestedatt--resource_group_details--service_template--components--resources))
- `serial_number` (String) The serial number of the component.
- `sub_type` (String) The sub-type of the component.
- `teardown` (Boolean) Indicates whether the component should be torn down.
- `type` (String) The type of the component.

<a id="nestedatt--resource_group_details--service_template--components--component_valid"></a>
### Nested Schema for `resource_group_details.service_template.components.component_valid`

Read-Only:

- `messages` (Attributes List) List of messages associated with the component validity. (see [below for nested schema](#nestedatt--resource_group_details--service_template--components--component_valid--messages))
- `valid` (Boolean) Indicates whether the component is valid.

<a id="nestedatt--resource_group_details--service_template--components--component_valid--messages"></a>
### Nested Schema for `resource_group_details.service_template.components.component_valid.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.



<a id="nestedatt--resource_group_details--service_template--components--resources"></a>
### Nested Schema for `resource_group_details.service_template.components.resources`

Read-Only:

- `display_name` (String) The display name for the resources.
- `guid` (String) The globally unique identifier (GUID) for the resources.
- `id` (String) The identifier for the resources.



<a id="nestedatt--resource_group_details--service_template--configuration"></a>
### Nested Schema for `resource_group_details.service_template.configuration`

Read-Only:

- `categories` (Attributes List) List of categories associated with the configuration. (see [below for nested schema](#nestedatt--resource_group_details--service_template--configuration--categories))
- `comparator` (String) Comparator used in the configuration.
- `controller_fqdd` (String) Fully Qualified Device Descriptor (FQDD) of the controller in the configuration.
- `disktype` (String) Type of disk in the configuration.
- `id` (String) Unique identifier for the configuration.
- `numberofdisks` (Number) Number of disks in the configuration.
- `raidlevel` (String) RAID level of the configuration.
- `virtual_disk_fqdd` (String) Fully Qualified Device Descriptor (FQDD) of the virtual disk in the configuration.

<a id="nestedatt--resource_group_details--service_template--configuration--categories"></a>
### Nested Schema for `resource_group_details.service_template.configuration.categories`

Read-Only:

- `device_type` (String) The type of device associated with the category.
- `display_name` (String) The display name of the category.
- `id` (String) The unique identifier for the category.



<a id="nestedatt--resource_group_details--service_template--firmware_repository"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository`

Read-Only:

- `bundle_count` (Number) The count of software bundles in the firmware repository.
- `component_count` (Number) The count of software components in the firmware repository.
- `created_by` (String) The user who created the firmware repository.
- `created_date` (String) The date when the firmware repository was created.
- `custom` (Boolean) Indicates whether the firmware repository is custom.
- `default_catalog` (Boolean) Indicates whether the firmware repository is the default catalog.
- `deployments` (Attributes List) List of deployments associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments))
- `disk_location` (String) The location on disk where the firmware repository is stored.
- `download_progress` (Number) The progress of the download for the firmware repository.
- `download_status` (String) The download status of the firmware repository.
- `embedded` (Boolean) Indicates whether the firmware repository is embedded.
- `extract_progress` (Number) The progress of the extraction for the firmware repository.
- `file_size_in_gigabytes` (Number) The size of the firmware repository file in gigabytes.
- `filename` (String) The filename of the firmware repository.
- `id` (String) The unique identifier of the firmware repository.
- `job_id` (String) The unique identifier of the job associated with the firmware repository.
- `md_5_hash` (String) The MD5 hash of the firmware repository.
- `minimal` (Boolean) Indicates whether the firmware repository is minimal.
- `name` (String) The name of the firmware repository.
- `needs_attention` (Boolean) Indicates whether the firmware repository needs attention.
- `password` (String) The password associated with the firmware repository.
- `rcmapproved` (Boolean) Indicates whether the firmware repository is RCM approved.
- `signature` (String) The signature of the firmware repository.
- `signed_key_source_location` (String) The source location of the signed key associated with the firmware repository.
- `software_bundles` (Attributes List) List of software bundles associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--software_bundles))
- `software_components` (Attributes List) List of software components associated with the firmware repository. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--software_components))
- `source_location` (String) The location of the source for the firmware repository.
- `source_type` (String) The type of the source for the firmware repository.
- `state` (String) The state of the firmware repository.
- `updated_by` (String) The user who last updated the firmware repository.
- `updated_date` (String) The date when the firmware repository was last updated.
- `user_bundle_count` (Number) The count of user-specific software bundles in the firmware repository.
- `username` (String) The username associated with the firmware repository.

<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments`

Read-Only:

- `all_users_allowed` (Boolean) Indicates whether all users are allowed for the deployment.
- `assigned_users` (Attributes List) List of users assigned to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--assigned_users))
- `brownfield` (Boolean) Indicates whether the deployment involves brownfield operations.
- `can_migratev_clsv_ms` (Boolean) Indicates whether migration of cluster virtual machines is allowed.
- `compliant` (Boolean) Indicates whether the deployment is compliant.
- `configuration_change` (Boolean) Specifies whether there has been a change in the deployment configuration.
- `created_by` (String) The user who created the deployment.
- `created_date` (String) The date when the deployment was created.
- `current_batch_count` (String) The current batch count during deployment.
- `current_step_count` (String) The current step count during deployment.
- `current_step_message` (String) The message associated with the current step during deployment.
- `custom_image` (String) The custom image used for deployment.
- `deployment_description` (String) The description of the deployment.
- `deployment_device` (Attributes List) List of devices associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_device))
- `deployment_finished_date` (String) The date when the deployment finished.
- `deployment_health_status_type` (String) The type of health status associated with the deployment.
- `deployment_name` (String) The name of the deployment.
- `deployment_scheduled_date` (String) The scheduled date for the deployment.
- `deployment_started_date` (String) The date when the deployment started.
- `deployment_valid` (Attributes) Details about the validity of the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_valid))
- `deployment_validation_response` (Attributes) Details about the validation response for the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_validation_response))
- `disruptive_firmware` (Boolean) Indicates whether disruptive firmware actions are allowed.
- `firmware_init` (Boolean) Indicates whether firmware initialization is performed during deployment.
- `firmware_repository_id` (String) The ID of the firmware repository associated with the deployment.
- `id` (String) The unique identifier of the deployment.
- `individual_teardown` (Boolean) Indicates whether to perform individual teardown for the deployment.
- `job_details` (Attributes List) List of job details associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--job_details))
- `license_repository` (Attributes) Details about the license repository associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--license_repository))
- `license_repository_id` (String) The ID of the license repository associated with the deployment.
- `lifecycle_mode` (Boolean) Indicates whether the deployment is in lifecycle mode.
- `lifecycle_mode_reasons` (List of String) List of reasons for the lifecycle mode of the deployment.
- `no_op` (Boolean) Indicates whether the deployment is a no-op (no operation).
- `number_of_deployments` (Number) The total number of deployments.
- `operation_data` (String) Additional data associated with the operation.
- `operation_status` (String) The status of the operation associated with the deployment.
- `operation_type` (String) The type of operation associated with the deployment.
- `original_deployment_id` (String) The ID of the original deployment.
- `overall_device_health` (String) The overall health status of the device in the deployment.
- `owner` (String) The owner of the deployment.
- `precalculated_device_health` (String) The precalculated health of devices associated with the deployment.
- `preconfigure_svm` (Boolean) Indicates whether to preconfigure SVM (Storage Virtual Machine).
- `preconfigure_svm_and_update` (Boolean) Indicates whether to preconfigure SVM and perform an update.
- `remove_service` (Boolean) Indicates whether the associated service should be removed.
- `retry` (Boolean) Indicates whether the deployment should be retried.
- `scale_up` (Boolean) Indicates whether the deployment involves scaling up.
- `schedule_date` (String) The date when the deployment is scheduled.
- `services_deployed` (String) Details about the services deployed during the deployment.
- `status` (String) The status of the deployment.
- `teardown` (Boolean) Indicates whether the deployment should be torn down.
- `teardown_after_cancel` (Boolean) Indicates whether teardown should occur after canceling the deployment.
- `template_valid` (Boolean) Indicates if the deployment template is valid.
- `total_batch_count` (String) The total number of batches involved in the deployment.
- `total_num_of_steps` (String) The total number of steps involved in the deployment.
- `update_server_firmware` (Boolean) Indicates whether to update server firmware during the deployment.
- `updated_by` (String) The user who last updated the deployment.
- `updated_date` (String) The date when the deployment was last updated.
- `use_default_catalog` (Boolean) Indicates whether to use the default catalog for the deployment.
- `vds` (Boolean) Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.
- `vms` (Attributes List) List of virtual machines associated with the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--vms))

<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--assigned_users"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.assigned_users`

Read-Only:

- `created_by` (String) The user who created the assigned user.
- `created_date` (String) The date when the assigned user was created.
- `domain_name` (String) The domain name of the assigned user.
- `email` (String) The email address of the assigned user.
- `enabled` (Boolean) Indicates whether the assigned user is enabled.
- `first_name` (String) The first name of the assigned user.
- `group_dn` (String) The distinguished name (DN) of the group associated with the assigned user.
- `group_name` (String) The name of the group associated with the assigned user.
- `id` (String) The unique identifier of the assigned user.
- `last_name` (String) The last name of the assigned user.
- `password` (String) The password associated with the assigned user.
- `phone_number` (String) The phone number of the assigned user.
- `role` (String) The role associated with the assigned user.
- `roles` (List of String) The roles associated with the assigned user.
- `system_user` (Boolean) Indicates whether the assigned user is a system user.
- `update_password` (Boolean) Indicates whether the user password needs to be updated.
- `updated_by` (String) The user who last updated the assigned user.
- `updated_date` (String) The date when the assigned user was last updated.
- `user_name` (String) The username of the assigned user.
- `user_preference` (String) The preferences of the assigned user.
- `user_seq_id` (Number) The sequential ID of the assigned user.


<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_device"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.deployment_device`

Read-Only:

- `brownfield` (Boolean) Indicates whether the deployment device is associated with a brownfield deployment.
- `brownfield_status` (String) The brownfield status of the deployment device.
- `cloud_link` (Boolean) Indicates whether the deployment device has a cloud link.
- `compliant_state` (String) The compliant state of the deployment device.
- `component_id` (String) The component ID associated with the deployment device.
- `current_ip_address` (String) The current IP address of the deployment device.
- `das_cache` (Boolean) Indicates whether the deployment device has Direct-Attached Storage (DAS) cache.
- `device_group_name` (String) The name of the device group associated with the deployment device.
- `device_health` (String) The health status of the deployment device.
- `device_state` (String) The state of the deployment device.
- `device_type` (String) The type of device associated with the deployment device.
- `health_message` (String) The health message associated with the deployment device.
- `ip_address` (String) The IP address of the deployment device.
- `log_dump` (String) The log dump information associated with the deployment device.
- `model` (String) The model of the deployment device.
- `puppet_cert_name` (String) The Puppet certificate name associated with the deployment device.
- `ref_id` (String) The reference ID associated with the deployment device.
- `ref_type` (String) The reference type associated with the deployment device.
- `service_tag` (String) The service tag associated with the deployment device.
- `status` (String) The status of the deployment device.
- `status_end_time` (String) The end time of the status for the deployment device.
- `status_message` (String) The status message associated with the deployment device.
- `status_start_time` (String) The start time of the status for the deployment device.


<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_valid"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.deployment_valid`

Read-Only:

- `messages` (Attributes List) List of messages related to the deployment. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_valid--messages))
- `valid` (Boolean) Indicates whether the deployment is valid.

<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_valid--messages"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.deployment_valid.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.



<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--deployment_validation_response"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.deployment_validation_response`

Read-Only:

- `can_deploy` (Boolean) Indicates whether the deployment can be executed.
- `disk_type_mismatch` (Boolean) Indicates whether there is a disk type mismatch in the deployment.
- `drives_per_storage_pool` (Number) The number of drives per storage pool in the deployment.
- `hostnames` (List of String) A list of hostnames associated with the deployment.
- `max_scalability` (Number) The maximum scalability of the deployment.
- `new_node_disk_types` (List of String) The disk types associated with new nodes in the deployment.
- `no_of_fault_sets` (Number) The number of fault sets in the deployment.
- `nodes` (Number) The number of nodes in the deployment.
- `nodes_per_fault_set` (Number) The number of nodes per fault set in the deployment.
- `number_of_service_volumes` (Number) The number of service volumes in the deployment.
- `protection_domain` (String) The protection domain associated with the deployment.
- `storage_pool_disk_type` (List of String) The disk types associated with each storage pool in the deployment.
- `storage_pools` (Number) The number of storage pools in the deployment.
- `virtual_machines` (Number) The number of virtual machines in the deployment.
- `warning_messages` (List of String) A list of warning messages associated with the deployment validation.


<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--job_details"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.job_details`

Read-Only:

- `component_id` (String) The unique identifier of the component associated with the job.
- `execution_id` (String) The unique identifier of the job execution.
- `level` (String) The log level of the job.
- `message` (String) The log message of the job.
- `timestamp` (String) The timestamp of the job execution.


<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--license_repository"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.license_repository`

Read-Only:

- `created_by` (String) The user who created the license repository.
- `created_date` (String) The date when the license repository was created.
- `disk_location` (String) The disk location of the license repository.
- `filename` (String) The filename associated with the license repository.
- `id` (String) The unique identifier of the license repository.
- `license_data` (String) The license data associated with the license repository.
- `name` (String) The name of the license repository.
- `state` (String) The state of the license repository.
- `type` (String) The type of the license repository.
- `updated_by` (String) The user who last updated the license repository.
- `updated_date` (String) The date when the license repository was last updated.


<a id="nestedatt--resource_group_details--service_template--firmware_repository--deployments--vms"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.deployments.vms`

Read-Only:

- `certificate_name` (String) The certificate name associated with the virtual machine (VM).
- `vm_ipaddress` (String) The IP address of the virtual machine (VM).
- `vm_manufacturer` (String) The manufacturer of the virtual machine (VM).
- `vm_model` (String) The model of the virtual machine (VM).
- `vm_service_tag` (String) The service tag associated with the virtual machine (VM).



<a id="nestedatt--resource_group_details--service_template--firmware_repository--software_bundles"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.software_bundles`

Read-Only:

- `bundle_date` (String) The date when the software bundle was created.
- `bundle_type` (String) The type of the software bundle.
- `created_by` (String) The user who initially created the software bundle.
- `created_date` (String) The date when the software bundle was initially created.
- `criticality` (String) The criticality level of the software bundle.
- `custom` (Boolean) Indicates whether the software bundle is custom.
- `description` (String) A brief description of the software bundle.
- `device_model` (String) The model of the device associated with the software bundle.
- `device_type` (String) The type of device associated with the software bundle.
- `fw_repository_id` (String) The identifier of the firmware repository associated with the software bundle.
- `id` (String) The unique identifier for the software bundle.
- `name` (String) The name of the software bundle.
- `needs_attention` (Boolean) Indicates whether the software bundle needs attention.
- `software_components` (Attributes List) List of software components associated with the software bundle. (see [below for nested schema](#nestedatt--resource_group_details--service_template--firmware_repository--software_bundles--software_components))
- `updated_by` (String) The user who last updated the software bundle.
- `updated_date` (String) The date when the software bundle was last updated.
- `user_bundle` (Boolean) Indicates whether the software bundle is a user-specific bundle.
- `user_bundle_hash_md_5` (String) The MD5 hash value of the user-specific software bundle.
- `user_bundle_path` (String) The path associated with the user-specific software bundle.
- `version` (String) The version of the software bundle.

<a id="nestedatt--resource_group_details--service_template--firmware_repository--software_bundles--software_components"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.software_bundles.software_components`

Read-Only:

- `category` (String) The category to which the component belongs.
- `component_id` (String) The identifier of the component.
- `component_type` (String) The type of the component.
- `created_by` (String) The user who created the component.
- `created_date` (String) The date when the component was created.
- `custom` (Boolean) Indicates whether the component is custom or not.
- `dell_version` (String) The version of the component according to Dell standards.
- `device_id` (String) The identifier of the device associated with the component.
- `firmware_repo_name` (String) The name of the firmware repository associated with the component.
- `hash_md_5` (String) The MD5 hash value of the component.
- `id` (String) The unique identifier for the software component.
- `ignore` (Boolean) Indicates whether the component should be ignored.
- `name` (String) The name of the software component.
- `needs_attention` (Boolean) Indicates whether the component needs attention.
- `operating_system` (String) The operating system associated with the component.
- `original_component_id` (String) The identifier of the original component.
- `original_version` (String) The original version of the component.
- `package_id` (String) The identifier of the package to which the component belongs.
- `path` (String) The path where the component is stored.
- `sub_device_id` (String) The sub-identifier of the device associated with the component.
- `sub_vendor_id` (String) The sub-identifier of the vendor associated with the component.
- `system_ids` (List of String) List of system IDs associated with the component.
- `updated_by` (String) The user who last updated the component.
- `updated_date` (String) The date when the component was last updated.
- `vendor_id` (String) The identifier of the vendor associated with the component.
- `vendor_version` (String) The version of the component according to the vendor's standards.



<a id="nestedatt--resource_group_details--service_template--firmware_repository--software_components"></a>
### Nested Schema for `resource_group_details.service_template.firmware_repository.software_components`

Read-Only:

- `category` (String) The category to which the component belongs.
- `component_id` (String) The identifier of the component.
- `component_type` (String) The type of the component.
- `created_by` (String) The user who created the component.
- `created_date` (String) The date when the component was created.
- `custom` (Boolean) Indicates whether the component is custom or not.
- `dell_version` (String) The version of the component according to Dell standards.
- `device_id` (String) The identifier of the device associated with the component.
- `firmware_repo_name` (String) The name of the firmware repository associated with the component.
- `hash_md_5` (String) The MD5 hash value of the component.
- `id` (String) The unique identifier for the software component.
- `ignore` (Boolean) Indicates whether the component should be ignored.
- `name` (String) The name of the software component.
- `needs_attention` (Boolean) Indicates whether the component needs attention.
- `operating_system` (String) The operating system associated with the component.
- `original_component_id` (String) The identifier of the original component.
- `original_version` (String) The original version of the component.
- `package_id` (String) The identifier of the package to which the component belongs.
- `path` (String) The path where the component is stored.
- `sub_device_id` (String) The sub-identifier of the device associated with the component.
- `sub_vendor_id` (String) The sub-identifier of the vendor associated with the component.
- `system_ids` (List of String) List of system IDs associated with the component.
- `updated_by` (String) The user who last updated the component.
- `updated_date` (String) The date when the component was last updated.
- `vendor_id` (String) The identifier of the vendor associated with the component.
- `vendor_version` (String) The version of the component according to the vendor's standards.



<a id="nestedatt--resource_group_details--service_template--license_repository"></a>
### Nested Schema for `resource_group_details.service_template.license_repository`

Read-Only:

- `created_by` (String) The user who created the license repository.
- `created_date` (String) The date when the license repository was created.
- `disk_location` (String) The disk location of the license repository.
- `filename` (String) The filename associated with the license repository.
- `id` (String) The unique identifier of the license repository.
- `license_data` (String) The license data associated with the license repository.
- `name` (String) The name of the license repository.
- `state` (String) The state of the license repository.
- `type` (String) The type of the license repository.
- `updated_by` (String) The user who last updated the license repository.
- `updated_date` (String) The date when the license repository was last updated.


<a id="nestedatt--resource_group_details--service_template--networks"></a>
### Nested Schema for `resource_group_details.service_template.networks`

Read-Only:

- `description` (String) The description of the network.
- `destination_ip_address` (String) The destination IP address for the network.
- `id` (String) The unique identifier for the network.
- `name` (String) The name of the network.
- `static` (Boolean) Boolean indicating if the network is static.
- `static_network_configuration` (Attributes) The static network configuration settings. (see [below for nested schema](#nestedatt--resource_group_details--service_template--networks--static_network_configuration))
- `type` (String) The type of the network.
- `vlan_id` (Number) The VLAN ID associated with the network.

<a id="nestedatt--resource_group_details--service_template--networks--static_network_configuration"></a>
### Nested Schema for `resource_group_details.service_template.networks.static_network_configuration`

Read-Only:

- `dns_suffix` (String) The DNS suffix for the static network configuration.
- `gateway` (String) The gateway for the static network configuration.
- `ip_address` (String) The IP address associated with the static network configuration.
- `ip_range` (Attributes List) List of IP ranges associated with the static network configuration. (see [below for nested schema](#nestedatt--resource_group_details--service_template--networks--static_network_configuration--ip_range))
- `primary_dns` (String) The primary DNS server for the static network configuration.
- `secondary_dns` (String) The secondary DNS server for the static network configuration.
- `static_route` (Attributes List) List of static routes associated with the static network configuration. (see [below for nested schema](#nestedatt--resource_group_details--service_template--networks--static_network_configuration--static_route))
- `subnet` (String) The subnet for the static network configuration.

<a id="nestedatt--resource_group_details--service_template--networks--static_network_configuration--ip_range"></a>
### Nested Schema for `resource_group_details.service_template.networks.static_network_configuration.ip_range`

Read-Only:

- `ending_ip` (String) The ending IP address of the range.
- `id` (String) The unique identifier for the IP range.
- `role` (String) The role associated with the IP range.
- `starting_ip` (String) The starting IP address of the range.


<a id="nestedatt--resource_group_details--service_template--networks--static_network_configuration--static_route"></a>
### Nested Schema for `resource_group_details.service_template.networks.static_network_configuration.static_route`

Read-Only:

- `destination_ip_address` (String) The IP address of the destination for the static route.
- `static_route_destination_network_id` (String) The ID of the destination network for the static route.
- `static_route_gateway` (String) The gateway for the static route.
- `static_route_source_network_id` (String) The ID of the source network for the static route.
- `subnet_mask` (String) The subnet mask for the static route.




<a id="nestedatt--resource_group_details--service_template--template_valid"></a>
### Nested Schema for `resource_group_details.service_template.template_valid`

Read-Only:

- `messages` (Attributes List) List of messages associated with the template validity. (see [below for nested schema](#nestedatt--resource_group_details--service_template--template_valid--messages))
- `valid` (Boolean) Indicates whether the template is valid.

<a id="nestedatt--resource_group_details--service_template--template_valid--messages"></a>
### Nested Schema for `resource_group_details.service_template.template_valid.messages`

Read-Only:

- `agent_id` (String) The identifier of the agent associated with the message.
- `category` (String) The category or type of the message.
- `correlation_id` (String) The identifier used to correlate related messages.
- `detailed_message` (String) A detailed version or description of the message.
- `display_message` (String) The message to be displayed or shown.
- `id` (String) The unique identifier for the message.
- `message_bundle` (String) The bundle or group to which the message belongs.
- `message_code` (String) The code associated with the message.
- `response_action` (String) The action to be taken in response to the message.
- `sequence_number` (Number) The sequence number of the message in a series.
- `severity` (String) The severity level of the message (e.g., INFO, WARNING, ERROR).
- `time_stamp` (String) The timestamp indicating when the message was generated.




<a id="nestedatt--resource_group_details--vms"></a>
### Nested Schema for `resource_group_details.vms`

Read-Only:

- `certificate_name` (String) The certificate name associated with the virtual machine (VM).
- `vm_ipaddress` (String) The IP address of the virtual machine (VM).
- `vm_manufacturer` (String) The manufacturer of the virtual machine (VM).
- `vm_model` (String) The model of the virtual machine (VM).
- `vm_service_tag` (String) The service tag associated with the virtual machine (VM).


