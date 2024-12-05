/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"fmt"
	"strings"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
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
func (d *sdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
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
	tflog.Info(ctx, "system Version"+system.System.SystemVersionName)

	// If it is < 4_0 there is no NVME/HostType value so we should skip this filter
	if !strings.Contains(system.System.SystemVersionName, "3_6") {
		// SDC endpoints returns both SDC and NVMe host. Need to filter out NVMe host
		n := 0
		for _, sdc := range sdcs {
			if sdc.HostType == "SdcHost" {
				sdcs[n] = sdc
				n++
			}
		}
		sdcs = sdcs[:n]
	}
	// Set state for filters
	if state.SdcFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.SdcFilter, sdcs)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering sdcs: %v please validate the filter", state.SdcFilter), err.Error(),
			)
			return
		}
		filteredSdc := []scaleiotypes.Sdc{}
		for _, val := range filtered {
			filteredSdc = append(filteredSdc, val.(scaleiotypes.Sdc))
		}
		sdcs = filteredSdc
	}
	state.Sdcs = helper.GetAllSdcState(ctx, *d.client, sdcs)

	state.ID = types.StringValue("sdc_datasource_id")

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
