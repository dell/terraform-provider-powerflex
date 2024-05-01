/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// HostResourceModel struct for CSV Data Processing
type HostResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	HostDetails        types.List   `tfsdk:"host_details"`
	OSFamily           types.String `tfsdk:"os_family"`
	Name               types.String `tfsdk:"name"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	PackagePath        types.String `tfsdk:"package_path"`
	DriverConfigPath   types.String `tfsdk:"driver_cfg_path"`
	Credential         types.List   `tfsdk:"credential"`
	IP                 types.String `tfsdk:"ip"`
	GUID               types.String `tfsdk:"guid"`
	MdmIPs             types.List   `tfsdk:"mdm_ips"` // list(string)
}

// HostDetailDataModel defines the struct for CSV Parse Data
type HostDetailModel struct {
	HostID             types.String `tfsdk:"host_id"`
	IP                 types.String `tfsdk:"ip"`
	OperatingSystem    types.String `tfsdk:"operating_system"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	HostName           types.String `tfsdk:"host_name"`
	SystemID           types.String `tfsdk:"system_id"`
	IsApproved         types.Bool   `tfsdk:"is_approved"`
	OnVMWare           types.Bool   `tfsdk:"on_vmware"`
	HostGUID           types.String `tfsdk:"host_guid"`
	MdmConnectionState types.String `tfsdk:"mdm_connection_state"`
}

// CredentialModel defines the struct for Connection Data
type CredentialModel struct {
	UserName types.String `tfsdk:"user_name"`
	Password types.String `tfsdk:"password"`
}
