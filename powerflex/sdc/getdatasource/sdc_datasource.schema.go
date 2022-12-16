package getdatasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description: ".",
	Attributes: map[string]schema.Attribute{
		"sdcid": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"sdcs": schema.ListNestedAttribute{
			Description: "List of sdcs.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "",
						Computed:    true,
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
						Computed:    true,
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
									Description: "Numeric identifier of the coffee ingredient.",
									Computed:    true,
								},
								"href": schema.StringAttribute{
									Description: "Numeric identifier of the coffee ingredient.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	},
}
