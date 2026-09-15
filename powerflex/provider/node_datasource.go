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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	client          *goscaleio.Client
	gatewayClient   *goscaleio.GatewayClient
	gatewayEndpoint string
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
		d.gatewayEndpoint = req.ProviderData.(*powerflexProvider).gatewayEndpoint
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// getAllNodesWithBearerAuth fetches all nodes using Bearer token authentication.
// This works around a bug in goscaleio SDK where GetAllNodes uses Basic auth
// for gateway versions other than "4.0", which fails on PFMP 5.1+ that requires Bearer tokens.
func (d *nodeDataSource) getAllNodesWithBearerAuth(ctx context.Context) ([]scaleiotypes.NodeDetails, error) {
	token, err := d.gatewayClient.NewTokenGeneration()
	if err != nil {
		return nil, fmt.Errorf("error generating token for node API: %s", err)
	}

	path := "/Api/V1/ManagedDevice"
	req, err := http.NewRequest(http.MethodGet, d.gatewayEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// #nosec G402
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node API returned status %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading node API response: %s", err)
	}

	var nodes []scaleiotypes.NodeDetails
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, fmt.Errorf("error parsing node API response: %s", err)
	}

	return nodes, nil
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

	// If SDK returns empty results or error and we have the endpoint, retry with Bearer auth.
	// This works around a goscaleio SDK bug where GetAllNodes uses Basic auth
	// for gateway versions != "4.0" (e.g., PFMP 5.1 reports version "5.1"),
	// but PFMP 5.1 requires Bearer token authentication.
	if (err != nil || len(nodeDetails) == 0) && d.gatewayEndpoint != "" {
		tflog.Info(ctx, "SDK GetAllNodes returned error or empty results, retrying with Bearer token auth")
		nodeDetails, err = d.getAllNodesWithBearerAuth(ctx)
	}

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
