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

// SystemModel maps the struct to System resource schema
type SystemModel struct {
	ID             types.String `tfsdk:"id"`
	RestrictedMode types.String `tfsdk:"restricted_mode"`
	SdcGuids       types.List   `tfsdk:"sdc_guids"`
	SdcApprovedIPs types.List   `tfsdk:"sdc_approved_ips"`
	SdcIDs         types.List   `tfsdk:"sdc_ids"`
	SdcNames       types.List   `tfsdk:"sdc_names"`
}

// SdcApprovedIPsModel maps the struct to SdcApprovedIPs schema
type SdcApprovedIPsModel struct {
	ID  types.String `tfsdk:"id"`
	IPs types.Set    `tfsdk:"ips"`
}
