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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// OsRepositoryResourceSchema - variable holds schema for Os Repository
var OsRepositoryResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the OS Repository entity of the PowerFlex Array. We can Create ,Read and Delete the os image repository using this resource.",
	MarkdownDescription: "This resource is used to manage the OS Repository entity of the PowerFlex Array. We can Create ,Read and Delete the os image repository using this resource.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the OS Repository",
			MarkdownDescription: "ID of the OS Repository",
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
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.OneOf("redhat7", "vmware_esxi", "sles", "windows2016", "windows2019"),
			},
		},
		"source_path": schema.StringAttribute{
			Description:         "Source path of the OS image",
			MarkdownDescription: "Source path of the OS image",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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
			Optional:            true,
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
			Optional:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the OS repository",
			MarkdownDescription: "Name of the OS repository",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"state": schema.StringAttribute{
			Description:         "State of the OS repository",
			MarkdownDescription: "State of the OS repository",
			Computed:            true,
		},
		"repo_type": schema.StringAttribute{
			Description:         "Type of the OS repository. Default is ISO",
			MarkdownDescription: "Type of the OS repository. Default is ISO",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("ISO"),
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
		"timeout": schema.Int64Attribute{
			Description:         "Describes the time in minutes to timeout the job.",
			MarkdownDescription: "Describes the time in minutes to timeout the job.",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(15),
			Validators: []validator.Int64{
				int64validator.AtLeast(15),
			},
		},
	},
}
