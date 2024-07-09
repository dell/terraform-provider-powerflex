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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSDCResource tests the SDC Expansion Operation
func TestAccSDCHostResourceNegative(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config:      ProviderConfigForTesting + SDCHostConfig1,
				ExpectError: regexp.MustCompile(`.*Password is required for Windows SDC.*`),
			},
			{
				Config:      ProviderConfigForTesting + SDCHostConfig2,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
			{
				Config:      ProviderConfigForTesting + SDCHostConfig3,
				ExpectError: regexp.MustCompile(`.*SDC Host with given IP already exists.*`),
			},
			{
				Config:      ProviderConfigForTesting + SDCHostConfig4,
				ExpectError: regexp.MustCompile(`.*Error while connecting sdc remote host.*`),
			},
		},
	})
}

func TestAccSDCHostResourcePositive(t *testing.T) {
	t.Skip("Skipping this test case for real environment")
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + SDCHostConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_host.sdc_windows", "performance_profile", "Compact"),
				),
			},
			//Update
			{
				Config: ProviderConfigForTesting + SDCHostConfig6,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_host.sdc_windows", "performance_profile", "HighPerformance"),
				),
			},
		},
	})
}

func TestAccSDCHostResourceRPMPositive(t *testing.T) {
	t.Skip("Skipping this test case for real environment")
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + SDCHostConfig7,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_host.sdc_linux_rpm", "performance_profile", "Compact"),
				),
			},
			//Update
			{
				Config: ProviderConfigForTesting + SDCHostConfig8,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_host.sdc_linux_rpm", "performance_profile", "HighPerformance"),
				),
			},
		},
	})
}

var SDCHostConfig1 = `
	
resource "powerflex_sdc_host" "sdc_windows" {
	ip =  "10.10.10.10"
	remote = {
	  user     = "Username"
	}
	os_family    = "windows"
	name         = "sdc-windows-10-118"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi" 
  }
`

var SDCHostConfig2 = `
	resource "powerflex_sdc_host" "sdc_windows" {
	remote = {
	  user     = "Username"
	  password = "Password"
	}
	os_family    = "windows"
	name         = "SDC_Test"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi"
  }
`

var SDCHostConfig3 = `
resource "powerflex_sdc_host" "sdc_windows" {
	ip =  "` + SdsResourceTestData.sdcHostExistingIP + `"
	remote = {
	  user     = "Username"
	  password = "Password"
	}
	os_family    = "windows"
	name          = "SDC_Test"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi"
  }
`

var SDCHostConfig4 = `
resource "powerflex_sdc_host" "sdc_windows" {
	ip =  "10.10.10.10"
	remote = {
	  user     = "Username"
	  password = "Password"
	}
	os_family    = "windows"
	name          = "SDC_Test"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi"
  }
`

var SDCHostConfig5 = `
resource "powerflex_sdc_host" "sdc_windows" {
	ip =  "` + SdsResourceTestData.sdcHostWinIP + `"
	remote = {
	  user     = "` + SdsResourceTestData.sdcHostWinUserName + `"
	  password =  "` + SdsResourceTestData.sdcHostWinPassword + `"
	  port = 5985
	}
	os_family    = "windows"
	name          = "SDC_Test"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi"
  }
`

var SDCHostConfig6 = `
resource "powerflex_sdc_host" "sdc_windows" {
	ip =  "` + SdsResourceTestData.sdcHostWinIP + `"
	remote = {
	  user     = "` + SdsResourceTestData.sdcHostWinUserName + `"
	  password =  "` + SdsResourceTestData.sdcHostWinPassword + `"
	  port = 5985
	}
	os_family    = "windows"
	name          = "SDC_Update"
	performance_profile = "HighPerformance"
	package_path = "/root/powerflex_packages/EMC-ScaleIO-sdc-3.6-200.105.msi"
  }
`

var SDCHostConfig7 = `
resource "powerflex_sdc_host" "sdc_linux_rpm" {
	ip =  "` + SdsResourceTestData.sdcHostRPMIP + `"
	remote = {
	  user     = "` + SdsResourceTestData.sdcHostRPMUserName + `"
	  password =  "` + SdsResourceTestData.sdcHostRPMPassword + `"
	  port = 5985
	}
	os_family    = "linux"
	name          = "SDC_Test"
	performance_profile = "Compact"
	package_path = "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm"
  }
`

var SDCHostConfig8 = `
resource "powerflex_sdc_host" "sdc_linux_rpm" {
	ip =  "` + SdsResourceTestData.sdcHostRPMIP + `"
	remote = {
	  user     = "` + SdsResourceTestData.sdcHostRPMUserName + `"
	  password =  "` + SdsResourceTestData.sdcHostRPMPassword + `"
	  port = 5985
	}
	os_family    = "linux"
	name          = "SDC_Update"
	performance_profile = "HighPerformance"
	package_path = "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm"
  }
`
