/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetResourceCredentials returns all the resource credentials
func GetResourceCredentials(ctx context.Context, system *goscaleio.System) ([]scaleiotypes.CredObj, error) {
	var resourceCredentials []scaleiotypes.CredObj
	creds, err := system.GetResourceCredentials()
	if err != nil {
		return nil, err
	}
	for _, val := range creds.Credentials {
		resourceCredentials = append(resourceCredentials, val.Credential)
	}
	return resourceCredentials, nil
}

// MapResourceCredentials a terraform mapped resource
func MapResourceCredentials(rcs []scaleiotypes.CredObj, state models.ResourceCredentialDataSourceModel) models.ResourceCredentialDataSourceModel {
	var resourceCredentialDetails []models.ResourceCredentialModel
	for _, val := range rcs {
		resourceCredentialDetails = append(resourceCredentialDetails,
			models.ResourceCredentialModel{
				ID:          types.StringValue(val.ID),
				Type:        types.StringValue(val.Type),
				CreateDate:  types.StringValue(val.CreateDate),
				CreatedBy:   types.StringValue(val.CreatedBy),
				UpdatedBy:   types.StringValue(val.UpdatedBy),
				UpdatedDate: types.StringValue(val.UpdatedDate),
				Label:       types.StringValue(val.Label),
				Domain:      types.StringValue(val.Domain),
				Username:    types.StringValue(val.Username),
			})
	}
	state.ResourceCredentialDetails = resourceCredentialDetails
	return state
}
