package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPackageResource(t *testing.T) {
	var packageTest = `
	resource "powerflex_package" "upload-test" {
		file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_Ubuntu22.04/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	var packageUpdateTest = `
	resource "powerflex_package" "upload-test" {
		file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_Ubuntu22.04/EMC-ScaleIO-mdm-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForGatewayTesting + packageTest,
				Check: resource.TestCheckTypeSetElemNestedAttrs("powerflex_package.upload-test", "package_details.*", map[string]string{
					"file_name": "EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar",
				}),
			},
			//Update
			{
				Config: ProviderConfigForGatewayTesting + packageUpdateTest,
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
		file_path = ["/root/powerflex_packages/abc.txt"]
	 }
	`

	var InvalidNameFile = `
	resource "powerflex_package" "upload-test" {
		file_path = ["/root/powerflex_packages/abc.rpm"]
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
				Config:      ProviderConfigForGatewayTesting + InvalidPath,
				ExpectError: regexp.MustCompile(`.*Error getting with file path.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + InvalidFile,
				ExpectError: regexp.MustCompile(`.*invalid file type, please provide valid file type.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + EmptyList,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + InvalidNameFile,
				ExpectError: regexp.MustCompile(`.*Error getting with file path.*`),
			},
		}})
}
