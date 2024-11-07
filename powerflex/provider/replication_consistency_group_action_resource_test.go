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

var ReplicationConsistencyGroupActionResourceConfigFailover = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Failover"
}
`

var ReplicationConsistencyGroupActionResourceConfigRestore = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Restore"
}
`

var ReplicationConsistencyGroupActionResourceConfigSync = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Sync"
}
`

var ReplicationConsistencyGroupActionResourceConfigReverse = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Reverse"
}
`

var ReplicationConsistencyGroupActionResourceConfigSwitchover = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Switchover"
}
`

var ReplicationConsistencyGroupActionResourceConfigSnapshot = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
    action = "Snapshot"
}
`

var ReplicationConsistencyGroupActionResourceConfigInvalidAction = `
resource "powerflex_replication_consistency_group_action" "example" {
    id = "` + RcgID + `"

    # Action to be performed on the replication consistency group.
    # Options are Failover, Restore, Sync, Reverse, Switchover and Snapshot (Default is Sync)
	action = "InvalidAction"
}
`

// Accptance Tests
func TestAccResourceAcceptanceReplicationConsistencyGroupActions(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Failover
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigFailover,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Restore
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigRestore,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceReplicationConsistencyGroupActions(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Action error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.RCGDoAction).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigFailover,
				ExpectError: regexp.MustCompile(`.*Error doing action Failover on replication consistency group*.`),
			},
			// Invalid Action
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigInvalidAction,
				Check:       resource.ComposeAggregateTestCheckFunc(),
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match*.`),
			},
			// Failover
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigFailover,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Restore
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigRestore,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Sync
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigSync,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Reverse
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigReverse,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Switchover
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigSwitchover,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Snapshot
			{
				Config: ProviderConfigForTesting + ReplicationConsistencyGroupActionResourceConfigSnapshot,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
