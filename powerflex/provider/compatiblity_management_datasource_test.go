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

func TestAccCompatibilityManagementDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Successful read
			{
				Config: ProviderConfigForTesting + compatibilityManagementData,
			},
			// Error doing the read
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetCompatibilityManagement, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + compatibilityManagementData,
				ExpectError: regexp.MustCompile(`.*Error in getting compatibility management details*.`),
			},
		},
	})
}

var compatibilityManagementData = `
data "powerflex_compatibility_management" "default" {
}
`
