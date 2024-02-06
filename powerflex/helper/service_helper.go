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
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		TemplateId:            types.StringValue(input.ServiceTemplate.ID),
		FirmwareId:            types.StringValue(input.FirmwareRepository.ID),
	}
}

// UpdateServiceState - function to update state file for Service resource.
func UpdateServiceState(deploymentResponse *scaleiotypes.ServiceResponse, plan models.ServiceResourceModel) (models.ServiceResourceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	state = GetServiceState(deploymentResponse)

	//state.NumberOfNode = plan.NumberOfNode
	return state, diags
}
