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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccFirmwareRepositoryResource(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			// Update fault set Test
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceModifyPos,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
		},
	})
}

func TestAccFirmwareRepositoryCreateNegative(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg,
				ExpectError: regexp.MustCompile(`.*Error: Approve cannot be set to false once it is set to true.*`),
			},
		},
	})
}

func TestAccFirmwareRepositoryUpdateNegative(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create fault set Test
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceCreate2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg1,
				ExpectError: regexp.MustCompile(`.*Error: Source Location cannot be updated.*`),
			},
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg2,
				ExpectError: regexp.MustCompile(`.*Error: Username cannot be updated.*`),
			},
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg3,
				ExpectError: regexp.MustCompile(`.*Error: Password cannot be updated.*`),
			},
		},
	})
}

var FirmwareRepositoryResourceCreate = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = false
}
`

var FirmwareRepositoryResourceModifyPos = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = true
}
`

var FirmwareRepositoryResourceCreate2 = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	username = "DummyUser"
	password = "DummyPassword"
	approve = false
}
`

var FirmwareRepositoryResourceModifyNeg1 = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "Source location updated"
	username = "DummyUser"
	password = "DummyPassword"
	approve = false
}
`

var FirmwareRepositoryResourceModifyNeg2 = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	username = "DummyUpdatedUser"
	password = "DummyPassword"
	approve = false
}
`

var FirmwareRepositoryResourceModifyNeg3 = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	username = "DummyUser"
	password = "DummyUpdatedPassword"
	approve = false
}
`

var FirmwareRepositoryResourceCreate3 = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = true
}
`

var FirmwareRepositoryResourceModifyNeg = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = false
}
`
