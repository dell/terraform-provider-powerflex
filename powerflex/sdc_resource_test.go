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

package powerflex

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.test", "sdc_details.*", map[string]string{
					"ip": GatewayDataPoints.tbIP,
				}),
			},
			//Update
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigUpdate,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.test", "sdc_details.*", map[string]string{
					"ip": GatewayDataPoints.sdcServerIP,
				}),
			},
			//Reaname
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigRename,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.test", "sdc_details.*", map[string]string{
					"name": time.Now().Weekday().String(),
					"ip":   GatewayDataPoints.secondaryMDMIP,
				}),
			},
			//Perormance Profile
			{
				Config: ProviderConfigForTesting + packageTest + SDCConfigPerProfile,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.test", "sdc_details.*", map[string]string{
					"ip":                  GatewayDataPoints.secondaryMDMIP,
					"performance_profile": "Compact",
				}),
			},
		},
	})
}

func TestAccSDCManagerResourceNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SDCConfigChangeName,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.name", "sdc_details.*", map[string]string{
					"name": time.Now().Weekday().String() + "1",
				}),
			},
			{
				Config: ProviderConfigForTesting + SDCConfigUpdateName,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc.name", "sdc_details.*", map[string]string{
					"name": time.Now().Weekday().String() + "2",
				}),
			},
			{
				Config:      ProviderConfigForTesting + WithoutIP,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
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
var SDCConfigChangeName = `
resource "powerflex_sdc" "name" {
	id   = "e3ce46c500000002"
  	name = "` + time.Now().Weekday().String() + `1"
}
`
var SDCConfigUpdateName = `
resource "powerflex_sdc" "name" {
	id   = "e3ce46c500000002"
  	name = "` + time.Now().Weekday().String() + `2"
}
`

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
			sdc_id = "id"
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
			performance_profile = "Compact"
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
			is_sdc = "NO"
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
			is_sdc = "NO"
			performance_profile = "HighPerformance"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
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
			is_sdc = "Yes"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
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
