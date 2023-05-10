package powerflex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDeviceDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + deviceData,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithID,
			},
		},
	})
}

var deviceData = `
data "powerflex_device" "dev" {
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
  }
`

var deviceDataWithID = `
data "powerflex_device" "dev" {
	id = "c7fc68a200000000"
}
`

// id = "c7fc68a200000000"
