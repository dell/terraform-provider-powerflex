package powerflex

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccStoragepoolResource
func TestAccStoragepoolResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {

		t.Skip("Dont run with units tests because it will try to create the context")

	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "qwert"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "st123new"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
		},
	})
}

// TestAccStoragepoolResourceInvalidCreate
func TestAccStoragepoolResourceInvalidCreate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {

		t.Skip("Dont run with units tests because it will try to create the context")

	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + CreateInvalidMediaType,
				ExpectError: regexp.MustCompile(`.*Could not create Storage Pool.*`),
			},
		},
	})
}

// TestAccStoragepoolResourceInvalidProtectionDomainID
func TestAccStoragepoolResourceInvalidProtectionDomainID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {

		t.Skip("Dont run with units tests because it will try to create the context")

	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + CreateInvalidProtectionDomainID,
				ExpectError: regexp.MustCompile(`.*Couldn't find protection domain.*`),
			},
		},
	})
}

var StoragePoolResourceCreate = `
resource "powerflex_storagepool" "storagepool" {
	name = "qwert"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
}
`

var StoragePoolResourceUpdate = `
resource "powerflex_storagepool" "storagepool" {
	name = "st123new"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`

var CreateInvalidMediaType = `
  resource "powerflex_storagepool" "storagepool" {
	name = "asd"
	protection_domain_id = "4eeb304600000000"
	media_type = "SAD"
	use_rmcache = true
	use_rfcache = true
}
`

var CreateInvalidProtectionDomainID = `
resource "powerflex_storagepool" "storagepool" {
	name = "asd"
	protection_domain_id = "123"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
