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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return SdsResourceTestData
	}

	SdsResourceTestData.SdsIP1 = setDefault(envMap["POWERFLEX_SDS_IP_1"], "tfacc_sds_1")
	SdsResourceTestData.SdsIP2 = setDefault(envMap["POWERFLEX_SDS_IP_2"], "tfacc_sds_2")
	SdsResourceTestData.SdsIP3 = setDefault(envMap["POWERFLEX_SDS_IP_3"], "tfacc_sds_3")
	SdsResourceTestData.SdsIP4 = setDefault(envMap["POWERFLEX_SDS_IP_4"], "tfacc_sds_4")
	SdsResourceTestData.SdsIP5 = setDefault(envMap["POWERFLEX_SDS_IP_5"], "tfacc_sds_5")
	SdsResourceTestData.SdsIP6 = setDefault(envMap["POWERFLEX_SDS_IP_6"], "tfacc_sds_6")
	SdsResourceTestData.SdsIP7 = setDefault(envMap["POWERFLEX_SDS_IP_7"], "tfacc_sds_7")
	SdsResourceTestData.SdsIP8 = setDefault(envMap["POWERFLEX_SDS_IP_8"], "tfacc_sds_8")
	SdsResourceTestData.SdsIP9 = setDefault(envMap["POWERFLEX_SDS_IP_9"], "tfacc_sds_9")
	SdsResourceTestData.SdsIP10 = setDefault(envMap["POWERFLEX_SDS_IP_10"], "tfacc_sds_10")
	SdsResourceTestData.SdsIP11 = setDefault(envMap["POWERFLEX_SDS_IP_11"], "tfacc_sds_11")
	SdsResourceTestData.SdcIP = setDefault(envMap["POWERFLEX_SDC_IP"], "tfacc_sdc_ip_1")
	SdsResourceTestData.SdcIP1 = setDefault(envMap["POWERFLEX_SDC_IP1"], "tfacc_sdc_ip_2")
	SdsResourceTestData.volName = setDefault(envMap["POWERFLEX_VOLUME_NAME"], "tfacc_volume_1")
	SdsResourceTestData.volName2 = setDefault(envMap["POWERFLEX_VOLUME_NAME_2"], "tfacc_volume_2")
	SdsResourceTestData.volName3 = setDefault(envMap["POWERFLEX_VOLUME_NAME_3"], "tfacc_volume_3")
	SdsResourceTestData.sdcName = setDefault(envMap["POWERFLEX_SDC_NAME"], "tfacc_sdc_name_1")
	SdsResourceTestData.sdcName2 = setDefault(envMap["POWERFLEX_SDC_NAME_2"], "tfacc_sdc_name_2")
	SdsResourceTestData.sdcName3 = setDefault(envMap["POWERFLEX_SDC_NAME_3"], "tfacc_sdc_name_3")
	SdsResourceTestData.sdcHostExistingIP = setDefault(envMap["POWERFLEX_SDC_HOST_EXISTING_IP"], "tfacc_sdc_host_existing_host_ip")
	SdsResourceTestData.sdcHostWinIP = setDefault(envMap["POWERFLEX_SDC_HOST_WINDOWS_IP"], "127.0.0.1")
	SdsResourceTestData.sdcHostWinUserName = setDefault(envMap["POWERFLEX_SDC_HOST_WINDOWS_USERNAME"], "tfacc_sdc_host_name")
	SdsResourceTestData.sdcHostWinPassword = setDefault(envMap["POWERFLEX_SDC_HOST_WINDOWS_PASSWORD"], "tfacc_sdc_host_password")
	SdsResourceTestData.sdcHostRPMIP = setDefault(envMap["POWERFLEX_SDC_HOST_LINUX_RPM_IP"], "127.0.0.1")
	SdsResourceTestData.sdcHostRPMUserName = setDefault(envMap["POWERFLEX_SDC_HOST_LINUX_RPM_USERNAME"], "tfacc_sdc_host_name")
	SdsResourceTestData.sdcHostRPMPassword = setDefault(envMap["POWERFLEX_SDC_HOST_LINUX_RPM_PASSWORD"], "tfacc_sdc_host_password")

	return SdsResourceTestData
}

