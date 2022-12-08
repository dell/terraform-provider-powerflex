package schemastructures

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDC CRUD operation Schema
var SDCCRUDSchema map[string]*schema.Schema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "Enter ID of Powerflex SDC of which name will be changed.",
		Required:    true,
		Sensitive:   true,
	},
	"name": {
		Type:        schema.TypeString,
		Description: "Enter Name of SDC to change.",
		Required:    true,
		Sensitive:   true,
	},
	"systemid": {
		Type:        schema.TypeString,
		Description: "Enter Name of SDC to change.",
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
			},
		},
	},
}
