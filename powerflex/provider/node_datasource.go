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
	client *goscaleio.Client
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

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
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

	// Fetch Node details if IDs are provided
	if !state.NodeIDs.IsNull() {
		nodeIDs := make([]string, 0)
		diags.Append(state.NodeIDs.ElementsAs(ctx, &nodeIDs, true)...)

		for _, nodeID := range nodeIDs {
			nodeDetails, err := d.client.GetNodeByID(nodeID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting node details using id %v", nodeID), err.Error(),
				)
				return
			}
			nodeModel = append(nodeModel, helper.GetNodeState(*nodeDetails))
		}
	} else if !state.IPAddresses.IsNull() {
		// Fetch Node details if IPs are provided
		IPAddresses := make([]string, 0)
		diags.Append(state.IPAddresses.ElementsAs(ctx, &IPAddresses, true)...)

		for _, ipAddress := range IPAddresses {
			nodeDetails, err := d.client.GetNodeByFilters("ipAddress", ipAddress)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting node details using ip %v", ipAddress), err.Error(),
				)
				return
			}
			nodeModel = append(nodeModel, helper.GetNodeState(nodeDetails[0]))
		}
	} else if !state.ServiceTags.IsNull() {
		// Fetch Node details if service tags are provided
		serviceTags := make([]string, 0)
		diags.Append(state.ServiceTags.ElementsAs(ctx, &serviceTags, true)...)

		for _, serviceTag := range serviceTags {
			nodeDetails, err := d.client.GetNodeByFilters("serviceTag", serviceTag)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting node details using service tag %v", serviceTag), err.Error(),
				)
				return
			}
			nodeModel = append(nodeModel, helper.GetNodeState(nodeDetails[0]))
		}
	} else if !state.NodePoolIDs.IsNull() {
		// Fetch Node details if node pool IDs are provided
		nodePoolIDs := make([]int64, 0)
		diags.Append(state.NodePoolIDs.ElementsAs(ctx, &nodePoolIDs, true)...)

		for _, nodePoolID := range nodePoolIDs {
			nodePoolDetails, err := d.client.GetNodePoolByID(int(nodePoolID))
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting node pool details using id %v", nodePoolID), err.Error(),
				)
				return
			}

			for _, node := range nodePoolDetails.ManagedDeviceList.ManagedDevices {
				nodeModel = append(nodeModel, helper.GetNodeState(node))
			}
		}
	} else if !state.NodePoolNames.IsNull() {
		// Fetch Node details if node pool names are provided
		nodePoolNames := make([]string, 0)
		diags.Append(state.NodePoolNames.ElementsAs(ctx, &nodePoolNames, true)...)

		for _, nodePoolName := range nodePoolNames {
			if nodePoolName == "Global" {
				nodeDetails, err := d.client.GetAllNodes()
				if err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Error in getting node details"), err.Error(),
					)
					return
				}

				for _, node := range nodeDetails {
					if node.DeviceGroupList.DeviceGroup[0].GroupName == "Global" {
						nodeModel = append(nodeModel, helper.GetNodeState(node))
					}
				}
			} else {
				nodePoolDetails, err := d.client.GetNodePoolByName(nodePoolName)
				if err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Error in getting node pool details using name %v", nodePoolName), err.Error(),
					)
					return
				}

				for _, node := range nodePoolDetails.ManagedDeviceList.ManagedDevices {
					nodeModel = append(nodeModel, helper.GetNodeState(node))
				}
			}
		}
	} else {
		nodeDetails, err := d.client.GetAllNodes()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in getting node details"), err.Error(),
			)
			return
		}

		for _, node := range nodeDetails {
			nodeModel = append(nodeModel, helper.GetNodeState(node))
		}
	}

	state.NodeDetails = nodeModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