func getNewSdcHostDataPointForTest() sdcHostDataPoints {
	var SdcHostDataPoints sdcHostDataPoints
	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return SdcHostDataPoints
	}

	SdcHostDataPoints.UbuntuIP = setDefault(envMap["POWERFLEX_SDC_IP_Ubuntu"], "127.0.0.1")
	SdcHostDataPoints.UbuntuUser = setDefault(envMap["POWERFLEX_SDC_USER_Ubuntu"], "ubuntuRoot")
	SdcHostDataPoints.UbuntuPassword = setDefault(envMap["POWERFLEX_SDC_PASSWORD_Ubuntu"], "secret")
	SdcHostDataPoints.UbuntuPort = setDefault(envMap["POWERFLEX_SDC_PORT_Ubuntu"], "2222")
	SdcHostDataPoints.UbuntuPkgPath = setDefault(envMap["POWERFLEX_SDC_PKG_PATH_Ubuntu"], "/tmp/tfaccsdc.tar")

	SdcHostDataPoints.EsxiIP = setDefault(envMap["POWERFLEX_SDC_IP_Esxi"], "127.0.0.1")
	SdcHostDataPoints.EsxiUser = setDefault(envMap["POWERFLEX_SDC_USER_Esxi"], "esxiRoot")
	SdcHostDataPoints.EsxiPassword = setDefault(envMap["POWERFLEX_SDC_PASSWORD_Esxi"], "secret")
	SdcHostDataPoints.EsxiPort = setDefault(envMap["POWERFLEX_SDC_PORT_Esxi"], "2222")
	SdcHostDataPoints.EsxiPkgPath = setDefault(envMap["POWERFLEX_SDC_PKG_PATH_Esxi"], "/tmp/tfaccsdc.tar")
	SdcHostDataPoints.CLS1 = setDefault(envMap["POWERFLEX_SDC_CLS_1"], "1.1.1.1,2.2.2.2")
	SdcHostDataPoints.CLS2 = setDefault(envMap["POWERFLEX_SDC_CLS_2"], "3.3.3.3,4.4.4.4")

	return SdcHostDataPoints
}

func getNewGatewayDataPointForTest() gatewayDataPoints {

	var GatewayDataPoints gatewayDataPoints

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return GatewayDataPoints
	}

	GatewayDataPoints.primaryMDMIP = setDefault(envMap["POWERFLEX_PRIMARY_MDM_IP"], "tfacc_primary_mdm_ip")
	GatewayDataPoints.secondaryMDMIP = setDefault(envMap["POWERFLEX_SECONDARY_MDM_IP"], "tfacc_secondary_mdm_ip")
	GatewayDataPoints.tbIP = setDefault(envMap["POWERFLEX_TB_IP"], "tfacc_tb_ip")
	GatewayDataPoints.sdcServerIP = setDefault(envMap["POWERFLEX_SDC_SERVER_IP"], "tfacc_sdc_server_ip")
	GatewayDataPoints.serverPassword = setDefault(envMap["POWERFLEX_SERVER_PASSWORD"], "tfacc_server_password")
	GatewayDataPoints.mdmPassword = setDefault(envMap["POWERFLEX_MDM_PASSWORD"], "tfacc_mdm_password")
	GatewayDataPoints.liaPassword = setDefault(envMap["POWERFLEX_LIA_PASSWORD"], "tfacc_lia_password")
	GatewayDataPoints.clusterPrimaryIP = setDefault(envMap["POWERFLEX_CLUSTER_IP_1"], "tfacc_cluster_ip_1")
	GatewayDataPoints.clusterSecondaryIP = setDefault(envMap["POWERFLEX_CLUSTER_IP_2"], "tfacc_cluster_ip_2")
	GatewayDataPoints.clusterTBIP = setDefault(envMap["POWERFLEX_CLUSTER_IP_3"], "tfacc_cluster_ip_3")
	GatewayDataPoints.virtualIP = setDefault(envMap["POWERFLEX_VIRTUAL_IP"], "tfacc_cluster_virtual_ip")
	GatewayDataPoints.virtualInterface = setDefault(envMap["POWERFLEX_VIRTUAL_INTERFACE"], "tfacc_cluster_virtual_interface")
	GatewayDataPoints.dataNetworkIP1 = setDefault(envMap["POWERFLEX_DATA_NETWORK_IP1"], "tfacc_cluster_data_network_ip1")
	GatewayDataPoints.dataNetworkIP2 = setDefault(envMap["POWERFLEX_DATA_NETWORK_IP2"], "tfacc_cluster_data_network_ip2")
	GatewayDataPoints.dataNetworkIP3 = setDefault(envMap["POWERFLEX_DATA_NETWORK_IP3"], "tfacc_cluster_data_network_ip3")
	return GatewayDataPoints
}

