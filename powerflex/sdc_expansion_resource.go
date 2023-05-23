package powerflex

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewSDCExpansionResource  resource for the SDC Expansion
func NewSDCExpansionResource() resource.Resource {
	return &sdcExpansionResource{}
}

type sdcExpansionResource struct {
	gatewayClient *goscaleio.GatewayClient
}

func (r *sdcExpansionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_expansion"
}

func (r *sdcExpansionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SDCExpansionResourceSchema
}

func (r *sdcExpansionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.gatewayClient = req.ProviderData.(*goscaleio.GatewayClient)
}

func (r *sdcExpansionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan CsvAndMdmDataModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// to make gateway available for installation
	queueOperationError := ResetInstallerQueue(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	tflog.Info(ctx, "Gateway Installer changed to idle phase before initiating process")

	parsecsvRespose, parseCSVError := ParseCSVOperation(ctx, &plan, r.gatewayClient)

	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while Parsing CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	mdmIP, mdmIPError := GetMDMIP(ctx, plan)
	if mdmIPError != nil {
		resp.Diagnostics.AddError(
			"Error while Getting MDM IP",
			"unexpected error: "+mdmIPError.Error(),
		)
		return
	}

	tflog.Info(ctx, "CSV File parsed ssucessfully")

	// Vaidate the MDM credentials
	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Details",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {

		tflog.Info(ctx, "MDM Details validated successfully")

		if !CheckForNewSDCIPs(strings.Split(parsecsvRespose.Message, ","), strings.Split(validateMDMResponse.Data, ",")) {
			installationError := InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

			if installationError != nil {
				resp.Diagnostics.AddError(
					"Error in Installation Process",
					"unexpected error: "+installationError.Error(),
				)
				return
			}

			validateMDMResponse, validateMDMError = ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)

			if validateMDMError != nil {
				resp.Diagnostics.AddError(
					"Error While Validating MDM Details",
					"unexpected error: "+validateMDMResponse.Message,
				)
				return
			}
		}

		data, dgs := updateState(validateMDMResponse, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "MDM Details updated to state file successfully")

		return

	}

	resp.Diagnostics.AddError(
		"Error While Validating MDM Credentials",
		"unexpected error: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
	)
	return

}

func updateState(gatewayResponse *goscaleio_types.GatewayResponse, plan CsvAndMdmDataModel) (CsvAndMdmDataModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	state.InstalledSDCIps = types.StringValue(gatewayResponse.Data)

	state.ID = types.StringValue("placeholder")

	return state, diags
}

func (r *sdcExpansionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve values from state
	var state CsvAndMdmDataModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mdmIP, mdmIPError := GetMDMIP(ctx, state)
	if mdmIPError != nil {
		resp.Diagnostics.AddError(
			"Error while Getting MDM IP",
			"unexpected error: "+mdmIPError.Error(),
		)
		return
	}

	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, state, r.gatewayClient, mdmIP)

	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Details",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {
		data, dgs := updateState(validateMDMResponse, state)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.AddError(
		"Error While Validating MDM Credentials",
		"unexpected error: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
	)
	return
}

func (r *sdcExpansionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CsvAndMdmDataModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// to make gateway available for installation
	queueOperationError := ResetInstallerQueue(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	tflog.Info(ctx, "Gateway Installer changed to idle phase before initiating process")

	parsecsvRespose, parseCSVError := ParseCSVOperation(ctx, &plan, r.gatewayClient)

	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while Parsing CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	mdmIP, mdmIPError := GetMDMIP(ctx, plan)
	if mdmIPError != nil {
		resp.Diagnostics.AddError(
			"Error while Getting MDM IP",
			"unexpected error: "+mdmIPError.Error(),
		)
		return
	}

	tflog.Info(ctx, "CSV File parsed ssucessfully")

	// Vaidate the MDM credentials
	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Details",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {

		tflog.Info(ctx, "MDM Details validated successfully")

		if !CheckForNewSDCIPs(strings.Split(parsecsvRespose.Message, ","), strings.Split(validateMDMResponse.Data, ",")) {
			installationError := InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

			if installationError != nil {
				resp.Diagnostics.AddError(
					"Error in Installation Process",
					"unexpected error: "+installationError.Error(),
				)
				return
			}

			validateMDMResponse, validateMDMError = ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)

			if validateMDMError != nil {
				resp.Diagnostics.AddError(
					"Error While Validating MDM Details",
					"unexpected error: "+validateMDMResponse.Message,
				)
				return
			}
		}

		data, dgs := updateState(validateMDMResponse, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "MDM Details updated to state file successfully")

		return

	}

	resp.Diagnostics.AddError(
		"Error While Validating MDM Credentials",
		"unexpected error: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
	)
	return

}

