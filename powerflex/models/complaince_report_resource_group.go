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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ComplianceReportResourceGroupDatasource defines the compliance report Datasource object.
type ComplianceReportResourceGroupDatasource struct {
	ComplianceReports      []ComplianceReportResourceGroup `tfsdk:"compliance_reports"`
	ID                     types.String                    `tfsdk:"id"`
	ResourceGroupID        types.String                    `tfsdk:"resource_group_id"`
	ComplianceReportFilter *ComplianceReportFilterType     `tfsdk:"filter"`
}

// ComplianceReportFilterType defines the filter for datasource
type ComplianceReportFilterType struct {
	IPAddresses types.Set  `tfsdk:"ip_addresses"`
	ServiceTags types.Set  `tfsdk:"service_tags"`
	Compliant   types.Bool `tfsdk:"compliant"`
	HostNames   types.Set  `tfsdk:"host_names"`
	ResourceIDs types.Set  `tfsdk:"resource_ids"`
}

// ComplianceReportResourceGroup defines the compliance report for a service.
type ComplianceReportResourceGroup struct {
	ServiceTag                 types.String                 `tfsdk:"service_tag"`
	IPAddress                  types.String                 `tfsdk:"ip_address"`
	FirmwareRepositoryName     types.String                 `tfsdk:"firmware_repository_name"`
	ComplianceReportComponents []ComplianceReportComponents `tfsdk:"firmware_compliance_report_components"`
	Compliant                  types.Bool                   `tfsdk:"compliant"`
	DeviceType                 types.String                 `tfsdk:"device_type"`
	Model                      types.String                 `tfsdk:"model"`
	Available                  types.Bool                   `tfsdk:"available"`
	ManagedState               types.String                 `tfsdk:"managed_state"`
	EmbeddedReport             types.Bool                   `tfsdk:"embedded_report"`
	DeviceState                types.String                 `tfsdk:"device_state"`
	ID                         types.String                 `tfsdk:"id"`
	HostName                   types.String                 `tfsdk:"host_name"`
	CanUpdate                  types.Bool                   `tfsdk:"can_update"`
}

// ComplianceReportComponents defines the components in the compliance report.
type ComplianceReportComponents struct {
	ID              types.String                         `tfsdk:"id"`
	Name            types.String                         `tfsdk:"name"`
	CurrentVersion  ComplianceReportComponentVersionInfo `tfsdk:"current_version"`
	TargetVersion   ComplianceReportComponentVersionInfo `tfsdk:"target_version"`
	Vendor          types.String                         `tfsdk:"vendor"`
	OperatingSystem types.String                         `tfsdk:"operating_system"`
	Compliant       types.Bool                           `tfsdk:"compliant"`
	Software        types.Bool                           `tfsdk:"software"`
	RPM             types.Bool                           `tfsdk:"rpm"`
	Oscapable       types.Bool                           `tfsdk:"os_compatible"`
}

// ComplianceReportComponentVersionInfo defines the version info in the compliance report component.
type ComplianceReportComponentVersionInfo struct {
	ID                 types.String `tfsdk:"id"`
	FirmwareName       types.String `tfsdk:"firmware_name"`
	FirmwareType       types.String `tfsdk:"firmware_type"`
	FirmwareVersion    types.String `tfsdk:"firmware_version"`
	FirmwareLastUpdate types.String `tfsdk:"firmware_last_update"`
	FirmwareLevel      types.String `tfsdk:"firmware_level"`
}
