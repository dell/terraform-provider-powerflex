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
	volume_id = "4577c84000000120"
}
`
var updateSnapshotRenamePosTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4577c84000000120"
}
`

var updateSnapshotResizePosTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4577c84000000120"
	size = 24
	capacity_unit="GB"
}
`

var updateSnapshotResizeNegTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create-1"
	volume_id = "4577c84000000120"
	size = 24
	capacity_unit="TB"
}
`

// snapshot-create-invalid is already created using UI on powerflex. so if we try to rename this to an existing snapshot name, it will throw an error.
var updateSnapshotRenameNegTest = `
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshot-create-invalid"
	// snapshot with name snapshot-create-invalid already exist  
	volume_id = "4577c84000000120"
}
`

var createSnapshotWithlowSizeNegTest = `
resource "powerflex_snapshot" "snapshots-create-with-low-size" {
	name = "snapshots-create-gamma"
	volume_name = "volume-ses" 
	// volume-ses has size of 16GB, so below 8 GB config of snapshot will be failing
	size = 8
	capacity_unit="GB"
}
`

var createSnapshotWithhighSizeNegTest = `
resource "powerflex_snapshot" "snapshots-create-with-high-size" {
	name = "snapshots-create-delta"
	volume_id = "4577c84000000120"
	size = 5
	capacity_unit="TB"
}
`

var createSnapshotAccessModeMapSdcPosTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4577c84000000120"
	access_mode = "ReadWrite"
  	size = 16
  	capacity_unit = "GB"
  	remove_mode = "INCLUDING_DESCENDANTS"
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
	volume_id = "4577c84000000120"
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
	volume_id = "4577c84000000120"
	access_mode = "ReadWrite"
  	size = 16
  	capacity_unit = "GB"
  	lock_auto_snapshot = true
  	remove_mode = "INCLUDING_DESCENDANTS"
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

var updateSnapshotMapSdcPosTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = "4577c84000000120"
	access_mode = "ReadWrite"
  size = 16
  capacity_unit = "GB"
  remove_mode = "INCLUDING_DESCENDANTS"
	sdc_list = [
    {	
			sdc_id = "c423b09800000003"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{	
			sdc_id = "c423b09900000004"
			limit_iops = 190
			limit_bw_in_mbps = 70
			access_mode = "NoAccess"
		},
    		{
			sdc_id = "c423b09a00000005"
			limit_iops = 82
			limit_bw_in_mbps = 17
			access_mode = "ReadOnly"
		},
	]
}
`

var createSnapshotAccessModeMapSdcNegTest = `
resource "powerflex_snapshot" "snapshots-create-access-mode-invalid-sdc-map" {
	name = "snapshots-create-zeta"
	volume_id = "4577c84000000120"
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
	volume_id = "4577c84000000120"
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
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create", "volume_id", "4577c84000000120"),
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
				Config: ProviderConfigForTesting + createSnapshotAccessModeMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.#", "1"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotInvalidAccessModePNegTest,
				ExpectError: regexp.MustCompile(`.*The command cannot be applied because the volume has read-write mappings*.`),
			},
			{
				Config:      ProviderConfigForTesting + updateSnapshotInvalidLockNegTest,
				ExpectError: regexp.MustCompile(`.*The specified volume is not an auto-snapshot and hence cannot be locked*.`),
			},
			{
				Config: ProviderConfigForTesting + updateSnapshotMapSdcPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.#", "3"),
				),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotAccessModeMapSdcNegTest,
				ExpectError: regexp.MustCompile(`.*Couldn't find SDC*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotLockedAutoSnapshotNegTest,
				ExpectError: regexp.MustCompile(`.*The specified volume is not an auto-snapshot and hence cannot be locked*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithInvalidVolumeID,
				ExpectError: regexp.MustCompile(`.*Error getting volume by id*.`),
			},
			{
				Config:      ProviderConfigForTesting + createSnapshotWithInvalideVolumeName,
				ExpectError: regexp.MustCompile(`.*Error getting volume by name*.`),
			},
		},
	})
}
