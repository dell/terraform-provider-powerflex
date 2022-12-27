package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description: "Fetches the list of Storagepool",
	Attributes: map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Description: "Last Updated",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "Gets the ID of Storagepool",
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: "Gets the System ID",
			Computed:    true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "Gets the Protection Domain ID for Storagepool",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "Returns the Name of Storagepool",
			Required:    true,
		},
		"media_type": schema.StringAttribute{
			Description: "Gets the Media Type",
			Required:    true,
		},
		"use_rmcache": schema.BoolAttribute{
			Description: "Gets the Read RAM Cache",
			Computed:    true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description: "Gets the RFCache",
			Computed:    true,
		},
		"links": schema.ListNestedAttribute{
			Description: "Specifies the links asscociated with Storagepool",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description: "Specifies the relationship with the Storagepool",
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: "Specifies the exact path to fetch the details",
						Computed:    true,
					},
				},
			},
		},
	},
}
