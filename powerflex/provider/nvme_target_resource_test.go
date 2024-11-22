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

var nvmeTargetResourceConfig = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameCreate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"	
		}
	]
	maintenance_state = "Inactive"
}
`

var nvmeTargetResourceConfigUpdate = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageAndHost"
		}
	]
	discovery_port = 8008
	nvme_port = 4421
}
`

var nvmeTargetResourceConfigUpdateEmptyName = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = ""
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [ 
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageAndHost"
		}
	]
	discovery_port = 8008
	nvme_port = 4421
}
`

var nvmeTargetResourceConfigUpdateDiffPd = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"
	protection_domain_id               = "other_pd_id"
	ip_list     = [ 
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageAndHost"
		}
	]
	discovery_port = 8008
	nvme_port = 4421
}
`

var nvmeTargetResourceConfigUpdateInvalidPd = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"	
	protection_domain_name               = "invalid_pd"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageAndHost"
		}
	]
	discovery_port = 8008
	nvme_port = 4421
}
`

var nvmeTargetResourceConfigUpdateStoragePort = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameCreate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		}
	]
	storage_port = 12201
}
`

var nvmeTargetResourceConfigUpdateNvmePort = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameCreate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		}
	]
	nvme_port = 4421
}
`

var nvmeTargetResourceConfigUpdateDiscoveryPort = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameCreate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		}
	]
	discovery_port = 8008
}
`

var nvmeTargetResourceConfigUpdateIPRole = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageOnly"
		}
	]
}
`
var nvmeTargetResourceConfigEnterMaintanance = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageOnly"
		}
	]
	maintenance_state = "Active"
}
`

var nvmeTargetResourceConfigExitMaintanance = `
resource "powerflex_nvme_target" "nvme_target_test" {
	name              = "` + NVMeTargetNameUpdate + `"
	protection_domain_id               = "` + ProtectionDomainID + `"
	ip_list     = [
		{
			ip = "` + NVMeTargetIP1 + `"
			role = "StorageAndHost"
		},
		{
			ip = "` + NVMeTargetIP2 + `"
			role = "StorageOnly"
		}
	]
	maintenance_state = "Inactive"
}
`

func TestAccNvmeTargetResource(t *testing.T) {
	resourceName := "powerflex_nvme_target.nvme_target_test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameCreate),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
				),
			},
			// import testing
			{
				ResourceName: resourceName,
				// Config: ProviderConfigForTesting + nvmeTargetResourceConfigImport,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update testing
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
				),
			},
			// update IP role testing
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfigUpdateIPRole,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.1.role", "StorageOnly"),
				),
			},
			// enter maintance testing
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfigEnterMaintanance,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "maintenance_state", "Active"),
				),
			},
			// exit maintance testing
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfigExitMaintanance,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "maintenance_state", "Inactive"),
				),
			},
			// rollback nvme target
			{
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameCreate),
				),
			},
		},
	})
}

func TestAccNvmeTargetResourceNegative(t *testing.T) {
	resourceName := "powerflex_nvme_target.nvme_target_test"
	var createFuncMocker *Mocker
	var maintenanceModeMocker *Mocker
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Negative testing for creating nvme target
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateEmptyName,
				ExpectError: regexp.MustCompile(`.*Name cannot be empty.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).CreateSdt).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not create NVMe target with name.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createFuncMocker = Mock((*goscaleio.ProtectionDomain).CreateSdt).Return(&scaleiotypes.SdtResp{ID: "1"}, nil).Build()
					FunctionMocker = Mock((*goscaleio.System).GetNvmeHostByID).Return(nil, fmt.Errorf("mock error")).Build()
					maintenanceModeMocker = Mock(helper.ToggleSdtMaintenanceMode).Return(nil).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe target with ID.*`),
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
					if maintenanceModeMocker != nil {
						maintenanceModeMocker.UnPatch()
					}

				},
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
			},
			// Negative testing for updating nvme target
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateEmptyName,
				ExpectError: regexp.MustCompile(`.*Name cannot be empty.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateDiffPd,
				ExpectError: regexp.MustCompile(`.*Protection domain ID cannot be updated.*`),
			},
			// failed to get protection domain
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).FindProtectionDomain).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateInvalidPd,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Protection domain by name.*`),
			},
			// update with invalid protection domain
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).FindProtectionDomain).Return(&scaleiotypes.ProtectionDomain{ID: "abcdefg", Name: "invlaid_pd"}, nil).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateInvalidPd,
				ExpectError: regexp.MustCompile(`.*Protection domain name does not match the original Protection domain.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetSdtDiscoveryPort).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateDiscoveryPort,
				ExpectError: regexp.MustCompile(`.*could not update discovery port.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetSdtNvmePort).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateNvmePort,
				ExpectError: regexp.MustCompile(`.*could not update NVMe port.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).RenameSdt).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*could not rename the NVMe target.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).AddSdtTargetIP).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*could not add target IP.*`),
			},
			// update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if createFuncMocker != nil {
						createFuncMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeTargetResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameUpdate),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).ModifySdtIPRole).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateIPRole,
				ExpectError: regexp.MustCompile(`.*could not update the role of IP.*`),
			},
			// remove target IP
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).RemoveSdtTargetIP).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*could not remove target IP.*`),
			},
			// rollback nvme target
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameCreate),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SetSdtStoragePort).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdateStoragePort,
				ExpectError: regexp.MustCompile(`.*could not update storage port.*`),
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

func TestAccNvmeTargetResourceHelperNegative(t *testing.T) {
	resourceName := "powerflex_nvme_target.nvme_target_test"
	var createFuncMocker *Mocker
	var getFuncMocker *Mocker
	var maintenanceModeMocker *Mocker
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
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not read NVMe target.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error in getting system instance.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetNewProtectionDomainEx).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createFuncMocker = Mock((*goscaleio.ProtectionDomain).CreateSdt).Return(&scaleiotypes.SdtResp{ID: "1"}, nil).Build()
					getFuncMocker = Mock((*goscaleio.System).GetSdtByID).Return(&scaleiotypes.Sdt{
						ID:   "1",
						Name: NVMeTargetNameCreate,
					}, nil).Build()
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
					maintenanceModeMocker = Mock(helper.ToggleSdtMaintenanceMode).Return(nil).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not parse NVMe target struct.*`),
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
					if maintenanceModeMocker != nil {
						maintenanceModeMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ToggleSdtMaintenanceMode).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigEnterMaintanance,
				ExpectError: regexp.MustCompile(`.*Could not set maintenance mode.*`),
			},
			// error during update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + nvmeTargetResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Could not parse NVMe target struct.*`),
			},
			// rollback nvme target
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + nvmeTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", NVMeTargetNameCreate),
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
