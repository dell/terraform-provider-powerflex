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

// NvmeHostResource is the model for NvmeHost Resource
type NvmeHostResource struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SystemID       types.String `tfsdk:"system_id"`
	Nqn            types.String `tfsdk:"nqn"`
	MaxNumPaths    types.String `tfsdk:"max_num_paths"`
	MaxNumSysPorts types.String `tfsdk:"max_num_sys_ports"`
}

// NvmeHostDataSource defines the model for NvmeHost Datasource
type NvmeHostDataSource struct {
	ID      types.String    `tfsdk:"id"`
	Details []NvmeHostModel `tfsdk:"nvme_host_details"`
	Filter  *IDNameFilter   `tfsdk:"filter"`
}

// IDNameFilter defines the model for filter IDs or Names
type IDNameFilter struct {
	Names []types.String `tfsdk:"names"`
	IDs   []types.String `tfsdk:"ids"`
}

// NvmeHostModel is the model for NVMe host
type NvmeHostModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SystemID       types.String `tfsdk:"system_id"`
	Nqn            types.String `tfsdk:"nqn"`
	MaxNumPaths    types.Int64  `tfsdk:"max_num_paths"`
	MaxNumSysPorts types.Int64  `tfsdk:"max_num_sys_ports"`
	Links          []LinkModel  `tfsdk:"links"`
}
