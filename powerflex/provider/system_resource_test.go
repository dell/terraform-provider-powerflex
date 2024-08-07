/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

func TestAccResourceSystem(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig6,
				ExpectError: regexp.MustCompile(".*Error getting SDC with ID.*"),
			},
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig9,
				ExpectError: regexp.MustCompile(".*Error getting SDC with GUID.*"),
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "Guid"),
				),
			},
			{
				ResourceName: "powerflex_system.test",
				ImportState:  true,
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_names.#", "1"),
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_guids.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig7,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "Guid"),
				),
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "None"),
				),
			},
		},
	})
}

func TestAccResourceSystem1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SystemResourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "Guid"),
				),
			},
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig4,
				ExpectError: regexp.MustCompile(".*Error getting SDC with Name.*"),
			},
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig5,
				ExpectError: regexp.MustCompile(".*Error getting SDC with ID.*"),
			},
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig6,
				ExpectError: regexp.MustCompile(".*Error getting SDC with ID.*"),
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig8,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_ids.#", "1"),
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_guids.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SystemResourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "None"),
				),
			},
		},
	})
}

var SearchSdcBasedOnName = `
	data "powerflex_sdc" "all" {
	}

	locals {
		matching_sdc = [for sdc in data.powerflex_sdc.all.sdcs : sdc if sdc.name == "terraform_sdc_do_not_delete"]
	}

	data "powerflex_sdc" "selected" {
		id = local.matching_sdc[0].id
	}`

var SystemResourceConfig1 = SearchSdcBasedOnName + `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_approved_ips = [
		{
			id = data.powerflex_sdc.selected.id,
			ips = ["` + MDMDataPoints.standByIP2 + `"]
		}
	]
}
`

var SystemResourceConfig7 = SearchSdcBasedOnName + `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_approved_ips = [
		{
			id = data.powerflex_sdc.selected.id,
			ips = ["` + MDMDataPoints.standByIP2 + `", "10.10.10.10"]
		}
	]
}
`

var SystemResourceConfig2 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_names = ["terraform_sdc_do_not_delete"]
}
`

var SystemResourceConfig3 = `
resource "powerflex_system" "test" {
	restricted_mode = "None"
}
`

var SystemResourceConfig4 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_names = ["invalid_sdc"]
}
`

var SystemResourceConfig5 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_ids = ["invalid_sdc"]
}
`

var SystemResourceConfig6 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_approved_ips = [
		{
			id = "invalid_sdc",
			ips = ["` + MDMDataPoints.standByIP2 + `"]
		}
	]
}
`

var SystemResourceConfig8 = SearchSdcBasedOnName + `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_ids = [data.powerflex_sdc.selected.id]
}
`
var SystemResourceConfig9 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_guids = ["invalid_guid"]
}
`
