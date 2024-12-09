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

var NodeDataSourceConfig1 = `
data "powerflex_node" "all" {
}
`

var NodeDataSourceConfig2 = `
data "powerflex_node" "example" {
	filter{
		ref_id  = ["scaleio-block-legacy-gateway"]
	}
}
`

var NodeDataSourceConfig3 = `
data "powerflex_node" "example2" {
	filter{
		ref_id  = ["scaleio-block-legacy-gateway"]
		state = ["READY"]
		ip_address = ["block-legacy-gateway"]
		service_tag = ["block-legacy-gateway"]
	}
}
`

// AT
func TestAccDatasourceAcceptanceNode(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceNode(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
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
					resource.TestCheckResourceAttr("data.powerflex_node.example", "node_details.0.ref_id", "scaleio-block-legacy-gateway"),
				),
			},
			{
				Config: ProviderConfigForTesting + NodeDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_node.example2", "node_details.0.ip_address", "block-legacy-gateway"),
					resource.TestCheckResourceAttr("data.powerflex_node.example2", "node_details.0.service_tag", "block-legacy-gateway"),
					resource.TestCheckResourceAttr("data.powerflex_node.example2", "node_details.0.ref_id", "scaleio-block-legacy-gateway"),
					resource.TestCheckResourceAttr("data.powerflex_node.example2", "node_details.0.state", "READY"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetAllNodes).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + NodeDataSourceConfig1,
				ExpectError: regexp.MustCompile(`.*Unable to Read node details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + NodeDataSourceConfig3,
				ExpectError: regexp.MustCompile(`.*Error in filtering node*.`),
			},
		},
	})
}
