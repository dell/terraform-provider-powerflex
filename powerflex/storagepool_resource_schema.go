package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description: "Manages storage pool resource",
	Attributes: map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Description:         "Last Updated",
			MarkdownDescription: "Last Updated",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			Description:         "ID of the Storage pool",
			MarkdownDescription: "ID of the Storage pool",
			Computed:            true,
		},
		"systemid": schema.StringAttribute{
			Description:         "ID of the system",
			MarkdownDescription: "ID of the system",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "ID of the Protection domain",
			MarkdownDescription: "ID of the Protection domain",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description:         "Name of the Protection domain.",
			MarkdownDescription: "Name of the Protection domain.",
			Optional:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the Storage pool",
			MarkdownDescription: "Name of the Storage pool",
			Required:            true,
		},
		"media_type": schema.StringAttribute{
			Description:         "Media Type",
			MarkdownDescription: "Media Type",
			Required:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"HDD",
				"SSD",
				"Transitional",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("HDD"),
			},
		},
		"use_rmcache": schema.BoolAttribute{
			Description:         "Enable/Disable RMcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RMcache on a specific storage pool",
			Optional:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description:         "Enable/Disable RFcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RFcache on a specific storage pool",
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
