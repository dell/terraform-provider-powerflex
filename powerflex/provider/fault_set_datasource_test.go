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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var FaultSetDataSourceConfig1 = `
data "powerflex_fault_set" "example" {						
}
`

// To-Do: Remove hard-coded values once fault set resource gets merged
var FaultSetDataSourceConfig2 = `
data "powerflex_fault_set" "example" {						
	fault_set_ids = ["` + FaultSetID + `"]
}
`

var FaultSetDataSourceConfig3 = `
data "powerflex_fault_set" "example" {						
	fault_set_names = ["terraform_fault_set"]
}
`

var FaultSetDataSourceConfig4 = `
data "powerflex_fault_set" "example" {						
	fault_set_ids = ["invalid"]
}
`

var FaultSetDataSourceConfig5 = `
data "powerflex_fault_set" "example" {						
	fault_set_names = ["invalid"]
}
`

func TestAccDatasourceFaultSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_fault_set.example", "fault_set_details.0.id", FaultSetID),
				),
			},
			{
				Config: ProviderConfigForTesting + FaultSetDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_fault_set.example", "fault_set_details.0.name", "terraform_fault_set"),
				),
			},
			{
				Config:      ProviderConfigForTesting + FaultSetDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting fault set details using id.*`),
			},
			{
				Config:      ProviderConfigForTesting + FaultSetDataSourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in getting fault set details using name.*`),
			},
		},
	})
}
