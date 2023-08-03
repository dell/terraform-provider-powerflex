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

	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetStandByMdmSetValue return the list for standby MDMs
func GetStandByMdmSetValue(mdms []models.StandByMdm) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	mdmInfoElemType := types.ObjectType{
		AttrTypes: GetStandByMdmType(),
	}

	if len(mdms) == 0 {
		return types.ListNull(mdmInfoElemType), diags
	}

	objectMdmInfos := []attr.Value{}
	for _, mdm := range mdms {
		obj := map[string]attr.Value{
			"id":                   mdm.ID,
			"name":                 mdm.Name,
			"port":                 mdm.Port,
			"role":                 mdm.Role,
			"ips":                  mdm.Ips,
			"management_ips":       mdm.ManagementIps,
			"allow_asymmetric_ips": mdm.AllowAsymmetricIps,
		}
		objVal, dgs := types.ObjectValue(GetStandByMdmType(), obj)
		diags = append(diags, dgs...)
		objectMdmInfos = append(objectMdmInfos, objVal)
	}
	mappedMdmInfoVal, dgs := types.ListValue(mdmInfoElemType, objectMdmInfos)
	diags = append(diags, dgs...)
	return mappedMdmInfoVal, diags
}

// GetStandByMdmType returns the type required for standby MDM
func GetStandByMdmType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"name":                 types.StringType,
		"port":                 types.Int64Type,
		"role":                 types.StringType,
		"ips":                  types.SetType{ElemType: types.StringType},
		"management_ips":       types.SetType{ElemType: types.StringType},
		"allow_asymmetric_ips": types.BoolType,
	}
}

// GetStandByMdmValue returns standby MDM object
func GetStandByMdmValue(ctx context.Context, mdm *goscaleio_types.Mdm, standby []models.StandByMdm) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	ips := []attr.Value{}
	for _, ip := range mdm.IPs {
		ips = append(ips, types.StringValue(ip))
	}
	ips1, _ := types.SetValue(types.StringType, ips)

	mgmtIps := []attr.Value{}
	for _, ip := range mdm.ManagementIPs {
		mgmtIps = append(mgmtIps, types.StringValue(ip))
	}
	mgmtIps1, _ := types.SetValue(types.StringType, ips)

	var asymmetricFlag bool

	for _, standbymdm := range standby {
		planIps := make([]string, 0)
		diags.Append(standbymdm.Ips.ElementsAs(ctx, &planIps, true)...)
		if CompareStringSlice(planIps, mdm.IPs) {
			asymmetricFlag = standbymdm.AllowAsymmetricIps.ValueBool()
		}
	}

	object, dia := types.ObjectValue(GetStandByMdmType(), map[string]attr.Value{
		"id":                   types.StringValue(mdm.ID),
		"name":                 types.StringValue(mdm.Name),
		"port":                 types.Int64Value(int64(mdm.Port)),
		"role":                 types.StringValue(mdm.Role),
		"ips":                  ips1,
		"management_ips":       mgmtIps1,
		"allow_asymmetric_ips": types.BoolValue(asymmetricFlag),
	})
	diags.Append(dia...)
	return object, dia
}

// UpdateMdmClusterState returns the state for MDM cluster
func UpdateMdmClusterState(ctx context.Context, mdmDetails *goscaleio_types.MdmCluster, plan *models.MdmResourceModel, perfProfile string) (*models.MdmResourceModel, diag.Diagnostics) {
	var state models.MdmResourceModel
	var dgs, diags diag.Diagnostics

	state.ID = types.StringValue(mdmDetails.ID)
	state.ClusterMode = types.StringValue(mdmDetails.ClusterMode)
	state.PerformanceProfile = types.StringValue(perfProfile)
	state.PrimaryMdm.ID = types.StringValue(mdmDetails.PrimaryMDM.ID)
	state.PrimaryMdm.Name = types.StringValue(mdmDetails.PrimaryMDM.Name)
	state.PrimaryMdm.Port = types.Int64Value(int64(mdmDetails.PrimaryMDM.Port))
	primaryMdmIps := []attr.Value{}
	for _, ip := range mdmDetails.PrimaryMDM.IPs {
		primaryMdmIps = append(primaryMdmIps, types.StringValue(ip))
	}
	state.PrimaryMdm.Ips, dgs = types.SetValue(types.StringType, primaryMdmIps)
	diags = append(diags, dgs...)

	primaryMdmMgmtIps := []attr.Value{}
	for _, ip := range mdmDetails.PrimaryMDM.ManagementIPs {
		primaryMdmMgmtIps = append(primaryMdmMgmtIps, types.StringValue(ip))
	}
	state.PrimaryMdm.ManagementIps, dgs = types.SetValue(types.StringType, primaryMdmMgmtIps)
	diags = append(diags, dgs...)

	secondaryMdmDetails := make([]models.Mdm, 0)
	for _, mdm := range mdmDetails.SecondaryMDM {
		secondaryMdmIps := []attr.Value{}
		for _, ip := range mdm.IPs {
			secondaryMdmIps = append(secondaryMdmIps, types.StringValue(ip))
		}
		mdmIps, dgs1 := types.SetValue(types.StringType, secondaryMdmIps)
		diags = append(diags, dgs1...)

		secondaryMdmMgmtIps := []attr.Value{}
		for _, ip := range mdm.ManagementIPs {
			secondaryMdmMgmtIps = append(secondaryMdmMgmtIps, types.StringValue(ip))
		}
		mgmtIps, dgs2 := types.SetValue(types.StringType, secondaryMdmMgmtIps)
		diags = append(diags, dgs2...)

		secondaryMdm := models.Mdm{
			ID:            types.StringValue(mdm.ID),
			Name:          types.StringValue(mdm.Name),
			Port:          types.Int64Value(int64(mdm.Port)),
			Ips:           mdmIps,
			ManagementIps: mgmtIps,
		}
		secondaryMdmDetails = append(secondaryMdmDetails, secondaryMdm)
	}
	state.SecondaryMdm = secondaryMdmDetails

	tbMdmDetails := make([]models.Mdm, 0)
	for _, mdm := range mdmDetails.TiebreakerMdm {
		tbMdmIps := []attr.Value{}
		for _, ip := range mdm.IPs {
			tbMdmIps = append(tbMdmIps, types.StringValue(ip))
		}
		mdmIps, dgs1 := types.SetValue(types.StringType, tbMdmIps)
		diags = append(diags, dgs1...)

		tbMdmMgmtIps := []attr.Value{}
		for _, ip := range mdm.ManagementIPs {
			tbMdmMgmtIps = append(tbMdmMgmtIps, types.StringValue(ip))
		}
		mgmtIps, dgs2 := types.SetValue(types.StringType, tbMdmMgmtIps)
		diags = append(diags, dgs2...)

		tbMdm := models.Mdm{
			ID:            types.StringValue(mdm.ID),
			Name:          types.StringValue(mdm.Name),
			Port:          types.Int64Value(int64(mdm.Port)),
			Ips:           mdmIps,
			ManagementIps: mgmtIps,
		}
		tbMdmDetails = append(tbMdmDetails, tbMdm)
	}
	state.TieBreakerMdm = tbMdmDetails

	var sbMdmDetails []attr.Value
	var sbMdmList []models.StandByMdm
	diags.Append(plan.StandByMdm.ElementsAs(ctx, &sbMdmList, true)...)
	for _, mdm := range mdmDetails.StandByMdm {
		model, dgs := GetStandByMdmValue(ctx, &mdm, sbMdmList)
		diags = append(diags, dgs...)
		sbMdmDetails = append(sbMdmDetails, model)
	}
	state.StandByMdm, _ = types.ListValue(types.ObjectType{AttrTypes: GetStandByMdmType()}, sbMdmDetails)

	return &state, diags
}

