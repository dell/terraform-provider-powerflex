package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

var createSnapshotAccessModeMapSdcPosTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
  	size = 16
  	capacity_unit = "GB"
  	remove_mode = "INCLUDING_DESCENDANTS"
	sdc_list = [
		{	
			sdc_id = "e3ce1fb500000000"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		},
	]
}
`

var updateSnapshotInvalidAccessModePNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadOnly"
	sdc_list = [
		{	
			sdc_id = "e3ce1fb500000000"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		},
	]
}
`

var updateSnapshotInvalidLockNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
  	size = 16
  	capacity_unit = "GB"
  	lock_auto_snapshot = true
  	remove_mode = "INCLUDING_DESCENDANTS"
	sdc_list = [
		{	
			sdc_id = "e3ce1fb500000000"
			limit_iops = 150
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		},
	]
}
`

var updateSnapshotMapSdcPosTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
	name = "snapshots-create-epsilon"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
  size = 16
  capacity_unit = "GB"
  remove_mode = "INCLUDING_DESCENDANTS"
	sdc_list = [
    {	
			sdc_id = "e3ce1fb500000000"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{	
			sdc_id = "e3ce46c500000002"
			limit_iops = 190
			limit_bw_in_mbps = 70
			access_mode = "NoAccess"
		},
    		{
			sdc_id = "e3ce46c600000003"
			limit_iops = 82
			limit_bw_in_mbps = 17
			access_mode = "ReadOnly"
		},
	]
}
`

var createSnapshotAccessModeMapSdcNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-access-mode-invalid-sdc-map" {
	name = "snapshots-create-zeta"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
	sdc_list = [
		{	
			sdc_id = "e3ce1fb500000000"
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

var createSnapshotLockedAutoSnapshotNegTest = createVolForSs + `
resource "powerflex_snapshot" "snapshots-create-locked-auto-invalid" {
	name = "snapshots-create-eta"
	volume_id = resource.powerflex_volume.ref-vol.id
	access_mode = "ReadWrite"
	lock_auto_snapshot = true
	sdc_list = [
		{	
			sdc_id = "e3ce1fb500000000"
			limit_iops = 200
			limit_bw_in_mbps = 40
			access_mode = "ReadWrite"
		},
		{
			sdc_id = "e3ce46c600000003"
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

func TestAccSnapshotResourceDuplicateSdc(t *testing.T) {
	createSsDuplicateSdcPos := createVolForSs + `
	resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
		name = "snapshots-create-epsilon"
		volume_id = resource.powerflex_volume.ref-vol.id
		access_mode = "ReadWrite"
		size = 16
		sdc_list = [
			{	
				sdc_id = "e3ce1fb500000000"
				limit_iops = 200
				limit_bw_in_mbps = 40
				access_mode = "ReadWrite"
			},
			{	
				sdc_id = "e3ce1fb500000000"
				limit_iops = 200
				limit_bw_in_mbps = 40
				access_mode = "ReadWrite"
			}
		]
	}
	`
	createSsDuplicateSdcInv := createVolForSs + `
	resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
		name = "snapshots-create-epsilon"
		volume_id = resource.powerflex_volume.ref-vol.id
		access_mode = "ReadWrite"
		size = 16
		sdc_list = [
			{	
				sdc_id = "e3ce1fb500000000"
				limit_iops = 200
				limit_bw_in_mbps = 40
				access_mode = "ReadWrite"
			},
			{	
				sdc_id = "e3ce1fb500000000"
				limit_iops = 200
				limit_bw_in_mbps = 45
				access_mode = "ReadWrite"
			}
		]
	}
	`

	createSsNoSdc := createVolForSs + `
	resource "powerflex_snapshot" "snapshots-create-access-mode-sdc-map" {
		name = "snapshots-create-epsilon"
		volume_id = resource.powerflex_volume.ref-vol.id
		access_mode = "ReadWrite"
		size = 16
		sdc_list = []
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + createSsDuplicateSdcInv,
				ExpectError: regexp.MustCompile(`.*Error: Duplicate SDC in list*.`),
			},
			{
				Config: ProviderConfigForTesting + createSsDuplicateSdcPos,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + createSsNoSdc,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_snapshot.snapshots-create-access-mode-sdc-map", "sdc_list.#", "0"),
				),
			},
		},
	})
}
