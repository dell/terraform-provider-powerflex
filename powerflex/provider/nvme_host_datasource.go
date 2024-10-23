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
	_ datasource.DataSource              = &nvmeHostDataSource{}
	_ datasource.DataSourceWithConfigure = &nvmeHostDataSource{}
)

// nvmeHostDataSource - for returning singleton holder with goscaleio client.
type nvmeHostDataSource struct {
	client *goscaleio.Client
}

// NvmeHostDataSource - function used to return NvmeHost DataSource provider with singleton values.
func NvmeHostDataSource() datasource.DataSource {
	return &nvmeHostDataSource{}
}

// Metadata - function used to define datasource metadata[reference in tf file].
func (d *nvmeHostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvme_host"
}

// Schema - function used to return NvmeHost datasource schema.
func (d *nvmeHostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NvmeHostDataSourceSchema
}

// Configure - function to call initial configurations before resource execution.
func (d *nvmeHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *nvmeHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.NvmeHostDataSource
	var models []models.NvmeHostDatasourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "[POWERFLEX] NvmeHostDataSource"+helper.PrettyJSON(state))

	system, err := helper.GetFirstSystem(d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex specific system",
			err.Error(),
		)
		return
	}

	nvmeHosts, err := helper.GetAllNvmeHosts(system)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex NVMe Hosts",
			err.Error(),
		)
		return
	}

	if state.Filter == nil {
		for _, host := range nvmeHosts {
			models = append(models, helper.GetNvmeHostState(host))
		}
	} else {
		filteredHosts, err := helper.GetDataSourceByValue(*state.Filter, nvmeHosts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in getting NVMe Host details  %v", state.Filter), err.Error(),
			)
			return
		}
		for i := 0; i < len(filteredHosts); i++ {
			m := filteredHosts[i].(scaleiotypes.NvmeHost)
			models = append(models, helper.GetNvmeHostState(m))
		}
	}

	nvmeHostState := models
	state.ID = types.StringValue("nvme_host_datasource")
	state.Details = nvmeHostState
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
