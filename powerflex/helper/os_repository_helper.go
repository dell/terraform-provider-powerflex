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
	"fmt"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetOsRepositoryID gets the OS repository ID
func GetOsRepositoryID(client *goscaleio.System, name string) (string, error) {
	osRepos, err := client.GetAllOSRepositories()
	if err != nil {
		return "", err
	}
	for _, osRepo := range osRepos {
		if osRepo.Name == name {
			return osRepo.ID, nil
		}
	}
	return "", fmt.Errorf("OS repository %s not found", name)
}

// CreateOSRepository creates an OS Repository
func CreateOSRepository(client *goscaleio.System, plan models.OSRepositoryResource) (*scaleiotypes.OSRepository, error) {
	createParam := &scaleiotypes.OSRepository{
		Name:       plan.Name.ValueString(),
		SourcePath: plan.SourcePath.ValueString(),
		RepoType:   plan.RepoType.ValueString(),
		ImageType:  plan.ImageType.ValueString(),
	}

	osRepoResp, err := client.CreateOSRepository(createParam)
	if err != nil {
		return nil, fmt.Errorf("Error creating OS repository: %s ", plan.Name.ValueString())
	}

	return osRepoResp, nil
}

// GetOSRepositoryByID gets the OS repository information by ID
func GetOSRepositoryByID(client *goscaleio.System, id string) (*scaleiotypes.OSRepository, error) {
	osRepo, err := client.GetOSRepositoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error getting OS repository with id: %s ", id)
	}
	return osRepo, nil
}

// UpdateOsRepositoryState updates the OS repository state
func UpdateOsRepositoryState(osRepo *scaleiotypes.OSRepository, plan models.OSRepositoryResource) models.OSRepositoryResource {
	state := plan
	state.ID = types.StringValue(osRepo.ID)
	state.Name = types.StringValue(osRepo.Name)
	state.RCMPath = types.StringValue(osRepo.RCMPath)
	state.BaseURL = types.StringValue(osRepo.BaseURL)
	state.CreatedBy = types.StringValue(osRepo.CreatedBy)
	state.CreatedDate = types.StringValue(osRepo.CreatedDate)
	state.FromWeb = types.BoolValue(osRepo.FromWeb)
	state.RazorName = types.StringValue(osRepo.RazorName)
	state.ImageType = types.StringValue(osRepo.ImageType)
	state.InUse = types.BoolValue(osRepo.InUse)
	state.UserName = types.StringValue(osRepo.UserName)
	state.Password = types.StringValue(osRepo.Password)
	state.RepoType = types.StringValue(osRepo.RepoType)
	state.SourcePath = types.StringValue(osRepo.SourcePath)
	state.State = types.StringValue(osRepo.State)
	return state
}
