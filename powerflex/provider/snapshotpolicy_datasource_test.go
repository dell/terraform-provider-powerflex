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
	"fmt"
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AT
func TestAccDatasourceAcceptanceSnapshotPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceSnapshotPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving all snapshot policices
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			//retrieving snapshot policy based single filter
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first snapshot policy to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp1", "snapshotpolicies.0.id", "896a535700000000"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp1", "snapshotpolicies.0.name", "snap-create-test"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp1", "snapshotpolicies.#", "1"),
				),
			},
			//retrieving snapshot policy based multiple filter
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first snapshot policy to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp2", "snapshotpolicies.0.id", "896a535700000000"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp2", "snapshotpolicies.0.name", "snap-create-test"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp2", "snapshotpolicies.0.last_auto_snapshot_failure_in_first_level", "false"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp2", "snapshotpolicies.0.system_id", "1250de83018c2d0f"),
					resource.TestCheckResourceAttr("data.powerflex_snapshot_policy.sp2", "snapshotpolicies.#", "1"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllSnapshotPolicies).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SnapshotPolicyDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Snapshot Policy*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SnapshotPolicyDataSourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error in filtering snapshot policies*.`),
			},
		},
	})
}

var SnapshotPolicyDataSourceAll = `
data "powerflex_snapshot_policy" "all" {
}
`

var SnapshotPolicyDataSourceConfig1 = `
data "powerflex_snapshot_policy" "sp1" {
	filter{
		id = ["896a535700000000"]
	}
}
`

var SnapshotPolicyDataSourceConfig2 = `
data "powerflex_snapshot_policy" "sp2" {
	filter{
		id = ["896a535700000000"]
		name = ["snap-create-test"]
		last_auto_snapshot_failure_in_first_level = false
		system_id = ["1250de83018c2d0f"]
	}
}
`
