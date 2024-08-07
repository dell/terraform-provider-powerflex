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

// SdcHostModel maps the resource schema data.
type SdcHostModel struct {
	ID     types.String `tfsdk:"id"`
	Remote types.Object `tfsdk:"remote"` // SdcHostRemoteModel
	Host   types.String `tfsdk:"ip"`
	Pkg    types.String `tfsdk:"package_path"`
	OS     types.String `tfsdk:"os_family"`

	// optional
	Name               types.String `tfsdk:"name"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	MdmIPs             types.List   `tfsdk:"clusters_mdm_ips"`
	UseRemotePath      types.Bool   `tfsdk:"use_remote_path"`

	// optional, os specific
	Esxi types.Object `tfsdk:"esxi"`

	// Optional Computed with defaults
	LinuxDrvCfg   types.String `tfsdk:"linux_drv_cfg"`
	WindowsDrvCfg types.String `tfsdk:"windows_drv_cfg"`

	//Computed
	SystemID           types.String `tfsdk:"system_id"`
	IsApproved         types.Bool   `tfsdk:"is_approved"`
	OnVMWare           types.Bool   `tfsdk:"on_vmware"`
	GUID               types.String `tfsdk:"guid"`
	MdmConnectionState types.String `tfsdk:"mdm_connection_state"`
}

// SdcHostRemoteModel maps the remote schema data.
type SdcHostRemoteModel struct {
	Port       string  `tfsdk:"port"`
	User       string  `tfsdk:"user"`
	Password   *string `tfsdk:"password"`
	PrivateKey *string `tfsdk:"private_key"`
	CaCert     *string `tfsdk:"certificate"`
	HostKey    *string `tfsdk:"host_key"`
	Dir        *string `tfsdk:"dir"`
}

// SdcHostEsxiModel maps the esxi resource schema data.
type SdcHostEsxiModel struct {
	GUID          types.String `tfsdk:"guid"`
	VerifyVibSign types.Bool   `tfsdk:"verify_vib_signature"`
}
