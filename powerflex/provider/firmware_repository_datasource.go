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
	_ datasource.DataSource              = &firmwareRepositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &firmwareRepositoryDataSource{}
)

// FirmwareRepositoryDataSource returns the firmware repository data source
func FirmwareRepositoryDataSource() datasource.DataSource {
	return &firmwareRepositoryDataSource{}
}

type firmwareRepositoryDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *firmwareRepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firmware_repository"
}

func (d *firmwareRepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = FirmwareRepositoryDataSourceSchema
}

func (d *firmwareRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *firmwareRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started firmware repository data source read method")
	var (
		state          models.FirmwareRepositoryDatasourceModel
		 firmwareRepositoryModel []models.FirmwareRepositoryDetails
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch firmware repository details if IDs are provided
	if !state.FirmwareRepositoryIDs.IsNull() {
		frIDs := make([]string, 0)
		diags.Append(state.FirmwareRepositoryIDs.ElementsAs(ctx, &frIDs, true)...)

		for _, frID := range frIDs {
			fr, err := d.gatewayClient.GetUploadComplianceDetailsUsingID(frID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting firmware repository details using id %v", frID), err.Error(),
				)
				return
			}
			firmwareRepositoryModel = append(firmwareRepositoryModel, helper.GetAllFirmwareRepositoryState(fr))
		}
	} else if !state.FirmwareRepositoryNames.IsNull() {
		frNames := make([]string, 0)
		diags.Append(state.FirmwareRepositoryNames.ElementsAs(ctx, &frNames, true)...)

		for _, frName := range frNames {
			fr, err := d.gatewayClient.GetFirmwareRepositoryDetailsUsingName(frName)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting firmware repository details using name %v", frName), err.Error(),
				)
				return
			}
			firmwareRepositoryModel = append(firmwareRepositoryModel, helper.GetAllFirmwareRepositoryState(fr))
		}
	} else {
		allfr, err := d.gatewayClient.GetAllUploadComplianceDetails()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting firmware repository details", err.Error(),
			)
			return
		}

		for _, frs := range *allfr {
			fr, err := d.gatewayClient.GetUploadComplianceDetailsUsingID(frs.ID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting firmware repository details using id %v", frs.ID), err.Error(),
				)
				return
			}
			firmwareRepositoryModel = append(firmwareRepositoryModel, helper.GetAllFirmwareRepositoryState(fr))
		}
	}

	state.FirmwareRepositoryDetails = firmwareRepositoryModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}
