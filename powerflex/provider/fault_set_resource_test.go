/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestAccFaultSetResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resourceName := "powerflex_fault_set.newFs"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test
			{
				Config: ProviderConfigForTesting + FaultSetResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_fault_set.newFs", "name", "fault-set-create"),
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
					resource.TestCheckResourceAttr("powerflex_fault_set.newFs", "name", "fault-set-update"),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFaultSetCreateNegative(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test negative
			{
				Config:      ProviderConfigForTesting + FaultSetResourceCreateNegative,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
			},
			// Create fault set Test negative
			{
				Config:      ProviderConfigForTesting + FaultSetResourceCreateNegative2,
				ExpectError: regexp.MustCompile(`.*Error creating fault set.*`),
			},
		},
	})
}

func TestAccFaultSetUpdateNegative(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test
			{
				Config: ProviderConfigForTesting + FaultSetResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_fault_set.newFs", "name", "fault-set-create"),
				),
			},
			{
				Config:      ProviderConfigForTesting + FaultSetResourceCreateNegative,
				ExpectError: regexp.MustCompile(`.*Error: Protection Domain ID cannot be updated.*`),
			},
			{
				Config:      ProviderConfigForTesting + FaultSetResourceUpdateNegative,
				ExpectError: regexp.MustCompile(`.*Error while updating name of fault set.*`),
			},
		},
	})
}

var FaultSetResourceCreate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-create"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetResourceUpdate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-update"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetResourceCreateNegative = `
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
