/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &nvmeTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &nvmeTargetDataSource{}
)

// nvmeTargetDataSource - for returning singleton holder with goscaleio client.
type nvmeTargetDataSource struct {
	client *goscaleio.Client
}

// NvmeTargetDataSource - function used to return NvmeTarget DataSource provider with singleton values.
func NvmeTargetDataSource() datasource.DataSource {
	return &nvmeTargetDataSource{}
}

// Metadata - function used to define datasource metadata[reference in tf file].
func (d *nvmeTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvme_target"
}

// Schema - function used to return NvmeTarget datasource schema.
func (d *nvmeTargetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NvmeTargetDataSourceSchema
}

// Configure - function to call initial configurations before resource execution.
func (d *nvmeTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *nvmeTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.NvmeTargetDataSource
	var filteredTargets []models.NvmeTargetDatasourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "[POWERFLEX] NvmeTargetDataSource"+helper.PrettyJSON(state))

	system, err := helper.GetFirstSystem(d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex specific system",
			err.Error(),
		)
		return
	}

	nvmeTargets, err := system.GetAllSdts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex NVMe Targets",
			err.Error(),
		)
		return
	}

	if state.Filter != nil {
		filteredHosts, err := helper.GetDataSourceByValue(*state.Filter, nvmeTargets)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting NVMe target details  %v", state.Filter), err.Error(),
			)
			return
		}
		nvmeTargets, err = helper.ConvertSlice[scaleiotypes.Sdt](filteredHosts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting NVMe target details  %v", state.Filter), err.Error(),
			)
			return
		}
	}

	for _, nvmeTarget := range nvmeTargets {
		nt := models.NvmeTargetDatasourceModel{}
		err = helper.CopyFields(ctx, nvmeTarget, &nt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting NVMe target details", err.Error(),
			)
			return
		}
		filteredTargets = append(filteredTargets, nt)
	}

	err = helper.SetAttachedNvmeHostInfo(system, filteredTargets)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting NVMe controller", err.Error(),
		)
	}

	state.ID = types.StringValue("nvme_target_datasource")
	state.Details = filteredTargets
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
