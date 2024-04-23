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
	_ datasource.DataSource              = &resourceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceGroupDataSource{}
)

// ResourceGroupDataSource returns the ResourceGroup data source
func ResourceGroupDataSource() datasource.DataSource {
	return &resourceGroupDataSource{}
}

type resourceGroupDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *resourceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

func (d *resourceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ResourceGroupDataSourceSchema
}

func (d *resourceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *resourceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started service data source read method")

	var (
		state        models.ResourceGroupDataSourceModel
		serviceModel []models.ResourceGroupModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Service details if IDs are provided
	if !state.ResourceGroupIDs.IsNull() {
		serviceIDs := make([]string, 0)
		diags.Append(state.ResourceGroupIDs.ElementsAs(ctx, &serviceIDs, true)...)

		for _, serviceID := range serviceIDs {
			serviceDetails, err := d.gatewayClient.GetServiceDetailsByID(serviceID, false)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting service details using id %v", serviceID), err.Error(),
				)
				return
			}
			serviceModel = append(serviceModel, helper.GetDataSourceResourceGroupState(*serviceDetails))
		}
	} else if !state.ResourceGroupNames.IsNull() {
		// Fetch Service details if names are provided
		Names := make([]string, 0)
		diags.Append(state.ResourceGroupNames.ElementsAs(ctx, &Names, true)...)

		for _, name := range Names {
			serviceDetails, err := d.gatewayClient.GetServiceDetailsByFilter("name", name)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting service details using Name %v", name), err.Error(),
				)
				return
			}
			serviceModel = append(serviceModel, helper.GetDataSourceResourceGroupState(serviceDetails[0]))
		}
	} else {
		//Fetch all the details
		serviceDetails, err := d.gatewayClient.GetAllServiceDetails()
		if err != nil {
			resp.Diagnostics.AddError("Error in getting service details", err.Error())
			return
		}

		for _, service := range serviceDetails {
			serviceModel = append(serviceModel, helper.GetDataSourceResourceGroupState(service))
		}
	}

	state.ResourceGroupDetails = serviceModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
