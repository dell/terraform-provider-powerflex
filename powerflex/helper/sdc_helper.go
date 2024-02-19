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

package helper

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SdcFilterType - Enum structure for filter types.
var SdcFilterType = struct {
	All    string
	ByName string
	ByID   string
}{
	All:    "All",
	ByName: "ByName",
	ByID:   "ByID",
}

// GetFilteredSdcState - function to filter sdc result from goscaleio.
func GetFilteredSdcState(sdcs *[]models.SdcModel, method string, name string, id string) *[]models.SdcModel {
	response := []models.SdcModel{}
	for _, sdcValue := range *sdcs {
		if method == SdcFilterType.ByName && name == sdcValue.Name.ValueString() {
			response = append(response, sdcValue)
		}
		if method == SdcFilterType.ByID && id == sdcValue.ID.ValueString() {
			response = append(response, sdcValue)
		}
	}
	return &response
}

// GetAllSdcState - function to return all sdc result from goscaleio.
func GetAllSdcState(ctx context.Context, client goscaleio.Client, sdcs []goscaleio_types.Sdc) *[]models.SdcModel {
	response := []models.SdcModel{}
	for _, sdcValue := range sdcs {
		sdcState := models.SdcModel{
			ID:                 types.StringValue(sdcValue.ID),
			Name:               types.StringValue(sdcValue.Name),
			SdcGUID:            types.StringValue(sdcValue.SdcGUID),
			SdcApproved:        types.BoolValue(sdcValue.SdcApproved),
			OnVMWare:           types.BoolValue(sdcValue.OnVMWare),
			SystemID:           types.StringValue(sdcValue.SystemID),
			SdcIP:              types.StringValue(sdcValue.SdcIP),
			MdmConnectionState: types.StringValue(sdcValue.MdmConnectionState),
		}

		for _, link := range sdcValue.Links {
			sdcState.Links = append(sdcState.Links, models.SdcLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}

		response = append(response, sdcState)
	}

	return &response
}

// GetSDCDetailType returns the SDC Detail type
func GetSDCDetailType() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":                  types.StringType,
		"username":            types.StringType,
		"password":            types.StringType,
		"operating_system":    types.StringType,
		"is_mdm_or_tb":        types.StringType,
		"is_sdc":              types.StringType,
		"performance_profile": types.StringType,
		"sdc_id":              types.StringType,
		"name":                types.StringType,
		"virtual_ips":         types.StringType,
		"virtual_ip_nics":     types.StringType,
		"data_network_ip":     types.StringType,
	}
}

// GetSDCStateDetailType returns the SDC Detail type
func GetSDCStateDetailType() map[string]attr.Type {
	return map[string]attr.Type{
		"sdc_id":               types.StringType,
		"ip":                   types.StringType,
		"operating_system":     types.StringType,
		"performance_profile":  types.StringType,
		"name":                 types.StringType,
		"system_id":            types.StringType,
		"sdc_approved":         types.BoolType,
		"on_vmware":            types.BoolType,
		"sdc_guid":             types.StringType,
		"mdm_connection_state": types.StringType,
	}
}

// GetSDCStateDetailValue returns the SDC state detail object
func GetSDCStateDetailValue(sdc models.SDCStateDataModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDCStateDetailType(), map[string]attr.Value{
		"ip":                   types.StringValue(sdc.IP.ValueString()),
		"operating_system":     types.StringValue(sdc.OperatingSystem.ValueString()),
		"performance_profile":  types.StringValue(sdc.PerformanceProfile.ValueString()),
		"name":                 types.StringValue(sdc.SDCName.ValueString()),
		"sdc_id":               types.StringValue(sdc.SDCID.ValueString()),
		"system_id":            types.StringValue(sdc.SystemID.ValueString()),
		"sdc_approved":         types.BoolValue(sdc.SdcApproved.ValueBool()),
		"on_vmware":            types.BoolValue(sdc.OnVMWare.ValueBool()),
		"sdc_guid":             types.StringValue(sdc.SdcGUID.ValueString()),
		"mdm_connection_state": types.StringValue(sdc.MdmConnectionState.ValueString()),
	})
}

