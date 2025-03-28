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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// CompatibilityManagementResourceSchema defines the schema for device datasource
var CompatibilityManagementResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to set the compatibility management from the PowerFlex array.",
	MarkdownDescription: "This resource is used to set the compatibility management from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the Compatibility Management Instance.",
			MarkdownDescription: "Unique identifier of the Compatibility Management Instance.",
			Computed:            true,
		},
		"source": schema.StringAttribute{
			Description:         "Source of the Compatibility Management Instance.",
			MarkdownDescription: "Source of the Compatibility Management Instance.",
			Computed:            true,
		},
		"repository_path": schema.StringAttribute{
			Description:         "Repository Path on your local machine to your compatibility gpg file ie: /example/example.gpg.",
			MarkdownDescription: "Repository Path on your local machine to your compatibility gpg file ie: /example/example.gpg.",
			Required:            true,
		},
		"current_version": schema.StringAttribute{
			Description:         "Current Version of the Compatibility Management Instance.",
			MarkdownDescription: "Current Version of the Compatibility Management Instance.",
			Computed:            true,
		},
		"available_version": schema.StringAttribute{
			Description:         "Available Version of the Compatibility Management Instance.",
			MarkdownDescription: "Available Version of the Compatibility Management Instance.",
			Computed:            true,
		},
		"compatibility_data": schema.StringAttribute{
			Description:         "Compatibility Data of the Compatibility Management Instance.",
			MarkdownDescription: "Compatibility Data of the Compatibility Management Instance.",
			Computed:            true,
		},
		"compatibility_data_bytes": schema.StringAttribute{
			Description:         "Compatibility Data Bytes of the Compatibility Management Instance.",
			MarkdownDescription: "Compatibility Data Bytes of the Compatibility Management Instance.",
			Computed:            true,
		},
	},
}
