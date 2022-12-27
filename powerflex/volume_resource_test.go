package powerflex

import (
	"regexp"
	"testing"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVolumeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create volume test
			{
				Config: ProviderConfigForTesting + createVolumeTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "volume-ses-test"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "size", "8"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "map_sdcs_id.0", "c423b09800000003"),
				),
			},
			// update volume test
			{
				Config: ProviderConfigForTesting + updateVolumeTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "volume-ses"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "size", "16"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "map_sdcs_id.0", "c423b09a00000005"),
				),
			},

			// createVolumeInvalidSizeTest
			{
				Config:      ProviderConfigForTesting + createVolumeInvalidSizeTest,
				ExpectError: regexp.MustCompile(`.*Size Must be in granularity of 8GB.*`),
			},

			// createVolumeInvalidCapacityUnitTest
			{
				Config: ProviderConfigForTesting + createVolumeInvalidCapacityUnitTest,
				ExpectError:  regexp.MustCompile(`.*Attribute capacity_unit value must be one of: ["\"GB\"" "\"TB\""]*.`),
			},
		},
	})
}

var createVolumeTest = `
resource "powerflex_volume" "avengers" {
	name = "volume-ses-test"
	storage_pool_name = "pool1"
	protection_domain_id = "4eeb304600000000"
	capacity_unit = "GB"
	size = 8
	map_sdcs_id = ["c423b09800000003"]
  }
`

var updateVolumeTest = `
resource "powerflex_volume" "avengers" {
	name = "volume-ses"
	storage_pool_id = "7630a24600000000"
	protection_domain_name = "domain1"
	size = 16
	map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
  }			
`
var createVolumeInvalidSizeTest = `
resource "powerflex_volume" "avengers-invalid-size" {
	name = "volume-ses-invalid-size-test"
	storage_pool_id = "7630a24600000000"
	protection_domain_name = "domain1"
	size = 10
	map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
  }		
`

var createVolumeInvalidCapacityUnitTest = `
resource "powerflex_volume" "avengers-invalid-capacity" {
	name = "volume-ses-invalid-capacity-unit"
	storage_pool_id = "pool1"
	protection_domain_name = "domain1"
	size = 8
	capacity_unit = "HB"
	map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
  }	
`