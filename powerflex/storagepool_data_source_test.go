package powerflex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestStoragePoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: StoragePoolDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: StoragePoolDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: StoragePoolDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: StoragePoolDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "protection_domain_name", "domain1"),
				),
			},
		},
	})
}

var StoragePoolDataSourceConfig1 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

data "powerflex_storagepool" "example1" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_name = ["pool2", "pool1"]
}
`

var StoragePoolDataSourceConfig2 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

data "powerflex_storagepool" "example2" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_id = ["7630a24600000000", "7630a24800000002"]
}
`

var StoragePoolDataSourceConfig3 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

data "powerflex_storagepool" "example3" {
	protection_domain_name = "domain1"
	storage_pool_name = ["pool2", "pool1"]
}
`
var StoragePoolDataSourceConfig4 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

data "powerflex_storagepool" "example4" {
	protection_domain_name = "domain1"
	storage_pool_id = ["7630a24600000000", "7630a24800000002"]
}
`
