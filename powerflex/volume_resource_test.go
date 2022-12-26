package powerflex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Powerflex client is properly configured.
	providerConfig = `
		provider "powerflex" {
		  username = ""
		  password = ""
		  host     = ""
		}
		`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerflex": providerserver.NewProtocol6WithError(New()),
	}
)

func TestAccVolumeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create volume test
			{
				Config: providerConfig + `
				resource "powerflex_volume" "avengers" {
					name = "avengers-ironman-3000"
					storage_pool_id = "7630a24600000000"
					protection_domain_id = "4eeb304600000000"
					capacity_unit = "GB"
					size = 10
					map_sdcs_id = ["c423b09800000003"]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "avengers-ironman-3000"),
					resource.TestCheckResourceAttr("powerflex.volume.avengers", "size", "10"),
					resource.TestCheckResourceAttr("powerflex.volume.avengers", "map_sdcs_id.0", "c423b09800000003"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "powerflex_volume.avengers",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update volume tes
			{
				Config: providerConfig + `
				resource "powerflex_volume" "avengers" {
					name = "avengers-ironman"
					storage_pool_id = "7630a24600000000"
					protection_domain_id = "4eeb304600000000"
					capacity_unit = "GB"
					size = 18
					map_sdcs_id = ["c423b09a00000005","c423b09900000004"]
				  }			
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_volume.avengers", "name", "avengers-ironman"),
					resource.TestCheckResourceAttr("powerflex.volume.avengers", "size", "18"),
					resource.TestCheckResourceAttr("powerflex.volume.avengers", "map_sdcs_id.0", "c423b09a00000005"),
				),
			},
		},
	})
}
