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
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example1", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example1", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example1", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example1", "protection_domain_id", "202a046600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "protection_domain_id", "202a046600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example3", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example3", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example3", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example3", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example5", "protection_domain_name", "domain1"),
				),
			},
		},
	})
}

var StoragePoolDataSourceConfig1 = `
data "powerflex_storage_pool" "example1" {
	protection_domain_id = "202a046600000000"
	storage_pool_names = ["pool2", "pool1"]
}
`

var StoragePoolDataSourceConfig2 = `
data "powerflex_storage_pool" "example2" {
	protection_domain_id = "202a046600000000"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`

var StoragePoolDataSourceConfig3 = `
data "powerflex_storage_pool" "example3" {
	protection_domain_name = "domain1"
	storage_pool_names = ["pool2", "pool1"]
}
`
var StoragePoolDataSourceConfig4 = `
data "powerflex_storage_pool" "example4" {
	protection_domain_name = "domain1"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`
var StoragePoolDataSourceConfig5 = `
data "powerflex_storage_pool" "example5" {
	protection_domain_name = "domain1"
}
`
