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

// DeviceModel defines the struct for device resource
type DeviceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	DevicePath               types.String `tfsdk:"device_path"`
	DeviceOriginalPath       types.String `tfsdk:"device_original_path"`
	ProtectionDomainName     types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID       types.String `tfsdk:"protection_domain_id"`
	StoragePoolName          types.String `tfsdk:"storage_pool_name"`
	StoragePoolID            types.String `tfsdk:"storage_pool_id"`
	SdsID                    types.String `tfsdk:"sds_id"`
	SdsName                  types.String `tfsdk:"sds_name"`
	MediaType                types.String `tfsdk:"media_type"`
	ExternalAccelerationType types.String `tfsdk:"external_acceleration_type"`
	DeviceCapacity           types.Int64  `tfsdk:"device_capacity"`
	DeviceCapacityInKB       types.Int64  `tfsdk:"device_capacity_in_kb"`
	DeviceState              types.String `tfsdk:"device_state"`
}
