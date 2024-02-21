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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
