/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceStoragePool(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example1", "protection_domain_id", ProtectionDomainID),
				),
			},
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.0.id", "c992bad600000005"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "storage_pools.1.id", "c98e26e500000000"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example2", "protection_domain_id", ProtectionDomainID),
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
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "storage_pools.0.id", "c992bad600000005"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "storage_pools.1.id", "c98e26e500000000"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.example4", "protection_domain_name", "domain1"),
				),
			},
		},
	})
}

var StoragePoolDataSourceConfig1 = `
data "powerflex_storage_pool" "example1" {
	protection_domain_id = "` + ProtectionDomainID + `"
	storage_pool_names = ["pool2", "pool1"]
}
`

var StoragePoolDataSourceConfig2 = `
data "powerflex_storage_pool" "example2" {
	protection_domain_id = "` + ProtectionDomainID + `"
	storage_pool_ids = ["c992bad600000005", "c98e26e500000000"]
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
	storage_pool_ids = ["c992bad600000005", "c98e26e500000000"]
}
`
