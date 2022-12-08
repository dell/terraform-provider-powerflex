package powerflex

import (
	datasources "terraform-provider-powerflex/powerflex/data-sources"
	schemastructures "terraform-provider-powerflex/powerflex/schema-structures"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:       schemastructures.AUTH_SCHEMA,
		ResourcesMap: map[string]*schema.Resource{
			// "powerflex_sdcs": datasources.ResourceSdcs(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"powerflex_sdcs":            datasources.DataSourceSdcs(),
			"powerflex_sdc_name_change": datasources.ResourceSdcs(),
		},
		ConfigureContextFunc: schemastructures.AuthConfigure,
	}
}
