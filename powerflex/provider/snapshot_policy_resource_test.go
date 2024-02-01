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
	"regexp"
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
			{
				Config: ProviderConfigForTesting + SPResourceUpdate2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "paused", "true"),
				),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdate3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "paused", "false"),
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
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-upadte-test"),
				),
			},
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVolUpdate2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-upadte-test"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicyResourceCreateFail(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				Config: ProviderConfigForTesting + SPResourceCreateFail,				
				ExpectError: regexp.MustCompile(`.*Error creating snapshot policy*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVolFail,
				ExpectError: regexp.MustCompile(`.*Error assigning volume to snapshot policy*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVolFail2,
				ExpectError: regexp.MustCompile(`.*Error assigning volume to snapshot policy*.`),
			},
		},
	})
}


func TestAccSnapshotPolicyResourceUpdateFail(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail,				
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "snap-upadte-fail"),
				),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail2,
				ExpectError: regexp.MustCompile(`.*Error assigning volume to snapshot policy*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail3,
				ExpectError: regexp.MustCompile(`.*Error while updating name of snapshot policy*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail4,
				ExpectError: regexp.MustCompile(`.*Error while updating auto snapshot creation cadence *.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail5,
				ExpectError: regexp.MustCompile(`.*Cannot Update Secure Snapshots after creation*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail6,
				ExpectError: regexp.MustCompile(`.*Cannot Update snapshot access mode after creation*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail7,
				ExpectError: regexp.MustCompile(`.*Error while updating auto snapshot creation cadence or num of retained snapshots*.`),
			},
			{
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail8,
				ExpectError: regexp.MustCompile(`.*Error while updating num of retained snapshots per level*.`),
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

var SPResourceUpdate2 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [2,4,6]
	auto_snapshot_creation_cadence_in_min = 5
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
  }
`
var SPResourceUpdate3 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [2,4,6]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
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
	volume_ids = ["edd2fb3100000007","edd2fb3200000008"]
  }
`

var SPResourceCreateWithVolUpdate = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-test"
	num_of_retained_snapshots_per_level = [2,4,6]
	auto_snapshot_creation_cadence_in_min = 6
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd322270000000a","edd2fb3200000008"]
	remove_mode = "Remove"
  }
`
var SPResourceCreateWithVolUpdate2 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-test"
	num_of_retained_snapshots_per_level = [2,4,6]
	auto_snapshot_creation_cadence_in_min = 6
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3200000008"]
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
resource "powerflex_snapshot_policy" "avengers-sp-create2" {
	name = "snap-create-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3100000010"]
  }
`

var SPResourceCreateWithVolFail2 = `
resource "powerflex_snapshot_policy" "avengers-sp-create3" {
	name = "snap-create-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["Invalid", "edd2fb3100000007"]
  }
`

var SPResourceUpdateWithVolFail = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceUpdateWithVolFail2 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "Invalid"]
  }
`

var SPResourceUpdateWithVolFail3 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap upadte fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`
var SPResourceUpdateWithVolFail4 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 1
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceUpdateWithVolFail5 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 1
	paused = true
	secure_snapshots = true
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceUpdateWithVolFail6 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 1
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadWrite"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceUpdateWithVolFail7 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1,4,6]
	auto_snapshot_creation_cadence_in_min = 1
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`

var SPResourceUpdateWithVolFail8 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1,4,6]
	auto_snapshot_creation_cadence_in_min = 5
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["edd2fb3100000007", "edd2fb3200000008"]
  }
`