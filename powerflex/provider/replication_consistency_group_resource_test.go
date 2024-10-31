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

var RPCCreate = `
resource "powerflex_replication_consistency_group" "example" {
    name = "tftest-rcg"
    protection_domain_id = "` + SourceProtectionDomainID + `"
    remote_protection_domain_id = "` + RemoteProtectionDomainID + `"
    destination_system_id = "` + DestinationSystemID + `"
  }
`

var RPCCreateUpdateOnlyFieldError = `
resource "powerflex_replication_consistency_group" "example" {
    name = "tftest-rcg"
    protection_domain_id = "` + SourceProtectionDomainID + `"
    remote_protection_domain_id = "` + RemoteProtectionDomainID + `"
    destination_system_id = "` + DestinationSystemID + `"
	freeze_state = "Frozen"
  }
`

var RPCUpdateName = `
resource "powerflex_replication_consistency_group" "example" {
    name = "tftest-rcg-update"
    protection_domain_id = "` + SourceProtectionDomainID + `"
    remote_protection_domain_id = "` + RemoteProtectionDomainID + `"
    destination_system_id = "` + DestinationSystemID + `"
  }
`
var RPCUpdateFieldWhichCannontBeUpdated = `
resource "powerflex_replication_consistency_group" "example" {
    name = "tftest-rcg"
    protection_domain_id = "updated-domain"
    remote_protection_domain_id = "updated-domain"
    destination_system_id = "updated-destination"
  }
`

var RPCFullUpdate = `
resource "powerflex_replication_consistency_group" "example" {
    name = "tftest-rcg-update"
    protection_domain_id = "` + SourceProtectionDomainID + `"
    remote_protection_domain_id = "` + RemoteProtectionDomainID + `"
    destination_system_id = "` + DestinationSystemID + `"
	freeze_state = "Frozen"
	rpo_in_seconds = 20
	local_activity_state = "Terminated"
	target_volume_access_mode = "ReadOnly"
	pause_mode = "Pause"
	curr_consist_mode = "Inconsistent"
  }
`

var RPCBulkImport = `
// Get all of the existing replication consistency groups
data "powerflex_replication_consistency_group" "existing" {
}

// Import all of the replication consistency groups
import {
    for_each = data.powerflex_replication_consistency_group.existing.replication_consistency_group_details
    to = powerflex_replication_consistency_group.this[each.key]
    id = each.value.id
}

// Add them to the terraform state
resource "powerflex_replication_consistency_group" "this" {
    count = length(data.powerflex_replication_consistency_group.existing.replication_consistency_group_details)
    name = data.powerflex_replication_consistency_group.existing.replication_consistency_group_details[count.index].name
    protection_domain_id = data.powerflex_replication_consistency_group.existing.replication_consistency_group_details[count.index].protection_domain_id
    remote_protection_domain_id = data.powerflex_replication_consistency_group.existing.replication_consistency_group_details[count.index].remote_protection_domain_id
    destination_system_id = data.powerflex_replication_consistency_group.existing.replication_consistency_group_details[count.index].destination_system_id
}
`

// Accptance Tests
func TestAccResourceAcceptanceReplicationConsistencyGroup(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Replication Consistency Group
			{
				Config: ProviderConfigForTesting + RPCCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Replication Consistency Group
			{
				Config: ProviderConfigForTesting + RPCUpdateName,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// check that import is working
			{
				ResourceName: "powerflex_replication_consistency_group.example",
				ImportState:  true,
			},
			// Delete is automatically tested
		},
	})
}

// Unit Tests

func TestAccResourceReplicationConsistencyGroupBulkImport(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// bulk Import success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + RPCBulkImport,
			},
		},
	})
}

func TestAccResourceReplicationConsistencyGroup(t *testing.T) {
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
					FunctionMocker = Mock(helper.CreateReplicationConsistencyGroup).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + RPCCreate,
				ExpectError: regexp.MustCompile(`.*Error creating replication consistency group*.`),
			},
			// Create Error when an update only field is changed
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + RPCCreateUpdateOnlyFieldError,
				ExpectError: regexp.MustCompile(`.*are not able to be modified from default values while creating a new replication consistency group*.`),
			},
			// Read Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetSpecificReplicationConsistencyGroup, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + RPCCreate,
				ExpectError: regexp.MustCompile(`.*Error reading replication consistency group*.`),
			},
			// Create success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + RPCCreate,
			},
			// check that import is working
			{
				ResourceName: "powerflex_replication_consistency_group.example",
				ImportState:  true,
			},
			// Read Error after update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetSpecificReplicationConsistencyGroup, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + RPCCreate,
				ExpectError: regexp.MustCompile(`.*Error reading replication consistency group*.`),
			},
			// Update Error from field which cannot be updated
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + RPCUpdateFieldWhichCannontBeUpdated,
				ExpectError: regexp.MustCompile(`.*protection_domain_id, remote_protection_domain_id, and destination_system_id cannot be updated*.`),
			},
			// Update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.RCGUpdates).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + RPCUpdateName,
				ExpectError: regexp.MustCompile(`.*Error updating replication consistency group*.`),
			},
			// Update Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + RPCFullUpdate,
			},
		},
	})
}
