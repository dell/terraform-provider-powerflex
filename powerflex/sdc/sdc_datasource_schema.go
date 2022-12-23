package sdcsource

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var sdcDatasourceSchemaDescriptions = struct {
	SdcDatasourceSchema string

	InputID       string
	InputSdcID    string
	InputSystemid string
	InputName     string

	Sdcs string // outpur slice

	LastUpdated        string
	SdcID              string
	SystemID           string
	Name               string
	SdcIP              string
	SdcApproved        string
	OnVMWare           string
	SdcGUID            string
	MdmConnectionState string
	Links              string
	LinksRel           string
	LinksHref          string
	Statistics         string
}{
	SdcDatasourceSchema: "",

	InputID:       "",
	InputSdcID:    "",
	InputSystemid: "",
	InputName:     "",

	Sdcs: "", // outpur slice

	LastUpdated:        "",
	SdcID:              "",
	SystemID:           "",
	Name:               "",
	SdcIP:              "",
	SdcApproved:        "",
	OnVMWare:           "",
	SdcGUID:            "",
	MdmConnectionState: "",
	Links:              "",
	LinksRel:           "",
	LinksHref:          "",
	Statistics:         "",
}

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description: sdcDatasourceSchemaDescriptions.SdcDatasourceSchema,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: sdcDatasourceSchemaDescriptions.InputID,
			Optional:    true,
		},
		"sdcid": schema.StringAttribute{
			Description: sdcDatasourceSchemaDescriptions.InputSdcID,
			Optional:    true,
			Computed:    true,
		},
		// "systemid": schema.StringAttribute{
		// 	Description: sdcDatasourceSchemaDescriptions.InputSystemid,
		// 	Required:    true,
		// },
		"name": schema.StringAttribute{
			Description: sdcDatasourceSchemaDescriptions.InputName,
			Optional:    true,
			Computed:    true,
		},
		"sdcs": schema.ListNestedAttribute{
			Description: sdcDatasourceSchemaDescriptions.Sdcs,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcID,
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.Name,
						Computed:    true,
					},
					"sdcguid": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcGUID,
						Computed:    true,
					},
					"onvmware": schema.BoolAttribute{
						Description: sdcDatasourceSchemaDescriptions.OnVMWare,
						Computed:    true,
					},
					"sdcapproved": schema.BoolAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcApproved,
						Computed:    true,
					},
					"systemid": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SystemID,
						Computed:    true,
					},
					"sdcip": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcIP,
						Computed:    true,
					},
					"mdmconnectionstate": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.MdmConnectionState,
						Computed:    true,
					},
					"statistics": schema.ObjectAttribute{
						Description: sdcDatasourceSchemaDescriptions.Statistics,
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
						Description: sdcDatasourceSchemaDescriptions.Links,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description: sdcDatasourceSchemaDescriptions.LinksRel,
									Computed:    true,
								},
								"href": schema.StringAttribute{
									Description: sdcDatasourceSchemaDescriptions.LinksHref,
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
