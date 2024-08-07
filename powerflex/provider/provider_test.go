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
	"log"
	"os"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

var ProviderConfigForTesting = ``
var FunctionMocker *Mocker

type sdsDataPoints struct {
	SdsIP1             string
	SdsIP2             string
	SdsIP3             string
	SdsIP4             string
	SdsIP5             string
	SdsIP6             string
	SdsIP7             string
	SdsIP8             string
	SdsIP9             string
	SdsIP10            string
	SdsIP11            string
	SdcIP              string
	SdcIP1             string
	volName            string
	volName2           string
	volName3           string
	sdcName            string
	sdcName2           string
	sdcName3           string
	sdcHostExistingIP  string
	sdcHostWinIP       string
	sdcHostWinUserName string
	sdcHostWinPassword string
	sdcHostRPMIP       string
	sdcHostRPMUserName string
	sdcHostRPMPassword string
}

type sdcHostDataPoints struct {
	UbuntuIP       string
	UbuntuUser     string
	UbuntuPassword string
	UbuntuPort     string
	UbuntuPkgPath  string
	EsxiIP         string
	EsxiUser       string
	EsxiPassword   string
	EsxiPort       string
	EsxiPkgPath    string
	MdmIPs         []string
	CLS1           string
	CLS2           string
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
	virtualIP          string
	virtualInterface   string
	dataNetworkIP1     string
	dataNetworkIP2     string
	dataNetworkIP3     string
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

type nodeDataPoints struct {
	NodeIP       string
	ServiceTag   string
	NodeID       string
	NodePoolID   string
	NodePoolName string
}

type templateDataPoints struct {
	TemplateID   string
	TemplateName string
}

type serviceDataPoints struct {
	ServiceID   string
	ServiceName string
}

type complianceReportDataPoints struct {
	ResourceGroupID string
	IPAddress       string
	Compliant       string
	HostName        string
	ResourceID      string
	ServiceTag      string
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
	SdsResourceTestData.sdcHostExistingIP = setDefault(os.Getenv("POWERFLEX_SDC_HOST_EXISTING_IP"), "tfacc_sdc_host_existing_host_ip")
	SdsResourceTestData.sdcHostWinIP = setDefault(os.Getenv("POWERFLEX_SDC_HOST_WINDOWS_IP"), "127.0.0.1")
	SdsResourceTestData.sdcHostWinUserName = setDefault(os.Getenv("POWERFLEX_SDC_HOST_WINDOWS_USERNAME"), "tfacc_sdc_host_name")
	SdsResourceTestData.sdcHostWinPassword = setDefault(os.Getenv("POWERFLEX_SDC_HOST_WINDOWS_PASSWORD"), "tfacc_sdc_host_password")
	SdsResourceTestData.sdcHostRPMIP = setDefault(os.Getenv("POWERFLEX_SDC_HOST_LINUX_RPM_IP"), "127.0.0.1")
	SdsResourceTestData.sdcHostRPMUserName = setDefault(os.Getenv("POWERFLEX_SDC_HOST_LINUX_RPM_USERNAME"), "tfacc_sdc_host_name")
	SdsResourceTestData.sdcHostRPMPassword = setDefault(os.Getenv("POWERFLEX_SDC_HOST_LINUX_RPM_PASSWORD"), "tfacc_sdc_host_password")

	return SdsResourceTestData
}

func getNewSdcHostDataPointForTest() sdcHostDataPoints {
	var SdcHostDataPoints sdcHostDataPoints
	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return SdcHostDataPoints
	}

	SdcHostDataPoints.UbuntuIP = setDefault(os.Getenv("POWERFLEX_SDC_IP_Ubuntu"), "127.0.0.1")
	SdcHostDataPoints.UbuntuUser = setDefault(os.Getenv("POWERFLEX_SDC_USER_Ubuntu"), "ubuntuRoot")
	SdcHostDataPoints.UbuntuPassword = setDefault(os.Getenv("POWERFLEX_SDC_PASSWORD_Ubuntu"), "secret")
	SdcHostDataPoints.UbuntuPort = setDefault(os.Getenv("POWERFLEX_SDC_PORT_Ubuntu"), "2222")
	SdcHostDataPoints.UbuntuPkgPath = setDefault(os.Getenv("POWERFLEX_SDC_PKG_PATH_Ubuntu"), "/tmp/tfaccsdc.tar")

	SdcHostDataPoints.EsxiIP = setDefault(os.Getenv("POWERFLEX_SDC_IP_Esxi"), "127.0.0.1")
	SdcHostDataPoints.EsxiUser = setDefault(os.Getenv("POWERFLEX_SDC_USER_Esxi"), "esxiRoot")
	SdcHostDataPoints.EsxiPassword = setDefault(os.Getenv("POWERFLEX_SDC_PASSWORD_Esxi"), "secret")
	SdcHostDataPoints.EsxiPort = setDefault(os.Getenv("POWERFLEX_SDC_PORT_Esxi"), "2222")
	SdcHostDataPoints.EsxiPkgPath = setDefault(os.Getenv("POWERFLEX_SDC_PKG_PATH_Esxi"), "/tmp/tfaccsdc.tar")
	SdcHostDataPoints.CLS1 = setDefault(os.Getenv("POWERFLEX_SDC_CLS_1"), "1.1.1.1,2.2.2.2")
	SdcHostDataPoints.CLS2 = setDefault(os.Getenv("POWERFLEX_SDC_CLS_2"), "3.3.3.3,4.4.4.4")

	return SdcHostDataPoints
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
	GatewayDataPoints.virtualIP = setDefault(os.Getenv("POWERFLEX_VIRTUAL_IP"), "tfacc_cluster_virtual_ip")
	GatewayDataPoints.virtualInterface = setDefault(os.Getenv("POWERFLEX_VIRTUAL_INTERFACE"), "tfacc_cluster_virtual_interface")
	GatewayDataPoints.dataNetworkIP1 = setDefault(os.Getenv("POWERFLEX_DATA_NETWORK_IP1"), "tfacc_cluster_data_network_ip1")
	GatewayDataPoints.dataNetworkIP2 = setDefault(os.Getenv("POWERFLEX_DATA_NETWORK_IP2"), "tfacc_cluster_data_network_ip2")
	GatewayDataPoints.dataNetworkIP3 = setDefault(os.Getenv("POWERFLEX_DATA_NETWORK_IP3"), "tfacc_cluster_data_network_ip3")
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

func getTemplateDataForTest() templateDataPoints {

	var TemplateDataPoints templateDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return TemplateDataPoints
	}

	TemplateDataPoints.TemplateID = setDefault(os.Getenv("POWERFLEX_TEMPLATE_ID"), "tfacc_node_id")
	TemplateDataPoints.TemplateName = setDefault(os.Getenv("POWERFLEX_TEMPLATE_NAME"), "tfacc_node_pool_name")
	return TemplateDataPoints
}

