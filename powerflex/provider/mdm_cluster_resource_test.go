/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var createMdmConfig = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var createMdmConfigUsingIps = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  ips = ["38.0.101.76"]
	}
	secondary_mdm = [
		{
			ips = ["38.0.101.75"]
		},
	]
	tiebreaker_mdm = [
		{
			ips = ["38.0.101.74"]
		},
	]
}
`

var updateOwnerMdmErrorConfig = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  id = "some-other-id"
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var createMdmConfigSwitchClusterMode = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "FiveNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var createMdmInvalidStandbyWithPortConfig = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
	standby_mdm = [
		{
			ips = ["` + MDMDataPoints.standByIP1 + `"]
			role = "Manager"
			port = 1234
		},
		{
			ips = ["` + MDMDataPoints.standByIP2 + `"]
			role = "TieBreaker"
			port = 1235
		},
	]
}
`

var createMdmConfigInvalidMdmIp = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  ips = ["invilid_ip"]
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var createMdmConfigInvalidSecondaryMdmIp = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
		id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
			ips = ["invilid_ip"]
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var createMdmConfigInvalidTiebreakerMdmIp = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
		id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
			id = "` + MDMDataPoints.secondaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
			ips = ["invilid_ip"]
		},
	]
}
`

var renameMdmConfig1 = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
		id = "` + MDMDataPoints.primaryMDMID + `"
	    name = "primary_mdm_renamed"
	}
	secondary_mdm = [
		{
			id = "` + MDMDataPoints.secondaryMDMID + `"
			name = "secondary_mdm_renamed"
		},
	]
	tiebreaker_mdm = [
		{
			id = "` + MDMDataPoints.tbID + `"
			name = "tb_mdm_renamed"
		},
	]
}
`

var renameMdmConfigNegative = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
		id = "` + MDMDataPoints.primaryMDMID + `"
	    name = "tb_mdm_renamed"
	}
	secondary_mdm = [
		{
			id = "` + MDMDataPoints.secondaryMDMID + `"
			name = "secondary_mdm_renamed"
		},
	]
	tiebreaker_mdm = [
		{
			id = "` + MDMDataPoints.tbID + `"
			name = "tb_mdm_renamed"
		},
	]
}
`

var renameMdmConfig2 = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "HighPerformance"
	primary_mdm = {
	  id = "` + MDMDataPoints.primaryMDMID + `"
	  name = "primary_mdm_renamed1"
	}
	secondary_mdm = [
		{
			id = "` + MDMDataPoints.secondaryMDMID + `"
			name = "secondary_mdm_renamed1"
		},
	]
	tiebreaker_mdm = [
		{
			id = "` + MDMDataPoints.tbID + `"
			name = "tb_mdm_renamed1"
		},
	]
}
`

var switchPrimaryMdm1 = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
	performance_profile = "Compact"
	primary_mdm = {
	  id = "` + MDMDataPoints.secondaryMDMID + `"
	}
	secondary_mdm = [
		{
	  		id = "` + MDMDataPoints.primaryMDMID + `"
		},
	]
	tiebreaker_mdm = [
		{
	  		id = "` + MDMDataPoints.tbID + `"
		},
	]
}
`

var switchPrimaryMdm2 = createMdmConfig

// UT
func TestAccResourceAcceptanceMdmCluster(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a unit test")
	}
	var mdmClusterResourceBlock = "powerflex_mdm_cluster.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 1 Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createMdmConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// 2Get MDM Cluster Details Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetMDMClusterDetails).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createMdmConfig,
				ExpectError: regexp.MustCompile(`.*Error getting MDM cluster details*.`),
			},
			// 3 Invalid Primary MDM Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + createMdmConfigInvalidMdmIp,
				ExpectError: regexp.MustCompile(`.*Please enter valid IP for primary MDM*.`),
			},
			// 4 Invalid Secondary MDM Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + createMdmConfigInvalidSecondaryMdmIp,
				ExpectError: regexp.MustCompile(`.*Please enter valid IP for secondary MDM*.`),
			},
			// 5 Invalid Tiebreaker MDM Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + createMdmConfigInvalidTiebreakerMdmIp,
				ExpectError: regexp.MustCompile(`.*Please enter valid IP for tiebreaker MDM*.`),
			},
			// 6 Create Success
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + createMdmConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "secondary_mdm.#", "1"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "tiebreaker_mdm.#", "1"),
				),
			},
			// 7 Change Owner Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).ChangeMdmOwnerShip).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + updateOwnerMdmErrorConfig,
				ExpectError: regexp.MustCompile(`.*Could not change MDM ownership with ID*.`),
			},
			// 8 Get System Error Update
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createMdmConfig,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// 9 Add Standby Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).AddStandByMdm).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createMdmInvalidStandbyWithPortConfig,
				ExpectError: regexp.MustCompile(`.*Could not add standby MDM with IP*.`),
			},
			// 10 Switch Cluster Mode Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).SwitchClusterMode).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + createMdmConfigSwitchClusterMode,
				ExpectError: regexp.MustCompile(`.*Could not expand the MDM cluster*.`),
			},
			// 11 Rename Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).RenameMdm).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + renameMdmConfig2,
				ExpectError: regexp.MustCompile(`.*Could not rename the MDM with ID*.`),
			},
			// 12 Rename
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + renameMdmConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "performance_profile", "HighPerformance"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.name", "primary_mdm_renamed1"),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id":   MDMDataPoints.secondaryMDMID,
						"name": "secondary_mdm_renamed1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "tiebreaker_mdm.*", map[string]string{
						"id":   MDMDataPoints.tbID,
						"name": "tb_mdm_renamed1",
					}),
				),
			},
			// 13 reset
			{
				Config: ProviderConfigForTesting + renameMdmConfig1,
			},
			// 14 Invalid Config
			{
				Config:      ProviderConfigForTesting + renameMdmConfigNegative,
				ExpectError: regexp.MustCompile("Could not rename the MDM"),
			},
		}})
}

// UT Second way of creating using ips
func TestAccResourceAcceptanceMdmClusterIps(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a unit test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create using ips
			{
				Config: ProviderConfigForTesting + createMdmConfigUsingIps,
			},
		}})
}

// AT
func TestAccResourceMdmCluster(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with Unit test, this is an acceptance test")
	}
	var mdmClusterResourceBlock = "powerflex_mdm_cluster.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + switchPrimaryMdm1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.secondaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.primaryMDMID,
					}),
				),
			},
			{
				Config: ProviderConfigForTesting + switchPrimaryMdm2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.secondaryMDMID,
					}),
				),
			},
		}})
}
