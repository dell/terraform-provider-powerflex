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

	"terraform-provider-powerflex/powerflex/models"

	"terraform-provider-powerflex/powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &complianceReportResourceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &complianceReportResourceGroupDataSource{}
)

// ComplianceReportResourceGroupDataSource returns the volume data source
func ComplianceReportResourceGroupDataSource() datasource.DataSource {
	return &complianceReportResourceGroupDataSource{}
}

type complianceReportResourceGroupDataSource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (d *complianceReportResourceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_report_resource_group"
}

func (d *complianceReportResourceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ComplianceReportResourceGroupSchema
}

func (d *complianceReportResourceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *complianceReportResourceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started compliance report data source read method")

	var state models.ComplianceReportResourceGroupDatasource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all Compliance Reports
	complianceReportList, err := d.gatewayClient.GetServiceComplianceDetails(state.ResourceGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting compliance report for resource group.",
			err.Error(),
		)
		return
	}
	complianceReports := make([]scaleiotypes.ComplianceReport, 0)
	if state.ComplianceReportFilter == nil {
		// add all reports
		complianceReports = append(complianceReports, complianceReportList...)
	} else {
		// Get filtered reports
		complianceReports, err = helper.GetFilteredComplianceReports(complianceReportList, *state.ComplianceReportFilter)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting compliance report for resource group for given filter.",
				err.Error(),
			)
			return
		}
	}

	state.ComplianceReports = helper.GetComplianceReportsModel(complianceReports)
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
