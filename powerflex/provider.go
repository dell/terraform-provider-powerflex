package powerflex

import (
	authmodel "terraform-provider-powerflex/models/auth"
	"terraform-provider-powerflex/powerflex/auth"
	"terraform-provider-powerflex/powerflex/sdc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:       authmodel.AuthSchemaModel,
		ResourcesMap: map[string]*schema.Resource{
			// "powerflex_sdcs": datasources.ResourceSdcs(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"powerflex_sdcs":            sdc.DataSourceSdcs(),
			"powerflex_sdc_name_change": sdc.ResourceSdcs(),
		},
		ConfigureContextFunc: auth.AuthHandler,
	}
}
