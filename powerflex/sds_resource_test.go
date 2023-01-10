package powerflex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSDSResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config: ProviderConfigForTesting + createSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "1"),
				),
			},
			// update sds test
			{
				Config: ProviderConfigForTesting + updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_02"),
					// resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list", "16"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
				),
			},
		},
	})
}

var createSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_01"
	ip_list = [
		"10.247.100.232"
	]
	protection_domain_id = "4eeb304600000000"
}
`

var updateSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_02"
	ip_list = [
		"10.247.100.232"
	]
	protection_domain_id = "4eeb304600000000"
}
`

// var createInvalidConfig = `

// resource "powerflex_storagepool" "storagepool" {
// 	name = "SP123"
// 	protection_domain_id = "4eeb304600000000"
// 	media_type = "HDD"
// }
// `

// var updateInvalidConfig = `
// resource "powerflex_sds" "sds" {
// 	name = "SDS_UPDATED"
// 	ip_list = [
// 		"10.247.101.60"
// 	]
// 	protection_domain_id = "4eeb304600000000"
// }
// `
