package powerflex

import (
	"os"
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
	sdcResourceTestData.name = "powerflex_sdc22"
	sdcResourceTestData.newname = "powerflex_sdc56"
	sdcResourceTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcResourceTestData.sdcip = ""
	sdcResourceTestData.systemid = "0e7a082862fedf0f"
}

func TestSdcResourceUpdate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// {
			// 	Config:      providerConfigForTesting + TestSdcResourceCreateBlock,
			// 	ExpectError: regexp.MustCompile(`.*SDC can not be added*`),
			// },
			{
				Config:            providerConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "c423b09800000003",
			},
			// // Update testing
			{
				Config: providerConfigForTesting + TestSdcResourceUpdateBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "system_id", sdcResourceTestData.systemid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdc_guid", sdcResourceTestData.sdcguid),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", ""),
					// resource.TestCheckResourceAttr("powerflex_sdc.sdc", "sdcip", sdcResourceTestData.sdcip),
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "links.#", sdcResourceTestData.noOflinks),
				),
			},
		},
	})
}

func TestSdcResourceCreateUpdate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigForTesting + TestSdcResourceCreateUpdateBlockS1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", "Block_S33"),
				),
			},
			// // Update testing
			{
				Config: providerConfigForTesting + TestSdcResourceCreateUpdateBlockS2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", "Block_S34"),
				),
			},
		},
	})
}

var (
	TestSdcResourceCreateBlock = `
	resource "powerflex_sdc" "sdc" {
		sdc_id = "c423b09800000003"
		name = "` + sdcResourceTestData.name + `"
	  }
	  `

	TestSdcResourceUpdateBlock = `
	  resource "powerflex_sdc" "sdc" {
		  sdc_id = "c423b09800000003"
		  name = "` + sdcResourceTestData.name + `"
		}
		`

	TestSdcResourceUpdateImportBlock = `
	resource "powerflex_sdc" "test_import" {
		id = "c423b09800000003"
		sdc_id = "c423b09800000003"
	  }
	  `
	TestSdcResourceCreateUpdateBlockS1 = `
	resource "powerflex_sdc" "sdc" {
		sdc_id = "c423b09900000004"
		name = "Block_S33"
	  }
	  `
	TestSdcResourceCreateUpdateBlockS2 = `
	resource "powerflex_sdc" "sdc" {
		sdc_id = "c423b09900000004"
		name = "Block_S34"
	  }
	  `
)
