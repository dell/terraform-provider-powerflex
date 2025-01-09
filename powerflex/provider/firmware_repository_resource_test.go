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
	"testing"

	. "github.com/bytedance/mockey"
	goscaleio "github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var FunctionMockerFirmwareRepository *Mocker

// AT
func TestAccResourceAcceptanceFirmwareRepository(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			{
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceModifyPos,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
		},
	})
}

// UT
func TestAccResourceFirmwareRepository(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 1 Create Gateway Connection Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).TestConnection).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				ExpectError: regexp.MustCompile(`.*Please provide valid credentials*.`),
			},
			// 2 Create Upload Compliance Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).UploadCompliance).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				ExpectError: regexp.MustCompile(`.*Could not Upload the compliance File*.`),
			},
			// 3 Create Read Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				ExpectError: regexp.MustCompile(`.*Could not get the Firmware Repository Details*.`),
			},
			// 4 Create Error State Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(&scaleiotypes.UploadComplianceTopologyDetails{
						State: "errors",
						ID:    "1234",
					}, nil).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				ExpectError: regexp.MustCompile(`.*Could not Upload the compliance File*.`),
			},
			// 5 Create Needs Approval State Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(&scaleiotypes.UploadComplianceTopologyDetails{
						State: "needsApproval",
						ID:    "1234",
					}, nil).Build()
					FunctionMockerFirmwareRepository = Mock((*goscaleio.GatewayClient).ApproveUnsignedFile).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceUpdate,
				ExpectError: regexp.MustCompile(`.*Could not approve the compliance File*.`),
			},
			// 6 Create Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerFirmwareRepository != nil {
						FunctionMockerFirmwareRepository.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceCreate3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			// 7 Import resource
			{
				ResourceName: "powerflex_firmware_repository.newFr",
				ImportState:  true,
			},
			// 8 Read Error after create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceUpdate,
				ExpectError: regexp.MustCompile(`.*Could not get the Firmware Repository Details*.`),
			},
			// 9 Update Needs Approval State Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(&scaleiotypes.UploadComplianceTopologyDetails{
						State:          "needsApproval",
						SourceLocation: SourceLocation,
						ID:             "1234",
					}, nil).Build()
					FunctionMockerFirmwareRepository = Mock((*goscaleio.GatewayClient).ApproveUnsignedFile).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceUpdate,
				ExpectError: regexp.MustCompile(`.*Could not Upload the compliance File*.`),
			},
			// 10 Update Error State Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerFirmwareRepository != nil {
						FunctionMockerFirmwareRepository.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetUploadComplianceDetails).Return(&scaleiotypes.UploadComplianceTopologyDetails{
						State:          "errors",
						SourceLocation: SourceLocation,
						ID:             "1234",
					}, nil).Build()
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceUpdate,
				ExpectError: regexp.MustCompile(`.*Could not Upload the compliance File*.`),
			},
			// 11 Update Without Approve Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceUpdateWithoutApprove,
				ExpectError: regexp.MustCompile(`.*Approve attribute needs to be updated*.`),
			},
			// 12 Update Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + FirmwareRepositoryResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_firmware_repository.newFr", "source_location", SourceLocation),
				),
			},
			// 13 Update Error Approve
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg,
				ExpectError: regexp.MustCompile(`.*Error: Approve cannot be set to false once it is set to true.*`),
			},
			// 14 Update Error Source Location
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg1,
				ExpectError: regexp.MustCompile(`.*Error: Source Location cannot be updated.*`),
			},
			// 15 Update Error Username
			{
				Config:      ProviderConfigForTesting + FirmwareRepositoryResourceModifyNeg2,
				ExpectError: regexp.MustCompile(`.*Error: Username cannot be updated.*`),
			},
			// 16 Update Error Password
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
	approve = false
	username = "DummyUser"
	password = "DummyPassword"
}
`

var FirmwareRepositoryResourceUpdate = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = true
	username = "DummyUser"
	password = "DummyPassword"
}
`

var FirmwareRepositoryResourceUpdateWithoutApprove = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = false
	username = "DummyUser"
	password = "DummyPassword"
	timeout = 41
}
`

var FirmwareRepositoryResourceModifyNeg = `
resource "powerflex_firmware_repository" "newFr" {
	source_location = "` + SourceLocation + `"
	approve = false
	username = "DummyUser"
	password = "DummyPassword"
}
`