// GetSDCDetailValue returns the SDC detail object
func GetSDCDetailValue(sdc models.SDCDetailDataModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDCDetailType(), map[string]attr.Value{
		"ip":                  types.StringValue(sdc.IP.ValueString()),
		"username":            types.StringValue(sdc.UserName.ValueString()),
		"password":            types.StringValue(sdc.Password.ValueString()),
		"operating_system":    types.StringValue(sdc.OperatingSystem.ValueString()),
		"is_mdm_or_tb":        types.StringValue(sdc.IsMdmOrTb.ValueString()),
		"is_sdc":              types.StringValue(sdc.IsSdc.ValueString()),
		"performance_profile": types.StringValue(sdc.PerformanceProfile.ValueString()),
		"name":                types.StringValue(sdc.SDCName.ValueString()),
		"sdc_id":              types.StringValue(sdc.SDCID.ValueString()),
	})
}

// UpdateState - function to update state file for SDC resource.
func UpdateState(sdcs []models.SDCStateDataModel, plan models.SdcResourceModel) (models.SdcResourceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	SDCAttrTypes := GetSDCStateDetailType()

	SDCElemType := types.ObjectType{
		AttrTypes: SDCAttrTypes,
	}

	objectSDCs := []attr.Value{}
	for _, sdc := range sdcs {
		objVal, dgs := GetSDCStateDetailValue(sdc)
		diags = append(diags, dgs...)
		objectSDCs = append(objectSDCs, objVal)
	}
	setSdcs, dgs := types.ListValue(SDCElemType, objectSDCs)
	diags = append(diags, dgs...)

	state.SDCStateDetails = setSdcs
	state.ID = types.StringValue("placeholder")

	return state, diags
}

// GetMDMIP function is used for fetch MDM IP from cluster details
func GetMDMIP(ctx context.Context, sdcDetails []models.SDCDetailDataModel) (string, error) {
	var mdmIP string

	for _, item := range sdcDetails {
		if strings.EqualFold(item.IsMdmOrTb.ValueString(), "Primary") {
			mdmIP = item.IP.ValueString()
			return mdmIP, nil
		}
	}
	return mdmIP, nil
}

// CheckForExpansion function is used for check for expansion
func CheckForExpansion(model []models.SDCDetailDataModel, stateSDCDetails []models.SDCStateDataModel) bool {
	performaneChangeSdc := false
	checkIP := make(map[string]bool)

	for _, element := range stateSDCDetails {
		checkIP[element.IP.ValueString()] = true
	}

	for _, item := range model {
		if item.Password.ValueString() != "" && !checkIP[item.IP.ValueString()] && strings.EqualFold(item.IsSdc.ValueString(), "Yes") {
			performaneChangeSdc = true
			break
		}
	}
	return performaneChangeSdc
}

// ParseCSVOperation function for Handling Parsing CSV Operation
func ParseCSVOperation(ctx context.Context, sdcDetails []models.SDCDetailDataModel, gatewayClient *goscaleio.GatewayClient) (*goscaleio_types.GatewayResponse, error) {

	var parseCSVResponse goscaleio_types.GatewayResponse

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Reading Current Directory is %s", err.Error())
	}
	// Create a csv writer
	filePath := filepath.Join(mydir, filepath.Clean("Minimal.csv"))
	file, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV is %s", err.Error())
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := make([]string, 0)
	var virtualIpFlag bool
	for _, item := range sdcDetails {
		if item.IsMdmOrTb.ValueString() == "Primary" && !item.VirtualIps.IsNull() {
			virtualIpFlag = true
			break
		}
	}

	if virtualIpFlag {
		header = []string{"IPs", "Username", "Password", "Operating System", "Is MDM/TB", "Virtual IPs", "Virtual IP NICs", "Is SDC", "perfProfileForSDC"}
	} else {
		header = []string{"IPs", "Username", "Password", "Operating System", "Is MDM/TB", "Is SDC", "perfProfileForSDC"}
	}
	err = writer.Write(header)
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
	}

	var sdcIPs []string

	for _, item := range sdcDetails {

		if item.Password.ValueString() != "" {
			var csvStruct models.CsvRow
			// Add mapped SDC
			if virtualIpFlag {
				csvStruct = models.CsvRow{
					IP:              item.IP.ValueString(),
					UserName:        item.UserName.ValueString(),
					Password:        item.Password.ValueString(),
					IsMdmOrTb:       item.IsMdmOrTb.ValueString(),
					OperatingSystem: item.OperatingSystem.ValueString(),
					IsSdc:           item.IsSdc.ValueString(),
					VirtualIps:      item.VirtualIps.ValueString(),
					VirtualIPNICs:   item.VirtualIPNICs.ValueString(),
				}
			} else {
				csvStruct = models.CsvRow{
					IP:              item.IP.ValueString(),
					UserName:        item.UserName.ValueString(),
					Password:        item.Password.ValueString(),
					IsMdmOrTb:       item.IsMdmOrTb.ValueString(),
					OperatingSystem: item.OperatingSystem.ValueString(),
					IsSdc:           item.IsSdc.ValueString(),
				}
			}

			if strings.EqualFold(csvStruct.IsSdc, "Yes") {
				sdcIPs = append(sdcIPs, csvStruct.IP)
			}

			if strings.EqualFold(item.PerformanceProfile.ValueString(), "HighPerformance") {
				csvStruct.PerformanceProfile = "High"
			}

			//Write the data row
			data := make([]string, 0)
			if virtualIpFlag {
				data = []string{csvStruct.IP, csvStruct.UserName, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.VirtualIps, csvStruct.VirtualIPNICs, csvStruct.IsSdc, csvStruct.PerformanceProfile} //, csvStruct.SDCName
			} else {
				data = []string{csvStruct.IP, csvStruct.UserName, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc, csvStruct.PerformanceProfile}
			}
			err = writer.Write(data)
			if err != nil {
				return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
			}
		}

	}
	writer.Flush()

	parsecsvRespose, parseCSVError := gatewayClient.ParseCSV(mydir + "/Minimal.csv")

	if parseCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("%s", parseCSVError.Error())
	}

	parsecsvRespose.Message = strings.Join(sdcIPs, ",")

	if parsecsvRespose.StatusCode != 200 {
		return &parseCSVResponse, fmt.Errorf("Meesage : %s, Error Cosde : %s", parsecsvRespose.Message, strconv.Itoa(parsecsvRespose.StatusCode))
	}

	return parsecsvRespose, nil
}

