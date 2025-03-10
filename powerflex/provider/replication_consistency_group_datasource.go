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

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &replicationConsistencyGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationConsistencyGroupDataSource{}
)

// ReplicationConsistencyGroupDataSource returns the ReplicationConsistencyGroup data source
func ReplicationConsistencyGroupDataSource() datasource.DataSource {
	return &replicationConsistencyGroupDataSource{}
}

type replicationConsistencyGroupDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *replicationConsistencyGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_consistency_group"
}

func (d *replicationConsistencyGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationConsistencyGroupDataSourceSchema
}

func (d *replicationConsistencyGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {

		d.client = req.ProviderData.(*powerflexProvider).client
	}

	if req.ProviderData.(*powerflexProvider).gatewayClient != nil {

		d.gatewayClient = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

func (d *replicationConsistencyGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ReplicationConsistencyGroupDataSourceModel
	// Get the state incase filters are set
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get All RCGs
	rcgs, err := helper.GetReplicationConsistencyGroups(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting Replication Consistency Groups details",
			err.Error(),
		)
		return
	}
	// Set state for filters
	if state.ReplicationConsistencyGroupFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.ReplicationConsistencyGroupFilter, rcgs)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering Replication Consistency Groups: %v please validate the filter", state.ReplicationConsistencyGroupFilter), err.Error(),
			)
			return
		}
		filteredPair := []scaleiotypes.ReplicationConsistencyGroup{}
		for _, val := range filtered {
			filteredPair = append(filteredPair, val.(scaleiotypes.ReplicationConsistencyGroup))
		}
		rcgs = filteredPair
	}
	mappedRps := helper.MapReplicationConsistencyGroupsState(rcgs, state)
	diagsState := resp.State.Set(ctx, mappedRps)
	resp.Diagnostics.Append(diagsState...)
}
