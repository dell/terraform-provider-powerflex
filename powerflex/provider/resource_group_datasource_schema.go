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

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceGroupDataSourceSchema defines the schema for ResourceGroup datasource
var ResourceGroupDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing ResourceGroup from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing ResourceGroup from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"resource_group_ids": schema.SetAttribute{
			Description:         "List of Resource Group IDs",
			MarkdownDescription: "List of Resource Group IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("service_names"),
				),
			},
		},
		"resource_group_names": schema.SetAttribute{
			Description:         "List of Resource Group names",
			MarkdownDescription: "List of Resource Group names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"resource_group_details": schema.SetNestedAttribute{
			Description:         "Resource Group details",
			MarkdownDescription: "Resource Group details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ResourceGroupResponseSchema()},
		},
	},
}

// ResourceGroupResponseSchema is a function that returns the schema for ResourceGroup details
func ResourceGroupResponseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the Resource Group",
			Description:         "Resource Group ID",
			Computed:            true,
		},
		"deployment_name": schema.StringAttribute{
			MarkdownDescription: "The name of the ResourceGroup",
			Description:         "Resource Group Name",
			Computed:            true,
		},
		"deployment_description": schema.StringAttribute{
			MarkdownDescription: "The description of the ResourceGroup",
			Description:         "Resource Group Description",
			Computed:            true,
		},
		"deployment_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validity of the Resource Group",
			Description:         "Deployment Validity",
			Computed:            true,
			Attributes:          DeploymentValidSchema(),
		},
		"retry": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment should be retried.",
			Description:         "Retry",
			Computed:            true,
		},
		"teardown": schema.BoolAttribute{
			MarkdownDescription: "teardown",
			Description:         "teardown",
			Computed:            true,
		},
		"teardown_after_cancel": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether teardown should occur after canceling the deployment.",
			Description:         "Indicates whether teardown should occur after canceling the deployment.",
			Computed:            true,
		},
		"remove_service": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the associated service should be removed.",
			Description:         "Indicates whether the associated service should be removed.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment was created.",
			Description:         "The date when the deployment was created.",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the deployment.",
			Description:         "The user who created the deployment.",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment was last updated.",
			Description:         "The date when the deployment was last updated.",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the deployment.",
			Description:         "The user who last updated the deployment.",
			Computed:            true,
		},
		"deployment_scheduled_date": schema.StringAttribute{
			MarkdownDescription: "The scheduled date for the deployment.",
			Description:         "The scheduled date for the deployment.",
			Computed:            true,
		},
		"deployment_started_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment started.",
			Description:         "The date when the deployment started.",
			Computed:            true,
		},
		"deployment_finished_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment finished.",
			Description:         "The date when the deployment finished.",
			Computed:            true,
		},
		"schedule_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment is scheduled.",
			Description:         "The date when the deployment is scheduled.",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The status of the deployment.",
			Description:         "The status of the deployment.",
			Computed:            true,
		},
		"compliant": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is compliant.",
			Description:         "Indicates whether the deployment is compliant.",
			Computed:            true,
		},
		"service_template": schema.SingleNestedAttribute{
			MarkdownDescription: "Template details",
			Description:         "Template details",
			Computed:            true,
			Attributes:          TemplateDetailSchema(),
		},
		"deployment_device": schema.ListNestedAttribute{
			MarkdownDescription: "List of devices associated with the deployment.",
			Description:         "List of devices associated with the deployment.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentDeviceSchema()},
		},
		"vms": schema.ListNestedAttribute{
			MarkdownDescription: "List of virtual machines associated with the deployment.",
			Description:         "List of virtual machines associated with the deployment.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VmsSchema()},
		},
		"update_server_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to update server firmware during the deployment.",
			Description:         "Indicates whether to update server firmware during the deployment.",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to use the default catalog for the deployment.",
			Description:         "Indicates whether to use the default catalog for the deployment.",
			Computed:            true,
		},
		"firmware_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the firmware repository used by the template.",
			Description:         "Details about the firmware repository used by the template.",
			Computed:            true,
			Attributes:          FirmwareRepositorySchema(),
		},
		"firmware_repository_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the firmware repository associated with the deployment.",
			Description:         "The ID of the firmware repository associated with the deployment.",
			Computed:            true,
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the license repository associated with the deployment.",
			Description:         "Details about the license repository associated with the deployment.",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"license_repository_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the license repository associated with the deployment.",
			Description:         "The ID of the license repository associated with the deployment.",
			Computed:            true,
		},
		"individual_teardown": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to perform individual teardown for the deployment.",
			Description:         "Indicates whether to perform individual teardown for the deployment.",
			Computed:            true,
		},
		"deployment_health_status_type": schema.StringAttribute{
			MarkdownDescription: "The type of health status associated with the deployment.",
			Description:         "The type of health status associated with the deployment.",
			Computed:            true,
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "List of users assigned to the deployment.",
			Description:         "List of users assigned to the deployment.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether all users are allowed for the deployment.",
			Description:         "Indicates whether all users are allowed for the deployment.",
			Computed:            true,
		},
		"owner": schema.StringAttribute{
			MarkdownDescription: "The owner of the deployment.",
			Description:         "The owner of the deployment.",
			Computed:            true,
		},
		"no_op": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is a no-op (no operation).",
			Description:         "Indicates whether the deployment is a no-op (no operation).",
			Computed:            true,
		},
		"firmware_init": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether firmware initialization is performed during deployment.",
			Description:         "Indicates whether firmware initialization is performed during deployment.",
			Computed:            true,
		},
		"disruptive_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether disruptive firmware actions are allowed.",
			Description:         "Indicates whether disruptive firmware actions are allowed.",
			Computed:            true,
		},
		"preconfigure_svm": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to preconfigure SVM (Storage Virtual Machine).",
			Description:         "Indicates whether to preconfigure SVM (Storage Virtual Machine).",
			Computed:            true,
		},
		"preconfigure_svm_and_update": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to preconfigure SVM and perform an update.",
			Description:         "Indicates whether to preconfigure SVM and perform an update.",
			Computed:            true,
		},
		"services_deployed": schema.StringAttribute{
			MarkdownDescription: "Details about the services deployed during the deployment.",
			Description:         "Details about the services deployed during the deployment.",
			Computed:            true,
		},
		"precalculated_device_health": schema.StringAttribute{
			MarkdownDescription: "The precalculated health of devices associated with the deployment.",
			Description:         "The precalculated health of devices associated with the deployment.",
			Computed:            true,
		},
		"lifecycle_mode_reasons": schema.ListAttribute{
			MarkdownDescription: "List of reasons for the lifecycle mode of the deployment.",
			Description:         "List of reasons for the lifecycle mode of the deployment.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"job_details": schema.ListNestedAttribute{
			MarkdownDescription: "List of job details associated with the deployment.",
			Description:         "List of job details associated with the deployment.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: JobDetailsSchema()},
		},
		"number_of_deployments": schema.Int64Attribute{
			MarkdownDescription: "The total number of deployments.",
			Description:         "The total number of deployments.",
			Computed:            true,
		},
		"operation_type": schema.StringAttribute{
			MarkdownDescription: "The type of operation associated with the deployment.",
			Description:         "The type of operation associated with the deployment.",
			Computed:            true,
		},
		"operation_status": schema.StringAttribute{
			MarkdownDescription: "The status of the operation associated with the deployment.",
			Description:         "The status of the operation associated with the deployment.",
			Computed:            true,
		},
		"operation_data": schema.StringAttribute{
			MarkdownDescription: "Additional data associated with the operation.",
			Description:         "Additional data associated with the operation.",
			Computed:            true,
		},
		"deployment_validation_response": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validation response for the deployment.",
			Description:         "Details about the validation response for the deployment.",
			Computed:            true,
			Attributes:          DeploymentValidationResponseSchema(),
		},
		"current_step_count": schema.StringAttribute{
			MarkdownDescription: "The current step count during deployment.",
			Description:         "The current step count during deployment.",
			Computed:            true,
		},
		"total_num_of_steps": schema.StringAttribute{
			MarkdownDescription: "The total number of steps involved in the deployment.",
			Description:         "The total number of steps involved in the deployment.",
			Computed:            true,
		},
		"current_step_message": schema.StringAttribute{
			MarkdownDescription: "The message associated with the current step during deployment.",
			Description:         "The message associated with the current step during deployment.",
			Computed:            true,
		},
		"custom_image": schema.StringAttribute{
			MarkdownDescription: "The custom image used for deployment.",
			Description:         "The custom image used for deployment.",
			Computed:            true,
		},
		"original_deployment_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the original deployment.",
			Description:         "The ID of the original deployment.",
			Computed:            true,
		},
		"current_batch_count": schema.StringAttribute{
			MarkdownDescription: "The current batch count during deployment.",
			Description:         "The current batch count during deployment.",
			Computed:            true,
		},
		"total_batch_count": schema.StringAttribute{
			MarkdownDescription: "The total number of batches involved in the deployment.",
			Description:         "The total number of batches involved in the deployment.",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment involves brownfield operations.",
			Description:         "Indicates whether the deployment involves brownfield operations.",
			Computed:            true,
		},
		"overall_device_health": schema.StringAttribute{
			MarkdownDescription: "The overall health status of the device in the deployment.",
			Description:         "The overall health status of the device in the deployment.",
			Computed:            true,
		},
		"vds": schema.BoolAttribute{
			MarkdownDescription: "Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.",
			Description:         "Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.",
			Computed:            true,
		},
		"scale_up": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment involves scaling up.",
			Description:         "Indicates whether the deployment involves scaling up.",
			Computed:            true,
		},
		"lifecycle_mode": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is in lifecycle mode.",
			Description:         "Indicates whether the deployment is in lifecycle mode.",
			Computed:            true,
		},
		"can_migratev_clsv_ms": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether migration of cluster virtual machines is allowed.",
			Description:         "Indicates whether migration of cluster virtual machines is allowed.",
			Computed:            true,
		},
		"template_valid": schema.BoolAttribute{
			MarkdownDescription: "Details about the validity of the template.",
			Description:         "Details about the validity of the template.",
			Computed:            true,
		},
		"configuration_change": schema.BoolAttribute{
			MarkdownDescription: "Specifies whether there has been a change in the deployment configuration.",
			Description:         "Specifies whether there has been a change in the deployment configuration.",
			Computed:            true,
		},
		"detail_message": schema.StringAttribute{
			MarkdownDescription: "Detailed Message",
			Description:         "Detailed Message",
			Computed:            true,
		},
		"timestamp": schema.StringAttribute{
			MarkdownDescription: "The timestamp indicating when the message was generated.",
			Description:         "The timestamp indicating when the message was generated.",
			Computed:            true,
		},
		"error": schema.StringAttribute{
			MarkdownDescription: "Error",
			Description:         "Error",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "Path",
			Description:         "Path",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "Messages",
			Description:         "Messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}
