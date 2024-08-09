/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var createSDSForTest = `
	resource "powerflex_sds" "sds" {
		name = "Tf_SDS_01"
		ip_list = [
			{
				ip = "` + SdsResourceTestData.SdsIP2 + `"
				role = "all"
			}
		]
		protection_domain_id = "` + ProtectionDomainID + `"
	}
	`

var createStoragePool = `
	resource "powerflex_storage_pool" "pre-req1"{
		name = "terraform-storage-pool"
		protection_domain_name = "domain1"
		media_type = "HDD"
	}
`

var AddDeviceWithSPID = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		media_type = "HDD"
		device_capacity = 300
	 }
`

func TestAccResourceDeviceWithSPID(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_capacity_in_kb", "314572800"),
				),
			},
		}})
}

func TestAccDeviceResourceWithSPName(t *testing.T) {
	var AddDeviceWithSPName = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_name = resource.powerflex_storage_pool.pre-req1.name
		protection_domain_name = "domain1"
		sds_id = powerflex_sds.sds.id
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
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "terraform-storage-pool"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "protection_domain_name", "domain1"),
				),
			},
		}})
}

func TestAccDeviceResourceWithPDID(t *testing.T) {
	var AddDeviceWithSPName = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_name =  resource.powerflex_storage_pool.pre-req1.name
		protection_domain_id = "` + ProtectionDomainID + `"
		sds_id = powerflex_sds.sds.id
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
					resource.TestCheckResourceAttr("powerflex_device.device-test", "storage_pool_name", "terraform-storage-pool"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "protection_domain_id", ProtectionDomainID),
				),
			},
		}})
}

func TestAccDeviceResourceWithSDSName(t *testing.T) {
	var AddDeviceWithSDSName = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device"
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_name = powerflex_sds.sds.name
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
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_name", "powerflex_sds.sds", "name"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
				),
			},
		}})
}

func TestAccDeviceNegative(t *testing.T) {
	var InvalidPath = createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidConfigWOPD = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_name = "akash-pool"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidPD = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_name = "akash-pool"
		protection_domain_name = "invalid"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSPID = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_id = "invalid"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSPName = `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_name = "invalid"
		protection_domain_name = "domain1"
		sds_name = "SDS_2"
		media_type = "HDD"
	 }
	`

	var InvalidSDSID = createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = "invalid"
		media_type = "HDD"
	 }
	`

	var InvalidSDSName = createStoragePool + `
	resource "powerflex_device" "device-test" {
		name = "terraform-device-invalid"
		device_path = "/dev/sdd"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_name = "invalid"
		media_type = "HDD"
	 }
	`

	var InvalidDeviceCapacity = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		media_type = "HDD"
		device_capacity = 500
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
			{
				Config:      ProviderConfigForTesting + InvalidDeviceCapacity,
				ExpectError: regexp.MustCompile("Error updating device capacity with ID"),
			},
		}})
}

func TestAccDeviceResourceUpdate(t *testing.T) {
	var RenameDevice = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		media_type = "HDD"
		name = "terraform-device-renamed"
	 }
`

	var UpdateDeviceMediaType = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		media_type = "SSD"
		name = "terraform-device-renamed"
	 }
`

	var UpdateDeviceExternalAccelerationType = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		external_acceleration_type = "Read"
		name = "terraform-device-renamed"
	 }
`

	var UpdateDeviceCapacity = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		external_acceleration_type = "Read"
		device_capacity = 300
		name = "terraform-device-renamed"
 }
`

	var UpdateDeviceCapacityNegative = createSDSForTest + createStoragePool + `
	resource "powerflex_device" "device-test" {
		device_path = "/dev/sdc"
		storage_pool_id = resource.powerflex_storage_pool.pre-req1.id
		sds_id = powerflex_sds.sds.id
		external_acceleration_type = "Read"
		device_capacity = 500
		name = "terraform-device-renamed"
 }
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPID,
			},
			{
				Config: ProviderConfigForTesting + RenameDevice,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device-renamed"),
				),
			},
			{
				Config:      ProviderConfigForTesting + UpdateDeviceMediaType,
				ExpectError: regexp.MustCompile("The device media type is not compatible with the Storage Pool media type"),
			},
			{
				Config: ProviderConfigForTesting + UpdateDeviceExternalAccelerationType,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device-renamed"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "external_acceleration_type", "Read"),
				),
			},
			{
				Config: ProviderConfigForTesting + UpdateDeviceCapacity,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_path", "/dev/sdc"),
					resource.TestCheckResourceAttrPair("powerflex_device.device-test", "sds_id", "powerflex_sds.sds", "id"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "name", "terraform-device-renamed"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "external_acceleration_type", "Read"),
					resource.TestCheckResourceAttr("powerflex_device.device-test", "device_capacity_in_kb", "314572800"),
				),
			},
			{
				Config:      ProviderConfigForTesting + UpdateDeviceCapacityNegative,
				ExpectError: regexp.MustCompile("Error updating device capacity with ID"),
			},
		}})
}

func TestAccDeviceResourceImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + AddDeviceWithSPID,
			},
			// Import resource
			{
				ResourceName: "powerflex_device.device-test",
				ImportState:  true,
			},
		}})
}
