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

package provider

import (
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ResourceCredentialDataSourceSchema defines the schema for ResourceCredential datasource
var ResourceCredentialDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to read the Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.",
	MarkdownDescription: "This datasource is used to read the Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "default datasource id",
			MarkdownDescription: "default datasource id",
			Computed:            true,
		},
		"resource_credential_details": schema.ListNestedAttribute{
			Description:         "List of Resource Credentials",
			MarkdownDescription: "List of Resource Credentials",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: resourceCredentialDataAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.ResourceCredentialFilter{})),
		},
	},
}

var resourceCredentialDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "Unique identifier of the resource credential instance.",
		MarkdownDescription: "Unique identifier of the resource credential instance.",
		Computed:            true,
	},
	"type": schema.StringAttribute{
		Description:         "Resource Credential type",
		MarkdownDescription: "Resource Credential type",
		Computed:            true,
	},
	"created_date": schema.StringAttribute{
		Description:         "Resource Credential created date",
		MarkdownDescription: "Resource Credential created date",
		Computed:            true,
	},
	"created_by": schema.StringAttribute{
		Description:         "Who the Resource Credential was created by",
		MarkdownDescription: "Who the Resource Credential was created by",
		Computed:            true,
	},
	"updated_by": schema.StringAttribute{
		Description:         "Who the Resource Credential was last updated by",
		MarkdownDescription: "Who the Resource Credential was last updated by",
		Computed:            true,
	},
	"updated_date": schema.StringAttribute{
		Description:         "Resource Credential updated date",
		MarkdownDescription: "Resource Credential updated date",
		Computed:            true,
	},
	"label": schema.StringAttribute{
		Description:         "Resource Credential label",
		MarkdownDescription: "Resource Credential label",
		Computed:            true,
	},
	"domain": schema.StringAttribute{
		Description:         "Resource Credential domain",
		MarkdownDescription: "Resource Credential domain",
		Computed:            true,
	},
	"username": schema.StringAttribute{
		Description:         "Resource Credential username",
		MarkdownDescription: "Resource Credential username",
		Computed:            true,
	},
}
