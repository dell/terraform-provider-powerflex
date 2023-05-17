package powerflex

import (
	"context"
	"time"
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSdcExpansionResource() resource.Resource {
	return &sdcExpansionResource{}
}

type sdcExpansionResource struct {
	gatewayClient *goscaleio.GatewayClient
}

type CsvAndMdmDataModel struct {
	CsvDetail    types.Set    `tfsdk:"csv_detail"`
	CsvComplete types.String `tfsdk:"csv_complete"`
	MdmIp        types.String `tfsdk:"mdm_ip"`
	MdmPassword  types.String `tfsdk:"mdm_password"`
	LiaPassword  types.String `tfsdk:"lia_password"`
}

// UploadPackageModel defines the struct for device resource
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
			"csv_complete": schema.StringAttribute{
				Description:         "The JSON data which is being received after parsing the csv.",
				MarkdownDescription: "The JSON data which is being received after parsing the csv.",
				Computed:            true,
			},
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

func (r *sdcExpansionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CsvAndMdmDataModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// to make gateway available for installation
	r.gatewayClient.AbortOperation()
	r.gatewayClient.ClearQueueCommand()
	r.gatewayClient.MoveToIdlePhase()

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Create a csv writer
	file, err := os.Create(mydir + "/Minimal.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"IPs", "Password", "Operating System", "Is MDM/TB", "Is SDC"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}
	csvItems := []CSVDataModel{}
	diags = plan.CsvDetail.ElementsAs(ctx, &csvItems, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			Ip:              si.Ip.ValueString(),
			Password:        si.Password.ValueString(),
			IsMdmOrTb:       si.IsMdmOrTb.ValueString(),
			OperatingSystem: si.OperatingSystem.ValueString(),
			IsSdc:           si.IsSdc.ValueString(),
		}
		//Write the data row
		data := []string{csvStruct.Ip, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc}
		err = writer.Write(data)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	parsecsvRespose, parseCSVError := r.gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while parsing the csv: ",
			"unexpected error: "+parsecsvRespose.Message,
		)
		//plan.CsvCompltete = types.StringValue("Not completed")
		return
	}
	if parsecsvRespose.StatusCode == 200 {
		plan.CsvComplete = types.StringValue("Completed")
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error parsing csv :"+parsecsvRespose.Message+" & Error Code :"+strconv.Itoa(parsecsvRespose.ErrorCode),
			"Status Code:"+strconv.Itoa(parsecsvRespose.StatusCode),
		)
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
						r.gatewayClient.AbortOperation()
						r.gatewayClient.ClearQueueCommand()
						r.gatewayClient.MoveToIdlePhase()
						couterForStopExecution = 10
					}

				} else {
					if checkForPhaseCompleted.Data == "Running" {
						couterForStopExecution++
					} else {
						if couterForStopExecution < 5 {
							r.gatewayClient.RetryPhase()
							couterForStopExecution = 5
						} else {
							r.gatewayClient.AbortOperation()
							r.gatewayClient.ClearQueueCommand()
							r.gatewayClient.MoveToIdlePhase()
							resp.Diagnostics.AddError("Errors in installation process",
								"Errors:"+checkForPhaseCompleted.Message)
							couterForStopExecution++
						}
					}
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Error in begin installation :"+beginInstallationResponse.Message+" & Error Code :"+strconv.Itoa(beginInstallationResponse.ErrorCode),
				"Status Code:"+strconv.Itoa(beginInstallationResponse.StatusCode),
			)
		}

	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error while validating mdm credentials :"+validateMDMResponse.Message+" & Error Code :"+strconv.Itoa(validateMDMResponse.ErrorCode),
			"Status Code:"+strconv.Itoa(validateMDMResponse.StatusCode),
		)
	}
	
}

func (r *sdcExpansionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *sdcExpansionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CsvAndMdmDataModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// to make gateway available for installation
	r.gatewayClient.AbortOperation()
	r.gatewayClient.ClearQueueCommand()
	r.gatewayClient.MoveToIdlePhase()

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Create a csv writer
	file, err := os.Create(mydir + "/Minimal.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"IPs", "Password", "Operating System", "Is MDM/TB", "Is SDC"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}
	csvItems := []CSVDataModel{}
	diags = plan.CsvDetail.ElementsAs(ctx, &csvItems, true)
	resp.Diagnostics.Append(diags...)
	for _, si := range csvItems {
		// Add mapped SDC
		csvStruct := CsvRow{
			Ip:              si.Ip.ValueString(),
			Password:        si.Password.ValueString(),
			IsMdmOrTb:       si.IsMdmOrTb.ValueString(),
			OperatingSystem: si.OperatingSystem.ValueString(),
			IsSdc:           si.IsSdc.ValueString(),
		}
		//Write the data row
		data := []string{csvStruct.Ip, csvStruct.Password, csvStruct.OperatingSystem, csvStruct.IsMdmOrTb, csvStruct.IsSdc}
		err = writer.Write(data)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	parsecsvRespose, parseCSVError := r.gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if parseCSVError != nil {
		resp.Diagnostics.AddError(
			"Error while parsing the csv: ",
			"unexpected error: "+parsecsvRespose.Message,
		)
		//plan.CsvCompltete = types.StringValue("Not completed")
		return
	}
	if parsecsvRespose.StatusCode == 200 {
		plan.CsvComplete = types.StringValue("Completed")
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error parsing csv :"+parsecsvRespose.Message+" & Error Code :"+strconv.Itoa(parsecsvRespose.ErrorCode),
			"Status Code:"+strconv.Itoa(parsecsvRespose.StatusCode),
		)
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
						r.gatewayClient.AbortOperation()
						r.gatewayClient.ClearQueueCommand()
						r.gatewayClient.MoveToIdlePhase()
						couterForStopExecution = 10
					}

				} else {
					if checkForPhaseCompleted.Data == "Running" {
						couterForStopExecution++
					} else {
						if couterForStopExecution < 5 {
							r.gatewayClient.RetryPhase()
							couterForStopExecution = 5
						} else {
							r.gatewayClient.AbortOperation()
							r.gatewayClient.ClearQueueCommand()
							r.gatewayClient.MoveToIdlePhase()
							resp.Diagnostics.AddError("Errors in installation process",
								"Errors:"+checkForPhaseCompleted.Message)
							couterForStopExecution++
						}
					}
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Error in begin installation :"+beginInstallationResponse.Message+" & Error Code :"+strconv.Itoa(beginInstallationResponse.ErrorCode),
				"Status Code:"+strconv.Itoa(beginInstallationResponse.StatusCode),
			)
		}

	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error while validating mdm credentials :"+validateMDMResponse.Message+" & Error Code :"+strconv.Itoa(validateMDMResponse.ErrorCode),
			"Status Code:"+strconv.Itoa(validateMDMResponse.StatusCode),
		)
	}

}

func (r *sdcExpansionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CsvAndMdmDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.gatewayClient.AbortOperation()
	r.gatewayClient.ClearQueueCommand()
	r.gatewayClient.MoveToIdlePhase()
	resp.State.RemoveResource(ctx)
}
