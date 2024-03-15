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
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
	"testing"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type pdDataPoints struct {
	name  string
	name2 string
	name3 string
	id    string
}

var protectiondomainTestData pdDataPoints = pdDataPoints{
	id:    ProtectionDomainID,
	name:  "domain1",
	name2: "domain2",
	name3: "domain_1",
}

var (
	ProtectionDomainDataSourceConfig1 string
	ProtectionDomainDataSourceConfig2 string
	ProtectionDomainDataSourceConfig3 string
	ProtectionDomainDataSourceConfig4 string
	ProtectionDomainDataSourceConfig5 string
	ProtectionDomainDataSourceConfig6 string
	ProtectionDomainDataSourceConfig7 string
	ProtectionDomainDataSourceConfig8 string
)

// TestAccProtectionDomainDataSource tests the protectiondomain data source
// where it fetches the protectiondomains based on protectiondomain id/name
// and if nothing is mentioned , then return all protectiondomains
func TestAccProtectionDomainDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving protection domain based on id
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.0.name", protectiondomainTestData.name),
				),
			},
			//retrieving protection domain based on name
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.0.name", protectiondomainTestData.name),
				),
			},
			//retrieving all the protection domains
			{
				Config: ProviderConfigForTesting + ProtectionDomainDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.0.name", protectiondomainTestData.name),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.1.name", protectiondomainTestData.name2),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.2.name", protectiondomainTestData.name3),
				),
			},
		},
	})
}

// TestNonNullPDConnInfo tests if non-null PDConnInfo are properly marshalled
func TestNonNullPDConnInfo(t *testing.T) {
	inputStr := "Dummy"
	input := scaleiotypes.PDConnInfo{
		ClientServerConnStatus: inputStr,
		DisconnectedClientID:   &inputStr,
		DisconnectedClientName: &inputStr,
		DisconnectedServerID:   &inputStr,
		DisconnectedServerName: &inputStr,
		DisconnectedServerIP:   &inputStr,
	}

	expectedOut := models.PdConnInfoModel{
		ClientServerConnStatus: types.StringValue(inputStr),
		DisconnectedClientID:   types.StringValue(inputStr),
		DisconnectedClientName: types.StringValue(inputStr),
		DisconnectedServerID:   types.StringValue(inputStr),
		DisconnectedServerName: types.StringValue(inputStr),
		DisconnectedServerIP:   types.StringValue(inputStr),
	}

	out := helper.PdConnInfoModelValue(input)

	if out != expectedOut {
		t.Fatalf("Error matching output and expected: %#v vs %#v", out, expectedOut)
	}

}

// TestNonNullReplicationCapacityMaxRatio tests that properl marshalling occurs when
// ReplicationCapacityMaxRatio field has non null value
func TestNonNullReplicationCapacityMaxRatio(t *testing.T) {
	inp := 10
	input := scaleiotypes.ProtectionDomain{
		ReplicationCapacityMaxRatio: &inp,
	}

	outList := helper.GetAllProtectionDomainState([]*scaleiotypes.ProtectionDomain{
		&input,
	})
	out := outList[0]
	if actual := out.ReplicationCapacityMaxRatio.ValueInt64(); actual != int64(inp) {
		t.Fatalf("Error matching output and expected: %#v vs %#v", actual, inp)
	}
}

func TestNegativeScenarios(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Combination.*`),
			},
			{
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceConfig5,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex ProtectionDomain by ID.*`),
			},
			{
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceConfig6,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex ProtectionDomain by name.*`),
			},
			{
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceConfig7,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex ProtectionDomain by ID.*`),
			},
			{
				Config:      ProviderConfigForTesting + ProtectionDomainDataSourceConfig8,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex ProtectionDomain by name.*`),
			},
		},
	})
}

func init() {
	// retrieve protection domain by id
	ProtectionDomainDataSourceConfig1 = `
	data "powerflex_protection_domain" "pd1" {						
		id = "` + protectiondomainTestData.id + `"
	}
	`
	// retrieve protection domain by name
	ProtectionDomainDataSourceConfig2 = `
	data "powerflex_protection_domain" "pd2" {			
		name = "domain1"
	}
	`
	// retrieve all protection domains
	ProtectionDomainDataSourceConfig3 = `
	data "powerflex_protection_domain" "pd3" {						
	}
	`

	// retrieve all protection domain by id and name
	ProtectionDomainDataSourceConfig4 = `
	data "powerflex_protection_domain" "pd4" {
		name = "domain1"
		id = "` + protectiondomainTestData.id + `"
	}
	`

	// retrieve protection domain by non existing id
	ProtectionDomainDataSourceConfig5 = `
	data "powerflex_protection_domain" "pd5" {
		id = "non_existing_id"
	}
	`

	// retrieve protection domain by non existing name
	ProtectionDomainDataSourceConfig6 = `
	data "powerflex_protection_domain" "pd6" {
		name = "non_existing_name"
	}
	`

	// retrieve protection domain by empty id
	ProtectionDomainDataSourceConfig7 = `
	data "powerflex_protection_domain" "pd7" {
		id = ""
	}
	`

	// retrieve protection domain by empty name
	ProtectionDomainDataSourceConfig8 = `
	data "powerflex_protection_domain" "pd8" {
		name = ""
	}
	`
}
