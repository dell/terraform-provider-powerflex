/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateClusterState - function to update state file for Cluster resource.
func UpdateClusterState(plan models.ClusterResourceModel, gatewayClient *goscaleio.GatewayClient, mdmIP string) (models.ClusterResourceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	clusteDetailResponse, error := GetClusterDetails(plan, gatewayClient, mdmIP, false)
	if error != nil {
		diags.AddError(error.Error(), error.Error())
		return plan, diags
	}

	//SDC Data
	SDCAttrTypes := GetSDCType()
	SDCElemType := types.ObjectType{
		AttrTypes: SDCAttrTypes,
	}

	//SDS Data
	SDSAttrTypes := GetSDSType()
	SDSElemType := types.ObjectType{
		AttrTypes: SDSAttrTypes,
	}

	// SDR Data
	SDRAttrTypes := GetSDRType()
	SDRElemType := types.ObjectType{
		AttrTypes: SDRAttrTypes,
	}

	// SDT Data
	SDTAttrTypes := GetSDTType()
	SDTElemType := types.ObjectType{
		AttrTypes: SDTAttrTypes,
	}

	//MDM Data
	MDMAttrTypes := GetMDMType()
	MDMElemType := types.ObjectType{
		AttrTypes: MDMAttrTypes,
	}

	//Protection Domain Data
	PDAttrTypes := GetPDType()
	PDElemType := types.ObjectType{
		AttrTypes: PDAttrTypes,
	}

	// Convert SdcList to SDCDataModel
	var sdcDataList []models.SDCModel
	for _, sdc := range clusteDetailResponse.ClusterDetails.SdcList {

		id, error := decimalToTwosComplementHex(sdc.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		sdcData := models.SDCModel{
			ID:   types.StringValue(id),
			GUID: types.StringValue(sdc.GUID),
			Name: types.StringValue(sdc.SdcName),
			IP:   types.StringValue(strings.Join(sdc.Node.NodeIPs, ",")),
		}
		sdcDataList = append(sdcDataList, sdcData)
	}

	objectSDCs := []attr.Value{}
	for _, sdc := range sdcDataList {
		objVal, dgs := GetSDCValue(sdc)
		diags = append(diags, dgs...)
		objectSDCs = append(objectSDCs, objVal)
	}
	setSdcs, dgs := types.SetValue(SDCElemType, objectSDCs)
	diags = append(diags, dgs...)

	state.SDCList = setSdcs

	state.ID = types.StringValue("placeholder")

	// Convert SdSList to SDSDataModel
	var sdsDataList []models.SDSModel
	objectSDSs := []attr.Value{}
	for _, sds := range clusteDetailResponse.ClusterDetails.SdsList {

		id, error := decimalToTwosComplementHex(sds.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		protectionDomainID, error := decimalToTwosComplementHex(sds.ProtectionDomainID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		var deviceDataList []models.DeviceDataModel

		for _, device := range sds.Devices {
			deviceData := models.DeviceDataModel{
				Name:            types.StringValue(device.DeviceName),
				Path:            types.StringValue(device.DevicePath),
				StoragePoolName: types.StringValue(device.StoragePool),
				MaxCapacity:     types.Int64Value(int64(device.MaxCapacityInKb)),
			}
			deviceDataList = append(deviceDataList, deviceData)
		}

		DeviceAttrTypes := GetDeviceType()
		DeviceElemType := types.ObjectType{
			AttrTypes: DeviceAttrTypes,
		}

		objectDevices := []attr.Value{}
		for _, device := range deviceDataList {
			objVal, dgs := GetDeviceValue(device)
			diags = append(diags, dgs...)
			objectDevices = append(objectDevices, objVal)
		}
		setDevices, dgs := types.SetValue(DeviceElemType, objectDevices)
		diags = append(diags, dgs...)

		sdsData := models.SDSModel{
			ID:                   types.StringValue(id),
			Name:                 types.StringValue(sds.SdsName),
			IP:                   types.StringValue(strings.Join(sds.Node.NodeIPs, ",")),
			AllIP:                types.StringValue(strings.Join(sds.AllIPs, ",")),
			SDSOnlyIP:            types.StringValue(strings.Join(sds.SdsOnlyIPs, ",")),
			SDSSDCIP:             types.StringValue(strings.Join(sds.SdcOnlyIPs, ",")),
			ProtectionDomainID:   types.StringValue(protectionDomainID),
			ProtectionDomainName: types.StringValue(sds.ProtectionDomain),
			FaultSet:             types.StringValue(sds.FaultSet),
			Devices:              setDevices,
		}

		sdsDataList = append(sdsDataList, sdsData)
	}

	for _, sds := range sdsDataList {
		objVal, dgs := GetSDSValue(sds)
		diags = append(diags, dgs...)
		objectSDSs = append(objectSDSs, objVal)
	}
	setSdSs, dgs := types.SetValue(SDSElemType, objectSDSs)
	diags = append(diags, dgs...)

	state.SDSList = setSdSs

	// Convert SdrList to SDRDataModel
	var sdrDataList []models.SDRModel
	objectSDRs := []attr.Value{}
	for _, sdr := range clusteDetailResponse.ClusterDetails.SdrList {

		id, error := decimalToTwosComplementHex(sdr.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		sdrData := models.SDRModel{
			ID:            types.StringValue(id),
			Name:          types.StringValue(sdr.SdrName),
			IP:            types.StringValue(strings.Join(sdr.Node.NodeIPs, ",")),
			AllIP:         types.StringValue(strings.Join(sdr.AllIPs, ",")),
			Port:          types.Int64Value(int64(sdr.SdrPort)),
			ApplicationIP: types.StringValue(strings.Join(sdr.ApplicationOnlyIPs, ",")),
			StorageIP:     types.StringValue(strings.Join(sdr.StorageOnlyIPs, ",")),
			ExternalIP:    types.StringValue(strings.Join(sdr.ExternalOnlyIPs, ",")),
		}

		sdrDataList = append(sdrDataList, sdrData)
	}

	for _, sdr := range sdrDataList {
		objVal, dgs := GetSDRValue(sdr)
		diags = append(diags, dgs...)
		objectSDRs = append(objectSDRs, objVal)
	}
	setSdrs, dgs := types.SetValue(SDRElemType, objectSDRs)
	diags = append(diags, dgs...)

	state.SDRList = setSdrs

	// Convert SdtList to SDTDataModel
	var sdtDataList []models.SDTModel
	objectSDTs := []attr.Value{}
	for _, sdt := range clusteDetailResponse.ClusterDetails.SdtList {

		id, error := decimalToTwosComplementHex(sdt.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		protectionDomainID, error := decimalToTwosComplementHex(sdt.ProtectionDomainID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		sdtData := models.SDTModel{
			ID:                   types.StringValue(id),
			Name:                 types.StringValue(sdt.SdtName),
			IP:                   types.StringValue(strings.Join(sdt.Node.NodeIPs, ",")),
			AllIP:                types.StringValue(strings.Join(sdt.AllIPs, ",")),
			StorageOnlyIPs:       types.StringValue(strings.Join(sdt.StorageOnlyIPs, ",")),
			HostOnlyIPs:          types.StringValue(strings.Join(sdt.HostOnlyIPs, ",")),
			ProtectionDomainID:   types.StringValue(protectionDomainID),
			ProtectionDomainName: types.StringValue(sdt.ProtectionDomain),
			StoragePort:          types.Int64Value(int64(sdt.StoragePort)),
			NvmePort:             types.Int64Value(int64(sdt.NvmePort)),
			DiscoveryPort:        types.Int64Value(int64(sdt.DiscoveryPort)),
		}

		sdtDataList = append(sdtDataList, sdtData)
	}

	for _, sdt := range sdtDataList {
		objVal, dgs := GetSDTValue(sdt)
		diags = append(diags, dgs...)
		objectSDTs = append(objectSDTs, objVal)
	}
	setSdts, dgs := types.SetValue(SDTElemType, objectSDTs)
	diags = append(diags, dgs...)

	state.SDTList = setSdts

	// Convert pdList to PDDataModel
	var pdDataList []models.ProtectionDomainDataModel
	objectPDs := []attr.Value{}

	for _, pd := range clusteDetailResponse.ClusterDetails.ProtectionDomains {

		var spDataList []models.StoragePoolDetailModel

		for _, sp := range pd.StoragePools {
			spData := models.StoragePoolDetailModel{
				Name:                                 types.StringValue(sp.Name),
				MediaType:                            types.StringValue(sp.MediaType),
				ExternalAcceleration:                 types.StringValue(sp.ExternalAccelerationType),
				DataLayout:                           types.StringValue(sp.DataLayout),
				ZeroPadding:                          types.StringValue(boolToYesNo(sp.ShouldApplyZeroPadding)),
				CompressionMethod:                    types.StringValue(sp.CompressionMethod),
				ReplicationJournalCapacityPercentage: types.Int64Value(int64(sp.RplJournalCapacity)),
			}
			spDataList = append(spDataList, spData)
		}

		SPAttrTypes := GetStoragePoolsType()
		SPElemType := types.ObjectType{
			AttrTypes: SPAttrTypes,
		}

		objectSPs := []attr.Value{}
		for _, sp := range spDataList {
			objVal, dgs := GetStoragePoolsValue(sp)
			diags = append(diags, dgs...)
			objectSPs = append(objectSPs, objVal)
		}
		setSPs, dgs := types.ListValue(SPElemType, objectSPs)
		diags = append(diags, dgs...)

		pdData := models.ProtectionDomainDataModel{
			Name:         types.StringValue(pd.Name),
			StoragePools: setSPs,
		}

		pdDataList = append(pdDataList, pdData)
	}

	for _, pd := range pdDataList {
		objVal, dgs := GetPDValue(pd)
		diags = append(diags, dgs...)
		objectPDs = append(objectPDs, objVal)
	}
	setPds, dgs := types.ListValue(PDElemType, objectPDs)
	diags = append(diags, dgs...)

	state.ProtectionDomains = setPds

	// Convert mdmList to MDMDataModel
	var mdmList []models.MDMModel
	objectMDMs := []attr.Value{}

	for _, mdm := range clusteDetailResponse.ClusterDetails.SlaveMdmSet {

		id, error := decimalToTwosComplementHex(mdm.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		mdmData := models.MDMModel{
			ID:           types.StringValue(id),
			Name:         types.StringValue(mdm.Name),
			IP:           types.StringValue(strings.Join(mdm.Node.NodeIPs, ",")),
			MDMIP:        types.StringValue(strings.Join(mdm.MdmIPs, ",")),
			MGMTIP:       types.StringValue(strings.Join(mdm.ManagementIPs, ",")),
			VirtualIPNIC: types.StringValue(strings.Join(mdm.VirtIPIntfsList, ",")),
			Role:         types.StringValue(string("Manager")),
			Mode:         types.StringValue(string("Secondary")),
		}

		if len(clusteDetailResponse.ClusterDetails.VirtualIPs) > 0 {
			mdmData.VirtualIP = types.StringValue(strings.Join(clusteDetailResponse.ClusterDetails.VirtualIPs, ","))
		}

		mdmList = append(mdmList, mdmData)
	}

	for _, mdm := range clusteDetailResponse.ClusterDetails.StandbyMdmSet {

		id, error := decimalToTwosComplementHex(mdm.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		mdmData := models.MDMModel{
			ID:           types.StringValue(id),
			Name:         types.StringValue(mdm.Name),
			IP:           types.StringValue(strings.Join(mdm.Node.NodeIPs, ",")),
			MDMIP:        types.StringValue(strings.Join(mdm.MdmIPs, ",")),
			MGMTIP:       types.StringValue(strings.Join(mdm.ManagementIPs, ",")),
			VirtualIPNIC: types.StringValue(strings.Join(mdm.VirtIPIntfsList, ",")),
			Role:         types.StringValue(string("Manager")),
			Mode:         types.StringValue(string("Standby")),
		}

		mdmList = append(mdmList, mdmData)
	}

	for _, mdm := range clusteDetailResponse.ClusterDetails.TbSet {

		id, error := decimalToTwosComplementHex(mdm.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		mdmData := models.MDMModel{
			ID:    types.StringValue(id),
			Name:  types.StringValue(mdm.Name),
			IP:    types.StringValue(strings.Join(mdm.Node.NodeIPs, ",")),
			MDMIP: types.StringValue(strings.Join(mdm.MdmIPs, ",")),
			Role:  types.StringValue(string("TieBreaker")),
			Mode:  types.StringValue(string("TieBreaker")),
		}

		mdmList = append(mdmList, mdmData)
	}

	for _, mdm := range clusteDetailResponse.ClusterDetails.StandbyTbSet {

		id, error := decimalToTwosComplementHex(mdm.ID)

		if error != nil {
			diags.AddError(error.Error(), error.Error())
			return plan, diags
		}

		mdmData := models.MDMModel{
			ID:    types.StringValue(id),
			Name:  types.StringValue(mdm.Name),
			IP:    types.StringValue(strings.Join(mdm.Node.NodeIPs, ",")),
			MDMIP: types.StringValue(strings.Join(mdm.MdmIPs, ",")),
			Role:  types.StringValue("TieBreaker"),
			Mode:  types.StringValue("Standby"),
		}

		mdmList = append(mdmList, mdmData)
	}

	mastMDM := clusteDetailResponse.ClusterDetails.MasterMdm

	id, error := decimalToTwosComplementHex(mastMDM.ID)

	if error != nil {
		diags.AddError(error.Error(), error.Error())
		return plan, diags
	}

	mdmData := models.MDMModel{
		ID:           types.StringValue(id),
		Name:         types.StringValue(mastMDM.Name),
		IP:           types.StringValue(getIP(mastMDM.ManagementIPs, mastMDM.Node.NodeIPs)),
		MGMTIP:       types.StringValue(strings.Join(mastMDM.ManagementIPs, ",")),
		MDMIP:        types.StringValue(strings.Join(mastMDM.MdmIPs, ",")),
		VirtualIPNIC: types.StringValue(strings.Join(mastMDM.VirtIPIntfsList, ",")),
		Role:         types.StringValue("Manager"),
		Mode:         types.StringValue("Primary"),
	}

	if len(clusteDetailResponse.ClusterDetails.VirtualIPs) > 0 {
		mdmData.VirtualIP = types.StringValue(strings.Join(clusteDetailResponse.ClusterDetails.VirtualIPs, ","))
	}

	mdmList = append(mdmList, mdmData)

	for _, mdm := range mdmList {
		objVal, dgs := GetMDMValue(mdm)
		diags = append(diags, dgs...)
		objectMDMs = append(objectMDMs, objVal)
	}
	setMDMs, dgs := types.SetValue(MDMElemType, objectMDMs)
	diags = append(diags, dgs...)

	state.MDMList = setMDMs

	return state, diags
}

func getIP(managementIPs []string, nodeIPs []string) string {
	if len(managementIPs) > 0 {
		return strings.Join(managementIPs, ",")
	}
	return strings.Join(nodeIPs, ",")
}

// GetPDType returns the Protection Domain Detail type
func GetPDType() map[string]attr.Type {
	return map[string]attr.Type{
		"name":              types.StringType,
		"storage_pool_list": types.ListType{ElemType: types.ObjectType{AttrTypes: GetStoragePoolsType()}},
	}
}

// GetPDValue returns the ProtectionDomain Detail model object value
func GetPDValue(pd models.ProtectionDomainDataModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetPDType(), map[string]attr.Value{
		"name":              types.StringValue(pd.Name.ValueString()),
		"storage_pool_list": pd.StoragePools,
	})
}

// GetStoragePoolsType returns the Storage Pool Detail type
func GetStoragePoolsType() map[string]attr.Type {
	return map[string]attr.Type{
		"media_type":            types.StringType,
		"external_acceleration": types.StringType,
		"name":                  types.StringType,
		"data_layout":           types.StringType,
		"compression_method":    types.StringType,
		"zero_padding":          types.StringType,
		"replication_journal_capacity_percentage": types.Int64Type,
	}
}

// GetStoragePoolsValue returns the Storage Pool Detail model object value
func GetStoragePoolsValue(sp models.StoragePoolDetailModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetStoragePoolsType(), map[string]attr.Value{
		"name":                  types.StringValue(sp.Name.ValueString()),
		"media_type":            types.StringValue(sp.MediaType.ValueString()),
		"external_acceleration": types.StringValue(sp.ExternalAcceleration.ValueString()),
		"data_layout":           types.StringValue(sp.DataLayout.ValueString()),
		"compression_method":    types.StringValue(sp.CompressionMethod.ValueString()),
		"zero_padding":          types.StringValue(sp.ZeroPadding.ValueString()),
		"replication_journal_capacity_percentage": types.Int64Value(sp.ReplicationJournalCapacityPercentage.ValueInt64()),
	})
}

// GetMDMType returns the MDM Detail type
func GetMDMType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"ip":             types.StringType,
		"name":           types.StringType,
		"mdm_ip":         types.StringType,
		"mgmt_ip":        types.StringType,
		"virtual_ip":     types.StringType,
		"virtual_ip_nic": types.StringType,
		"role":           types.StringType,
		"mode":           types.StringType,
	}
}

// GetMDMValue returns the MDM Detail model object value
func GetMDMValue(mdm models.MDMModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetMDMType(), map[string]attr.Value{
		"id":             types.StringValue(mdm.ID.ValueString()),
		"ip":             types.StringValue(mdm.IP.ValueString()),
		"name":           types.StringValue(mdm.Name.ValueString()),
		"mdm_ip":         types.StringValue(mdm.MDMIP.ValueString()),
		"mgmt_ip":        types.StringValue(mdm.MGMTIP.ValueString()),
		"virtual_ip":     types.StringValue(mdm.VirtualIP.ValueString()),
		"virtual_ip_nic": types.StringValue(mdm.VirtualIPNIC.ValueString()),
		"role":           types.StringValue(mdm.Role.ValueString()),
		"mode":           types.StringValue(mdm.Mode.ValueString()),
	})
}

// GetSDRType returns the SDR Detail type
func GetSDRType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"ip":              types.StringType,
		"name":            types.StringType,
		"port":            types.Int64Type,
		"application_ips": types.StringType,
		"storage_ips":     types.StringType,
		"external_ips":    types.StringType,
		"all_ips":         types.StringType,
	}
}

// GetSDRValue returns the SDR Detail model object value
func GetSDRValue(sdr models.SDRModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDRType(), map[string]attr.Value{
		"id":              types.StringValue(sdr.ID.ValueString()),
		"ip":              types.StringValue(sdr.IP.ValueString()),
		"name":            types.StringValue(sdr.Name.ValueString()),
		"port":            types.Int64Value(sdr.Port.ValueInt64()),
		"application_ips": types.StringValue(sdr.ApplicationIP.ValueString()),
		"storage_ips":     types.StringValue(sdr.StorageIP.ValueString()),
		"external_ips":    types.StringValue(sdr.ExternalIP.ValueString()),
		"all_ips":         types.StringValue(sdr.AllIP.ValueString()),
	})
}

// GetDeviceType returns the Device Detail type
func GetDeviceType() map[string]attr.Type {
	return map[string]attr.Type{
		"path":               types.StringType,
		"storage_pool":       types.StringType,
		"name":               types.StringType,
		"max_capacity_in_kb": types.Int64Type,
	}
}

// GetDeviceValue returns the Device Detail model object value
func GetDeviceValue(device models.DeviceDataModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetDeviceType(), map[string]attr.Value{
		"path":               types.StringValue(device.Path.ValueString()),
		"storage_pool":       types.StringValue(device.StoragePoolName.ValueString()),
		"name":               types.StringValue(device.Name.ValueString()),
		"max_capacity_in_kb": types.Int64Value(device.MaxCapacity.ValueInt64()),
	})
}

