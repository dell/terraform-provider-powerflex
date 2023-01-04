package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description: "Fetches the list of Storagepool",
	Attributes: map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Description:         "Last Updated",
			MarkdownDescription: "Last Updated",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			Description:         "Gets the ID of Storagepool",
			MarkdownDescription: "Gets the ID of Storagepool",
			Computed:            true,
		},
		"systemid": schema.StringAttribute{
			Description:         "Gets the System ID",
			MarkdownDescription: "Gets the System ID",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "Gets the Protection Domain ID for Storagepool",
			MarkdownDescription: "Gets the Protection Domain ID for Storagepool",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description:         "Gets the Protection Domain Name for Storagepool",
			MarkdownDescription: "Gets the Protection Domain Name for Storagepool",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"name": schema.StringAttribute{
			Description:         "Returns the Name of Storagepool",
			MarkdownDescription: "Returns the Name of Storagepool",
			Required:            true,
		},
		"media_type": schema.StringAttribute{
			Description:         "Gets the Media Type",
			MarkdownDescription: "Gets the Media Type",
			Required:            true,
		},
		"use_rmcache": schema.BoolAttribute{
			Description:         "Gets the Read RAM Cache",
			MarkdownDescription: "Gets the Read RAM Cache",
			Optional:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description:         "Gets the RFCache",
			MarkdownDescription: "Gets the RFCache",
			Optional:            true,
		},
		"links": schema.ListNestedAttribute{
			Description:         "Specifies the links asscociated with Storagepool",
			MarkdownDescription: "Specifies the links asscociated with Storagepool",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description:         "Specifies the relationship with the Storagepool",
						MarkdownDescription: "Specifies the relationship with the Storagepool",
						Computed:            true,
					},
					"href": schema.StringAttribute{
						Description:         "Specifies the exact path to fetch the details",
						MarkdownDescription: "Specifies the exact path to fetch the details",
						Computed:            true,
					},
				},
			},
		},
	},
}
