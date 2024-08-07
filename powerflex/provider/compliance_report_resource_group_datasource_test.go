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
	"fmt"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ComplianceGetAllConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
}
`

var ComplianceFilterIPAddressConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
		ip_addresses = ["` + ComplianceReportDataPoints.IPAddress + `"]
	 }
}
`

var ComplianceFilterCompliantConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
		compliant = true
	 }
}
`
var ComplianceFilterHostNamesConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
		 host_names = ["` + ComplianceReportDataPoints.HostName + `"]
	 }
}
`
var ComplianceFilterServiceTagsConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
	 	service_tags = ["` + ComplianceReportDataPoints.ServiceTag + `"]
	 }
}
`

var ComplianceFilterResourceIdsConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
		resource_ids = ["` + ComplianceReportDataPoints.ResourceID + `"]
	 }
}
`
var ComplianceGetAllError = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
}
`
var ComplianceFilterMultiple = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	 resource_group_id = "` + ComplianceReportDataPoints.ResourceGroupID + `"
	 filter {
		ip_addresses = ["` + ComplianceReportDataPoints.IPAddress + `"]
		compliant = true
	 }
}
`

func TestAccComplianceReportDataSource(t *testing.T) {
	t.Skip("Skipping this test case, only use on 4.x or greater")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ComplianceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.#"),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterIPAddressConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.0.ip_address", ComplianceReportDataPoints.IPAddress),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterHostNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.0.host_name", ComplianceReportDataPoints.HostName),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterServiceTagsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.0.service_tag", ComplianceReportDataPoints.ServiceTag),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterResourceIdsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.0.id", ComplianceReportDataPoints.ResourceID),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterCompliantConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.#"),
				),
			},
			{
				Config: ProviderConfigForTesting + ComplianceFilterMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_compliance_report_resource_group.complianceReport", "compliance_reports.#"),
				),
			},
		},
	})
}

func TestAccComplianceReportDataSourceNegative(t *testing.T) {
	t.Skip("Skipping this test case, only use on 4.x or greater")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + ComplianceGetAllError,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFilteredComplianceReports).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ComplianceFilterCompliantConfig,
				ExpectError: regexp.MustCompile(`.*Error in getting compliance report for resource group for given filter.*`),
			},
		},
	})
}
