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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var create8gbVol = `
resource "powerflex_volume" "ref-vol"{
	name = "tfaccp-ssvol-test"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 8
}
`

var create16gbVol = `
resource "powerflex_volume" "ref-vol-16gb"{
	name = "tfaccp-16gb-ssvol-test"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 16
}
`

var createVolForSs = create8gbVol + create16gbVol

var createSnapshotPosTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_id = resource.powerflex_volume.ref-vol.id
}
`
var updateSnapshotRenamePosTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = resource.powerflex_volume.ref-vol.id
}
`

var updateSnapshotResizePosTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = resource.powerflex_volume.ref-vol.id
	size = 24
	capacity_unit="GB"
	access_mode = "ReadWrite"
}
`

var updateSnapshotResizeNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = resource.powerflex_volume.ref-vol.id
	size = 24
	capacity_unit="TB"
}
`

var updateSnapshotRenameNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-before" {
	name = "snapshot-create-invalid"
	volume_id = resource.powerflex_volume.ref-vol.id
}

resource "powerflex_snapshot" "snapshots-create" {
	name = resource.powerflex_snapshot.snapshots-create-before.name
	// snapshot with name snapshot-create-invalid already exist  
	volume_id = resource.powerflex_volume.ref-vol.id
}
`

var createSnapshotWithlowSizeNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-with-low-size" {
	name = "snapshots-create-gamma"
	volume_name = resource.powerflex_volume.ref-vol-16gb.name 
	size = 8
	capacity_unit="GB"
}
`

var createSnapshotWithhighSizeNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-with-high-size" {
	name = "snapshots-create-delta"
	volume_id = resource.powerflex_volume.ref-vol.id
	size = 5
	capacity_unit="TB"
}
`

var createSnapshotLockedAutoSnapshotNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-locked-auto-invalid" {
	name = "snapshots-create-eta"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
	lock_auto_snapshot = true
}
`

var createSnapshotWithInvalidVolumeID = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_id = "inv"
}
`

var createSnapshotWithInvalideVolumeName = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_name = "inv"
}
`

var createSnapshotWithInvalidSnapshotName = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create" {
	name = "creating-snapshot-with-more-than-31-characters"
	volume_id = resource.powerflex_volume.ref-vol.id
}
`

func TestAccSnapshotResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSnapshotPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-alpha"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "access_mode", "ReadOnly"),
					resource.TestCheckResourceAttrPair("powerflex_snapshot.snapshots-create", "volume_id", "powerflex_volume.ref-vol", "id"),
				),
			},
			// check that import is working
			{
				ResourceName: "powerflex_snapshot.snapshots-create",
				ImportState:  true,
				// TODO // ImportStateVerify: true,
			},
			{
				Config: ProviderConfigForTesting + updateSnapshotRenamePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-1"),
				),
			},
			{
				Config: ProviderConfigForTesting + updateSnapshotResizePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "size", "24"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "access_mode", "ReadWrite"),
				),
			},
			// check that import is working
			{
				ResourceName: "powerflex_snapshot.snapshots-create",
				ImportState:  true,
				// TODO // ImportStateVerify: true,
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotResizeNegTest,
				ExpectError: regexp.MustCompile(`.*Requested volume size exceeds the volume allocation limit*.`),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotRenameNegTest,
				ExpectError: regexp.MustCompile(`.*Volume name already in use*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithlowSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Volume capacity can only be increased*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithhighSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Requested volume size exceeds the volume allocation limit*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotLockedAutoSnapshotNegTest,
				ExpectError: regexp.MustCompile(`.*The specified volume is not an auto-snapshot and hence cannot be locked*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithInvalidVolumeID,
				ExpectError: regexp.MustCompile(`.*Error creating snapshot*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithInvalideVolumeName,
				ExpectError: regexp.MustCompile(`.*Error getting volume by name*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithInvalidSnapshotName,
				ExpectError: regexp.MustCompile(`.*given name exceeds the allowed length of 31 characters*.`),
			},
		},
	})
}

func TestAccSnapshotResourceDependant(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + CreateSnapshotWithVolumeName,
			},
			{
				Config: ProviderConfigForTesting + UpdateSnapshotWithVolumeName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-alpha"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "volume_name", "tfaccp-ssvol-test2"),
				),
			},
		},
	},
	)
}

func TestAccSnapshotPolicyResourceUnit(t *testing.T) {

	t.Skip("Unit Test Case")

	var SPResourceCreateWithVol = `
		resource "powerflex_snapshot_policy" "avengers-sp-create" {
			name = "sample_snap_policy_1"
			num_of_retained_snapshots_per_level = [5]
			auto_snapshot_creation_cadence_in_min = 50400
			paused = true
			secure_snapshots = false
			snapshot_access_mode = "ReadOnly"
			volume_ids = ["5f54577100000004","5f5437c200000003"]
		}
	`

	var SPResourceCreateWithVolUpdate = `
		resource "powerflex_snapshot_policy" "avengers-sp-create" {
			name = "sample_snap_policy_update"
			num_of_retained_snapshots_per_level = [2,4,6]
			auto_snapshot_creation_cadence_in_min = 6
			paused = true
			secure_snapshots = false
			snapshot_access_mode = "ReadOnly"
			volume_ids = ["5f54577100000004","5f5437c200000003"]
			remove_mode = "Remove"
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVol,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "sample_snap_policy_1"),
				),
			},
			// Update snapshot policy
			{
				Config: ProviderConfigForTesting + SPResourceCreateWithVolUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot_policy.avengers-sp-create", "name", "sample_snap_policy_update"),
				),
			},
		},
	})
}

var CreateSnapshotWithVolumeName = `
resource "powerflex_volume" "ref-vol"{
	name = "tfaccp-ssvol-test"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 8
}
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_name = resource.powerflex_volume.ref-vol.name
}
`
var UpdateSnapshotWithVolumeName = `
resource "powerflex_volume" "ref-vol"{
	name = "tfaccp-ssvol-test2"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1"
	size = 8
}
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_name = resource.powerflex_volume.ref-vol.name
}
`
