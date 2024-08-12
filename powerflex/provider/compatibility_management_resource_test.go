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
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceCompatibilityManagement(t *testing.T) {

	t.Skip("Skipping this test case, only use on 4.x or greater")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error doing the create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.SetCompatibilityManagement).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + compatibilityManagementResource,
				ExpectError: regexp.MustCompile(`.*Error setting compatibility management*.`),
			},
			// Error with invalid repository path
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + compatibilityManagementResourceInvalidPath,
				ExpectError: regexp.MustCompile(`.*Could not read repository file, make sure path to gpg file is correct*.`),
			},
			// Error with invalid read
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetCompatibilityManagement).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + compatibilityManagementResource,
				ExpectError: regexp.MustCompile(`.*Error in getting compatibility management details*.`),
			},
			// Successful Set
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + compatibilityManagementResource,
			},
		}})
}

var compatibilityManagementResource = `
resource "powerflex_compatibility_management" "test" {
    repository_path = "../resource-test/gpg_compliance_management/cm-20231005-01.gpg"
}
`
var compatibilityManagementResourceInvalidPath = `
resource "powerflex_compatibility_management" "test" {
    repository_path = "../fake/path/fake.gpg"
}
`
