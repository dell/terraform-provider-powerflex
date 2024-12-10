/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// Accptance Tests
func TestAccDatasourceAcceptanceVolume(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving volume based on id
			{
				Config: ProviderConfigForTesting + VolumeDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccDatasourceVolume(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Get All Volumes
			{
				Config: ProviderConfigForTesting + VolumeDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Filter Single
			{
				Config: ProviderConfigForTesting + VolumeDataSourceNames,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_volume.name-filter", "volumes.0.name", "block-volume-physical-deploy"),
				),
			},
			// Filter Multiple
			{
				Config: ProviderConfigForTesting + VolumeDataSourceMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_volume.multiple-filter", "volumes.0.vtree_id", "e1e939c600000001"),
					resource.TestCheckResourceAttr("data.powerflex_volume.multiple-filter", "volumes.0.time_stamp_is_accurate", "false"),
					resource.TestCheckResourceAttr("data.powerflex_volume.multiple-filter", "volumes.1.vtree_id", "e1e939c800000003"),
					resource.TestCheckResourceAttr("data.powerflex_volume.multiple-filter", "volumes.1.time_stamp_is_accurate", "false"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetAllVolumes).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + VolumeDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Volumes*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + VolumeDataSourceNames,
				ExpectError: regexp.MustCompile(`.*Error in filtering volumes*.`),
			},
		},
	})
}

var VolumeDataSourceAll = `
data "powerflex_volume" "all" {						
}
`

var VolumeDataSourceNames = `
data "powerflex_volume" "name-filter" {	
	filter {
		name = ["block-volume-physical-deploy", "Nas_68691eb600000000_ClstVol"]
	}					
}
`

var VolumeDataSourceMultiple = `
data "powerflex_volume" "multiple-filter" {	
	filter {
		vtree_id = ["e1e939c800000003", "e1e939c600000001"]
		time_stamp_is_accurate = false
	}					
}
`
