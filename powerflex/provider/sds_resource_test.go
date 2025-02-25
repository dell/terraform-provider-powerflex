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

var protectionDomainID1 = ProtectionDomainIDSds

var FaultSetCreate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-create"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

func TestAccResourceSDSa(t *testing.T) {
	var createSDSTest = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_01"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			{
				ip = "10.10.10.1"
				role = "sdcOnly"
			}
		]
		performance_profile = "Compact"
		rmcache_enabled = true
		rmcache_size_in_mb = 156
		drl_mode = "NonVolatile"
		protection_domain_id = "` + protectionDomainID1 + `"
		fault_set_id = resource.powerflex_fault_set.newFs.id
	}
	`
	var updateSDSTest = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_02"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			}
		]
		drl_mode = "Volatile"
		performance_profile = "HighPerformance"
		rmcache_size_in_mb = 256
		rmcache_enabled = true
		rfcache_enabled = false
		protection_domain_id = "` + protectionDomainID1 + `"
		fault_set_id = resource.powerflex_fault_set.newFs.id
	}
	`

	var updateSDSTestUpdatePdError = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_02"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			}
		]
		drl_mode = "Volatile"
		performance_profile = "HighPerformance"
		rmcache_size_in_mb = 256
		rmcache_enabled = true
		rfcache_enabled = false
		protection_domain_name = "updated-pd-error"
		fault_set_id = resource.powerflex_fault_set.newFs.id
	}
	`

	var updateSDSTest2 = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_02"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			}
		]
		performance_profile = "Compact"
		rmcache_enabled = false
		rfcache_enabled = true
		protection_domain_id = "` + protectionDomainID1 + `"
		fault_set_id = resource.powerflex_fault_set.newFs.id
	}
	`
	var updateFaultSet = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
		depends_on = [powerflex_fault_set.newFs]
		name = "Tf_SDS_02"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			},
			{
				ip = "10.10.10.2"
				role = "sdcOnly"
			}
		]
		performance_profile = "Compact"
		rmcache_enabled = false
		rfcache_enabled = true
		protection_domain_id = "` + protectionDomainID1 + `"
		fault_set_id = "id-change-invalid"
	}
	`
	createSDSTestManyInvalid := `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_01"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			},
			{
				ip = "10.10.10.1"
				role = "sdcOnly"
			},
			{
				ip = "10.10.10.1"
				role = "sdsOnly"
			},
			{
				ip = "10.10.10.2"
				role = "sdcOnly"
			}
		]
		protection_domain_id = "` + protectionDomainID1 + `"
	}
	`
	sdsConfig := `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_01"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
		]
		protection_domain_name = "domain1"
		%s
		%s
	}
	`
	rcDisabled := "rmcache_enabled = \"false\""
	rcSize := "rmcache_size_in_mb = 200"
	var rmcacheDisabled = `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_01"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			{
				ip = "10.10.10.1"
				role = "sdcOnly"
			}
		]
		performance_profile = "Compact"
		rmcache_enabled = false
		rmcache_size_in_mb = 128
		rfcache_enabled = true
		drl_mode = "NonVolatile"
		protection_domain_id = "` + protectionDomainID1 + `"
	}
	`
	var updateProtectionDomainID = `
	resource "powerflex_sds" "sds" {
		name = "Terraform_SDS"
		protection_domain_id = ""
		ip_list = [
		{
			ip = "` + SdsResourceTestData.SdsIP2 + `"
			role = "all"
		}
		]
	}
