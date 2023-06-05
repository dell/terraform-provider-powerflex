package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var sdcDetailsDescriptions = struct {
	SdcResourceSchema  string
	LastUpdated        string
	ID                 string
	SystemID           string
	Name               string
	SdcIP              string
	SdcApproved        string
	OnVMWare           string
	SdcGUID            string
	MdmConnectionState string
	Links              string
	LinksRel           string
	LinksHref          string
}{
	SdcResourceSchema:  "This resource can be used to manage Storage Data Clients on a PowerFlex array.",
	LastUpdated:        "The Last updated timestamp of the SDC.",
	ID:                 "ID of the SDC to manage. This can be retrieved from the PowerFlex website. Cannot be updated.",
	SystemID:           "The System ID of the fetched SDC.",
	Name:               "Name of the SDC to manage.",
	SdcIP:              "The IP of the fetched SDC.",
	SdcApproved:        "If the fetched SDC is approved.",
	OnVMWare:           "If the fetched SDC is on vmware.",
	SdcGUID:            "The GUID of the fetched SDC.",
	MdmConnectionState: "The MDM connection status of the fetched SDC.",
	Links:              "The Links of the fetched SDC.",
	LinksRel:           "The Links-Rel of the fetched SDC.",
	LinksHref:          "The Links-HREF of the fetched SDC.",
}

// SDCDataModel struct for CSV Data Processing
type SDCDataModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ClusterDetails types.Set    `tfsdk:"cluster_details"`
	MdmPassword    types.String `tfsdk:"mdm_password"`
	LiaPassword    types.String `tfsdk:"lia_password"`
	InstalledSDCs  types.Set    `tfsdk:"installed_sdcs"`
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

// SDCManagerResourceSchema - variable holds schema for SDC Management
var SDCManagerResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to add the SDC in PowerFlex Cluster.",
	MarkdownDescription: "This resource can be used to add the SDC in PowerFlex Cluster.",
	Attributes: map[string]schema.Attribute{
		"cluster_details": csvSchema,
		"name": schema.StringAttribute{
			Description: sdcDetailsDescriptions.Name,
			Optional:    true,
		},
		"mdm_password": schema.StringAttribute{
			Description:         "MDM Password to connect MDM Server.",
			MarkdownDescription: "MDM Password to connect MDM Server.",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"lia_password": schema.StringAttribute{
			Description:         "LIA Password to connect MDM Server.",
			MarkdownDescription: "LIA Password to connect MDM Server.",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"installed_sdcs": sdcDetailsSchema,
		"id": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: sdcDetailsDescriptions.ID,
		},
	},
}

// sdcDetailsSchema - variable holds sdc Details
var sdcDetailsSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDC Details.",
	Optional:            true,
	MarkdownDescription: "List of SDC Details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: sdcDetailsDescriptions.ID,
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: sdcDetailsDescriptions.LastUpdated,
			},
			"name": schema.StringAttribute{
				Description: sdcDetailsDescriptions.Name,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sdc_guid": schema.StringAttribute{
				Description: sdcDetailsDescriptions.SdcGUID,
				Computed:    true,
			},
			"on_vmware": schema.BoolAttribute{
				Description: sdcDetailsDescriptions.OnVMWare,
				Computed:    true,
			},
			"sdc_approved": schema.BoolAttribute{
				Description: sdcDetailsDescriptions.SdcApproved,
				Computed:    true,
			},
			"system_id": schema.StringAttribute{
				Description: sdcDetailsDescriptions.SystemID,
				Computed:    true,
			},
			"sdc_ip": schema.StringAttribute{
				Description: sdcDetailsDescriptions.SdcIP,
				Computed:    true,
			},
			"mdm_connection_state": schema.StringAttribute{
				Description: sdcDetailsDescriptions.MdmConnectionState,
				Computed:    true,
			},
			"links": schema.ListNestedAttribute{
				Description: sdcDetailsDescriptions.Links,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rel": schema.StringAttribute{
							Description: sdcDetailsDescriptions.LinksRel,
							Computed:    true,
						},
						"href": schema.StringAttribute{
							Description: sdcDetailsDescriptions.LinksHref,
							Computed:    true,
						},
					},
				},
			},
		},
	},
}

// csvSchema - variable holds schema for CSV Param Details
var csvSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDC Expansion Server Details.",
	Optional:            true,
	MarkdownDescription: "List of SDC Expansion Server Details.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description:         "IP of the node",
				Optional:            true,
				MarkdownDescription: "IP of the node",
			},
			"username": schema.StringAttribute{
				Description:         "Username of the node",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Username of the node",
				PlanModifiers: []planmodifier.String{
					stringDefault("root"),
				},
			},
			"password": schema.StringAttribute{
				Description:         "Password of the node",
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password of the node",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"operating_system": schema.StringAttribute{
				Description:         "Operating System on the node",
				Optional:            true,
				MarkdownDescription: "Operating System on the node",
			},
			"is_mdm_or_tb": schema.StringAttribute{
				Description:         "Whether this works as MDM or Tie Breaker,The acceptable value is `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
				Optional:            true,
				MarkdownDescription: "Whether this works as MDM or Tie Breaker,The acceptable value is `Primary`, `Secondary`, `TB`, `Standby` or blank. Default value is blank",
			},
			"is_sdc": schema.StringAttribute{
				Description:         "whether this node is SDC or not,The acceptable value is `Yes` or `No`",
				Optional:            true,
				MarkdownDescription: "whether this node is SDC or not,The acceptable value is `Yes` or `No`.",
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
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
