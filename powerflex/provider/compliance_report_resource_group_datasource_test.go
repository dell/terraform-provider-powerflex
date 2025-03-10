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
}
`

var ComplianceFilterIPAddressConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	filter {
		ip_address = ["` + ComplianceReportDataPoints.IPAddress + `"]
	}
}
`
var ComplianceFilterCompliantConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	filter {
		compliant = true
	}
}
`
var ComplianceFilterHostNamesConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	filter {
		host_name = ["` + ComplianceReportDataPoints.HostName + `"]
	}
}
`
var ComplianceFilterServiceTagsConfig = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	filter {
	    service_tag = ["` + ComplianceReportDataPoints.ServiceTag + `"]
	}
}
`
var ComplianceGetAllError = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
}
`
var ComplianceFilterMultiple = `
data "powerflex_compliance_report_resource_group" "complianceReport" {
	filter {
		ip_address = ["` + ComplianceReportDataPoints.IPAddress + `"]
		compliant = true
		host_name = ["` + ComplianceReportDataPoints.HostName + `"]
	}
}
`

// AT
func TestAccDatasourceAcceptanceComplianceReport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ComplianceGetAllConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// UT
func TestAccComplianceReportDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + ComplianceFilterCompliantConfig,
				ExpectError: regexp.MustCompile(`.*Error in getting compliance report for resource group for given filter.*`),
			},
		},
	})
}
