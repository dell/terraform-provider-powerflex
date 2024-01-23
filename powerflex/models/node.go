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

// NodeDataSourceModel maps the struct to Node data source schema
type NodeDataSourceModel struct {
	NodeIDs       types.Set    `tfsdk:"node_ids"`
	ServiceTags   types.Set    `tfsdk:"service_tags"`
	IPAddresses   types.Set    `tfsdk:"ip_addresses"`
	NodePoolIDs   types.Set    `tfsdk:"node_pool_ids"`
	NodePoolNames types.Set    `tfsdk:"node_pool_names"`
	NodeDetails   []NodeModel  `tfsdk:"node_details"`
	ID            types.String `tfsdk:"id"`
}

// NodeModel maps the struct to Node schema
type NodeModel struct {
	RefID               types.String    `tfsdk:"ref_id"`
	IPAddress           types.String    `tfsdk:"ip_address"`
	CurrentIPAddress    types.String    `tfsdk:"current_ip_address"`
	ServiceTag          types.String    `tfsdk:"service_tag"`
	Model               types.String    `tfsdk:"model"`
	DeviceType          types.String    `tfsdk:"device_type"`
	DiscoverDeviceType  types.String    `tfsdk:"discover_device_type"`
	DisplayName         types.String    `tfsdk:"display_name"`
	ManagedState        types.String    `tfsdk:"managed_state"`
	State               types.String    `tfsdk:"state"`
	InUse               types.Bool      `tfsdk:"in_use"`
	CustomFirmware      types.Bool      `tfsdk:"custom_firmware"`
	NeedsAttention      types.Bool      `tfsdk:"needs_attention"`
	Manufacturer        types.String    `tfsdk:"manufacturer"`
	SystemID            types.String    `tfsdk:"system_id"`
	Health              types.String    `tfsdk:"health"`
	HealthMessage       types.String    `tfsdk:"health_message"`
	OperatingSystem     types.String    `tfsdk:"operating_system"`
	NumberOfCPUs        types.Int64     `tfsdk:"number_of_cpus"`
	Nics                types.Int64     `tfsdk:"nics"`
	MemoryInGB          types.Int64     `tfsdk:"memory_in_gb"`
	ComplianceCheckDate types.String    `tfsdk:"compliance_check_date"`
	DiscoveredDate      types.String    `tfsdk:"discovered_date"`
	DeviceGroupList     DeviceGroupList `tfsdk:"device_group_list"`
	CredID              types.String    `tfsdk:"cred_id"`
	Compliance          types.String    `tfsdk:"compliance"`
	FailuresCount       types.Int64     `tfsdk:"failures_count"`
	Facts               types.String    `tfsdk:"facts"`
	PuppetCertName      types.String    `tfsdk:"puppet_cert_name"`
	FlexosMaintMode     types.Int64     `tfsdk:"flex_os_maint_mode"`
	EsxiMaintMode       types.Int64     `tfsdk:"esxi_maint_mode"`
}

// DeviceGroupList defines struct for devices
type DeviceGroupList struct {
	DeviceGroup []DeviceGroup `tfsdk:"device_group"`
}

// DeviceGroup defines struct for nodepool
type DeviceGroup struct {
	GroupSeqID       int           `tfsdk:"group_seq_id"`
	GroupName        string        `tfsdk:"group_name"`
	GroupDescription string        `tfsdk:"group_description"`
	CreatedDate      string        `tfsdk:"created_date"`
	CreatedBy        string        `tfsdk:"created_by"`
	UpdatedDate      string        `tfsdk:"updated_date"`
	UpdatedBy        string        `tfsdk:"updated_by"`
	GroupUserList    GroupUserList `tfsdk:"group_user_list"`
}

// GroupUserList defines struct for group users
type GroupUserList struct {
	TotalRecords int          `tfsdk:"total_records"`
	GroupUsers   []GroupUsers `tfsdk:"group_users"`
}

// GroupUsers defines struct for group user
type GroupUsers struct {
	UserSeqID string `tfsdk:"user_seq_id"`
	UserName  string `tfsdk:"user_name"`
	FirstName string `tfsdk:"first_name"`
	LastName  string `tfsdk:"last_name"`
	Role      string `tfsdk:"role"`
	Enabled   bool   `tfsdk:"enabled"`
}