// GetSDSType returns the SDS Detail type
func GetSDSType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"ip":                     types.StringType,
		"name":                   types.StringType,
		"all_ips":                types.StringType,
		"sds_only_ips":           types.StringType,
		"sds_sdc_ips":            types.StringType,
		"protection_domain_id":   types.StringType,
		"protection_domain_name": types.StringType,
		"fault_set":              types.StringType,
		"devices":                types.SetType{ElemType: types.ObjectType{AttrTypes: GetDeviceType()}},
	}
}

// GetSDSValue returns the SDS Detail model object value
func GetSDSValue(sds models.SDSModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDSType(), map[string]attr.Value{
		"id":                     types.StringValue(sds.ID.ValueString()),
		"ip":                     types.StringValue(sds.IP.ValueString()),
		"name":                   types.StringValue(sds.Name.ValueString()),
		"all_ips":                types.StringValue(sds.AllIP.ValueString()),
		"sds_only_ips":           types.StringValue(sds.SDSOnlyIP.ValueString()),
		"sds_sdc_ips":            types.StringValue(sds.SDSSDCIP.ValueString()),
		"protection_domain_id":   types.StringValue(sds.ProtectionDomainID.ValueString()),
		"protection_domain_name": types.StringValue(sds.ProtectionDomainName.ValueString()),
		"fault_set":              types.StringValue(sds.FaultSet.ValueString()),
		"devices":                sds.Devices,
	})
}

