package sdcprovider

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sdcDataSource{}
	_ datasource.DataSourceWithConfigure = &sdcDataSource{}
)

func SDCDataSource() datasource.DataSource {
	return &sdcDataSource{}
}

type sdcDataSource struct {
	client *goscaleio.Client
}

type sdcDataSourceModel struct {
	Sdcs     []sdcModel   `tfsdk:"sdcs"`
	ID       types.String `tfsdk:"sdcid"`
	SystemID types.String `tfsdk:"systemid"`
}

type sdcModel struct {
	ID                 types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"systemid"`
	SdcIP              types.String   `tfsdk:"sdcip"`
	SdcApproved        types.Bool     `tfsdk:"sdcapproved"`
	OnVMWare           types.Bool     `tfsdk:"onvmware"`
	SdcGUID            types.String   `tfsdk:"sdcguid"`
	MdmConnectionState types.String   `tfsdk:"mdmconnectionstate"`
	Name               types.String   `tfsdk:"name"`
	Links              []sdcLinkModel `tfsdk:"links"`
}

type sdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

func PrettyJSON(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}

func (d *sdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

func (d *sdcDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the list of coffees.",
		Attributes: map[string]tfsdk.Attribute{
			"sdcid": {
				Type:        types.StringType,
				Description: "Enter ID of Powerflex SDC. [Default/empty will all sdc present in given system]",
				Required:    true,
				Sensitive:   true,
			},
			"systemid": {
				Type:        types.StringType,
				Description: "Enter System ID of Powerflex System. [Default/empty will be any first system in list]",
				Required:    true,
				Sensitive:   true,
			},
			"sdcs": {
				Description: "List of coffees.",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description: "Numeric identifier of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"name": {
						Description: "Product name of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"sdcguid": {
						Description: "Fun tagline for the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"onvmware": {
						Description: "Product description of the coffee.",
						Type:        types.BoolType,
						Computed:    true,
					},
					"sdcapproved": {
						Description: "Product description of the coffee.",
						Type:        types.BoolType,
						Computed:    true,
					},
					"systemid": {
						Description: "Suggested cost of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"sdcip": {
						Description: "URI for an image of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"mdmconnectionstate": {
						Description: "URI for an image of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"links": {
						Description: "List of ingredients in the coffee.",
						Computed:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"rel": {
								Description: "Numeric identifier of the coffee ingredient.",
								Type:        types.StringType,
								Computed:    true,
							},
							"href": {
								Description: "Numeric identifier of the coffee ingredient.",
								Type:        types.StringType,
								Computed:    true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *sdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *sdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "(d *sdcDataSource) Read")
	var state sdcDataSourceModel
	diags := req.Config.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	tflog.Info(ctx, "[POWERFLEX] sdcDataSourceModel"+PrettyJSON((state)))
	system, err := d.client.FindSystem(state.SystemID.ValueString(), "", "")
	// tflog.Debug(ctx, "system"+PrettyJSON(system))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems Coffees",
			err.Error(),
		)
		return
	}
	sdcs, err := system.GetSdc()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex sdcs",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, sdcValue := range sdcs {
		sdcState := sdcModel{
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
			sdcState.Links = append(sdcState.Links, sdcLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}

		state.Sdcs = append(state.Sdcs, sdcState)
	}

	// Set state

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
