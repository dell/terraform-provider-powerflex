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

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccClusterResource tests the SDC Expansion Operation
func TestAccClusterResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			// {
			// 	Config:      ProviderConfigForTesting + ClusterConfig1,
			// 	ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			// },
			//Import
			{
				Config:        ProviderConfigForTesting + importClusterTest,
				ImportState:   true,
				ImportStateId: "10.247.103.161,Password123,Password123",
				ResourceName:  "powerflex_cluster.test",
			},
		},
	})
}

var importClusterTest = `
resource "powerflex_cluster" "test"  {
	
	}
`

var ClusterConfig1 = `
resource "powerflex_cluster" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	allow_non_secure_communication_with_lia= true
	allow_non_secure_communication_with_mdm= true
	cluster = [
	{
		ips= "10.247.103.161",
		username= "root",
		password= "dangerous",
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
	 {
		ips= "10.247.103.162",
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
		ips= "10.247.103.160",
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
