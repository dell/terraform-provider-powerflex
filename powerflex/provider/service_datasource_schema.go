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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TemplateDataSourceSchema defines the schema for template datasource
var ServiceDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing services from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing templates from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"service_ids": schema.SetAttribute{
			Description:         "List of service IDs",
			MarkdownDescription: "List of service IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ConflictsWith(
					path.MatchRoot("service_names"),
				),
			},
		},
		"service_names": schema.SetAttribute{
			Description:         "List of service names",
			MarkdownDescription: "List of service names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ConflictsWith(
					path.MatchRoot("service_ids"),
				),
			},
		},
		"service_details": schema.SetNestedAttribute{
			Description:         "Service details",
			MarkdownDescription: "Service details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ServiceResponseSchema()},
		},
	},
}

// ServiceResponseSchema is a function that returns the schema for ServiceResponse
func ServiceResponseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the service",
			Description:         "Service ID",
			Computed:            true,
		},
		"deployment_name": schema.StringAttribute{
			MarkdownDescription: "The name of the service",
			Description:         "Service Name",
			Computed:            true,
		},
		"deployment_description": schema.StringAttribute{
			MarkdownDescription: "The description of the service",
			Description:         "Service Description",
			Computed:            true,
		},
		"deployment_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validity of the service",
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
			MarkdownDescription: "teardown_after_cancel",
			Description:         "teardown_after_cancel",
			Computed:            true,
		},
		"remove_service": schema.BoolAttribute{
			MarkdownDescription: "remove_service",
			Description:         "remove_service",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"deployment_scheduled_date": schema.StringAttribute{
			MarkdownDescription: "deployment_scheduled_date",
			Description:         "deployment_scheduled_date",
			Computed:            true,
		},
		"deployment_started_date": schema.StringAttribute{
			MarkdownDescription: "deployment_started_date",
			Description:         "deployment_started_date",
			Computed:            true,
		},
		"deployment_finished_date": schema.StringAttribute{
			MarkdownDescription: "deployment_finished_date",
			Description:         "deployment_finished_date",
			Computed:            true,
		},
		"schedule_date": schema.StringAttribute{
			MarkdownDescription: "schedule_date",
			Description:         "schedule_date",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "status",
			Description:         "status",
			Computed:            true,
		},
		"compliant": schema.BoolAttribute{
			MarkdownDescription: "compliant",
			Description:         "compliant",
			Computed:            true,
		},
		"service_template": schema.SingleNestedAttribute{
			MarkdownDescription: "service_template",
			Description:         "service_template",
			Computed:            true,
			Attributes:          TemplateDetailSchema(),
		},
		"deployment_device": schema.ListNestedAttribute{
			MarkdownDescription: "deployment_device",
			Description:         "deployment_device",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentDeviceSchema()},
		},
		"vms": schema.ListNestedAttribute{
			MarkdownDescription: "vms",
			Description:         "vms",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VmsSchema()},
		},
		"update_server_firmware": schema.BoolAttribute{
			MarkdownDescription: "update_server_firmware",
			Description:         "update_server_firmware",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "use_default_catalog",
			Description:         "use_default_catalog",
			Computed:            true,
		},
		"firmware_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "firmware_repository",
			Description:         "firmware_repository",
			Computed:            true,
			Attributes:          FirmwareRepositorySchema(),
		},
		"firmware_repository_id": schema.StringAttribute{
			MarkdownDescription: "firmware_repository_id",
			Description:         "firmware_repository_id",
			Computed:            true,
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "license_repository",
			Description:         "license_repository",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"license_repository_id": schema.StringAttribute{
			MarkdownDescription: "license_repository_id",
			Description:         "license_repository_id",
			Computed:            true,
		},
		"individual_teardown": schema.BoolAttribute{
			MarkdownDescription: "individual_teardown",
			Description:         "individual_teardown",
			Computed:            true,
		},
		"deployment_health_status_type": schema.StringAttribute{
			MarkdownDescription: "deployment_health_status_type",
			Description:         "deployment_health_status_type",
			Computed:            true,
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "assigned_users",
			Description:         "assigned_users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "all_users_allowed",
			Description:         "all_users_allowed",
			Computed:            true,
		},
		"owner": schema.StringAttribute{
			MarkdownDescription: "owner",
			Description:         "owner",
			Computed:            true,
		},
		"no_op": schema.BoolAttribute{
			MarkdownDescription: "no_op",
			Description:         "no_op",
			Computed:            true,
		},
		"firmware_init": schema.BoolAttribute{
			MarkdownDescription: "firmware_init",
			Description:         "firmware_init",
			Computed:            true,
		},
		"disruptive_firmware": schema.BoolAttribute{
			MarkdownDescription: "disruptive_firmware",
			Description:         "disruptive_firmware",
			Computed:            true,
		},
		"preconfigure_svm": schema.BoolAttribute{
			MarkdownDescription: "preconfigure_svm",
			Description:         "preconfigure_svm",
			Computed:            true,
		},
		"preconfigure_svm_and_update": schema.BoolAttribute{
			MarkdownDescription: "preconfigure_svm_and_update",
			Description:         "preconfigure_svm_and_update",
			Computed:            true,
		},
		"services_deployed": schema.StringAttribute{
			MarkdownDescription: "services_deployed",
			Description:         "services_deployed",
			Computed:            true,
		},
		"precalculated_device_health": schema.StringAttribute{
			MarkdownDescription: "precalculated_device_health",
			Description:         "precalculated_device_health",
			Computed:            true,
		},
		"lifecycle_mode_reasons": schema.ListAttribute{
			MarkdownDescription: "lifecycle_mode_reasons",
			Description:         "lifecycle_mode_reasons",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"job_details": schema.ListNestedAttribute{
			MarkdownDescription: "job_details",
			Description:         "job_details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: JobDetailsSchema()},
		},
		"number_of_deployments": schema.Int64Attribute{
			MarkdownDescription: "number_of_deployments",
			Description:         "number_of_deployments",
			Computed:            true,
		},
		"operation_type": schema.StringAttribute{
			MarkdownDescription: "operation_type",
			Description:         "operation_type",
			Computed:            true,
		},
		"operation_status": schema.StringAttribute{
			MarkdownDescription: "operation_status",
			Description:         "operation_status",
			Computed:            true,
		},
		"operation_data": schema.StringAttribute{
			MarkdownDescription: "operation_data",
			Description:         "operation_data",
			Computed:            true,
		},
		"deployment_validation_response": schema.SingleNestedAttribute{
			MarkdownDescription: "deployment_validation_response",
			Description:         "deployment_validation_response",
			Computed:            true,
			Attributes:          DeploymentValidationResponseSchema(),
		},
		"current_step_count": schema.StringAttribute{
			MarkdownDescription: "current_step_count",
			Description:         "current_step_count",
			Computed:            true,
		},
		"total_num_of_steps": schema.StringAttribute{
			MarkdownDescription: "total_num_of_steps",
			Description:         "total_num_of_steps",
			Computed:            true,
		},
		"current_step_message": schema.StringAttribute{
			MarkdownDescription: "current_step_message",
			Description:         "current_step_message",
			Computed:            true,
		},
		"custom_image": schema.StringAttribute{
			MarkdownDescription: "custom_image",
			Description:         "custom_image",
			Computed:            true,
		},
		"original_deployment_id": schema.StringAttribute{
			MarkdownDescription: "original_deployment_id",
			Description:         "original_deployment_id",
			Computed:            true,
		},
		"current_batch_count": schema.StringAttribute{
			MarkdownDescription: "current_batch_count",
			Description:         "current_batch_count",
			Computed:            true,
		},
		"total_batch_count": schema.StringAttribute{
			MarkdownDescription: "total_batch_count",
			Description:         "total_batch_count",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "brownfield",
			Description:         "brownfield",
			Computed:            true,
		},
		"overall_device_health": schema.StringAttribute{
			MarkdownDescription: "overall_device_health",
			Description:         "overall_device_health",
			Computed:            true,
		},
		"vds": schema.BoolAttribute{
			MarkdownDescription: "vds",
			Description:         "vds",
			Computed:            true,
		},
		"scale_up": schema.BoolAttribute{
			MarkdownDescription: "scale_up",
			Description:         "scale_up",
			Computed:            true,
		},
		"lifecycle_mode": schema.BoolAttribute{
			MarkdownDescription: "lifecycle_mode",
			Description:         "lifecycle_mode",
			Computed:            true,
		},
		"can_migratev_clsv_ms": schema.BoolAttribute{
			MarkdownDescription: "can_migratev_clsv_ms",
			Description:         "can_migratev_clsv_ms",
			Computed:            true,
		},
		"template_valid": schema.BoolAttribute{
			MarkdownDescription: "template_valid",
			Description:         "template_valid",
			Computed:            true,
		},
		"configuration_change": schema.BoolAttribute{
			MarkdownDescription: "configuration_change",
			Description:         "configuration_change",
			Computed:            true,
		},
		"detail_message": schema.StringAttribute{
			MarkdownDescription: "detail_message",
			Description:         "detail_message",
			Computed:            true,
		},
		"timestamp": schema.StringAttribute{
			MarkdownDescription: "timestamp",
			Description:         "timestamp",
			Computed:            true,
		},
		"error": schema.StringAttribute{
			MarkdownDescription: "error",
			Description:         "error",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "path",
			Description:         "path",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "messages",
			Description:         "messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}
