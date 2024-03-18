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

var VolPreq1 = `
resource "powerflex_volume" "pre-req1-vtree"{
	name = "terraform-vol-vtree-datasource"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 8
}`

var VolPreq2 = `
resource "powerflex_volume" "pre-req2-vtree"{
	name = "terraform-vol-vtree-datasource1"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 8
}`

var VTreeDataSourceConfig4 = `
data "powerflex_vtree" "example" {						
	vtree_ids = ["27ca983400000000"]
}
`

var VTreeDataSourceConfig2 = VolPreq1 + VolPreq2 + `
data "powerflex_vtree" "example" {						
	volume_ids = [resource.powerflex_volume.pre-req1-vtree.id, resource.powerflex_volume.pre-req2-vtree.id]
}
`

var VTreeDataSourceConfig3 = VolPreq1 + VolPreq2 + `
data "powerflex_vtree" "example" {						
	volume_names = [resource.powerflex_volume.pre-req1-vtree.name, resource.powerflex_volume.pre-req2-vtree.name]
}
`

var VTreeDataSourceConfig1 = `
data "powerflex_vtree" "example" {						
}
`

var VTreeDataSourceConfig5 = `
data "powerflex_vtree" "example" {						
	vtree_ids = ["invalid"]
}
`

var VTreeDataSourceConfig6 = `
data "powerflex_vtree" "example" {						
	volume_ids = ["invalid"]
}
`

var VTreeDataSourceConfig7 = `
data "powerflex_vtree" "example" {						
	volume_names = ["invalid"]
}
`

func TestAccVTreeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + VolPreq1,
			},
			{
				Config: ProviderConfigForTesting + VolPreq2,
			},
			{
				Config: ProviderConfigForTesting + VTreeDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + VTreeDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_vtree.example", "vtree_details.#", "2"),
				),
			},
			{
				Config: ProviderConfigForTesting + VTreeDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_vtree.example", "vtree_details.#", "2"),
				),
			},
			{
				Config: ProviderConfigForTesting + VTreeDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_vtree.example", "vtree_details.#", "1"),
				),
			},
			{
				Config:      ProviderConfigForTesting + VTreeDataSourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in getting vTree details using vTree IDs.*`),
			},
			{
				Config:      ProviderConfigForTesting + VTreeDataSourceConfig6,
				ExpectError: regexp.MustCompile(`.*Error in getting vTree details with volume ID.*`),
			},
			{
				Config:      ProviderConfigForTesting + VTreeDataSourceConfig7,
				ExpectError: regexp.MustCompile(`.*Error in getting volume details with volume name.*`),
			},
		},
	})
}
