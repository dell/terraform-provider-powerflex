/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// TestAccSDCResource tests the SDC Expansion Operation
func TestAccSDCHostResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	if SdcHostResourceTestData.UbuntuIP == "127.0.0.1" {
		os.WriteFile("/tmp/tfaccsdc.tar", []byte("Dummy SDC package"), 0644)
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
						drv_cfg_path = "/dummy/drv_cfg-3.6.500.106-esx7.x"
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
