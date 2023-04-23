package powerflex

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVolumeResource(t *testing.T) {
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

	var modifyPlanSdcNameNegTest = `
	resource "powerflex_volume" "avengers-volume-sdc-map-name"{
		name = "avengers-volume-sdc-map-name"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		sdc_list = [
		  {
			sdc_name = "invalid-sdc-name"
		  }
		]
	}`

	var modifyPlanSdcIDNegTest = `
	resource "powerflex_volume" "avengers-volume-sdc-map-id"{
		name = "avengers-volume-sdc-map-id"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		sdc_list = [
		  {
			sdc_id = "invalid-sdc-id"
		  }
		]
	}`

	var modifyPlanVolumeInvalidMapLimit = `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create-lm"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned" 
		access_mode = "ReadOnly"
		sdc_list = [
				  {
				   sdc_name = "` + SdsResourceTestData.sdcName2 + `"
				   limit_iops = 9
				   limit_bw_in_mbps = 122
				   access_mode = "ReadWrite"
			   },
	
		]
	  }
	`

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

	var createVolumeWithSdcConfigNegTest = `
	resource "powerflex_volume" "avengers-volume----create"{
		name = "avengers-volume----create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned" 
		access_mode = "ReadOnly" # sdc_can't be mapped to volume with access mode readonly
		sdc_list = [
		  {
				   sdc_name = "` + SdsResourceTestData.sdcName + `"
				   limit_iops = 119
				   limit_bw_in_mbps = 19
				   access_mode = "ReadWrite"
			   },
		]
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
		sdc_list = [
		  {
				   sdc_name = "` + SdsResourceTestData.sdcName + `"
				   limit_iops = 119
				   limit_bw_in_mbps = 19
				   access_mode = "ReadOnly"
			   },
		]
	  }
	`

	var updateVolumePosTest = `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned" 
		access_mode = "ReadWrite"
		sdc_list = [
				  {
				   sdc_name = "` + SdsResourceTestData.sdcName + `"
				   limit_iops = 328
				   limit_bw_in_mbps = 28
				   access_mode = "ReadOnly"
			   },
			   {
				sdc_name = "` + SdsResourceTestData.sdcName2 + `"
				limit_iops = 129
				limit_bw_in_mbps = 17
				access_mode = "ReadWrite"
			   },
			   {
				sdc_id = "e3ce46c600000003"
				limit_iops = 38
				limit_bw_in_mbps = 28
				access_mode = "NoAccess"
			   }
		]
	  }
	`

	var updateVolumeUnmapPosTest = `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1" #pool1 have medium granularity
		size = 8
		use_rm_cache = true 
		volume_type = "ThickProvisioned" 
		access_mode = "ReadWrite"
		sdc_list = [
				  {
				   sdc_name = "` + SdsResourceTestData.sdcName + `"
				   limit_iops = 328
				   limit_bw_in_mbps = 28
				   access_mode = "ReadOnly"
			   },
	
			   {
				sdc_id = "e3ce46c600000003"
				limit_iops = 38
				limit_bw_in_mbps = 28
				access_mode = "NoAccess"
			   }
		]
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
				Config:      ProviderConfigForTesting + modifyPlanSdcNameNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting sdc with the name*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanSdcIDNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting sdc name from sdc id*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanVolumeInvalidMapLimit,
				ExpectError: regexp.MustCompile(`.*Error setting the limit iops*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithInvalidCompressionMethodNegTest,
				ExpectError: regexp.MustCompile(`.*Error creating volume*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithSdcConfigNegTest,
				ExpectError: regexp.MustCompile(`.*Error mapping sdc*.`),
			},
			{
				Config: ProviderConfigForTesting + createVolumePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "sdc_list.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + updateVolumePosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "sdc_list.#", "3"),
				),
			},
			{
				Config: ProviderConfigForTesting + updateVolumeUnmapPosTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "sdc_list.#", "2"),
				),
			},
			{
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

func TestAccVolumeResourceImport(t *testing.T) {
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
		},
	})
}

func TestAccVolumeResourceDuplicateSDC(t *testing.T) {
	createDuplicateInv := `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		access_mode = "ReadWrite"
		sdc_list = [
			{
				sdc_name = "` + SdsResourceTestData.sdcName2 + `"
				limit_iops = 119
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			},
			{
				sdc_name = "` + SdsResourceTestData.sdcName2 + `"
				limit_iops = 119
				limit_bw_in_mbps = 19
				access_mode = "ReadWrite"
			}
		]
	  }
	`
	createDuplicatePos := `
	resource "powerflex_volume" "avengers-volume-create"{
		name = "avengers-volume-create"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		access_mode = "ReadWrite"
		sdc_list = [
			{
				sdc_name = "` + SdsResourceTestData.sdcName + `"
				limit_iops = 119
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			},
			{
				sdc_name = "` + SdsResourceTestData.sdcName + `"
				limit_iops = 119
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + createDuplicateInv,
				ExpectError: regexp.MustCompile(`.*Duplicate SDS*.`),
			},
			{
				Config: ProviderConfigForTesting + createDuplicatePos,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "sdc_list.#", "1"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVolumeResourceUnknown(t *testing.T) {
	if SdsResourceTestData.SdcIP == "" {
		t.Fatal("POWERFLEX_SDC_IP must be set for TestAccVolumeResourceUnknown")
	}
	tfVars := fmt.Sprintf(`
	locals {
		sdc_ip = "%s"
	}
	`, SdsResourceTestData.SdcIP)
	createVolUnk := tfVars + `
	data "powerflex_sdc" "all" {
	}
	
	data "powerflex_protection_domain" "pd" {
		 name = "domain1"
	}

	provider "random" {
	}
	
	resource "random_integer" "sdc_ind" {
	  min = 0
	  max = 0
	}

	locals {
		ips = [local.sdc_ip]
		matching_sdc = [for sdc in data.powerflex_sdc.all.sdcs : sdc if sdc.sdc_ip == local.ips[random_integer.sdc_ind.result]]
	}
	
	resource "powerflex_volume" "avengers-volume-create"{
	  name = "tf-volume-create"
	  protection_domain_name = data.powerflex_protection_domain.pd.protection_domains[0].name
	  storage_pool_name = "pool1"
	  size = 8
	  access_mode = "ReadWrite"
	  sdc_list = [
		{
		  sdc_id = local.matching_sdc[0].id
		  limit_iops = 119
		  limit_bw_in_mbps = 19
		  access_mode = "ReadOnly"
		}
	  ]
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createVolUnk,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers-volume-create", "sdc_list.#", "1"),
				),
				// TODO
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
