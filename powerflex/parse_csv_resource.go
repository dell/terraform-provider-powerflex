package powerflex

import (
	"context"
	"fmt"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"os"
	"encoding/csv"
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
	CsvDetail              types.Set    `tfsdk:"csv_detail"`
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

func (r *parsecsvResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parse_csv"
}

func (r *parsecsvResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to parse csv.",
		MarkdownDescription: "This resource can be used to parse csv.",
		Attributes: map[string]schema.Attribute{
			"csv_detail": csvSchema,
		},
	}
}

var csvSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	Computed:            true,
	Optional:            true,
	MarkdownDescription: "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	DeprecationMessage:  "Please use sdc_volumes_mapping resource for mapping SDCs to the volumes/snapshots. This attribute will be removed in future release.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "ip of the node",
				Required: true,
				MarkdownDescription: "ip of the node",
			},
			"password": schema.StringAttribute{
				Description:         "The Path of the directory of csv file",
				Required:            true,
				MarkdownDescription: "The Path of the directory of csv file",
			},
			"operating_system": schema.StringAttribute{
				Description:         "The Path of the directory of csv file",
				Required:            true,
				MarkdownDescription: "The Path of the directory of csv file",
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "The Path of the directory of csv file",
				Required:            true,
				MarkdownDescription: "The Path of the directory of csv file",
			},
			"is_sdc": schema.StringAttribute{
				Description:         "The Path of the directory of csv file",
				Required:            true,
				MarkdownDescription: "The Path of the directory of csv file",
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
	
	// Create a csv writer
	file, err := os.Create("Minimal.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"IPs","Password","Operating System","Is MDM/TB","Is SDC"}
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
			Ip: fmt.Sprintf("%v",si.Ip),
			Password: fmt.Sprintf("%v",si.Password),
			IsMdmOrTb: fmt.Sprintf("%v",si.IsMdmOrTb),
			OperatingSystem: fmt.Sprintf("%v",si.OperatingSystem),
			IsSdc: fmt.Sprintf("%v",si.IsSdc),

		}
		//Write the data row
		data := []string{csvStruct.Ip,csvStruct.Password,csvStruct.OperatingSystem,csvStruct.IsMdmOrTb,csvStruct.IsSdc}
		err = writer.Write(data)
		if err != nil {
			panic(err)
		}
	}

	err3 := r.gatewayClient.ParseCSV("Minimal.csv")
	// if err3 != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error getting with file path: "+plan.FilePath.ValueString(),
	// 		"unexpected error: "+err3.Error(),
	// 	)
	// 	return
	// }

	// if response.StatusCode == 200 {
	// 	res, err3 := r.gatewayClient.GetPackgeDetails()
	// 	if err3 != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error getting pacckage details:",
	// 			"unexpected error: "+err3.Error(),
	// 		)
	// 		return
	// 	}

		// Set refreshed state
		// data, dgs := updateUploadPackageState(res, plan)
		// resp.Diagnostics.Append(dgs...)

		// diags = resp.State.Set(ctx, data)
		// resp.Diagnostics.Append(diags...)
	//}
}

func (r *parsecsvResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *parsecsvResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *parsecsvResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
