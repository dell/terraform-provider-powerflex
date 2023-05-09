package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUploadPackageResource(t *testing.T) {
	var uploadPackageTest = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = ["/home/krunal/Work/Software/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	var uploadPackageUpdateTest = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = ["/home/krunal/Work/Software/PowerFlex_3.6.700.103_Ubuntu22.04/EMC-ScaleIO-mdm-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	// /root/powerflex_packages/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForGatewayTesting + uploadPackageTest,
				Check:  resource.TestCheckResourceAttr("powerflex_uploadPackage.upload-test", "package_details.0.file_name", "EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"),
			},
			//Update
			{
				Config: ProviderConfigForGatewayTesting + uploadPackageUpdateTest,
				Check:  resource.TestCheckResourceAttr("powerflex_uploadPackage.upload-test", "package_details.0.file_name", "EMC-ScaleIO-mdm-3.6-700.103.Ubuntu.22.04.x86_64.tar"),
			},
		}})
}

func TestAccUploadPackageNegative(t *testing.T) {
	var InvalidPath = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = ["/home/Software/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
	 }
	`

	var InvalidFile = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = ["/home/krunal/Work/Software/abc.txt"]
	 }
	`

	var InvalidNameFile = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = ["/home/krunal/Work/Software/abc.rpm"]
	 }
	`

	var EmptyList = `
	resource "powerflex_uploadPackage" "upload-test" {
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
