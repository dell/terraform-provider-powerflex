package powerflextesting

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type dataPoints struct {
	noOfSdc   string
	name      string
	systemid  string
	sdcguid   string
	sdcip     string
	noOflinks string
}

var sdcTestData dataPoints

func init() {
	sdcTestData.noOfSdc = "1"
	sdcTestData.noOflinks = "4"
	sdcTestData.name = "powerflex_sdc13"
	sdcTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcTestData.sdcip = "10.247.96.90"
	sdcTestData.systemid = "0e7a082862fedf0f"
}

func TestAccCoffeesDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config: ProviderConfigForTesting + `data "powerflex_sdc" "selected" {
						sdcid = "c423b09800000003"
						id = ""
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", sdcTestData.noOfSdc),
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.systemid", sdcTestData.systemid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdcguid", sdcTestData.sdcguid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", sdcTestData.name),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdcip", sdcTestData.sdcip),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.links.#", sdcTestData.noOflinks),
				),
			},
		},
	})
}
