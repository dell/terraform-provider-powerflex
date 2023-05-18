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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSDCExpansionResource() resource.Resource {
	return &sdcExpansionResource{}
}

type sdcExpansionResource struct {
	gatewayClient *goscaleio.GatewayClient
}

type CsvAndMdmDataModel struct {
	ID              types.String `tfsdk:"id"`
	CsvDetail       types.Set    `tfsdk:"csv_detail"`
	MdmIp           types.String `tfsdk:"mdm_ip"`
	MdmPassword     types.String `tfsdk:"mdm_password"`
	LiaPassword     types.String `tfsdk:"lia_password"`
	InstalledSDCIps types.String `tfsdk:"installed_sdc_ips"`
}

// CSVDataModel defines the struct for CSV Parse Data
type CSVDataModel struct {
	Ip              types.String `tfsdk:"ip"`
	Password        types.String `tfsdk:"password"`
	OperatingSystem types.String `tfsdk:"operating_system"`
	IsMdmOrTb       types.String `tfsdk:"is_mdm_or_tb"`
	IsSdc           types.String `tfsdk:"is_sdc"`
}

type CsvRow struct {
	Ip              string
	Password        string
	OperatingSystem string
	IsMdmOrTb       string
	IsSdc           string
}

func (r *sdcExpansionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_expansion"
}

func (r *sdcExpansionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to add the SDC.",
		MarkdownDescription: "This resource can be used to add the SDC.",
		Attributes: map[string]schema.Attribute{
			"csv_detail": csvSchema,
			"mdm_ip": schema.StringAttribute{
				Description:         "The JSON data which is being received after parsing the csv.",
				MarkdownDescription: "The JSON data which is being received after parsing the csv.",
				Required:            true,
			},
			"mdm_password": schema.StringAttribute{
				Description:         "The JSON data which is being received after parsing the csv.",
				MarkdownDescription: "The JSON data which is being received after parsing the csv.",
				Required:            true,
			},
			"lia_password": schema.StringAttribute{
				Description:         "The JSON data which is being received after parsing the csv.",
				MarkdownDescription: "The JSON data which is being received after parsing the csv.",
				Required:            true,
			},
			"installed_sdc_ips": schema.StringAttribute{
				Description:         "List of installed SDC IPs",
				Computed:            true,
				MarkdownDescription: "List of installed SDC IPs",
			},
			"id": schema.StringAttribute{
				Description:         "The ID of the package.",
				Computed:            true,
				MarkdownDescription: "The ID of the package.",
			},
		},
	}
}

var csvSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	Required:            true,
	MarkdownDescription: "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "ip of the node",
				Required:            true,
				MarkdownDescription: "ip of the node",
			},
			"password": schema.StringAttribute{
				Description:         "Password on the node",
				Required:            true,
				MarkdownDescription: "Password on the node",
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Required:            true,
				MarkdownDescription: "Operating System on the node",
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "Whether this works as MDM or Tie Breaker",
				Required:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker",
			},
			"is_sdc": schema.StringAttribute{
				Description:         "whether this node is SDC or not",
				Required:            true,
				MarkdownDescription: "whether this node is SDC or not",
			},
		},
	},
}

func (r *sdcExpansionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.gatewayClient = req.ProviderData.(*goscaleio.GatewayClient)
}

