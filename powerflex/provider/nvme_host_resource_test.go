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

	goscaleio "github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var nvmeHostResourceConfig = `
resource "powerflex_nvme_host" "nvme_host_test" {
		name              = "nvme_acc_client1002"
		nqn               = "nqn.2014-08.org.nvmexpress:uuid:a10e4d56-a2c0-4cab-9a0a-9a7a4ebb8c0e"
		max_num_paths     = 4
		max_num_sys_ports = 10
}
`

var nvmeHostResourceConfigUpdate = `
resource "powerflex_nvme_host" "nvme_host_test" {
		name              = "nvme_acc_client1002_update"
		nqn               = "nqn.2014-08.org.nvmexpress:uuid:a10e4d56-a2c0-4cab-9a0a-9a7a4ebb8c0e"
		max_num_paths     = 8
		max_num_sys_ports = 8
}
`

func TestAccNvmeHostResource(t *testing.T) {
	resourceName := "powerflex_nvme_host.nvme_host_test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "nvme_acc_client1002"),
					resource.TestCheckResourceAttr(resourceName, "nqn", "nqn.2014-08.org.nvmexpress:uuid:a10e4d56-a2c0-4cab-9a0a-9a7a4ebb8c0e"),
					resource.TestCheckResourceAttr(resourceName, "max_num_paths", "4"),
					resource.TestCheckResourceAttr(resourceName, "max_num_sys_ports", "10"),
				),
			},
			// import testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update testing
			{
				Config: ProviderConfigForTesting + nvmeHostResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "nvme_acc_client1002_update"),
					resource.TestCheckResourceAttr(resourceName, "max_num_paths", "8"),
					resource.TestCheckResourceAttr(resourceName, "max_num_sys_ports", "8"),
				),
			},
			// rollback nvme host
			{
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "nvme_acc_client1002"),
					resource.TestCheckResourceAttr(resourceName, "max_num_paths", "4"),
					resource.TestCheckResourceAttr(resourceName, "max_num_sys_ports", "10"),
				),
			},
		},
	})
}

func TestAccNvmeHostResourceNegative(t *testing.T) {
	resourceName := "powerflex_nvme_host.nvme_host_test"
	var createFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Negative testing for creating nvme host
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).CreateNvmeHost).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not create NVMe host with name.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createFuncMocker = Mock((*goscaleio.System).CreateNvmeHost).Return(&scaleiotypes.NvmeHostResp{ID: "1"}, nil).Build()
					FunctionMocker = Mock((*goscaleio.System).GetNvmeHostByID).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe host with ID.*`),
			},
			// create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if createFuncMocker != nil {
						createFuncMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
			},
			// Negative testing for updating nvme host
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).ChangeNvmeHostName).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Could not update NVMe host name.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).ChangeNvmeHostMaxNumPaths).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Could not update max_num_paths.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).ChangeNvmeHostMaxNumSysPorts).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Could not update max_num_sys_ports.*`),
			},
			// rollback nvme host
			{
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "nvme_acc_client1002"),
					resource.TestCheckResourceAttr(resourceName, "max_num_paths", "4"),
					resource.TestCheckResourceAttr(resourceName, "max_num_sys_ports", "10"),
				),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if createFuncMocker != nil {
				createFuncMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccNvmeHostResourceHelperNegative(t *testing.T) {
	resourceName := "powerflex_nvme_host.nvme_host_test"
	var createFuncMocker *Mocker
	var getFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe host param.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createFuncMocker = Mock((*goscaleio.System).CreateNvmeHost).Return(&scaleiotypes.NvmeHostResp{ID: "1"}, nil).Build()
					getFuncMocker = Mock((*goscaleio.System).GetNvmeHostByID).Return(&scaleiotypes.NvmeHost{
						ID:   "1",
						Name: "host",
					}, nil).Build()
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe host.*`),
			},
			// create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if createFuncMocker != nil {
						createFuncMocker.UnPatch()
					}
					if getFuncMocker != nil {
						getFuncMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
			},
			// error during update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeHostResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe host struct.*`),
			},
			// rollback nvme host
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeHostResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "nvme_acc_client1002"),
					resource.TestCheckResourceAttr(resourceName, "max_num_paths", "4"),
					resource.TestCheckResourceAttr(resourceName, "max_num_sys_ports", "10"),
				),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if createFuncMocker != nil {
				createFuncMocker.UnPatch()
			}
			if getFuncMocker != nil {
				getFuncMocker.UnPatch()
			}
			return nil
		},
	})
}