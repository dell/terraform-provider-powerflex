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
	"fmt"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceFaultSet(t *testing.T) {
	resourceName := "powerflex_fault_set.newFs"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test
			{
				Config: ProviderConfigForTesting + FaultSetResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_fault_set.newFs", "name", "fault-set-create-test"),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update fault set Test
			{
				Config: ProviderConfigForTesting + FaultSetResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_fault_set.newFs", "name", "fault-set-update-test"),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to invalid protection domain should fail
			{
				Config:      ProviderConfigForTesting + FaultSetResourceInvalidPNegative,
				ExpectError: regexp.MustCompile(`.*Error: Protection Domain ID cannot be updated.*`),
			},
			// Should show failure if unable to update the name of the fault set
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ModifyFaultSetName, OptGeneric).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetResourceUpdateNegative,
				ExpectError: regexp.MustCompile(`.*Error while updating name of fault set.*`),
			},
		},
	})
}

func TestAccResourceFaultSetCreateNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test negative
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetNewProtectionDomainEx).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetResourceInvalidPNegative,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
			},
			// Create fault set Test negative
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateFaultSet, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FaultSetResourceCreateNegative2,
				ExpectError: regexp.MustCompile(`.*Error creating fault set.*`),
			},
		},
	})
}

var FaultSetResourceCreate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-create-test"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetResourceUpdate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-update-test"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetResourceInvalidPNegative = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-create"
	protection_domain_id = "Invalid"
}
`

var FaultSetResourceCreateNegative2 = `
resource "powerflex_fault_set" "newFs" {
	name = "fault set create"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetResourceUpdateNegative = `
resource "powerflex_fault_set" "newFs" {
	name = "fault set update"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`
