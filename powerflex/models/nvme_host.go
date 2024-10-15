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

// NvmeHostResourceModel is the model for NvmeHost Resource
type NvmeHostResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SystemID       types.String `tfsdk:"system_id" json:"systemId"`
	Nqn            types.String `tfsdk:"nqn"`
	MaxNumPaths    types.Int64  `tfsdk:"max_num_paths" json:"maxNumPaths"`
	MaxNumSysPorts types.Int64  `tfsdk:"max_num_sys_ports" json:"maxNumSysPorts"`
}

// NvmeHostDataSource defines the model for NvmeHost Datasource
type NvmeHostDataSource struct {
	ID      types.String              `tfsdk:"id"`
	Details []NvmeHostDatasourceModel `tfsdk:"nvme_host_details"`
	Filter  *NvmeHostFilter           `tfsdk:"filter"`
}

// NvmeHostFilter defines the model for NvmeHost filter
type NvmeHostFilter struct {
	Name           []types.String `tfsdk:"name"`
	ID             []types.String `tfsdk:"id"`
	SystemID       []types.String `tfsdk:"system_id"`
	Nqn            []types.String `tfsdk:"nqn"`
	MaxNumPaths    []types.Int64  `tfsdk:"max_num_paths"`
	MaxNumSysPorts []types.Int64  `tfsdk:"max_num_sys_ports"`
}

// NvmeHostDatasourceModel is the datasource model for NVMe host
type NvmeHostDatasourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SystemID       types.String `tfsdk:"system_id"`
	Nqn            types.String `tfsdk:"nqn"`
	MaxNumPaths    types.Int64  `tfsdk:"max_num_paths"`
	MaxNumSysPorts types.Int64  `tfsdk:"max_num_sys_ports"`
	Links          []LinkModel  `tfsdk:"links"`
}