// ValidateMDMOperation function for Validate the MDM credentials
func ValidateMDMOperation(ctx context.Context, model models.SdcResourceModel, gatewayClient *goscaleio.GatewayClient, mdmIP string) (*goscaleio_types.GatewayResponse, error) {
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": model.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{mdmIP}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData
	jsonres, _ := json.Marshal(mapData)

	validateMDMResponse, validateMDMError := gatewayClient.ValidateMDMDetails(jsonres)
	if validateMDMError != nil {
		return validateMDMResponse, fmt.Errorf("%s", validateMDMError.Error())
	}

	return validateMDMResponse, nil
}

// InstallationOperations function for begin instllation process
func InstallationOperations(ctx context.Context, model models.SdcResourceModel, gatewayClient *goscaleio.GatewayClient, parsecsvRespose *goscaleio_types.GatewayResponse) error {

	beginInstallationResponse, installationError := gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", model.MdmPassword.ValueString(), model.LiaPassword.ValueString(), true, true, false, true)

	if installationError != nil {
		return fmt.Errorf("Error while begin installation is %s", installationError.Error())
	}

	if beginInstallationResponse.StatusCode == 200 {
		currentPhase := "query"
		couterForStopExecution := 0

		tflog.Info(ctx, "Gateway Installation Begin, Current Phase - Query")

		for couterForStopExecution <= 5 {

			time.Sleep(1 * time.Minute)

			checkForPhaseCompleted, _ := gatewayClient.CheckForCompletionQueueCommands(currentPhase)

			if checkForPhaseCompleted.Data == "Completed" {
				couterForStopExecution = 0

				if currentPhase != "configure" {
					moveToNextPhaseResponse, err := gatewayClient.MoveToNextPhase()

					if err != nil {
						return fmt.Errorf("Error while moving to next phase is %s", err.Error())
					}

					if moveToNextPhaseResponse.StatusCode == 200 {
						if currentPhase == "query" {
							currentPhase = "upload"
							tflog.Info(ctx, "Gateway Installation phase changed to Upload")
						} else if currentPhase == "upload" {
							currentPhase = "install"
							tflog.Info(ctx, "Gateway Installation phase changed to Install")
						} else if currentPhase == "install" {
							currentPhase = "configure"
							tflog.Info(ctx, "Gateway Installation phase changed to Configure")
						}
					} else {
						return fmt.Errorf("Messsage: %s, Error Code: %s", moveToNextPhaseResponse.Message, strconv.Itoa(moveToNextPhaseResponse.StatusCode))
					}
				} else {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
					}

					return nil
				}

			} else if checkForPhaseCompleted.Data == "Running" {
				couterForStopExecution++

				tflog.Info(ctx, "Gateway Installation operations are still running")

				if couterForStopExecution == 5 {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
					}

					return fmt.Errorf("Time Out,Some Operations of Installer running from since long")
				}

			} else {
				return fmt.Errorf("Error During Installation is %s", checkForPhaseCompleted.Message)
			}
		}
	} else {
		return fmt.Errorf("Message: %s, Error Code: %s", beginInstallationResponse.Message, strconv.Itoa(beginInstallationResponse.StatusCode))
	}

	return nil
}

