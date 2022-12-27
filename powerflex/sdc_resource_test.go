package powerflex

import (
	"os"
	"regexp"
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
	sdcResourceTestData.sdcip = ""
	sdcResourceTestData.systemid = "0e7a082862fedf0f"
}

func TestSdcResourceCreate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigForTesting + `
				resource "powerflex_sdc" "sdc" {
					sdc_id = "c423b09800000003"
					name = "` + sdcResourceTestData.name + `"
				  }
				  `,
				ExpectError: regexp.MustCompile(`.*SDC can not be added*`),
			},
		},
	})
}
func TestSdcResourceUpdate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigForTesting + `
				resource "powerflex_sdc" "sdc" {
					sdc_id = "c423b09800000003"
					name = "` + sdcResourceTestData.name + `"
				  }
				  `,
				ExpectError: regexp.MustCompile(`.*SDC can not be added*`),
			},
			{
				ResourceName:        "powerflex_sdc.sdc",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: "c423b09800000003",
			},
			// Update testing
			{
				Config: providerConfigForTesting + `
				resource "powerflex_sdc" "sdc" {
					sdc_id = "c423b09800000003"
					name = "` + sdcResourceTestData.name + `"
				  }
				  `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "system_id", sdcResourceTestData.systemid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdc_guid", sdcResourceTestData.sdcguid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", sdcResourceTestData.newname),
					// resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdcip", sdcResourceTestData.sdcip),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "links.#", sdcResourceTestData.noOflinks),
				),
			},
		},
	})
}
