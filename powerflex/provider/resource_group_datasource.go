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
	//Fetch all the details
	serviceDetails, err := d.gatewayClient.GetAllServiceDetails()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read service details", err.Error())
		return
	}
	// If resource group filter is provided, filter the service details
	if state.ResourceGroupFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.ResourceGroupFilter, serviceDetails)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering resource groups: %v please validate the filter", state.ResourceGroupFilter), err.Error(),
			)
			return
		}
		//convert filtered to []scaleiotypes.ServiceResponse
		filteredServiceDetails := []scaleiotypes.ServiceResponse{}
		for _, val := range filtered {
			filteredServiceDetails = append(filteredServiceDetails, val.(scaleiotypes.ServiceResponse))
		}
		serviceDetails = filteredServiceDetails
	}

	for _, service := range serviceDetails {
		serviceModel = append(serviceModel, helper.GetDataSourceResourceGroupState(service))
	}
	state.ResourceGroupDetails = serviceModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
