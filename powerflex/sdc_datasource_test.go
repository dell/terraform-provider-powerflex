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
	sdcTestData.sdcguid = "C87ACC43-298B-4AD3-A95F-344FE83192C6"
	sdcTestData.sdcip = SdsResourceTestData.SdcIP
	sdcTestData.systemid = "09a186f8167ebe0f"
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
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "terraform_sdc"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_ip", sdcTestData.sdcip),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.links.#", sdcTestData.noOflinks),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByEmptyBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of sdc returned
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
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "id", ""),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "name", ""),
				),
			},
		},
	})
}

func TestSdcDataSourceByName(t *testing.T) {
	var TestSdcDataSourceByNameBlock = `data "powerflex_sdc" "selected" {
		name = "` + SdsResourceTestData.volName + `"
	}`
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
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "name", SdsResourceTestData.volName),
				),
			},
		},
	})
}

var (
	TestSdcDataSourceBlockOnlyID = `data "powerflex_sdc" "selected" {
		id = "e3ce1fb500000000"
	}`
	TestSdcDataSourceByEmptyIDNeg = `data "powerflex_sdc" "selected" {
		id = ""
	}`
	TestSdcDataSourceBlockBothNeg = `data "powerflex_sdc" "selected" {
		id = ""
		name = ""
	}`

	TestSdcDataSourceByEmptyNameBlock = `data "powerflex_sdc" "selected" {
		name = ""
	}`
	TestSdcDataSourceByEmptyBlock = `data "powerflex_sdc" "selected" {
	}`
)
