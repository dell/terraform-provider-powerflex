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

// NvmeTargetResourceModel is the model for NvmeTarget Resource
type NvmeTargetResourceModel struct {
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	SystemID                          types.String `tfsdk:"system_id" json:"systemId"`
	ProtectionDomainID                types.String `tfsdk:"protection_domain_id" json:"protectionDomainId"`
	IPList                            []IPList     `tfsdk:"ip_list" json:"ipList"`
	StoragePort                       types.Int64  `tfsdk:"storage_port" json:"storagePort"`
	NvmePort                          types.Int64  `tfsdk:"nvme_port" json:"nvmePort"`
	DiscoveryPort                     types.Int64  `tfsdk:"discovery_port" json:"discoveryPort"`
	SdtState                          types.String `tfsdk:"sdt_state" json:"sdtState"`
	MdmConnectionState                types.String `tfsdk:"mdm_connection_state" json:"mdmConnectionState"`
	MembershipState                   types.String `tfsdk:"membership_state" json:"membershipState"`
	FaultSetID                        types.String `tfsdk:"fault_set_id" json:"faultSetId"`
	SoftwareVersionInfo               types.String `tfsdk:"software_version_info" json:"softwareVersionInfo"`
	MaintenanceState                  types.String `tfsdk:"maintenance_state" json:"maintenanceState"`
	AuthenticationError               types.String `tfsdk:"authentication_error" json:"authenticationError"`
	PersistentDiscoveryControllersNum types.Int64  `tfsdk:"persistent_discovery_controllers_num" json:"persistentDiscoveryControllersNum"`
}

// NvmeTargetDataSource defines the model for NvmeTarget Datasource
type NvmeTargetDataSource struct {
	ID      types.String                `tfsdk:"id"`
	Details []NvmeTargetDatasourceModel `tfsdk:"nvme_target_details"`
	Filter  *NvmeTargetFilter           `tfsdk:"filter"`
}

// NvmeTargetFilter defines the model for NvmeTarget filter
type NvmeTargetFilter struct {
	ID                                []types.String `tfsdk:"id"`
	Name                              []types.String `tfsdk:"name"`
	SystemID                          []types.String `tfsdk:"system_id"`
	ProtectionDomainIDn               []types.String `tfsdk:"protection_domain_id"`
	StoragePort                       []types.Int64  `tfsdk:"storage_port"`
	NvmePort                          []types.Int64  `tfsdk:"nvme_port"`
	DiscoveryPort                     []types.Int64  `tfsdk:"discovery_port"`
	SdtState                          []types.String `tfsdk:"sdt_state"`
	MdmConnectionState                []types.String `tfsdk:"mdm_connection_state"`
	MembershipState                   []types.String `tfsdk:"membership_state"`
	FaultSetID                        []types.String `tfsdk:"fault_set_id"`
	SoftwareVersionInfo               []types.String `tfsdk:"software_version_info"`
	MaintenanceState                  []types.String `tfsdk:"maintenance_state"`
	AuthenticationError               []types.String `tfsdk:"authentication_error"`
	PersistentDiscoveryControllersNum []types.Int64  `tfsdk:"persistent_discovery_controllers_num"`
}

// NvmeTargetDatasourceModel is the datasource model for NVMe target
type NvmeTargetDatasourceModel struct {
	ID                                types.String         `tfsdk:"id"`
	Name                              types.String         `tfsdk:"name"`
	SystemID                          types.String         `tfsdk:"system_id" json:"systemId"`
	ProtectionDomainID                types.String         `tfsdk:"protection_domain_id" json:"protectionDomainId"`
	IPList                            []IPList             `tfsdk:"ip_list" json:"ipList"`
	StoragePort                       types.Int64          `tfsdk:"storage_port" json:"storagePort"`
	NvmePort                          types.Int64          `tfsdk:"nvme_port" json:"nvmePort"`
	DiscoveryPort                     types.Int64          `tfsdk:"discovery_port" json:"discoveryPort"`
	SdtState                          types.String         `tfsdk:"sdt_state" json:"sdtState"`
	MdmConnectionState                types.String         `tfsdk:"mdm_connection_state" json:"mdmConnectionState"`
	MembershipState                   types.String         `tfsdk:"membership_state" json:"membershipState"`
	FaultSetID                        types.String         `tfsdk:"fault_set_id" json:"faultSetId"`
	SoftwareVersionInfo               types.String         `tfsdk:"software_version_info" json:"softwareVersionInfo"`
	MaintenanceState                  types.String         `tfsdk:"maintenance_state" json:"maintenanceState"`
	AuthenticationError               types.String         `tfsdk:"authentication_error" json:"authenticationError"`
	PersistentDiscoveryControllersNum types.Int64          `tfsdk:"persistent_discovery_controllers_num" json:"persistentDiscoveryControllersNum"`
	CertificateInfo                   CertificateInfoModel `tfsdk:"certificate_info" json:"certificateInfo"`
	Links                             []LinkModel          `tfsdk:"links"`
	HostList                          []NvmeHostList       `tfsdk:"host_list"`
}

// NvmeHostList defines the model for NvmeHostList
type NvmeHostList struct {
	HostIP      types.String `tfsdk:"host_ip"`
	IsConnected types.Bool   `tfsdk:"is_connected"`
	HostName    types.String `tfsdk:"host_name"`
	HostID      types.String `tfsdk:"host_id"`
	SysPortIP   types.String `tfsdk:"sys_port_ip"`
}
