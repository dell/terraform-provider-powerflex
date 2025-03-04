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
	"context"
	"fmt"
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var FunctionMockerSdcHostResource *Mocker
var FunctionMockerSdcHostResourceDelete *Mocker
var FunctionMockerSdcHostResourceSetParams *Mocker

var linuxSdcUt = fmt.Sprintf(`
resource powerflex_sdc_host sdc {
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
`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword)

var windowsSdcNoRemoteUt = `
resource powerflex_sdc_host sdc {
	ip = "1.1.1.1"
	os_family = "windows"
	remote = {
		port = "123"
		user = "user"
		private_key = "key"
	}
	name = "sdc-windows"
	package_path = "/tmp/tfaccsdc1.tar"
}
`

var windowsSdcSuccessUt = `
resource powerflex_sdc_host sdc {
	ip = "1.1.1.1"
	os_family = "windows"
	remote = {
		port = "123"
		user = "user"
		password = "pass"
	}
	name = "sdc-windows"
	package_path = "/tmp/tfaccsdc1.tar"
}
`

var windowsSdcSuccessUpdateUt = `
resource powerflex_sdc_host sdc {
	ip = "1.1.1.1"
	clusters_mdm_ips = ["1.1.1.1","1.1.1.2"]
	os_family = "windows"
	remote = {
		port = "123"
		user = "user"
		password = "pass"
	}
	name = "sdc-windows-update"
	package_path = "/tmp/tfaccsdc1.tar"
}
`

var osUpdateErrorUt = `
resource powerflex_sdc_host sdc {
	ip = "1.1.1.1"
	os_family = "linux"
	remote = {
		port = "123"
		user = "user"
		password = "pass"
	}
	name = "sdc-windows"
	package_path = "/tmp/tfaccsdc1.tar"
}
`

var valFake, _ = types.ObjectValue(
	map[string]attr.Type{
		"port":        types.StringType,
		"user":        types.StringType,
		"password":    types.StringType,
		"dir":         types.StringType,
		"host_key":    types.StringType,
		"private_key": types.StringType,
		"certificate": types.StringType,
	},
	map[string]attr.Value{
		"port":        types.StringValue("123"),
		"user":        types.StringValue("user"),
		"password":    types.StringValue("pass"),
		"dir":         types.StringNull(),
		"host_key":    types.StringNull(),
		"private_key": types.StringNull(),
		"certificate": types.StringNull(),
	},
)

var listVal, _ = types.ListValueFrom(context.TODO(), types.StringType, []string{"1.1.1.1", "1.1.1.2"})

var sdcWindowsFakeModel = models.SdcHostModel{
	Remote:             valFake,
	MdmIPs:             types.ListNull(types.StringType),
	ID:                 types.StringValue("1.1.1.1"),
	Host:               types.StringValue("1.1.1.1"),
	OS:                 types.StringValue("windows"),
	LinuxDrvCfg:        types.StringValue("/opt/emc/scaleio/sdc/bin/"),
	WindowsDrvCfg:      types.StringValue("C:\\Program Files\\EMC\\scaleio\\sdc\\bin\\"),
	UseRemotePath:      types.BoolValue(false),
	Name:               types.StringValue("sdc-windows"),
	Pkg:                types.StringValue("/tmp/tfaccsdc1.tar"),
	PerformanceProfile: types.StringValue("default"),
	MdmConnectionState: types.StringValue("connected"),
	OnVMWare:           types.BoolValue(false),
	GUID:               types.StringValue("1234"),
	IsApproved:         types.BoolValue(true),
	Esxi: types.ObjectNull(map[string]attr.Type{
		"guid":                 types.StringType,
		"verify_vib_signature": types.BoolType,
	}),
}

var sdcWindowsFakeUpdateModel = models.SdcHostModel{
	Remote:             valFake,
	MdmIPs:             listVal,
	ID:                 types.StringValue("1.1.1.1"),
	Host:               types.StringValue("1.1.1.1"),
	OS:                 types.StringValue("windows"),
	LinuxDrvCfg:        types.StringValue("/opt/emc/scaleio/sdc/bin/"),
	WindowsDrvCfg:      types.StringValue("C:\\Program Files\\EMC\\scaleio\\sdc\\bin\\"),
	UseRemotePath:      types.BoolValue(false),
	Name:               types.StringValue("sdc-windows-update"),
	Pkg:                types.StringValue("/tmp/tfaccsdc1.tar"),
	PerformanceProfile: types.StringValue("default"),
	MdmConnectionState: types.StringValue("connected"),
	OnVMWare:           types.BoolValue(false),
	GUID:               types.StringValue("1234"),
	IsApproved:         types.BoolValue(true),
	Esxi: types.ObjectNull(map[string]attr.Type{
		"guid":                 types.StringType,
		"verify_vib_signature": types.BoolType,
	}),
}

