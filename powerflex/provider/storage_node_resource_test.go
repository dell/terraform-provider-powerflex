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

// StorageNode resource configs for testing
var StorageNodeResourceCreate = fmt.Sprintf(`
resource "powerflex_storage_node" "test" {
	name                = "terraform-test-storage-node"
	protection_domain_id = "%s"
	ip_list {
		ip   = "%s"
		role = "all"
	}
}
`, ProtectionDomainID, SdsResourceTestData.SdsIP1)

var StorageNodeResourceUpdate = fmt.Sprintf(`
resource "powerflex_storage_node" "test" {
	name                = "terraform-test-storage-node-updated"
	protection_domain_id = "%s"
	ip_list {
		ip   = "%s"
		role = "all"
	}
}
`, ProtectionDomainID, SdsResourceTestData.SdsIP1)

func TestAccResourceStorageNodeCreateUpdate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storage Node
			{
				Config: ProviderConfigForTesting + StorageNodeResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_node.test", "name", "terraform-test-storage-node"),
					resource.TestCheckResourceAttr("powerflex_storage_node.test", "protection_domain_id", ProtectionDomainID),
					resource.TestCheckResourceAttrSet("powerflex_storage_node.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "powerflex_storage_node.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Name
			{
				Config: ProviderConfigForTesting + StorageNodeResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_node.test", "name", "terraform-test-storage-node-updated"),
				),
			},
		},
	})
}

func TestAccResourceStorageNodeInvalidConfig(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}

	var invalidConfig = `
	resource "powerflex_storage_node" "test" {
		name                 = "invalid-node"
		protection_domain_id = "invalid-pd-id"
		ip_list {
			ip   = "10.0.0.1"
			role = "all"
		}
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidConfig,
				ExpectError: regexp.MustCompile(`.*Error Creating Storage Node.*`),
			},
		},
	})
}
