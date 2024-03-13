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

package helper

import (
	"context"
	"strconv"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetServiceState converts scaleiotypes.Template to models.Template
func GetServiceState(input *scaleiotypes.ServiceResponse) models.ServiceResourceModel {
	return models.ServiceResourceModel{
		ID:                    types.StringValue(input.ID),
		DeploymentName:        types.StringValue(input.DeploymentName),
		DeploymentDescription: types.StringValue(input.DeploymentDescription),
		Status:                types.StringValue(input.Status),
		Nodes:                 types.Int64Value(int64(input.ServiceTemplate.ServerCount)),
		Compliant:             types.BoolValue(input.Compliant),
		TemplateID:            types.StringValue(input.ServiceTemplate.OriginalTemplateID),
		TemplateName:          types.StringValue(input.ServiceTemplate.TemplateName),
		FirmwareID:            types.StringValue(input.FirmwareRepository.ID),
	}
}

// UpdateServiceState - function to update state file for Service resource.
func UpdateServiceState(deploymentResponse *scaleiotypes.ServiceResponse, plan models.ServiceResourceModel) (models.ServiceResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := GetServiceState(deploymentResponse)

	state.DeploymentTimeout = plan.DeploymentTimeout

	state.CloneFromHost = plan.CloneFromHost

	if plan.ServersInInventory.ValueString() != "" {
		state.ServersInInventory = plan.ServersInInventory
	} else {
		state.ServersInInventory = types.StringValue("keep")
	}

	if plan.ServersManagedState.ValueString() != "" {
		state.ServersManagedState = plan.ServersManagedState
	} else {
		state.ServersManagedState = types.StringValue("unmanaged")
	}

	return state, diags
}

// HandleServiceDeployment - function to handle service deployment.
func HandleServiceDeployment(ctx context.Context, deploymentResponse *scaleiotypes.ServiceResponse, plan models.ServiceResourceModel, gatewayClient *goscaleio.GatewayClient) (*scaleiotypes.ServiceResponse, diag.Diagnostics) {

	var diags diag.Diagnostics

	deploymentID := deploymentResponse.ID

	couterForStopExecution := 0

	var deploymentTimeout int

	deploymentTimeout, _ = strconv.Atoi(plan.DeploymentTimeout.String())

	deadLineCount := deploymentTimeout / 5

	for couterForStopExecution <= deadLineCount {

		time.Sleep(5 * time.Minute)

		deploymentResponse, err := gatewayClient.GetServiceDetailsByID(deploymentID, true)
		if err != nil {
			diags.AddError(
				"Error in getting service details",
				err.Error(),
			)
			return nil, diags
		}

		tflog.Info(ctx, "Service Deployment Status is ::"+deploymentResponse.Status)

		if deploymentResponse.Status == "complete" {

			tflog.Info(ctx, "Service Details updated to state file successfully")

			return deploymentResponse, diags
		} else if deploymentResponse.Status == "error" {

			var errorMsg string

			for _, details := range deploymentResponse.JobDetails {
				if details.Level == "error" {
					errorMsg += details.Message + "\n"
				}
			}

			if errorMsg != "" {
				diags.AddError("Error in deploying service", errorMsg)
			}

			return nil, diags
		}

		couterForStopExecution++

	}

	diags.AddError("Timed Out For Getting the Deployemnt Status", "Timed Out")

	return nil, diags
}

// GetDataSourceServiceState  is the function to update the state file
func GetDataSourceServiceState(input scaleiotypes.ServiceResponse) models.ServiceModel {
	return models.ServiceModel{
		ID:                           types.StringValue(input.ID),
		DeploymentName:               types.StringValue(input.DeploymentName),
		DeploymentDescription:        types.StringValue(input.DeploymentDescription),
		DeploymentValid:              GetDeploymentValid(input.DeploymentValid),
		Retry:                        types.BoolValue(input.Retry),
		Teardown:                     types.BoolValue(input.Teardown),
		TeardownAfterCancel:          types.BoolValue(input.TeardownAfterCancel),
		RemoveService:                types.BoolValue(input.RemoveService),
		CreatedDate:                  types.StringValue(input.CreatedDate),
		CreatedBy:                    types.StringValue(input.CreatedBy),
		UpdatedDate:                  types.StringValue(input.UpdatedDate),
		UpdatedBy:                    types.StringValue(input.UpdatedBy),
		DeploymentScheduledDate:      types.StringValue(input.DeploymentScheduledDate),
		DeploymentStartedDate:        types.StringValue(input.DeploymentStartedDate),
		DeploymentFinishedDate:       types.StringValue(input.DeploymentFinishedDate),
		ScheduleDate:                 types.StringValue(input.ScheduleDate),
		Status:                       types.StringValue(input.Status),
		Compliant:                    types.BoolValue(input.Compliant),
		ServiceTemplate:              GetTemplateState(input.ServiceTemplate),
		DeploymentDevice:             GetDeploymentDeviceList(input.DeploymentDevice),
		Vms:                          GetVmsList(input.Vms),
		UpdateServerFirmware:         types.BoolValue(input.UpdateServerFirmware),
		UseDefaultCatalog:            types.BoolValue(input.UseDefaultCatalog),
		FirmwareRepository:           GetFirmwareRepository(input.FirmwareRepository),
		FirmwareRepositoryID:         types.StringValue(input.FirmwareRepositoryID),
		LicenseRepository:            GetLicenseRepository(input.LicenseRepository),
		LicenseRepositoryID:          types.StringValue(input.LicenseRepositoryID),
		IndividualTeardown:           types.BoolValue(input.IndividualTeardown),
		DeploymentHealthStatusType:   types.StringValue(input.DeploymentHealthStatusType),
		AssignedUsers:                GetAssignedUsersList(input.AssignedUsers),
		AllUsersAllowed:              types.BoolValue(input.AllUsersAllowed),
		Owner:                        types.StringValue(input.Owner),
		NoOp:                         types.BoolValue(input.NoOp),
		FirmwareInit:                 types.BoolValue(input.FirmwareInit),
		DisruptiveFirmware:           types.BoolValue(input.DisruptiveFirmware),
		PreconfigureSVM:              types.BoolValue(input.PreconfigureSVM),
		PreconfigureSVMAndUpdate:     types.BoolValue(input.PreconfigureSVMAndUpdate),
		ServicesDeployed:             types.StringValue(input.ServicesDeployed),
		PrecalculatedDeviceHealth:    types.StringValue(input.PrecalculatedDeviceHealth),
		LifecycleModeReasons:         GeStringList(input.LifecycleModeReasons),
		JobDetails:                   GetJobDetailsList(input.JobDetails),
		NumberOfDeployments:          types.Int64Value(int64(input.NumberOfDeployments)),
		OperationType:                types.StringValue(input.OperationType),
		OperationStatus:              types.StringValue(input.OperationStatus),
		OperationData:                types.StringValue(input.OperationData),
		DeploymentValidationResponse: GetDeploymentValidationResponse(input.DeploymentValidationResponse),
		CurrentStepCount:             types.StringValue(input.CurrentStepCount),
		TotalNumOfSteps:              types.StringValue(input.TotalNumOfSteps),
		CurrentStepMessage:           types.StringValue(input.CurrentStepMessage),
		CustomImage:                  types.StringValue(input.CustomImage),
		OriginalDeploymentID:         types.StringValue(input.OriginalDeploymentID),
		CurrentBatchCount:            types.StringValue(input.CurrentBatchCount),
		TotalBatchCount:              types.StringValue(input.TotalBatchCount),
		Brownfield:                   types.BoolValue(input.Brownfield),
		OverallDeviceHealth:          types.StringValue(input.OverallDeviceHealth),
		Vds:                          types.BoolValue(input.Vds),
		ScaleUp:                      types.BoolValue(input.ScaleUp),
		LifecycleMode:                types.BoolValue(input.LifecycleMode),
		CanMigratevCLSVMs:            types.BoolValue(input.CanMigratevCLSVMs),
		TemplateValid:                types.BoolValue(input.TemplateValid),
		ConfigurationChange:          types.BoolValue(input.ConfigurationChange),
		DetailMessage:                types.StringValue(input.DetailMessage),
		Timestamp:                    types.StringValue(input.Timestamp),
		Error:                        types.StringValue(input.Error),
		Path:                         types.StringValue(input.Path),
		Messages:                     GetMessagesList(input.Messages),
	}
}