// GetSDTType returns the SDT Detail type
func GetSDTType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"ip":                     types.StringType,
		"name":                   types.StringType,
		"all_ips":                types.StringType,
		"protection_domain_id":   types.StringType,
		"protection_domain_name": types.StringType,
		"storage_only_ips":       types.StringType,
		"host_only_ips":          types.StringType,
		"storage_port":           types.Int64Type,
		"nvme_port":              types.Int64Type,
		"discovery_port":         types.Int64Type,
	}
}

// GetSDTValue returns the SDT Detail model object value
func GetSDTValue(sdt models.SDTModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDTType(), map[string]attr.Value{
		"id":                     types.StringValue(sdt.ID.ValueString()),
		"ip":                     types.StringValue(sdt.IP.ValueString()),
		"name":                   types.StringValue(sdt.Name.ValueString()),
		"all_ips":                types.StringValue(sdt.AllIP.ValueString()),
		"protection_domain_id":   types.StringValue(sdt.ProtectionDomainID.ValueString()),
		"protection_domain_name": types.StringValue(sdt.ProtectionDomainName.ValueString()),
		"storage_only_ips":       types.StringValue(sdt.StorageOnlyIPs.ValueString()),
		"host_only_ips":          types.StringValue(sdt.HostOnlyIPs.ValueString()),
		"storage_port":           types.Int64Value(sdt.StoragePort.ValueInt64()),
		"nvme_port":              types.Int64Value(sdt.NvmePort.ValueInt64()),
		"discovery_port":         types.Int64Value(sdt.DiscoveryPort.ValueInt64()),
	})
}

