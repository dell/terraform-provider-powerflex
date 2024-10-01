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

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
)

var (
	_ datasource.DataSource              = &replicationPairsDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationPairsDataSource{}
)

// ReplicationPairsDataSource returns the ReplicationPairs data source
func ReplicationPairsDataSource() datasource.DataSource {
	return &replicationPairsDataSource{}
}

type replicationPairsDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *replicationPairsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_pair"
}

func (d *replicationPairsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationPairsDataSourceSchema
}

func (d *replicationPairsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *replicationPairsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ReplicationPairDataSourceModel
	// Get the state incase filters are set
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get All Replication Pairs
	rps, err := helper.GetReplicationPairs(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting Replication Pairs details",
			err.Error(),
		)
		return
	}
	// Set state for filters
	if state.ReplicationPairFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.ReplicationPairFilter, rps)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering Replication Pairs: %v please validate the filter", state.ReplicationPairFilter), err.Error(),
			)
			return
		}
		filteredPair := []scaleiotypes.ReplicationPair{}
		for _, val := range filtered {
			filteredPair = append(filteredPair, val.(scaleiotypes.ReplicationPair))
		}
		rps = filteredPair
	}
	mappedRps := helper.MapReplicationPairsState(rps, state)
	diagsState := resp.State.Set(ctx, mappedRps)
	resp.Diagnostics.Append(diagsState...)
}
