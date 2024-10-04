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

var ReplicationPairReadAll = `
data "powerflex_replication_pair" "rp" {}
`

var ReplicationPairReadFilterID = `
data "powerflex_replication_pair" "rp" {
    filter {
        id = ["rp-mock-data-id"]
    }
}
`

var ReplicationPairReadFilterInitCopy = `
data "powerflex_replication_pair" "rp" {
    filter {
        user_requested_pause_transmit_init_copy = false
    }
}
`

var ReplicationPairReadFilterRemoteMb = `
data "powerflex_replication_pair" "rp" {
    filter {
        remote_capacity_in_mb = [8192]
    }
}
`

var ReplicationPairReadFilterMultiple = `
data "powerflex_replication_pair" "rp" {
    filter {
        name = ["rp-name-1", "rp-name-2"]
        remote_capacity_in_mb = [8192]
    }
}
`

// Accptance Tests
func TestAccDatasourceAcceptanceReplicationPair(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ReplicationPairReadAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccDatasourceReplicationPair(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ReplicationPairReadAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + ReplicationPairReadFilterID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.0.id", "rp-mock-data-id"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationPairReadFilterInitCopy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.0.user_requested_pause_transmit_init_copy", "false"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationPairReadFilterRemoteMb,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.0.remote_capacity_in_mb", "8192"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationPairReadFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.0.remote_capacity_in_mb", "8192"),
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.0.name", "rp-name-1"),
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.1.remote_capacity_in_mb", "8192"),
					resource.TestCheckResourceAttr("data.powerflex_replication_pair.rp", "replication_pair_details.1.name", "rp-name-2"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetReplicationPairs).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairReadAll,
				ExpectError: regexp.MustCompile(`.*Error in getting Replication Pairs details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairReadFilterMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering Replication Pairs*.`),
			},
		},
	})
}
