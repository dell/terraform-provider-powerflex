/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
)

var (
	_ datasource.DataSource              = &resourceCredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceCredentialDataSource{}
)

// ResourceCredentialDataSource returns the resourceCredentialGroup data source
func ResourceCredentialDataSource() datasource.DataSource {
	return &resourceCredentialDataSource{}
}

type resourceCredentialDataSource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *resourceCredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_credential"
}

func (d *resourceCredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ResourceCredentialDataSourceSchema
}

func (d *resourceCredentialDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {

		d.client = req.ProviderData.(*powerflexProvider).client
	}

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

func (d *resourceCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ResourceCredentialDataSourceModel
	// Get the state incase filters are set
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rcs, err := helper.GetResourceCredentials(ctx, d.system)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting resource credentials",
			err.Error(),
		)
		return
	}

	// Set state for filters
	if state.ResourceCredentialFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.ResourceCredentialFilter, rcs)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering resource credentials: %v please validate the filter", state.ResourceCredentialFilter), err.Error(),
			)
			return
		}
		filteredCreds := []scaleiotypes.CredObj{}
		for _, val := range filtered {
			filteredCreds = append(filteredCreds, val.(scaleiotypes.CredObj))
		}
		rcs = filteredCreds
	}

	mapped := helper.MapResourceCredentials(rcs, state)
	mapped.ID = types.StringValue("resource_credential_datasource_id")
	diagsState := resp.State.Set(ctx, mapped)
	resp.Diagnostics.Append(diagsState...)
	if resp.Diagnostics.HasError() {
		return
	}
}