// TestAccResourceSDCUT UT tests
func TestAccResourceSDCHostUT(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 1 Create with wrong package path negative
			{
				Config:      ProviderConfigForTesting + linuxSdcUt,
				ExpectError: regexp.MustCompile(`.*Error uploading package*`),
			},
			// 2 Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMockerSdcHostResourceDelete = Mock((*helper.SdcHostResource).DeleteWindows).Return(nil).Build()
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + linuxSdcUt,
				ExpectError: regexp.MustCompile(`.*Error in getting system instance on the PowerFlex cluster*.`),
			},
			// 3 Windows Password Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + windowsSdcNoRemoteUt,
				ExpectError: regexp.MustCompile(`.*Password is required for Windows SDC*`),
			},
			// 4 Read Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}

					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + windowsSdcSuccessUt,
				ExpectError: regexp.MustCompile(`.*Error reading SDC state*`),
			},
			// 5 SetParams Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}

					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(sdcWindowsFakeModel, nil).Build()
					FunctionMockerSdcHostResourceSetParams = Mock((*helper.SdcHostResource).SetSDCParams).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + windowsSdcSuccessUt,
				ExpectError: regexp.MustCompile(`.*Error setting SDC parameters*`),
			},
			// 6 Success Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}
					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResourceSetParams = Mock((*helper.SdcHostResource).SetSDCParams).Return(nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(sdcWindowsFakeModel,
						nil).Build()
				},
				Config: ProviderConfigForTesting + windowsSdcSuccessUt,
			},
			// 7 Read Error Update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}

					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + windowsSdcSuccessUt,
				ExpectError: regexp.MustCompile(`.*Error refreshing SDC state*`),
			},
			// 8 IP Update Error Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}

					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(sdcWindowsFakeModel, nil).Build()
				},
				Config:      ProviderConfigForTesting + linuxSdcUt,
				ExpectError: regexp.MustCompile(`.*SDC IP cannot be updated through this resource*`),
			},
			// 9 Os Update Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}

					FunctionMocker = Mock((*helper.SdcHostResource).CreateWindows).Return(nil, nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(sdcWindowsFakeModel, nil).Build()
				},
				Config:      ProviderConfigForTesting + osUpdateErrorUt,
				ExpectError: regexp.MustCompile(`.*Error updating SDC*`),
			},
			// 10 Success Update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerSdcHostResource != nil {
						FunctionMockerSdcHostResource.UnPatch()
					}
					if FunctionMockerSdcHostResourceSetParams != nil {
						FunctionMockerSdcHostResourceSetParams.UnPatch()
					}
					FunctionMocker = Mock((*helper.SdcHostResource).UpdateWindowsMdms).Return(nil).Build()
					FunctionMockerSdcHostResource = Mock((*helper.SdcHostResource).ReadSDCHost).Return(sdcWindowsFakeUpdateModel,
						nil).Build()
				},
				Config: ProviderConfigForTesting + windowsSdcSuccessUpdateUt,
			},
		},
	})
}

// TestAccResourceSDCUbuntu tests the SDC Expansion Operation on Ubuntu
func TestAccResourceSDCHostUbuntu(t *testing.T) {
	t.Skip("Skipping this test case for real environment")

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
					clusters_mdm_ips = ["%s", "%s"]
				}
				`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword,
					SdcHostResourceTestData.UbuntuPkgPath, SdcHostResourceTestData.CLS1, SdcHostResourceTestData.CLS2),
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

// TestAccResourceSDCHostEsxiNeg tests the SDC Expansion Operation on Esxi Negative Validations
func TestAccResourceSDCHostsxiNeg(t *testing.T) {

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

// TestAccResourceSDCHostEsxi tests the SDC Expansion Operation on Esxi
func TestAccResourceSDCHostEsxi(t *testing.T) {
	t.Skip("Skipping this test case for real environment")

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
					clusters_mdm_ips = ["%s", "%s"]
				}
				`, SdcHostResourceTestData.EsxiIP, SdcHostResourceTestData.EsxiPort, SdcHostResourceTestData.EsxiUser, SdcHostResourceTestData.EsxiPassword,
					SdcHostResourceTestData.EsxiPkgPath, SdcHostResourceTestData.CLS1, SdcHostResourceTestData.CLS2),
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
