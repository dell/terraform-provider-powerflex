package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var sdcDatasourceSchemaDescriptions = struct {
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
	LastUpdated:        "The Last updated timestamp of the fetched SDC.",
	ID:                 "The ID of the fetched SDC.",
	SystemID:           "The System ID of the fetched SDC.",
	Name:               "The name of the fetched SDC.",
	SdcIP:              "The IP of the fetched SDC.",
	SdcApproved:        "If the fetched SDC is approved.",
	OnVMWare:           "If the fetched SDC is on vmware.",
	SdcGUID:            "The GUID of the fetched SDC.",
	MdmConnectionState: "The MDM connection status of the fetched SDC.",
	Links:              "The Links of the fetched SDC.",
	LinksRel:           "The Links-Rel of the fetched SDC.",
	LinksHref:          "The Links-HREF of the fetched SDC.",
}

// SDCDataSourceScheme is variable for schematic for SDC Data Source
var SDCDataSourceScheme schema.Schema = schema.Schema{
	Description: "This data-source can be used to fetch information related to Storage Data Clients from a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "ID of the SDC to fetch." +
				" Conflicts with 'name'",
			MarkdownDescription: "ID of the SDC to fetch." +
				" Conflicts with `name`",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("name")),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"name": schema.StringAttribute{
			Description: "Name of the SDC to fetch." +
				" Conflicts with 'id'",
			MarkdownDescription: "Name of the SDC to fetch." +
				" Conflicts with `id`",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"sdcs": schema.ListNestedAttribute{
			Description: "List of fetched SDCs.",
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
