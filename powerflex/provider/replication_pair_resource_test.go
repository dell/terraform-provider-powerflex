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

var ReplicationPairCreate = `
resource "powerflex_replication_pair" "rp" {
    name                              = "` + RpName + `"
    source_volume_id                  = "` + RpSourceVolumeID + `"
    destination_volume_id             = "` + RpDestinationVolumeID + `"
    replication_consistency_group_id  = "` + RcgID + `"
    pause_initial_copy = true
}
`

var ReplicationPairUpdate = `
resource "powerflex_replication_pair" "rp" {
    name                              = "` + RpName + `"
    source_volume_id                  = "` + RpSourceVolumeID + `"
    destination_volume_id             = "` + RpDestinationVolumeID + `"
    replication_consistency_group_id  = "` + RcgID + `"
    pause_initial_copy = false
}
`

var ReplicationPairInvalidUpdate = `
resource "powerflex_replication_pair" "rp" {
    name                              = "invalid"
    source_volume_id                  = "invalid"
    destination_volume_id             = "invalid"
    replication_consistency_group_id  = "invalid"
    pause_initial_copy = true
}
`

// Accptance Tests
func TestAccResourceAcceptanceReplicationPair(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Replication Pair
			{
				Config: ProviderConfigForTesting + ReplicationPairCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Replication Pair
			{
				Config: ProviderConfigForTesting + ReplicationPairCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceReplicationPair(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateReplicationPair).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairCreate,
				ExpectError: regexp.MustCompile(`.*Error creating replication pair*.`),
			},
			// Pausing Replication in create method Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.PauseReplicationPair, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairCreate,
				ExpectError: regexp.MustCompile(`.*Error pausing replication pair after create*.`),
			},
			// Reading Replication Pair Error after create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetSpecificReplicationPair, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairCreate,
				ExpectError: regexp.MustCompile(`.*Error reading replication pair*.`),
			},
			// Successful create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + ReplicationPairCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// check that import is working
			{
				ResourceName: "powerflex_replication_pair.rp",
				ImportState:  true,
			},
			// Update Invalid fields should show Error
			{
				Config:      ProviderConfigForTesting + ReplicationPairInvalidUpdate,
				ExpectError: regexp.MustCompile(`.*name, source_volume_id, replication_consistency_group_id, and destination_volume_id cannot be updated*.`),
			},
			// Resume Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ResumeReplicationPair, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairUpdate,
				ExpectError: regexp.MustCompile(`.*Error resuming replication pair, only avaiable during initial copy*.`),
			},
			// Update Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + ReplicationPairUpdate,
			},
			// Pause Update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.PauseReplicationPair, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairCreate,
				ExpectError: regexp.MustCompile(`.*Error pausing replication pair, only avaiable during initial copy*.`),
			},
			// Reading Replication Pair Error after update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetSpecificReplicationPair, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ReplicationPairCreate,
				ExpectError: regexp.MustCompile(`.*Error reading replication pair*.`),
			},
		},
	})
}