// GetSDCType returns the SDC Detail type
func GetSDCType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"ip":   types.StringType,
		"name": types.StringType,
		"guid": types.StringType,
	}
}

// GetSDCValue returns the SDC Detail model object value
func GetSDCValue(sds models.SDCModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetSDCType(), map[string]attr.Value{
		"id":   types.StringValue(sds.ID.ValueString()),
		"ip":   types.StringValue(sds.IP.ValueString()),
		"name": types.StringValue(sds.Name.ValueString()),
		"guid": types.StringValue(sds.GUID.ValueString()),
	})
}

// GetMDMIPFromClusterDetails function is used for fetch MDM IP from Cluster Details
func GetMDMIPFromClusterDetails(clusterInstallationDetailsDataModel []models.ClusterModel) (string, error) {
	var mdmIP string

	for _, item := range clusterInstallationDetailsDataModel {
		if strings.EqualFold(item.IsMdmOrTb.ValueString(), "Primary") {
			if item.MDMMgmtIP.ValueString() != "" {
				mdmIP = item.MDMMgmtIP.ValueString()
			} else if item.MDMIP.ValueString() != "" {
				mdmIP = item.MDMIP.ValueString()
			} else {
				mdmIP = item.IP.ValueString()
			}
			return mdmIP, nil
		}
	}
	return mdmIP, nil
}

