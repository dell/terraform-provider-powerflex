package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var sdcResourceSchemaDescriptions = struct {
	SdcResourceSchema  string
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
	SdcResourceSchema:  "",
	LastUpdated:        "last updated timestamp.",
	ID:                 "SDC ID.",
	SystemID:           "System ID.",
	Name:               "SDC Name.",
	SdcIP:              "SDC IP.",
	SdcApproved:        "SDC Approved.",
	OnVMWare:           "On VMware.",
	SdcGUID:            "SDC GUID.",
	MdmConnectionState: "MDM Connection state.",
	Links:              "Links.",
	LinksRel:           "Links Rel.",
	LinksHref:          "Links HREF.",
}

// SDCReourceSchema - varible holds schema for SDC resource
var SDCReourceSchema schema.Schema = schema.Schema{
	Description: sdcResourceSchemaDescriptions.SdcResourceSchema,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Required:    true,
			Description: sdcResourceSchemaDescriptions.ID,
		},
		"last_updated": schema.StringAttribute{
			Computed:    true,
			Description: sdcResourceSchemaDescriptions.LastUpdated,
		},
		"name": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.Name,
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"sdc_guid": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SdcGUID,
			Computed:    true,
		},
		"on_vmware": schema.BoolAttribute{
			Description: sdcResourceSchemaDescriptions.OnVMWare,
			Computed:    true,
		},
		"sdc_approved": schema.BoolAttribute{
			Description: sdcResourceSchemaDescriptions.SdcApproved,
			Computed:    true,
		},
		"system_id": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SystemID,
			Computed:    true,
		},
		"sdc_ip": schema.StringAttribute{
			Description: sdcResourceSchemaDescriptions.SdcIP,
			Computed:    true,
		},
		"mdm_connection_state": schema.StringAttribute{
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
