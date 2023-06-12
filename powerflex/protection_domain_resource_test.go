/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package powerflex

import (
	"context"
	"regexp"
	"testing"

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccPDResource(t *testing.T) {
	createPDTest := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
	}
	`
	// too many characters in the name
	updatePDNameTestNeg := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd**********%^$()@!#$%~5555555555555555555555555555555555555555555555"
	}
	`

	updatePDFGLM1 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		fgl_metadata_cache_enabled = true
		fgl_default_metadata_cache_size = 1024
	}
	`
	updatePDFGLM2 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		fgl_metadata_cache_enabled = false
	}
	`
	updatePDFGLM3 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		fgl_metadata_cache_enabled = true
		fgl_default_metadata_cache_size = 5*1024
	}
	`

	// fgl_default_metadata_cache_size must be multiple of 1024
	updatePDFGLMNeg := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		fgl_metadata_cache_enabled = true
		fgl_default_metadata_cache_size = 1025
	}
	`
	// fgl_default_metadata_cache_size cannot be set when fgl_metadata_cache_enabled is false
	updatePDFGLMNeg2 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		fgl_metadata_cache_enabled = false
		fgl_default_metadata_cache_size = 1024
	}
	`

	updatePDdeactivate := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		active = false
	}
	`

	updatePDreactivate := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		active = true
	}
	`

	updatePDIopsTest := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		protected_maintenance_mode_network_throttling_in_kbps = 10*1024
		rebuild_network_throttling_in_kbps = 10*1024
		rebalance_network_throttling_in_kbps = 10*1024
		vtree_migration_network_throttling_in_kbps = 10*1024
	}
	`

	// we are setting all throttling as unlimited while overall_io_network_throttling_in_kbps is limited
	updatePDIopsTestNeg1 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		overall_io_network_throttling_in_kbps = 100*1024
	}
	`

	// we are setting vtree_migration_network_throttling_in_kbps less than 5MB
	updatePDIopsTestNeg2 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		protected_maintenance_mode_network_throttling_in_kbps = 10*1024
		rebuild_network_throttling_in_kbps = 10*1024
		rebalance_network_throttling_in_kbps = 10*1024
		vtree_migration_network_throttling_in_kbps = 4*1024
		overall_io_network_throttling_in_kbps = 100*1024
	}
	`

	updatePDRfCacheTest := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		rf_cache_enabled = true
		rf_cache_operational_mode = "ReadAndWrite"
		rf_cache_page_size_kb = 16
		rf_cache_max_io_size_kb = 32
	}
	`
	// rf cache sizes need to be powers of 2
	updatePDRfCacheTestNeg1 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		rf_cache_enabled = true
		rf_cache_operational_mode = "ReadAndWrite"
		rf_cache_page_size_kb = 10
		rf_cache_max_io_size_kb = 10
	}
	`

	// rf cache sizes need to be powers of 2 within ranges
	updatePDRfCacheTestNeg2 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		rf_cache_enabled = true
		rf_cache_operational_mode = "ReadAndWrite"
		rf_cache_page_size_kb = 512
		rf_cache_max_io_size_kb = 512
	}
	`

	// rf cache params cannot be set when rf cache is not enabled
	updatePDRfCacheDisableTestNeg := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		rf_cache_enabled = false
		rf_cache_operational_mode = "ReadAndWrite"
		rf_cache_page_size_kb = 16
		rf_cache_max_io_size_kb = 32
	}
	`

	updatePDRfCacheDisableTest := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		rf_cache_enabled = false
	}
	`

	resourceName := "powerflex_protection_domain.pd"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config: ProviderConfigForTesting + createPDTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rebuild_network_throttling_in_kbps", "0"),
					resource.TestCheckResourceAttr(resourceName, "rebalance_network_throttling_in_kbps", "0"),
					resource.TestCheckResourceAttr(resourceName, "vtree_migration_network_throttling_in_kbps", "0"),
					resource.TestCheckResourceAttr(resourceName, "overall_io_network_throttling_in_kbps", "0"),
					resource.TestCheckResourceAttr(resourceName, "protected_maintenance_mode_network_throttling_in_kbps", "0"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updatePDNameTestNeg,
				ExpectError: regexp.MustCompile(".*name exceeds the allowed length of 31 character.*"),
			},
			{
				// this does not guarantee that fgl metadata was not true from before
				Config: ProviderConfigForTesting + updatePDFGLM1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "fgl_metadata_cache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "fgl_default_metadata_cache_size", "1024"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ProviderConfigForTesting + updatePDFGLM2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "fgl_metadata_cache_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "fgl_default_metadata_cache_size", "1024"),
				),
			},
			{
				// this is just to check that update of caching from false to true works
				Config: ProviderConfigForTesting + updatePDFGLM3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "fgl_metadata_cache_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "fgl_default_metadata_cache_size", "5120"),
				),
			},
			{
				Config:      ProviderConfigForTesting + updatePDFGLMNeg,
				ExpectError: regexp.MustCompile(".*The metadata cache size is out of range, or uses a wrong granularity.*"),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      ProviderConfigForTesting + updatePDFGLMNeg2,
				ExpectError: regexp.MustCompile(".*can be set only when fgl_metadata_cache_enabled is set to true.*"),
			},
			{
				Config: ProviderConfigForTesting + updatePDdeactivate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "state", "Inactive"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ProviderConfigForTesting + updatePDreactivate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "state", "Active"),
				),
			},
			{
				Config: ProviderConfigForTesting + updatePDIopsTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "rebuild_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rebalance_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "vtree_migration_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "overall_io_network_throttling_in_kbps", "0"),
					resource.TestCheckResourceAttr(resourceName, "protected_maintenance_mode_network_throttling_in_kbps", "10240"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      ProviderConfigForTesting + updatePDIopsTestNeg1,
				ExpectError: regexp.MustCompile(".*must be set to a value less than.*"),
			},
			{
				Config:      ProviderConfigForTesting + updatePDIopsTestNeg2,
				ExpectError: regexp.MustCompile(".*Each limit must be more than 5 MB.*"),
			},
			{
				Config: ProviderConfigForTesting + updatePDRfCacheTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "rebuild_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rebalance_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_page_size_kb", "16"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_max_io_size_kb", "32"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_enabled", "true"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      ProviderConfigForTesting + updatePDRfCacheTestNeg1,
				ExpectError: regexp.MustCompile(".*Invalid RFcache page size. Valid values are powers of 2 between 4KB and 64KB.*"),
			},
			{
				Config:      ProviderConfigForTesting + updatePDRfCacheTestNeg2,
				ExpectError: regexp.MustCompile(".*Invalid RFcache page size. Valid values are powers of 2 between 4KB and 64KB.*"),
			},
			{
				Config:      ProviderConfigForTesting + updatePDRfCacheDisableTestNeg,
				ExpectError: regexp.MustCompile(".*can be set only when rf_cache_enabled is set to true.*"),
			},
			{
				Config: ProviderConfigForTesting + updatePDRfCacheDisableTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "rebuild_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rebalance_network_throttling_in_kbps", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_page_size_kb", "16"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_max_io_size_kb", "32"),
					resource.TestCheckResourceAttr(resourceName, "rf_cache_enabled", "false"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDResource2(t *testing.T) {
	// name contains the illegal character &
	createPDTestNeg := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd**********%^$()@!#$%~5555555555555555555555555555555555555555555555"
	}
	`

	validatePDEmptyNameNeg := `
	resource "powerflex_protection_domain" "pd" {
		name = ""
	}
	`

	// a failure on update after create
	createPDTestNeg2 := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		active = false
		rf_cache_enabled = false
		rf_cache_operational_mode = "ReadAndWrite"
		rf_cache_page_size_kb = 5
		rf_cache_max_io_size_kb = 32
	}
	`
	createPDdeactive := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_1"
		active = false
	}
	`
	updatePDname := `
	resource "powerflex_protection_domain" "pd" {
		name = "test_acc_pd_2"
		active = false
	}
	`

	resourceName := "powerflex_protection_domain.pd"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + validatePDEmptyNameNeg,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(".*Attribute name string length must be at least 1.*"),
			},
			{
				Config:      ProviderConfigForTesting + createPDTestNeg,
				ExpectError: regexp.MustCompile(".*name exceeds the allowed length of 31 character.*"),
			},
			{
				Config:      ProviderConfigForTesting + createPDTestNeg2,
				ExpectError: regexp.MustCompile(".*can be set only when rf_cache_enabled is set to true.*"),
			},
			{
				Config: ProviderConfigForTesting + createPDdeactive,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_1"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "state", "Inactive"),
				),
			},
			{
				Config: ProviderConfigForTesting + updatePDname,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_acc_pd_2"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "state", "Inactive"),
				),
			},
			// check that import is working
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  ProviderConfigForTesting + updatePDname,
				Destroy: true,
			},
		},
	})
}

func TestApiLinks(t *testing.T) {
	tlinks := []*scaleiotypes.Link{
		{
			Rel:  "self",
			HREF: "/api/delete",
		},
	}
	links := getLinkTfList(tlinks)
	a, d := getLinksFromTfList(context.TODO(), links)
	assert.Equal(t, false, d.HasError())
	assert.NotEmpty(t, a)
	assert.Equal(t, tlinks, a)
}
