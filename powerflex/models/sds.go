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
	"context"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SDSResourceModel maps the resource schema data.
type SdsResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	ProtectionDomainID           types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName         types.String `tfsdk:"protection_domain_name"`
	IPList                       types.Set    `tfsdk:"ip_list"`
	Port                         types.Int64  `tfsdk:"port"`
	SdsState                     types.String `tfsdk:"sds_state"`
	MembershipState              types.String `tfsdk:"membership_state"`
	MdmConnectionState           types.String `tfsdk:"mdm_connection_state"`
	DrlMode                      types.String `tfsdk:"drl_mode"`
	RmcacheEnabled               types.Bool   `tfsdk:"rmcache_enabled"`
	RmcacheSizeInMB              types.Int64  `tfsdk:"rmcache_size_in_mb"`
	RfcacheEnabled               types.Bool   `tfsdk:"rfcache_enabled"`
	RmcacheFrozen                types.Bool   `tfsdk:"rmcache_frozen"`
	IsOnVMware                   types.Bool   `tfsdk:"is_on_vmware"`
	FaultSetID                   types.String `tfsdk:"fault_set_id"`
	NumOfIoBuffers               types.Int64  `tfsdk:"num_of_io_buffers"`
	RmcacheMemoryAllocationState types.String `tfsdk:"rmcache_memory_allocation_state"`
	PerformanceProfile           types.String `tfsdk:"performance_profile"`
}

// SDS IP object
type SdsIPModel struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

// Conversion of list of IPs from tf model to go type
func (sds *SdsResourceModel) GetIPList(ctx context.Context) []*scaleiotypes.SdsIP {
	iplist := []*scaleiotypes.SdsIP{}
	var ipModellist []SdsIPModel
	sds.IPList.ElementsAs(ctx, &ipModellist, false)
	for _, v := range ipModellist {
		sdsIP := scaleiotypes.SdsIP{
			IP:   v.IP.ValueString(),
			Role: v.Role.ValueString(),
		}
		iplist = append(iplist, &sdsIP)
	}
	return iplist
}
