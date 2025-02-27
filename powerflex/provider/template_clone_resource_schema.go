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
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var pattern = "^.+$"
var regex = regexp.MustCompile(pattern)

// TemplateCloneResourceSchema defines the schema for template datasource
var TemplateCloneResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to clone the existing templates from the PowerFlex array.",
	MarkdownDescription: "This resource is used to clone the existing templates from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"template_name": schema.StringAttribute{
			Description:         "Template Clone Resources's template name.",
			MarkdownDescription: "Template Clone Resources's template name.",
			Computed:            false,
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.RegexMatches(regex, "template_name must not be empty"),
			},
		},
		"original_template_id": schema.StringAttribute{
			Description:         "Template Clone Resources's original template ID.",
			MarkdownDescription: "Template Clone Resources's original template ID.",
			Computed:            false,
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.RegexMatches(regex, "original_template_id must not be empty"),
			},
		},
	},
}
