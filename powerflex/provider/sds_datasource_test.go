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
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AT
func TestAccDatasourceAcceptanceSds(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SdsDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceSds(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// no filter
			{
				Config: ProviderConfigForTesting + SdsDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// single filter
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.0.name", "Tf_SDS_01_DV"),
				),
			},
			// multiple filter
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.name", "Tf_SDS_01_DV"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.id", "0db7306f00000003"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.maintenance_state", "NoMaintenance"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.performance_profile", "Custom"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SdsDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Error in getting system instance*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetAllSds).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SdsDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Sdses*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + SdsDataSourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error in filtering SDS*.`),
			},
		},
	})
}

var SdsDataSourceAll = `
data "powerflex_sds" "all" {
}
`

var SdsDataSourceConfig1 = `
data "powerflex_sds" "example1" {
	filter {
		name = ["Tf_SDS_01_DV"]
	}
}
`

var SdsDataSourceConfig2 = `
data "powerflex_sds" "example2" {
	filter {
		id = ["0db7306f00000003"]
		name = ["Tf_SDS_01_DV"]
		maintenance_state = ["NoMaintenance"]
		performance_profile = ["Custom"]
	}
}
`
