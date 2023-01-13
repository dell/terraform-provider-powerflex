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
			// update sds name test
			{
				Config: ProviderConfigForTesting + updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
				),
			},
		},
	})
}

func TestAccSDSResourceManyIP(t *testing.T) {
	createSDSTestMany := `
		resource "powerflex_sds" "sds" {
			name = "Tf_SDS_01"
			ip_list = [
				{
					ip = "10.247.100.232"
					role = "all"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.2"
					role = "sdcOnly"
				}
			]
			protection_domain_id = "4eeb304600000000"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config: ProviderConfigForTesting + createSDSTestMany,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "3"),
				),
			},
		},
	})
}

var createSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_01"
	ip_list = [
		{
			ip = "10.247.100.232"
			role = "all"
		}
	]
	rmcache_enabled = true
	rmcache_size_in_kb = 256000
	# num_of_io_buffers = 4
	drl_mode = "Volatile"
	protection_domain_id = "4eeb304600000000"
}
`

var updateSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_02"
	ip_list = [
		{
			ip = "10.247.100.232"
			role = "all"
		}
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
