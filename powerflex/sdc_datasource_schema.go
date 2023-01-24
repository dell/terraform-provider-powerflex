package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var sdcDatasourceSchemaDescriptions = struct {
	SdcDatasourceSchema string

	InputID    string
	InputSdcID string
	// InputSystemid string
	InputName string

	Sdcs string // outpur slice

	LastUpdated        string
	ID                 string
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
}{
	SdcDatasourceSchema: "",

	InputID:    "Input ID required only for testing.",
	InputSdcID: "Input SDC id to search for.",
	// InputSystemid: "",
	InputName: "SDC input sdc name to search for.",

	Sdcs: "result SDCs.", // outpur slice

	LastUpdated:        "SDC result last updated timestamp.",
	ID:                 "SDC ID.",
	SystemID:           "SDC System ID.",
	Name:               "SDC name.",
	SdcIP:              "SDC IP.",
	SdcApproved:        "SDC is approved.",
	OnVMWare:           "SDC is onvmware.",
	SdcGUID:            "SDC GUID.",
	MdmConnectionState: "SDC MDM connection status.",
	Links:              "SDC Links.",
	LinksRel:           "SDC Links-Rel.",
	LinksHref:          "SDC Links-HREF.",
}

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description: sdcDatasourceSchemaDescriptions.SdcDatasourceSchema,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: sdcDatasourceSchemaDescriptions.InputSdcID,
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("name")),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"name": schema.StringAttribute{
			Description: sdcDatasourceSchemaDescriptions.InputName,
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("id")),
			},
		},
		"sdcs": schema.ListNestedAttribute{
			Description: sdcDatasourceSchemaDescriptions.Sdcs,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.ID,
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.Name,
						Computed:    true,
					},
					"sdc_guid": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcGUID,
						Computed:    true,
					},
					"on_vmware": schema.BoolAttribute{
						Description: sdcDatasourceSchemaDescriptions.OnVMWare,
						Computed:    true,
					},
					"sdc_approved": schema.BoolAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcApproved,
						Computed:    true,
					},
					"system_id": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SystemID,
						Computed:    true,
					},
					"sdc_ip": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.SdcIP,
						Computed:    true,
					},
					"mdm_connection_state": schema.StringAttribute{
						Description: sdcDatasourceSchemaDescriptions.MdmConnectionState,
						Computed:    true,
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
