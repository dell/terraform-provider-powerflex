/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FaultSetResourceSchema - variable holds schema for Fault set
var FaultSetResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Fault Set entity of PowerFlex Array. We can Create, Update and Delete the fault set using this resource. We can also import an existing fault set from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Fault Set entity of PowerFlex Array. We can Create, Update and Delete the fault set using this resource. We can also import an existing fault set from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the Fault Set",
			MarkdownDescription: "ID of the Fault Set",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "ID of the Protection Domain under which the fault set will be created." +
				" Cannot be updated.",
			MarkdownDescription: "ID of the Protection Domain under which the fault set will be created." +
				" Cannot be updated.",
			Required: true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the Fault set",
			MarkdownDescription: "Name of the Fault set",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	},
}
