package powerflex

import (
	"regexp"
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
			{
				Config:      ProviderConfigForTesting + deviceDataWithNameInvalid,
				ExpectError: regexp.MustCompile("Error getting device with Name"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithPathInvalid,
				ExpectError: regexp.MustCompile("Error getting device with CurrentPath"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithIDInvalid,
				ExpectError: regexp.MustCompile("Error getting device with ID"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithProtectionDomainNameInvalid,
				ExpectError: regexp.MustCompile("Error in getting protection domain details with ID"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithStoragePoolNameInvalid,
				ExpectError: regexp.MustCompile("Error in getting storage pool details with name"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithStoragePoolIDInvalid,
				ExpectError: regexp.MustCompile("Error getting storage pool instance with ID"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithSdsIDInvalid,
				ExpectError: regexp.MustCompile("Could not get SDS by ID"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithSdsNameInvalid,
				ExpectError: regexp.MustCompile("Error in getting sds details with name"),
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

var deviceDataWithNameInvalid = `
data "powerflex_device" "dev" {
	name = "invalid"
}
`

var deviceDataWithPathInvalid = `
data "powerflex_device" "dev" {
	current_path = "invalid"
}
`

var deviceDataWithIDInvalid = `
data "powerflex_device" "dev" {
	id = "Invalid"
}
`

var deviceDataWithProtectionDomainNameInvalid = `
data "powerflex_device" "dev" {
	protection_domain_name = "Invalid"
	storage_pool_name = "pool1"
  }
`

var deviceDataWithStoragePoolNameInvalid = `
data "powerflex_device" "dev" {
	protection_domain_name = "domain1"
	storage_pool_name = "Invalid"
  }
`

var deviceDataWithStoragePoolIDInvalid = `
data "powerflex_device" "dev" {
	storage_pool_id = "Invalid"
}
`

var deviceDataWithSdsIDInvalid = `
data "powerflex_device" "dev" {
	sds_id = "invalid"
}
`

var deviceDataWithSdsNameInvalid = `
data "powerflex_device" "dev" {
	sds_name = "invalid"
}
`
