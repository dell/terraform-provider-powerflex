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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// GetAllOsRepositories gets all OS repositories
func GetAllOsRepositories(client *goscaleio.System) ([]scaleiotypes.OSRepository, error) {
	osRepos, err := client.GetAllOSRepositories()
	if err != nil {
		return nil, fmt.Errorf("Error getting OS repositories")
	}
	// Filter ISO repos
	var filteredOSRepos []scaleiotypes.OSRepository
	for _, osRepo := range osRepos {
		if osRepo.RepoType == "ISO" {
			filteredOSRepos = append(filteredOSRepos, osRepo)
		}
	}
	return filteredOSRepos, nil
}

// GetOSRepositoriesByNames gets the OS Repositories filtered by names
func GetOSRepositoriesByNames(client *goscaleio.System, names []basetypes.StringValue) ([]scaleiotypes.OSRepository, error) {
	// Get all OS repositories from the client
	osRepos, err := client.GetAllOSRepositories()
	if err != nil {
		return nil, fmt.Errorf("error getting all OS repositories: %w", err)
	}

	var filteredRepos []scaleiotypes.OSRepository
	nameSet := make(map[string]struct{})

	// create name set
	for _, name := range names {
		nameSet[name.ValueString()] = struct{}{}
	}

	// Filter the repositories based on the names provided
	for _, osRepo := range osRepos {
		if _, exists := nameSet[osRepo.Name]; exists {
			filteredRepos = append(filteredRepos, osRepo)
		}
	}

	return filteredRepos, nil
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

// GetAllOSRepositoryState sets the state for the OS repository datasource.
func GetAllOSRepositoryState(input scaleiotypes.OSRepository) models.OSRepositoryModel {
	return models.OSRepositoryModel{
		ID:          types.StringValue(input.ID),
		Name:        types.StringValue(input.Name),
		RCMPath:     types.StringValue(input.RCMPath),
		BaseURL:     types.StringValue(input.BaseURL),
		CreatedBy:   types.StringValue(input.CreatedBy),
		CreatedDate: types.StringValue(input.CreatedDate),
		FromWeb:     types.BoolValue(input.FromWeb),
		RazorName:   types.StringValue(input.RazorName),
		ImageType:   types.StringValue(input.ImageType),
		InUse:       types.BoolValue(input.InUse),
		UserName:    types.StringValue(input.UserName),
		Password:    types.StringValue(input.Password),
		RepoType:    types.StringValue(input.RepoType),
		SourcePath:  types.StringValue(input.SourcePath),
		State:       types.StringValue(input.State),
		Metadata:    UpdateOsRepoMetadata(input.Metadata),
	}
}

// UpdateOsRepoMetadata sets the state for osRepoMetadata
func UpdateOsRepoMetadata(osRepoMetadata scaleiotypes.OSRepositoryMetadata) models.OSRepositoryMetadataModel {
	state := models.OSRepositoryMetadataModel{}
	for _, repo := range osRepoMetadata.Repos {
		repoModel := models.OSRepositoryMetadataRepoModel{
			BasePath:    types.StringValue(repo.BasePath),
			Description: types.StringValue(repo.Description),
			GPGKey:      types.StringValue(repo.GPGKey),
			Name:        types.StringValue(repo.Name),
			OSPackages:  types.BoolValue(repo.OSPackages),
		}
		state.Repos = append(state.Repos, repoModel)
	}
	return state
}
