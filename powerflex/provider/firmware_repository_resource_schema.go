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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FirmwareRepositoryResourceSchema - variable holds schema for Firmware Repository
var FirmwareRepositoryResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Firmware Repository entity of PowerFlex Array. We can Create and Update the firmware repository using this resource. As part of create operation, we can upload the compliance file and as part of update we can approve the unsigned files.We can also import an existing firmware repository from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Firmware Repository entity of PowerFlex Array. We can Create and Update the firmware repository using this resource. As part of create operation, we can upload the compliance file and as part of update we can approve the unsigned files.We can also import an existing firmware repository from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the Firmware Repository",
			MarkdownDescription: "ID of the Firmware Repository",
			Computed:            true,
		},
		"source_location": schema.StringAttribute{
			Description: "Specfies the path from where it will download the compliance file." +
				" Cannot be updated.",
			MarkdownDescription: "Specfies the path from where it will download the compliance file." +
				" Cannot be updated.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"username": schema.StringAttribute{
			Description: "Username is only used if specifying a CIFS share" +
				" Cannot be updated.",
			MarkdownDescription: "Username is only used if specifying a CIFS share" +
				" Cannot be updated.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.AlsoRequires(
					path.MatchRoot("password"),
				),
			},
		},
		"password": schema.StringAttribute{
			Description: "Password is only used if specifying a CIFS share" +
				" Cannot be updated.",
			MarkdownDescription: "Password is only used if specifying a CIFS share" +
				" Cannot be updated.",
			Optional:  true,
			Sensitive: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.AlsoRequires(
					path.MatchRoot("username"),
				),
			},
		},
		"approve": schema.BoolAttribute{
			Description:         "Whether to approve the unsigned file or not.",
			MarkdownDescription: "Whether to approve the unsigned file or not.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"timeout": schema.Int64Attribute{
			Description:         "Describes the time in minutes to timeout the job.",
			MarkdownDescription: "Describes the time in minutes to timeout the job.",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(40),
			Validators: []validator.Int64{
				int64validator.AtLeast(40),
			},
		},
		"name": schema.StringAttribute{
			Description:         "Name of the Firmware Repository",
			MarkdownDescription: "Name of the Firmware Repository",
			Computed:            true,
		},
		"disk_location": schema.StringAttribute{
			Description:         "Disk location of the Firmware Repository",
			MarkdownDescription: "Disk location of the Firmware Repository",
			Computed:            true,
		},
		"file_name": schema.StringAttribute{
			Description:         "File Name",
			MarkdownDescription: "File Name",
			Computed:            true,
		},
		"default_catalog": schema.BoolAttribute{
			Description:         "Whether this Firmware Repository is set to default or not.",
			MarkdownDescription: "Whether this Firmware Repository is set to default or not.",
			Computed:            true,
		},
	},
}
