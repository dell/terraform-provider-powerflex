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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceGroupResourceModel is the tfsdk model of Resource Group Response
type ResourceGroupResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	DeploymentName        types.String `tfsdk:"deployment_name"`
	DeploymentDescription types.String `tfsdk:"deployment_description"`
	TemplateID            types.String `tfsdk:"template_id"`
	FirmwareID            types.String `tfsdk:"firmware_id"`
	CloneFromHost         types.String `tfsdk:"clone_from_host"`
	Nodes                 types.Int64  `tfsdk:"nodes"`
	Status                types.String `tfsdk:"status"`
	Compliant             types.Bool   `tfsdk:"compliant"`
	TemplateName          types.String `tfsdk:"template_name"`
	DeploymentTimeout     types.Int64  `tfsdk:"deployment_timeout"`
	ServersInInventory    types.String `tfsdk:"servers_in_inventory"`
	ServersManagedState   types.String `tfsdk:"servers_managed_state"`
}

// ResourceGroupDataSourceModel is the tfsdk model of ResourceGroup data source schema
type ResourceGroupDataSourceModel struct {
	ResourceGroupIDs     types.Set            `tfsdk:"resource_group_ids"`
	ResourceGroupNames   types.Set            `tfsdk:"resource_group_names"`
	ResourceGroupDetails []ResourceGroupModel `tfsdk:"resource_group_details"`
	ID                   types.String         `tfsdk:"id"`
}

// ResourceGroupModel is the tfsdk model of ResourceGroup Details
type ResourceGroupModel struct {
	ID                           types.String                 `tfsdk:"id"`
	DeploymentName               types.String                 `tfsdk:"deployment_name"`
	DeploymentDescription        types.String                 `tfsdk:"deployment_description"`
	DeploymentValid              DeploymentValid              `tfsdk:"deployment_valid"`
	Retry                        types.Bool                   `tfsdk:"retry"`
	Teardown                     types.Bool                   `tfsdk:"teardown"`
	TeardownAfterCancel          types.Bool                   `tfsdk:"teardown_after_cancel"`
	RemoveService                types.Bool                   `tfsdk:"remove_service"`
	CreatedDate                  types.String                 `tfsdk:"created_date"`
	CreatedBy                    types.String                 `tfsdk:"created_by"`
	UpdatedDate                  types.String                 `tfsdk:"updated_date"`
	UpdatedBy                    types.String                 `tfsdk:"updated_by"`
	DeploymentScheduledDate      types.String                 `tfsdk:"deployment_scheduled_date"`
	DeploymentStartedDate        types.String                 `tfsdk:"deployment_started_date"`
	DeploymentFinishedDate       types.String                 `tfsdk:"deployment_finished_date"`
	ScheduleDate                 types.String                 `tfsdk:"schedule_date"`
	Status                       types.String                 `tfsdk:"status"`
	Compliant                    types.Bool                   `tfsdk:"compliant"`
	ServiceTemplate              TemplateModel                `tfsdk:"service_template"`
	DeploymentDevice             []DeploymentDevice           `tfsdk:"deployment_device"`
	Vms                          []Vms                        `tfsdk:"vms"`
	UpdateServerFirmware         types.Bool                   `tfsdk:"update_server_firmware"`
	UseDefaultCatalog            types.Bool                   `tfsdk:"use_default_catalog"`
	FirmwareRepository           FirmwareRepository           `tfsdk:"firmware_repository"`
	FirmwareRepositoryID         types.String                 `tfsdk:"firmware_repository_id"`
	LicenseRepository            LicenseRepository            `tfsdk:"license_repository"`
	LicenseRepositoryID          types.String                 `tfsdk:"license_repository_id"`
	IndividualTeardown           types.Bool                   `tfsdk:"individual_teardown"`
	DeploymentHealthStatusType   types.String                 `tfsdk:"deployment_health_status_type"`
	AssignedUsers                []AssignedUsers              `tfsdk:"assigned_users"`
	AllUsersAllowed              types.Bool                   `tfsdk:"all_users_allowed"`
	Owner                        types.String                 `tfsdk:"owner"`
	NoOp                         types.Bool                   `tfsdk:"no_op"`
	FirmwareInit                 types.Bool                   `tfsdk:"firmware_init"`
	DisruptiveFirmware           types.Bool                   `tfsdk:"disruptive_firmware"`
	PreconfigureSVM              types.Bool                   `tfsdk:"preconfigure_svm"`
	PreconfigureSVMAndUpdate     types.Bool                   `tfsdk:"preconfigure_svm_and_update"`
	ServicesDeployed             types.String                 `tfsdk:"services_deployed"`
	PrecalculatedDeviceHealth    types.String                 `tfsdk:"precalculated_device_health"`
	LifecycleModeReasons         []types.String               `tfsdk:"lifecycle_mode_reasons"`
	JobDetails                   []JobDetails                 `tfsdk:"job_details"`
	NumberOfDeployments          types.Int64                  `tfsdk:"number_of_deployments"`
	OperationType                types.String                 `tfsdk:"operation_type"`
	OperationStatus              types.String                 `tfsdk:"operation_status"`
	OperationData                types.String                 `tfsdk:"operation_data"`
	DeploymentValidationResponse DeploymentValidationResponse `tfsdk:"deployment_validation_response"`
	CurrentStepCount             types.String                 `tfsdk:"current_step_count"`
	TotalNumOfSteps              types.String                 `tfsdk:"total_num_of_steps"`
	CurrentStepMessage           types.String                 `tfsdk:"current_step_message"`
	CustomImage                  types.String                 `tfsdk:"custom_image"`
	OriginalDeploymentID         types.String                 `tfsdk:"original_deployment_id"`
	CurrentBatchCount            types.String                 `tfsdk:"current_batch_count"`
	TotalBatchCount              types.String                 `tfsdk:"total_batch_count"`
	Brownfield                   types.Bool                   `tfsdk:"brownfield"`
	OverallDeviceHealth          types.String                 `tfsdk:"overall_device_health"`
	Vds                          types.Bool                   `tfsdk:"vds"`
	ScaleUp                      types.Bool                   `tfsdk:"scale_up"`
	LifecycleMode                types.Bool                   `tfsdk:"lifecycle_mode"`
	CanMigratevCLSVMs            types.Bool                   `tfsdk:"can_migratev_clsv_ms"`
	TemplateValid                types.Bool                   `tfsdk:"template_valid"`
	ConfigurationChange          types.Bool                   `tfsdk:"configuration_change"`
	DetailMessage                types.String                 `tfsdk:"detail_message"`
	Timestamp                    types.String                 `tfsdk:"timestamp"`
	Error                        types.String                 `tfsdk:"error"`
	Path                         types.String                 `tfsdk:"path"`
	Messages                     []Messages                   `tfsdk:"messages"`
}
