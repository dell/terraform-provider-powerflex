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
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccServiceResource tests the Service Resource
func TestAccServiceResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig1,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			//Update
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
		},
	})
}

// TestAccServiceResource tests the Service Resource
func TestAccServiceResourceImport(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Import
			{
				Config:        ProviderConfigForTesting + importServiceTest,
				ImportState:   true,
				ImportStateId: "8aaaee208da6a8bc018dc256df790c4b",
				ResourceName:  "powerflex_service.service",
				Check: resource.ComposeAggregateTestCheckFunc(
					validateServiceDetails,
				),
			},
			//Update
			{
				Config:      ProviderConfigForTesting + ServiceResourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
		},
	})
}

func validateServiceDetails(state *terraform.State) error {
	// Retrieve the resource instance
	serviceResource, ok := state.RootModule().Resources["powerflex_service.service"]
	if !ok {
		return fmt.Errorf("Failed to find powerflex_service.service in state")
	}

	// Get the value of the "deployment_name" attribute from the resource instance
	deploymentName, ok := serviceResource.Primary.Attributes["deployment_name"]
	if !ok {
		return fmt.Errorf("sdc_list attribute not found in state")
	}

	// Check if the length of the sdc_list is greater than 0
	if deploymentName == "ABC" {
		return fmt.Errorf("sdc_list attribute length is not greater than 0")
	}

	return nil
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
	firmware_id = "8aaaee208c8c467e018cd37813250614"
  }
`

var ServiceResourceConfig2 = `
resource "powerflex_service" "service" {
	deployment_name = "Test-Create-Update"
	deployment_description = "Test Service-Update"
	template_id = "453c41eb-d72a-4ed1-ad16-bacdffbdd766"
	firmware_id = "8aaaee208c8c467e018cd37813250614"
	nodes = 5
}
`
