package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description: "Manages storage pool resource",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the Storage pool",
			MarkdownDescription: "ID of the Storage pool",
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
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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
		},
		"use_rmcache": schema.BoolAttribute{
			Description:         "Enable/Disable RMcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RMcache on a specific storage pool",
			Optional:            true,
			Computed:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description:         "Enable/Disable RFcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RFcache on a specific storage pool",
			Optional:            true,
			Computed:            true,
		},
	},
}
