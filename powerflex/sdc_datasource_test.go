package powerflex

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type sdcDataPoints struct {
	noOfSdc   string
	name      string
	systemid  string
	sdcguid   string
	sdcip     string
	noOflinks string
}

var sdcTestData sdcDataPoints

var providerConfigForTesting = `
provider "powerflex" {
	username = ""
	password = ""
	endpoint = ""
	insecure = true
}
`

func init() {
	sdcTestData.noOfSdc = "1"
	sdcTestData.noOflinks = "4"
	sdcTestData.name = "powerflex_sdc20"
	sdcTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcTestData.sdcip = ""
	sdcTestData.systemid = "0e7a082862fedf0f"
}

func TestSdcDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config: providerConfigForTesting + `data "powerflex_sdc" "selected" {
						sdc_id = "c423b09800000003"
						id = ""
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", sdcTestData.noOfSdc),
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.system_id", sdcTestData.systemid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_guid", sdcTestData.sdcguid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", sdcTestData.name),
					// resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_ip", sdcTestData.sdcip),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.links.#", sdcTestData.noOflinks),
				),
			},
		},
	})
}
