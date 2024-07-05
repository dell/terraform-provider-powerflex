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
	"regexp"
	"testing"

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

var addStandby = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
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
		},
		{
			ips = ["` + MDMDataPoints.standByIP2 + `"]
			role = "TieBreaker"
		},
	]
}
`

var removeStandBy = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "ThreeNodes"
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
	standby_mdm = []
}
`

var expandCluster = `
resource "powerflex_mdm_cluster" "test" {
	cluster_mode = "FiveNodes"
	primary_mdm = {
		id = "` + MDMDataPoints.primaryMDMID + `"
	}
	secondary_mdm = [
		{
			id = "` + MDMDataPoints.secondaryMDMID + `"
		},
		{
			ips = ["` + MDMDataPoints.standByIP1 + `"]
		}
	]
	tiebreaker_mdm = [
		{
			id = "` + MDMDataPoints.tbID + `"
		},
		{
			ips = ["` + MDMDataPoints.standByIP2 + `"]
		}
	]
	standby_mdm = []
}
`

var reduceCluster = addStandby

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

func TestAccMdmClusterResource(t *testing.T) {
	var mdmClusterResourceBlock = "powerflex_mdm_cluster.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + createMdmConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "secondary_mdm.#", "1"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "tiebreaker_mdm.#", "1"),
				),
			},
			{
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
			{
				Config: ProviderConfigForTesting + renameMdmConfig1,
			},
			{
				Config:      ProviderConfigForTesting + renameMdmConfigNegative,
				ExpectError: regexp.MustCompile("Could not rename the MDM"),
			},
		}})
}

func TestAccMdmClusterSwitchClusterMode(t *testing.T) {
	var mdmClusterResourceBlock = "powerflex_mdm_cluster.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + addStandby,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.secondaryMDMID,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "tiebreaker_mdm.*", map[string]string{
						"id": MDMDataPoints.tbID,
					}),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "standby_mdm.#", "2"),
				),
			},
			{
				Config: ProviderConfigForTesting + expandCluster,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.secondaryMDMID,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "tiebreaker_mdm.*", map[string]string{
						"id": MDMDataPoints.tbID,
					}),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "standby_mdm.#", "0"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "secondary_mdm.#", "2"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "tiebreaker_mdm.#", "2"),
				),
			},
			{
				Config: ProviderConfigForTesting + expandCluster,
			},
			{
				Config: ProviderConfigForTesting + reduceCluster,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.secondaryMDMID,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "tiebreaker_mdm.*", map[string]string{
						"id": MDMDataPoints.tbID,
					}),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "standby_mdm.#", "2"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "secondary_mdm.#", "1"),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "tiebreaker_mdm.#", "1"),
				),
			},
			{
				Config: ProviderConfigForTesting + removeStandBy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "primary_mdm.id", MDMDataPoints.primaryMDMID),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "secondary_mdm.*", map[string]string{
						"id": MDMDataPoints.secondaryMDMID,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(mdmClusterResourceBlock, "tiebreaker_mdm.*", map[string]string{
						"id": MDMDataPoints.tbID,
					}),
					resource.TestCheckResourceAttr(mdmClusterResourceBlock, "standby_mdm.#", "0"),
				),
			},
		}})
}

func TestAccMdmClusterSwitchPrimaryMdm(t *testing.T) {
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
