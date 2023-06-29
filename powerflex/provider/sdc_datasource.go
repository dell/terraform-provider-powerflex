/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sdcDataSource{}
	_ datasource.DataSourceWithConfigure = &sdcDataSource{}
)

// sdcDataSource - for returning singleton holder with goscaleio client.
type sdcDataSource struct {
	client *goscaleio.Client
}

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
	var state models.SdcDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "[POWERFLEX] sdcDataSourceModel"+helper.PrettyJSON((state)))

	system, err := helper.GetFirstSystem(d.client)

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
	searchFilter := helper.SdcFilterType.All
	if !state.Name.IsNull() {
		searchFilter = helper.SdcFilterType.ByName
	}
	if !state.ID.IsNull() {
		searchFilter = helper.SdcFilterType.ByID
	}

	allSdcWithStats := helper.GetAllSdcState(ctx, *d.client, sdcs)

	if len(*allSdcWithStats) == 0 {
		resp.Diagnostics.AddError("SDCs are not installed on the PowerFlex cluster.",
			"SDCs are not installed on the PowerFlex cluster.",
		)
		return
	}

	if searchFilter == helper.SdcFilterType.All {
		state.Sdcs = *allSdcWithStats
	} else {
		filterResult := helper.GetFilteredSdcState(allSdcWithStats, searchFilter, state.Name.ValueString(), state.ID.ValueString())
		if len(*filterResult) == 0 {
			resp.Diagnostics.AddError("Couldn't find SDC.",
				"Couldn't find SDC.",
			)
			return
		}
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
