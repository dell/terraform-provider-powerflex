package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CsvAndMdmDataModel struct for CSV Data Processing
type CsvAndMdmDataModel struct {
	ID              types.String `tfsdk:"id"`
	ClusterDetails  types.Set    `tfsdk:"cluster_details"`
	MdmPassword     types.String `tfsdk:"mdm_password"`
	LiaPassword     types.String `tfsdk:"lia_password"`
	InstalledSDCIps types.String `tfsdk:"installed_sdc_ips"`
}

// CSVDataModel defines the struct for CSV Parse Data
type CSVDataModel struct {
	IP                 types.String `tfsdk:"ip"`
	UserName           types.String `tfsdk:"username"`
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
	UserName           string
	Password           string
	OperatingSystem    string
	IsMdmOrTb          string
	IsSdc              string
	PerformanceProfile string
	SDCName            string
}

// SDCExpansionResourceSchema - variable holds schema for SDC Expansion
var SDCExpansionResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to add the SDC in PowerFlex Cluster.",
	MarkdownDescription: "This resource can be used to add the SDC in PowerFlex Cluster.",
	Attributes: map[string]schema.Attribute{
		"cluster_details": csvSchema,
		"mdm_password": schema.StringAttribute{
			Description:         "MDM Password to connect MDM Server.",
			MarkdownDescription: "MDM Password to connect MDM Server.",
			Required:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"lia_password": schema.StringAttribute{
			Description:         "LIA Password to connect MDM Server.",
			MarkdownDescription: "LIA Password to connect MDM Server.",
			Required:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"installed_sdc_ips": schema.StringAttribute{
			Description:         "List of installed SDC IPs",
			Computed:            true,
			MarkdownDescription: "List of installed SDC IPs",
		},
		"id": schema.StringAttribute{
			Description: "The ID of the package.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "The ID of the package.",
		},
	},
}

// csvSchema - variable holds schema for CSV Param Details
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
			"username": schema.StringAttribute{
				Description:         "Username of the node",
				Optional:            true,
				MarkdownDescription: "Username of the node",
				PlanModifiers: []planmodifier.String{
					stringDefault("root"),
				},
			},
			"password": schema.StringAttribute{
				Description:         "Password of the node",
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Password of the node",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Required:            true,
				MarkdownDescription: "Operating System on the node",
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "Whether this works as MDM or Tie Breaker,The acceptable value is `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
				Optional:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker,The acceptable value is `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
				PlanModifiers: []planmodifier.String{
					stringDefault(" "),
				},
			},
			"is_sdc": schema.StringAttribute{
				Description:         "whether this node is SDC or not,The acceptable value is `Yes` or `No`",
				Required:            true,
				MarkdownDescription: "whether this node is SDC or not,The acceptable value is `Yes` or `No`.",
				Validators: []validator.String{stringvalidator.OneOf(
					"Yes",
					"No",
				)},
			},
			"performance_profile": schema.StringAttribute{
				Description:         "Performance Profile of SDC, The acceptable value is `High` or `Compact`.",
				Optional:            true,
				MarkdownDescription: "Performance Profile of SDC, The acceptable value is `High` or `Compact`.",
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
