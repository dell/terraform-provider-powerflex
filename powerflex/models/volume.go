/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// VolumeResourceModel maps the resource schema data.
type VolumeResourceModel struct {
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	StoragePoolName      types.String `tfsdk:"storage_pool_name"`
	StoragePoolID        types.String `tfsdk:"storage_pool_id"`
	VolumeType           types.String `tfsdk:"volume_type"`
	UseRmCache           types.Bool   `tfsdk:"use_rm_cache"`
	CompressionMethod    types.String `tfsdk:"compression_method"`
	Size                 types.Int64  `tfsdk:"size"`
	CapacityUnit         types.String `tfsdk:"capacity_unit"`
	Name                 types.String `tfsdk:"name"`
	SizeInKb             types.Int64  `tfsdk:"size_in_kb"`
	ID                   types.String `tfsdk:"id"`
	AccessMode           types.String `tfsdk:"access_mode"`
	RemoveMode           types.String `tfsdk:"remove_mode"`
}

// SDCItemize maps the sdc_list schema data
type SDCItemize struct {
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName       types.String `tfsdk:"sdc_name"`
	AccessMode    types.String `tfsdk:"access_mode"`
}
