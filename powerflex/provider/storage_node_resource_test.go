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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStorageNodeInvalidAction(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidConfig = `
	resource "powerflex_storage_node" "test" {
		protection_domain_id = "invalid-pd-id"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidConfig,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
		},
	})
}

func TestAccResourceStorageNodeInvalidIPRole(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidIPRoleConfig = `
	resource "powerflex_storage_node" "test" {
		protection_domain_id = "pd-123456"
		ip_list {
			ip   = "192.168.1.1"
			role = "InvalidRole"
		}
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidIPRoleConfig,
				ExpectError: regexp.MustCompile(`.*value must be one of.*`),
			},
		},
	})
}
