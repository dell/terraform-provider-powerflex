package sdcprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SDCDataSourceScheme tfsdk.Schema = tfsdk.Schema{
	Description: ".",
	Attributes: map[string]tfsdk.Attribute{
		"sdcid": {
			Type:        types.StringType,
			Description: "",
			Required:    true,
			Sensitive:   true,
		},
		"systemid": {
			Type:        types.StringType,
			Description: "",
			Required:    true,
			Sensitive:   true,
		},
		"sdcs": {
			Description: ".",
			Computed:    true,
			Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				"id": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"name": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"sdcguid": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"onvmware": {
					Description: "",
					Type:        types.BoolType,
					Computed:    true,
				},
				"sdcapproved": {
					Description: ".",
					Type:        types.BoolType,
					Computed:    true,
				},
				"systemid": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"sdcip": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"mdmconnectionstate": {
					Description: "",
					Type:        types.StringType,
					Computed:    true,
				},
				"links": {
					Description: "",
					Computed:    true,
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"rel": {
							Description: "",
							Type:        types.StringType,
							Computed:    true,
						},
						"href": {
							Description: "",
							Type:        types.StringType,
							Computed:    true,
						},
					}),
				},
			}),
		},
	},
}
