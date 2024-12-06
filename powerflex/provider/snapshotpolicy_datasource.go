/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &snapshotPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &snapshotPolicyDataSource{}
)

// SnapshotPolicyDataSource returns the snapshot policy data source
func SnapshotPolicyDataSource() datasource.DataSource {
	return &snapshotPolicyDataSource{}
}

type snapshotPolicyDataSource struct {
	client *goscaleio.Client
}

func (d *snapshotPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_policy"
}

func (d *snapshotPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotPolicyDataSourceSchema
}

func (d *snapshotPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

func (d *snapshotPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started snapshot policy data source read method")
	var (
		state models.SnapshotPolicyDataSourceModel
		sps   []scaleiotypes.SnapshotPolicy
		err   error
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	sps, err = helper.GetAllSnapshotPolicies(d.client)
	//check if there is any error while getting the snapshot policy
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Snapshot Policy",
			err.Error(),
		)
		return
	}
	// Gets all snapshot policies if no filter is provided
	if state.SnapshotPolicyFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.SnapshotPolicyFilter, sps)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering snapshot policies: %v please validate the filter", state.SnapshotPolicyFilter), err.Error(),
			)
			return
		}
		//convert filtered to []scaleiotypes.SnapshotPolicy
		filteredSps := []scaleiotypes.SnapshotPolicy{}
		for _, val := range filtered {
			filteredSps = append(filteredSps, val.(scaleiotypes.SnapshotPolicy))
		}
		sps = filteredSps
	}

	state.SnapshotPolicies = helper.UpdateSnapshotPolicyState(sps)
	state.ID = types.StringValue("dummyID")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
