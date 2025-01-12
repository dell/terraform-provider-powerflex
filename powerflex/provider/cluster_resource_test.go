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
	"encoding/json"
	"fmt"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccResourceCluster tests the SDC Expansion Operation
func TestAccResourceClusterA(t *testing.T) {
	t.Skip("Skipping this test case")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForTesting + packageTest + ClusterConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					validateSDCLength,
				),
			},
		},
	})
}

func validateSDCLength(state *terraform.State) error {
	// Retrieve the resource instance
	clusterResource, ok := state.RootModule().Resources["powerflex_cluster.test"]
	if !ok {
		return fmt.Errorf("Failed to find powerflex_cluster.test in state")
	}

	// Get the value of the "sdc_list" attribute from the resource instance
	sdcListValue, ok := clusterResource.Primary.Attributes["sdc_list"]
	if !ok {
		return fmt.Errorf("sdc_list attribute not found in state")
	}

	// Parse the sdc_list value into a list
	var sdcList []interface{}
	if err := json.Unmarshal([]byte(sdcListValue), &sdcList); err != nil {
		return fmt.Errorf("Failed to unmarshal sdc_list attribute: %s", err)
	}

	// Check if the length of the sdc_list is greater than 0
	if len(sdcList) <= 0 {
		return fmt.Errorf("sdc_list attribute length is not greater than 0")
	}

	return nil
}

func TestAccResourceClusterValidation(t *testing.T) {
	var FunctionMockerClusterCreate *Mocker
	var FunctionMockerClusterMdmIp *Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ParseClusterCSVOperation).Return(&goscaleio_types.GatewayResponse{
						Message: "",
						Data:    "",
					}, nil).Build()
					FunctionMockerClusterCreate = Mock(helper.ClusterInstallationOperations).Return(nil).Build()
					FunctionMockerClusterMdmIp = Mock(helper.GetMDMIPFromClusterDetails).Return("tfacc_cluster_ip_1", nil).Build()
				},
				Config: ProviderConfigForTesting + ClusterConfigValidator,
				Check: resource.ComposeAggregateTestCheckFunc(
					validateSDCLength,
				),
				ExpectError: regexp.MustCompile(`.*strconv.ParseUint: parsing*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerClusterCreate != nil {
						FunctionMockerClusterCreate.UnPatch()
					}
					if FunctionMockerClusterMdmIp != nil {
						FunctionMockerClusterMdmIp.UnPatch()
					}
					FunctionMocker = Mock(helper.ResetInstallerQueue).Return(fmt.Errorf("Mock Error")).Build()
				},
				Config: ProviderConfigForTesting + ClusterConfigValidator,
				Check: resource.ComposeAggregateTestCheckFunc(
					validateSDCLength,
				),
				ExpectError: regexp.MustCompile(`.*Error Clearing Queue*`),
			},
			//Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerClusterCreate != nil {
						FunctionMockerClusterCreate.UnPatch()
					}
					if FunctionMockerClusterMdmIp != nil {
						FunctionMockerClusterMdmIp.UnPatch()
					}
					FunctionMocker = Mock(helper.ParseClusterCSVOperation).Return(&goscaleio_types.GatewayResponse{
						Message: "",
						Data:    "",
					}, nil).Build()
					FunctionMockerClusterCreate = Mock(helper.ClusterInstallationOperations).Return(fmt.Errorf("Mock Error")).Build()
					FunctionMockerClusterMdmIp = Mock(helper.GetMDMIPFromClusterDetails).Return("tfacc_cluster_ip_1", nil).Build()
				},
				Config: ProviderConfigForTesting + ClusterConfigValidator,
				Check: resource.ComposeAggregateTestCheckFunc(
					validateSDCLength,
				),
				ExpectError: regexp.MustCompile(`.*Error in Installation Process*`),
			},
			//Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerClusterCreate != nil {
						FunctionMockerClusterCreate.UnPatch()
					}
					if FunctionMockerClusterMdmIp != nil {
						FunctionMockerClusterMdmIp.UnPatch()
					}
					FunctionMocker = Mock(helper.ParseClusterCSVOperation).Return(&goscaleio_types.GatewayResponse{
						Message: "",
						Data:    "",
					}, nil).Build()
					FunctionMockerClusterCreate = Mock(helper.ClusterInstallationOperations).Return(nil).Build()
					FunctionMockerClusterMdmIp = Mock(helper.GetMDMIPFromClusterDetails).Return(nil, fmt.Errorf("Mock Error")).Build()
				},
				Config: ProviderConfigForTesting + ClusterConfigValidator,
				Check: resource.ComposeAggregateTestCheckFunc(
					validateSDCLength,
				),
				ExpectError: regexp.MustCompile(`.*Error in Fecthing Primary MDM IP Before Installation*`),
			},
			//Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if FunctionMockerClusterCreate != nil {
						FunctionMockerClusterCreate.UnPatch()
					}
					if FunctionMockerClusterMdmIp != nil {
						FunctionMockerClusterMdmIp.UnPatch()
					}
				},
				Config:      ProviderConfigForTesting + ClusterValidationConfig1,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value*`),
			},
			//Create
			{
				Config:      ProviderConfigForTesting + ClusterValidationConfig2,
				ExpectError: regexp.MustCompile(`.*Missing required argument*`),
			},
			//Create
			{
				Config:      ProviderConfigForTesting + ClusterValidationConfig3,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value*`),
			},
			//Create
			{
				Config:      ProviderConfigForTesting + ClusterValidationConfig4,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Length*`),
			},
			//Create
			{
				Config:      ProviderConfigForTesting + ClusterValidationConfig5,
				ExpectError: regexp.MustCompile(`.*Please configure replication_journal_capacity_percentage for SDR*`),
			},
			//Create
			{
				Config:      ProviderConfigForTesting + ClusterValidationConfig6,
				ExpectError: regexp.MustCompile(`.*Error During Installation*`),
			},
			//Import
			{
				Config:        ProviderConfigForTesting + importClusterTest,
				ImportState:   true,
				ImportStateId: "1.1.1.1,Password",
				ResourceName:  "powerflex_cluster.test",
				ExpectError:   regexp.MustCompile(`.*Please provide valid Input Details*`),
			},
			//Import
			{
				Config:        ProviderConfigForTesting + importClusterTest,
				ImportState:   true,
				ImportStateId: "1.1.1.1,Password,Password",
				ResourceName:  "powerflex_cluster.test",
				ExpectError:   regexp.MustCompile(`.*Error Getting Cluster Details*`),
			},
		},
	})
}

