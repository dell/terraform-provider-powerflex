package powerflex

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
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
	sdcResourceTestData.name = "powerflex_sdc26"
	sdcResourceTestData.newname = ""
	sdcResourceTestData.sdcguid = "0877AE5E-BDBF-4E87-A002-218D9F883896"
	sdcResourceTestData.sdcip = ""
	sdcResourceTestData.systemid = "0e7a082862fedf0f"
}

func TestSdcResourceUpdate(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "c423b09800000003",
			},
			// Update testing
			{
				Config:      ProviderConfigForTesting + TestSdcResourceUpdateBlock,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Length*`),
			},
		},
	})
}

func TestSdcResourceUpdateWrongID(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "c423b098000000034343",
				ExpectError:       regexp.MustCompile(`.*Unable to Read Powerflex systems-sdcs Read*`),
			},
		},
	})
}

func TestSdcResourceUpdateSameName(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            ProviderConfigForTesting + TestSdcResourceUpdateImportBlock,
				ResourceName:      "powerflex_sdc.test_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "c423b09800000003",
			},
			// Update testing
			{
				Config:      ProviderConfigForTesting + TestSdcResourceUpdateBlockSameName,
				ExpectError: regexp.MustCompile(`.*Unable to Change name Powerflex sdc*`),
			},
		},
	})
}

func TestSdcResourceCreateUpdate(t *testing.T) {
	var TestSdcResourceCreateUpdateBlockS1 = `
	resource "powerflex_sdc" "sdc" {
		id = "c423b09900000004"
		name = "` + SdsResourceTestData.volName3 +`"
	  }
	  `
	var TestSdcResourceCreateUpdateBlockS2 = `
	resource "powerflex_sdc" "sdc" {
		id = "c423b09900000004"
		name = "` + SdsResourceTestData.volName2 +`"
	  }
	  `
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TestSdcResourceCreateUpdateBlockS1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", SdsResourceTestData.volName3),
				),
			},
			// // Update testing
			{
				Config: ProviderConfigForTesting + TestSdcResourceCreateUpdateBlockS2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sdc.sdc", "name", SdsResourceTestData.volName2),
				),
			},
		},
	})
}

var (
	TestSdcResourceCreateBlock = `
	resource "powerflex_sdc" "sdc" {
		id = "c423b09800000003"
		name = "` + sdcResourceTestData.name + `"
	  }
	  `

	TestSdcResourceUpdateBlock = `
	  resource "powerflex_sdc" "sdc" {
		  id = "c423b09800000003"
		  name = "` + sdcResourceTestData.newname + `"
		}
		`

	TestSdcResourceUpdateBlockSameName = `
	resource "powerflex_sdc" "sdc" {
		id = "c423b09800000003"
		name = "powerflex_sdc26"
		}
		`

	TestSdcResourceUpdateImportBlock = `
	resource "powerflex_sdc" "test_import" {
		id = "c423b09800000003"
		name = "powerflex_sdc26"
	  }
	  `


)
