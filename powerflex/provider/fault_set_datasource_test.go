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
func TestAccDatasourceAcceptanceFaultSet(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceFaultSet(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Single Filter
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceSingle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_fault_set.filter-single", "fault_set_details.0.protection_domain_id", "898f009900000000"),
				),
			},
			// Multi Filter
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_fault_set.filter-multiple", "fault_set_details.0.protection_domain_id", "898f009900000000"),
					resource.TestCheckResourceAttr("data.powerflex_fault_set.filter-multiple", "fault_set_details.0.name", "fs1"),
					resource.TestCheckResourceAttr("data.powerflex_fault_set.filter-multiple", "fault_set_details.1.protection_domain_id", "898f27a900000001"),
					resource.TestCheckResourceAttr("data.powerflex_fault_set.filter-multiple", "fault_set_details.1.name", "terraform_fault_set"),
				),
			},
			// Read System error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// Read Fault Set error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllFaultSets, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Error in getting Fault Set details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetDataSourceMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering fault sets*.`),
			},
		},
	})
}

var FaultSetDataSourceAll = `
data "powerflex_fault_set" "all" {
	
}
`

var FaultSetDataSourceSingle = `
data "powerflex_fault_set" "filter-single" {
	filter {
		protection_domain_id = ["898f009900000000"]
	}
}
`

var FaultSetDataSourceMultiple = `
data "powerflex_fault_set" "filter-multiple" {
	filter {
		protection_domain_id = ["898f009900000000", "898f27a900000001"]
		name = ["fs1", "terraform_fault_set"]	
	}
}
`
