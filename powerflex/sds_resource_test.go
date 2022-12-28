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
				Config: createSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "size", "8"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
				),
			},
			// update sds test
			{
				Config: updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list", "16"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
				),
			},
		},
	})
}

var createSDSTest = `
provider "powerflex" {
  username = "` + username + `"
  password = "` + password + `"
  endpoint = "` + endpoint + `"
  insecure = true
}

resource "powerflex_sds" "sds" {
	name = "SDS_01"
	ip_list = [
		"10.247.101.60"
	  ]
	protection_domain_id = "4eeb304600000000"
}
`

var updateSDSTest = `
provider "powerflex" {
  username = "` + username + `"
  password = "` + password + `"
  endpoint = "` + endpoint + `"
  insecure = true
}

resource "powerflex_sds" "sds" {
	name = "SDS_01"
	ip_list = [
		"10.247.101.60"
	  ]
	protection_domain_id = "4eeb304600000000"
}
`

var createInvalidConfig = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}		

  resource "powerflex_storagepool" "storagepool" {
	name = "SP123"
	  protection_domain_id = "4eeb304600000000"
	  media_type = "HDD"
}
`

var updateInvalidConfig = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

resource "powerflex_sds" "sds" {
	name = "SDS_UPDATED"
	ip_list = [
		"10.247.101.60"
	  ]
	protection_domain_id = "4eeb304600000000"
}	
`
