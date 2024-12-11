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
	"fmt"
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AT
func TestAccDatasourceAcceptanceStoragePool(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceStoragePool(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Get All
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Filter Single
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceFilterSingle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.filter", "storage_pools.0.name", "Terraform_pool"),
				),
			},
			// Filter Multiple
			{
				Config: ProviderConfigForTesting + StoragePoolDataSourceFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.filter", "storage_pools.0.rebalance_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("data.powerflex_storage_pool.filter", "storage_pools.0.use_rm_cache", "false"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllStoragePools).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + StoragePoolDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Storage Groups*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + StoragePoolDataSourceFilterMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering storage pools*.`),
			},
		},
	})
}

var StoragePoolDataSourceAll = `
data "powerflex_storage_pool" "all" {

}
`

var StoragePoolDataSourceFilterSingle = `
data "powerflex_storage_pool" "filter" {
	filter {
		name = ["Terraform_pool"]
	}
}
`

var StoragePoolDataSourceFilterMultiple = `
data "powerflex_storage_pool" "filter" {
	filter {
		rebalance_io_priority_policy = ["favorAppIos"]
		use_rm_cache = false
	}
}
`
