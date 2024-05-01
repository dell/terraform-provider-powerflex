/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HostReourceSchema - varible holds schema for Host resource
var HostReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Host entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex Host using this resource. We can also Import an existing Host from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Host entity of PowerFlex Array. We can Create, Update and Delete the PowerFlex Host using this resource. We can also Import an existing Host from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"host_details": hostStateDetailSchema,
		"credential":   connectionSchema,
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "Placeholder",
			MarkdownDescription: "Placeholder",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ip": schema.StringAttribute{
			Required:            true,
			Description:         "IP of the Host",
			MarkdownDescription: "IP of the Host",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"os_family": schema.StringAttribute{
			Required:            true,
			Description:         "OS Family of the Host",
			MarkdownDescription: "OS Family of the Host",
			Validators: []validator.String{
				stringvalidator.OneOf("linux", "windows", "esxi"),
			},
		},
		"name": schema.StringAttribute{
			Required:            true,
			Description:         "Name of the Host",
			MarkdownDescription: "Name of the Host",
		},
		"guid": schema.StringAttribute{
			Optional:            true,
			Description:         "GUID of the Host",
			MarkdownDescription: "GUID of the Host",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"performance_profile": schema.StringAttribute{
			Required:            true,
			Description:         "Performance Profile of the Host",
			MarkdownDescription: "Performance Profile of the Host",
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"HighPerformance",
				"Compact",
			)},
		},
		"package_path": schema.StringAttribute{
			Required:            true,
			Description:         "Package Path of the Host",
			MarkdownDescription: "Package Path of the Host",
		},
		"driver_cfg_path": schema.StringAttribute{
			Optional:            true,
			Description:         "Driver Path of the Host",
			MarkdownDescription: "Driver Path of the Host",
		},
		"mdm_ips": schema.ListAttribute{
			Description:         "List of MDM IPs to be assigned to the SDC.",
			MarkdownDescription: "List of MDM IPs to be assigned to the SDC.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

// hostStateDetailSchema holds details about Host state
var hostStateDetailSchema schema.ListNestedAttribute = schema.ListNestedAttribute{
	Description:         "List of Host state details.",
	Computed:            true,
	MarkdownDescription: "List of Host state details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "IP of the node.",
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "IP of the node.",
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Computed:            true,
				MarkdownDescription: "Operating System on the node",
			},
			"performance_profile": schema.StringAttribute{
				Description:         "Performance Profile of Host.",
				Computed:            true,
				MarkdownDescription: "Performance Profile of Host.",
			},
			"host_id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the Host to manage. This can be retrieved from the Datasource and PowerFlex Server.",
				MarkdownDescription: "ID of the Host to manage. This can be retrieved from the Datasource and PowerFlex Server.",
			},
			"host_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the HOST to manage.",
				MarkdownDescription: "Name of the HOST to manage.",
			},
			"host_guid": schema.StringAttribute{
				Description:         "GUID of the HOST",
				MarkdownDescription: "GUID of the HOST",
				Computed:            true,
			},
			"on_vmware": schema.BoolAttribute{
				Description:         "Is Host on VMware",
				MarkdownDescription: "Is Host on VMware",
				Computed:            true,
			},
			"is_approved": schema.BoolAttribute{
				Description:         "Is Host Approved",
				MarkdownDescription: "Is Host Approved",
				Computed:            true,
			},
			"system_id": schema.StringAttribute{
				Description:         "System ID of the Host",
				MarkdownDescription: "System ID of the Host",
				Computed:            true,
			},
			"mdm_connection_state": schema.StringAttribute{
				Description:         "MDM Connection State",
				MarkdownDescription: "MDM Connection State",
				Computed:            true,
			},
		},
	},
}

// connectionSchema holds details about Connection state
var connectionSchema schema.ListNestedAttribute = schema.ListNestedAttribute{
	Description:         "List of Connection state details.",
	Required:            true,
	MarkdownDescription: "List of Connection state details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"user_name": schema.StringAttribute{
				Required:            true,
				Description:         "UserName of the Connection to manage",
				MarkdownDescription: "UserName of the Connection to manage",
			},
			"password": schema.StringAttribute{
				Description:         "Password of the Connection to manage",
				MarkdownDescription: "Password of the Connection to manage",
				Required:            true,
			},
		},
	},
}
