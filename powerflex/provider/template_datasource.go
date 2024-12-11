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
	_ datasource.DataSource              = &templateDataSource{}
	_ datasource.DataSourceWithConfigure = &templateDataSource{}
)

// TemplateDataSource returns the template data source
func TemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *templateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *templateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = TemplateDataSourceSchema
}

func (d *templateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started template data source read method")

	var (
		state         models.TemplateDataSourceModel
		templateModel []models.TemplateModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateDetails, err := d.gatewayClient.GetAllTemplates()
	if err != nil {
		resp.Diagnostics.AddError("Error in getting template details", err.Error())
		return
	}

	// Fetch Template details if IDs are provided
	if state.TemplateFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.TemplateFilter, templateDetails)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering Template: %v please validate the filter", state.TemplateDetails), err.Error(),
			)
			return
		}
		filteredTemplate := []scaleiotypes.TemplateDetails{}
		for _, val := range filtered {
			filteredTemplate = append(filteredTemplate, val.(scaleiotypes.TemplateDetails))
		}
		templateDetails = filteredTemplate
	}

	for _, template := range templateDetails {
		templateModel = append(templateModel, helper.GetTemplateState(template))
	}

	state.TemplateDetails = templateModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
