package sdcprovider

import (
	"context"

	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
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

func (d *sdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

func (d *sdcDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return SDCDataSourceScheme, nil
}

func (d *sdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *sdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sdcDataSourceModel
	diags := req.Config.Get(ctx, &state)
	tflog.Info(ctx, "[POWERFLEX] sdcDataSourceModel"+helper.PrettyJSON((state)))
	system, err := d.client.FindSystem(state.SystemID.ValueString(), "", "")
	// tflog.Debug(ctx, "system"+helper.PrettyJSON(system))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems sdcs",
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

	// Set state
	state.Sdcs = getAllSdcState(sdcs)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getAllSdcState(sdcs []scaleiotypes.Sdc) (response []sdcModel) {
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

		response = append(response, sdcState)
	}

	return
}
