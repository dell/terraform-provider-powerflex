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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccServiceResource tests the Service Resource
func TestAccServiceResource(t *testing.T) {
	t.Skip("Skipping this test case")
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig1,
				ExpectError: regexp.MustCompile(`.*"firmware_id" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig2,
				ExpectError: regexp.MustCompile(`.*"template_id" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig3,
				ExpectError: regexp.MustCompile(`.*"deployment_name" is required.*`),
			},
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig4,
				ExpectError: regexp.MustCompile(`.*Service Template Not Found.*`),
			},
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in deploying service.*`),
			},
			//Import
			{
				Config:        ProviderConfigForTesting + importServiceTest,
				ImportState:   true,
				ImportStateId: "WRONG",
				ResourceName:  "powerflex_service.service",
				ExpectError:   regexp.MustCompile(`.*Couldn't find service with the given filter.*`),
			},
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig6,
				ExpectError: regexp.MustCompile(`.*No need to pass clone_from_host during the deployment of service.*`),
			},
		},
	})
}

func TestAccServiceResourcePositive(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + ServiceResourceCreateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_service.service", "deployment_name", "Block-Storage-Hardware"),
				),
			},
			//Update
			{
				Config: ProviderConfigForTesting + ServiceResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_service.service", "deployment_name", "Block-Storage-Hardware-Update"),
				),
			},
		},
	})
}

var importServiceTest = `
resource "powerflex_service" "service"  {
	
}
`

var ServiceResourceConfig1 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ServiceResourceConfig2 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
}
`

var ServiceResourceConfig3 = `
resource "powerflex_service" "service" {
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ServiceResourceConfig4 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "8aaa80658cd602e0018cd996a1c91bdc"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
}
`

var ServiceResourceConfig5 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "WRONG"
	template_id = "ddedf050-c429-4114-b563-3818965481d8"
}
`

var ServiceResourceConfig6 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	firmware_id = "WRONG"
	clone_from_host = "ABCD"
	template_id = "ddedf050-c429-4114-b563-3818965481d8"
}
`

var ServiceResourceCreateConfig = `
resource "powerflex_service" "service" {
	deployment_name = "Block-Storage-Hardware"
	deployment_description = "Block-Storage-Hardware"
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`

var ServiceResourceUpdateConfig = `
resource "powerflex_service" "service" {
	deployment_name = "Block-Storage-Hardware-Update"
	deployment_description = "Block-Storage-Hardware-Update"
	template_id = "4f4b69de-debb-4a5f-8f3f-44aca8259596"
	firmware_id = "8aaa804a8b4d6b5a018b4d77a75900e9"
}
`
