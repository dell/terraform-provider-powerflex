/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package provider

import (
	"fmt"
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ProtectionDomainDataSourceAll = `
data "powerflex_protection_domain" "all" {}
`

var ProtectionDomainDataSourceFilterSingle = `
data "powerflex_protection_domain" "filter-single" {
  filter {
		name = ["domain1"]
  }
}
`

var ProtectionDomainDataSourceFilterMultiple = `
data "powerflex_protection_domain" "filter-multi" {
	filter {
		system_id = ["1250de83018c2d0f"]
		rf_cache_enabled = true
	}
}
`

// AT
func TestAccDatasourceAcceptanceProtectionDomain(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccDatasourceProtectionDomain(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is an Unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// All
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceAll,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Filter on a single value
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceFilterSingle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.filter-single", "protection_domains.0.name", "domain1"),
				),
			},
			// Filter on multiple values
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.filter-multi", "protection_domains.0.system_id", "1250de83018c2d0f"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.filter-multi", "protection_domains.0.rf_cache_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.filter-multi", "protection_domains.1.system_id", "1250de83018c2d0f"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.filter-multi", "protection_domains.1.rf_cache_enabled", "true"),
				),
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetProtectionDomains).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceAll,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex ProtectionDomain*.`),
			},
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceFilterMultiple,
				ExpectError: regexp.MustCompile(`.*Error in filtering protection domains*.`),
			},
		},
	})
}
