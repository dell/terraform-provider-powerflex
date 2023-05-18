package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SDCExpansionResourceSchema - varible holds schema for SDC Expansion
var SDCExpansionResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to add the SDC.",
	MarkdownDescription: "This resource can be used to add the SDC.",
	Attributes: map[string]schema.Attribute{
		"csv_detail": csvSchema,
		"mdm_ip": schema.StringAttribute{
			Description:         "The JSON data which is being received after parsing the csv.",
			MarkdownDescription: "The JSON data which is being received after parsing the csv.",
			Required:            true,
		},
		"mdm_password": schema.StringAttribute{
			Description:         "The JSON data which is being received after parsing the csv.",
			MarkdownDescription: "The JSON data which is being received after parsing the csv.",
			Required:            true,
		},
		"lia_password": schema.StringAttribute{
			Description:         "The JSON data which is being received after parsing the csv.",
			MarkdownDescription: "The JSON data which is being received after parsing the csv.",
			Required:            true,
		},
		"installed_sdc_ips": schema.StringAttribute{
			Description:         "List of installed SDC IPs",
			Computed:            true,
			MarkdownDescription: "List of installed SDC IPs",
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the package.",
			Computed:            true,
			MarkdownDescription: "The ID of the package.",
		},
	},
}

// csvSchema - varible holds schema for CSV Param Details
var csvSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	Required:            true,
	MarkdownDescription: "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "ip of the node",
				Required:            true,
				MarkdownDescription: "ip of the node",
			},
			"password": schema.StringAttribute{
				Description:         "Password on the node",
				Required:            true,
				MarkdownDescription: "Password on the node",
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Required:            true,
				MarkdownDescription: "Operating System on the node",
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "Whether this works as MDM or Tie Breaker",
				Required:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker",
			},
			"is_sdc": schema.StringAttribute{
				Description:         "whether this node is SDC or not",
				Required:            true,
				MarkdownDescription: "whether this node is SDC or not",
			},
		},
	},
}
