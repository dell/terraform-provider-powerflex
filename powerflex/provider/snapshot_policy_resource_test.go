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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccSnapshotPolicyResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy test
			{
				Config: ProviderConfigForTesting + SPResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-create-test"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicyResourceUpadte(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVol,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-create-test"),
				),
			},
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVolUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-create-test"),
				),
			},
		},
	})
}

var SPResourceCreate = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
  }
`

var SPResourceCreateFail = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap create fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
  }
`

var SPResourceCreateWithVolFail = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceCreateWithVol = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3200000008"]
  }
`

var SPResourceCreateWithVolUpdate = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_id = ["edd322270000000a", "edd2fb3200000008"]
	remove_mode = "Remove"
  }
`
