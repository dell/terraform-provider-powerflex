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
				Config: ProviderConfigForTesting + devicesData,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithName,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithPath,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithID,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithStoragePoolName,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithStoragePoolID,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithSdsName,
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithSdsID,
			},
		},
	})
}

var devicesData = `
data "powerflex_device" "dev" {
}
`

var deviceDataWithName = `
data "powerflex_device" "dev" {
	name = "device_1"
}
`

var deviceDataWithPath = `
data "powerflex_device" "dev" {
	current_path = "/dev/sdb"
}
`

var deviceDataWithID = `
data "powerflex_device" "dev" {
	id = "c7fc68a200000000"
}
`
var deviceDataWithStoragePoolName = `
data "powerflex_device" "dev" {
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
  }
`

var deviceDataWithStoragePoolID = `
data "powerflex_device" "dev" {
	storage_pool_id = "c98e26e500000000"
}
`

var deviceDataWithSdsName = `
data "powerflex_device" "dev" {
	sds_name = "SDS_2"
}
`

var deviceDataWithSdsID = `
data "powerflex_device" "dev" {
	sds_id = "db2c37000000000"
}
`
