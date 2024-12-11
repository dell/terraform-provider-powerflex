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

var PeerMdmDataSourceAll = `
data "powerflex_peer_system" "test" {
}
`
var PeerMdmDataSourceByID = `
data "powerflex_peer_system" "test" {
	# this datasource supports filters like peer mdm ids, names, coupling RC, etc.
	filter {
		id = ["` + PeerMdmID + `"]
	}
  }
`

var PeerMdmDataSourceMultipleFilters = `
data "powerflex_peer_system" "test" {
	# this datasource supports filters like peer mdm ids, names, coupling RC, etc.
	filter {
		id = ["` + PeerMdmID + `"]
	    name = ["` + PeerMdmName + `"]
		coupling_rc = ["` + PeerMdmCouplingRC + `"]
		port = ["` + PeerMdmPort + `"]
	}
  }
`

// Accptance Tests
func TestAccDatasourceAcceptancePeerMdm(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + PeerMdmDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccDatasourcePeerMdm(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + PeerMdmDataSourceAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_peer_system.test", "peer_system_details.#"),
				),
			},
			{
				Config: ProviderConfigForTesting + PeerMdmDataSourceByID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_peer_system.test", "peer_system_details.0.id", PeerMdmID),
				),
			},
			{
				Config: ProviderConfigForTesting + PeerMdmDataSourceMultipleFilters,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_peer_system.test", "peer_system_details.0.id", PeerMdmID),
					resource.TestCheckResourceAttr("data.powerflex_peer_system.test", "peer_system_details.0.name", PeerMdmName),
					resource.TestCheckResourceAttr("data.powerflex_peer_system.test", "peer_system_details.0.coupling_rc", PeerMdmCouplingRC),
					resource.TestCheckResourceAttr("data.powerflex_peer_system.test", "peer_system_details.0.port", PeerMdmPort),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetPeerMdms).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerMdmDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Error in getting Peer MDM details*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerMdmDataSourceMultipleFilters,
				ExpectError: regexp.MustCompile(`.*Error in filtering Peer MDMs*.`),
			},
		},
	})
}
