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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &nodeDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeDataSource{}
)

// NodeDataSource returns the node data source
func NodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}

type nodeDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NodeDataSourceSchema
}

func (d *nodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *nodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started node data source read method")

	var (
		state     models.NodeDataSourceModel
		nodeModel []models.NodeModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeDetails, err := d.gatewayClient.GetAllNodes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read node details", err.Error(),
		)
		return
	}

	// If filter is present
	if state.NodeFilter != nil {
		// Get filtered nodes
		nodesFiltered, err := helper.GetDataSourceByValue(*state.NodeFilter, nodeDetails)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering node: %v please validate the filter", state.NodeFilter),
				err.Error(),
			)
			return
		}
		// Convert filtered nodes to node details
		nodeDetailFiltered := []scaleiotypes.NodeDetails{}
		for _, val := range nodesFiltered {
			nodeDetailFiltered = append(nodeDetailFiltered, val.(scaleiotypes.NodeDetails))
		}
		nodeDetails = nodeDetailFiltered
	}

	for _, node := range nodeDetails {
		nodeModel = append(nodeModel, helper.GetNodeState(node))
	}

	state.NodeDetails = nodeModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
