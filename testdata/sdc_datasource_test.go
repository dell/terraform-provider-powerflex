package powerflextesting

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config: ProviderConfigForTesting + `data "powerflex_sdc" "selected" {
						systemid = "0e7a082862fedf0f"
						sdcid = "c423b09800000003"
						id = "c423b09800000003"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", "1"),
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.systemid", "0e7a082862fedf0f"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdcguid", "0877AE5E-BDBF-4E87-A002-218D9F883896"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "LGLW6090"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdcip", "10.247.96.90"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.links.#", "4"),
				),
			},
		},
	})
}
