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
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

var ProviderConfigForTesting = ``

type sdsDataPoints struct {
	SdsIP1   string
	SdsIP2   string
	SdsIP3   string
	SdsIP4   string
	SdsIP5   string
	SdsIP6   string
	SdsIP7   string
	SdsIP8   string
	SdsIP9   string
	SdsIP10  string
	SdsIP11  string
	SdcIP    string
	SdcIP1   string
	volName  string
	volName2 string
	volName3 string
	sdcName  string
	sdcName2 string
	sdcName3 string
}

type gatewayDataPoints struct {
	primaryMDMIP       string
	secondaryMDMIP     string
	tbIP               string
	serverPassword     string
	mdmPassword        string
	liaPassword        string
	sdcServerIP        string
	clusterPrimaryIP   string
	clusterSecondaryIP string
	clusterTBIP        string
}

type mdmDataPoints struct {
	primaryMDMIP   string
	secondaryMDMIP string
	tbIP           string
	primaryMDMID   string
	secondaryMDMID string
	tbID           string
	standByIP1     string
	standByIP2     string
}

func getNewSdsDataPointForTest() sdsDataPoints {
	var SdsResourceTestData sdsDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return SdsResourceTestData
	}
	SdsResourceTestData.SdsIP1 = setDefault(os.Getenv("POWERFLEX_SDS_IP_1"), "tfacc_sds_1")
	SdsResourceTestData.SdsIP2 = setDefault(os.Getenv("POWERFLEX_SDS_IP_2"), "tfacc_sds_2")
	SdsResourceTestData.SdsIP3 = setDefault(os.Getenv("POWERFLEX_SDS_IP_3"), "tfacc_sds_3")
	SdsResourceTestData.SdsIP4 = setDefault(os.Getenv("POWERFLEX_SDS_IP_4"), "tfacc_sds_4")
	SdsResourceTestData.SdsIP5 = setDefault(os.Getenv("POWERFLEX_SDS_IP_5"), "tfacc_sds_5")
	SdsResourceTestData.SdsIP6 = setDefault(os.Getenv("POWERFLEX_SDS_IP_6"), "tfacc_sds_6")
	SdsResourceTestData.SdsIP7 = setDefault(os.Getenv("POWERFLEX_SDS_IP_7"), "tfacc_sds_7")
	SdsResourceTestData.SdsIP8 = setDefault(os.Getenv("POWERFLEX_SDS_IP_8"), "tfacc_sds_8")
	SdsResourceTestData.SdsIP9 = setDefault(os.Getenv("POWERFLEX_SDS_IP_9"), "tfacc_sds_9")
	SdsResourceTestData.SdsIP10 = setDefault(os.Getenv("POWERFLEX_SDS_IP_10"), "tfacc_sds_10")
	SdsResourceTestData.SdsIP11 = setDefault(os.Getenv("POWERFLEX_SDS_IP_11"), "tfacc_sds_11")
	SdsResourceTestData.SdcIP = setDefault(os.Getenv("POWERFLEX_SDC_IP"), "tfacc_sdc_ip_1")
	SdsResourceTestData.SdcIP1 = setDefault(os.Getenv("POWERFLEX_SDC_IP1"), "tfacc_sdc_ip_2")
	SdsResourceTestData.volName = setDefault(os.Getenv("POWERFLEX_VOLUME_NAME"), "tfacc_volume_1")
	SdsResourceTestData.volName2 = setDefault(os.Getenv("POWERFLEX_VOLUME_NAME_2"), "tfacc_volume_2")
	SdsResourceTestData.volName3 = setDefault(os.Getenv("POWERFLEX_VOLUME_NAME_3"), "tfacc_volume_3")
	SdsResourceTestData.sdcName = setDefault(os.Getenv("POWERFLEX_SDC_NAME"), "tfacc_sdc_name_1")
	SdsResourceTestData.sdcName2 = setDefault(os.Getenv("POWERFLEX_SDC_NAME_2"), "tfacc_sdc_name_2")
	SdsResourceTestData.sdcName3 = setDefault(os.Getenv("POWERFLEX_SDC_NAME_3"), "tfacc_sdc_name_3")

	return SdsResourceTestData
}

