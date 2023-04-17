package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSDCVolumesResource(t *testing.T) {
	var MapSDCVolumesResource = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
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

	var AddVolumesToSDC = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			name = "` + SDCMappingResourceName2 + `"
			volume_list = [
			{
				volume_id = "edb2a2cb00000002"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadOnly"
			},
			{
				volume_name = "terraform-vol1"
				limit_iops = 140
				limit_bw_in_mbps = 19
				access_mode = "ReadWrite"
			}	
		]
	 }
	`

	var ChangeSDCVolumesResource = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
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
			// Map SDC to volume test
			{
				Config: ProviderConfigForTesting + MapSDCVolumesResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "id", "e3ce1fb600000001"),
					resource.TestCheckResourceAttr("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_id":        "edb2a2cb00000002",
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
						"volume_id":        "edb2a2cb00000002",
						"volume_name":      "terraform-vol",
						"access_mode":      "ReadOnly",
						"limit_iops":       "140",
						"limit_bw_in_mbps": "19",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sdc_volumes_mapping.map-sdc-volumes-test", "volume_list.*", map[string]string{
						"volume_id":        "edb2a2ca00000003",
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
						"volume_id":        "edb2a2cb00000002",
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
						"volume_id":        "edb2a2ca00000003",
						"volume_name":      "terraform-vol1",
						"access_mode":      "ReadWrite",
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
	var InvalidLimits = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_id = "edb2a2cb00000002"
				limit_iops = 10
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var IncorrectAccessMode = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
		id = "` + SDCMappingResourceID2 + `"
		volume_list = [
		{
			volume_name = "terraform-vol"
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
	var CreateSDCVolumesResource = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var UpdateAccessMode = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadWrite"
			}
		]
	 }
	`
	var UpdateMapNegative = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			},
			{
				volume_name = "terraform-vol"
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadWrite"
			}
		]
	 }
	`
	var UpdateLimitsNegative = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
				limit_iops = 120
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			},
			{
				volume_name = "terraform-vol"
				limit_iops = 10
				limit_bw_in_mbps = 20
				access_mode = "ReadOnly"
			}
		]
	 }
	`
	var UpdateExistingLimitsNegative = `
	resource "powerflex_sdc_volumes_mapping" "map-sdc-volumes-test" {
			id = "` + SDCMappingResourceID2 + `"
			volume_list = [
			{
				volume_name = "terraform-vol1"
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
						"volume_id":        "edb2a2ca00000003",
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
