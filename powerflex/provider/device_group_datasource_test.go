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

func TestAccDataSourceDeviceGroupAll(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var listAllConfig = `
	data "powerflex_device_group" "all" {
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + listAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_device_group.all", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceDeviceGroupById(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var byIdConfig = `
	data "powerflex_device_group" "test" {
		id = "tfacc_device_group_id"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + byIdConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device_group.test", "id", "tfacc_device_group_id"),
					resource.TestCheckResourceAttr("data.powerflex_device_group.test", "device_groups.0.id", "tfacc_device_group_id"),
					resource.TestCheckResourceAttr("data.powerflex_device_group.test", "device_groups.0.name", "terraform-test-device-group"),
				),
			},
		},
	})
}

func TestAccDataSourceDeviceGroupInvalidID(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidIDConfig = `
	data "powerflex_device_group" "test" {
		id = "invalid-id-12345"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidIDConfig,
				ExpectError: regexp.MustCompile(`.*Error Reading Device Group.*`),
			},
		},
	})
}
