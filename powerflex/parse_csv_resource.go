package powerflex

import (
	"context"
	//"time"
	"time"
	//"fmt"

	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	//goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewParsecsvResource() resource.Resource {
	return &parsecsvResource{}
}

type parsecsvResource struct {
	gatewayClient *goscaleio.GatewayClient
}

type CsvListModel struct {
	CsvDetail    types.Set    `tfsdk:"csv_detail"`
	CsvCompltete types.String `tfsdk:"csv_complete"`
	MdmIp        types.String `tfsdk:"mdm_ip"`
	MdmPassword  types.String `tfsdk:"mdm_password"`
	LiaPassword  types.String `tfsdk:"lia_password"`
}

// UploadPackageModel defines the struct for device resource
type ParseCSVModel struct {
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

type PhaseDetail struct {
	Phase CurrentPhase `json:"phase"`
}

type CurrentPhase struct {
	Name string `json:"name"`
}

func (r *parsecsvResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parse_csv"
}

func (r *parsecsvResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to parse csv.",
		MarkdownDescription: "This resource can be used to parse csv.",
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

func (r *parsecsvResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.gatewayClient = req.ProviderData.(*goscaleio.GatewayClient)
}

func (r *parsecsvResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CsvListModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	csvItems := []ParseCSVModel{}
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
	parsecsvRespose, err3 := r.gatewayClient.ParseCSV(mydir + "/Minimal.csv")
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error parsing csv: ",
			"unexpected error: "+err3.Error(),
		)
		//plan.CsvCompltete = types.StringValue("Not completed")
		return
	}
	if parsecsvRespose.StatusCode == 200 {
		plan.CsvCompltete = types.StringValue("Completed")
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error parsing csv :"+parsecsvRespose.Message+" & Error Code :"+strconv.Itoa(parsecsvRespose.ErrorCode),
			"Status Code:"+strconv.Itoa(parsecsvRespose.StatusCode),
		)
	}

	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": plan.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = []string{plan.MdmIp.ValueString()}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": false,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData
	jsonres, _ := json.Marshal(mapData)
	response, err4 := r.gatewayClient.ValidateMDMDetails(jsonres)
	if err4 != nil {
		resp.Diagnostics.AddError(
			"Error validating details: ",
			"unexpected error: "+err3.Error(),
		)
		return
	}
	if response.StatusCode == 200 {

		//var phaseData PhaseDetail
		response3, _ := r.gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", plan.MdmPassword.ValueString(), plan.LiaPassword.ValueString(),true)
		// phaseStatusResponseBefore,_:= r.gatewayClient.GetInstallerPhaseDetails()
		// json.Unmarshal([]byte(phaseStatusResponseBefore.Message), &phaseData)
		// currentPhaseName := phaseData.Phase.Name
		// resp.Diagnostics.AddWarning(
		// 	"current pahse name :"+currentPhaseName,
		// 	currentPhaseName,
		// )
		 if response3.StatusCode == 0 {
			time.Sleep(7*time.Minute)
			moveToNextPhaseResponse, _ := r.gatewayClient.MoveToNextPhase()
			time.Sleep(15*time.Second)
			moveToNextPhaseResponse2, _ := r.gatewayClient.MoveToNextPhase()
			time.Sleep(1*time.Minute)
			moveToNextPhaseResponse3, _ := r.gatewayClient.MoveToNextPhase()
			if moveToNextPhaseResponse.StatusCode == 200 {
				plan.CsvCompltete = types.StringValue("Completed")
				diags = resp.State.Set(ctx, &plan)
				resp.Diagnostics.Append(diags...)
			} else {
				resp.Diagnostics.AddError(
					"Error in moving to next phase :"+moveToNextPhaseResponse.Message+" & Error Code :"+strconv.Itoa(moveToNextPhaseResponse.ErrorCode),
					"Status Code:"+strconv.Itoa(response3.StatusCode),
				)
			}
			if moveToNextPhaseResponse2.StatusCode == 200 {
				plan.CsvCompltete = types.StringValue("Completed")
				diags = resp.State.Set(ctx, &plan)
				resp.Diagnostics.Append(diags...)
			} else {
				resp.Diagnostics.AddError(
					"Error in moving to next phase :"+moveToNextPhaseResponse.Message+" & Error Code :"+strconv.Itoa(moveToNextPhaseResponse.ErrorCode),
					"Status Code:"+strconv.Itoa(response3.StatusCode),
				)
			}
			if moveToNextPhaseResponse3.StatusCode == 200 {
				plan.CsvCompltete = types.StringValue("Completed")
				diags = resp.State.Set(ctx, &plan)
				resp.Diagnostics.Append(diags...)
			} else {
				resp.Diagnostics.AddError(
					"Error in moving to next phase :"+moveToNextPhaseResponse.Message+" & Error Code :"+strconv.Itoa(moveToNextPhaseResponse.ErrorCode),
					"Status Code:"+strconv.Itoa(response3.StatusCode),
				)
			}
		
			// i:=0
			// details, _ := r.gatewayClient.GetInstallerPhaseDetails()
			// currentPhase := details.NextPhase.Name
			// for i<3{
			// 	time.Sleep(5*time.Second)
			// 	phaseAfterFiveSec,_:= r.gatewayClient.GetInstallerPhaseDetails()
			// 	if currentPhase != phaseAfterFiveSec.Phase.Name {
			// 		i= i+1
			// 		moveToNextPhaseResponse, _ := r.gatewayClient.MoveToNextPhase()
			// 		currentPhase = phaseAfterFiveSec.Phase.Name
			// 		if moveToNextPhaseResponse.StatusCode == 200 {
			// 			plan.CsvCompltete = types.StringValue("Completed")
			// 			diags = resp.State.Set(ctx, &plan)
			// 			resp.Diagnostics.Append(diags...)
			// 		} else {
			// 			resp.Diagnostics.AddError(
			// 				"Error in moving to next phase :"+moveToNextPhaseResponse.Message+" & Error Code :"+strconv.Itoa(moveToNextPhaseResponse.ErrorCode),
			// 				"Status Code:"+strconv.Itoa(response3.StatusCode),
			// 			)
			// 			break
			// 		}
					
			// 	}

			// }
			
		} else {

			resp.Diagnostics.AddError(
				"Error in begin installation :"+response3.Message+" & Error Code :"+strconv.Itoa(response3.ErrorCode),
				"Status Code:"+strconv.Itoa(response3.StatusCode),
			)
		}

	} else {
		//plan.CsvCompltete = types.StringValue("Unsuccessful")
		resp.Diagnostics.AddError(
			"Error while validating mdm credentials :"+response.Message+" & Error Code :"+strconv.Itoa(response.ErrorCode),
			"Status Code:"+strconv.Itoa(response.StatusCode),
		)
	}

}

func (r *parsecsvResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *parsecsvResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *parsecsvResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
