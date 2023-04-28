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
	resourceName := "powerflex_storage_pool.storagepool"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", protection_domain_id),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "true"),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool_new"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", protection_domain_id),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "true"),
				),
			},
			// check that import is creating correct state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", protection_domain_id),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config: ProviderConfigForTesting + StoragePoolResourceCreateRMCacheFalse,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "false"),
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
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match.*`),
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
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
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
			Config: ProviderConfigForTesting + StoragePoolResourceCreateRMCacheFalse,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
				resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_name", "domain1"),
				resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
				resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "false"),
				resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "false"),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", protection_domain_id),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config:      ProviderConfigForTesting + CreateInvalidProtectionDomainID,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain.*`),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "name", "storage_pool"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", protection_domain_id),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "use_rfcache", "true"),
				),
			},
			// Update Storagepool Test
			{
				Config:      ProviderConfigForTesting + CreateInvalidName,
				ExpectError: regexp.MustCompile(`.*Error while updating name of Storagepool.*`),
			},
		},
	})
}

func TestAccStoragepoolResourceNegativeCases(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test Negative
			{
				Config:      ProviderConfigForTesting + CreateExistingStoragePoolName,
				ExpectError: regexp.MustCompile(`.*Error creating Storage Pool.*`),
			},
		},
	})
}

var StoragePoolResourceCreate = `
resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "` + protection_domain_id + `"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
}
`
var StoragePoolResourceCreateRMCacheFalse = `
resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_name = "domain1"
	media_type = "HDD"
	use_rmcache = false
	use_rfcache = false
}
`
var StoragePoolResourceUpdate = `
resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool_new"
	protection_domain_id = "` + protection_domain_id + `"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
var CreateInvalidMediaType = `
  resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "` + protection_domain_id + `"
	media_type = "HSD"
	use_rmcache = true
	use_rfcache = true
}
`
var CreateInvalidProtectionDomainID = `
resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "123"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
var CreateInvalidName = `
resource "powerflex_storage_pool" "storagepool" {
	name = "Terraform_POWERFLEX_storage_pool_33"
	protection_domain_id = "` + protection_domain_id + `"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`

var CreateExistingStoragePoolName = `
resource "powerflex_storage_pool" "storagepool1" {
	name = "pool1"
	protection_domain_id = "` + protection_domain_id + `"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