// GetMDMIPFromMDMList function is used for fetch MDM IP from MDM List Details
func GetMDMIPFromMDMList(mdmDataModel []models.MDMModel) (string, error) {
	var mdmIP string

	for _, item := range mdmDataModel {
		if strings.EqualFold(item.Mode.ValueString(), "Primary") {
			mdmIP = item.IP.ValueString()
			return mdmIP, nil
		}
	}
	return mdmIP, nil
}

// ParseClusterCSVOperation function for Handling Parsing CSV Operation
func ParseClusterCSVOperation(ctx context.Context, gatewayClient *goscaleio.GatewayClient, clusterDataModel []models.ClusterModel, storagePoolDataModel []models.StoragePoolDataModel) (*goscaleio_types.GatewayResponse, error) {

	var parseCSVResponse goscaleio_types.GatewayResponse

	//Create a csv file from the input given by the user
	mydir, err := os.Getwd()
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Reading Current Directory is %s", err.Error())
	}
	// Create a csv writer
	filePath := filepath.Join(mydir, filepath.Clean("Minimal.csv"))
	file, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV is %s", err.Error())
	}
	defer file.Close()
	writer := NewCustomCSVWriter(file)

	// Write the header row
	header := []string{"IPs", "Username", "Password", "Operating System", "Is MDM/TB", "MDM Mgmt IP", "MDM IPs", "MDM Name", "perfProfileForMDM", "Virtual IPs", "Virtual IP NICs", "Is SDS", "SDS Name", "SDS All IPs", "SDS-SDS Only IPs", "SDS-SDC Only IPs", "Protection Domain", "Fault Set", "SDS Storage Device List", "StoragePool List", "SDS Storage Device Names", "perfProfileForSDS", "Is SDC", "perfProfileForSDC", "SDC Name", "RFcache", "RFcache SSD Device List", "Is SDR", "SDR Name", "SDR Port", "SDR Application IPs", "SDR Storage IPs", "SDR External IPs", "SDR All IPs", "perfProfileForSDR", "Is SDT", "SDT Name", "SDT All IPs"}
	var headerIndicesToWrite []int

	for i := range header {
		headerIndicesToWrite = append(headerIndicesToWrite, i)
	}

	// Create a slice to hold the filtered header with non-empty values
	var filteredHeader []string

	for _, item := range clusterDataModel {

		data := []string{
			item.IP.ValueString(),
			item.UserName.ValueString(),
			item.Password.ValueString(),
			item.OperatingSystem.ValueString(),
			item.IsMdmOrTb.ValueString(),
			item.MDMMgmtIP.ValueString(),
			item.MDMIP.ValueString(),
			item.MDMName.ValueString(),
			item.PerfProfileForMDM.ValueString(),
			item.VirtualIPs.ValueString(),
			item.VirtualIPNICs.ValueString(),
			item.IsSds.ValueString(),
			item.SDSName.ValueString(),
			item.SDSAllIPs.ValueString(),
			item.SDSToSDSOnlyIPs.ValueString(),
			item.SDSToSDCOnlyIPs.ValueString(),
			item.ProtectionDomain.ValueString(),
			item.FaultSet.ValueString(),
			item.SDSStorageDeviceList.ValueString(),
			item.StoragePoolList.ValueString(),
			item.SDSStorageDeviceNames.ValueString(),
			item.PerfProfileForSDS.ValueString(),
			item.IsSdc.ValueString(),
			item.PerfProfileForSDC.ValueString(),
			item.SDCName.ValueString(),
			item.IsRFCache.ValueString(),
			item.RFcacheSSDDeviceList.ValueString(),
			item.IsSdr.ValueString(),
			item.SDRName.ValueString(),
			item.SDRPort.ValueString(),
			item.SDRApplicationIPs.ValueString(),
			item.SDRStorageIPs.ValueString(),
			item.SDRExternalIPs.ValueString(),
			item.SDRAllIPS.ValueString(),
			item.PerfProfileForSDR.ValueString(),
			item.IsSdt.ValueString(),
			item.SDTName.ValueString(),
			item.SDTAllIPs.ValueString(),
		}

		// Check which columns have non-empty values in the current row
		var columnsWithValues []int
		for i, value := range data {

			if i == 2 {
				//we have to add the Password column no matter if it's value is empty or not.
				columnsWithValues = append(columnsWithValues, i)
				filteredHeader = append(filteredHeader, header[i])
			} else if i != 2 && value != "" {
				columnsWithValues = append(columnsWithValues, i)
				// Add the corresponding header to the filteredHeader
				filteredHeader = append(filteredHeader, header[i])
			}

		}

		// Update the header indices to write based on the current row's non-empty columns
		headerIndicesToWrite = intersect(headerIndicesToWrite, columnsWithValues)
	}
	// Remove duplicates from the filteredHeader
	filteredHeader = removeDuplicates(filteredHeader)

	err = writer.Write(filteredHeader)
	if err != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
	}

	// Write the values for each data row according to the filtered headers
	for _, item := range clusterDataModel {
		var data []string
		for _, header := range filteredHeader {
			data = append(data, getFieldFromItem(item, header))
		}

		err = writer.Write(data)
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
		}
	}

	//checking if storage_pool data is available than only write data for that.
	if len(storagePoolDataModel) > 0 {
		// Add a new line with commas
		err = writer.Write([]string{strings.Repeat(",", len(filteredHeader))})
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
		}

		// Add a blank line after writing each data row
		err = writer.Write([]string{"Storage Pool Configuration", strings.Repeat(",", len(filteredHeader)-1)})
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
		}

		// Write the Storage Pool header row
		headerForStorage := []string{"ProtectionDomain", "StoragePool", "Media Type", "External Acceleration", "Data Layout", "Zero Padding", "Compression Method", "Replication journal capacity percentage"}
		var storageHeaderIndicesToWrite []int

		for i := range headerForStorage {
			storageHeaderIndicesToWrite = append(storageHeaderIndicesToWrite, i)
		}

		// Create a slice to hold the filtered header with non-empty values
		var filteredStorageHeader []string

		for _, item := range storagePoolDataModel {

			data := []string{
				item.ProtectionDomain.ValueString(),
				item.StoragePool.ValueString(),
				item.MediaType.ValueString(),
				item.ExternalAcceleration.ValueString(),
				item.DataLayout.ValueString(),
				item.ZeroPadding.ValueString(),
				item.CompressionMethod.ValueString(),
				item.ReplicationJournalCapacityPercentage.ValueString(),
			}

			// Check which columns have non-empty values in the current row
			var columnsWithValues []int
			for i, value := range data {
				if value != "" {
					columnsWithValues = append(columnsWithValues, i)
					// Add the corresponding header to the filteredStorageHeader
					filteredStorageHeader = append(filteredStorageHeader, headerForStorage[i])
				}
			}

			// Update the header indices to write based on the current row's non-empty columns
			storageHeaderIndicesToWrite = intersect(storageHeaderIndicesToWrite, columnsWithValues)
		}
		// Remove duplicates from the filteredStorageHeader
		filteredStorageHeader = removeDuplicates(filteredStorageHeader)

		err = writer.Write(append(filteredStorageHeader, strings.Repeat(",", len(filteredHeader)-len(filteredStorageHeader))))
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Writing Temp CSV is %s", err.Error())
		}

		// Write the values for each data row according to the filtered headers
		for _, item := range storagePoolDataModel {
			var data []string
			for _, header := range filteredStorageHeader {
				data = append(data, getFieldFromStorage(item, header))
			}

			data = append(data, strings.Repeat(",", len(filteredHeader)-len(filteredStorageHeader)))

			err = writer.Write(data)
			if err != nil {
				return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
			}
		}

		err = writer.Flush()
		if err != nil {
			return &parseCSVResponse, fmt.Errorf("Error While Creating Temp CSV File is %s", err.Error())
		}
	}

	parsecsvRespose, parseCSVError := gatewayClient.ParseCSV(mydir + "/Minimal.csv")

	deletCSVError := os.Remove(mydir + "/Minimal.csv")
	if deletCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("Error While Deleting Temp CSV File is %s", deletCSVError.Error())
	}

	if parseCSVError != nil {
		return &parseCSVResponse, fmt.Errorf("%s", parseCSVError.Error())
	}

	if parsecsvRespose.StatusCode != 200 {
		return &parseCSVResponse, fmt.Errorf("Meesage : %s, Error Cosde : %s", parsecsvRespose.Message, strconv.Itoa(parsecsvRespose.StatusCode))
	}

	return parsecsvRespose, nil
}

