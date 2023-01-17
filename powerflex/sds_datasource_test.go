package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSdsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.0.name", "SDS_Test_Tf"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "sds_details.0.id", "6ae199c500000007"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example1", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.name", "SDS_Test_Tf"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "sds_details.0.id", "6ae199c500000007"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example2", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.0.name", "SDS_Test_Tf"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "sds_details.0.id", "6ae199c500000007"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example3", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.0.name", "SDS_Test_Tf"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "sds_details.0.id", "6ae199c500000007"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example4", "protection_domain_name", "domain1"),
				),
			},
			{
				Config: ProviderConfigForTesting + SdsDataSourceConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_sds.example5", "sds_details.#", "10"),
					resource.TestCheckResourceAttr("data.powerflex_sds.example5", "protection_domain_id", "4eeb304600000000"),
				),
			},
			{
				Config:      ProviderConfigForTesting + SdsDataSourceConfig6,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain.*`),
			},
			{
				Config:      ProviderConfigForTesting + SdsDataSourceConfig7,
				ExpectError: regexp.MustCompile(`.*Unable to read SDS.*`),
			},
		},
	})
}

var SdsDataSourceConfig1 = `
data "powerflex_sds" "example1" {
	protection_domain_id = "4eeb304600000000"
	sds_names = ["SDS_Test_Tf"]
}
`

var SdsDataSourceConfig2 = `
data "powerflex_sds" "example2" {
	protection_domain_id = "4eeb304600000000"
	sds_ids = ["6ae199c500000007"]
}
`

var SdsDataSourceConfig3 = `
data "powerflex_sds" "example3" {
	protection_domain_name = "domain1"
	sds_names = ["SDS_Test_Tf"]
}
`

var SdsDataSourceConfig4 = `
data "powerflex_sds" "example4" {
	protection_domain_name = "domain1"
	sds_ids = ["6ae199c500000007"]
}
`

var SdsDataSourceConfig5 = `
data "powerflex_sds" "example5" {
	protection_domain_id = "4eeb304600000000"
}
`

var SdsDataSourceConfig6 = `
data "powerflex_sds" "example6" {
	protection_domain_id = "invalid_domain"
}
`

var SdsDataSourceConfig7 = `
data "powerflex_sds" "example7" {
	protection_domain_id = "4eeb304600000000"
	sds_ids = ["invalid_sds_id"]
}
`
