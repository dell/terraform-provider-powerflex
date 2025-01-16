/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var allResourceCreds = `data "powerflex_resource_credential" "test" {}`
var idResourceCreds = `data "powerflex_resource_credential" "test" {
    filter {
        id = ["11124caa-79eb-41f5-86d1-0915a667f987"]
    }
}`

// Accptance Tests
func TestAccDatasourceAcceptanceResourceCredentials(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Accpetance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + allResourceCreds,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests

func TestAccDatasourceResourceCredentials(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + allResourceCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_resource_credential.test", "resource_credential_details.#"),
				),
			},
			{
				Config: ProviderConfigForTesting + idResourceCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_resource_credential.test", "resource_credential_details.0.id", "11124caa-79eb-41f5-86d1-0915a667f987"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + allResourceCreds,
				ExpectError: regexp.MustCompile(`.*Error in getting system instance on the PowerFlex cluster*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetResourceCredentials).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + allResourceCreds,
				ExpectError: regexp.MustCompile(`.*Error in getting resource credentials*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + idResourceCreds,
				ExpectError: regexp.MustCompile(`.*Error in filtering resource credentials*.`),
			},
		},
	})
}
