/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AT
func TestAccDatasourceAcceptanceDevice(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + devicesData,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceDevice(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this for Unit tests")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + devicesData,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Single Filter
			{
				Config: ProviderConfigForTesting + devicesDataSingleFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.filter-single", "device_model.0.name", "/dev/sdj"),
				),
			},
			// Multi Filter
			{
				Config: ProviderConfigForTesting + devicesDataMultiFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_device.filter-multi", "device_model.0.aggregated_state", "NeverFailed"),
					resource.TestCheckResourceAttr("data.powerflex_device.filter-multi", "device_model.0.cache_look_ahead_active", "true"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllDevices, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + devicesData,
				ExpectError: regexp.MustCompile(`.*Error getting all devices in the system*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + devicesDataMultiFilter,
				ExpectError: regexp.MustCompile(`.*Error in filtering devices*.`),
			},
		},
	})
}

var devicesData = `
data "powerflex_device" "all" {
}
`

var devicesDataSingleFilter = `
data "powerflex_device" "filter-single" {
	filter {
		name = ["/dev/sdj"]
	}
}
`

var devicesDataMultiFilter = `
data "powerflex_device" "filter-multi" {
	filter {
		aggregated_state = ["NeverFailed"]
		cache_look_ahead_active = true
	}
}
`