func QueueOpeartions(gatewayClient *goscaleio.GatewayClient) error {

	_, err := gatewayClient.AbortOperation()

	if err != nil {
		return fmt.Errorf("Error while Aborting Operation: %s", err.Error())
	}
	_, err = gatewayClient.ClearQueueCommand()

	if err != nil {
		return fmt.Errorf("Error while Clearing Queue: %s", err.Error())
	}

	_, err = gatewayClient.MoveToIdlePhase()

	if err != nil {
		return fmt.Errorf("Error while Move to Ideal Phase: %s", err.Error())
	}

	return nil
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

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Readin Current Directory",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}
	// Create a csv writer
	file, err := os.Create(mydir + "/Minimal.csv")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Creating Temp CSV",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"IPs", "Password", "Operating System", "Is MDM/TB", "Is SDC"}
	err = writer.Write(header)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Writing Temp CSV",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	csvItems := []CSVDataModel{}
	diags = plan.CsvDetail.ElementsAs(ctx, &csvItems, true)
	resp.Diagnostics.Append(diags...)

	for _, item := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			Ip:              item.Ip.ValueString(),
			Password:        item.Password.ValueString(),
			IsMdmOrTb:       item.IsMdmOrTb.ValueString(),
			OperatingSystem: item.OperatingSystem.ValueString(),
			IsSdc:           item.IsSdc.ValueString(),
		}
		//Write the data row
		data := []string{csvStruct.Ip, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc}
		err = writer.Write(data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error While Creating Temp CSV File",
				"unexpected error: "+queueOperationError.Error(),
			)
			return
		}
	}
	writer.Flush()

	parsecsvRespose, parseCSVError := r.gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error While Parsing the CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	if parsecsvRespose.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error While Parsing CSV",
			"Error While Parsing CSV :"+parsecsvRespose.Message+" & Status Code :"+strconv.Itoa(parsecsvRespose.StatusCode),
		)
		return
	}

	// Vaidate the MDM credentials
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": plan.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{plan.MdmIp.ValueString()}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData
	jsonres, _ := json.Marshal(mapData)
	validateMDMResponse, validateMDMError := r.gatewayClient.ValidateMDMDetails(jsonres)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error validating details: ",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {
		// begin instllation process
		beginInstallationResponse, _ := r.gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", plan.MdmPassword.ValueString(), plan.LiaPassword.ValueString(), true)

		if beginInstallationResponse.StatusCode == 200 {
			currentPhase := "query"
			couterForStopExecution := 0

			for couterForStopExecution <= 5 {

				time.Sleep(1 * time.Minute)

				checkForPhaseCompleted, _ := r.gatewayClient.CheckForCompletionQueueCommands(currentPhase)

				if checkForPhaseCompleted.Data == "Completed" {
					couterForStopExecution = 0

					if currentPhase != "configure" {
						moveToNextPhaseResponse, err := r.gatewayClient.MoveToNextPhase()

						if err != nil {
							resp.Diagnostics.AddError(
								"Error while moving to next phase",
								"unexpected error: "+moveToNextPhaseResponse.Message,
							)
							return
						}
						if moveToNextPhaseResponse.StatusCode == 200 {
							if currentPhase == "query" {
								currentPhase = "upload"
							} else if currentPhase == "upload" {
								currentPhase = "install"
							} else if currentPhase == "install" {
								currentPhase = "configure"
							}
						}
					} else {
						// to make gateway available for installation
						queueOperationError := QueueOpeartions(r.gatewayClient)
						if queueOperationError != nil {
							resp.Diagnostics.AddError(
								"Error Clearing Queue",
								"unexpected error: "+queueOperationError.Error(),
							)
							return
						}

						couterForStopExecution = 10

						//Fetching Latest Data Updating to State
						validateMDMResponse, validateMDMError := r.gatewayClient.ValidateMDMDetails(jsonres)
						if validateMDMError != nil {
							resp.Diagnostics.AddError(
								"Error validating details: ",
								"unexpected error: "+validateMDMResponse.Message,
							)
							return
						}

						data, dgs := updateState(validateMDMResponse, plan)
						resp.Diagnostics.Append(dgs...)

						diags = resp.State.Set(ctx, data)
						resp.Diagnostics.Append(diags...)
					}

				} else {
					if checkForPhaseCompleted.Data == "Running" {
						couterForStopExecution++
					} else {
						if couterForStopExecution < 5 {
							r.gatewayClient.RetryPhase()
							couterForStopExecution = 5
						} else {
							// to make gateway available for installation
							queueOperationError := QueueOpeartions(r.gatewayClient)
							if queueOperationError != nil {
								resp.Diagnostics.AddError(
									"Error Clearing Queue",
									"unexpected error: "+queueOperationError.Error(),
								)
								return
							}

							resp.Diagnostics.AddError("Errors in installation process",
								"Errors:"+checkForPhaseCompleted.Message)

							return
						}
					}
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Error in begin installation",
				"Error in begin installation :"+beginInstallationResponse.Message+" & Status Code :"+strconv.Itoa(beginInstallationResponse.StatusCode),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Credentials",
			"Error While Validating MDM Credentials: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
		)
		return
	}

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

	// Vaidate the MDM credentials
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": state.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{state.MdmIp.ValueString()}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData

	jsonres, jsonerr := json.Marshal(mapData)
	if jsonerr != nil {
		resp.Diagnostics.AddError(
			"Error while validating MDM inputs",
			"unexpected error: "+jsonerr.Error(),
		)
		return
	}

	validateMDMResponse, validateMDMError := r.gatewayClient.ValidateMDMDetails(jsonres)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error validating details: ",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {
		// Set refreshed state
		data, dgs := updateState(validateMDMResponse, state)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.AddError(
		"Error While Validating MDM Details",
		"Error Meesage: "+validateMDMResponse.Message+" Error Code:"+strconv.FormatInt(int64(validateMDMResponse.StatusCode), 10),
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

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Readin Current Directory",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}
	// Create a csv writer
	file, err := os.Create(mydir + "/Minimal.csv")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Creating Temp CSV",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"IPs", "Password", "Operating System", "Is MDM/TB", "Is SDC"}
	err = writer.Write(header)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error While Writing Temp CSV",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	csvItems := []CSVDataModel{}
	diags = plan.CsvDetail.ElementsAs(ctx, &csvItems, true)
	resp.Diagnostics.Append(diags...)

	for _, item := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			Ip:              item.Ip.ValueString(),
			Password:        item.Password.ValueString(),
			IsMdmOrTb:       item.IsMdmOrTb.ValueString(),
			OperatingSystem: item.OperatingSystem.ValueString(),
			IsSdc:           item.IsSdc.ValueString(),
		}
		//Write the data row
		data := []string{csvStruct.Ip, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc}
		err = writer.Write(data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error While Creating Temp CSV File",
				"unexpected error: "+queueOperationError.Error(),
			)
			return
		}
	}
	writer.Flush()

	parsecsvRespose, parseCSVError := r.gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error While Parsing the CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	if parsecsvRespose.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error While Parsing CSV",
			"Error While Parsing CSV :"+parsecsvRespose.Message+" & Status Code :"+strconv.Itoa(parsecsvRespose.StatusCode),
		)
		return
	}

	// Vaidate the MDM credentials
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": plan.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{plan.MdmIp.ValueString()}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData
	jsonres, _ := json.Marshal(mapData)
	validateMDMResponse, validateMDMError := r.gatewayClient.ValidateMDMDetails(jsonres)
	if validateMDMError != nil {
		resp.Diagnostics.AddError(
			"Error validating details: ",
			"unexpected error: "+validateMDMResponse.Message,
		)
		return
	}

	if validateMDMResponse.StatusCode == 200 {
		// begin instllation process
		beginInstallationResponse, _ := r.gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", plan.MdmPassword.ValueString(), plan.LiaPassword.ValueString(), true)

		if beginInstallationResponse.StatusCode == 200 {
			currentPhase := "query"
			couterForStopExecution := 0

			for couterForStopExecution <= 5 {

				time.Sleep(1 * time.Minute)

				checkForPhaseCompleted, _ := r.gatewayClient.CheckForCompletionQueueCommands(currentPhase)

				if checkForPhaseCompleted.Data == "Completed" {
					couterForStopExecution = 0

					if currentPhase != "configure" {
						moveToNextPhaseResponse, err := r.gatewayClient.MoveToNextPhase()

						if err != nil {
							resp.Diagnostics.AddError(
								"Error while moving to next phase",
								"unexpected error: "+moveToNextPhaseResponse.Message,
							)
							return
						}
						if moveToNextPhaseResponse.StatusCode == 200 {
							if currentPhase == "query" {
								currentPhase = "upload"
							} else if currentPhase == "upload" {
								currentPhase = "install"
							} else if currentPhase == "install" {
								currentPhase = "configure"
							}
						}
					} else {
						// to make gateway available for installation
						queueOperationError := QueueOpeartions(r.gatewayClient)
						if queueOperationError != nil {
							resp.Diagnostics.AddError(
								"Error Clearing Queue",
								"unexpected error: "+queueOperationError.Error(),
							)
							return
						}

						couterForStopExecution = 10

						//Fetching Latest Data Updating to State
						validateMDMResponse, validateMDMError := r.gatewayClient.ValidateMDMDetails(jsonres)
						if validateMDMError != nil {
							resp.Diagnostics.AddError(
								"Error validating details: ",
								"unexpected error: "+validateMDMResponse.Message,
							)
							return
						}

						data, dgs := updateState(validateMDMResponse, plan)
						resp.Diagnostics.Append(dgs...)

						diags = resp.State.Set(ctx, data)
						resp.Diagnostics.Append(diags...)
					}

				} else {
					if checkForPhaseCompleted.Data == "Running" {
						couterForStopExecution++
					} else {
						if couterForStopExecution < 5 {
							r.gatewayClient.RetryPhase()
							couterForStopExecution = 5
						} else {
							// to make gateway available for installation
							queueOperationError := QueueOpeartions(r.gatewayClient)
							if queueOperationError != nil {
								resp.Diagnostics.AddError(
									"Error Clearing Queue",
									"unexpected error: "+queueOperationError.Error(),
								)
								return
							}

							resp.Diagnostics.AddError("Errors in installation process",
								"Errors:"+checkForPhaseCompleted.Message)

							return
						}
					}
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Error in begin installation",
				"Error in begin installation :"+beginInstallationResponse.Message+" & Status Code :"+strconv.Itoa(beginInstallationResponse.StatusCode),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Error While Validating MDM Credentials",
			"Error While Validating MDM Credentials: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
		)
		return
	}

}

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
