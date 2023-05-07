package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUploadPackageResource(t *testing.T) {
	var uploadPackageTest = `
	resource "powerflex_uploadPackage" "upload-test" {
		file_path = "/home/krunal/Work/Software/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"
	 }
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForGatewayTesting + uploadPackageTest,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		}})
}

func TestAccUploadPackageNegative(t *testing.T) {
	var InvalidPath = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_id = "c98f115200000004"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidConfigWOPD = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_name = "akash-pool"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidPD = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_name = "akash-pool"
		protection_domain_name = "invalid"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSPID = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_id = "invalid"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSPName = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_name = "invalid"
		protection_domain_name = "domain1"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSDSID = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_id = "c98f115200000004"
		sds_id = "invalid"
		media_type = "HDD"
	 }
	`

	var InvalidSDSName = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdd"
		storage_pool_id = "c98f115200000004"
		sds_name = "invalid"
		media_type = "HDD"
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + InvalidPath,
				ExpectError: regexp.MustCompile("Error adding device with path"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidConfigWOPD,
				ExpectError: regexp.MustCompile("Please provide protection_domain_name or protection_domain_id with storage_pool_name"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidPD,
				ExpectError: regexp.MustCompile("Error in getting protection domain"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidSPID,
				ExpectError: regexp.MustCompile("Error in getting storage pool details"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidSPName,
				ExpectError: regexp.MustCompile("Error in getting storage pool details"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidSDSID,
				ExpectError: regexp.MustCompile("Error in getting sds details"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidSDSName,
				ExpectError: regexp.MustCompile("Error in getting sds details"),
			},
		}})
}
