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
	"os"
	"regexp"
	"testing"

	"fmt"

	. "github.com/bytedance/mockey"
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var cloneTemplateCreate = `
resource "powerflex_template_clone" "example" {
    original_template_id = "` + OriginalTemplateID + `"
	template_name = "` + TemplateName + `- COPY"
}
`

var cloneTemplateExistingTemplateName = `
resource "powerflex_template_clone" "example" {
	original_template_id = "` + OriginalTemplateID + `"
	template_name = "` + TemplateName + `"
}
`

var cloneTemplateInvalidTemplateName = `
resource "powerflex_template_clone" "example" {
    original_template_id = "` + OriginalTemplateID + `"
	template_name = ""
}
`
var cloneTemplateInvalidOriginalTemplateID = `
resource "powerflex_template_clone" "example" {
    original_template_id = "-12345"
	template_name = "` + TemplateName + `- COPY"
}
`

var cloneTemplateEmptyOriginalTemplateID = `
resource "powerflex_template_clone" "example" {
    original_template_id = ""
	template_name = "` + TemplateName + `- COPY"
}
`

// AT
func TestAccResourceAcceptanceTemplateClone(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create
			{
				Config: ProviderConfigForTesting + cloneTemplateCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccResourceTemplateClone(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + cloneTemplateInvalidOriginalTemplateID,
				ExpectError: regexp.MustCompile(`.*Template not found*.`),
			},
			{
				Config:      ProviderConfigForTesting + cloneTemplateEmptyOriginalTemplateID,
				ExpectError: regexp.MustCompile(`.*string length must be at least 1*.`),
			},
			{
				Config:      ProviderConfigForTesting + cloneTemplateInvalidTemplateName,
				ExpectError: regexp.MustCompile(`.*string length must be at least 1*.`),
			},
			{
				Config:      ProviderConfigForTesting + cloneTemplateEmptyOriginalTemplateID,
				ExpectError: regexp.MustCompile(`.*string length must be at least 1*.`),
			},
			{
				Config:      ProviderConfigForTesting + cloneTemplateInvalidTemplateName,
				ExpectError: regexp.MustCompile(`.*string length must be at least 1*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.GatewayClient).CloneTemplate).Return(fmt.Errorf("Error While Cloning Template: Template already exists please use a different name")).Build()
				},
				Config:      ProviderConfigForTesting + cloneTemplateExistingTemplateName,
				ExpectError: regexp.MustCompile(`.*Template already exists please use a different\nname*.`),
			},
			// Create Successfully
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + cloneTemplateCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_template_clone.example", "template_name", TemplateName+"- COPY"),
					resource.TestCheckResourceAttr("powerflex_template_clone.example", "original_template_id", OriginalTemplateID),
				),
			},
		},
	})
}
