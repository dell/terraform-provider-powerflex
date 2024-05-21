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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var FRDataSourceConfig1 = `
data "powerflex_firmware_repository" "test" {
	firmware_repository_ids = ["8aaa3fda8f5c2609018f854266e12865", "8aaa3fda8f5c2609018f857b6c0d2ede"]
	}
`

var FRDataSourceConfig2 = `
data "powerflex_firmware_repository" "test" {
	firmware_repository_names = ["PowerFlex 4.5.0.0 (14)", "PowerFlex 4.5.0.0 (16)"]
	}
`

var FRDataSourceConfig3 = `
data "powerflex_firmware_repository" "test" {
}
`

var FRDataSourceConfig4 = `
data "powerflex_firmware_repository" "test" {
	firmware_repository_ids = ["8aaa3fda8f5c2609018f854266e12865", "Invalid"]
}
`

var FRDataSourceConfig5 = `
data "powerflex_firmware_repository" "test" {
	firmware_repository_names = ["Invalid", "PowerFlex 4.5.0.0 (16)"]
}
`

func TestAccFRDataSource(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FRDataSourceConfig1,
			},
			{
				Config: ProviderConfigForTesting + FRDataSourceConfig2,
			},
			{
				Config: ProviderConfigForTesting + FRDataSourceConfig3,
			},
		},
	})
}

func TestAccFRDataSourceNegative(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + FRDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting firmware repository details using id Invalid.*`),
			},
			{
				Config:      ProviderConfigForTesting + FRDataSourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in getting firmware repository details using name Invalid.*`),
			},
		},
	})
}
