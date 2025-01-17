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
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// UT TestAccResourceGroupResource tests the ResourceGroup Resource
func TestAccResourceGroupResourceNegative(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig1,
				ExpectError: regexp.MustCompile(`.*"firmware_id" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig2,
				ExpectError: regexp.MustCompile(`.*"template_id" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig3,
				ExpectError: regexp.MustCompile(`.*"deployment_name" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig4,
				ExpectError: regexp.MustCompile(`.*Service Template Not Found.*`),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in deploying ResourceGroup.*`),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceConfig6,
				ExpectError: regexp.MustCompile(`.*No need to pass clone_from_host during the deployment of ResourceGroup*`),
			},
		},
	})
}

// UT
func TestAccResourceGroupResourcePositive(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + ResourceGroupResourceCreateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_resource_group.data", "deployment_name", "Block-Storage-Hardware"),
				),
			},
			// check that import is working
			{
				ResourceName: "powerflex_resource_group.data",
				ImportState:  true,
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetServiceDetailsByID).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ResourceGroupResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*Error in getting ResourceGroup details*.`),
			},
			// Error updating node value
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + ResourceGroupResourceUpdateNodeError,
				ExpectError: regexp.MustCompile(`.*please validate your inputs*.`),
			},
			// Error updating template id value
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceUpdateTemplateError,
				ExpectError: regexp.MustCompile(`.*please validate your inputs*.`),
			},
			// Error updating firmware id value
			{
				Config:      ProviderConfigForTesting + ResourceGroupResourceUpdateFirmwareError,
				ExpectError: regexp.MustCompile(`.*please validate your inputs*.`),
			},
			// Update error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).UpdateService).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ResourceGroupResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*Error in deploying service*.`),
			},
			// Update success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + ResourceGroupResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_resource_group.data", "deployment_name", "Block-Storage-Hardware-Update"),
				),
			},
		},
	})
}

var ResourceGroupResourceConfig1 = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ResourceGroupResourceConfig2 = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
}
`

var ResourceGroupResourceConfig3 = `
resource "powerflex_resource_group" "data" {
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ResourceGroupResourceConfig4 = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ResourceGroupResourceConfig5 = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "WRONG"
	template_id = "ddedf050-c429-4114-b563-3818965481d8"
}
`

var ResourceGroupResourceConfig6 = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "WRONG"
	clone_from_host = "ABCD"
	template_id = "ddedf050-c429-4114-b563-3818965481d8"
}
`

var ResourceGroupResourceUpdateNodeError = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Block-Storage-Hardware-Update"
	deployment_description = "Block-Storage-Hardware-Update"
	deployment_timeout = 10
	nodes = 1
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`

var ResourceGroupResourceUpdateTemplateError = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Block-Storage-Hardware-Update"
	deployment_description = "Block-Storage-Hardware-Update"
	deployment_timeout = 10
	template_id = "invalid"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`

var ResourceGroupResourceUpdateFirmwareError = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Block-Storage-Hardware-Update"
	deployment_description = "Block-Storage-Hardware-Update"
	deployment_timeout = 10
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "invalid-firware-id"
}
`

var ResourceGroupResourceCreateConfig = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Block-Storage-Hardware"
	deployment_description = "Block-Storage-Hardware"
	deployment_timeout = 10
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`

var ResourceGroupResourceUpdateConfig = `
resource "powerflex_resource_group" "data" {
	deployment_name = "Block-Storage-Hardware-Update"
	deployment_description = "Block-Storage-Hardware-Update"
	deployment_timeout = 10
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`
