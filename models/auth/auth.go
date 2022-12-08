package models

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// AuthSchemaModel is a schema for authentication inputs for goscaleio
var AuthSchemaModel map[string]*schema.Schema = map[string]*schema.Schema{
	"host": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_HOST", nil),
	},
	"username": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Username",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_USERNAME", nil),
	},
	"password": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Password",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_PASSWORD", nil),
	},
	"insecure": {
		Type:        schema.TypeString,
		Description: "Add Insecure Value[true/false]",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_INSECURE", "true"), // anshuman check default to set
	},
	"usecerts": {
		Type:        schema.TypeString,
		Description: "Add Use Certificates Value[true/false]",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_USECERTS", "true"), // anshuman check default to set
	},
	"pfxm_version": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Verion",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_VERSION", ""), // anshuman check default to set
	},
}
