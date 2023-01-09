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
				ExpectError: regexp.MustCompile(`.*Error creating Storage Pool*`),
			},
		},
	})
}

// func TestStoragepoolResourceResourceImport(t *testing.T) {
// 	os.Setenv("TF_ACC", "1")
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:            ProviderConfigForTesting + TestStoragepoolResourceImportBlock,
// 				ResourceName:      "powerflex_storagepool.storagepool",
// 				ImportState:       true,
// 				ImportStateVerify: false,
// 				ImportStateId:     "7630a24800000002",
// 			},
// 		},
// 	})
// }

func TestStoragepoolResourceByAnika(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + CreateSPWOSPName,
				ExpectError: regexp.MustCompile(`.*Missing required argument*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWOPDIDOrName,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Combination*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWOMediaType,
				ExpectError: regexp.MustCompile(`.*Missing required argument*.`),
			},
			//bug
			// {
			// 	Config:      ProviderConfigForTesting + CreateSPWithEmptySPName,
			// 	ExpectError: regexp.MustCompile(`.*Missing required argument*.`),
			// },
			{
				Config:      ProviderConfigForTesting + CreateSPWithEmptyPDID,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithEmptyMediaType,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithSpacesInName,
				ExpectError: regexp.MustCompile(`.*Error creating Storage Pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithInvalidPDID,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithInvalidMediaType,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithExistingSPName,
				ExpectError: regexp.MustCompile(`.*Error creating Storage Pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithInvalidRmCache,
				ExpectError: regexp.MustCompile(`.*Incorrect attribute value type*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithInvalidRfCache,
				ExpectError: regexp.MustCompile(`.*Incorrect attribute value type*.`),
			},
			//bug
			// {
			// 	Config:      ProviderConfigForTesting + CreateSPWithRequiredParams,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "name", "ses-storage-pool6"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "protection_domain_id", "4eeb304600000000"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.storagepool", "use_rfcache", "false"),
			// 	),
			// },		
			{
				Config:      ProviderConfigForTesting + CreateSPWithRequiredAndOptionalParams,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "name", "ses-storage-pool7"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rfcache", "true"),
				),
			},
			//Bug
			// {
			// 	Config:      ProviderConfigForTesting + CreateSPWithPDNameandOneOptionalRmCache,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "name", "ses-storage-pool8"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "protection_domain_id", "4eeb304600000000"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rfcache", "false"),
			// 	),
			// },
				//Bug
			// {
			// 	Config:      ProviderConfigForTesting + CreateSPWithPDNameandOneOptionalRfCache,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "name", "ses-storage-pool8"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "protection_domain_id", "4eeb304600000000"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp2", "use_rfcache", "false"),
			// 	),
			// },

			{
				Config:      ProviderConfigForTesting + CreateSPWithExistingSPNameButDifferentPD,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp6", "name", "pool1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp6", "protection_domain_name", "domain2"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp6", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp6", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp6", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithSpecialCharactersinSPName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp7", "name", "ses_sp!@#$%^&*():;><,.?~"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp7", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp7", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp7", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp7", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp8", "name", "ses-storage-pool11"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp8", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp8", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp8", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp8", "use_rfcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp9", "name", "ses-storage-pool12"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp9", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp9", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp9", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp9", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifySPWithExistingSPName,
				ExpectError: regexp.MustCompile(`.*Error getting while updating Storagepool*.`),
			},
			//Was not able to test - Same issue as Discussed with Akash
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool2,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp10", "name", "ses-storage-pool13"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp10", "protection_domain_name", "domain1"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp10", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp10", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp10", "use_rfcache", "false"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + ModifySPWithInvalidRfCache,
			// 	ExpectError: regexp.MustCompile(`.*Inappropriate value for attribute.*`),
			// },
			//Was not able to test - Same issue as Discussed with Akash
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool3,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp11", "name", "ses-storage-pool14"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp11", "protection_domain_name", "domain1"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp11", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp11", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp11", "use_rfcache", "true"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + ModifySPWithInvalidRmCache,
			// 	ExpectError: regexp.MustCompile(`.*Inappropriate value for attribute.*`),
			// },
			//Was not able to test - Same issue as Discussed with Akash
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool4,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp12", "name", "ses-storage-pool15"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp12", "protection_domain_name", "domain1"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp12", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp12", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp12", "use_rfcache", "true"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + ModifySPWithInvalidMediaType,
			// 	ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match.*`),
			// },
			//Was not able to test - Same issue as Discussed with Akash
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool5,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp13", "name", "ses-storage-pool16"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp13", "protection_domain_name", "domain1"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp13", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp13", "use_rmcache", "false"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp13", "use_rfcache", "true"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + ModifySPWithEmptyMediaType,
			// 	ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match.*`),
			// },
			//Bug
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool6,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp16", "name", "ses-storage-pool19"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp16", "protection_domain_id", "4eeb304600000000"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp16", "media_type", "HDD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp16", "use_rmcache", "true"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp16", "use_rfcache", "false"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + ModifySPWithEmptyName,
			// 	ExpectError: regexp.MustCompile(`.*Error getting while updating Storagepool.*`),
			// },
			{
				Config:      ProviderConfigForTesting + CreateStoragePool7,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "name", "ses-storage-pool-20"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifySPNameUsingPDID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "name", "ses-storage-pool-new"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp17", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool8,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "name", "ses-storage-pool21"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyMediaType,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "name", "ses-storage-pool21"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp18", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool9,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "name", "ses-storage-pool-22"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyRmCache,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "name", "ses-storage-pool-22"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp19", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool10,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "name", "ses-storage-pool-23"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyRfCache,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "name", "ses-storage-pool-23"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp20", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool11,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "name", "ses-storage-pool-24"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyMultipleParams,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "name", "ses-storage-pool-new2"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "media_type", "Transitional"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp21", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool12,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "name", "ses-storage-pool-25"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyMultipleParamsPartII,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "name", "ses-storage-pool-new3"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "media_type", "HDD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp22", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool13,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp23", "name", "ses-storage-pool-25"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp23", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp23", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp23", "use_rmcache", "false"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp23", "use_rfcache", "true"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyMultipleParamsPartIII,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithEmptyPDName,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithInvalidPDName,
				ExpectError: regexp.MustCompile(`.*Error getting Protection Domain*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithBothPDIDAndName,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Combination*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateSPWithTransitionalMediaType,
				ExpectError: regexp.MustCompile(`.*Error creating Storage Pool*.`),
			},
			{
				Config:      ProviderConfigForTesting + CreateStoragePool14,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "name", "ses-storage-pool31"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "media_type", "SSD"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "use_rfcache", "false"),
				),
			},
			{
				Config:      ProviderConfigForTesting + ModifyMediaTypeSSDToTransitional,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "name", "ses-storage-pool31"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "protection_domain_name", "domain1"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "media_type", "Transitional"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "use_rmcache", "true"),
					resource.TestCheckResourceAttr("powerflex_storagepool.sp28", "use_rfcache", "false"),
				),
			},
			// {
			// 	Config:      ProviderConfigForTesting + CreateStoragePool15,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp29", "name", "ses-storage-pool32"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp29", "protection_domain_name", "domain1"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp29", "media_type", "SSD"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp29", "use_rmcache", "true"),
			// 		resource.TestCheckResourceAttr("powerflex_storagepool.sp29", "use_rfcache", "false"),
			// 	),
			// },
			// {
			// 	Config:      ProviderConfigForTesting + AttachVolumeToSP,
			// 	ExpectError: regexp.MustCompile(`.*absb*.`),
			// },
			
			
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
resource "powerflex_storagepool" "storagepool2" {
	name = "storage_pool"
	protection_domain_id = "4eeb304600000000"
	media_type = "HDD"
	use_rmcache = true
	use_rfcache = true
  }
`

// var TestStoragepoolResourceImportBlock = `
// resource "powerflex_storagepool" "storagepool" {
// 	name = "storage_pool"
// 	id = "7630a24800000002"
// 	protection_domain_id = "4eeb304600000000"
// 	media_type = "HDD"
// 	use_rmcache = true
// 	use_rfcache = true
//   }
// `

var CreateSPWOSPName = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"

}
`

var CreateSPWOPDIDOrName = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	name = "test-storage-pool1"
}
`
var CreateSPWOMediaType = `
resource "powerflex_storagepool" "sp1" {
	name= "test-storage-pool1"
	protection_domain_name = "domain1"

}
`

// var CreateSPWithEmptySPName = `
// resource "powerflex_storagepool" "sp1" {
// 	media_type= "HDD"
// 	protection_domain_id = "4eeb304600000000"
// 	name = ""
// }
// `

var CreateSPWithEmptyPDID = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = ""
	name = "ses-storage-pool1"
}
`

var CreateSPWithEmptyMediaType = `
resource "powerflex_storagepool" "sp1" {
	media_type= ""
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool1"
}
`

var CreateSPWithSpacesInName = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses storage pool 1"
}
`

var CreateSPWithInvalidPDID = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb30460000"
	name = "ses-storage-pool1"
}
`

var CreateSPWithInvalidMediaType = `
resource "powerflex_storagepool" "sp1" {
	media_type= "Media Type"
	protection_domain_name = "4eeb304600000000"
	name = "ses-storage-pool1"
}
`

var CreateSPWithExistingSPName = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "pool1"
}
`

