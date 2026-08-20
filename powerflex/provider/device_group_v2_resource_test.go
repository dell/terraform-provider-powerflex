/*
Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// DeviceGroup resource configs for testing
var DeviceGroupResourceCreate = fmt.Sprintf(`
resource "powerflex_device_group_v2" "test" {
	name                 = "terraform-test-device-group"
	protection_domain_id = "%s"
}
`, ProtectionDomainID)

var DeviceGroupResourceCreateMock = fmt.Sprintf(`
resource "powerflex_device_group_v2" "test" {
	name                 = "terraform-test-device-group"
	protection_domain_id = "%s"
}
`, ProtectionDomainID)

var DeviceGroupResourceUpdate = fmt.Sprintf(`
resource "powerflex_device_group_v2" "test" {
	name                 = "terraform-test-device-group-updated"
	protection_domain_id = "%s"
}
`, ProtectionDomainID)

func TestAccResourceDeviceGroupCreateUpdate(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Device Group
			{
				Config: ProviderConfigForTesting + DeviceGroupResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "name", "terraform-test-device-group"),
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "protection_domain_id", ProtectionDomainID),
					resource.TestCheckResourceAttrSet("powerflex_device_group_v2.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "powerflex_device_group_v2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Name
			{
				Config: ProviderConfigForTesting + DeviceGroupResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "name", "terraform-test-device-group-updated"),
				),
			},
		},
	})
}

func TestAccResourceDeviceGroupInvalidConfig(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidConfig = `
	resource "powerflex_device_group_v2" "test" {
		name                 = "invalid-group"
		protection_domain_id = "invalid-pd-id"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidConfig,
				ExpectError: regexp.MustCompile(`.*Error.*`),
			},
		},
	})
}

func TestAccResourceDeviceGroupCreateUpdateMock(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var mockCreate = `
	resource "powerflex_device_group_v2" "test" {
		name                 = "terraform-test-device-group"
		protection_domain_id = "tfacc_protection_domain_id"
	}
	`

	var mockUpdate = `
	resource "powerflex_device_group_v2" "test" {
		name                 = "terraform-test-device-group-updated"
		protection_domain_id = "tfacc_protection_domain_id"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Device Group
			{
				Config: ProviderConfigForTesting + mockCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "name", "terraform-test-device-group"),
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "protection_domain_id", "tfacc_protection_domain_id"),
					resource.TestCheckResourceAttrSet("powerflex_device_group_v2.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "powerflex_device_group_v2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Name
			{
				Config: ProviderConfigForTesting + mockUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device_group_v2.test", "name", "terraform-test-device-group-updated"),
				),
			},
		},
	})
}
