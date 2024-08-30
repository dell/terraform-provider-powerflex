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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	SDCMappingResourceName2 = "terraform_sdc_do_not_delete"
	SDCVolName = "tf-unknown-test-donot-delete"
}

var getSDCID = `
	data "powerflex_sdc" "all" {
	}

	locals {
		matching_sdc = [for sdc in data.powerflex_sdc.all.sdcs : sdc if sdc.name == "terraform_sdc_do_not_delete"]
	}
`

var createVolRO = `
	resource "powerflex_volume" "pre-req1"{
		name = "terraform-vol"
		protection_domain_name = "domain1"
		storage_pool_name = "pool1"
		size = 8
		access_mode = "ReadOnly"
	}
`

var createVolNeg = `
	resource "powerflex_volume" "pre-req1"{
		name = "terraform-vol-neg"
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

func TestAccResourceSDCVolumes(t *testing.T) {
	var MapSDCVolumesResource = createVolRO + createVolRW + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
			volume_list = [
			{
				volume_name = resource.powerflex_volume.pre-req1.name
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

	var ChangeSDCVolumesResource = createVolRO + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "terraform_sdc_do_not_delete"),
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
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "terraform_sdc_do_not_delete"),
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
				ResourceName: "powerflex_sdc_volumes_mapping.map-sdc-volumes-test",
				ImportState:  true,
			},
			// Unmap volume from SDC
			{
				Config: ProviderConfigForTesting + MapSDCVolumesResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "terraform_sdc_do_not_delete"),
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
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "terraform_sdc_do_not_delete"),
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
	var NonExistingVolumeByID = getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var NonExistingVolumeByName = getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var InvalidLimits = createVolNeg + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
		}})
}

func TestAccSDCVolumesResourceUpdate(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an ACC test")
	}
	var CreateSDCVolumesResource = createVolRW + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var UpdateAccessMode = createVolRW + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var UpdateMapNegative = createVolRW + createVolRO + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var UpdateLimitsNegative = createVolRW + createVolRO + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
	var UpdateExistingLimitsNegative = createVolRW + createVolRO + getSDCID + `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = local.matching_sdc[0].id
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
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "terraform_sdc_do_not_delete"),
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
