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

// SdcVolumeMappingResourceModel defines struct SDC volume mapping resource
type SdcVolumeMappingResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	VolumeList types.List   `tfsdk:"volume_list"`
}

// SdcVolumeModel defines struct volume mapping data
type SdcVolumeModel struct {
	VolumeID   types.String `tfsdk:"volume_id"`
	VolumeName types.String `tfsdk:"volume_name"`
	IOPSLimit  types.Int64  `tfsdk:"limit_iops"`
	BWLimit    types.Int64  `tfsdk:"limit_bw_in_mbps"`
	AccessMode types.String `tfsdk:"access_mode"`
}
