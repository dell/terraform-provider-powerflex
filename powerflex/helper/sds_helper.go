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

package helper

import (
	"context"

	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SdsIPListDiff get difference between sets of IP in state and plan
func SdsIPListDiff(ctx context.Context, plan, state *models.SdsResourceModel) (toAdd, toRmv, changed, common []*scaleiotypes.SdsIP) {
	plist, slist := plan.GetIPList(ctx), state.GetIPList(ctx)
	type ipObj struct {
		pip *scaleiotypes.SdsIP
		sip *scaleiotypes.SdsIP
	}
	vmap := make(map[string]*ipObj)
	for _, pip := range plist {
		vmap[pip.IP] = &ipObj{pip, nil}
	}
	for _, sip := range slist {
		if mip, ok := vmap[sip.IP]; ok {
			mip.sip = sip
		} else {
			vmap[sip.IP] = &ipObj{nil, sip}
		}
	}
	toAdd, toRmv, common, changed = make([]*scaleiotypes.SdsIP, 0), make([]*scaleiotypes.SdsIP, 0),
		make([]*scaleiotypes.SdsIP, 0), make([]*scaleiotypes.SdsIP, 0)
	for _, mip := range vmap {
		if mip.sip != nil {
			if mip.pip != nil {
				if mip.pip.Role == mip.sip.Role {
					common = append(common, mip.pip)
				} else {
					changed = append(changed, mip.pip)
				}
			} else {
				toRmv = append(toRmv, mip.sip)
			}
		} else {
			toAdd = append(toAdd, mip.pip)
		}
	}
	return toAdd, toRmv, changed, common
}

func UpdateSdsState(sds *scaleiotypes.Sds, plan models.SdsResourceModel) (models.SdsResourceModel, diag.Diagnostics) {
	state := plan
	state.ID = types.StringValue(sds.ID)
	state.Name = types.StringValue(sds.Name)
	state.ProtectionDomainID = types.StringValue(sds.ProtectionDomainID)
	state.Port = types.Int64Value(int64(sds.Port))
	state.SdsState = types.StringValue(sds.SdsState)
	state.MembershipState = types.StringValue(sds.MembershipState)
	state.MdmConnectionState = types.StringValue(sds.MdmConnectionState)
	state.DrlMode = types.StringValue(sds.DrlMode)
	state.RmcacheEnabled = types.BoolValue(sds.RmcacheEnabled)
	state.RmcacheSizeInMB = types.Int64Value(int64(sds.RmcacheSizeInKb) / 1024)
	state.RfcacheEnabled = types.BoolValue(sds.RfcacheEnabled)
	state.RmcacheFrozen = types.BoolValue(sds.RmcacheFrozen)
	state.IsOnVMware = types.BoolValue(sds.IsOnVMware)
	state.FaultSetID = types.StringValue(sds.FaultSetID)
	state.NumOfIoBuffers = types.Int64Value(int64(sds.NumOfIoBuffers))
	state.RmcacheMemoryAllocationState = types.StringValue(sds.RmcacheMemoryAllocationState)
	state.PerformanceProfile = types.StringValue(sds.PerformanceProfile)

	IPAttrTypes := map[string]attr.Type{
		"ip":   types.StringType,
		"role": types.StringType,
	}
	IPElemType := types.ObjectType{
		AttrTypes: IPAttrTypes,
	}

	objectIPs := []attr.Value{}
	var diags diag.Diagnostics
	for _, ip := range sds.IPList {
		obj := map[string]attr.Value{
			"ip":   types.StringValue(ip.IP),
			"role": types.StringValue(ip.Role),
		}
		objVal, dgs := types.ObjectValue(IPAttrTypes, obj)
		diags = append(diags, dgs...)
		objectIPs = append(objectIPs, objVal)
	}
	setVal, dgs := types.SetValue(IPElemType, objectIPs)
	diags = append(diags, dgs...)
	state.IPList = setVal

	return state, diags
}
