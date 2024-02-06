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
		},
	})
}

var ServiceResourceConfig1 = `
resource "powerflex_service" "service-create" {
	deployment_name = "Test-Create"
	deployment_desc = "Test Service"
	template_id = "8150d563-639d-464e-80c4-a435ed10f132"
	firmware_id = "8aaaee208c8c467e018cd37813250614"
  }
`
