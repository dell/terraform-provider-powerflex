/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var FaultSetCreate = `
resource "powerflex_fault_set" "newFs" {
	name = "fault-set-create-sds"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

var FaultSetUpdate = `
resource "powerflex_fault_set" "newFs1" {
	name = "fault-set-update-sds"
	protection_domain_id = "` + protectionDomainID1 + `"
}
`

func TestAccSDSResource(t *testing.T) {
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
			},
			{
				ip = "10.10.10.2"
				role = "sdcOnly"
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

	var updateSDSTest2 = FaultSetCreate + `
	resource "powerflex_sds" "sds" {
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
		fault_set_id = resource.powerflex_fault_set.newFs.id
	}
	`
	var updateFaultSet = FaultSetCreate + FaultSetUpdate + `
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
		fault_set_id = resource.powerflex_fault_set.newFs1.id
	}
	`

	resourceName := "powerflex_sds.sds"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
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
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update sds name
			// update sds ips from all, sdcOnly to sdsOnly, all
			// increase rmcache
			// disable rfcache
			// Enable high performance profile
			{
				Config: ProviderConfigForTesting + updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_size_in_mb", "256"),
					resource.TestCheckResourceAttr(resourceName, "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rfcache_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "performance_profile", "HighPerformance"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "sdsOnly",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_list.*", map[string]string{
						"ip":   "10.10.10.2",
						"role": "sdcOnly",
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
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
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
		},
	})
}

func TestAccSDSResourceDuplicateIP(t *testing.T) {
	createSDSTestManyValid := `
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
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.2"
					role = "sdcOnly"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
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
			// create sds test valid
			{
				Config: ProviderConfigForTesting + createSDSTestManyValid,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "3"),
				),
			},
			// modify sds test invalid
			{
				Config:      ProviderConfigForTesting + createSDSTestManyInValid,
				ExpectError: regexp.MustCompile(`.*The IP .* is configured with .*roles.*`),
			},
		},
	})
}

func TestAccSDSResourceRmCache(t *testing.T) {
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
	rcEnabled, rcDisabled, rcUnknown := "rmcache_enabled = \"true\"", "rmcache_enabled = \"false\"", ""
	rcSize, rcSizeUnknown, rcSizeIncreased := "rmcache_size_in_mb = 200", "", "rmcache_size_in_mb = 300"
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
			// Check that SDS can be created with correct rmcache settings
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcEnabled, rcSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "200"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
				),
			},
			// Check that SDS cannot be updated with wrong rmcache settings
			{
				Config:      ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcDisabled, rcSize),
				ExpectError: regexp.MustCompile(".*Read Ram cache must be enabled in order to configure its size.*"),
			},
			// Check that SDS can be updated with correct rmcache settings
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcUnknown, rcSizeIncreased),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "300"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
				),
			},
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcDisabled, rcSizeUnknown),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "false"),
				),
			},
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(sdsConfig, rcEnabled, rcSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "200"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
				),
			},
		},
	})
}

func TestAccSDSResourceCreateWithoutIP(t *testing.T) {
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

func TestAccSDSResourceCreateWithBadRole(t *testing.T) {
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

func TestAccSDSResourceCreateWithBadPerformanceProfile(t *testing.T) {
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

func TestAccSDSResourceCreateWithoutPD(t *testing.T) {
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

func TestAccSDSResourceCreateWithoutName(t *testing.T) {
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

func TestSDSResourceCreateNegative(t *testing.T) {
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
	var createInvalidCharName = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS 1"
			ip_list = [
				{
					ip = "` + SdsResourceTestData.SdsIP2 + `"
					role = "all"
				}
			]
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	var createInvalidCharName1 = `
		resource "powerflex_sds" "sds" {
			name = "Terraform_SDS_more_than_31_chars!!"
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
				Config:      ProviderConfigForTesting + createInvalidCharName,
				ExpectError: regexp.MustCompile(`.*Could not create SDS with name Terraform_SDS 1.*`),
			},
			{
				Config:      ProviderConfigForTesting + createInvalidCharName1,
				ExpectError: regexp.MustCompile(`.*Could not create SDS with name Terraform_SDS_more_than_31_chars!!.*`),
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

func TestSDSResourceCreateSpecialChar(t *testing.T) {
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

func TestSDSResourceCreateMandatoryParams(t *testing.T) {
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

func TestSDSResourceModifyRole(t *testing.T) {
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

func TestSDSResourceModifyRoleAddIP(t *testing.T) {
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

func TestSDSResourceAddIP(t *testing.T) {
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

func TestSDSResourceRemoveIP(t *testing.T) {
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

func TestSDSResourceModifyInvalid(t *testing.T) {
	var createSDSTest = `
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
		rfcache_enabled = true
		drl_mode = "NonVolatile"
		protection_domain_id = "` + protectionDomainID1 + `"
	}
	`
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
	var invalidRMCacheMaxSize = `
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
			rmcache_size_in_mb = 3912
			rfcache_enabled = true
			drl_mode = "NonVolatile"
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	var invalidRMCacheMinSize = `
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
			rmcache_size_in_mb = 127
			rfcache_enabled = true
			drl_mode = "NonVolatile"
			protection_domain_id = "` + protectionDomainID1 + `"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", protectionDomainID1),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "2"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "156"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rfcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "performance_profile", "Compact"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   SdsResourceTestData.SdsIP2,
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   "10.10.10.1",
						"role": "sdcOnly",
					}),
				),
			},
			// {
			// 	Config:      ProviderConfigForTesting + invalidRFCache,
			// 	ExpectError: regexp.MustCompile(`.*Invalid reference.*`),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + invalidRMCache,
			// 	ExpectError: regexp.MustCompile(`.*Invalid reference.*`),
			// },
			{
				Config:      ProviderConfigForTesting + rmcacheDisabled,
				ExpectError: regexp.MustCompile(`.*rmcache_size_in_mb cannot be specified while rmcache_enabled is not set to true.*`),
			},
			{
				Config:      ProviderConfigForTesting + invalidRMCacheMaxSize,
				ExpectError: regexp.MustCompile(`.*Could not change SDS Read Ram Cache size to 3912.*`),
			},
			{
				Config:      ProviderConfigForTesting + invalidRMCacheMinSize,
				ExpectError: regexp.MustCompile(`.*Could not change SDS Read Ram Cache size to 127.*`),
			},
			// {
			// 	Config:      ProviderConfigForTesting + invalidRMCacheFloatSize,
			// 	ExpectError: regexp.MustCompile(`.*Int64 Type Validation Error.*`),
			// },
		},
	})
}

func TestSDSResourceRename(t *testing.T) {
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

func TestSDSResoureUpdateProtectionDomainIDNegative(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createSDSMandatoryParams,
			},
			{
				Config:      ProviderConfigForTesting + updateProtectionDomainID,
				ExpectError: regexp.MustCompile(`.*Protection domain ID cannot be updated.*`),
			},
		},
	})
}

func TestSDSResoureUpdateProtectionDomainName(t *testing.T) {
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
