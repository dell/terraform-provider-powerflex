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
	_ datasource.DataSource              = &osRepositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &osRepositoryDataSource{}
)

// OSRepositoryDataSource returns the OS repository data source
func OSRepositoryDataSource() datasource.DataSource {
	return &osRepositoryDataSource{}
}

type osRepositoryDataSource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *osRepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_os_repository"
}

func (d *osRepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OSRepositoryDatasourceSchema
}

func (d *osRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}
	d.system = system
}

// Read refreshes the Terraform state with the latest data.
func (d *osRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started OS repository data source read method")
	var (
		state               models.OSRepositoryDataSource
		osRepositoriesModel []models.OSRepositoryModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	osRepositories, err := helper.GetAllOsRepositories(d.system)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting OS repository details", err.Error(),
		)
		return
	}

	if state.OSRepoFilter == nil {
		for _, osRepo := range osRepositories {
			osRepositoriesModel = append(osRepositoriesModel, helper.GetAllOSRepositoryState(osRepo))
		}
	} else {

		osRepo, err := helper.GetDataSourceByValue(*state.OSRepoFilter, osRepositories)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in getting OS repository details  %v", state.OSRepoFilter), err.Error(),
			)
			return
		}
		for i := 0; i < len(osRepo); i++ {
			var osRepoCast scaleiotypes.OSRepository = osRepo[i].(scaleiotypes.OSRepository)

			osRepositoriesModel = append(osRepositoriesModel, helper.GetAllOSRepositoryState(osRepoCast))
		}

	}
	state.OSRepositories = osRepositoriesModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
