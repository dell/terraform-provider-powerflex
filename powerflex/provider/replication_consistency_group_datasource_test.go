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

var ReplicationReplicationConsistanceGroupReadAll = `
data "powerflex_replication_consistancy_group" "rcg" {}
`
var ReplicationConsistanceGroupReadFilterID = `
data "powerflex_replication_consistancy_group" "rcg" {
    filter {
        id = ["rcg-mock-data-id"]
    }
}
`

var ReplicationConsistanceGroupReadFilterSnapCreationInProgress = `
data "powerflex_replication_consistancy_group" "rcg" {
    filter {
        snap_creation_in_progress = false
    }
}
`

var ReplicationConsistanceGroupReadFilterRpoInSeconds = `
data "powerflex_replication_consistancy_group" "rcg" {
    filter {
        rpo_in_seconds = [500]
    }
}
`

var ReplicationConsistanceGroupReadFilterMultiple = `
data "powerflex_replication_consistancy_group" "rcg" {
    filter {
        name = ["tftest-rcg", "rcg-name-2"]
        inactive_reason = [11]
    }
}
`

// Accptance Tests
func TestAccDatasourceAcceptanceReplicationConsistanceGroup(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ReplicationReplicationConsistanceGroupReadAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccDatasourceReplicationConsistanceGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with Accpetance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ReplicationReplicationConsistanceGroupReadAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + ReplicationConsistanceGroupReadFilterID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.0.id", "rcg-mock-data-id"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationConsistanceGroupReadFilterSnapCreationInProgress,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.0.snap_creation_in_progress", "false"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationConsistanceGroupReadFilterRpoInSeconds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.0.rpo_in_seconds", "500"),
				),
			},
			{
				Config: ProviderConfigForTesting + ReplicationConsistanceGroupReadFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.0.inactive_reason", "11"),
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.0.name", "tftest-rcg"),
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.1.inactive_reason", "11"),
					resource.TestCheckResourceAttr("data.powerflex_replication_consistancy_group.rcg", "replication_consistency_group_details.1.name", "rcg-name-2"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetReplicationConsistancyGroups).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationReplicationConsistanceGroupReadAll,
				ExpectError: regexp.MustCompile(`.*Error in getting Replication Consistency Groups details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationConsistanceGroupReadFilterMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering Replication Consistency Groups*.`),
			},
		},
	})
}
