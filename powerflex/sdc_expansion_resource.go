package powerflex

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	queueOperationError := QueueOpeartions(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	parsecsvRespose, parseCSVError := ParseCSVOperation(ctx, plan, r.gatewayClient)

	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while Parsing CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	// Vaidate the MDM credentials
	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, plan, r.gatewayClient)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Details",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {

		installationError := InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

		if installationError != nil {
			resp.Diagnostics.AddError(
				"Error in Installation Process",
				"unexpected error: "+installationError.Error(),
			)
			return
		}

		validateMDMResponse, validateMDMError = ValidateMDMOperation(ctx, plan, r.gatewayClient)

		if validateMDMError != nil {
			resp.Diagnostics.AddError(
				"Error While Validating MDM Details",
				"unexpected error: "+validateMDMResponse.Message,
			)
			return
		}

		data, dgs := updateState(validateMDMResponse, plan)
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

	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, state, r.gatewayClient)

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
	queueOperationError := QueueOpeartions(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	parsecsvRespose, parseCSVError := ParseCSVOperation(ctx, plan, r.gatewayClient)

	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while Parsing CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	// Vaidate the MDM credentials
	validateMDMResponse, validateMDMError := ValidateMDMOperation(ctx, plan, r.gatewayClient)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Details",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {

		installationError := InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

		if installationError != nil {
			resp.Diagnostics.AddError(
				"Error in Installation Process",
				"unexpected error: "+installationError.Error(),
			)
			return
		}

		validateMDMResponse, validateMDMError = ValidateMDMOperation(ctx, plan, r.gatewayClient)

		if validateMDMError != nil {
			resp.Diagnostics.AddError(
				"Error While Validating MDM Details",
				"unexpected error: "+validateMDMResponse.Message,
			)
			return
		}

		data, dgs := updateState(validateMDMResponse, plan)
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

// Working as PowerFlex Gateway Server Cleanup
func (r *sdcExpansionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CsvAndMdmDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queueOperationError := QueueOpeartions(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// QueueOpeartions function for the Abort, Clear and Move To Idle Execution
func QueueOpeartions(gatewayClient *goscaleio.GatewayClient) error {

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
func ParseCSVOperation(ctx context.Context, plan CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient) (*goscaleio_types.GatewayResponse, error) {

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
	header := []string{"IPs", "Password", "Operating System", "Is MDM/TB", "Is SDC", "perfProfileForSDC", "SDC Name"}
	err = writer.Write(header)
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
	}

	csvItems := []CSVDataModel{}
	diags := plan.CsvDetail.ElementsAs(ctx, &csvItems, true)
	if diags.HasError() {
		return &parseCSVResponse, fmt.Errorf("Error While Parse CSV Data  is %s", diags.Errors())
	}

	for _, item := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			IP:                 item.IP.ValueString(),
			Password:           item.Password.ValueString(),
			IsMdmOrTb:          item.IsMdmOrTb.ValueString(),
			OperatingSystem:    item.OperatingSystem.ValueString(),
			IsSdc:              item.IsSdc.ValueString(),
			PerformanceProfile: item.PerformanceProfile.ValueString(),
			SDCName:            item.SDCName.ValueString(),
		}
		//Write the data row
		data := []string{csvStruct.IP, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc, csvStruct.PerformanceProfile, csvStruct.SDCName}
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
func ValidateMDMOperation(ctx context.Context, plan CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient) (*goscaleio_types.GatewayResponse, error) {
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": plan.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{plan.MdmIP.ValueString()}

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
func InstallationOperations(ctx context.Context, plan CsvAndMdmDataModel, gatewayClient *goscaleio.GatewayClient, parsecsvRespose *goscaleio_types.GatewayResponse) error {
	beginInstallationResponse, _ := gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", plan.MdmPassword.ValueString(), plan.LiaPassword.ValueString(), true)

	if beginInstallationResponse.StatusCode == 200 {
		currentPhase := "query"
		couterForStopExecution := 0

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
						} else if currentPhase == "upload" {
							currentPhase = "install"
						} else if currentPhase == "install" {
							currentPhase = "configure"
						}
					} else {
						return fmt.Errorf("Messsage: %s, Error Code: %s", moveToNextPhaseResponse.Message, strconv.Itoa(moveToNextPhaseResponse.StatusCode))
					}
				} else {
					// to make gateway available for installation
					queueOperationError := QueueOpeartions(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
					}

					couterForStopExecution = 10

					return nil
				}

			} else {
				if checkForPhaseCompleted.Data == "Running" {
					couterForStopExecution++

					if couterForStopExecution == 5 {
						// to make gateway available for installation
						queueOperationError := QueueOpeartions(gatewayClient)
						if queueOperationError != nil {
							return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
						}

						return fmt.Errorf("Time Out,Some Operations of Installer running from since long")
					}

				} else {
					if couterForStopExecution < 5 {

						gatewayClient.RetryPhase()

						couterForStopExecution = 5
					} else {
						// to make gateway available for installation
						queueOperationError := QueueOpeartions(gatewayClient)
						if queueOperationError != nil {
							return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
						}

						return fmt.Errorf("Errors in installation process is %s", queueOperationError.Error())
					}
				}
			}
		}
	} else {
		return fmt.Errorf("Message: %s, Error Code: %s", beginInstallationResponse.Message, strconv.Itoa(beginInstallationResponse.StatusCode))
	}

	return nil
}