// getFieldFromStorage Function to get the field value from the StoragePoolDataModel
func getFieldFromStorage(item models.StoragePoolDataModel, header string) string {
	switch header {
	case "ProtectionDomain":
		return item.ProtectionDomain.ValueString()
	case "StoragePool":
		return item.StoragePool.ValueString()
	case "Media Type":
		return item.MediaType.ValueString()
	case "External Acceleration":
		return item.ExternalAcceleration.ValueString()
	case "Data Layout":
		return item.DataLayout.ValueString()
	case "Zero Padding":
		return item.ZeroPadding.ValueString()
	case "Compression Method":
		return item.CompressionMethod.ValueString()
	case "Replication journal capacity percentage":
		return item.ReplicationJournalCapacityPercentage.ValueString()
	default:
		return "" // Return empty string for unknown headers
	}
}

// getFieldFromItem Function to get the field value from the ClusterDataModel
func getFieldFromItem(item models.ClusterModel, header string) string {
	switch header {
	case "IPs":
		return "\"" + item.IP.ValueString() + "\""
	case "Username":
		return item.UserName.ValueString()
	case "Password":
		return item.Password.ValueString()
	case "Operating System":
		return item.OperatingSystem.ValueString()
	case "Is MDM/TB":
		return "\"" + item.IsMdmOrTb.ValueString() + "\""
	case "MDM Mgmt IP":
		return "\"" + item.MDMMgmtIP.ValueString() + "\""
	case "MDM IPs":
		return "\"" + item.MDMIP.ValueString() + "\""
	case "MDM Name":
		return item.MDMName.ValueString()
	case "perfProfileForMDM":
		return item.PerfProfileForMDM.ValueString()
	case "Virtual IPs":
		return "\"" + item.VirtualIPs.ValueString() + "\""
	case "Virtual IP NICs":
		return "\"" + item.VirtualIPNICs.ValueString() + "\""
	case "Is SDS":
		return item.IsSds.ValueString()
	case "SDS Name":
		return item.SDSName.ValueString()
	case "SDS All IPs":
		return "\"" + item.SDSAllIPs.ValueString() + "\""
	case "SDS-SDS Only IPs":
		return "\"" + item.SDSToSDSOnlyIPs.ValueString() + "\""
	case "SDS-SDC Only IPs":
		return "\"" + item.SDSToSDCOnlyIPs.ValueString() + "\""
	case "Protection Domain":
		return item.ProtectionDomain.ValueString()
	case "Fault Set":
		return item.FaultSet.ValueString()
	case "SDS Storage Device List":
		return "\"" + item.SDSStorageDeviceList.ValueString() + "\""
	case "StoragePool List":
		return "\"" + item.StoragePoolList.ValueString() + "\""
	case "SDS Storage Device Names":
		return "\"" + item.SDSStorageDeviceNames.ValueString() + "\""
	case "perfProfileForSDS":
		return item.PerfProfileForSDS.ValueString()
	case "Is SDC":
		return item.IsSdc.ValueString()
	case "perfProfileForSDC":
		return item.PerfProfileForSDC.ValueString()
	case "SDC Name":
		return item.SDCName.ValueString()
	case "RFcache":
		return item.IsRFCache.ValueString()
	case "RFcache SSD Device List":
		return "\"" + item.RFcacheSSDDeviceList.ValueString() + "\""
	case "Is SDR":
		return item.IsSdr.ValueString()
	case "SDR Name":
		return item.SDRName.ValueString()
	case "SDR Port":
		return "\"" + item.SDRPort.ValueString() + "\""
	case "SDR Application IPs":
		return "\"" + item.SDRApplicationIPs.ValueString() + "\""
	case "SDR Storage IPs":
		return "\"" + item.SDRStorageIPs.ValueString() + "\""
	case "SDR External IPs":
		return "\"" + item.SDRExternalIPs.ValueString() + "\""
	case "SDR All IPs":
		return "\"" + item.SDRAllIPS.ValueString() + "\""
	case "perfProfileForSDR":
		return item.PerfProfileForSDR.ValueString()
	case "Is SDT":
		return item.IsSdt.ValueString()
	case "SDT Name":
		return item.SDTName.ValueString()
	case "SDT All IPs":
		return "\"" + item.SDTAllIPs.ValueString() + "\""
	default:
		return "" // Return empty string for unknown headers
	}
}

