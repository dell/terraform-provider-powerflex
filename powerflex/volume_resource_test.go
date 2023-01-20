package powerflex

import (
	"regexp"
	"testing"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

var createVolumeWithInvalidCompressionMethodNegTest = `
resource "powerflex_volume" "avengers-volume---create"{
	name = "avengers-volume---create"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1" #pool1 have medium granularity
	size = 8
	use_rm_cache = true 
	volume_type = "ThickProvisioned"
	compression_method = "Normal" 
	# this config will result in error for storage pool with medium granularity
  }
`

var createVolumeWithInvalidSdcConfigNegTest = `
resource "powerflex_volume" "avengers-volume----create"{
	name = "avengers-volume----create"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1" #pool1 have medium granularity
	size = 8
	use_rm_cache = true 
	volume_type = "ThickProvisioned"
	# compression_method = "Normal" 
	access_mode = "ReadOnly" # sdc_can't be mapped to volume with access mode readonly
	sdc_list = [
	  {
			   sdc_name = "LGLW6092"
			   limit_iops = 119
			   limit_bw_in_mbps = 19
			   access_mode = "ReadOnly"
		   },
	]
  }
`

func TestAccVolumeResource(t *testing.T) {
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
				Config:      ProviderConfigForTesting + modifyPlanProtectionDomainIdNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting protection domain with id*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanStoragePoolNameNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting storage pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanStoragePoolIdNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting storage pool with id*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanSdcNameNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting sdc with the name*.`),
			},
			{
				Config:      ProviderConfigForTesting + modifyPlanSdcIdNegTest,
				ExpectError: regexp.MustCompile(`.*Error getting sdc name from sdc id*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithInvalidCompressionMethodNegTest,
				ExpectError: regexp.MustCompile(`.*Error creating volume*.`),
			},
			{
				Config:      ProviderConfigForTesting + createVolumeWithInvalidSdcConfigNegTest,
				ExpectError: regexp.MustCompile(`.*Failure Message*.`),
			},
		},
	})
}
