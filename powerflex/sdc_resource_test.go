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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

type resourceDataPoints struct {
	noOfSdc   string
	name      string
	newname   string
	systemid  string
	sdcguid   string
	sdcip     string
	noOflinks string
}

var sdcResourceTestData resourceDataPoints

func init() {
	sdcResourceTestData.noOfSdc = "1"
	sdcResourceTestData.noOflinks = "4"
	sdcResourceTestData.name = "terraform_sdc_create"
	sdcResourceTestData.newname = "terraform_rename"
	sdcResourceTestData.sdcguid = "C87ACC43-298B-4AD3-A95F-344FE83192C6"
	sdcResourceTestData.sdcip = os.Getenv("POWERFLEX_SDC_IP")
	sdcResourceTestData.systemid = "09a186f8167ebe0f"
}

func TestSdcResourceUpdate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "e3ce1fb500000000",
			},
			// Update testing
			{
				Config:      ProviderConfigForTesting + TestSdcResourceUpdateBlock,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Length*`),
			},
		},
	})
}

func TestSdcResourceUpdateWrongID(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "e3ce1fb5000000004343",
				ExpectError:       regexp.MustCompile(`.*Unable to Read Powerflex systems-sdcs Read*`),
			},
		},
	})
}

func TestSdcResourceUpdateSameName(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "e3ce1fb500000000",
			},
			// Update testing
			{
				Config:      ProviderConfigForTesting + TestSdcResourceUpdateBlockSameName,
				ExpectError: regexp.MustCompile(`.*Unable to Change name Powerflex sdc*`),
			},
		},
	})
}

func TestSdcResourceCreateUpdate(t *testing.T) {
	var TestSdcResourceCreateUpdateBlockS1 = `
	resource "powerflex_sdc" "sdc" {
		id = "e3ce1fb500000000"
		name = "` + SdsResourceTestData.sdcName + "-create" + `"
	  }
	  `
	var TestSdcResourceCreateUpdateBlockS2 = `
	resource "powerflex_sdc" "sdc" {
		id = "e3ce1fb500000000"
		name = "` + SdsResourceTestData.sdcName + `"
	  }
	  `
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TestSdcResourceCreateUpdateBlockS1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", SdsResourceTestData.sdcName+"-create"),
				),
			},
			// // Update testing
			{
				Config: ProviderConfigForTesting + TestSdcResourceCreateUpdateBlockS2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", SdsResourceTestData.sdcName),
				),
			},
		},
	})
}

var (
	TestSdcResourceCreateBlock = `
	resource "powerflex_sdc" "sdc" {
		id = "e3ce1fb500000000"
		name = "` + sdcResourceTestData.name + `"
	  }
	  `

	TestSdcResourceUpdateBlock = `
	  resource "powerflex_sdc" "sdc" {
		  id = "e3ce1fb500000000"
		  name = "` + sdcResourceTestData.newname + `"
		}
		`

	TestSdcResourceUpdateBlockSameName = `
	resource "powerflex_sdc" "test_import" {
		id = "e3ce1fb500000000"
		name = "terraform_sdc"
		}
		`

	TestSdcResourceUpdateImportBlock = `
	resource "powerflex_sdc" "test_import" {
		id = "e3ce1fb500000000"
		name = "terraform_sdc"
	  }
	  `
)
