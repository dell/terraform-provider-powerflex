package powerflex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type resourceDataPoints struct {
	noOfSdc   string
	name      string
	newname   string
	systemid  string
	sdcguid   string
	sdcip     string
	noOflinks string
}

var sdcResourceTestData resourceDataPoints

func init() {
	sdcResourceTestData.noOfSdc = "1"
	sdcResourceTestData.noOflinks = "4"
	sdcResourceTestData.name = "powerflex_sdc21"
	sdcResourceTestData.newname = "powerflex_sdc56"
	sdcResourceTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcResourceTestData.sdcip = "10.247.96.90"
	sdcResourceTestData.systemid = "0e7a082862fedf0f"
}

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfigForTesting + `
				resource "powerflex_sdc" "sdc" {
					sdcid = "c423b09800000003"
					name = "` + sdcResourceTestData.name + `"
				  }
				  `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "systemid", sdcResourceTestData.systemid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdcguid", sdcResourceTestData.sdcguid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", sdcResourceTestData.name),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdcip", sdcResourceTestData.sdcip),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "links.#", sdcResourceTestData.noOflinks),
				),
			},
		},
	})
}
