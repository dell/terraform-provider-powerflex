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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool_new"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
		},
	})
}
func TestAccStoragepoolResourceUpdateRMCache(t *testing.T) {
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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceCreateRMCacehFalse,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "false"),
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
func TestAccStoragepoolResourceVariousCases(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	tests := []resource.TestStep{
		{
			Config: ProviderConfigForTesting + StoragePoolResourceCreateRMCacehFalse,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
				resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
				resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
				resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "false"),
				resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "false"),
			),
		},
	}

	for i := range tests {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps:                    []resource.TestStep{tests[i]},
		})
	}
}
func TestAccStoragepoolResourceInvalidUpdateProtectionDomainID(t *testing.T) {
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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config:      ProviderConfigForTesting + CreateInvalidProtectionDomainID,
				ExpectError: regexp.MustCompile(`.*Unable to find protection domain.*`),
			},
		},
	})
}
func TestAccStoragepoolResourceInvalidUpdateName(t *testing.T) {
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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config:      ProviderConfigForTesting + CreateInvalidName,
				ExpectError: regexp.MustCompile(`.*Could not get Storagepool.*`),
			},
		},
	})
}
func TestAccStoragepoolResourceInvalidUpdateMediaType(t *testing.T) {
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
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config:      ProviderConfigForTesting + CreateInvalidMediaType,
				ExpectError: regexp.MustCompile(`.*Could not get Storagepool.*`),
			},
		},
	})
}

var StoragePoolResourceCreate = `
resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
}
`
var StoragePoolResourceCreateRMCacehFalse = `
resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = false
	use_rfcache = false
}
`
var StoragePoolResourceUpdate = `
resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool_new"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
var CreateInvalidMediaType = `
  resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
	media_type = "SSD"
	use_rmcache = true
	use_rfcache = true
}
`
var CreateInvalidProtectionDomainID = `
resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "123"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
var CreateInvalidName = `
resource "powerflex_storagepool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
