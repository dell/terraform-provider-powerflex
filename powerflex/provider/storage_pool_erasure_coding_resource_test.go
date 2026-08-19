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

func TestAccResourceStoragePoolErasureCodingInvalidPool(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidPoolConfig = `
	resource "powerflex_storage_pool_erasure_coding" "test" {
		storage_pool_id       = "invalid-pool-id"
		erasure_coding_policy = "rs_2_1"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidPoolConfig,
				ExpectError: regexp.MustCompile(`.*Error Setting Erasure Coding Policy.*`),
			},
		},
	})
}

func TestAccResourceStoragePoolErasureCodingInvalidPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}

	var invalidPolicyConfig = `
	resource "powerflex_storage_pool_erasure_coding" "test" {
		storage_pool_id       = "some-pool-id"
		erasure_coding_policy = "invalid_policy"
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + invalidPolicyConfig,
				ExpectError: regexp.MustCompile(`.*value must be one of.*`),
			},
		},
	})
}
