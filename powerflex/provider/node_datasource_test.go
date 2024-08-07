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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NodeDataSourceConfig1 = `
data "powerflex_node" "example" {						
}
`

var NodeDataSourceConfig2 = `
data "powerflex_node" "example" {
	node_ids = ["` + NodeDataPoints.NodeID + `"]					
}
`

var NodeDataSourceConfig3 = `
data "powerflex_node" "example" {
	ip_addresses = ["` + NodeDataPoints.NodeIP + `"]						
}
`

var NodeDataSourceConfig4 = `
data "powerflex_node" "example" {
	service_tags = ["` + NodeDataPoints.ServiceTag + `"]						
}
`

var NodeDataSourceConfig5 = `
data "powerflex_node" "example" {
	node_pool_ids = ["` + NodeDataPoints.NodePoolID + `"]						
}
`

var NodeDataSourceConfig6 = `
data "powerflex_node" "example" {
	node_pool_names = ["` + NodeDataPoints.NodePoolName + `", "Global"]						
}
`

var NodeDataSourceConfig7 = `
data "powerflex_node" "example" {
	node_ids = ["invalid"]					
}
`

func TestAccDatasourceNode(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.0.ref_id", NodeDataPoints.NodeID),
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.0.ip_address", NodeDataPoints.NodeIP),
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.0.service_tag", NodeDataPoints.ServiceTag),
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig5,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig6,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config:      ProviderConfigForTesting + NodeDataSourceConfig7,
				ExpectError: regexp.MustCompile(`.*Error in getting node details using id*.`),
			},
		},
	})
}