`

	resourceName := "powerflex_sds.sds"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 1 Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createSDSTest,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// 2 Get ProtectionDomain Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetNewProtectionDomainEx).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createSDSTest,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			// 3 Create SDS error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).CreateSdsWithParams).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not create SDS with name*.`),
			},
			// 4 Get FindSds error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).FindSds).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createSDSTest,
				ExpectError: regexp.MustCompile(`.*Error getting SDS after creation*.`),
			},
			// 5 performance profile error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsPerformanceProfile).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not set SDS Performance Profile settings to*.`),
			},
			// 6 create sds success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + createSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_size_in_mb", "156"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "performance_profile", "Compact"),
					resource.TestCheckResourceAttrPair("powerflex_sds.sds", "fault_set_id", "powerflex_fault_set.newFs", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_list.*", map[string]string{
						"ip":   "10.10.10.1",
						"role": "sdcOnly",
					}),
				),
			},
			// 7 check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 8 SetSDSIPRole error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSDSIPRole).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Error updating IP*.`),
			},
			// 9 Get Sds after update error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetSdsByID).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not get SDS*.`),
			},
			// 10 Get Protection Domain after update error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).FindProtectionDomain).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTestUpdatePdError,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Protection domain by name*.`),
			},
			// 11 Get Protection Domain Ex after update error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetNewProtectionDomainEx).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			// 12 Set Name error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsName).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not rename SDS*.`),
			},
			// 13 SetSdsDrlMode error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsDrlMode).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not change SDS DRL Mode to*.`),
			},
			// 14 SetSdsRmCache error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsRmCache).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest2,
				ExpectError: regexp.MustCompile(`.*Could not change SDS Read Ram Cache settings to*.`),
			},
			// 15 SetSdsRmCacheSize error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsRmCacheSize).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not change SDS Read Ram Cache size to*.`),
			},
			// 16 SetSdsRfCache error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsRfCache).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not change SDS Rf Cache settings to*.`),
			},
			// SetSdsPerformanceProfile error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.ProtectionDomain).SetSdsPerformanceProfile).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateSDSTest,
				ExpectError: regexp.MustCompile(`.*Could not set SDS Performance Profile settings to*.`),
			},
			// update sds name
			// update sds ips from all, sdcOnly to sdsOnly, all
			// increase rmcache
			// disable rfcache
			// Enable high performance profile
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_size_in_mb", "256"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rfcache_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "performance_profile", "HighPerformance"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "sdsOnly",
					}),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// disable sds rmcache
			// re-enable rfcache
			// Disable high performance profile
			{
				Config: ProviderConfigForTesting + updateSDSTest2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_size_in_mb", "256"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rfcache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "performance_profile", "Compact"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateFaultSet,
				ExpectError: regexp.MustCompile(`.*Fault set ID cannot be updated.*`),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// this threefold import state check more or less validates that import functionality works

			// Update Error Tests

			// modify sds test invalid to many Ips
			{
				Config:      ProviderConfigForTesting + createSDSTestManyInvalid,
				ExpectError: regexp.MustCompile(`.*The IP .* is configured with .*roles.*`),
			},
			// Check that SDS cannot be updated with wrong rmcache settings
			{
				Config:      ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcDisabled, rcSize),
				ExpectError: regexp.MustCompile(".*Read Ram cache must be enabled in order to configure its size.*"),
			},
			// RM cache disabled, but rmcache_size_in_mb set
			{
				Config:      ProviderConfigForTesting + rmcacheDisabled,
				ExpectError: regexp.MustCompile(`.*rmcache_size_in_mb cannot be specified while rmcache_enabled is not set to true.*`),
			},
			// Try to update the domain
			{
				Config:      ProviderConfigForTesting + updateProtectionDomainID,
				ExpectError: regexp.MustCompile(`.*Protection domain ID cannot be updated.*`),
			},
		},
	})
}

func TestAccResourceSDSDuplicateIP(t *testing.T) {
	createSDSTestManyInValid := `
		resource "powerflex_sds" "sds" {
			name = "Tf_SDS_01"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "sdsOnly"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.1"
					role = "sdsOnly"
				},
				{
					ip = "10.10.10.2"
					role = "sdcOnly"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test invalid
			{
				Config:      ProviderConfigForTesting + createSDSTestManyInValid,
				ExpectError: regexp.MustCompile(`.*The IP .* is configured with .*roles.*`),
			},
		},
	})
}

func TestAccResourceSDSRmCache(t *testing.T) {
	sdsConfig := `
		resource "powerflex_sds" "sds" {
			name = "Tf_SDS_01"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_name = "domain1"
			%s
			%s
		}
		`
	rcDisabled, rcUnknown := "rmcache_enabled = \"false\"", ""
	rcSize := "rmcache_size_in_mb = 200"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check that SDS cannot be created with wrong rmcache settings
			{
				Config:      ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcDisabled, rcSize),
				ExpectError: regexp.MustCompile(".*Read Ram cache must be enabled in order to configure its size.*"),
			},
			{
				Config:      ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcUnknown, rcSize),
				ExpectError: regexp.MustCompile(".*Read Ram cache must be enabled in order to configure its size.*"),
			},
		},
	})
}

func TestAccResourceSDSCreateWithoutIP(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*ip_list.*`),
			},
		},
	})
}

func TestAccResourceSDSCreateWithBadRole(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
				{
					ip = "10.10.10.1"
					role = "invalidRole"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*role.*`),
			},
		},
	})
}

func TestAccResourceSDSCreateWithBadPerformanceProfile(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			performance_profile = "inv"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*performance_profile.*`),
			},
		},
	})
}

func TestAccResourceSDSCreateWithoutPD(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*protection_domain.*`),
			},
		},
	})
}

func TestAccResourceSDSCreateWithoutName(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*name.*`),
			},
		},
	})
}