// CheckForNewSDCIPs function to check SDC Alredy Installed or not
func CheckForNewSDCIPs(newSDCIPS []string, installedSDCIPs []string) bool {
	checkset := make(map[string]bool)
	for _, element := range newSDCIPS {
		checkset[element] = true
	}
	for _, value := range installedSDCIPs {
		if checkset[value] {
			delete(checkset, value)
		}
	}
	return len(checkset) == 0 //this implies that set is subset of superset
}

// GetSDCState - function to return sdc result from goscaleio.
func GetSDCState(sdc goscaleio_types.Sdc) models.SDCStateDataModel {
	var model models.SDCStateDataModel

	model.OperatingSystem = types.StringValue(sdc.OSType)
	model.SDCID = types.StringValue(sdc.ID)
	model.SDCName = types.StringValue(sdc.Name)
	model.SdcGUID = types.StringValue(sdc.SdcGUID)
	model.SdcApproved = types.BoolValue(sdc.SdcApproved)
	model.OnVMWare = types.BoolValue(sdc.OnVMWare)
	model.SystemID = types.StringValue(sdc.SystemID)
	model.PerformanceProfile = types.StringValue(sdc.PerfProfile)
	model.IP = types.StringValue(sdc.SdcIP)
	model.MdmConnectionState = types.StringValue(sdc.MdmConnectionState)
	return model
}

// CheckForSDCName - check for the SDC Name already exist or not
func CheckForSDCName(system *goscaleio.System, sdcDetail models.SDCDetailDataModel) (bool, error) {

	sdcData, err := system.FindSdc("Name", sdcDetail.SDCName.ValueString())

	if err == nil && sdcData.Sdc.ID == sdcDetail.SDCID.ValueString() {
		return true, fmt.Errorf("SDC Name:%s already exist", sdcDetail.SDCName.ValueString())
	}

	return false, nil
}

// FindDeletedSDC function to find deleted SDC Details in Plan
func FindDeletedSDC(state []models.SDCStateDataModel, plan []models.SDCDetailDataModel) []models.SDCStateDataModel {
	difference := []models.SDCStateDataModel{}
	checkID := make(map[string]bool)
	checkIP := make(map[string]bool)

	for _, element := range state {
		checkID[element.SDCID.ValueString()] = true
		checkIP[element.IP.ValueString()] = true
	}

	for _, obj1 := range state {
		found := false
		for _, obj2 := range plan {
			if (obj2.IP.ValueString() != "" || obj2.DataNetworkIP.ValueString() != "") && (checkIP[obj2.IP.ValueString()] || checkIP[obj2.DataNetworkIP.ValueString()]) {
				found = true
				delete(checkIP, obj2.IP.ValueString())
				delete(checkID, obj2.SDCID.ValueString())
				break
			} else if obj2.SDCID.ValueString() != "" && checkID[obj2.SDCID.ValueString()] {
				found = true
				delete(checkIP, obj2.IP.ValueString())
				delete(checkID, obj2.SDCID.ValueString())
				break
			}
		}
		if !found {
			difference = append(difference, obj1)
		}
	}

	return difference
}

// GetSdcsValue returns the SDC list for the plan
func GetSdcsValue(planSdcs []models.SDCDetailDataModel) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	sdcInfoElemType := types.ObjectType{
		AttrTypes: GetSDCDetailType(),
	}

	objectSdcInfos := []attr.Value{}
	for _, sdc := range planSdcs {
		obj := map[string]attr.Value{
			"ip":                  sdc.IP,
			"username":            sdc.UserName,
			"password":            sdc.Password,
			"operating_system":    sdc.OperatingSystem,
			"is_mdm_or_tb":        sdc.IsMdmOrTb,
			"is_sdc":              sdc.IsSdc,
			"performance_profile": sdc.PerformanceProfile,
			"sdc_id":              sdc.SDCID,
			"name":                sdc.SDCName,
			"virtual_ips":         sdc.VirtualIps,
			"virtual_ip_nics":     sdc.VirtualIPNICs,
			"data_network_ip":     sdc.DataNetworkIP,
		}
		objVal, dgs := types.ObjectValue(GetSDCDetailType(), obj)
		diags = append(diags, dgs...)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	sdcInfoVal, dgs := types.ListValue(sdcInfoElemType, objectSdcInfos)
	diags = append(diags, dgs...)
	return sdcInfoVal, diags
}
