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

// SDCDataSource - function used to return SDC DataSource provider with singleton values.
func SDCDataSource() datasource.DataSource {
	return &sdcDataSource{}
}

// sdcFilterType - Enum structure for filter types.
var sdcFilterType = struct {
	ALL     string
	BY_NAME string
	BY_ID   string
}{
	ALL:     "ALL",
	BY_NAME: "BY_NAME",
	BY_ID:   "BY_ID",
}

// sdcDataSource - for returning singleton holder with goscaleio client.
type sdcDataSource struct {
	client *goscaleio.Client
}

// sdcDataSourceModel - for returning result to terraform.
type sdcDataSourceModel struct {
	Sdcs     []sdcModel   `tfsdk:"sdcs"`
	ID       types.String `tfsdk:"sdcid"`
	SystemID types.String `tfsdk:"systemid"`
	Name     types.String `tfsdk:"name"`
}

// sdcModel - MODEL for SDC data returned by goscaleio.
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

// sdcLinkModel - MODEL for SDC Links data returned by goscaleio.
type sdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// Metadata - function used to define datasource metadata[referance in tf file].
func (d *sdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

// GetSchema - function used to return SDC datasource schema.
func (d *sdcDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return SDCDataSourceScheme, nil
}

// Configure - function to call initial configurations before resource execution.
func (d *sdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

// Read - function to read sdc values from goscaleio.
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
	searchFilter := sdcFilterType.ALL
	if state.Name.ValueString() != "" {
		searchFilter = sdcFilterType.BY_NAME
	}
	if state.ID.ValueString() != "" {
		searchFilter = sdcFilterType.BY_ID
	}
	if state.Name.ValueString() != "" && state.ID.ValueString() != "" {
		searchFilter = sdcFilterType.ALL
	}

	if searchFilter == sdcFilterType.ALL {
		state.Sdcs = getAllSdcState(sdcs)
	} else {
		state.Sdcs = getFilteredSdcState(ctx, sdcs, searchFilter, state.Name.ValueString(), state.ID.ValueString())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// getFilteredSdcState - function to filter sdc result from goscaleio.
func getFilteredSdcState(ctx context.Context, sdcs []scaleiotypes.Sdc, method string, name string, id string) (response []sdcModel) {
	tflog.Debug(ctx, "[POWERFLEX] searchFilter getFilteredSdcState method "+method+" name "+name+" id "+id)
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
		// tflog.Debug(ctx, "[POWERFLEX] searchFilter getFilteredSdcState sdcValue.Name "+sdcValue.Name)
		// tflog.Debug(ctx, "[POWERFLEX] searchFilter getFilteredSdcState sdcValue.ID "+sdcValue.ID+" -- need "+id)
		if method == sdcFilterType.BY_NAME && name == sdcValue.Name {
			response = append(response, sdcState)
		}
		if method == sdcFilterType.BY_ID && id == sdcValue.ID {
			response = append(response, sdcState)
		}

	}

	return
}

// getAllSdcState - function to return all sdc result from goscaleio.
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
