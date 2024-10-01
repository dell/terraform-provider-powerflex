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

// OSRepositoryResource is the model for OS Repository Resource
type OSRepositoryResource struct {
	ID          types.String `tfsdk:"id"`
	CreatedDate types.String `tfsdk:"created_date"`
	ImageType   types.String `tfsdk:"image_type"`
	SourcePath  types.String `tfsdk:"source_path"`
	RazorName   types.String `tfsdk:"razor_name"`
	InUse       types.Bool   `tfsdk:"in_use"`
	UserName    types.String `tfsdk:"username"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Password    types.String `tfsdk:"password"`
	Name        types.String `tfsdk:"name"`
	State       types.String `tfsdk:"state"`
	RepoType    types.String `tfsdk:"repo_type"`
	RCMPath     types.String `tfsdk:"rcm_path"`
	BaseURL     types.String `tfsdk:"base_url"`
	FromWeb     types.Bool   `tfsdk:"from_web"`
	Timeout     types.Int64  `tfsdk:"timeout"`
}

// OSRepositoryDataSource defines the model for Os Repository Datasource
type OSRepositoryDataSource struct {
	ID             types.String        `tfsdk:"id"`
	OSRepositories []OSRepositoryModel `tfsdk:"os_repositories"`
	// filter
	OSRepoFilter *OSRepoFilter `tfsdk:"filter"`
}

// OSRepoFilter defines the model for filters used for OSRepositoryDataSource
type OSRepoFilter struct {
	ID          []types.String `tfsdk:"id"`
	CreatedDate []types.String `tfsdk:"created_date"`
	ImageType   []types.String `tfsdk:"image_type"`
	SourcePath  []types.String `tfsdk:"source_path"`
	RazorName   []types.String `tfsdk:"razor_name"`
	InUse       types.Bool     `tfsdk:"in_use"`
	UserName    []types.String `tfsdk:"username"`
	CreatedBy   []types.String `tfsdk:"created_by"`
	Name        []types.String `tfsdk:"name"`
	State       []types.String `tfsdk:"state"`
	RepoType    []types.String `tfsdk:"repo_type"`
	RCMPath     []types.String `tfsdk:"rcm_path"`
	BaseURL     []types.String `tfsdk:"base_url"`
	FromWeb     types.Bool     `tfsdk:"from_web"`
}

// OSRepositoryModel is the model for OS Repository Model
type OSRepositoryModel struct {
	ID          types.String              `tfsdk:"id"`
	CreatedDate types.String              `tfsdk:"created_date"`
	ImageType   types.String              `tfsdk:"image_type"`
	SourcePath  types.String              `tfsdk:"source_path"`
	RazorName   types.String              `tfsdk:"razor_name"`
	InUse       types.Bool                `tfsdk:"in_use"`
	UserName    types.String              `tfsdk:"username"`
	CreatedBy   types.String              `tfsdk:"created_by"`
	Password    types.String              `tfsdk:"password"`
	Name        types.String              `tfsdk:"name"`
	State       types.String              `tfsdk:"state"`
	RepoType    types.String              `tfsdk:"repo_type"`
	RCMPath     types.String              `tfsdk:"rcm_path"`
	BaseURL     types.String              `tfsdk:"base_url"`
	FromWeb     types.Bool                `tfsdk:"from_web"`
	Metadata    OSRepositoryMetadataModel `tfsdk:"metadata"`
}

// OSRepositoryMetadataModel is the model for OS Repository Metadata
type OSRepositoryMetadataModel struct {
	Repos []OSRepositoryMetadataRepoModel `tfsdk:"repos"`
}

// OSRepositoryMetadataRepoModel is the model for OSRepositoryMetadataRepo
type OSRepositoryMetadataRepoModel struct {
	BasePath    types.String `tfsdk:"base_path"`
	Description types.String `tfsdk:"description"`
	GPGKey      types.String `tfsdk:"gpg_key"`
	Name        types.String `tfsdk:"name"`
	OSPackages  types.Bool   `tfsdk:"os_packages"`
}
