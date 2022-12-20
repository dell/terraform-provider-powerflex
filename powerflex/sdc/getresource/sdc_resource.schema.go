package getresource

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var sdcResourceSchemaDescriptions = struct {
	SdcResourceSchema  string
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
}{
	SdcResourceSchema:  "",
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
}

// SDCReourceSchema - varible holds schema for SDC resource
var SDCReourceSchema schema.Schema = schema.Schema{
	Description: sdcResourceSchemaDescriptions.SdcResourceSchema,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "",
		},
		"last_updated": schema.StringAttribute{
			Computed:    true,
			Description: sdcResourceSchemaDescriptions.LastUpdated,
		},
		"sdcid": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SdcID,
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.Name,
			Required:    true,
		},
		"sdcguid": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SdcGUID,
			Computed:    true,
		},
		"onvmware": schema.BoolAttribute{
			Description: sdcResourceSchemaDescriptions.SdcID,
			Computed:    true,
		},
		"sdcapproved": schema.BoolAttribute{
			Description: sdcResourceSchemaDescriptions.SdcApproved,
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SystemID,
			Required:    true,
		},
		"sdcip": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SdcIP,
			Computed:    true,
		},
		"mdmconnectionstate": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.MdmConnectionState,
			Computed:    true,
		},
		"links": schema.ListNestedAttribute{
			Description: sdcResourceSchemaDescriptions.Links,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description: sdcResourceSchemaDescriptions.LinksRel,
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: sdcResourceSchemaDescriptions.LinksHref,
						Computed:    true,
					},
				},
			},
		},
	},
}
