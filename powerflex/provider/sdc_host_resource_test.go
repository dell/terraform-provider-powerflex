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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccSDCResource tests the SDC Expansion Operation
func TestAccSDCHostResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	cfg := ProviderConfigForTesting + SDCHostConfig1
	// t.Log(cfg)

	if SdcHostResourceTestData.UbuntuIP == "localhost" {
		os.WriteFile("/tmp/tfaccsdc.tar", []byte("Dummy SDC package"), 0644)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: cfg,
				// ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			//Import
			// {
			// 	Config:        ProviderConfigForTesting + importTest,
			// 	ImportState:   true,
			// 	ImportStateId: "123",
			// 	ResourceName:  "powerflex_sdc.test",
			// 	ExpectError:   regexp.MustCompile(`.*Unable to Find SDC.*`),
			// },
			// //Create with Packages
			// {
			// 	Config: ProviderConfigForTesting + packageTest + SDCConfig2,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.2.ip", GatewayDataPoints.tbIP),
			// 	),
			// },
			// //Update
			// {
			// 	Config: ProviderConfigForTesting + packageTest + SDCConfigUpdate,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.3.ip", GatewayDataPoints.sdcServerIP),
			// 	),
			// },
			// //Rename
			// {
			// 	Config: ProviderConfigForTesting + packageTest + SDCConfigRename,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.ip", GatewayDataPoints.secondaryMDMIP),
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.name", time.Now().Weekday().String()),
			// 	),
			// },
			// //Performance Profile
			// {
			// 	Config: ProviderConfigForTesting + packageTest + SDCConfigPerProfile,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.ip", GatewayDataPoints.secondaryMDMIP),
			// 		resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.2.performance_profile", "HighPerformance"),
			// 	),
			// },
		},
	})
}

var SDCHostConfig1 = fmt.Sprintf(`
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
	package_path = "/tmp/tfaccsdc.tar"  # "EMC-ScaleIO-sdc-3.6-700.103.Ubuntu.22.04.x86_64.tar"
	mdm_ips = ["10.247.100.214", "10.247.66.67"]
}
`, SdcHostResourceTestData.UbuntuIP, SdcHostResourceTestData.UbuntuPort, SdcHostResourceTestData.UbuntuUser, SdcHostResourceTestData.UbuntuPassword)
