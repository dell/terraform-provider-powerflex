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
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	"github.com/dell/goscaleio"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var nvmeTargetDatasourceConfig = `
data "powerflex_nvme_target" "nvme_target_datasource" {
	
}
`

var nvmeTargetDatasourceConfigFilter = `
data "powerflex_nvme_target" "nvme_target_datasource" {
	filter {
		name = ["mock-name-1"]
	}
}
`

func TestAccDatasourceNvmeTarget(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + nvmeTargetDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.#", "2"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.0.name", "mock-name-1"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.0.storage_port", "12200"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.0.mdm_connection_state", "Connected"),
				),
			},
			{
				Config: ProviderConfigForTesting + nvmeTargetDatasourceConfigFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.0.storage_port", "12200"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_target.nvme_target_datasource", "nvme_target_details.0.mdm_connection_state", "Connected"),
				),
			},
		},
	})
}

func TestAccDatasourceNvmeTargetNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex specific system*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.SetAttachedNvmeHostInfo).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting NVMe controller*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetDatasourceConfigFilter,
				ExpectError: regexp.MustCompile(`.*Error getting NVMe target details*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetAllSdts).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex NVMe Targets*.`),
			},
		},
	})
}