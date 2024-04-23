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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ServiceReourceSchema - variable holds schema for Service resource
var ServiceReourceSchema schema.Schema = schema.Schema{
	DeprecationMessage:  "Use Resource Group Datasource instead",
	Description:         "This resource is used to manage the Service entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex Service using this resource. We can also Import an existing Service from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Service entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex Service using this resource. We can also Import an existing Service from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"deployment_name": schema.StringAttribute{
			Description:         "Deployment Name",
			MarkdownDescription: "Deployment Name",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"deployment_description": schema.StringAttribute{
			Description:         "Deployment Description",
			MarkdownDescription: "Deployment Description",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"template_id": schema.StringAttribute{
			Description:         "Published Template ID",
			MarkdownDescription: "Published Template ID",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"template_name": schema.StringAttribute{
			Description:         "Service Template Name",
			MarkdownDescription: "Service Template Name",
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"firmware_id": schema.StringAttribute{
			Description:         "Firmware ID",
			MarkdownDescription: "Firmware ID",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"clone_from_host": schema.StringAttribute{
			Description:         "Resource to Duplicate From Host",
			MarkdownDescription: "Resource to Duplicate From Host",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"nodes": schema.Int64Attribute{
			MarkdownDescription: "Number of Nodes",
			Description:         "Number of Nodes",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"deployment_timeout": schema.Int64Attribute{
			MarkdownDescription: "Deployment Timeout, It should be in multiples of 5",
			Description:         "Deployment Timeout, It should be in multiples of 5",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(60),
			Validators: []validator.Int64{
				int64validator.AtLeast(10),
			},
		},
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "Deployment ID",
			MarkdownDescription: "Deployment ID",
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Deployment Status",
			Description:         "Deployment Status",
			Computed:            true,
		},
		"compliant": schema.BoolAttribute{
			MarkdownDescription: "Deployment Compliant Status",
			Description:         "Deployment Compliant Status",
			Computed:            true,
		},
		"servers_in_inventory": schema.StringAttribute{
			MarkdownDescription: "After Delete the Service, Servers in inventory `Keep` or `Remove`.  Default value is `Keep`",
			Description:         "After Delete the Service, Servers in inventory `Keep` or `Remove`.  Default value is `Keep`",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"remove",
				"keep",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("keep"),
			},
		},
		"servers_managed_state": schema.StringAttribute{
			MarkdownDescription: "After Delete the Service, Servers's state `Managed` or `Unmanaged`. Default value is `Unmanaged`.",
			Description:         "After Delete the Service, Servers's state `Managed` or `Unmanaged`. Default value is `Unmanaged`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"managed",
				"unmanaged",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("unmanaged"),
			},
		},
	},
}
