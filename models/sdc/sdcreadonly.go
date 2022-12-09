package models

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SDCReadOnlyModel schema for SDC Data resource
var SDCReadOnlyModel map[string]*schema.Schema = map[string]*schema.Schema{
	"sdcid": {
		Type:        schema.TypeString,
		Description: "Enter ID of Powerflex SDC. [Default/empty will all sdc present in given system]",
		Required:    true,
		Sensitive:   true,
	},
	"systemid": {
		Type:        schema.TypeString,
		Description: "Enter System ID of Powerflex System. [Default/empty will be any first system in list]",
		Required:    true,
		Sensitive:   true,
	},
	"sdcs": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcguid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcapproved": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"onvmware": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"systemid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mdmconnectionstate": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"links": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"rel": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"href": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	},
}
