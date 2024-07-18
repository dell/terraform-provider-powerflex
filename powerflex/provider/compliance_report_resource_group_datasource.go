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
	"strconv"

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

	if !state.IPAddress.IsNull() {
		// Filter based on IP address
		ipAddresses := make([]string, 0)
		diags.Append(state.IPAddress.ElementsAs(ctx, &ipAddresses, true)...)

		for _, ip := range ipAddresses {
			complianceReport, err := helper.GetFilteredComplianceReport(complianceReportList, "IpAddress", ip)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting compliance report for resource group for given filter.",
					err.Error(),
				)
				return
			}
			if complianceReport != nil {
				complianceReports = append(complianceReports, *complianceReport)
			}
		}
	} else if !state.ServiceTags.IsNull() {
		// Filter based on service tags
		serviceTags := make([]string, 0)
		diags.Append(state.ServiceTags.ElementsAs(ctx, &serviceTags, true)...)
		for _, serviceTag := range serviceTags {
			complianceReport, err := helper.GetFilteredComplianceReport(complianceReportList, "ServiceTag", serviceTag)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting compliance report for resource group for given filter.",
					err.Error(),
				)
				return
			}
			if complianceReport != nil {
				complianceReports = append(complianceReports, *complianceReport)
			}
		}
	} else if !state.HostNames.IsNull() {
		// Filter based on host names
		hostNames := make([]string, 0)
		diags.Append(state.HostNames.ElementsAs(ctx, &hostNames, true)...)
		for _, hostName := range hostNames {
			complianceReport, err := helper.GetFilteredComplianceReport(complianceReportList, "HostName", hostName)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting compliance report for resource group for given filter.",
					err.Error(),
				)
				return
			}
			if complianceReport != nil {
				complianceReports = append(complianceReports, *complianceReport)
			}
		}
	} else if !state.ResourceIDs.IsNull() {
		// Filter based on ids
		resourceIDs := make([]string, 0)
		diags.Append(state.ResourceIDs.ElementsAs(ctx, &resourceIDs, true)...)
		for _, resourceID := range resourceIDs {
			complianceReport, err := helper.GetFilteredComplianceReport(complianceReportList, "ID", resourceID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting compliance report for resource group for given filter.",
					err.Error(),
				)
				return
			}
			if complianceReport != nil {
				complianceReports = append(complianceReports, *complianceReport)
			}
		}
	} else if !state.Compliant.IsNull() {
		// Filter based on Compliant status
		compliant := state.Compliant.ValueBool()
		complianceReport, err := helper.GetFilteredComplianceReport(complianceReportList, "Compliant", strconv.FormatBool(compliant))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting compliance report for resource group for given filter.",
				err.Error(),
			)
			return
		}
		if complianceReport != nil {
			complianceReports = append(complianceReports, *complianceReport)
		}
	} else {
		// add all reports
		complianceReports = append(complianceReports, complianceReportList...)
	}

	state.ComplianceReports = helper.GetComplianceReportsModel(complianceReports)
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
