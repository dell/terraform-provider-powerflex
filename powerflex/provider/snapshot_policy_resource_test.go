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
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var localMockerSnapshotPolicyRead *Mocker

func TestAccResourceSnapshotPolicyA(t *testing.T) {
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

func TestAccResourceSnapshotPolicyUpdateFail(t *testing.T) {
	sp := make([]*scaleiotypes.SnapshotPolicy, 1)
	sp[0] = &scaleiotypes.SnapshotPolicy{
		Name:                             "snap-upadte-fail",
		SnapshotAccessMode:               "ReadOnly",
		AutoSnapshotCreationCadenceInMin: 5,
		NumOfRetainedSnapshotsPerLevel:   []int{1},
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateSnapshotPolicy).Return("sucess-create-mock", nil).Build()
					localMockerSnapshotPolicyRead = Mock(helper.GetSnapshotPolicy, OptGeneric).Return(sp, nil).Build()
				},
				Config: ProviderConfigForTesting + SPResourceUpdateWithVolFail,
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.AssignVolumeToSnapshotPolicy, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail2,
				ExpectError: regexp.MustCompile(`.*Error assigning volume to snapshot policy*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ModifySnapshotPolicy, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail4,
				ExpectError: regexp.MustCompile(`.*Error while updating auto snapshot creation cadence *.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.RenameSnapshotPolicy, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail3,
				ExpectError: regexp.MustCompile(`.*Error while updating name of snapshot policy*.`),
			},

			{
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail5,
				ExpectError: regexp.MustCompile(`.*Cannot Update Secure Snapshots after creation*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail6,
				ExpectError: regexp.MustCompile(`.*Cannot Update snapshot access mode after creation*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ModifySnapshotPolicy, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail7,
				ExpectError: regexp.MustCompile(`.*Error while updating auto snapshot creation cadence or num of retained snapshots*.`),
			},
			{
				Config:      ProviderConfigForTesting + SPResourceUpdateWithVolFail8,
				ExpectError: regexp.MustCompile(`.*Error while updating num of retained snapshots per level*.`),
			},
		},
	})
}

func TestAccResourceSnapshotPolicyCreateFail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create snapshot policy
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMockerSnapshotPolicyRead != nil {
						localMockerSnapshotPolicyRead.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateSnapshotPolicy, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceCreateFail,
				ExpectError: regexp.MustCompile(`.*Error creating snapshot policy*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.AssignVolumeToSnapshotPolicy, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SPResourceCreateWithVolFail,
				ExpectError: regexp.MustCompile(`.*Error assigning volume to snapshot policy*.`),
			},
		},
	})
}

func TestAccResourceSnapshotPolicyUpadte(t *testing.T) {
	var SPResourceCreateWithVol = createVol1 + createVol2 + createVol3 + `
		resource "powerflex_snapshot_policy" "avengers-sp-create" {
			name = "snap-create-test"
			num_of_retained_snapshots_per_level = [1]
			auto_snapshot_creation_cadence_in_min = 5
			paused = false
			secure_snapshots = false
			snapshot_access_mode = "ReadOnly"
			volume_ids = [resource.powerflex_volume.pre-req1.id,resource.powerflex_volume.pre-req2.id]
		}
	`

	var SPResourceCreateWithVolUpdate = createVol1 + createVol2 + createVol3 + `
		resource "powerflex_snapshot_policy" "avengers-sp-create" {
			name = "snap-upadte-test"
			num_of_retained_snapshots_per_level = [2,4,6]
			auto_snapshot_creation_cadence_in_min = 6
			paused = false
			secure_snapshots = false
			snapshot_access_mode = "ReadOnly"
			volume_ids = [resource.powerflex_volume.pre-req3.id, resource.powerflex_volume.pre-req2.id]
			remove_mode = "Remove"
		}
	`

	var SPResourceCreateWithVolUpdate2 = createVol1 + createVol2 + createVol3 + `
		resource "powerflex_snapshot_policy" "avengers-sp-create" {
			name = "snap-upadte-test"
			num_of_retained_snapshots_per_level = [2,4,6]
			auto_snapshot_creation_cadence_in_min = 6
			paused = false
			secure_snapshots = false
			snapshot_access_mode = "ReadOnly"
			volume_ids = [resource.powerflex_volume.pre-req2.id]
		}
	`
	t.Log(os.Getenv("TF_ACC"))
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
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

var createVol1 = `
	resource "powerflex_volume" "pre-req1"{
		name = "terraform-vol1"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
	}
`
var createVol2 = `
	resource "powerflex_volume" "pre-req2"{
		name = "terraform-vol2"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
	}
`

var createVol3 = `
	resource "powerflex_volume" "pre-req3"{
		name = "terraform-vol3"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
	}
`

var SPResourceCreate = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
  }
`
var SPResourceUpdate3 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-create-test"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 5
	paused = false
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f5437c200000003", "edd2fb3100000010"]
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
	volume_ids = ["Invalid", "5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004", "Invalid"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
  }
`
var SPResourceUpdateWithVolFail4 = `
resource "powerflex_snapshot_policy" "avengers-sp-create" {
	name = "snap-upadte-fail"
	num_of_retained_snapshots_per_level = [1]
	auto_snapshot_creation_cadence_in_min = 4
	paused = true
	secure_snapshots = false
	snapshot_access_mode = "ReadOnly"
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
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
	volume_ids = ["5f54577100000004","5f5437c200000003"]
  }
`
