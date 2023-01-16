package powerflex

import (
	"testing"
	"regexp"
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
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example1", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example2", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.0.name", "pool2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "storage_pools.1.name", "pool1"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example3", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.0.id", "7630a24600000000"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "storage_pools.1.id", "7630a24800000002"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example4", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example5", "storage_pools.#", "8"),
					resource.TestCheckResourceAttr("data.powerflex_storagepool.example5", "protection_domain_name", "domain1"),
				),
			 },				
		},
	})
}

func TestStoragePoolDataSourceNegativeScenario(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyPDName,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyPDID,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceIncorrectPDID,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceIncorrectPDName,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptySPID,
				ExpectError: regexp.MustCompile(`.*Unable to read storage pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptySPName,
				ExpectError: regexp.MustCompile(`.*Unable to read storage pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceIncorrectSPName,
				ExpectError: regexp.MustCompile(`.*Unable to read storage pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceIncorrectSPID,
				ExpectError: regexp.MustCompile(`.*Unable to read storage pool*.`),
			},

			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyDataBlock,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyDataBlock2,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyDataBlock3,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceEmptyDataBlock4,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceBothIncorrect,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceBothIncorrect2,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceBothIncorrect3,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceBothIncorrect4,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceWithEmptySPNamePartII,
				ExpectError: regexp.MustCompile(`.*No storage pools found for the specified protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceWithEmptySPIDPartII,
				ExpectError: regexp.MustCompile(`.*No storage pools found for the specified protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPDataSourceWOPDIDAndName,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Combination.*"),
			},
		},
	})
}

var StoragePoolDataSourceConfig1 = `
data "powerflex_storagepool" "example1" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_names = ["pool2", "pool1"]
}
`

var StoragePoolDataSourceConfig2 = `
data "powerflex_storagepool" "example2" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`

var StoragePoolDataSourceConfig3 = `
data "powerflex_storagepool" "example3" {
	protection_domain_name = "domain1"
	storage_pool_names = ["pool2", "pool1"]
}
`
var StoragePoolDataSourceConfig4 = `
data "powerflex_storagepool" "example4" {
	protection_domain_name = "domain1"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`
var StoragePoolDataSourceConfig5 = `
data "powerflex_storagepool" "example5" {
	protection_domain_name = "domain1"
}
`
var SPDataSourceEmptyPDName = `
data "powerflex_storagepool" "example5" {
	protection_domain_name = ""
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`

var SPDataSourceEmptyPDID = `
data "powerflex_storagepool" "example6" {
	protection_domain_id = ""
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`

var SPDataSourceIncorrectPDID = `
data "powerflex_storagepool" "example7" {
	protection_domain_id = "4eeb30460000"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`

var SPDataSourceIncorrectPDName = `
data "powerflex_storagepool" "example8" {
	protection_domain_name = "abcde"
	storage_pool_ids = ["7630a24600000000", "7630a24800000002"]
}
`
var SPDataSourceEmptySPID = `
data "powerflex_storagepool" "example9" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_ids = [""]
}
`

var SPDataSourceEmptySPName = `
data "powerflex_storagepool" "example10" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_names = [""]
}
`

var SPDataSourceIncorrectSPName = `
data "powerflex_storagepool" "example11" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_names = ["abcde"]
}
`

var SPDataSourceIncorrectSPID = `
data "powerflex_storagepool" "example12" {
	protection_domain_id = "4eeb304600000000"
	storage_pool_ids = ["7630a246000"]
}
`

var SPDataSourceEmptyDataBlock = `
data "powerflex_storagepool" "example13" {
	protection_domain_id = ""
	storage_pool_ids = [""]
}
`
var SPDataSourceEmptyDataBlock2 = `
data "powerflex_storagepool" "example14" {
	protection_domain_name = ""
	storage_pool_names = [""]
}
`

var SPDataSourceEmptyDataBlock3 = `
data "powerflex_storagepool" "example15" {
	protection_domain_id = ""
	storage_pool_names = [""]
}
`

var SPDataSourceEmptyDataBlock4 = `
data "powerflex_storagepool" "example16" {
	protection_domain_name = ""
	storage_pool_ids = [""]
}
`
var SPDataSourceBothIncorrect = `
data "powerflex_storagepool" "example17" {
	protection_domain_id = "4eeb3046000"
	storage_pool_ids = ["7630a24600"]
}
`

var SPDataSourceBothIncorrect2 = `
data "powerflex_storagepool" "example18" {
	protection_domain_name = "abcde"
	storage_pool_names = ["fghij"]
}
`

var SPDataSourceBothIncorrect3 = `
data "powerflex_storagepool" "example19" {
	protection_domain_id = "4eeb30460000"
	storage_pool_names = ["fghij"]
}
`

var SPDataSourceBothIncorrect4 = `
data "powerflex_storagepool" "example20" {
	protection_domain_name = "abcde"
	storage_pool_ids = ["7630a2460000"]
}
`

var SPDataSourceWithEmptySPNamePartII = `
data "powerflex_storagepool" "example20" {
	protection_domain_name = "domain1"
	storage_pool_names = []
}
`

var SPDataSourceWithEmptySPIDPartII = `
data "powerflex_storagepool" "example20" {
	protection_domain_name = "domain1"
	storage_pool_ids = []
}
`

var SPDataSourceWOPDIDAndName = `
data "powerflex_storagepool" "example21" {
	storage_pool_names = ["pool1"]
}
`