var CreateSPWithInvalidRmCache = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool1"        
	use_rmcache = "use_rmCache"
}
`

var CreateSPWithInvalidRfCache = `
resource "powerflex_storagepool" "sp1" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool1"
	use_rmcache = true
	use_rfcache = "use_rfcache"
}
`

// var CreateSPWithRequiredParams = `
// resource "powerflex_storagepool" "sp1" {
// 	media_type= "HDD"
// 	protection_domain_id = "4eeb304600000000"
// 	name = "ses-storage-pool6"
// }
//`

var CreateSPWithRequiredAndOptionalParams = `
resource "powerflex_storagepool" "sp2" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool7"
	use_rmcache = true
	use_rfcache = true
}
`
// var CreateSPWithPDNameandOneOptionalRmCache = `
// resource "powerflex_storagepool" "sp3" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool8"
// 	use_rmcache = false
//  }
// `

// var CreateSPWithPDNameandOneOptionalRfCache = `
// resource "powerflex_storagepool" "sp4" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool9"
// 	use_rfcache = false
// }
// `

var CreateSPWithExistingSPNameButDifferentPD = `
resource "powerflex_storagepool" "sp6" {
	media_type= "SSD"
	protection_domain_name = "domain2"
	name = "pool1"
	use_rmcache = true
	use_rfcache = false
}
`