// GetClusterDetails function for Validate the MDM credentials
func GetClusterDetails(model models.ClusterResourceModel, gatewayClient *goscaleio.GatewayClient, mdmIP string, requireJSONOutput bool) (*goscaleio_types.GatewayResponse, error) {
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": model.MdmPassword.ValueString(),
	}
	mapData["mdmIps"] = strings.Split(mdmIP, ",")

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": model.AllowNonSecureCommunicationWithMdm.ValueBool(),
		"allowNonSecureCommunicationWithLia": model.AllowNonSecureCommunicationWithLia.ValueBool(),
		"disableNonMgmtComponentsAuth":       model.DisableNonMgmtComponentsAuth.ValueBool(),
	}
	mapData["securityConfiguration"] = secureData
	jsonreq, _ := json.Marshal(mapData)

	validateMDMResponse, validateMDMError := gatewayClient.GetClusterDetails(jsonreq, requireJSONOutput)
	if validateMDMError != nil {
		return validateMDMResponse, fmt.Errorf("%s", validateMDMError.Error())
	} else if validateMDMResponse.StatusCode >= 300 {
		return validateMDMResponse, fmt.Errorf("%s, Please Validate Entered Details", validateMDMResponse.Message)
	}

	return validateMDMResponse, nil
}

// ClusterInstallationOperations function for begin instllation process
func ClusterInstallationOperations(ctx context.Context, model models.ClusterResourceModel, gatewayClient *goscaleio.GatewayClient, parsecsvRespose *goscaleio_types.GatewayResponse) error {

	beginInstallationResponse, installationError := gatewayClient.BeginInstallation(parsecsvRespose.Data, "admin", model.MdmPassword.ValueString(), model.LiaPassword.ValueString(), model.AllowNonSecureCommunicationWithMdm.ValueBool(), model.AllowNonSecureCommunicationWithLia.ValueBool(), model.DisableNonMgmtComponentsAuth.ValueBool(), false)

	if installationError != nil {
		return fmt.Errorf("Error while begin installation is %s", installationError.Error())
	}

	if beginInstallationResponse.StatusCode == 200 {
		currentPhase := "query"
		couterForStopExecution := 0

		tflog.Info(ctx, "Gateway Installation Begin, Current Phase - Query")

		for couterForStopExecution <= 5 {

			time.Sleep(1 * time.Minute)

			checkForPhaseCompleted, _ := gatewayClient.CheckForCompletionQueueCommands(currentPhase)

			if checkForPhaseCompleted.Data == "Completed" {
				couterForStopExecution = 0

				if currentPhase != "configure" {
					moveToNextPhaseResponse, err := gatewayClient.MoveToNextPhase()

					if err != nil {
						return fmt.Errorf("Error while moving to next phase is %s", err.Error())
					}

					if moveToNextPhaseResponse.StatusCode == 200 {
						if currentPhase == "query" {
							currentPhase = "upload"
							tflog.Info(ctx, "Gateway Installation phase changed to Upload")
						} else if currentPhase == "upload" {
							currentPhase = "install"
							tflog.Info(ctx, "Gateway Installation phase changed to Install")
						} else if currentPhase == "install" {
							currentPhase = "configure"
							tflog.Info(ctx, "Gateway Installation phase changed to Configure")
						}
					} else {
						return fmt.Errorf("Messsage: %s, Error Code: %s", moveToNextPhaseResponse.Message, strconv.Itoa(moveToNextPhaseResponse.StatusCode))
					}
				} else {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Installation is %s", queueOperationError.Error())
					}

					return nil
				}

			} else if checkForPhaseCompleted.Data == "Running" {
				couterForStopExecution++

				tflog.Info(ctx, "Gateway Installation operations are still running in phase "+currentPhase)

				if couterForStopExecution == 6 {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Installation in phase %s is %s", currentPhase, queueOperationError.Error())
					}

					return fmt.Errorf("Time Out,Some Operations of Installer running from since long")
				}

			} else {
				return fmt.Errorf("Error During Installation is %s", checkForPhaseCompleted.Message)
			}
		}
	} else {
		return fmt.Errorf("Message: %s, Error Code: %s", beginInstallationResponse.Message, strconv.Itoa(beginInstallationResponse.StatusCode))
	}

	return nil
}