func getMdmDataPointsForTest() mdmDataPoints {

	var MDMDataPoints mdmDataPoints

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return MDMDataPoints
	}

	MDMDataPoints.primaryMDMIP = setDefault(envMap["POWERFLEX_PRIMARY_MDM_IP"], "tfacc_primary_mdm_ip")
	MDMDataPoints.secondaryMDMIP = setDefault(envMap["POWERFLEX_SECONDARY_MDM_IP"], "tfacc_secondary_mdm_ip")
	MDMDataPoints.tbIP = setDefault(envMap["POWERFLEX_TB_IP"], "tfacc_tb_ip")
	MDMDataPoints.primaryMDMID = setDefault(envMap["POWERFLEX_PRIMARY_MDM_ID"], "tfacc_primary_mdm_id")
	MDMDataPoints.secondaryMDMID = setDefault(envMap["POWERFLEX_SECONDARY_MDM_ID"], "tfacc_secondary_mdm_id")
	MDMDataPoints.tbID = setDefault(envMap["POWERFLEX_TB_ID"], "tfacc_tb_id")
	MDMDataPoints.standByIP1 = setDefault(envMap["POWERFLEX_STANDBY_MDM_IP1"], "1.1.1.1")
	MDMDataPoints.standByIP2 = setDefault(envMap["POWERFLEX_STANDBY_MDM_IP2"], "1.1.1.2")

	return MDMDataPoints
}

func getTemplateDataForTest() templateDataPoints {

	var TemplateDataPoints templateDataPoints

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return TemplateDataPoints
	}

	TemplateDataPoints.TemplateID = setDefault(envMap["POWERFLEX_TEMPLATE_ID"], "tfacc_node_id")
	TemplateDataPoints.TemplateName = setDefault(envMap["POWERFLEX_TEMPLATE_NAME"], "tfacc_node_pool_name")
	return TemplateDataPoints
}

func getServiceDataForTest() serviceDataPoints {

	var ServiceDataPoints serviceDataPoints

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return ServiceDataPoints
	}

	ServiceDataPoints.ServiceID = setDefault(envMap["POWERFLEX_SERVICE_ID"], "tfacc_service_id")
	ServiceDataPoints.ServiceName = setDefault(envMap["POWERFLEX_SERVICE_NAME"], "tfacc_service_name")
	return ServiceDataPoints
}

func getComplianceReportDataForTest() complianceReportDataPoints {
	var ComplianceReportDataPoints complianceReportDataPoints
	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return ComplianceReportDataPoints
	}

	ComplianceReportDataPoints.Compliant = setDefault(envMap["POWERFLEX_COMP_REP_COMPLIANT"], "true")
	ComplianceReportDataPoints.ResourceID = setDefault(envMap["POWERFLEX_COMP_REP_ID"], "tfacc_compliance_report_id")
	ComplianceReportDataPoints.ResourceGroupID = setDefault(envMap["POWERFLEX_COMP_REP_GROUP_ID"], "tfacc_compliance_report_group_id")
	ComplianceReportDataPoints.HostName = setDefault(envMap["POWERFLEX_COMP_REP_HOST_NAME"], "tfacc_compliance_report_host_name")
	ComplianceReportDataPoints.ServiceTag = setDefault(envMap["POWERFLEX_COMP_REP_SERVICE_TAG"], "tfacc_compliance_report_service_tag")
	ComplianceReportDataPoints.IPAddress = setDefault(envMap["POWERFLEX_COMP_REP_IP_ADDRESS"], "tfacc_compliance_report_ip_address")
	return ComplianceReportDataPoints
}

func getNodeDataForTest() nodeDataPoints {

	var NodeDataPoints nodeDataPoints

	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return NodeDataPoints
	}

	NodeDataPoints.NodeIP = setDefault(envMap["POWERFLEX_NODE_IP"], "tfacc_node_ip")
	NodeDataPoints.ServiceTag = setDefault(envMap["POWERFLEX_SERVICE_TAG"], "tfacc_service_tag")
	NodeDataPoints.NodeID = setDefault(envMap["POWERFLEX_NODE_ID"], "tfacc_node_id")
	NodeDataPoints.NodePoolID = setDefault(envMap["POWERFLEX_NODE_POOL_ID"], "tfacc_node_pool_id")
	NodeDataPoints.NodePoolName = setDefault(envMap["POWERFLEX_NODE_POOL_NAME"], "tfacc_node_pool_name")
	return NodeDataPoints
}

