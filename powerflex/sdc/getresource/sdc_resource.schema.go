package getresource

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SDCReourceSchema - varible holds schema for SDC resource
var SDCReourceSchema schema.Schema = schema.Schema{
	Description: "",
	Attributes: map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Computed: true,
		},
		"sdcid": schema.StringAttribute{
			Description: "",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "",
			Required:    true,
		},
		"sdcguid": schema.StringAttribute{
			Description: "",
			Computed:    true,
		},
		"onvmware": schema.BoolAttribute{
			Description: "",
			Computed:    true,
		},
		"sdcapproved": schema.BoolAttribute{
			Description: ".",
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: "",
			Required:    true,
		},
		"sdcip": schema.StringAttribute{
			Description: "",
			Computed:    true,
		},
		"mdmconnectionstate": schema.StringAttribute{
			Description: "",
			Computed:    true,
		},
		"links": schema.ListNestedAttribute{
			Description: "",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	},
}
