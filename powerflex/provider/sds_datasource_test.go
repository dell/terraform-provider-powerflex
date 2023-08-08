/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

var protectionDomainID1 = ProtectionDomainID

func TestSdsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.0.name", "SDS_1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.0.id", "0db7306f00000003"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "protection_domain_id", protectionDomainID1),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.name", "SDS_1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.id", "0db7306f00000003"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "protection_domain_id", protectionDomainID1),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.0.name", "SDS_1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.0.id", "0db7306f00000003"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.0.name", "SDS_1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.0.id", "0db7306f00000003"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example5", "protection_domain_id", protectionDomainID1),
				),
			},
			{
				Config:      ProviderConfigForTesting + SdsDataSourceConfig6,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain.*`),
			},
			{
				Config:      ProviderConfigForTesting + SdsDataSourceConfig7,
				ExpectError: regexp.MustCompile(`.*Unable to read SDS.*`),
			},
		},
	})
}

var SdsDataSourceConfig1 = `
data "powerflex_sds" "example1" {
	protection_domain_id = "` + protectionDomainID1 + `"
	sds_names = ["SDS_1"]
}
`

var SdsDataSourceConfig2 = `
data "powerflex_sds" "example2" {
	protection_domain_id = "` + protectionDomainID1 + `"
	sds_ids = ["0db7306f00000003"]
}
`

var SdsDataSourceConfig3 = `
data "powerflex_sds" "example3" {
	protection_domain_name = "domain1"
	sds_names = ["SDS_1"]
}
`

var SdsDataSourceConfig4 = `
data "powerflex_sds" "example4" {
	protection_domain_name = "domain1"
	sds_ids = ["0db7306f00000003"]
}
`

var SdsDataSourceConfig5 = `
data "powerflex_sds" "example5" {
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var SdsDataSourceConfig6 = `
data "powerflex_sds" "example6" {
	protection_domain_id = "invalid_domain"
}
`

var SdsDataSourceConfig7 = `
data "powerflex_sds" "example7" {
	protection_domain_id = "` + protectionDomainID1 + `"
	sds_ids = ["invalid_sds_id"]
}
`
