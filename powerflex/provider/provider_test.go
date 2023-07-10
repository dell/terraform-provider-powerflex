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
	primaryMDMIP   string
	secondaryMDMIP string
	tbIP           string
	serverPassword string
	mdmPassword    string
	liaPassword    string
	sdcServerIP    string
}

func getNewSdsDataPointForTest() sdsDataPoints {
	var SdsResourceTestData sdsDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return SdsResourceTestData
	}
	SdsResourceTestData.SdsIP1 = os.Getenv("POWERFLEX_SDS_IP_1")
	SdsResourceTestData.SdsIP2 = os.Getenv("POWERFLEX_SDS_IP_2")
	SdsResourceTestData.SdsIP3 = os.Getenv("POWERFLEX_SDS_IP_3")
	SdsResourceTestData.SdsIP4 = os.Getenv("POWERFLEX_SDS_IP_4")
	SdsResourceTestData.SdsIP5 = os.Getenv("POWERFLEX_SDS_IP_5")
	SdsResourceTestData.SdsIP6 = os.Getenv("POWERFLEX_SDS_IP_6")
	SdsResourceTestData.SdsIP7 = os.Getenv("POWERFLEX_SDS_IP_7")
	SdsResourceTestData.SdsIP8 = os.Getenv("POWERFLEX_SDS_IP_8")
	SdsResourceTestData.SdsIP9 = os.Getenv("POWERFLEX_SDS_IP_9")
	SdsResourceTestData.SdsIP10 = os.Getenv("POWERFLEX_SDS_IP_10")
	SdsResourceTestData.SdsIP11 = os.Getenv("POWERFLEX_SDS_IP_11")
	SdsResourceTestData.SdcIP = os.Getenv("POWERFLEX_SDC_IP")
	SdsResourceTestData.SdcIP1 = os.Getenv("POWERFLEX_SDC_IP1")
	SdsResourceTestData.volName = os.Getenv("POWERFLEX_VOLUME_NAME")
	SdsResourceTestData.volName2 = os.Getenv("POWERFLEX_VOLUME_NAME_2")
	SdsResourceTestData.volName3 = os.Getenv("POWERFLEX_VOLUME_NAME_3")
	SdsResourceTestData.sdcName = os.Getenv("POWERFLEX_SDC_NAME")
	SdsResourceTestData.sdcName2 = os.Getenv("POWERFLEX_SDC_NAME_2")
	SdsResourceTestData.sdcName3 = os.Getenv("POWERFLEX_SDC_NAME_3")

	return SdsResourceTestData
}

func getNewGatewayDataPointForTest() gatewayDataPoints {

	var GatewayDataPoints gatewayDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return GatewayDataPoints
	}

	GatewayDataPoints.primaryMDMIP = os.Getenv("POWERFLEX_PRIMARY_MDM_IP")
	GatewayDataPoints.secondaryMDMIP = os.Getenv("POWERFLEX_SECONDARY_MDM_IP")
	GatewayDataPoints.tbIP = os.Getenv("POWERFLEX_TB_IP")
	GatewayDataPoints.sdcServerIP = os.Getenv("POWERFLEX_SDC_SERVER_IP")
	GatewayDataPoints.serverPassword = os.Getenv("POWERFLEX_SERVER_PASSWORD")
	GatewayDataPoints.mdmPassword = os.Getenv("POWERFLEX_MDM_PASSWORD")
	GatewayDataPoints.liaPassword = os.Getenv("POWERFLEX_LIA_PASSWORD")

	return GatewayDataPoints
}

var SdsResourceTestData = getNewSdsDataPointForTest()
var GatewayDataPoints = getNewGatewayDataPointForTest()
var SDCMappingResourceID2 = os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_ID2")
var SDCMappingResourceName2 = os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME2")
var SDCVolName = os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME")
var SdsID = os.Getenv("POWERFLEX_DEVICE_SDS_ID")

func init() {
	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	username := os.Getenv("POWERFLEX_USERNAME")
	password := os.Getenv("POWERFLEX_PASSWORD")
	endpoint := os.Getenv("POWERFLEX_ENDPOINT")

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
	if v := os.Getenv("POWERFLEX_USERNAME"); v == "" {
		t.Fatal("POWERFLEX_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("POWERFLEX_PASSWORD"); v == "" {
		t.Fatal("POWERFLEX_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("POWERFLEX_ENDPOINT"); v == "" {
		t.Fatal("POWERFLEX_ENDPOINT must be set for acceptance tests")
	}
}
