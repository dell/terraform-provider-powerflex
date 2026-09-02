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
	_ datasource.DataSource              = &templateDataSource{}
	_ datasource.DataSourceWithConfigure = &templateDataSource{}
)

// TemplateDataSource returns the template data source
func TemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	client          *goscaleio.Client
	gatewayClient   *goscaleio.GatewayClient
	gatewayEndpoint string
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
		d.gatewayEndpoint = req.ProviderData.(*powerflexProvider).gatewayEndpoint
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// getAllTemplatesWithBearerAuth fetches all templates using Bearer token authentication.
// This works around a bug in goscaleio SDK where GetAllTemplates uses Basic auth
// for gateway versions other than "4.0", which fails on PFMP 5.1+ that requires Bearer tokens.
func (d *templateDataSource) getAllTemplatesWithBearerAuth(ctx context.Context) ([]scaleiotypes.TemplateDetails, error) {
	token, err := d.gatewayClient.NewTokenGeneration()
	if err != nil {
		return nil, fmt.Errorf("error generating token for template API: %s", err)
	}

	path := "/Api/V1/template"
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
		return nil, fmt.Errorf("template API returned status %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading template API response: %s", err)
	}

	var templates scaleiotypes.TemplateDetailsFilter
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("error parsing template API response: %s", err)
	}

	return templates.TemplateDetails, nil
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

	// Try SDK method first
	templateDetails, err := d.gatewayClient.GetAllTemplates()

	// If SDK returns empty results and we have the endpoint, retry with Bearer auth.
	// This works around a goscaleio SDK bug where GetAllTemplates uses Basic auth
	// for gateway versions != "4.0" (e.g., PFMP 5.1 reports version "5.1"),
	// but PFMP 5.1 requires Bearer token authentication.
	if err == nil && len(templateDetails) == 0 && d.gatewayEndpoint != "" {
		tflog.Info(ctx, "SDK GetAllTemplates returned empty results, retrying with Bearer token auth")
		templateDetails, err = d.getAllTemplatesWithBearerAuth(ctx)
	}

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