// Working as PowerFlex Gateway Server Cleanup
func (r *sdcExpansionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CsvAndMdmDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queueOperationError := ResetInstallerQueue(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// GetMDMIP function is used for fetch MDM IP from cluster details
func GetMDMIP(ctx context.Context, model CsvAndMdmDataModel) (string, error) {
	var mdmIP string
	csvItems := []CSVDataModel{}
	diags := model.ClusterDetails.ElementsAs(ctx, &csvItems, true)

	if diags.HasError() {
		return mdmIP, fmt.Errorf("Error While Parse CSV Data  is %s", diags.Errors())
	}

	for _, item := range csvItems {
		if strings.EqualFold(item.IsMdmOrTb.ValueString(), "primary") {
			mdmIP = item.IP.ValueString()
			return mdmIP, nil
		}
	}
	return mdmIP, nil
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
func ParseCSVOperation(ctx context.Context, model *CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient) (*goscaleio_types.GatewayResponse, error) {

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
	header := []string{"IPs", "Username", "Password", "Operating System", "Is MDM/TB", "Is SDC", "perfProfileForSDC", "SDC Name"}
	err = writer.Write(header)
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
	}

	csvItems := []CSVDataModel{}
	diags := model.ClusterDetails.ElementsAs(ctx, &csvItems, true)
	if diags.HasError() {
		return &parseCSVResponse, fmt.Errorf("Error While Parse CSV Data  is %s", diags.Errors())
	}

	var newSDCIPs []string

	for _, item := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			IP:                 item.IP.ValueString(),
			UserName:           item.UserName.ValueString(),
			Password:           item.Password.ValueString(),
			IsMdmOrTb:          item.IsMdmOrTb.ValueString(),
			OperatingSystem:    item.OperatingSystem.ValueString(),
			IsSdc:              item.IsSdc.ValueString(),
			PerformanceProfile: item.PerformanceProfile.ValueString(),
			SDCName:            item.SDCName.ValueString(),
		}

		if !strings.EqualFold(csvStruct.IsMdmOrTb, "primary") && !strings.EqualFold(csvStruct.IsMdmOrTb, "secondary") && !strings.EqualFold(csvStruct.IsMdmOrTb, "tb") {
			newSDCIPs = append(newSDCIPs, csvStruct.IP)
		}

		//Write the data row
		data := []string{csvStruct.IP, csvStruct.UserName, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc, csvStruct.PerformanceProfile, csvStruct.SDCName}
		err = writer.Write(data)
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
		}
	}
	writer.Flush()

	parsecsvRespose, parseCSVError := gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if parseCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("%s", parseCSVError.Error())
	}

	if len(newSDCIPs) == 0 {
		return &parseCSVResponse, fmt.Errorf("No SDC Expansion Details are provided")
	}

	parsecsvRespose.Message = strings.Join(newSDCIPs, ",")

	deletCSVError := os.Remove(mydir + "/Minimal.csv")
	if deletCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Deleting Temp CSV File is %s", deletCSVError.Error())
	}

	if parsecsvRespose.StatusCode != 200 {
		return &parseCSVResponse, fmt.Errorf("Meesage : %s, Error Cosde : %s", parsecsvRespose.Message, strconv.Itoa(parsecsvRespose.StatusCode))
	}

	return parsecsvRespose, nil
}

// ValidateMDMOperation function for Vaidate the MDM credentials
func ValidateMDMOperation(ctx context.Context, model CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient, mdmIP string) (*goscaleio_types.GatewayResponse, error) {
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
func InstallationOperations(ctx context.Context, model CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient, parsecsvRespose *goscaleio_types.GatewayResponse) error {
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
				if couterForStopExecution < 5 {

					tflog.Info(ctx, "Some Gateway Installation operations are failed, Retrying...")

					_, err := gatewayClient.RetryPhase()

					if err != nil {
						return fmt.Errorf("Error while retrying failure is %s", err.Error())
					}

					couterForStopExecution = 5
				} else {

					return fmt.Errorf("Error During Installation is %s", checkForPhaseCompleted.Message)
				}
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
