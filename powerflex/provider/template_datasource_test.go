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
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var TemplateDataSourceConfig1 = `
data "powerflex_template" "example" {
}
`

var TemplateDataSourceConfig2 = `
data "powerflex_template" "example" {
	filter{
		id = ["c44cb500-020f-4562-9456-42ec1eb5f9b2"]
	}
}
`

var TemplateDataSourceConfig3 = `
data "powerflex_template" "example" {
	filter{
		id = ["c44cb500-020f-4562-9456-42ec1eb5f9b2"]
		template_name = ["block-only"]
		cluster_count = [1]
		in_configuration = false
	}
}
`

func TestAccDatasourceAcceptanceTemplate(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TemplateDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDatasourceTemplate(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TemplateDataSourceConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + TemplateDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.id", "c44cb500-020f-4562-9456-42ec1eb5f9b2"),
				),
			},
			{
				Config: ProviderConfigForTesting + TemplateDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.id", "c44cb500-020f-4562-9456-42ec1eb5f9b2"),
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.template_name", "block-only"),
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.cluster_count", "1"),
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.in_configuration", "false"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).GetAllTemplates).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + TemplateDataSourceConfig1,
				ExpectError: regexp.MustCompile(`.*Error in getting template details*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + TemplateDataSourceConfig3,
				ExpectError: regexp.MustCompile(`.*Error in filtering Template*.`),
			},
		},
	})
}