var CreateSPWithSpecialCharactersinSPName = `
resource "powerflex_storagepool" "sp7" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses_sp!@#$%^&*():;><,.?~"
	use_rmcache = true
	use_rfcache = false
}
`

var CreateStoragePool = `
resource "powerflex_storagepool" "sp8" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool11"
	use_rmcache = false
	use_rfcache = true
}
resource "powerflex_storagepool" "sp9" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool12"
	use_rmcache = false
	use_rfcache = true
}
`

var ModifySPWithExistingSPName = `
resource "powerflex_storagepool" "sp8" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool12" 
	use_rmcache = false
	use_rfcache = true
}
`
//Was not able to test - Same issue as Discussed with Akash
// var CreateStoragePool2 = `
// resource "powerflex_storagepool" "sp10" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool13"
// 	use_rmcache = false
// 	use_rfcache = false
// }
// `

// var ModifySPWithInvalidRfCache = `
// resource "powerflex_storagepool" "sp10" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool13"
// 	use_rmcache = false
// 	use_rfcache = "use_rfcache"
// }
// `

//Was not able to test - Same issue as Discussed with Akash
// var CreateStoragePool3 = `
// resource "powerflex_storagepool" "sp11" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool14"
// 	use_rmcache = false
// 	use_rfcache = true
// }
// `

// var ModifySPWithInvalidRmCache = `
// resource "powerflex_storagepool" "sp11" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool14"
// 	use_rmcache = "use_rmcache"
// 	use_rfcache = true
// }
// `

//Was not able to test - Same issue as Discussed with Akash
// var CreateStoragePool4 = `
// resource "powerflex_storagepool" "sp12" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool15"
// 	use_rmcache = false
// 	use_rfcache = true
// }
// `

// var ModifySPWithInvalidMediaType = `
// resource "powerflex_storagepool" "sp12" {
// 	media_type= "MediaType"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool15"
// 	use_rmcache = false
// 	use_rfcache = true
// }
// `

//Was not able to test - Same issue as Discussed with Akash
// var CreateStoragePool5 = `
// resource "powerflex_storagepool" "sp13" {
// 	media_type= "HDD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool16"
// 	use_rmcache = false
// 	use_rfcache = true
// }
// `

// var ModifySPWithEmptyMediaType = `
// resource "powerflex_storagepool" "sp13" {
// 	media_type= ""
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool16"
// 	use_rmcache = false
// 	use_rfcache = true
// }
// `

// var CreateStoragePool6 = `
// resource "powerflex_storagepool" "sp16" {
// 	media_type= "HDD"
// 	protection_domain_id = "4eeb304600000000"
// 	name = "ses-storage-pool19"
// 	use_rmcache = true
// 	use_rfcache = false
// }
// `

// var ModifySPWithEmptyName = `
// resource "powerflex_storagepool" "sp16" {
// 	media_type= "HDD"
// 	protection_domain_id = "4eeb304600000000"
// 	name = ""
// 	use_rmcache = true
// 	use_rfcache = false
// }
// `

