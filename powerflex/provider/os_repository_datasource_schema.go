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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// OSRepositoryDatasourceSchema - variable holds schema for Os Repository Datasource
var OSRepositoryDatasourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing OS Repository from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing OS Repository from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the OS Repository Datasource",
			MarkdownDescription: "ID of the OS Repository Datasource",
			Computed:            true,
		},
		"os_repositories": schema.ListNestedAttribute{
			Description:         "List of OS Repository Models",
			MarkdownDescription: "List of OS Repository Models",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: OSRepositoryModelSchema,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.OSRepoFilter{})),
		},
	},
}

// OSRepositoryModelSchema  variable holds schema for OS Repository
var OSRepositoryModelSchema map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "ID of the OS Repository ",
		MarkdownDescription: "ID of the OS Repository ",
		Computed:            true,
	},
	"created_date": schema.StringAttribute{
		Description:         "Date of creation of the OS Repository",
		MarkdownDescription: "Date of creation of the OS Repository",
		Computed:            true,
	},
	"image_type": schema.StringAttribute{
		Description:         "Type of the OS image. Supported types are redhat7, vmware_esxi, sles, windows2016, windows2019",
		MarkdownDescription: "Type of the OS image. Supported types are redhat7, vmware_esxi, sles, windows2016, windows2019",
		Computed:            true,
	},
	"source_path": schema.StringAttribute{
		Description:         "Source path of the OS image",
		MarkdownDescription: "Source path of the OS image",
		Computed:            true,
	},
	"razor_name": schema.StringAttribute{
		Description:         "Name of the Razor",
		MarkdownDescription: "Name of the Razor",
		Computed:            true,
	},
	"in_use": schema.BoolAttribute{
		Description:         "Whether the OS repository is in use or not",
		MarkdownDescription: "Whether the OS repository is in use or not",
		Computed:            true,
	},
	"username": schema.StringAttribute{
		Description:         "Username of the OS repository",
		MarkdownDescription: "Username of the OS repository",
		Computed:            true,
	},
	"created_by": schema.StringAttribute{
		Description:         "User who created the OS repository",
		MarkdownDescription: "User who created the OS repository",
		Computed:            true,
	},
	"password": schema.StringAttribute{
		Description:         "Password of the OS repository",
		MarkdownDescription: "Password of the OS repository",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the OS repository",
		MarkdownDescription: "Name of the OS repository",
		Computed:            true,
	},
	"state": schema.StringAttribute{
		Description:         "State of the OS repository",
		MarkdownDescription: "State of the OS repository",
		Computed:            true,
	},
	"repo_type": schema.StringAttribute{
		Description:         "Type of the OS repository. Default is ISO",
		MarkdownDescription: "Type of the OS repository. Default is ISO",
		Computed:            true,
	},
	"rcm_path": schema.StringAttribute{
		Description:         "Path of the RCM",
		MarkdownDescription: "Path of the RCM",
		Computed:            true,
	},
	"base_url": schema.StringAttribute{
		Description:         "Base URL of the OS repository",
		MarkdownDescription: "Base URL of the OS repository",
		Computed:            true,
	},
	"from_web": schema.BoolAttribute{
		Description:         "Whether the OS repository is from the web or not",
		MarkdownDescription: "Whether the OS repository is from the web or not",
		Computed:            true,
	},
	"metadata": schema.SingleNestedAttribute{
		Description:         "Metadata of the OS Repository",
		MarkdownDescription: "Metadata of the OS Repository",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"repos": schema.ListNestedAttribute{
				Description:         "List of OS Repository Metadata Repos",
				MarkdownDescription: "List of OS Repository Metadata Repos",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getOsRepoMetadataReposSchema(),
				},
			},
		},
	},
}

// getOsRepoMetadataReposSchema gets the schema for OS metadata Repos
func getOsRepoMetadataReposSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"base_path": schema.StringAttribute{
			Description:         "Base path of the OS Repository Metadata Repo",
			MarkdownDescription: "Base path of the OS Repository Metadata Repo",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			Description:         "Description of the OS Repository Metadata Repo",
			MarkdownDescription: "Description of the OS Repository Metadata Repo",
			Computed:            true,
		},
		"gpg_key": schema.StringAttribute{
			Description:         "GPG key of the OS Repository Metadata Repo",
			MarkdownDescription: "GPG key of the OS Repository Metadata Repo",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the OS Repository Metadata Repo",
			MarkdownDescription: "Name of the OS Repository Metadata Repo",
			Computed:            true,
		},
		"os_packages": schema.BoolAttribute{
			Description:         "Whether the OS Repository Metadata Repo has OS packages or not",
			MarkdownDescription: "Whether the OS Repository Metadata Repo has OS packages or not",
			Computed:            true,
		},
	}
}
