/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

var FRDataSourceAll = `
data "powerflex_firmware_repository" "all" {}
`

var FRDataSourceSingle = `
data "powerflex_firmware_repository" "filter-single" {
	filter {
		bundle_count = [0]
	}
}
`

var FRDataSourceMultiple = `
data "powerflex_firmware_repository" "filter-multiple" {
	filter {
		minimal = false
		name = ["Intelligent Catalog 45.373.00"]
	}
}
`
var FRDataSourceRegex = `
data "powerflex_firmware_repository" "filter-regex" {
	filter {
		name = ["^Intelligent.*$"]
	}
}
`

// AT
func TestAccDatasourceAcceptanceFR(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FRDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceFR(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Get All
			{
				Config: ProviderConfigForTesting + FRDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Filter Single
			{
				Config: ProviderConfigForTesting + FRDataSourceSingle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_firmware_repository.filter-single", "firmware_repository_details.0.bundle_count", "0"),
				),
			},
			// Filter Multiple
			{
				Config: ProviderConfigForTesting + FRDataSourceMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_firmware_repository.filter-multiple", "firmware_repository_details.0.minimal", "false"),
					resource.TestCheckResourceAttr("data.powerflex_firmware_repository.filter-multiple", "firmware_repository_details.0.name", "Intelligent Catalog 45.373.00"),
				),
			},
			// Filter Regex
			{
				Config: ProviderConfigForTesting + FRDataSourceRegex,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_firmware_repository.filter-regex", "firmware_repository_details.0.name", "Intelligent Catalog 45.373.00"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirmwareRepositories).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FRDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Error in getting firmware repositories*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FRDataSourceMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering firmware repositories*.`),
			},
		},
	})
}
