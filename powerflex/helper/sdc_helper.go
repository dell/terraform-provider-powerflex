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
		"sdc_id":               types.StringType,
		"ip":                   types.StringType,
		"username":             types.StringType,
		"password":             types.StringType,
		"operating_system":     types.StringType,
		"is_mdm_or_tb":         types.StringType,
		"is_sdc":               types.StringType,
		"performance_profile":  types.StringType,
		"name":                 types.StringType,
		"system_id":            types.StringType,
		"sdc_approved":         types.BoolType,
		"on_vmware":            types.BoolType,
		"sdc_guid":             types.StringType,
		"mdm_connection_state": types.StringType,
	}
}

// GetSDCDetailValue returns the SDC Detail model object value
func GetSDCDetailValue(sdc models.SDCDetailDataModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDCDetailType(), map[string]attr.Value{
		"sdc_id":               types.StringValue(sdc.SDCID.ValueString()),
		"ip":                   types.StringValue(sdc.IP.ValueString()),
		"username":             types.StringValue(sdc.UserName.ValueString()),
		"password":             types.StringValue(sdc.Password.ValueString()),
		"operating_system":     types.StringValue(sdc.OperatingSystem.ValueString()),
		"is_mdm_or_tb":         types.StringValue(sdc.IsMdmOrTb.ValueString()),
		"is_sdc":               types.StringValue(sdc.IsSdc.ValueString()),
		"performance_profile":  types.StringValue(sdc.PerformanceProfile.ValueString()),
		"name":                 types.StringValue(sdc.SDCName.ValueString()),
		"system_id":            types.StringValue(sdc.SystemID.ValueString()),
		"sdc_approved":         types.BoolValue(sdc.SdcApproved.ValueBool()),
		"on_vmware":            types.BoolValue(sdc.OnVMWare.ValueBool()),
		"sdc_guid":             types.StringValue(sdc.SdcGUID.ValueString()),
		"mdm_connection_state": types.StringValue(sdc.MdmConnectionState.ValueString()),
	})
}

// UpdateState - function to update state file for SDC resource.
func UpdateState(sdcs []models.SDCDetailDataModel, plan models.SdcResourceModel) (models.SdcResourceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	SDCAttrTypes := GetSDCDetailType()

	SDCElemType := types.ObjectType{
		AttrTypes: SDCAttrTypes,
	}

	objectSDCs := []attr.Value{}
	for _, sdc := range sdcs {
		objVal, dgs := GetSDCDetailValue(sdc)
		diags = append(diags, dgs...)
		objectSDCs = append(objectSDCs, objVal)
	}
	setSdcs, dgs := types.ListValue(SDCElemType, objectSDCs)
	diags = append(diags, dgs...)

	state.SDCDetails = setSdcs

	if plan.ID.ValueString() != "" && len(strings.Split(plan.ID.ValueString(), ",")) == 1 {
		state.ID = plan.ID
	} else {
		state.ID = types.StringValue("placeholder")
	}

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
func CheckForExpansion(model []models.SDCDetailDataModel) bool {
	performaneChangeSdc := false

	for _, item := range model {
		if item.Password.ValueString() != "" && strings.EqualFold(item.IsSdc.ValueString(), "Yes") {
			performaneChangeSdc = true
			break
		}
	}
	return performaneChangeSdc
}

// ResetInstallerQueue function for the Abort, Clear and Move To Idle Execution
func ResetInstallerQueue(gatewayClient *goscaleio.GatewayClient) error {

	_, err := gatewayClient.AbortOperation()

	if err != nil {
		return fmt.Errorf("Error while Aborting Operation is %s", err.Error())
	}
	_, err = gatewayClient.ClearQueueCommand()

	if err != nil {
		return fmt.Errorf("Error while Clearing Queue is %s", err.Error())
	}

	_, err = gatewayClient.MoveToIdlePhase()

	if err != nil {
		return fmt.Errorf("Error while Move to Ideal Phase is %s", err.Error())
	}

	return nil
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
	file, err := os.Create(mydir + "/Minimal.csv")
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV is %s", err.Error())
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"IPs", "Username", "Password", "Operating System", "Is MDM/TB", "Is SDC", "perfProfileForSDC"}
	err = writer.Write(header)
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
	}

	var sdcIPs []string

	for _, item := range sdcDetails {

		if item.Password.ValueString() != "" {
			// Add mapped SDC
			csvStruct := models.CsvRow{
				IP:              item.IP.ValueString(),
				UserName:        item.UserName.ValueString(),
				Password:        item.Password.ValueString(),
				IsMdmOrTb:       item.IsMdmOrTb.ValueString(),
				OperatingSystem: item.OperatingSystem.ValueString(),
				IsSdc:           item.IsSdc.ValueString(),
			}

			if strings.EqualFold(csvStruct.IsSdc, "Yes") {
				sdcIPs = append(sdcIPs, csvStruct.IP)
			}

			if strings.EqualFold(item.PerformanceProfile.ValueString(), "HighPerformance") {
				csvStruct.PerformanceProfile = "High"
			}

			//Write the data row
			data := []string{csvStruct.IP, csvStruct.UserName, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc, csvStruct.PerformanceProfile} //, csvStruct.SDCName
			err = writer.Write(data)
			if err != nil {
				return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
			}
		}

	}
	writer.Flush()

	parsecsvRespose, parseCSVError := gatewayClient.ParseCSV(mydir + "/Minimal.csv")

	deletCSVError := os.Remove(mydir + "/Minimal.csv")
	if deletCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Deleting Temp CSV File is %s", deletCSVError.Error())
	}

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

	beginInstallationResponse, installationError := gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", model.MdmPassword.ValueString(), model.LiaPassword.ValueString(), true)

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

					couterForStopExecution = 10

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
func GetSDCState(sdc goscaleio_types.Sdc, model models.SDCDetailDataModel) (response models.SDCDetailDataModel) {

	if sdc.ID != "" {
		model.SDCID = types.StringValue(sdc.ID)

		model.SDCName = types.StringValue(sdc.Name)

		if sdc.SdcGUID != "" {
			model.SdcGUID = types.StringValue(sdc.SdcGUID)
		}

		model.SdcApproved = types.BoolValue(sdc.SdcApproved)

		model.OnVMWare = types.BoolValue(sdc.OnVMWare)

		if sdc.SystemID != "" {
			model.SystemID = types.StringValue(sdc.SystemID)
		}

		model.PerformanceProfile = types.StringValue(sdc.PerfProfile)

		if sdc.SdcIP != "" {
			model.IP = types.StringValue(sdc.SdcIP)
		}

		if sdc.MdmConnectionState != "" {
			model.MdmConnectionState = types.StringValue(sdc.MdmConnectionState)
		}
	}

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
func FindDeletedSDC(state, plan []models.SDCDetailDataModel) []models.SDCDetailDataModel {
	difference := []models.SDCDetailDataModel{}

	for _, obj1 := range state {
		found := false
		for _, obj2 := range plan {
			if obj2.IP.ValueString() != "" && obj1.IP == obj2.IP {
				found = true
				break
			} else if obj2.SDCID.ValueString() != "" && obj1.SDCID == obj2.SDCID {
				found = true
				break
			} else if obj2.SDCName.ValueString() != "" && obj1.SDCName == obj2.SDCName {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, obj1)
		}
	}

	return difference
}
