package powerflex

import (
	"regexp"
	"testing"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var createSnapshotPosTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-alpha"
	volume_id = "4578b32d000000e9"
}
`
var updateSnapshotRenamePosTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4578b32d000000e9"
}
`

var updateSnapshotResizePosTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4578b32d000000e9"
	size = 24
	capacity_unit="GB"
}
`

var updateSnapshotResizeNegTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4578b32d000000e9"
	size = 24
	capacity_unit="TB"
}
`

// snapshot-create-invalid is already created using UI on powerflex. so if we try to rename this to an existing snapshot name, it will throw an error.
var updateSnapshotRenameNegTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshot-create-invalid"
	volume_id = "4578b32d000000e9"
}
`

var createSnapshotWithNonExistentVolumeNegTest = `
resource "powerflex_snapshot" "snapshots-create-without-volume-id" {
	name = "snapshots-create-beta"
	volume_id = "abc"
}
`

var createSnapshotWithlowSizeNegTest = `
resource "powerflex_snapshot" "snapshots-create-with-low-size" {
	name = "snapshots-create-gamma"
	volume_id = "4578b32d000000e9"
	size = 5
	capacity_unit="GB"
}
`

var createSnapshotWithhighSizeNegTest = `
resource "powerflex_snapshot" "snapshots-create-with-high-size" {
	name = "snapshots-create-delta"
	volume_id = "4578b32d000000e9"
	size = 5
	capacity_unit="TB"
}
`

var createSnapshotAccessModeMapSdcPosTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	map_sdcs_id = ["c423b09800000003"]
}
`
var updateSnapshotInvalidAccessModePNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadOnly"
	map_sdcs_id = ["c423b09800000003"]
}
`

var updateSnapshotInvalidLockNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	locked_auto_snapshot = true
	map_sdcs_id = ["c423b09800000003"]
}
`
var updateSnapshotMapSdcPosTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
}
`

var createSnapshotAccessModeMapSdcNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-invalid-sdc-map" {
	name = "snapshots-create-zeta"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	map_sdcs_id = ["c423b09800000003","c423"]
}
`

var createSnapshotLockedAutoSnapshotNegTest = `
resource "powerflex_snapshot" "snapshots-create-locked-auto-invalid" {
	name = "snapshots-create-eta"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	map_sdcs_id = ["c423b09800000003"]
	locked_auto_snapshot = true
}
`

func TestAccSnapshotResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSnapshotPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-alpha"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "access_mode", "ReadOnly"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "volume_id", "4578b32d000000e9"),
				),
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
				),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotResizeNegTest,
				ExpectError: regexp.MustCompile(`.*Error setting snapshot size*.`),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotRenameNegTest,
				ExpectError: regexp.MustCompile(`.*Error renaming the snapshot*.`),
			},

			{
				Config:      ProviderConfigForTesting + createSnapshotWithoutVolumeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not create snapshot*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithlowSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not set the size for snapshot below volume size*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithhighSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not set snapshot size, unexpected err*.`),
			},
			{
				Config: ProviderConfigForTesting + createSnapshotAccessModeMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "map_sdcs_id.0", "c423b09800000003"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotInvalidAccessModePNegTest,
				ExpectError: regexp.MustCompile(`.*Could not set the Snapshot Access Mode*.`),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotInvalidLockNegTest,
				ExpectError: regexp.MustCompile(`.*Error Locking Auto Snapshots*.`),
			},
			{
				Config: ProviderConfigForTesting + updateSnapshotMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "map_sdcs_id.0", "c423b09a00000005"),
				),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotAccessModeMapSdcNegTest,
				ExpectError: regexp.MustCompile(`.*Error Mapping Snapshot to SDCs*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotLockedAutoSnapshotNegTest,
				ExpectError: regexp.MustCompile(`.*Error Locking Auto Snapshots*.`),
			},
		},
	})
}
