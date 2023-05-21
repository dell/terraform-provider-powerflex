package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CsvAndMdmDataModel struct for CSV Data Processing
type CsvAndMdmDataModel struct {
	ID              types.String `tfsdk:"id"`
	CsvDetail       types.Set    `tfsdk:"csv_detail"`
	MdmIP           types.String `tfsdk:"mdm_ip"`
	MdmPassword     types.String `tfsdk:"mdm_password"`
	LiaPassword     types.String `tfsdk:"lia_password"`
	InstalledSDCIps types.String `tfsdk:"installed_sdc_ips"`
}

// CSVDataModel defines the struct for CSV Parse Data
type CSVDataModel struct {
	IP                 types.String `tfsdk:"ip"`
	Password           types.String `tfsdk:"password"`
	OperatingSystem    types.String `tfsdk:"operating_system"`
	IsMdmOrTb          types.String `tfsdk:"is_mdm_or_tb"`
	IsSdc              types.String `tfsdk:"is_sdc"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	SDCName            types.String `tfsdk:"sdc_name"`
}

// CsvRow desfines the srtuct for the CSV Data
type CsvRow struct {
	IP                 string
	Password           string
	OperatingSystem    string
	IsMdmOrTb          string
	IsSdc              string
	PerformanceProfile string
	SDCName            string
}

// SDCExpansionResourceSchema - varible holds schema for SDC Expansion
var SDCExpansionResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to add the SDC in PowerFlex Cluster.",
	MarkdownDescription: "This resource can be used to add the SDC in PowerFlex Cluster.",
	Attributes: map[string]schema.Attribute{
		"csv_detail": csvSchema,
		"mdm_ip": schema.StringAttribute{
			Description:         "MDM Server IPs. User can provide Primary and Secondary MDM IP comma seperated",
			MarkdownDescription: "MDM Server IPs. User can provide Primary and Secondary MDM IP comma seperated",
			Required:            true,
		},
		"mdm_password": schema.StringAttribute{
			Description:         "MDM Password to connect MDM Server.",
			MarkdownDescription: "MDM Password to connect MDM Server.",
			Required:            true,
		},
		"lia_password": schema.StringAttribute{
			Description:         "LIA Password to connect MDM Server.",
			MarkdownDescription: "LIA Password to connect MDM Server.",
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

// csvSchema - varible holds schema for CSV Param Details
var csvSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDC Expansion Server Details.",
	Required:            true,
	MarkdownDescription: "List of SDC Expansion Server Details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "IP of the node",
				Required:            true,
				MarkdownDescription: "IP of the node",
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
				Optional:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker",
				PlanModifiers: []planmodifier.String{
					stringDefault(" "),
				},
			},
			"is_sdc": schema.StringAttribute{
				Description:         "whether this node is SDC or not",
				Required:            true,
				MarkdownDescription: "whether this node is SDC or not",
			},
			"performance_profile": schema.StringAttribute{
				Description:         "Performance Profile of SDC",
				Optional:            true,
				MarkdownDescription: "Performance Profile of SDC",
				Validators: []validator.String{stringvalidator.OneOf(
					"High",
					"Compact",
				)},
			},
			"sdc_name": schema.StringAttribute{
				Description:         "Name of the SDC",
				Optional:            true,
				MarkdownDescription: "Name of the SDC",
			},
		},
	},
}
