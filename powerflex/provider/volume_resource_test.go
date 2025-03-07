/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceVolume(t *testing.T) {
	var modifyPlanProtectionDomainNegTest = `
	resource "powerflex_volume" "avengers-volume-protection-domain-name"{
		name = "avengers-volume-protection-domain-name"
		protection_domain_name = "invalid-domain-name"
		storage_pool_id = "pool1"
		size = 8
	}`

	var modifyPlanProtectionDomainIDNegTest = `
	resource "powerflex_volume" "avengers-volume-protection-domain-id"{
		name = "avengers-volume-protection-domain-id"
		protection_domain_id = "invalid-domain-id"
		storage_pool_id = "pool1"
		size = 8
	}`

	var modifyPlanStoragePoolNameNegTest = `
	resource "powerflex_volume" "avengers-volume-storage-pool-name"{
		name = "avengers-volume-storage-pool-name"
		protection_domain_name = "domain1"
		storage_pool_name = "invalid-pool-name"
		size = 8
	}`

	var modifyPlanStoragePoolIDNegTest = `
	resource "powerflex_volume" "avengers-volume-storage-pool-id"{
		name = "avengers-volume-storage-pool-id"
		protection_domain_name = "domain1"
		storage_pool_id = "invalid-pool-id"
		size = 8
	}`

	var createVolumeWithInvalidCompressionMethodNegTest = `
	resource "powerflex_volume" "avengers-volume---create"{
		name = "avengers-volume---create0101010101010100101010101010101010101"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned"
	  }
	`

	var createVolumePosTest = `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned" 
		access_mode = "ReadWrite"
	  }
	`

	var modifyVolumePosTest = `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = false
		volume_type = "ThickProvisioned" 
		access_mode = "ReadOnly"
	  }
	`

	var createVolumePos01Test = `
	resource "powerflex_volume" "avengers-volume-create-01"{
		name = "avengers-volume-create-01"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned"
		access_mode = "ReadWrite"
	}
	`

	var updateVolumeRenameNegTest = `
	resource "powerflex_volume" "avengers-volume-create-01"{
		name = "avengers-volume-create-101010101010101010101"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned"
		access_mode = "ReadWrite"
	}
	`

	var updateVolumeSizeNegTest = `
	resource "powerflex_volume" "avengers-volume-create-01"{
		name = "avengers-volume-create-01"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		capacity_unit = "TB"
		use_rm_cache = true 
		volume_type = "ThickProvisioned"
		access_mode = "ReadWrite"
	}
	`

	var createVolumeCompressionMethodNegTest = `
	resource "powerflex_volume" "avengers-volume-create-compression"{
		name = "volume-create-compression"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 16
		compression_method = "None"
	}
	`

	var createVolumeTypePosTest = `
	resource "powerflex_volume" "avengers-volume-create-volume-type"{
		name = "volume-create-volume-type"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 16
		volume_type = "ThinProvisioned"
	}
	`

	var updateVolumeTypeNegTest = `
	resource "powerflex_volume" "avengers-volume-create-volume-type"{
		name = "volume-create-volume-type"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 16
		volume_type = "ThickProvisioned"
	}
	`
	var createVolumeWithInvalidSizeNegTest = `
	resource "powerflex_volume" "volume-size-invalid"{
		name = "volume-with-invalid-size"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 9
		volume_type = "ThinProvisioned"
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// modify plan test
			{
				Config:      ProviderConfigForTesting + modifyPlanProtectionDomainNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting protection domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanProtectionDomainIDNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting protection domain with id*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanStoragePoolNameNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting storage pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanStoragePoolIDNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting storage pool with id*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithInvalidCompressionMethodNegTest,
				ExpectError: regexp.MustCompile(`.*Error creating volume*.`),
			},
			// Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// Get Storage Pool Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetStoragePoolInstance).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Error getting storage pool with id*.`),
			},
			// Get Volume Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.StoragePool).GetVolume).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Error getting volume after creation*.`),
			},
			// Create Successfully
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + createVolumePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "name", "avengers-volume-create"),
				),
			},
			// Read Error after update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.Client).GetVolume).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + modifyVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Error getting volume*.`),
			},
			// Update PD Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).FindProtectionDomain).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + modifyVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Unable to read name of protection domain of ID*.`),
			},
			// Update RM Cache Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.Volume).SetVolumeUseRmCache).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + modifyVolumePosTest,
				ExpectError: regexp.MustCompile(`.*Error setting the use rm cache*.`),
			},
			// Update Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + modifyVolumePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "name", "avengers-volume-create"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "use_rm_cache", "false"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "access_mode", "ReadOnly"),
				),
			},
			// Update access mode Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.Volume).SetVolumeAccessModeLimit).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createVolumePos01Test,
				ExpectError: regexp.MustCompile(`.*Error setting access mode on volume*.`),
			},
			// Other Attirbutes update success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + createVolumePos01Test,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create-01", "size", "8"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateVolumeRenameNegTest,
				ExpectError: regexp.MustCompile(`.*Error renaming the volume*.`),
			},
			{
				Config:      ProviderConfigForTesting + updateVolumeSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Error setting the volume size*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeCompressionMethodNegTest,
				ExpectError: regexp.MustCompile(`.*error setting the compression method*.`),
			},
			{
				Config: ProviderConfigForTesting + createVolumeTypePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create-volume-type", "volume_type", "ThinProvisioned"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updateVolumeTypeNegTest,
				ExpectError: regexp.MustCompile(`.*volume type cannot be update after volume creation*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithInvalidSizeNegTest,
				ExpectError: regexp.MustCompile(`.*Size Must be in granularity of 8GB*.`),
			},
		},
	})
}

func TestAccResourceVolumeImport(t *testing.T) {
	resourceName := "powerflex_volume.tf_create"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create dummy volume
			{
				Config: ProviderConfigForTesting + `
				resource "powerflex_volume" "tf_create"{
					name = "volume-create-tf"
					protection_domain_name = "domain1"
					storage_pool_name = "pool1"
					size = 8
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "volume-create-tf"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr(resourceName, "storage_pool_name", "pool1"),
					resource.TestCheckResourceAttr(resourceName, "size", "8"),
				),
			},
			// check that import is working
			{
				ResourceName: resourceName,
				ImportState:  true,
				// TODO // ImportStateVerify: true,
			},
			// Change name and increase size of volume
			{
				Config: ProviderConfigForTesting + `
				resource "powerflex_volume" "tf_create"{
					name = "volume-update-tf"
					protection_domain_name = "domain1"
					storage_pool_name = "pool1"
					size = 16
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "volume-update-tf"),
					resource.TestCheckResourceAttr(resourceName, "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr(resourceName, "storage_pool_name", "pool1"),
					resource.TestCheckResourceAttr(resourceName, "size", "16"),
				),
			},
			// check that import is working
			{
				ResourceName: resourceName,
				ImportState:  true,
				// TODO // ImportStateVerify: true,
			},
			// Try to import with an empty id should fail
			{
				ResourceName: "powerflex_volume.tf_create",
				ImportState:  true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					return "", nil
				},
				ExpectError: regexp.MustCompile(`.*Please provide valid volume ID*.`),
			},
		},
	})
}
