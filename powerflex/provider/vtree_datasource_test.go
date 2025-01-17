/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
func TestAccDatasourceAcceptanceVTree(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + VtreeGetAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceVTree(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + VtreeGetAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Single Filter
			{
				Config: ProviderConfigForTesting + VtreeGetSingleFiltered,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_vtree.filter", "vtree_details.0.name", "block-volume-physical-deploy"),
				),
			},
			// Multi Filter
			{
				Config: ProviderConfigForTesting + VtreeGetMultipleFiltered,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_vtree.filter", "vtree_details.0.storage_pool_id", "68691eb600000000"),
					resource.TestCheckResourceAttr("data.powerflex_vtree.filter", "vtree_details.0.data_layout", "MediumGranularity"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllVTrees, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + VtreeGetAll,
				ExpectError: regexp.MustCompile(`.*Error in getting vTree details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + VtreeGetMultipleFiltered,
				ExpectError: regexp.MustCompile(`.*Error in filtering vtrees*.`),
			},
		},
	})
}

var VtreeGetAll = `
data "powerflex_vtree" "all" {
	
}
`

var VtreeGetSingleFiltered = `
data "powerflex_vtree" "filter" {
	filter {
		name = ["block-volume-physical-deploy"]
	}
}
`

var VtreeGetMultipleFiltered = `
data "powerflex_vtree" "filter" {
	filter {
		storage_pool_id = ["68691eb600000000"]
		data_layout = ["MediumGranularity"]
	}
}
`