var CreateStoragePool7 = `
resource "powerflex_storagepool" "sp17" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-20"
	use_rmcache = true
	use_rfcache = true
}
`

var ModifySPNameUsingPDID = `
resource "powerflex_storagepool" "sp17" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-new"
	use_rmcache = true
	use_rfcache = true
}
`

var CreateStoragePool8 = `
resource "powerflex_storagepool" "sp18" {
	media_type= "HDD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool21"
	use_rmcache = false
	use_rfcache = false
}
`

var ModifyMediaType = `
resource "powerflex_storagepool" "sp18" {
	media_type= "SSD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool21"
	use_rmcache = false
	use_rfcache = false
}
`

var CreateStoragePool9 = `
resource "powerflex_storagepool" "sp19" {
	media_type= "SSD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-22"
	use_rmcache = false
	use_rfcache = true
}
`

var ModifyRmCache = `
resource "powerflex_storagepool" "sp19" {
	media_type= "SSD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-22"
	use_rmcache = true
	use_rfcache = true
}
`

var CreateStoragePool10 = `
resource "powerflex_storagepool" "sp20" {
	media_type= "HDD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool-23"
	use_rmcache = true
	use_rfcache = false
}
`

var ModifyRfCache = `
resource "powerflex_storagepool" "sp20" {
	media_type= "HDD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool-23"
	use_rmcache = true
	use_rfcache = true
}
`

var CreateStoragePool11 = `
resource "powerflex_storagepool" "sp21" {
	media_type= "HDD"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-24"
	use_rmcache = true
	use_rfcache = true
}
`

var ModifyMultipleParams = `
resource "powerflex_storagepool" "sp21" {
	media_type= "Transitional"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool-new2"
   use_rmcache = false
   use_rfcache = false
}
`

var CreateStoragePool12 = `
resource "powerflex_storagepool" "sp22" {
	media_type= "SSD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool-25"
	use_rmcache = false
	use_rfcache = true
}
`

var ModifyMultipleParamsPartII = `
resource "powerflex_storagepool" "sp22" {
	media_type= "HDD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool-new3"
	use_rmcache = true
	use_rfcache = false
}
`
var CreateStoragePool13 = `
resource "powerflex_storagepool" "sp23" {
	media_type= "SSD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool-25"
	use_rmcache = false
	use_rfcache = true
}
`

var ModifyMultipleParamsPartIII = `
resource "powerflex_storagepool" "sp23" {
	media_type= "HDD"
	protection_domain_name = "abcdehdv"
	name = "ses-storage-pool-new3"
	use_rmcache = true
	use_rfcache = false
}
`

var CreateSPWithEmptyPDName = `
resource "powerflex_storagepool" "sp24" {
	media_type= "HDD"
	protection_domain_name = ""
	name = "ses-storage-pool-27"
}
`

var CreateSPWithInvalidPDName = `
resource "powerflex_storagepool" "sp25" {
	media_type= "HDD"
	protection_domain_name = "test_pd_name"
	name = "ses-storage-pool28"
}
`

var CreateSPWithBothPDIDAndName = `
resource "powerflex_storagepool" "sp26" {
	media_type= "HDD"
	protection_domain_name = "domain1"
	protection_domain_id = "4eeb304600000000"
	name = "ses-storage-pool29"
}
`

var CreateSPWithTransitionalMediaType = `
resource "powerflex_storagepool" "sp27" {
	media_type= "Transitional"
	protection_domain_name = "domain1"
	name = "ses-storage-pool30"
}
`

var CreateStoragePool14 = `
resource "powerflex_storagepool" "sp28" {
	media_type= "SSD"
	protection_domain_name = "domain1"
	name = "ses-storage-pool31"
	use_rmcache = true
	use_rfcache = false
}
`

var ModifyMediaTypeSSDToTransitional = `
resource "powerflex_storagepool" "sp28" {
	media_type= "Transitional"
	protection_domain_name = "domain1"
	name = "ses-storage-pool31"
	use_rmcache = true
	use_rfcache = false
}
`

// var CreateStoragePool15 = `
// resource "powerflex_storagepool" "sp29" {
// 	media_type= "SSD"
// 	protection_domain_name = "domain1"
// 	name = "ses-storage-pool32"
// 	use_rmcache = true
// 	use_rfcache = false
// }
// `

// var AttachVolumeToSP = `
// resource "powerflex_volume" "avengers" {
// 	name = "abcde"
// 	storage_pool_name = "ses-storage-pool32"
// 	protection_domain_name = "domain1"
// 	capacity_unit = "GB"
// 	size = 8
//   }
// `