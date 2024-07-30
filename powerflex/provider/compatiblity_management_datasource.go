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

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
)

var (
	_ datasource.DataSource              = &compatibilityManagementDataSource{}
	_ datasource.DataSourceWithConfigure = &compatibilityManagementDataSource{}
)

// CompatibilityManagementDataSource returns the Compatibility Management datasource
func CompatibilityManagementDataSource() datasource.DataSource {
	return &compatibilityManagementDataSource{}
}

type compatibilityManagementDataSource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *compatibilityManagementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compatibility_management"
}

func (d *compatibilityManagementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CompatibilityManagementDataSourceSchema
}

func (d *compatibilityManagementDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client

	system, err := helper.GetFirstSystem(d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	d.system = system
}

func (d *compatibilityManagementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.CompatibilityManagementDatasourceModel

	// Get the Compatibility Management details
	cm, err := helper.GetCompatibilityManagement(ctx, d.system)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting compatibility management details",
			err.Error(),
		)
		return
	}

	// Set state
	state = helper.MapCompatibilityManagementState(ctx, cm)
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
