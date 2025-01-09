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
	"fmt"
	"os"
	"regexp"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AT
func TestAccResourceSystem(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
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

// UT
func TestAccResourceSystemUT(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 1 set mode Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetRestrictedMode).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error changing restricted mode*.`),
			},
			// 2 get instance Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.Client).GetInstance).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error in getting system instance on the PowerFlex cluster*.`),
			},
			// 3 Error getting SDC with ID
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).FindSdc).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfigInvalidApprovedIps,
				ExpectError: regexp.MustCompile(`.*Error getting SDC with Name*.`),
			},
			// 4 Error getting SDC with Name
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetApprovedIps).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfigInvalidApprovedIps,
				ExpectError: regexp.MustCompile(`.*Error getting SDC with ID*.`),
			},
			// 5 Create Successfully
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + SystemResourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "restricted_mode", "Guid"),
				),
			},
			// 6 Import resource
			{
				ResourceName: "powerflex_system.test",
				ImportState:  true,
			},
			// 7 get instance Error
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig4,
				ExpectError: regexp.MustCompile(".*Error getting SDC with Name.*"),
			},
			// 8 get instance Error
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig5,
				ExpectError: regexp.MustCompile(".*Error getting SDC with ID.*"),
			},
			// 9 get instance Error
			{
				Config:      ProviderConfigForTesting + SystemResourceConfig6,
				ExpectError: regexp.MustCompile(".*Error getting SDC with ID.*"),
			},
			// 10 get instance after update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.Client).GetInstance).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfig8,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex specific system*.`),
			},
			// 11 set mode update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetRestrictedMode).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SystemResourceConfig3,
				ExpectError: regexp.MustCompile(`.*Error changing restricted mode*.`),
			},
			// 12 update successfully
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + SystemResourceConfig8,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_ids.#", "1"),
					resource.TestCheckResourceAttr("powerflex_system.test", "sdc_guids.#", "1"),
				),
			},
			// 13 update successfully
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
	data "powerflex_sdc" "selected" {
		filter {
			name = ["Terraform_sdc1"]
		}
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

var SystemResourceConfigInvalidApprovedIps = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_names = ["terraform_sdc_do_not_delete"]
	sdc_approved_ips = [{
		id = "invalid_sdc",
		ips = ["1.1.1.1"]
	}]
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
	sdc_ids = [data.powerflex_sdc.selected.sdcs[0].id]
}
`
var SystemResourceConfig9 = `
resource "powerflex_system" "test" {
	restricted_mode = "Guid"
	sdc_guids = ["invalid_guid"]
}
`