// ClusterUninstallationOperations function for begin uninstllation process
func ClusterUninstallationOperations(ctx context.Context, model models.ClusterResourceModel, gatewayClient *goscaleio.GatewayClient, parsecsvRespose *goscaleio_types.GatewayResponse) error {

	clusterMapData, jsonParseError := jsonToMap(parsecsvRespose.Data)
	if jsonParseError != nil {
		return fmt.Errorf("Error while begin uninstallation is %s", jsonParseError.Error())
	}

	sdcResData := clusterMapData["sdcList"].([]interface{})

	if len(sdcResData) > 0 {
		var sdcFinalData []interface{}

		for _, sdcNode := range sdcResData {
			sdcNode := sdcNode.(map[string]interface{})
			node := sdcNode["node"].(map[string]interface{})
			nodeIPs := node["nodeIPs"].([]interface{})
			if nodeIPs[0].(string) != "N/A" {
				sdcFinalData = append(sdcFinalData, sdcNode)
			}
		}

		clusterMapData["sdcList"] = sdcFinalData
	}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": model.AllowNonSecureCommunicationWithMdm,
		"allowNonSecureCommunicationWithLia": model.AllowNonSecureCommunicationWithLia,
		"disableNonMgmtComponentsAuth":       model.DisableNonMgmtComponentsAuth,
	}

	clusterMapData["securityConfiguration"] = secureData

	clusterJSONData, jsonParseError := json.Marshal(clusterMapData)
	if jsonParseError != nil {
		return fmt.Errorf("Error while begin uninstallation is %s", jsonParseError.Error())
	}

	beginUninstallationResponse, uninstallationError := gatewayClient.UninstallCluster(string(clusterJSONData), "admin", model.MdmPassword.ValueString(), model.LiaPassword.ValueString(), model.AllowNonSecureCommunicationWithMdm.ValueBool(), model.AllowNonSecureCommunicationWithLia.ValueBool(), model.DisableNonMgmtComponentsAuth.ValueBool(), false)

	if uninstallationError != nil {
		return fmt.Errorf("Error while begin uninstallation is %s", uninstallationError.Error())
	}

	if beginUninstallationResponse.StatusCode == 200 {
		currentPhase := "query"
		couterForStopExecution := 0

		tflog.Info(ctx, "Gateway Uninstallation Begin, Current Phase - Query")

		for couterForStopExecution <= 5 {

			time.Sleep(1 * time.Minute)

			checkForPhaseCompleted, _ := gatewayClient.CheckForCompletionQueueCommands(currentPhase)

			if checkForPhaseCompleted.Data == "Completed" {
				couterForStopExecution = 0

				if currentPhase != "clean" {
					moveToNextPhaseResponse, err := gatewayClient.MoveToNextPhase()

					if err != nil {
						return fmt.Errorf("Error while moving to next phase is %s", err.Error())
					}

					if moveToNextPhaseResponse.StatusCode == 200 {
						if currentPhase == "query" {
							currentPhase = "clean"
							tflog.Info(ctx, "Gateway uninstallation phase changed to Clean")
						}
					} else {
						return fmt.Errorf("Messsage: %s, Error Code: %s", moveToNextPhaseResponse.Message, strconv.Itoa(moveToNextPhaseResponse.StatusCode))
					}
				} else {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue After uninstallation is %s", queueOperationError.Error())
					}

					return nil
				}

			} else if checkForPhaseCompleted.Data == "Running" {
				couterForStopExecution++

				tflog.Info(ctx, "Gateway Uninstallation operations are still running")

				if couterForStopExecution == 5 {
					// to make gateway available for installation
					queueOperationError := ResetInstallerQueue(gatewayClient)
					if queueOperationError != nil {
						return fmt.Errorf("Error Clearing Queue During Uninstallation is %s", queueOperationError.Error())
					}

					return fmt.Errorf("Time Out,Some Operations of uninstall running from since long")
				}

			} else {
				return fmt.Errorf("Error During Uninstallation is %s", checkForPhaseCompleted.Message)
			}
		}
	} else {
		return fmt.Errorf("Message: %s, Error Code: %s", beginUninstallationResponse.Message, strconv.Itoa(beginUninstallationResponse.StatusCode))
	}

	return nil
}

// removeDuplicates Helper function to remove duplicates from a slice of strings
func removeDuplicates(s []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, item := range s {
		if !unique[item] {
			unique[item] = true
			result = append(result, item)
		}
	}
	return result
}

// intersect Helper function to intersect slices
func intersect(a, b []int) []int {
	m := make(map[int]bool)
	var intersection []int

	for _, val := range a {
		m[val] = true
	}

	for _, val := range b {
		if m[val] {
			intersection = append(intersection, val)
		}
	}

	return intersection
}

// decimalToTwosComplementHex function to convert decimal id to hex value
func decimalToTwosComplementHex(decimalStr string) (string, error) {
	decimalInt, err := strconv.ParseInt(decimalStr, 10, 64)
	if err != nil {
		return "", err
	}

	// Calculate the two's complement value
	twosComplement := uint64(decimalInt) & ((1 << 64) - 1)

	// Convert the value to a hexadecimal string
	twosComplementHex := fmt.Sprintf("%016X", twosComplement)

	return twosComplementHex, nil
}

// boolToYesNo function to convert bool value to Yes/No
func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// jsonToMap used for convert json to map
func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
