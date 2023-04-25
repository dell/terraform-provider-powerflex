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
	sdcTestData.sdcip = "10.247.66.194"
}

func TestSdcDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceBlockOnlyID,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first sdc to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "Terraform_sdc1"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.sdc_ip", sdcTestData.sdcip),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByEmptyBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.powerflex_sdc.selected", "sdcs.*", map[string]string{
						"sdcs.0.id":   "e3ce1fb500000000",
						"sdcs.0.name": "terraform_sdc",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.powerflex_sdc.selected", "sdcs.*", map[string]string{
						"sdcs.1.id":   "e3ce1fb600000001",
						"sdcs.1.name": "Terraform_sdc1",
					}),
				),
			},
			{
				Config: ProviderConfigForTesting + TestSdcDataSourceByNameBlock,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.name", "Terraform_sdc"),
					resource.TestCheckResourceAttr("data.powerflex_sdc.selected", "sdcs.0.id", "e3ce1fb500000000"),
				),
			},
		},
	})
}

func TestSdcDataSourceNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// Error here = https://github.com/hashicorp/terraform-plugin-sdk/pull/1077
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceByEmptyIDNeg,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceBlockBothNeg,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Combination.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceByEmptyNameBlock,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
			{
				Config:      ProviderConfigForTesting + TestSdcDataSourceInvalidName,
				ExpectError: regexp.MustCompile(".*Couldn't find SDC.*"),
			},
		},
	})
}

var (
	TestSdcDataSourceBlockOnlyID = `data "powerflex_sdc" "selected" {
		id = "e3ce1fb600000001"
	}`

	TestSdcDataSourceByEmptyIDNeg = `data "powerflex_sdc" "selected" {
		id = ""
	}`

	TestSdcDataSourceBlockBothNeg = `data "powerflex_sdc" "selected" {
		id = "e3ce1fb600000001"
		name = "Terraform_sdc1"
	}`

	TestSdcDataSourceByEmptyNameBlock = `data "powerflex_sdc" "selected" {
		name = ""
	}`

	TestSdcDataSourceByEmptyBlock = `data "powerflex_sdc" "selected" {
	}`

	TestSdcDataSourceByNameBlock = `data "powerflex_sdc" "selected" {
		name = "Terraform_sdc"
	}`

	TestSdcDataSourceInvalidName = `data "powerflex_sdc" "selected" {
		name = "Terraform_sdc11"
	}`
)
