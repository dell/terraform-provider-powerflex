package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var createStoragePool = `
	resource "powerflex_storage_pool" "pre-req1"{
		name = "terraform-storage-pool"
		protection_domain_name = "domain1"
		media_type = "HDD"
	}
`

func TestAccDeviceResourceWithSPID(t *testing.T) {
	var AddDeviceWithSPID = createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = "0db2c37100000001"
		media_type = "HDD"
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "terraform-storage-pool"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "sds_id", "0db2c37100000001"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
				),
			},
		}})
}

func TestAccDeviceResourceWithSPName(t *testing.T) {
	var AddDeviceWithSPName = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_name = "pool1"
		protection_domain_name = "domain1"
		sds_id = "0db2c37100000001"
		media_type = "HDD"
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "pool1"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "sds_id", "0db2c37100000001"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "protection_domain_name", "domain1"),
				),
			},
		}})
}

func TestAccDeviceResourceWithPDID(t *testing.T) {
	var AddDeviceWithSPName = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_name = "pool1"
		protection_domain_id = "202a046600000000"
		sds_id = "0db2c37100000001"
		media_type = "HDD"
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "pool1"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "sds_id", "0db2c37100000001"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "protection_domain_id", "202a046600000000"),
				),
			},
		}})
}

func TestAccDeviceResourceWithSDSName(t *testing.T) {
	var AddDeviceWithSDSName = createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		protection_domain_id = "202a046600000000"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSDSName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "terraform-storage-pool"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "sds_name", "SDS_2"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "protection_domain_id", "202a046600000000"),
				),
			},
		}})
}

func TestAccDeviceNegative(t *testing.T) {
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
