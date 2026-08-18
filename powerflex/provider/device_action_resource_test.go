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

func TestAccResourceDeviceActionInvalidAction(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}

	var invalidActionConfig = `
	resource "powerflex_device_action" "test" {
		device_id = "invalid-device-id"
		action    = "activate"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidActionConfig,
				ExpectError: regexp.MustCompile(`.*Error Activating Device.*`),
			},
		},
	})
}

func TestAccResourceDeviceActionInvalidActionType(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}

	var invalidActionTypeConfig = `
	resource "powerflex_device_action" "test" {
		device_id = "some-device-id"
		action    = "invalid_action"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidActionTypeConfig,
				ExpectError: regexp.MustCompile(`.*value must be one of.*`),
			},
		},
	})
}

func TestAccResourceDeviceActionSetCapacityMissingParam(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}

	var missingCapacityConfig = `
	resource "powerflex_device_action" "test" {
		device_id = "some-device-id"
		action    = "set_capacity_limit"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + missingCapacityConfig,
				ExpectError: regexp.MustCompile(`.*Missing Required Parameter.*`),
			},
		},
	})
}
