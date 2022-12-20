package getdatasource

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description: ".",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "",
			Optional:    true,
		},
		"sdcid": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: "",
			Required:    true,
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
					"statistics": schema.ObjectAttribute{
						Description: "",
						Computed:    true,
						AttributeTypes: map[string]attr.Type{
							"numofmappedvolumes": types.Int64Type,
							"volumeids":          types.ListType{ElemType: types.StringType},
							"userdatareadbwc": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
							"userdatawritebwc": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
							"userdatatrimbwc": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
							"userdatasdcreadlatency": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
							"userdatasdcwritelatency": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
							"userdatasdctrimlatency": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"totalweightinkb": types.Int64Type,
									"numoccured":      types.Int64Type,
									"numseconds":      types.Int64Type,
								},
							},
						},
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
