package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStoragepoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test
			{
				Config: StoragePoolResourceCreateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "Storage_91"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
				),
			},
			// Update Storagepool Test
			{
				Config: StoragePoolResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "Storage_91"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
				),
			},

			// createVolumeInvalidSizeTest
			{
				Config:      StoragePoolResourceCreateInvalidConfig,
				ExpectError: regexp.MustCompile(`.*Incorrect Name format.*`),
			},

			// createVolumeInvalidCapacityUnitTest
			{
				Config:      StoragePoolResourceUpdateInvalidConfig,
				ExpectError: regexp.MustCompile(`.*Media Type is required to have either HDD(Default) or SSD*.`),
			},
		},
	})
}

var StoragePoolResourceCreateConfig = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

resource "powerflex_storagepool" "storagepool" {
	name = "Storage_91"
	  protection_domain_id = "4eeb304600000000"
	  media_type = "HDD"
}
`

var StoragePoolResourceUpdateConfig = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

resource "powerflex_storagepool" "storagepool" {
		name = "Storage_91"
  		protection_domain_id = "4eeb304600000000"
  		media_type = "HDD"
  }
`

var StoragePoolResourceCreateInvalidConfig = `
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

var StoragePoolResourceUpdateInvalidConfig = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}

  resource "powerflex_storagepool" "storagepool" {
	name = "Storage_919"
	  protection_domain_id = "4eeb304600000000"
	  media_type = "SDD"
}
`
