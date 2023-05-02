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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", "4eeb304600000000"),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", "4eeb304600000000"),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", "4eeb304600000000"),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", "4eeb304600000000"),
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
					resource.TestCheckResourceAttr("powerflex_storage_pool.storagepool", "protection_domain_id", "4eeb304600000000"),
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

func TestAccStoragepoolResourceInvalidConfig(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test Negative
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig1,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig2,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig3,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig4,
				ExpectError: regexp.MustCompile(`.*rm_cache_write_handling_mode cannot be specified.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig5,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig6,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig7,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig8,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig9,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolWithInvalidConfig10,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
		},
	})
}

func TestAccStoragepoolResourceManyAttributes(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test Negative
			{
				Config: ProviderConfigForTesting + CreateStoragePoolWithAllAttributesConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "name", "storagepool1"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "use_rfcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "zero_padding_enabled", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "capacity_alert_high_threshold", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "capacity_alert_critical_threshold", "77"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "protected_maintenance_mode_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "protected_maintenance_mode_num_of_concurrent_ios_per_device", "7"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "protected_maintenance_mode_bw_limit_per_device_in_kbps", "1028"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebalance_enabled", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebalance_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebalance_num_of_concurrent_ios_per_device", "7"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebalance_bw_limit_per_device_in_kbps", "1032"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "vtree_migration_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "vtree_migration_num_of_concurrent_ios_per_device", "7"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "vtree_migration_bw_limit_per_device_in_kbps", "1030"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "spare_percentage", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rm_cache_write_handling_mode", "Passthrough"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebuild_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "rebuild_rebalance_parallelism", "5"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp1", "fragmentation", "false"),
				),
			},

			{
				Config: ProviderConfigForTesting + CreateStoragePoolWithAllAttributesConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "name", "storagepool2"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "use_rfcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "zero_padding_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "capacity_alert_high_threshold", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "capacity_alert_critical_threshold", "77"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "protected_maintenance_mode_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "protected_maintenance_mode_num_of_concurrent_ios_per_device", "9"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rebalance_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rebalance_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rebalance_num_of_concurrent_ios_per_device", "8"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "vtree_migration_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "vtree_migration_num_of_concurrent_ios_per_device", "10"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "spare_percentage", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rm_cache_write_handling_mode", "Cached"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rebuild_enabled", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "rebuild_rebalance_parallelism", "6"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp2", "fragmentation", "true"),
				),
			},
		},
	})
}

func TestAccStoragepoolResourceCapacityAlertInvalidValue(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test Negative
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolInvalidAttributesValue1,
				ExpectError: regexp.MustCompile(`.*Could not set capacity alert high threshold.*`),
			},
		},
	})
}

func TestAccStoragepoolResourceInvalidAttributesValue(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Storagepool Test Negative
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolInvalidAttributesValue2,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolInvalidAttributesValue3,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePoolInvalidAttributesValue4,
				ExpectError: regexp.MustCompile(`.*Attribute Error.*`),
			},
		},
	})
}

func TestAccStoragePoolResourceUpdate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + CreateUpdateStoragePoolWithAllAttributesConfig1,
			},
			{
				Config: ProviderConfigForTesting + UpdateStoragePoolWithAllAttributesConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "name", "storagepool3"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "use_rfcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "zero_padding_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "capacity_alert_high_threshold", "69"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "capacity_alert_critical_threshold", "72"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "protected_maintenance_mode_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "protected_maintenance_mode_num_of_concurrent_ios_per_device", "9"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "protected_maintenance_mode_bw_limit_per_device_in_kbps", "2042"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_num_of_concurrent_ios_per_device", "9"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_bw_limit_per_device_in_kbps", "2047"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "vtree_migration_io_priority_policy", "favorAppIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "vtree_migration_num_of_concurrent_ios_per_device", "7"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "vtree_migration_bw_limit_per_device_in_kbps", "1803"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "spare_percentage", "77"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rm_cache_write_handling_mode", "Cached"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebuild_rebalance_parallelism", "10"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "fragmentation", "true"),
				),
			},
			{
				Config: ProviderConfigForTesting + UpdateateStoragePoolWithAllAttributesConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "name", "storagepool5"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "use_rfcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "zero_padding_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "capacity_alert_high_threshold", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "capacity_alert_critical_threshold", "77"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "protected_maintenance_mode_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "protected_maintenance_mode_num_of_concurrent_ios_per_device", "9"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebalance_num_of_concurrent_ios_per_device", "8"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "vtree_migration_io_priority_policy", "limitNumOfConcurrentIos"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "vtree_migration_num_of_concurrent_ios_per_device", "10"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "spare_percentage", "66"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rm_cache_write_handling_mode", "Cached"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "rebuild_rebalance_parallelism", "6"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3", "fragmentation", "true"),
				),
			},
			{
				Config: ProviderConfigForTesting + CreateUpdateStoragePoolCacheAttribute,
			},
			{
				Config: ProviderConfigForTesting + UpdateStoragePoolCacheAttribute,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "name", "storagepool6"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "use_rfcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "replication_journal_capacity", "8"),
					resource.TestCheckResourceAttr("powerflex_storage_pool.sp3cache", "rebuild_enabled", "true"),
				),
			},
		},
	},
	)
}

var StoragePoolResourceCreate = `
resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
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
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`
var CreateInvalidMediaType = `
  resource "powerflex_storage_pool" "storagepool" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
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
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`

var CreateExistingStoragePoolName = `
resource "powerflex_storage_pool" "storagepool1" {
	name = "pool1"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`

// with protected maintenance mode policy as limitNumOfConcurrentIos , user can't pass bw limit per device
var CreateStoragePoolWithInvalidConfig1 = `
resource "powerflex_storage_pool" "sp1" {
	name                 = "storagepool1"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	protected_maintenance_mode_io_priority_policy = "limitNumOfConcurrentIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  }
`

// with rebalance policy as limitNumOfConcurrentIos , user can't pass bw limit per device
var CreateStoragePoolWithInvalidConfig2 = `
resource "powerflex_storage_pool" "sp2" {
	name                 = "storagepool2"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	rebalance_io_priority_policy = "limitNumOfConcurrentIos"
  	rebalance_num_of_concurrent_ios_per_device = 7
  	rebalance_bw_limit_per_device_in_kbps = 1032
  }
`

// with vtree migration policy as limitNumOfConcurrentIos , user can't pass bw limit per device
var CreateStoragePoolWithInvalidConfig3 = `
resource "powerflex_storage_pool" "sp3" {
	name                 = "storagepool3"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
  	vtree_migration_num_of_concurrent_ios_per_device = 7
  	vtree_migration_bw_limit_per_device_in_kbps = 1030
  }
`

// with rm_cache disabled, user can't configure the write handling mode
var CreateStoragePoolWithInvalidConfig4 = `
resource "powerflex_storage_pool" "sp4" {
	name                 = "storagepool4"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	use_rmcache = false
	rm_cache_write_handling_mode = "Passthrough"
  }
`

var CreateStoragePoolWithInvalidConfig5 = `
resource "powerflex_storage_pool" "sp5" {
	name                 = "storagepool5"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	protected_maintenance_mode_io_priority_policy = "unlimited"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  }
`

var CreateStoragePoolWithInvalidConfig6 = `
resource "powerflex_storage_pool" "sp6" {
	name                 = "storagepool6"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	rebalance_io_priority_policy = "unlimited"
  	rebalance_num_of_concurrent_ios_per_device = 7
  	rebalance_bw_limit_per_device_in_kbps = 1032
  }
`

var CreateStoragePoolWithInvalidConfig7 = `
resource "powerflex_storage_pool" "sp7" {
	name                 = "storagepool7"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	vtree_migration_io_priority_policy = "unlimited"
  	vtree_migration_num_of_concurrent_ios_per_device = 7
  	vtree_migration_bw_limit_per_device_in_kbps = 1030
  }
`

var CreateStoragePoolWithInvalidConfig8 = `
resource "powerflex_storage_pool" "sp8" {
	name                 = "storagepool8"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  }
`

var CreateStoragePoolWithInvalidConfig9 = `
resource "powerflex_storage_pool" "sp9" {
	name                 = "storagepool9"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	rebalance_num_of_concurrent_ios_per_device = 7
	rebalance_bw_limit_per_device_in_kbps = 1032
  }
`

var CreateStoragePoolWithInvalidConfig10 = `
resource "powerflex_storage_pool" "sp10" {
	name                 = "storagepool10"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	vtree_migration_num_of_concurrent_ios_per_device = 7
	vtree_migration_bw_limit_per_device_in_kbps = 1030
  }
`

// Note: I'm not testing replication_journal_capacity because it needs some changes in the delete functionality
var CreateStoragePoolWithAllAttributesConfig1 = `
resource "powerflex_storage_pool" "sp1" {
	name                 = "storagepool1"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	capacity_alert_high_threshold = 66
	capacity_alert_critical_threshold = 77
	zero_padding_enabled = false
	protected_maintenance_mode_io_priority_policy = "favorAppIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
	rebalance_enabled = false
	use_rmcache = true
	use_rfcache = true
	rebalance_io_priority_policy = "favorAppIos"
	rebalance_num_of_concurrent_ios_per_device = 7
	rebalance_bw_limit_per_device_in_kbps = 1032
	vtree_migration_io_priority_policy = "favorAppIos"
	vtree_migration_num_of_concurrent_ios_per_device = 7
	vtree_migration_bw_limit_per_device_in_kbps = 1030
	spare_percentage = 66
	rm_cache_write_handling_mode = "Passthrough"
	rebuild_enabled = true
	rebuild_rebalance_parallelism = 5
	fragmentation = false
  }
`

// Note: I'm not testing replication_journal_capacity attribute because it needs some changes in the delete functionality
var CreateStoragePoolWithAllAttributesConfig2 = `
resource "powerflex_storage_pool" "sp2" {
	name                 = "storagepool2"
	protection_domain_name = "domain1"
	media_type  = "SSD"
	capacity_alert_high_threshold = 66
	capacity_alert_critical_threshold = 77
	zero_padding_enabled = true
	protected_maintenance_mode_io_priority_policy = "limitNumOfConcurrentIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 9
	rebalance_enabled = true
	use_rmcache = true
	use_rfcache = false
	rebalance_io_priority_policy = "limitNumOfConcurrentIos"
	rebalance_num_of_concurrent_ios_per_device = 8
	vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
	vtree_migration_num_of_concurrent_ios_per_device = 10
	spare_percentage = 66
	rm_cache_write_handling_mode = "Cached"
	rebuild_enabled = false
	rebuild_rebalance_parallelism = 6
	fragmentation = true
  }
`

// high threshold can't be greater than critical threshold
var CreateStoragePoolInvalidAttributesValue1 = `
resource "powerflex_storage_pool" "sp1" {
	name                 = "storagepool1"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	capacity_alert_high_threshold = 78
	capacity_alert_critical_threshold = 77
  }
`

// num of concurrent IOs and bandwidth limit can't be passed when policy is not mentioned
var CreateStoragePoolInvalidAttributesValue2 = `
resource "powerflex_storage_pool" "sp2" {
	name                 = "storagepool2"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
  }
`

// num of concurrent IOs and bandwidth limit can't be passed when policy is not mentioned
var CreateStoragePoolInvalidAttributesValue3 = `
resource "powerflex_storage_pool" "sp3" {
	name                 = "storagepool3"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	rebalance_num_of_concurrent_ios_per_device = 7
	rebalance_bw_limit_per_device_in_kbps = 1032
  }
`

// num of concurrent IOs and bandwidth limit can't be passed when policy is not mentioned
var CreateStoragePoolInvalidAttributesValue4 = `
resource "powerflex_storage_pool" "sp4" {
	name                 = "storagepool4"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	vtree_migration_num_of_concurrent_ios_per_device = 7
	vtree_migration_bw_limit_per_device_in_kbps = 1030
  }
`

var CreateUpdateStoragePoolWithAllAttributesConfig1 = `
resource "powerflex_storage_pool" "sp3" {
	name                 = "storagepool3"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	capacity_alert_high_threshold = 66
	capacity_alert_critical_threshold = 77
	zero_padding_enabled = false
	protected_maintenance_mode_io_priority_policy = "favorAppIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 7
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 1028
	rebalance_enabled = false
	use_rmcache = true
	use_rfcache = true
	rebalance_io_priority_policy = "favorAppIos"
	rebalance_num_of_concurrent_ios_per_device = 7
	rebalance_bw_limit_per_device_in_kbps = 1032
	vtree_migration_io_priority_policy = "favorAppIos"
	vtree_migration_num_of_concurrent_ios_per_device = 7
	vtree_migration_bw_limit_per_device_in_kbps = 1030
	spare_percentage = 66
	rm_cache_write_handling_mode = "Passthrough"
	rebuild_enabled = true
	rebuild_rebalance_parallelism = 5
	fragmentation = false
  }
`

var UpdateStoragePoolWithAllAttributesConfig1 = `
resource "powerflex_storage_pool" "sp3" {
	name                 = "storagepool3"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	capacity_alert_high_threshold = 69
	capacity_alert_critical_threshold = 72
	zero_padding_enabled = true
	protected_maintenance_mode_io_priority_policy = "favorAppIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 9
	protected_maintenance_mode_bw_limit_per_device_in_kbps = 2042
	rebalance_enabled = true
	use_rmcache = true
	use_rfcache = false
	rebalance_io_priority_policy = "favorAppIos"
	rebalance_num_of_concurrent_ios_per_device = 9
	rebalance_bw_limit_per_device_in_kbps = 2047
	vtree_migration_io_priority_policy = "favorAppIos"
	vtree_migration_num_of_concurrent_ios_per_device = 7
	vtree_migration_bw_limit_per_device_in_kbps = 1803
	spare_percentage = 77
	rm_cache_write_handling_mode = "Cached"
	rebuild_rebalance_parallelism = 10
	fragmentation = true
  }
`

// Note: I'm not testing replication_journal_capacity attribute because it needs some changes in the delete functionality
var UpdateateStoragePoolWithAllAttributesConfig2 = `
resource "powerflex_storage_pool" "sp3" {
	name                 = "storagepool5"
	protection_domain_name = "domain1"
	media_type  = "SSD"
	capacity_alert_high_threshold = 66
	capacity_alert_critical_threshold = 77
	zero_padding_enabled = true
	protected_maintenance_mode_io_priority_policy = "limitNumOfConcurrentIos"
	protected_maintenance_mode_num_of_concurrent_ios_per_device = 9
	rebalance_enabled = true
	use_rmcache = true
	use_rfcache = false
	rebalance_io_priority_policy = "limitNumOfConcurrentIos"
	rebalance_num_of_concurrent_ios_per_device = 8
	vtree_migration_io_priority_policy = "limitNumOfConcurrentIos"
	vtree_migration_num_of_concurrent_ios_per_device = 10
	spare_percentage = 66
	rm_cache_write_handling_mode = "Cached"
	rebuild_rebalance_parallelism = 6
	fragmentation = true
  }
`

var CreateUpdateStoragePoolCacheAttribute = `
resource "powerflex_storage_pool" "sp3cache" {
	name                 = "storagepool6"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	use_rmcache = false
	replication_journal_capacity = 12
	rebuild_enabled = false
	zero_padding_enabled = true
}
`

var UpdateStoragePoolCacheAttribute = `
resource "powerflex_storage_pool" "sp3cache" {
	name                 = "storagepool6"
	protection_domain_name = "domain1"
	media_type  = "HDD"
	use_rfcache = true 
	use_rmcache = true
	// replication_journal_capacity = 8
	rebuild_enabled = true
}
`
