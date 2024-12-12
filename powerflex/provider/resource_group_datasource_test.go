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
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ResourceGroupDataSourceConfig = `
data "powerflex_resource_group" "example" {
}
`

var ResourceGroupDataSourceConfig2 = `
data "powerflex_resource_group" "example" {
	filter{
		id  = ["8aaa804a8b4d6b5a018b5c71a64f7052"]
	}
}
`

var ResourceGroupDataSourceConfig3 = `
data "powerflex_resource_group" "example" {
	filter{
		id  = ["8aaa804a8b4d6b5a018b5c71a64f7052"]
		deployment_name = ["Block-Storage-Hardware"]
		number_of_deployments = [0]
        compliant = true
	}
}
`

// AT
func TestAccDatasourceAcceptanceResourceGroup(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceResourceGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "resource_group_details.0.id", "8aaa804a8b4d6b5a018b5c71a64f7052"),
				),
			},
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "resource_group_details.0.id", "8aaa804a8b4d6b5a018b5c71a64f7052"),
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "resource_group_details.0.deployment_name", "Block-Storage-Hardware"),
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "resource_group_details.0.number_of_deployments", "0"),
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "resource_group_details.0.compliant", "true"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetAllServiceDetails).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ResourceGroupDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read service details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ResourceGroupDataSourceConfig3,
				ExpectError: regexp.MustCompile(`.*Error in filtering resource groups*.`),
			},
		},
	})
}
