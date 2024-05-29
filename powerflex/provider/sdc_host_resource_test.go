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
	"strings"
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, "55", SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, "invalid",
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			// Update mdm ip negative
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
					mdm_ips = ["%s", "10.0.6.8"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				ExpectError: regexp.MustCompile(`.*mdm_ips cannot be changed.*`),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			// Update mdm ip negative
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
					mdm_ips = ["%s", "10.0.6.8"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				ExpectError: regexp.MustCompile(`.*mdm_ips cannot be changed.*`),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				ExpectError: regexp.MustCompile(`.*package cannot be changed.*`),
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
					mdm_ips = ["%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, strings.Join(SdcHostResourceTestData.MdmIPs, `", "`)),
				ExpectError: regexp.MustCompile(`.*OS cannot be changed.*`),
			},
		},
	})
}
