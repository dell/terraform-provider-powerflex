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

package helper

import (
	"context"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// UpdateSystemState updates the state of a system
func UpdateSystemState(plan *models.SystemModel, system *scaleiotypes.System, r *goscaleio.System) (state *models.SystemModel, diags diag.Diagnostics) {
	var tfStateApprovedIPs []attr.Value
	state = plan
	state.ID = types.StringValue(system.ID)
	state.RestrictedMode = types.StringValue(system.RestrictedSdcMode)
	sdcGuids := make([]string, 0)
	diags.Append(plan.SdcGuids.ElementsAs(context.TODO(), &sdcGuids, true)...)

	if len(sdcGuids) > 0 {
		stateGuids := make([]string, 0)

		for _, guid := range sdcGuids {
			sdc, err := r.FindSdc("SdcGUID", guid)
			if err != nil {
				diags.AddError(
					"Error getting SDC with GUID: ",
					"unexpected error: "+err.Error(),
				)
				return state, diags
			}
			if sdc.Sdc.SdcApproved {
				stateGuids = append(stateGuids, sdc.Sdc.SdcGUID)
			}
		}
		state.SdcGuids, _ = types.ListValueFrom(context.TODO(), types.StringType, stateGuids)
	}

	if !plan.SdcApprovedIPs.IsNull() {
		planApprovedIPs := make([]models.SdcApprovedIPsModel, 0)
		diags.Append(plan.SdcApprovedIPs.ElementsAs(context.TODO(), &planApprovedIPs, true)...)

		for _, sdcPlan := range planApprovedIPs {
			sdc, err := r.FindSdc("ID", sdcPlan.ID.ValueString())
			if err != nil {
				diags.AddError(
					"Error getting SDC with ID: "+sdcPlan.ID.ValueString(),
					"unexpected error: "+err.Error(),
				)
				continue
			}

			model, dgs := GetSdcApprovedIPsValue(sdc)
			diags = append(diags, dgs...)
			tfStateApprovedIPs = append(tfStateApprovedIPs, model)

		}
		state.SdcApprovedIPs, _ = types.ListValue(types.ObjectType{AttrTypes: GetSdcApprovedIPsType()}, tfStateApprovedIPs)
	}
	return state, diags
}

// GetSdcApprovedIPsType returns a map with the keys "id" and "ips" and their corresponding types.
func GetSdcApprovedIPsType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":  types.StringType,
		"ips": types.SetType{ElemType: types.StringType},
	}
}

// GetSdcApprovedIPsValue retrieves the SDC approved IPs value from the given Sdc object
func GetSdcApprovedIPsValue(sdc *goscaleio.Sdc) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	ips := []attr.Value{}

	for _, ip := range sdc.Sdc.SdcApprovedIPs {
		ips = append(ips, types.StringValue(ip))
	}
	ips1, _ := types.SetValue(types.StringType, ips)

	object, dia := types.ObjectValue(GetSdcApprovedIPsType(), map[string]attr.Value{
		"id":  types.StringValue(sdc.Sdc.ID),
		"ips": ips1,
	})
	diags.Append(dia...)
	return object, dia
}
