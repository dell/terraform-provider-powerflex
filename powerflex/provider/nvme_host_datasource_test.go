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
	"github.com/dell/goscaleio"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var nvmeHostDatasourceConfig = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	filter {
		names = ["nvme_acc_client1001"]
	}
}
`
var nvmeHostDatasourceNoFilterConfig = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	
}
`

var nvmeHostDatasourceNonExistentConfig = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	filter {
		names = ["nvme_acc_NonExistent"]
	}
}
`

func TestAccDatasourceNvmeHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + nvmeHostDatasourceNoFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + nvmeHostDatasourceNonExistentConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.#", "0"),
				),
			},
			{
				Config: ProviderConfigForTesting + nvmeHostDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.0.name", "nvme_acc_client1001"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.0.max_num_paths", "4"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.0.max_num_sys_ports", "10"),
				),
			},
		},
	})
}

func TestAccDatasourceNvmeHostNegative(t *testing.T) {
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
				Config:      ProviderConfigForTesting + nvmeHostDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex specific system*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetAllNvmeHosts).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex NVMe Hosts*.`),
			},
		},
	})
}
