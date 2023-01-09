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
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		},
	]
}
`
var updateSnapshotInvalidAccessModePNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadOnly"
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		},
	]
}
`

var updateSnapshotInvalidLockNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	lock_auto_snapshot = true
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadOnly"
		},
	]
}
`
var updateSnapshotMapSdcPosTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{
			sdc_id = "c423b09a00000005"
			limit_iops = 90
			limit_bw_in_mbps = 9
			access_mode = "ReadOnly"
		},
	]
}
`

// c423b09800000003

var createSnapshotAccessModeMapSdcNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-invalid-sdc-map" {
	name = "snapshots-create-zeta"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{
			sdc_id = "c4zav"
			limit_iops = 90
			limit_bw_in_mbps = 9
			access_mode = "ReadOnly"
		},
	]
}
`

var createSnapshotLockedAutoSnapshotNegTest = `
resource "powerflex_snapshot" "snapshots-create-locked-auto-invalid" {
	name = "snapshots-create-eta"
	volume_id = "4578b32d000000e9"
	access_mode = "ReadWrite"
	lock_auto_snapshot = true
	sdc_list = [
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{
			sdc_id = "c423b09a00000005"
			limit_iops = 90
			limit_bw_in_mbps = 9
			access_mode = "ReadOnly"
		},
	]
}
`

func TestAccSnapshotResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// test-1
				Config: ProviderConfigForTesting + createSnapshotPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-alpha"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "access_mode", "ReadOnly"),
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "volume_id", "4578b32d000000e9"),
				),
			},
			{
				// test-2
				Config: ProviderConfigForTesting + updateSnapshotRenamePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "name", "snapshots-create-1"),
				),
			},

			{
				// test-3
				Config: ProviderConfigForTesting + updateSnapshotResizePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "size", "24"),
				),
			},
			{
				// test-4
				Config:      ProviderConfigForTesting + updateSnapshotResizeNegTest,
				ExpectError: regexp.MustCompile(`.*Error setting snapshot size*.`),
			},
			{
				// test-5
				Config:      ProviderConfigForTesting + updateSnapshotRenameNegTest,
				ExpectError: regexp.MustCompile(`.*Error renaming the snapshot*.`),
			},

			{
				// test-6
				Config:      ProviderConfigForTesting + createSnapshotWithNonExistentVolumeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not create snapshot*.`),
			},
			{
				// test-7
				Config:      ProviderConfigForTesting + createSnapshotWithlowSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not set the size for snapshot below volume size*.`),
			},
			{
				// test-8
				Config:      ProviderConfigForTesting + createSnapshotWithhighSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Could not set snapshot size, unexpected err*.`),
			},
			{
				// test-9
				ExpectNonEmptyPlan: true,
				Config:             ProviderConfigForTesting + createSnapshotAccessModeMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.0.sdc_id", "c423b09900000004"),
				),
			},
			{
				// test-10
				ExpectNonEmptyPlan: true,
				Config:             ProviderConfigForTesting + updateSnapshotInvalidAccessModePNegTest,
				ExpectError:        regexp.MustCompile(`.*Could not set the Snapshot Access Mode*.`),
			},
			{
				// test-11
				ExpectNonEmptyPlan: true,
				Config:             ProviderConfigForTesting + updateSnapshotInvalidLockNegTest,
				ExpectError:        regexp.MustCompile(`.*Error Locking Auto Snapshots*.`),
			},
			{
				// test-12
				ExpectNonEmptyPlan: true,
				Config:             ProviderConfigForTesting + updateSnapshotMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.0.sdc_id", "c423b09900000004"),
				),
			},
			{
				// test-13
				Config:      ProviderConfigForTesting + createSnapshotAccessModeMapSdcNegTest,
				ExpectError: regexp.MustCompile(`.*Error Mapping Snapshot to SDCs*.`),
			},
			{
				// test-14
				Config:      ProviderConfigForTesting + createSnapshotLockedAutoSnapshotNegTest,
				ExpectError: regexp.MustCompile(`.*Error Locking Auto Snapshots*.`),
			},
		},
	})
}
