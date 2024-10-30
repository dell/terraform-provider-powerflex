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

var PeerSystemCreate = `
resource "powerflex_peer_system" "system_1" {  
    name = "tfacc_peer_mdm_coupling_rc"
    peer_system_id = "` + DestinationSystemID + `"
    ip_list = ["` + GatewayDataPoints.primaryMDMIP + `","` + GatewayDataPoints.secondaryMDMIP + `","` + GatewayDataPoints.tbIP + `"] 
    source_primary_mdm_information = {
         # Required fields
           ip = ""
           ssh_username = ""
           ssh_password = ""
           management_ip = ""
           management_username = ""
           management_password = ""
           #ssh_port = "22"
    }
    destination_primary_mdm_information = {
        # Required fields
        ip = ""
        ssh_username = ""
        ssh_password = ""
        management_ip = ""
        management_username = ""
        management_password = ""
        #ssh_port = "22"
   }
}
`

var PeerSystemCreateAddCert = `
resource "powerflex_peer_system" "system_1" {  
    name = "tfacc_peer_mdm_coupling_rc"
    peer_system_id = "` + DestinationSystemID + `"
    ip_list = ["` + GatewayDataPoints.primaryMDMIP + `","` + GatewayDataPoints.secondaryMDMIP + `","` + GatewayDataPoints.tbIP + `"] 
    add_certificate = true
    source_primary_mdm_information = {
         # Required fields
           ip = ""
           ssh_username = ""
           ssh_password = ""
           management_ip = ""
           management_username = ""
           management_password = ""
           #ssh_port = "22"
    }
    destination_primary_mdm_information = {
        # Required fields
        ip = ""
        ssh_username = ""
        ssh_password = ""
        management_ip = ""
        management_username = ""
        management_password = ""
        #ssh_port = "22"
   }
}
`

var PeerSystemUpdate = `
resource "powerflex_peer_system" "system_1" {  
    name = "tfacc_peer_update"
    peer_system_id = "` + DestinationSystemID + `"
    ip_list = ["` + GatewayDataPoints.primaryMDMIP + `","` + GatewayDataPoints.secondaryMDMIP + `","` + GatewayDataPoints.tbIP + `"] 
    source_primary_mdm_information = {
        # Required fields
          ip = ""
          ssh_username = ""
          ssh_password = ""
          management_ip = ""
          management_username = ""
          management_password = ""
          #ssh_port = "22"
   }
   destination_primary_mdm_information = {
       # Required fields
       ip = ""
       ssh_username = ""
       ssh_password = ""
       management_ip = ""
       management_username = ""
       management_password = ""
       #ssh_port = "22"
  }
}
`

var PeerSystemUpdateAll = `
resource "powerflex_peer_system" "system_1" {  
    name = "tfacc_peer_update"
    peer_system_id = "` + DestinationSystemID + `"
    ip_list = ["` + GatewayDataPoints.primaryMDMIP + `","` + GatewayDataPoints.secondaryMDMIP + `","` + GatewayDataPoints.tbIP + `"] 
    port = 7612
    perf_profile = "Compact"
    source_primary_mdm_information = {
        # Required fields
          ip = ""
          ssh_username = ""
          ssh_password = ""
          management_ip = ""
          management_username = ""
          management_password = ""
          #ssh_port = "22"
   }
   destination_primary_mdm_information = {
       # Required fields
       ip = ""
       ssh_username = ""
       ssh_password = ""
       management_ip = ""
       management_username = ""
       management_password = ""
       #ssh_port = "22"
  }
}
`

var PeerSystemUpdateInvalidFields = `
resource "powerflex_peer_system" "system_1" {  
    name = "tfacc_peer_update"
    peer_system_id = "invalid-update"
    ip_list = ["` + GatewayDataPoints.primaryMDMIP + `","` + GatewayDataPoints.secondaryMDMIP + `","` + GatewayDataPoints.tbIP + `"] 
    port = 7612
    perf_profile = "Compact"
    source_primary_mdm_information = {
        # Required fields
          ip = ""
          ssh_username = ""
          ssh_password = ""
          management_ip = ""
          management_username = ""
          management_password = ""
          #ssh_port = "22"
   }
   destination_primary_mdm_information = {
       # Required fields
       ip = ""
       ssh_username = ""
       ssh_password = ""
       management_ip = ""
       management_username = ""
       management_password = ""
       #ssh_port = "22"
  }
}
`

var PeerSystemBulkImport = `
// Get all of the exiting peer systems
data "powerflex_peer_system" "all_current_peer_systems" {
}
// Import all of the peers
import {
    for_each = data.powerflex_peer_system.all_current_peer_systems.peer_system_details
    to = powerflex_peer_system.imported_peer_systems[each.key]
    id = each.value.id
}

// Add them to the terraform state
resource "powerflex_peer_system" "imported_peer_systems" {
    count = length(data.powerflex_peer_system.all_current_peer_systems.peer_system_details)
    name = data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].name
    peer_system_id = data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].peer_system_id
    ip_list = [
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[0].ip,
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[1].ip,
      data.powerflex_peer_system.all_current_peer_systems.peer_system_details[count.index].ip_list[2].ip,
    ]
	source_primary_mdm_information = {
        # Required fields
          ip = ""
          ssh_username = ""
          ssh_password = ""
          management_ip = ""
          management_username = ""
          management_password = ""
          #ssh_port = "22"
   }
   destination_primary_mdm_information = {
       # Required fields
       ip = ""
       ssh_username = ""
       ssh_password = ""
       management_ip = ""
       management_username = ""
       management_password = ""
       #ssh_port = "22"
  }
}

`

// Accptance Tests
func TestAccResourceAcceptancePeerSystems(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Peer System
			{
				Config: ProviderConfigForTesting + PeerSystemCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Peer System
			{
				Config: ProviderConfigForTesting + PeerSystemUpdate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// check that import is working
			{
				ResourceName: "powerflex_peer_system.system_1",
				ImportState:  true,
			},
			// Delete is automatically tested
		},
	})
}

// Unit Tests

func TestAccResourcePeerSystemBulkImport(t *testing.T) {
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
				Config: ProviderConfigForTesting + PeerSystemBulkImport,
			},
		},
	})
}

// Unit Test
func TestAccResourcePeerSystemA(t *testing.T) {
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
					FunctionMocker = Mock(helper.CreatePeerSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerSystemCreate,
				ExpectError: regexp.MustCompile(`.*Error creating peer system*.`),
			},
			// Read Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetPeerSystem, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerSystemCreate,
				ExpectError: regexp.MustCompile(`.*Error reading peer system*.`),
			},
			// Add Cert Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.AddCertificate).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerSystemCreateAddCert,
				ExpectError: regexp.MustCompile(`.*Error adding certificate to trust store from destination primary mdm to source primary mdm*.`),
			},
			// Create successfully
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + PeerSystemCreate,
			},
			// check that import is working
			{
				ResourceName: "powerflex_peer_system.system_1",
				ImportState:  true,
			},
			// Read for update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetPeerSystem, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerSystemUpdateAll,
				ExpectError: regexp.MustCompile(`.*Error reading peer system*.`),
			},
			// Update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.PeerSystemUpdate).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + PeerSystemUpdateAll,
				ExpectError: regexp.MustCompile(`.*Error updating peer system*.`),
			},
			// Update Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + PeerSystemUpdateAll,
			},
			// Update Invalid Field Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + PeerSystemUpdateInvalidFields,
				ExpectError: regexp.MustCompile(`.*peer_system_id cannot be updated*.`),
			},
		},
	})
}
