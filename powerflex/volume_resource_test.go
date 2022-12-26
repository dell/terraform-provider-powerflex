package powerflex

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVolumeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create volume test
			{
				Config: ProviderConfigForTesting + `
				resource "powerflex_volume" "avengers" {
					name = "volume-test-01"
					storage_pool_id = "7630a24600000000"
					protection_domain_id = "4eeb304600000000"
					capacity_unit = "GB"
					size = 8
					map_sdcs_id = ["c423b09800000003"]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "volume-test-01"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "size", "8"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "map_sdcs_id.0", "c423b09800000003"),
				),
			},
			// update volume tes
			{
				Config: ProviderConfigForTesting + `
				resource "powerflex_volume" "avengers" {
					name = "volume-test"
					storage_pool_id = "7630a24600000000"
					protection_domain_id = "4eeb304600000000"
					size = 16
					map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
				  }			
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "volume-test"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "size", "16"),
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "map_sdcs_id.0", "c423b09a00000005"),
				),
			},
		},
	})
}
