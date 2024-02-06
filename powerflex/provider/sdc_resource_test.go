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
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccSDCResource tests the SDC Expansion Operation
func TestAccSDCResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config:      ProviderConfigForTesting + SDCConfig1,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			//Import
			{
				Config:        ProviderConfigForTesting + importTest,
				ImportState:   true,
				ImportStateId: "123",
				ResourceName:  "powerflex_sdc.test",
				ExpectError:   regexp.MustCompile(`.*Unable to Find SDC.*`),
			},
			//Create with Packages
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.2.ip", GatewayDataPoints.tbIP),
				),
			},
			//Update
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.3.ip", GatewayDataPoints.sdcServerIP),
				),
			},
			//Rename
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigRename,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.ip", GatewayDataPoints.secondaryMDMIP),
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.name", time.Now().Weekday().String()),
				),
			},
			//Performance Profile
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigPerProfile,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.1.ip", GatewayDataPoints.secondaryMDMIP),
					resource.TestCheckResourceAttr("powerflex_sdc.test", "sdc_details.2.performance_profile", "HighPerformance"),
				),
			},
		},
	})
}

func TestAccSDCResourceMultiNameOperation(t *testing.T) {

	var sdcID string // Variable to store the SDC IP

	// Convert the current time to milliseconds since the Unix epoch
	milliseconds1 := time.Now().UnixNano() / 1e6

	time.Sleep(5 * time.Second)

	// Convert the current time to milliseconds since the Unix epoch
	milliseconds2 := time.Now().UnixNano() / 1e6

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByName,

				Check: resource.ComposeAggregateTestCheckFunc(

					// Store the SDC ID in the variable
					func(state *terraform.State) error {
						var index string

						attrs, ok := state.RootModule().Resources["data.powerflex_sdc.selected"].Primary.Attributes, true
						if !ok {
							return fmt.Errorf("data resource not found in state")
						}

						for key, value := range attrs {
							if strings.Contains(key, "ip") && value == "N/A" {
								index = strings.Split(key, ".")[1]
							}
						}

						sdcID = attrs["sdcs."+index+".id"]
						return nil
					},
				),
			},
		}})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + `resource "powerflex_sdc" "update" {
					mdm_password = "MDM_Password"
					lia_password = "LIA_Password"
					sdc_details = [
					  {
						sdc_id = "` + sdcID + `"
						name   = "sdc_` + fmt.Sprint(milliseconds1) + `"
						performance_profile   = "HighPerformance"
						is_sdc = "Yes"
					  },
					]
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.update", "sdc_details.0.name", "sdc_"+fmt.Sprint(milliseconds1)),
				),
			},
			{
				Config: ProviderConfigForTesting + `resource "powerflex_sdc" "update" {
					mdm_password = "MDM_Password"
					lia_password = "LIA_Password"
					sdc_details = [
					  {
						sdc_id = "` + sdcID + `"
						name   = "sdc_` + fmt.Sprint(milliseconds2) + `"
						performance_profile   = "HighPerformance"
						is_sdc = "Yes"
					  },
					]
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.update", "sdc_details.0.name", "sdc_"+fmt.Sprint(milliseconds2)),
				),
			},
		}})
}

func TestAccSDCResourceNegative(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + WithoutIP,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Length.*`),
			},
			{
				Config:      ProviderConfigForTesting + WithoutPrimary,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForTesting + WithoutSecondary,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForTesting + WithoutTB,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForTesting + WrongMDMCred,
				ExpectError: regexp.MustCompile(`.*Error While Validating MDM Credentials.*`),
			},
		}})
}

var importTest = `
resource "powerflex_sdc" "test"  {
	
	}
`

var packageTest = `
resource "powerflex_package" "upload-test" {
	file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm"]
	}
`
var TestSdcDataSourceByName = `data "powerflex_sdc" "selected" {
}`

var SDCConfig1 = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "NO"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "No"
	    },
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "Yes"
   		},
	]
}
`

var SDCConfigPerProfile = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "NO"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "No"
			performance_profile = "HighPerformance"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "No"
   		},
	]
}
`

var SDCConfigRename = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "No"
			name = "` + time.Now().Weekday().String() + `"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "No"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "No"
   		},
	]
}
`

var SDCConfig2 = `
resource "powerflex_sdc" "test" {

	depends_on = [
		powerflex_package.upload-test
	]

	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"

	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "No"
			performance_profile = "HighPerformance"
			name                = "sdc_expansion_test1"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "No"
	    },
	]
}
`

var SDCConfigUpdate = `
resource "powerflex_sdc" "test" {

	depends_on = [
		powerflex_package.upload-test
	]
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "NO"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "No"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "No"
			performance_profile = "HighPerformance"
   		},
	]

}
`

var WithoutPrimary = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
			},
	]
}
`

var WithoutSecondary = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
			},
	]
}
`

var WithoutTB = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
			},
	]
}
`

var WithoutIP = `
resource "powerflex_sdc" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = ""
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
			},
	]
}
`

var WrongMDMCred = `
resource "powerflex_sdc" "test" {
	mdm_password =  "ABCD"
	lia_password= "ABCD"
	sdc_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
		},
		{
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
			},
	]
}
`