var packageTest = `
resource "powerflex_package" "upload-test" {
	file_path = ["../resource-test/powerflex_packages/EMC-ScaleIO-lia-3.6-700.103.Ubuntu.22.04.x86_64.tar"]
 }
`

var importClusterTest = `
resource "powerflex_cluster" "test"  {
	
	}
`

var ClusterConfigValidator = `
resource "powerflex_cluster" "test" {

	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	disable_non_mgmt_components_auth= false
	cluster = [
	{
		ips= "` + GatewayDataPoints.clusterPrimaryIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "HighPerformance",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterSecondaryIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterTBIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		is_sds= "No",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`

var ClusterConfig1 = `
resource "powerflex_cluster" "test" {

	depends_on = [
		powerflex_package.upload-test
	]

	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	disable_non_mgmt_components_auth= false
	cluster = [
	{
		ips= "` + GatewayDataPoints.clusterPrimaryIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "HighPerformance",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterSecondaryIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterTBIP + `",
		username= "root",
		password = "` + GatewayDataPoints.serverPassword + `"
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		is_sds= "No",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`

var ClusterValidationConfig6 = `
resource "powerflex_cluster" "test" {

	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = [
	{
		ips= "` + GatewayDataPoints.clusterPrimaryIP + `",
		username= "root",
		password= "` + GatewayDataPoints.serverPassword + `",
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterSecondaryIP + `",
		username= "root",
		password= "` + GatewayDataPoints.serverPassword + `",
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "` + GatewayDataPoints.clusterTBIP + `",
		username= "root",
		password= "` + GatewayDataPoints.serverPassword + `",
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds3",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`

var ClusterValidationConfig1 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = []
	storage_pools = []
}
`

var ClusterValidationConfig2 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
}
`

var ClusterValidationConfig3 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = [
	{
		ips= "1.1.1.1",
		username= "root",
		password= "Password",
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
		sdr_all_ips= ""
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`

var ClusterValidationConfig4 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = [
	{
		ips= "10.10.10.10",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		perf_profile_for_mdm= "HighPerformance",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "HighPerformance",
		ia_rfcache= "No",
		is_sdr= "No",
		sdr_all_ips= ""
	 },
	 {
		ips= "10.10.10.11",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
		sdr_all_ips= ""
	 },
	 {
		ips= "10.10.10.12",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds3",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
		sdr_all_ips= ""
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`

var ClusterValidationConfig5 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = [
	{
		ips= "10.10.10.10",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "primary",
		perf_profile_for_mdm= "ABCD",
		is_sds= "yes",
		sds_name= "sds1",
		is_sdc= "yes",
		sdc_name= "sdc1",
		perf_profile_for_sdc= "ABCD",
		ia_rfcache= "No",
		is_sdr= "Yes",
		sdr_all_ips= "1.1.1.1"
	 },
	 {
		ips= "10.10.10.11",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "Secondary",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds2",
		is_sdc= "yes",
		sdc_name= "sdc2",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	 {
		ips= "10.10.10.12",
		username= "root",
		password= "dangerous",
		operating_system= "linux",
		is_mdm_or_tb= "TB",
		perf_profile_for_mdm= "compact",
		is_sds= "yes",
		sds_name= "sds3",
		is_sdc= "yes",
		sdc_name= "sdc3",
		perf_profile_for_sdc= "compact",
		ia_rfcache= "No",
		is_sdr= "No",
	 },
	]
	storage_pools = [
		{
			media_type = "HDD"
		}	
	]
}
`