func TestResourceSDSCreateNegative(t *testing.T) {
	var createWOName = `
		resource "powerflex_sds" "sds" {
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	var createEmptyName = `
		resource "powerflex_sds" "sds" {
			name = ""
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	var createWOPD = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
		}
		`
	var createInvalidPD = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_id = "4eeb304600000011"
		}
		`
	var createPDNameID = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
			protection_domain_name = "domain1"
		}
		`
	var createWOIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	var createWOIPRole = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
			}
			]
		}
		`
	var createWithSDSOnlyRole = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			}
			]
		}
		`
	var createWithSDCOnlyRole = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + createWOName,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
			{
				Config:      ProviderConfigForTesting + createEmptyName,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Length.*`),
			},
			{
				Config:      ProviderConfigForTesting + createWOPD,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Combination.*`),
			},
			{
				Config:      ProviderConfigForTesting + createInvalidPD,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
			},
			{
				Config:      ProviderConfigForTesting + createPDNameID,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Combination.*`),
			},
			{
				Config:      ProviderConfigForTesting + createWOIP,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
			{
				Config:      ProviderConfigForTesting + createWOIPRole,
				ExpectError: regexp.MustCompile(`.*Incorrect attribute value type.*`),
			},
			{
				Config:      ProviderConfigForTesting + createWithSDSOnlyRole,
				ExpectError: regexp.MustCompile(`.*Could not create SDS.*`),
			},
			{
				Config:      ProviderConfigForTesting + createWithSDCOnlyRole,
				ExpectError: regexp.MustCompile(`.*Could not create SDS.*`),
			},
		},
	})
}

func TestAccResourceSDSCreateSpecialChar(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var createSDSSpecialChar = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS_!@#$%^&*"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSSpecialChar,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS_!@#$%^&*"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
				),
			},
		},
	})
}

func TestAccResourceSDSCreateMandatoryParams(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var createSDSMandatoryParams = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSMandatoryParams,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
		},
	})
}

func TestAccResourceSDSModifyRole(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var addSDSIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			{
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	var modifyRolefromsdcOnlytoall = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "all"
			  }
			]
		}
		`
	var modifyRolefromsdcOnlytosdsOnly = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdsOnly"
			  }
			]
		}
		`
	var modifyRolefromalltosdsOnly = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			  }
			]
		}
		`
	var modifyRolefromsdsOnlytosdcOnly = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdcOnly"
			},
			{
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	var modifyRolefromsdsOnlytoall = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			  }
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + addSDSIP,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + modifyRolefromsdcOnlytoall,
				ExpectError: regexp.MustCompile(`.*Error updating IP.*`),
			},
			{
				Config:      ProviderConfigForTesting + modifyRolefromsdcOnlytosdsOnly,
				ExpectError: regexp.MustCompile(`.*Error updating IP.*`),
			},
			{
				Config: ProviderConfigForTesting + modifyRolefromalltosdsOnly,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "sdsOnly",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + modifyRolefromsdsOnlytosdcOnly,
				ExpectError: regexp.MustCompile(`.*Error updating IP.*`),
			},
			{
				Config: ProviderConfigForTesting + modifyRolefromsdsOnlytoall,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
		},
	})
}

func TestAccResourceSDSModifyRoleAddIP(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var addSDSSingleIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			]
		}	
		`
	var modifyRolefromalltosdcOnlySingleIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdcOnly"
			  },
			]
		}	
		`
	var modifyRolefromalltosdsOnlySingleIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			  },
			]
		}
		`
	var addSDSIPWithRoleAll = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP10 + `"
				role = "all"
			  }
			]
		}
		`
	var addSDSIPWithRolesdsOnly = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP10 + `"
				role = "sdsOnly"
			  }
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + addSDSSingleIP,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + modifyRolefromalltosdcOnlySingleIP,
				ExpectError: regexp.MustCompile(`.*Error updating IP.*`),
			},
			{
				Config:      ProviderConfigForTesting + modifyRolefromalltosdsOnlySingleIP,
				ExpectError: regexp.MustCompile(`.*Error updating IP.*`),
			},
			{
				Config:      ProviderConfigForTesting + addSDSIPWithRoleAll,
				ExpectError: regexp.MustCompile(`.*Error adding IP.*`),
			},
			{
				Config:      ProviderConfigForTesting + addSDSIPWithRolesdsOnly,
				ExpectError: regexp.MustCompile(`.*Error adding IP.*`),
			},
		},
	})
}

func TestAccResourceSDSAddIP(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var createSDSMandatoryParams = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
		`
	var addSDSIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			{
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	var addNonexistingSDSIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP1 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP3 + `"
					role = "sdcOnly"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP4 + `"
				role = "sdcOnly"
			  }
			]
		}
		`
	var addMoreThan8IP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP1 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP3 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP4 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP5 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP6 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP7 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP8 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP9 + `"
					role = "sdcOnly"
			  }
			]
		}
		`
	var addOccupiedSDSIP = `
		resource "powerflex_sds" "sds_1" {
			name = "Terraform_SDS_1"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
			  } 
			]
		}
		resource "powerflex_sds" "sds" {
			depends_on = [
				powerflex_sds.sds_1
			]
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP1 + `"
					role = "sdcOnly"
			  },
			  {
					ip = "` + SdsResourceTestData.SdsIP3 + `"
					role = "sdcOnly"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP4 + `"
				role = "sdcOnly"
			  },
			  {
				ip = "` + SdsResourceTestData.SdsIP11 + `"
				role = "sdcOnly"
			  }
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSMandatoryParams,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
			{
				Config: ProviderConfigForTesting + addSDSIP,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
			// {
			// 	Config:      ProviderConfigForTesting + addSDSWORole,
			// 	ExpectError: regexp.MustCompile("`.*Incorrect attribute value type.*`"),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + addSDSInvalidRole,
			// 	ExpectError: regexp.MustCompile("`.*Invalid Attribute Value Match.*`"),
			// },
			{
				Config: ProviderConfigForTesting + addNonexistingSDSIP,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP3,
						"role": "sdcOnly",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP4,
						"role": "sdcOnly",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + addOccupiedSDSIP,
				ExpectError: regexp.MustCompile(`.*The SDS IP address and port already in use.*`),
			},
			{
				Config:      ProviderConfigForTesting + addMoreThan8IP,
				ExpectError: regexp.MustCompile(`.*Error adding IP.*`),
			},
		},
	})
}