// StandByMdmDifference returns the diff between plan and state standby MDMs
func StandByMdmDifference(ctx context.Context, state, plan []models.StandByMdm) ([]models.StandByMdm, diag.Diagnostics) {
	difference := make([]models.StandByMdm, 0)
	var dia diag.Diagnostics

	for _, mdm1 := range state {
		found := false
		for _, mdm2 := range plan {
			stateIps := make([]string, 0)
			planIps := make([]string, 0)
			dia.Append(mdm1.Ips.ElementsAs(ctx, &stateIps, true)...)
			dia.Append(mdm2.Ips.ElementsAs(ctx, &planIps, true)...)
			if CompareStringSlice(stateIps, planIps) {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, mdm1)
		}
	}
	return difference, dia
}

// CheckForSwitchCluster checks if new standby MDMs needs to be added
func CheckForSwitchCluster(ctx context.Context, standby []models.StandByMdm, stateSecondary, stateTb []models.Mdm) ([]models.StandByMdm, diag.Diagnostics) {
	standByFinal := make([]models.StandByMdm, 0)
	var dia diag.Diagnostics
	for _, mdm := range standby {
		found := false
		sbIps := make([]string, 0)
		dia.Append(mdm.Ips.ElementsAs(ctx, &sbIps, true)...)
		if mdm.Role.ValueString() == goscaleio_types.Manager {
			for _, sec := range stateSecondary {
				if !sec.Ips.IsUnknown() {
					secIps := make([]string, 0)
					dia.Append(sec.Ips.ElementsAs(ctx, &secIps, true)...)
					if CompareStringSlice(sbIps, secIps) {
						found = true
						break
					}
				} else if mdm.ID.ValueString() == sec.ID.ValueString() {
					found = true
					break
				}
			}
		} else {
			for _, tb := range stateTb {
				if !tb.Ips.IsUnknown() {
					tbIps := make([]string, 0)
					dia.Append(tb.Ips.ElementsAs(ctx, &tbIps, true)...)
					if CompareStringSlice(sbIps, tbIps) {
						found = true
						break
					}
				} else if mdm.ID.ValueString() == tb.ID.ValueString() {
					found = true
					break
				}
			}
		}
		if !found {
			standByFinal = append(standByFinal, mdm)
		}
	}
	return standByFinal, dia
}

// GetMdmIPMap returns the map with IDs as value
func GetMdmIPMap(mdmDetails *goscaleio_types.MdmCluster) map[string]string {
	ipmap := make(map[string]string)

	ipmap[mdmDetails.PrimaryMDM.IPs[0]] = mdmDetails.PrimaryMDM.ID

	for _, mdm := range mdmDetails.SecondaryMDM {
		ipmap[mdm.IPs[0]] = mdm.ID
	}

	for _, mdm := range mdmDetails.TiebreakerMdm {
		ipmap[mdm.IPs[0]] = mdm.ID
	}

	for _, mdm := range mdmDetails.StandByMdm {
		ipmap[mdm.IPs[0]] = mdm.ID
	}
	return ipmap
}

// CheckforExistingName checks for the given name in map
func CheckforExistingName(plan models.MdmResourceModel, names map[string]string) bool {
	if _, ok := names[plan.PrimaryMdm.Name.ValueString()]; ok {
		return false
	}
	return true
}
