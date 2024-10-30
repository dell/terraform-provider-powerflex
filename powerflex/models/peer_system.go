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

import "github.com/hashicorp/terraform-plugin-framework/types"

// PeerMdmModel model for Peer Mdm
type PeerMdmModel struct {
	ID                  types.String    `tfsdk:"id"`
	Name                types.String    `tfsdk:"name"`
	Port                types.Int64     `tfsdk:"port"`
	PeerSystemID        types.String    `tfsdk:"peer_system_id"`
	SystemID            types.String    `tfsdk:"system_id"`
	SoftwareVersionInfo types.String    `tfsdk:"software_version_info"`
	MembershipState     types.String    `tfsdk:"membership_state"`
	PerfProfile         types.String    `tfsdk:"perf_profile"`
	NetworkType         types.String    `tfsdk:"network_type"`
	CouplingRC          types.String    `tfsdk:"coupling_rc"`
	IPList              []*IPListNoRole `tfsdk:"ip_list"`
}

// IPListNoRole model for Peer Mdm
type IPListNoRole struct {
	IP types.String `tfsdk:"ip"`
}

// PeerMdmDataSourceModel model for Peer Mdm
type PeerMdmDataSourceModel struct {
	PeerMdmDetails []PeerMdmModel `tfsdk:"peer_system_details"`
	ID             types.String   `tfsdk:"id"`
	PeerMdmFilter  *PeerMdmFilter `tfsdk:"filter"`
}

// PeerMdmFilter defines the model for filters used for PeerMdmDataSource
type PeerMdmFilter struct {
	ID                  []types.String `tfsdk:"id"`
	Name                []types.String `tfsdk:"name"`
	Port                []types.Int64  `tfsdk:"port"`
	PeerSystemID        []types.String `tfsdk:"peer_system_id"`
	SystemID            []types.String `tfsdk:"system_id"`
	SoftwareVersionInfo []types.String `tfsdk:"software_version_info"`
	MembershipState     []types.String `tfsdk:"membership_state"`
	PerfProfile         []types.String `tfsdk:"perf_profile"`
	NetworkType         []types.String `tfsdk:"network_type"`
	CouplingRC          []types.String `tfsdk:"coupling_rc"`
}

// PeerMdmResourceModel model for Peer Mdm
type PeerMdmResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	Name                      types.String   `tfsdk:"name"`
	Port                      types.Int64    `tfsdk:"port"`
	PeerSystemID              types.String   `tfsdk:"peer_system_id"`
	SystemID                  types.String   `tfsdk:"system_id"`
	SoftwareVersionInfo       types.String   `tfsdk:"software_version_info"`
	MembershipState           types.String   `tfsdk:"membership_state"`
	PerfProfile               types.String   `tfsdk:"perf_profile"`
	NetworkType               types.String   `tfsdk:"network_type"`
	CouplingRC                types.String   `tfsdk:"coupling_rc"`
	DestinationPrimaryMdmInfo types.Object   `tfsdk:"destination_primary_mdm_information"`
	SourcePrimaryMdmInfo      types.Object   `tfsdk:"source_primary_mdm_information"`
	AddCertificate            types.Bool     `tfsdk:"add_certificate"`
	IPList                    []types.String `tfsdk:"ip_list"`
}

// PrimaryMdmInfo model for Primary a Mdm
type PrimaryMdmInfo struct {
	IP                 string  `tfsdk:"ip"`
	Username           string  `tfsdk:"ssh_username"`
	Password           *string `tfsdk:"ssh_password"`
	Port               string  `tfsdk:"ssh_port"`
	ManagementIP       string  `tfsdk:"management_ip"`
	ManagementUsername string  `tfsdk:"management_username"`
	ManagementPassword *string `tfsdk:"management_password"`
}
