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
	"fmt"
	"os"
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dell/goscaleio"
)

// GetCompatibilityManagement Does the read operation for compatibility management
func GetCompatibilityManagement(ctx context.Context, client *goscaleio.System) (*scaleiotypes.CompatibilityManagement, error) {
	return client.GetCompatibilityManagement()
}

// SetCompatibilityManagement Does the read operation for compatibility management
func SetCompatibilityManagement(ctx context.Context, client *goscaleio.System, plan models.CompatibilityManagementDatasourceModel) (*scaleiotypes.CompatibilityManagement, error) {

	// Get the byte array of the file
	data, err := os.ReadFile(plan.RepositoryPath.ValueString())
	if err != nil {
		return nil, fmt.Errorf("Could not read repository file, make sure path to gpg file is correct: %s", err.Error())
	}

	return client.SetCompatibilityManagement(&scaleiotypes.CompatibilityManagementPost{
		CompatibilityDataBytes: data,
		RepositoryPath:         plan.RepositoryPath.ValueString(),
		// Source is always hardcoded to local, that is the only supported value in PowerFlex
		Source: "local",
	})
}

// MapCompatibilityManagementState Does the mapping of compatibility management
func MapCompatibilityManagementState(ctx context.Context, cm *scaleiotypes.CompatibilityManagement) models.CompatibilityManagementDatasourceModel {
	var state models.CompatibilityManagementDatasourceModel
	state.ID = types.StringValue(cm.ID)
	state.AvailabieVersion = types.StringValue(cm.AvailableVersion)
	state.CurrentVersion = types.StringValue(cm.CurrentVersion)
	state.RepositoryPath = types.StringValue(cm.RepositoryPath)
	state.Source = types.StringValue(cm.Source)
	state.CompatibilityData = types.StringValue(cm.CompatibilityData)
	state.CompatibilityDataBytes = types.StringValue(cm.CompatibilityDataBytes)
	return state
}
