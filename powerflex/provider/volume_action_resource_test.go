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

func TestAccResourceVolumeActionInvalidAction(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidActionConfig = `
	resource "powerflex_volume_action" "test" {
		volume_id = "invalid-volume-id"
		action    = "refresh"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidActionConfig,
				ExpectError: regexp.MustCompile(`.*Error Refreshing Volume.*`),
			},
		},
	})
}

func TestAccResourceVolumeActionInvalidActionType(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidActionTypeConfig = `
	resource "powerflex_volume_action" "test" {
		volume_id = "some-volume-id"
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

func TestAccResourceVolumeActionRestoreMissingSnapshot(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var missingSnapshotConfig = `
	resource "powerflex_volume_action" "test" {
		volume_id = "some-volume-id"
		action    = "restore"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + missingSnapshotConfig,
				ExpectError: regexp.MustCompile(`.*Missing Required Parameter.*`),
			},
		},
	})
}

func TestAccResourceVolumeActionMapMissingHost(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var missingHostConfig = `
	resource "powerflex_volume_action" "test" {
		volume_id = "some-volume-id"
		action    = "map_to_host"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + missingHostConfig,
				ExpectError: regexp.MustCompile(`.*Missing Required Parameter.*`),
			},
		},
	})
}
