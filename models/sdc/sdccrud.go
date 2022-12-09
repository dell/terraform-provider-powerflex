package models

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDCCRUDModel is schema for SDC CRUD operation
var SDCCRUDModel map[string]*schema.Schema = map[string]*schema.Schema{
	"sdcid": {
		Type:        schema.TypeString,
		Description: "Enter ID of Powerflex SDC of which name will be changed.",
		Optional:    true,
	},
	"name": {
		Type:        schema.TypeString,
		Description: "Enter Name of SDC to change.",
		Optional:    true,
		Computed:    true,
	},
	"systemid": {
		Type:        schema.TypeString,
		Description: "Enter Name of SDC to change.",
		Optional:    true,
		Computed:    true,
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
