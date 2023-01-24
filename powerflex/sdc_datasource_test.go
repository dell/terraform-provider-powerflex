package powerflex

import (
	"os"
	"regexp"
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

func init() {
	sdcTestData.noOfSdc = "1"
	sdcTestData.noOflinks = "4"
	sdcTestData.name = ""
	sdcTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcTestData.sdcip = "10.247.96.90"
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
				Config: ProviderConfigForTesting + TestSdcDataSourceBlockOnlyID,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", sdcTestData.noOfSdc),
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.system_id", sdcTestData.systemid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_guid", sdcTestData.sdcguid),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", ""),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_ip", sdcTestData.sdcip),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.links.#", sdcTestData.noOflinks),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByEmptyBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", "151"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "id", ""),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "name", ""),
				),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceByEmptyIDNeg,
				ExpectError: regexp.MustCompile(".*id.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceBlockBothNeg,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Combination.*"),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByEmptyNameBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", "147"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "id", ""),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "name", ""),
				),
			},
		},
	})
}

func TestSdcDataSourceByName(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByNameBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "name", "LGLW6092"),
				),
			},
		},
	})
}

var (
	TestSdcDataSourceBlockOnlyID = `data "powerflex_sdc" "selected" {
		id = "c423b09800000003"
	}`
	TestSdcDataSourceByEmptyIDNeg = `data "powerflex_sdc" "selected" {
		id = ""
	}`
	TestSdcDataSourceBlockBothNeg = `data "powerflex_sdc" "selected" {
		id = ""
		name = ""
	}`
	TestSdcDataSourceByNameBlock = `data "powerflex_sdc" "selected" {
		name = "LGLW6092"
	}`
	TestSdcDataSourceByEmptyNameBlock = `data "powerflex_sdc" "selected" {
		name = ""
	}`
	TestSdcDataSourceByEmptyBlock = `data "powerflex_sdc" "selected" {
	}`
)