func getServiceDataForTest() serviceDataPoints {

	var ServiceDataPoints serviceDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return ServiceDataPoints
	}

	ServiceDataPoints.ServiceID = setDefault(os.Getenv("POWERFLEX_SERVICE_ID"), "tfacc_service_id")
	ServiceDataPoints.ServiceName = setDefault(os.Getenv("POWERFLEX_SERVICE_NAME"), "tfacc_service_name")
	return ServiceDataPoints
}

func getComplianceReportDataForTest() complianceReportDataPoints {
	var ComplianceReportDataPoints complianceReportDataPoints
	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return ComplianceReportDataPoints
	}
	ComplianceReportDataPoints.Compliant = setDefault(os.Getenv("POWERFLEX_COMP_REP_COMPLIANT"), "true")
	ComplianceReportDataPoints.ResourceID = setDefault(os.Getenv("POWERFLEX_COMP_REP_ID"), "tfacc_compliance_report_id")
	ComplianceReportDataPoints.ResourceGroupID = setDefault(os.Getenv("POWERFLEX_COMP_REP_GROUP_ID"), "tfacc_compliance_report_group_id")
	ComplianceReportDataPoints.HostName = setDefault(os.Getenv("POWERFLEX_COMP_REP_HOST_NAME"), "tfacc_compliance_report_host_name")
	ComplianceReportDataPoints.ServiceTag = setDefault(os.Getenv("POWERFLEX_COMP_REP_SERVICE_TAG"), "tfacc_compliance_report_service_tag")
	ComplianceReportDataPoints.IPAddress = setDefault(os.Getenv("POWERFLEX_COMP_REP_IP_ADDRESS"), "tfacc_compliance_report_ip_address")
	return ComplianceReportDataPoints
}

