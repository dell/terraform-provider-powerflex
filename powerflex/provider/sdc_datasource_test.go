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
func TestAccDatasourceAcceptanceSdc(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + sdcDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceSdc(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Get All
			{
				Config: ProviderConfigForTesting + sdcDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Filter Single
			{
				Config: ProviderConfigForTesting + sdcDataSourceFilterSingle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sdc.filtered-single", "sdcs.0.name", "Terraform_sdc1"),
				),
			},
			// Filter multiple
			{
				Config: ProviderConfigForTesting + sdcDataSourceFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sdc.filtered-multiple", "sdcs.0.name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.filtered-multiple", "sdcs.0.system_id", "1250de83018c2d0f"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.filtered-multiple", "sdcs.1.name", "terraform_sdc_do_not_delete"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.filtered-multiple", "sdcs.1.system_id", "1250de83018c2d0f"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + sdcDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex specific system*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + sdcDataSourceFilterSingle,
				ExpectError: regexp.MustCompile(`.*Error in filtering sdcs*.`),
			},
		},
	})
}

var sdcDataSourceAll = `
	data "powerflex_sdc" "all" {

	}
`

var sdcDataSourceFilterSingle = `
	data "powerflex_sdc" "filtered-single" {
		filter {
			name = ["Terraform_sdc1"]
		}
	}
`

var sdcDataSourceFilterMultiple = `
	data "powerflex_sdc" "filtered-multiple" {
		filter {
			system_id = ["1250de83018c2d0f"]
			name = ["Terraform_sdc1", "terraform_sdc_do_not_delete"]
		}
	}
`
