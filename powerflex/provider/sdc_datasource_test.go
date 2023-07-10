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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type sdcDataPoints struct {
	noOfSdc   string
	name      string
	systemid  string
	sdcguid   string
	sdcip     string
	noOflinks string
}

var sdcTestData sdcDataPoints

func init() {
	sdcTestData.noOfSdc = "1"
	sdcTestData.noOflinks = "4"
	sdcTestData.name = ""
	sdcTestData.sdcguid = "123"
	sdcTestData.systemid = "456"
	sdcTestData.sdcip = os.Getenv("POWERFLEX_SDC_IP1")
}

func TestSdcDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceBlockOnlyID,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_ip", sdcTestData.sdcip),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByEmptyBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.powerflex_sdc.selected", "sdcs.*", map[string]string{
						"id":   "e3d01ba200000001",
						"name": "terraform_sdc",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.powerflex_sdc.selected", "sdcs.*", map[string]string{
						"id":   "e3d01ba100000000",
						"name": "Terraform_sdc1",
					}),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByNameBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "terraform_sdc"),
				),
			},
		},
	})
}

func TestSdcDataSourceNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceByEmptyIDNeg,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceBlockBothNeg,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Combination.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceByEmptyNameBlock,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceInvalidName,
				ExpectError: regexp.MustCompile(".*Couldn't find SDC.*"),
			},
		},
	})
}

var (
	TestSdcDataSourceBlockOnlyID = `data "powerflex_sdc" "selected" {
		id = "e3d01ba100000000"
	}`

	TestSdcDataSourceByEmptyIDNeg = `data "powerflex_sdc" "selected" {
		id = ""
	}`

	TestSdcDataSourceBlockBothNeg = `data "powerflex_sdc" "selected" {
		id = "e3d01ba100000000"
		name = "Terraform_sdc1"
	}`

	TestSdcDataSourceByEmptyNameBlock = `data "powerflex_sdc" "selected" {
		name = ""
	}`

	TestSdcDataSourceByEmptyBlock = `data "powerflex_sdc" "selected" {
	}`

	TestSdcDataSourceByNameBlock = `data "powerflex_sdc" "selected" {
		name = "terraform_sdc"
	}`

	TestSdcDataSourceInvalidName = `data "powerflex_sdc" "selected" {
		name = "Terraform_sdc11"
	}`
)
