/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPackageResource(t *testing.T) {
	var packageTest = `
	resource "powerflex_package" "upload-test" {
		file_path = ["../resource-test/powerflex_packages/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	var packageUpdateTest = `
	resource "powerflex_package" "upload-test" {
		file_path = ["../resource-test/powerflex_packages/EMC-ScaleIO-mdm-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + packageTest,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_package.upload-test", "package_details.*", map[string]string{
					"file_name": "EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar",
				}),
			},
			//Update
			{
				Config: ProviderConfigForTesting + packageUpdateTest,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_package.upload-test", "package_details.*", map[string]string{
					"file_name": "EMC-ScaleIO-mdm-3.6-700.103.Ubuntu.22.04.x86_64.tar",
				}),
			},
		}})
}

func TestAccPackageNegative(t *testing.T) {
	var InvalidPath = `
	resource "powerflex_package" "upload-test" {
		file_path = ["/home/Software/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	var InvalidFile = `
	resource "powerflex_package" "upload-test" {
		file_path = ["../resource-test/powerflex_packages/abc.txt"]
	 }
	`

	var InvalidNameFile = `
	resource "powerflex_package" "upload-test" {
		file_path = ["../resource-test/powerflex_packages/abc.rpm"]
	 }
	`

	var EmptyList = `
	resource "powerflex_package" "upload-test" {
		file_path = []
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + InvalidPath,
				ExpectError: regexp.MustCompile(`.*Error getting with file path.*`),
			},
			{
				Config:      ProviderConfigForTesting + InvalidFile,
				ExpectError: regexp.MustCompile(`.*invalid file type, please provide valid file type.*`),
			},
			{
				Config:      ProviderConfigForTesting + EmptyList,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value.*`),
			},
			{
				Config:      ProviderConfigForTesting + InvalidNameFile,
				ExpectError: regexp.MustCompile(`.*Error getting with file path.*`),
			},
		}})
}