func TestAccResourceSDSRemoveIP(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var addSDSIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			},
			{
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	var modifyRolefromalltosdsOnly = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			},
			{
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			}
			]
		}
		`
	var removesdcOnlyIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP1 + `"
				role = "sdcOnly"
			  },
			]
		}
		`
	var removesdsOnlyIP = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "sdsOnly"
			  },
			]
		}
		`
	var removesdcOnlyIP1 = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  },
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + modifyRolefromalltosdsOnly,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "sdsOnly",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + removesdcOnlyIP,
				ExpectError: regexp.MustCompile(`.*Error removing IP.*`),
			},
			{
				Config:      ProviderConfigForTesting + removesdsOnlyIP,
				ExpectError: regexp.MustCompile(`.*Error removing IP.*`),
			},
			{
				Config: ProviderConfigForTesting + addSDSIP,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP1,
						"role": "sdcOnly",
					}),
				),
			},
			{
				Config: ProviderConfigForTesting + removesdcOnlyIP1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
		},
	})
}

func TestAccResourceSDSRename(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var createSDSMandatoryParams = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
		`
	var renameSDSExistingName = `
		resource "powerflex_sds" "sds" {
			name = "SDS_1"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  }
			]
		}
		`
	var renameSDSInvalidName = `
		resource "powerflex_sds" "sds" {
			name = "node 1"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  }
			]
		}
		`
	var renameSDS = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS_renamed"
			protection_domain_id = "` + protectionDomainID1 + `"
			ip_list = [
			  {
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			  }
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSMandatoryParams,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + renameSDSExistingName,
				ExpectError: regexp.MustCompile(`.*Could not rename SDS.*`),
			},
			{
				Config:      ProviderConfigForTesting + renameSDSInvalidName,
				ExpectError: regexp.MustCompile(`.*Could not rename SDS.*`),
			},
			{
				Config: ProviderConfigForTesting + renameSDS,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS_renamed"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
		},
	})
}

func TestAccResourceSDSUpdateProtectionDomainName(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var createProtectionDomain = `
	resource "powerflex_protection_domain" "pd" {
		name = "Terraform-AccTest-PD"
	 }
	`

	var updateProtectionDomain = `
	resource "powerflex_protection_domain" "pd" {
		name = "Terraform-AccTest-PD-Renamed"
	 }
	`

	var createSDSMandatoryParams = createProtectionDomain + `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_name = resource.powerflex_protection_domain.pd.name
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
	}
	`

	var updateProtectionDomainName = updateProtectionDomain + `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS"
			protection_domain_name = resource.powerflex_protection_domain.pd.name
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
	`

	var updateProtectionDomainNameNegative = updateProtectionDomain + `
		resource "powerflex_sds" "sds" {
			depends_on = [resource.powerflex_protection_domain.pd]
			name = "Terraform_SDS"
			protection_domain_name = "random_name"
			ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
			]
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSMandatoryParams,
			},
			{
				Config: ProviderConfigForTesting + updateProtectionDomainName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Terraform_SDS"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_name", "Terraform-AccTest-PD-Renamed"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateProtectionDomainNameNegative,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex Protection domain by name.*`),
			},
		},
	})
}