var globalEnvMap = getEnvMap()
var SdsResourceTestData = getNewSdsDataPointForTest()
var SdcHostResourceTestData = getNewSdcHostDataPointForTest()
var GatewayDataPoints = getNewGatewayDataPointForTest()
var SDCMappingResourceID2 = setDefault(globalEnvMap["POWERFLEX_SDC_VOLUMES_MAPPING_ID2"], "tfacc_sdc_volumes_mapping_id2")
var SDCMappingResourceName2 = setDefault(globalEnvMap["POWERFLEX_SDC_VOLUMES_MAPPING_NAME2"], "tfacc_sdc_volumes_mapping_name2")
var SDCVolName = setDefault(globalEnvMap["POWERFLEX_SDC_VOLUMES_MAPPING_NAME"], "tfacc_sdc_volumes_mapping_name")
var SdsID = setDefault(globalEnvMap["POWERFLEX_DEVICE_SDS_ID"], "tfacc_device_sds_id")
var MDMDataPoints = getMdmDataPointsForTest()
var StoragePoolName = setDefault(globalEnvMap["POWERFLEX_STORAGE_POOL_NAME"], "tfacc_storage_pool_name")
var ProtectionDomainID = setDefault(globalEnvMap["POWERFLEX_PROTECTION_DOMAIN_ID"], "tfacc_protection_domain_id")
var ProtectionDomainIDSds = setDefault(globalEnvMap["POWERFLEX_PROTECTION_DOMAIN_ID_SDS"], "tfacc_protection_domain_id_sds")
var username = setDefault(globalEnvMap["POWERFLEX_USERNAME"], "test")
var password = setDefault(globalEnvMap["POWERFLEX_PASSWORD"], "test")
var endpoint = setDefault(globalEnvMap["POWERFLEX_ENDPOINT"], "http://localhost:3002")
var insecure = setDefault(globalEnvMap["POWERFLEX_INSECURE"], "false")
var NodeDataPoints = getNodeDataForTest()
var TemplateDataPoints = getTemplateDataForTest()
var ServiceDataPoints = getServiceDataForTest()
var SourceLocation = setDefault(globalEnvMap["POWERFLEX_SOURCE_LOCATION"], "tfacc_source_location")
var FirmwareRepoID1 = setDefault(globalEnvMap["POWERFLEX_FIRMWARE_REPO_ID1"], "tfacc_firmware_repo_id1")
var FirmwareRepoID2 = setDefault(globalEnvMap["POWERFLEX_FIRMWARE_REPO_ID2"], "tfacc_firmware_repo_id2")
var FirmwareRepoName1 = setDefault(globalEnvMap["POWERFLEX_FIRMWARE_REPO_NAME1"], "tfacc_firmware_repo_name1")
var FirmwareRepoName2 = setDefault(globalEnvMap["POWERFLEX_FIRMWARE_REPO_NAME2"], "tfacc_firmware_repo_name2")
var ComplianceReportDataPoints = getComplianceReportDataForTest()
var OSRepoSourcePath = setDefault(globalEnvMap["POWERFLEX_OS_REPO_SOURCE_PATH"], "tfacc_os_repo_source_path")
var OSRepoID1 = setDefault(globalEnvMap["POWERFLEX_OS_REPO_ID1"], "tfacc_os_repo_id1")
var OSRepoName1 = setDefault(globalEnvMap["POWERFLEX_OS_REPO_NAME1"], "tfacc_os_repo_name1")
var OSRepoState = setDefault(globalEnvMap["POWERFLEX_OS_REPO_STATE"], "available")
var OSRepoImageType = setDefault(globalEnvMap["POWERFLEX_OS_REPO_IMAGE_TYPE"], "vmware_esxi")
var OSRepoType = setDefault(globalEnvMap["POWERFLEX_OS_REPO_TYPE"], "ISO")
var OSRepoCreatedBy = setDefault(globalEnvMap["POWERFLEX_OS_REPO_CREATED_BY"], "system")
var FaultSetID = setDefault(globalEnvMap["POWERFLEX_FAULT_SET_ID"], "1EE7752911111112")

func getEnvMap() map[string]string {
	envMap, err := loadEnvFile("powerflex.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return envMap
	}
	return envMap
}

func init() {
	ProviderConfigForTesting = fmt.Sprintf(`
		provider "powerflex" {
			username = "%s"
			password = "%s"
			endpoint = "%s"
			insecure = %s
		}
	`, username, password, endpoint, insecure)
	// Set the specific TF_ACC test environment
	os.Setenv("TF_ACC", globalEnvMap["TF_ACC"])
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

func loadEnvFile(path string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
