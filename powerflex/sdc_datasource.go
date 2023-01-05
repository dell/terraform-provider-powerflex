package powerflex

import (
	"context"

	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

// Metadata - function used to define datasource metadata[referance in tf file].
func (d *sdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

// GetSchema - function used to return SDC datasource schema.
func (d *sdcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SDCDataSourceScheme
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

	system, err := getFirstSystem(d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex specific system",
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
	searchFilter := sdcFilterType.All
	if state.Name.ValueString() != "" {
		searchFilter = sdcFilterType.ByName
	}
	if state.ID.ValueString() != "" {
		searchFilter = sdcFilterType.ByID
	}

	allSdcWithStats, _ := getAllSdcState(ctx, *d.client, sdcs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Statics for sdc id = "+state.ID.ValueString()+", name = "+state.Name.ValueString(),
			err.Error(),
		)
		return
	}

	if searchFilter == sdcFilterType.All {
		state.Sdcs = *allSdcWithStats
	} else {
		filterResult := getFilteredSdcState(allSdcWithStats, searchFilter, state.Name.ValueString(), state.ID.ValueString())
		state.Sdcs = *filterResult
	}

	state.Name = types.StringValue(state.Name.ValueString())
	state.ID = types.StringValue(state.ID.ValueString())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}