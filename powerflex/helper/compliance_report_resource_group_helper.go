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

package helper

import (
	"fmt"
	"strconv"
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetComplianceReportsModel converts []ComplianceReport to []ComplianceReportResourceGroup
func GetComplianceReportsModel(complianceReports []scaleiotypes.ComplianceReport) []models.ComplianceReportResourceGroup {
	complainceReportModel := make([]models.ComplianceReportResourceGroup, len(complianceReports))

	for i, report := range complianceReports {
		complainceReportModel[i] = models.ComplianceReportResourceGroup{
			ServiceTag:                 types.StringValue(report.ServiceTag),
			IPAddress:                  types.StringValue(report.IPAddress),
			FirmwareRepositoryName:     types.StringValue(report.FirmwareRepositoryName),
			ComplianceReportComponents: GetComplianceReportComponents(report.ComplianceReportComponents),
			Compliant:                  types.BoolValue(report.Compliant),
			DeviceType:                 types.StringValue(report.DeviceType),
			Model:                      types.StringValue(report.Model),
			Available:                  types.BoolValue(report.Available),
			ManagedState:               types.StringValue(report.ManagedState),
			EmbeddedReport:             types.BoolValue(report.EmbeddedReport),
			DeviceState:                types.StringValue(report.DeviceState),
			ID:                         types.StringValue(report.ID),
			HostName:                   types.StringValue(report.HostName),
			CanUpdate:                  types.BoolValue(report.CanUpdate),
		}
	}
	return complainceReportModel
}

// GetComplianceReportComponents converts list of scaleiotypes.Components to list of models.Components
func GetComplianceReportComponents(components []scaleiotypes.ComplianceReportComponents) []models.ComplianceReportComponents {
	componentModels := make([]models.ComplianceReportComponents, len(components))
	for i, component := range components {
		componentModels[i] = models.ComplianceReportComponents{
			ID:              types.StringValue(component.ID),
			Name:            types.StringValue(component.Name),
			CurrentVersion:  GetComplianceReportComponentVersionInfo(component.CurrentVersion),
			TargetVersion:   GetComplianceReportComponentVersionInfo(component.TargetVersion),
			Vendor:          types.StringValue(component.Vendor),
			OperatingSystem: types.StringValue(component.OperatingSystem),
			Compliant:       types.BoolValue(component.Compliant),
			Oscapable:       types.BoolValue(component.Oscapable),
			Software:        types.BoolValue(component.Software),
			RPM:             types.BoolValue(component.RPM),
		}
	}
	return componentModels
}

// GetComplianceReportComponentVersionInfo converts scaleiotypes.ComplianceReportComponentVersionInfo to models.ComplianceReportComponentVersionInfo
func GetComplianceReportComponentVersionInfo(versionInfo scaleiotypes.ComplianceReportComponentVersionInfo) models.ComplianceReportComponentVersionInfo {
	return models.ComplianceReportComponentVersionInfo{
		ID:                 types.StringValue(versionInfo.ID),
		FirmwareName:       types.StringValue(versionInfo.FirmwareName),
		FirmwareType:       types.StringValue(versionInfo.FirmwareType),
		FirmwareVersion:    types.StringValue(versionInfo.FirmwareVersion),
		FirmwareLastUpdate: types.StringValue(versionInfo.FirmwareLastUpdate),
		FirmwareLevel:      types.StringValue(versionInfo.FirmwareLevel),
	}
}

// GetFilteredComplianceReport returns the first compliance report that matches the filter and value
func GetFilteredComplianceReport(complianceReports []scaleiotypes.ComplianceReport, filter, value string) (*scaleiotypes.ComplianceReport, error) {
	for _, report := range complianceReports {
		switch filter {
		case "IpAddress":
			if report.IPAddress == value {
				return &report, nil
			}
		case "ServiceTag":
			if report.ServiceTag == value {
				return &report, nil
			}
		case "HostName":
			if report.HostName == value {
				return &report, nil
			}
		case "ID":
			if report.ID == value {
				return &report, nil
			}
		case "Compliant":
			compliant, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("invalid value for Compliant filter: %w", err)
			}
			if report.Compliant == compliant {
				return &report, nil
			}
		}
	}
	return nil, fmt.Errorf("no compliance report found matching the filter: %s with value: %s", filter, value)
}
