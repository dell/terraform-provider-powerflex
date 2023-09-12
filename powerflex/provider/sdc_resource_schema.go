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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SDCReourceSchema - varible holds schema for SDC resource
var SDCReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to Manage the SDC in PowerFlex Cluster.",
	MarkdownDescription: "This resource can be used to Manage the SDC in PowerFlex Cluster.",
	Attributes: map[string]schema.Attribute{
		"sdc_details": sdcDetailSchema,
		"mdm_password": schema.StringAttribute{
			Description:         "MDM Password to connect MDM Server.",
			MarkdownDescription: "MDM Password to connect MDM Server.",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"lia_password": schema.StringAttribute{
			Description:         "LIA Password to connect MDM Server.",
			MarkdownDescription: "LIA Password to connect MDM Server.",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "Placeholder",
			MarkdownDescription: "Placeholder",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

// sdcDetailSchema - variable holds schema for CSV Param Details
var sdcDetailSchema schema.ListNestedAttribute = schema.ListNestedAttribute{
	Description:         "List of SDC Expansion Server Details.",
	Optional:            true,
	Computed:            true,
	MarkdownDescription: "List of SDC Expansion Server Details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "IP of the node. Conflict with `sdc_id`",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "IP of the node. Conflict with `sdc_id`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_id")),
				},
			},
			"username": schema.StringAttribute{
				Description:         "Username of the node",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Username of the node",
				PlanModifiers: []planmodifier.String{
					helper.StringDefault("root"),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Description:         "Password of the node",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Password of the node",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.AlsoRequires(path.MatchRoot("lia_password")),
					stringvalidator.AlsoRequires(path.MatchRoot("mdm_password")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Operating System on the node",
				PlanModifiers: []planmodifier.String{
					helper.StringDefault("linux"),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "Whether this works as MDM or Tie Breaker,The acceptable value are `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker,The acceptable value are `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_sdc": schema.StringAttribute{
				Description:         "Whether this node is to operate as an SDC or not. The acceptable values are `Yes` and `No`. Default value is `Yes`.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this node is to operate as an SDC or not. The acceptable values are `Yes` and `No`. Default value is `Yes`.",
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					"Yes",
					"No",
				)},
				PlanModifiers: []planmodifier.String{
					helper.StringDefault("Yes"),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"performance_profile": schema.StringAttribute{
				Description:         "Performance Profile of SDC, The acceptable value are `HighPerformance` or `Compact`.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Performance Profile of SDC, The acceptable value are `HighPerformance` or `Compact`.",
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					"HighPerformance",
					"Compact",
				)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sdc_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "ID of the SDC to manage. This can be retrieved from the Datasource and PowerFlex Server. Cannot be updated. Conflict with `ip`",
				MarkdownDescription: "ID of the SDC to manage. This can be retrieved from the Datasource and PowerFlex Server. Cannot be updated. Conflict with `ip`",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("ip")),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Name of the SDC to manage.",
				MarkdownDescription: "Name of the SDC to manage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(31),
				},
			},
			"sdc_guid": schema.StringAttribute{
				Description:         models.SdcResourceSchemaDescriptions.SdcGUID,
				MarkdownDescription: models.SdcResourceSchemaDescriptions.SdcGUID,
				Computed:            true,
			},
			"on_vmware": schema.BoolAttribute{
				Description:         models.SdcResourceSchemaDescriptions.OnVMWare,
				MarkdownDescription: models.SdcResourceSchemaDescriptions.OnVMWare,
				Computed:            true,
			},
			"sdc_approved": schema.BoolAttribute{
				Description:         models.SdcResourceSchemaDescriptions.SdcApproved,
				MarkdownDescription: models.SdcResourceSchemaDescriptions.SdcApproved,
				Computed:            true,
			},
			"system_id": schema.StringAttribute{
				Description:         models.SdcResourceSchemaDescriptions.SystemID,
				MarkdownDescription: models.SdcResourceSchemaDescriptions.SystemID,
				Computed:            true,
			},
			"mdm_connection_state": schema.StringAttribute{
				Description:         models.SdcResourceSchemaDescriptions.MdmConnectionState,
				MarkdownDescription: models.SdcResourceSchemaDescriptions.MdmConnectionState,
				Computed:            true,
			},
		},
	},
}
