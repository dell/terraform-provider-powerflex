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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev1", "device_model.#", "3"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev2", "name", "device_1"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithPath,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev3", "current_path", "/dev/sdb"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev4", "id", "c7fc68a200000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithStoragePoolName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev5", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("data.powerflex_device.dev5", "storage_pool_name", "pool1"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithStoragePoolID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev6", "storage_pool_id", "c98e26e500000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithSdsName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev7", "sds_name", "SDS_2"),
				),
			},
			{
				Config: ProviderConfigForTesting + deviceDataWithSdsID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.dev8", "sds_id", "db2c37000000000"),
				),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithNameInvalid,
				ExpectError: regexp.MustCompile("Error getting device with Name"),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithPathInvalid,
				ExpectError: regexp.MustCompile("Error getting device with Current Path"),
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
			{
				Config:      ProviderConfigForTesting + deviceDataWithPDIDInvalid,
				ExpectError: regexp.MustCompile("Please provide protection_domain_id with storage_pool_name."),
			},
			{
				Config:      ProviderConfigForTesting + deviceDataWithPDNameInvalid,
				ExpectError: regexp.MustCompile("Please provide protection_domain_name with storage_pool_name."),
			},
		},
	})
}

var devicesData = `
data "powerflex_device" "dev1" {
}
`

var deviceDataWithName = `
data "powerflex_device" "dev2" {
	name = "device_1"
}
`

var deviceDataWithPath = `
data "powerflex_device" "dev3" {
	current_path = "/dev/sdb"
}
`

var deviceDataWithID = `
data "powerflex_device" "dev4" {
	id = "c7fc68a200000000"
}
`
var deviceDataWithStoragePoolName = `
data "powerflex_device" "dev5" {
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
  }
`

var deviceDataWithStoragePoolID = `
data "powerflex_device" "dev6" {
	storage_pool_id = "c98e26e500000000"
}
`

var deviceDataWithSdsName = `
data "powerflex_device" "dev7" {
	sds_name = "SDS_2"
}
`

var deviceDataWithSdsID = `
data "powerflex_device" "dev8" {
	sds_id = "db2c37000000000"
}
`

var deviceDataWithNameInvalid = `
data "powerflex_device" "dev9" {
	name = "invalid"
}
`

var deviceDataWithPathInvalid = `
data "powerflex_device" "dev10" {
	current_path = "invalid"
}
`

var deviceDataWithIDInvalid = `
data "powerflex_device" "dev11" {
	id = "Invalid"
}
`

var deviceDataWithProtectionDomainNameInvalid = `
data "powerflex_device" "dev12" {
	protection_domain_name = "Invalid"
	storage_pool_name = "pool1"
  }
`

var deviceDataWithStoragePoolNameInvalid = `
data "powerflex_device" "dev13" {
	protection_domain_name = "domain1"
	storage_pool_name = "Invalid"
  }
`

var deviceDataWithStoragePoolIDInvalid = `
data "powerflex_device" "dev14" {
	storage_pool_id = "Invalid"
}
`

var deviceDataWithSdsIDInvalid = `
data "powerflex_device" "dev15" {
	sds_id = "invalid"
}
`

var deviceDataWithSdsNameInvalid = `
data "powerflex_device" "dev16" {
	sds_name = "invalid"
}
`

var deviceDataWithPDIDInvalid = `
data "powerflex_device" "dev17" {
	protection_domain_id = "202a046600000000"
}
`

var deviceDataWithPDNameInvalid = `
data "powerflex_device" "dev17" {
	protection_domain_name = "domain1"
}
`