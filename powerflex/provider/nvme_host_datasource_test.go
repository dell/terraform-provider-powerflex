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
	goscaleio_types "github.com/dell/goscaleio/types/v1"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var nvmeHostDatasourceConfig = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	
}
`

var nvmeHostDatasourceConfigFilter = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	filter {
		name = ["` + NVMeHostName + `"]
	}
}
`

var nvmeHostDatasourceNonExistentConfig = `
data "powerflex_nvme_host" "nvme_host_datasource" {
	filter {
		name = ["nvme_acc_NonExistent"]
	}
}
`

func TestAccNvmeHostDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + nvmeHostDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					listCountGreaterThan("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details", 0),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetAllNvmeHosts).Return([]goscaleio_types.NvmeHost{
						{ID: "mock-id", HostType: "NVMeHost", SystemID: "mock-system-id", Name: "", Nqn: "mock-nqn", MaxNumPaths: 4, MaxNumSysPorts: 10},
					}, nil).Build()
				},
				Config: ProviderConfigForTesting + nvmeHostDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.0.name", "NVMeHost:mock-id"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeHostDatasourceNonExistentConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.#", "0"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeHostDatasourceConfigFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_nvme_host.nvme_host_datasource", "nvme_host_details.0.name", NVMeHostName),
				),
			},
		},
	})
}

func TestAccNvmeHostDatasourceNegative(t *testing.T) {
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
				Config:      ProviderConfigForTesting + nvmeHostDatasourceConfigFilter,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex specific system*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetAllNvmeHosts).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostDatasourceConfigFilter,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex NVMe Hosts*.`),
			},
		},
	})
}
