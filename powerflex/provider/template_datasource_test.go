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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var TemplateDataSourceConfig1 = `
data "powerflex_template" "example" {						
}
`

var TemplateDataSourceConfig2 = `
data "powerflex_template" "example" {
	template_ids = ["` + TemplateDataPoints.TemplateID + `"]					
}
`

var TemplateDataSourceConfig3 = `
data "powerflex_template" "example" {
	template_names = ["` + TemplateDataPoints.TemplateName + `"]						
}
`

var TemplateDataSourceConfig4 = `
data "powerflex_template" "example" {
	template_ids = ["invalid"]					
}
`

func TestAccTemplateDataSource(t *testing.T) {
	t.Skip("Skipping this test case")
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
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.id", TemplateDataPoints.TemplateID),
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + TemplateDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.0.template_name", TemplateDataPoints.TemplateName),
					resource.TestCheckResourceAttr("data.powerflex_template.example", "template_details.#", "1"),
				),
			},
			{
				Config:      ProviderConfigForTesting + TemplateDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting template details using id*.`),
			},
		},
	})
}
