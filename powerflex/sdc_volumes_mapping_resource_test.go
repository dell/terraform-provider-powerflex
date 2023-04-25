package powerflex

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	SDCMappingResourceID2 = "e3ce1fb600000001"
	SDCMappingResourceName2 = "Terraform_sdc1"
	SDCVolName = "terraform-vol"

}

var createVolRO = `
	resource "powerflex_volume" "pre-req1"{
		name = "terraform-vol"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		access_mode = "ReadOnly"
	}
`

var createVolRW = `
	resource "powerflex_volume" "pre-req2"{
		name = "terraform-vol1"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		access_mode = "ReadWrite"
	}
`

func TestAccSDCVolumesResource(t *testing.T) {
	var MapSDCVolumesResource = createVolRO + createVolRW + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	 }
	`

	var AddVolumesToSDC = createVolRO + createVolRW + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			name = "` + SDCMappingResourceName2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			},
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadWrite"
			}	
		]
	 }
	`

	var ChangeSDCVolumesResource = createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Map SDC to volume test
			{
				Config: ProviderConfigForTesting + MapSDCVolumesResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol",
						"access_mode":      "ReadOnly",
						"limit_iops":       "140",
						"limit_bw_in_mbps": "19",
					}),
				),
			},
			// Map additional volume to SDC
			{
				Config: ProviderConfigForTesting + AddVolumesToSDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol",
						"access_mode":      "ReadOnly",
						"limit_iops":       "140",
						"limit_bw_in_mbps": "19",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol1",
						"access_mode":      "ReadWrite",
						"limit_iops":       "140",
						"limit_bw_in_mbps": "19",
					}),
				),
			},
			// Import resource
			{
				ResourceName:      "powerflex_sdc_volumes_mapping.map-sdc-volumes-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Unmap volume from SDC
			{
				Config: ProviderConfigForTesting + MapSDCVolumesResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol",
						"access_mode":      "ReadOnly",
						"limit_iops":       "140",
						"limit_bw_in_mbps": "19",
					}),
				),
			},
			// Modify limits and access mode
			{
				Config: ProviderConfigForTesting + ChangeSDCVolumesResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol",
						"access_mode":      "ReadOnly",
						"limit_iops":       "120",
						"limit_bw_in_mbps": "20",
					}),
				),
			},
		},
	})
}

func TestAccSDCVolumesResourceNegative(t *testing.T) {
	var NonExistingSDCByID = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "invalid-sdc"
			volume_list = [
			{
				volume_id = "edb2a2cb00000002"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var NonExistingSDCByName = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			name = "invalid-sdc"
			volume_list = [
			{
				volume_id = "edb2a2cb00000002"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var NonExistingVolumeByID = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = "invalid-vol"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var NonExistingVolumeByName = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "invalid-vol"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var InvalidLimits = createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 10
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var IncorrectAccessMode = createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
		id = "` + SDCMappingResourceID2 + `"
		volume_list = [
		{
			volume_id = resource.powerflex_volume.pre-req1.id
			limit_iops = 120
			limit_bw_in_mbps = 20
			access_mode = "ReadWrite"
		}
	]
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + NonExistingSDCByID,
				ExpectError: regexp.MustCompile("Error getting SDC with ID"),
			},
			{
				Config:      ProviderConfigForTesting + NonExistingSDCByName,
				ExpectError: regexp.MustCompile("Error getting SDC with name"),
			},
			{
				Config:      ProviderConfigForTesting + NonExistingVolumeByID,
				ExpectError: regexp.MustCompile("Error getting volume with ID"),
			},
			{
				Config:      ProviderConfigForTesting + NonExistingVolumeByName,
				ExpectError: regexp.MustCompile("Error getting volume with name"),
			},
			{
				Config:      ProviderConfigForTesting + InvalidLimits,
				ExpectError: regexp.MustCompile("Error setting limits to sdc"),
			},
			{
				Config:      ProviderConfigForTesting + IncorrectAccessMode,
				ExpectError: regexp.MustCompile("Error mapping sdc"),
			},
		}})
}

func TestAccSDCVolumesResourceUpdate(t *testing.T) {
	var CreateSDCVolumesResource = createVolRW + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var UpdateAccessMode = createVolRW + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadWrite"
			}
		]
	 }
	`
	var UpdateMapNegative = createVolRW + createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			},
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadWrite"
			}
		]
	 }
	`
	var UpdateLimitsNegative = createVolRW + createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			},
			{
				volume_id = resource.powerflex_volume.pre-req1.id
				limit_iops = 10
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var UpdateExistingLimitsNegative = createVolRW + createVolRO + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = resource.powerflex_volume.pre-req2.id
				limit_iops = 10
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + CreateSDCVolumesResource,
			},
			{
				Config: ProviderConfigForTesting + UpdateAccessMode,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-vol1",
						"access_mode":      "ReadWrite",
						"limit_iops":       "120",
						"limit_bw_in_mbps": "20",
					}),
				),
			},
			{
				Config:      ProviderConfigForTesting + UpdateMapNegative,
				ExpectError: regexp.MustCompile("Error mapping volume to sdc"),
			},
			{
				Config:      ProviderConfigForTesting + UpdateLimitsNegative,
				ExpectError: regexp.MustCompile("Error setting limits to sdc"),
			},
			{
				Config:      ProviderConfigForTesting + UpdateExistingLimitsNegative,
				ExpectError: regexp.MustCompile("Error setting limits to sdc"),
			},
		}})
}

func TestAccSDCResourceUnknown(t *testing.T) {
	if SdsResourceTestData.SdcIP1 == "" {
		t.Fatal("POWERFLEX_SDC_IP1 must be set for TestAccSDCResourceUnknown")
	}

	if SDCVolName == "" {
		t.Fatal("POWERFLEX_SDC_VOLUMES_MAPPING_NAME must be set for TestAccSDCResourceUnknown")
	}

	tfVars := fmt.Sprintf(`
	locals {
		sdc_ip = "%s"
	}
	`, SdsResourceTestData.SdcIP1)

	tfVars1 := fmt.Sprintf(`
	locals {
		volume_name = "%s"
	}
	`, SDCVolName)

	createSDCVolMapUnk := tfVars + tfVars1 + `
	data "powerflex_sdc" "all" {
	}

	data "powerflex_volume" "all_vol" {
	}

	provider "random" {
	}
	
	resource "random_integer" "sdc_ind" {
	  min = 0
	  max = 0
	}
	locals {
		ips = [local.sdc_ip]
		vols = [local.volume_name]
		matching_sdc = [for sdc in data.powerflex_sdc.all.sdcs : sdc if sdc.sdc_ip == local.ips[random_integer.sdc_ind.result]]
		matching_vol = [for vol in data.powerflex_volume.all_vol.volumes : vol if vol.name == local.vols[random_integer.sdc_ind.result]]
	}

	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
		id = local.matching_sdc[0].id
		volume_list = [
		{
			volume_id = local.matching_vol[0].id
			limit_iops = 120
			limit_bw_in_mbps = 20
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
				Config: ProviderConfigForTesting + createSDCVolMapUnk,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_name":      "terraform-unknown-test",
						"access_mode":      "ReadOnly",
						"limit_iops":       "120",
						"limit_bw_in_mbps": "20",
					}),
				),
			},
		},
	})
}
