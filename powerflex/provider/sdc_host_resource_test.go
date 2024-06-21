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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccSDCResourceUbuntu tests the SDC Expansion Operation on Ubuntu
func TestAccSDCHostResourceUbuntu(t *testing.T) {
	t.Skip("Skipping this test case for real environment")
	os.Setenv("TF_ACC", "1")
	if SdcHostResourceTestData.UbuntuIP == "127.0.0.1" {
		err := os.WriteFile("/tmp/tfaccsdc.tar", []byte("Dummy SDC package"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create with wrong package path negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "/tmp/tfaccsdc1.tar"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword),
				ExpectError: regexp.MustCompile(`.*no such file or directory.*`),
			},
			//Create with unsupported os negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "mac"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				ExpectError: regexp.MustCompile(`.*Attribute os_family value must be one of.*`),
			},
			//Create with wrong port negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, "55", SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				ExpectError: regexp.MustCompile(`.*connection[[:space:]]refused.*`),
			},
			//Create with wrong password negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, "invalid",
					SdcHostResourceTestData.UbuntuPkgPath),
				ExpectError: regexp.MustCompile(`.*unable[[:space:]]to[[:space:]]authenticate.*`),
			},
			//Create
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},

			// Import with wrong IP
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				ImportState:   true,
				ImportStateId: "16.16.16.16",
				ResourceName:  "powerflex_sdc_host.sdc",
				// ImportStateVerifyIgnore: []string{"package_path", "remote"},
				ExpectError: regexp.MustCompile(`.*error finding SDC by IP.*`),
			},
			// Import
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				ImportState:             true,
				ImportStateId:           SdcHostResourceTestData.UbuntuIP,
				ResourceName:            "powerflex_sdc_host.sdc",
				ImportStateVerifyIgnore: []string{"package_path", "remote"},
			},
			// Update
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu2"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			// Update ip negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "10.10.10.10"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu2"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*SDC IP cannot be updated through this resource.*`),
			},
			// Update package negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-ubuntu2"
					package_path = "/dummy/tfaccsdc2.tar"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword),
				ExpectError: regexp.MustCompile(`.*package cannot be changed.*`),
			},
			// Update os negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "esxi"
					esxi = {
						guid = "esxi-guid"
					}
					name = "sdc-ubuntu2"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath),
				ExpectError: regexp.MustCompile(`.*OS cannot be changed.*`),
			},
		},
	})
}

// TestAccSDCHostResourceEsxiNeg tests the SDC Expansion Operation on Esxi Negative Validations
func TestAccSDCHostResourceEsxiNeg(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create without esxi block negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				ExpectError: regexp.MustCompile(`.*Esxi block is required for esxi SDC.*`),
			},
			//Create without guid negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				ExpectError: regexp.MustCompile(`.*attribute "guid" is required.*`),
			},
		},
	})
}

// TestAccSDCHostResourceEsxi tests the SDC Expansion Operation on Esxi
func TestAccSDCHostResourceEsxi(t *testing.T) {
	t.Skip("Skipping this test case for real environment")
	os.Setenv("TF_ACC", "1")
	if SdcHostResourceTestData.EsxiIP == "127.0.0.1" {
		err := os.WriteFile("/tmp/tfaccsdc.zip", []byte("Dummy SDC package"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	randomGUID := `
	resource "random_uuid" "sdc_guid" {
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			//Create with wrong package path negative
			{
				Config: ProviderConfigForTesting + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = "dummy"
					}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "/tmp/tfaccsdc1.zip"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword),
				ExpectError: regexp.MustCompile(`.*no such file or directory.*`),
			},
			//Create
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},

			// Import with wrong IP
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				ImportState:   true,
				ImportStateId: "16.16.16.16",
				ResourceName:  "powerflex_sdc_host.sdc",
				// ImportStateVerifyIgnore: []string{"package_path", "remote"},
				ExpectError: regexp.MustCompile(`.*error finding SDC by IP.*`),
			},
			// Import
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				ImportState:             true,
				ImportStateId:           SdcHostResourceTestData.EsxiIP,
				ResourceName:            "powerflex_sdc_host.sdc",
				ImportStateVerifyIgnore: []string{"package_path", "remote"},
			},
			// Update
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi2"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			// Update package negative
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi2"
					package_path = "/dummy/tfaccsdc2.tar"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword),
				ExpectError: regexp.MustCompile(`.*package cannot be changed.*`),
			},
			// Update guid negative
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = "invalid"
					}
					os_family = "esxi"
					name = "sdc-esxi2"
					package_path = "/dummy/tfaccsdc2.tar"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*ESXi SDC details cannot be updated.*`),
			},
			// Update vib ignore negative
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
						verify_vib_signature = false
					}
					os_family = "esxi"
					name = "sdc-esxi2"
					package_path = "/dummy/tfaccsdc2.tar"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*ESXi SDC details cannot be updated.*`),
			},
			// Update IP negative
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					ip = "10.10.10.10"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					esxi = {
						guid = random_uuid.sdc_guid.result
					}
					os_family = "esxi"
					name = "sdc-esxi2"
					package_path = "/dummy/tfaccsdc2.tar"
				}
				`, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*SDC IP cannot be updated through this resource.*`),
			},
			// Update os negative
			{
				Config: ProviderConfigForTesting + randomGUID + fmt.Sprintf(`
				resource powerflex_sdc_host sdc {
					# depends_on = [ terraform_data.ubuntu_scini ]
					ip = "%s"
					remote = {
						port = "%s"
						user = "%s"
						password = "%s"
					}
					os_family = "linux"
					name = "sdc-esxi2"
					package_path = "%s"
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath),
				ExpectError: regexp.MustCompile(`.*OS cannot be changed.*`),
			},
		},
	})
}