func getNodeDataForTest() nodeDataPoints {

	var NodeDataPoints nodeDataPoints

	err := godotenv.Load("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return NodeDataPoints
	}

	NodeDataPoints.NodeIP = setDefault(os.Getenv("POWERFLEX_NODE_IP"), "tfacc_node_ip")
	NodeDataPoints.ServiceTag = setDefault(os.Getenv("POWERFLEX_SERVICE_TAG"), "tfacc_service_tag")
	NodeDataPoints.NodeID = setDefault(os.Getenv("POWERFLEX_NODE_ID"), "tfacc_node_id")
	NodeDataPoints.NodePoolID = setDefault(os.Getenv("POWERFLEX_NODE_POOL_ID"), "tfacc_node_pool_id")
	NodeDataPoints.NodePoolName = setDefault(os.Getenv("POWERFLEX_NODE_POOL_NAME"), "tfacc_node_pool_name")
	return NodeDataPoints
}

var SdsResourceTestData = getNewSdsDataPointForTest()
var SdcHostResourceTestData = getNewSdcHostDataPointForTest()
var GatewayDataPoints = getNewGatewayDataPointForTest()
var SDCMappingResourceID2 = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_ID2"), "tfacc_sdc_volumes_mapping_id2")
var SDCMappingResourceName2 = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME2"), "tfacc_sdc_volumes_mapping_name2")
var SDCVolName = setDefault(os.Getenv("POWERFLEX_SDC_VOLUMES_MAPPING_NAME"), "tfacc_sdc_volumes_mapping_name")
var SdsID = setDefault(os.Getenv("POWERFLEX_DEVICE_SDS_ID"), "tfacc_device_sds_id")
var MDMDataPoints = getMdmDataPointsForTest()
var StoragePoolName = setDefault(os.Getenv("POWERFLEX_STORAGE_POOL_NAME"), "tfacc_storage_pool_name")
var ProtectionDomainID = setDefault(os.Getenv("POWERFLEX_PROTECTION_DOMAIN_ID"), "tfacc_protection_domain_id")
var username = setDefault(os.Getenv("POWERFLEX_USERNAME"), "test")
var password = setDefault(os.Getenv("POWERFLEX_PASSWORD"), "test")
var endpoint = setDefault(os.Getenv("POWERFLEX_ENDPOINT"), "http://localhost:3002")
var insecure = setDefault(os.Getenv("POWERFLEX_INSECURE"), "false")
var NodeDataPoints = getNodeDataForTest()
var TemplateDataPoints = getTemplateDataForTest()
var ServiceDataPoints = getServiceDataForTest()
var SourceLocation = setDefault(os.Getenv("POWERFLEX_SOURCE_LOCATION"), "tfacc_source_location")
var FirmwareRepoID1 = setDefault(os.Getenv("POWERFLEX_FIRMWARE_REPO_ID1"), "tfacc_firmware_repo_id1")
var FirmwareRepoID2 = setDefault(os.Getenv("POWERFLEX_FIRMWARE_REPO_ID2"), "tfacc_firmware_repo_id2")
var FirmwareRepoName1 = setDefault(os.Getenv("POWERFLEX_FIRMWARE_REPO_NAME1"), "tfacc_firmware_repo_name1")
var FirmwareRepoName2 = setDefault(os.Getenv("POWERFLEX_FIRMWARE_REPO_NAME2"), "tfacc_firmware_repo_name2")
var ComplianceReportDataPoints = getComplianceReportDataForTest()
var OSRepoSourcePath = setDefault(os.Getenv("POWERFLEX_OS_REPO_SOURCE_PATH"), "tfacc_os_repo_source_path")

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
			insecure = %s
		}
	`, username, password, endpoint, insecure)
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
	// Make sure to unpatch before each new test is run
	if FunctionMocker != nil {
		FunctionMocker.UnPatch()
	}
}

// if there is no os setting set, then use the default value
func setDefault(osInput string, defaultStr string) string {
	if osInput == "" {
		return defaultStr
	}
	return osInput
}