func getNewGatewayDataPointForTest() gatewayDataPoints {

	var GatewayDataPoints gatewayDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return GatewayDataPoints
	}

	GatewayDataPoints.primaryMDMIP = setDefault(os.Getenv("POWERFLEX_PRIMARY_MDM_IP"), "tfacc_primary_mdm_ip")
	GatewayDataPoints.secondaryMDMIP = setDefault(os.Getenv("POWERFLEX_SECONDARY_MDM_IP"), "tfacc_secondary_mdm_ip")
	GatewayDataPoints.tbIP = setDefault(os.Getenv("POWERFLEX_TB_IP"), "tfacc_tb_ip")
	GatewayDataPoints.sdcServerIP = setDefault(os.Getenv("POWERFLEX_SDC_SERVER_IP"), "tfacc_sdc_server_ip")
	GatewayDataPoints.serverPassword = setDefault(os.Getenv("POWERFLEX_SERVER_PASSWORD"), "tfacc_server_password")
	GatewayDataPoints.mdmPassword = setDefault(os.Getenv("POWERFLEX_MDM_PASSWORD"), "tfacc_mdm_password")
	GatewayDataPoints.liaPassword = setDefault(os.Getenv("POWERFLEX_LIA_PASSWORD"), "tfacc_lia_password")
	GatewayDataPoints.clusterPrimaryIP = setDefault(os.Getenv("POWERFLEX_CLUSTER_IP_1"), "tfacc_cluster_ip_1")
	GatewayDataPoints.clusterSecondaryIP = setDefault(os.Getenv("POWERFLEX_CLUSTER_IP_2"), "tfacc_cluster_ip_2")
	GatewayDataPoints.clusterTBIP = setDefault(os.Getenv("POWERFLEX_CLUSTER_IP_3"), "tfacc_cluster_ip_3")

	return GatewayDataPoints
}

func getMdmDataPointsForTest() mdmDataPoints {

	var MDMDataPoints mdmDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return MDMDataPoints
	}

	MDMDataPoints.primaryMDMIP = setDefault(os.Getenv("POWERFLEX_PRIMARY_MDM_IP"), "tfacc_primary_mdm_ip")
	MDMDataPoints.secondaryMDMIP = setDefault(os.Getenv("POWERFLEX_SECONDARY_MDM_IP"), "tfacc_secondary_mdm_ip")
	MDMDataPoints.tbIP = setDefault(os.Getenv("POWERFLEX_TB_IP"), "tfacc_tb_ip")
	MDMDataPoints.primaryMDMID = setDefault(os.Getenv("POWERFLEX_PRIMARY_MDM_ID"), "tfacc_primary_mdm_id")
	MDMDataPoints.secondaryMDMID = setDefault(os.Getenv("POWERFLEX_SECONDARY_MDM_ID"), "tfacc_secondary_mdm_id")
	MDMDataPoints.tbID = setDefault(os.Getenv("POWERFLEX_TB_ID"), "tfacc_tb_id")
	MDMDataPoints.standByIP1 = setDefault(os.Getenv("POWERFLEX_STANDBY_MDM_IP1"), "1.1.1.1")
	MDMDataPoints.standByIP2 = setDefault(os.Getenv("POWERFLEX_STANDBY_MDM_IP2"), "1.1.1.2")

	return MDMDataPoints
}

var SdsResourceTestData = getNewSdsDataPointForTest()
var GatewayDataPoints = getNewGatewayDataPointForTest()
var SDCMappingResourceID2 = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_ID2"), "tfacc_sdc_volumes_mapping_id2")
var SDCMappingResourceName2 = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME2"), "tfacc_sdc_volumes_mapping_name2")
var SDCVolName = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME"), "tfacc_sdc_volumes_mapping_name")
var SdsID = setDefault(os.Getenv("POWERFLEX_DEVICE_SDS_ID"), "tfacc_device_sds_id")
var MDMDataPoints = getMdmDataPointsForTest()
var StoragePoolName = setDefault(os.Getenv("POWERFLEX_STORAGE_POOL_NAME"), "tfacc_storage_pool_name")
var ProtectionDomainID = setDefault(os.Getenv("POWERFLEX_PROTECTION_DOMAIN_ID"), "tfacc_protection_domain_id")
var Volume1 = setDefault(os.Getenv("POWERFLEX_VOLUME1"), "tfacc_volume_1")
var Volume2 = setDefault(os.Getenv("POWERFLEX_VOLUME2"), "tfacc_volume_2")
var Volume3 = setDefault(os.Getenv("POWERFLEX_VOLUME3"), "tfacc_volume_3")
var username = setDefault(os.Getenv("POWERFLEX_USERNAME"), "test")
var password = setDefault(os.Getenv("POWERFLEX_PASSWORD"), "test")
var endpoint = setDefault(os.Getenv("POWERFLEX_ENDPOINT"), "http://localhost:3002")

func init() {
	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	ProviderConfigForTesting = fmt.Sprintf(`
		provider "powerflex" {
			username = "%s"
			password = "%s"
			endpoint = "%s"
		}
	`, username, password, endpoint)
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerflex": providerserver.NewProtocol6WithError(New()),
	}
)

func testAccPreCheck(t *testing.T) {
	if v := username; v == "" {
		t.Fatal("POWERFLEX_USERNAME must be set for acceptance tests")
	}

	if v := password; v == "" {
		t.Fatal("POWERFLEX_PASSWORD must be set for acceptance tests")
	}

	if v := endpoint; v == "" {
		t.Fatal("POWERFLEX_ENDPOINT must be set for acceptance tests")
	}
}

// if there is no os setting set, then use the default value
func setDefault(osInput string, defaultStr string) string {
	if osInput == "" {
		return defaultStr
	}
	return osInput
}
