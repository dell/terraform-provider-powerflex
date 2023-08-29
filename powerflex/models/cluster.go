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

// ClusterResourceModel struct for CSV Data Processing
type ClusterResourceModel struct {
	//Input Fields
	ID                                 types.String `tfsdk:"id"`
	Cluster                            types.List   `tfsdk:"cluster"`
	StoragePools                       types.List   `tfsdk:"storage_pools"`
	MdmPassword                        types.String `tfsdk:"mdm_password"`
	LiaPassword                        types.String `tfsdk:"lia_password"`
	AllowNonSecureCommunicationWithMdm types.Bool   `tfsdk:"allow_non_secure_communication_with_mdm"`
	AllowNonSecureCommunicationWithLia types.Bool   `tfsdk:"allow_non_secure_communication_with_lia"`
	DisableNonMgmtComponentsAuth       types.Bool   `tfsdk:"disable_non_mgmt_components_auth"`
	MDMList                            types.Set    `tfsdk:"mdm_list"`
	SDSList                            types.Set    `tfsdk:"sds_list"`
	SDCList                            types.Set    `tfsdk:"sdc_list"`
	SDRList                            types.Set    `tfsdk:"sdr_list"`
	ProtectionDomains                  types.List   `tfsdk:"protection_domains"`
}

// ClusterModel defines the struct for Cluster Details Data
type ClusterModel struct {
	IP                    types.String `tfsdk:"ips"`
	UserName              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	OperatingSystem       types.String `tfsdk:"operating_system"`
	IsMdmOrTb             types.String `tfsdk:"is_mdm_or_tb"`
	MDMIP                 types.String `tfsdk:"mdm_ips"`
	MDMMgmtIP             types.String `tfsdk:"mdm_mgmt_ip"`
	MDMName               types.String `tfsdk:"mdm_name"`
	PerfProfileForMDM     types.String `tfsdk:"perf_profile_for_mdm"`
	VirtualIPs            types.String `tfsdk:"virtual_ips"`
	VirtualIPNICs         types.String `tfsdk:"virtual_ip_nics"`
	IsSds                 types.String `tfsdk:"is_sds"`
	SDSName               types.String `tfsdk:"sds_name"`
	SDSAllIPs             types.String `tfsdk:"sds_all_ips"`
	SDSToSDSOnlyIPs       types.String `tfsdk:"sds_to_sds_only_ips"`
	SDSToSDCOnlyIPs       types.String `tfsdk:"sds_to_sdc_only_ips"`
	ProtectionDomain      types.String `tfsdk:"protection_domain"`
	FaultSet              types.String `tfsdk:"fault_set"`
	SDSStorageDeviceList  types.String `tfsdk:"sds_storage_device_list"`
	StoragePoolList       types.String `tfsdk:"storage_pool_list"`
	SDSStorageDeviceNames types.String `tfsdk:"sds_storage_device_names"`
	PerfProfileForSDS     types.String `tfsdk:"perf_profile_for_sds"`
	IsSdc                 types.String `tfsdk:"is_sdc"`
	PerfProfileForSDC     types.String `tfsdk:"perf_profile_for_sdc"`
	SDCName               types.String `tfsdk:"sdc_name"`
	IsRFCache             types.String `tfsdk:"is_rfcache"`
	RFcacheSSDDeviceList  types.String `tfsdk:"rf_cache_ssd_device_list"`
	IsSdr                 types.String `tfsdk:"is_sdr"`
	SDRName               types.String `tfsdk:"sdr_name"`
	SDRPort               types.String `tfsdk:"sdr_port"`
	SDRApplicationIPs     types.String `tfsdk:"sdr_application_ips"`
	SDRStorageIPs         types.String `tfsdk:"sdr_storage_ips"`
	SDRExternalIPs        types.String `tfsdk:"sdr_external_ips"`
	SDRAllIPS             types.String `tfsdk:"sdr_all_ips"`
	PerfProfileForSDR     types.String `tfsdk:"perf_profile_for_sdr"`
}

// StoragePoolDataModel desfines the srtuct for the Storage Pool Data
type StoragePoolDataModel struct {
	ProtectionDomain                     types.String `tfsdk:"protection_domain"`
	StoragePool                          types.String `tfsdk:"storage_pool"`
	MediaType                            types.String `tfsdk:"media_type"`
	ExternalAcceleration                 types.String `tfsdk:"external_acceleration"`
	DataLayout                           types.String `tfsdk:"data_layout"`
	ZeroPadding                          types.String `tfsdk:"zero_padding"`
	CompressionMethod                    types.String `tfsdk:"compression_method"`
	ReplicationJournalCapacityPercentage types.String `tfsdk:"replication_journal_capacity_percentage"`
}

// MDMModel desfines the srtuct for the MDM Data
type MDMModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	IP           types.String `tfsdk:"ip"`
	MDMIP        types.String `tfsdk:"mdm_ip"`
	MGMTIP       types.String `tfsdk:"mgmt_ip"`
	VirtualIP    types.String `tfsdk:"virtual_ip"`
	VirtualIPNIC types.String `tfsdk:"virtual_ip_nic"`
	Role         types.String `tfsdk:"role"`
	Mode         types.String `tfsdk:"mode"`
}

// SDSModel desfines the srtuct for the SDS Data
type SDSModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	IP                   types.String `tfsdk:"ip"`
	AllIP                types.String `tfsdk:"all_ips"`
	SDSOnlyIP            types.String `tfsdk:"sds_only_ips"`
	SDSSDCIP             types.String `tfsdk:"sds_sdc_ips"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	FaultSet             types.String `tfsdk:"fault_set"`
	Devices              types.Set    `tfsdk:"devices"`
}

// DeviceDataModel desfines the srtuct for the Device Data
type DeviceDataModel struct {
	Name            types.String `tfsdk:"name"`
	Path            types.String `tfsdk:"path"`
	StoragePoolName types.String `tfsdk:"storage_pool"`
	MaxCapacity     types.Int64  `tfsdk:"max_capacity_in_kb"`
}

// SDCModel desfines the srtuct for the SDC Data
type SDCModel struct {
	ID   types.String `tfsdk:"id"`
	GUID types.String `tfsdk:"guid"`
	Name types.String `tfsdk:"name"`
	IP   types.String `tfsdk:"ip"`
}

// SDRModel desfines the srtuct for the SDR Data
type SDRModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	IP            types.String `tfsdk:"ip"`
	Port          types.Int64  `tfsdk:"port"`
	ApplicationIP types.String `tfsdk:"application_ips"`
	StorageIP     types.String `tfsdk:"storage_ips"`
	ExternalIP    types.String `tfsdk:"external_ips"`
	AllIP         types.String `tfsdk:"all_ips"`
}

// ProtectionDomainDataModel desfines the srtuct for the Protection Domain Data
type ProtectionDomainDataModel struct {
	Name         types.String `tfsdk:"name"`
	StoragePools types.List   `tfsdk:"storage_pool_list"`
}

// StoragePoolDetailModel desfines the srtuct for the StoragePool Data
type StoragePoolDetailModel struct {
	Name                                 types.String `tfsdk:"name"`
	MediaType                            types.String `tfsdk:"media_type"`
	ExternalAcceleration                 types.String `tfsdk:"extern_alacceleration"`
	DataLayout                           types.String `tfsdk:"data_layout"`
	ZeroPadding                          types.String `tfsdk:"zero_padding"`
	CompressionMethod                    types.String `tfsdk:"compression_method"`
	ReplicationJournalCapacityPercentage types.Int64  `tfsdk:"replication_journal_capacity_percentage"`
}

// ClusterCsvRow desfines the srtuct for the Cluster CSV Data
type ClusterCsvRow struct {
	IP                    string
	UserName              string
	Password              string
	OperatingSystem       string
	IsMdmOrTb             string
	MDMIP                 string
	MDMMgmtIP             string
	MDMName               string
	PerfProfileForMDM     string
	VirtualIPs            string
	VirtualIPNICs         string
	IsSds                 string
	SDSName               string
	SDSAllIPs             string
	SDSToSDSOnlyIPs       string
	SDSToSDCOnlyIPs       string
	ProtectionDomain      string
	FaultSet              string
	SDSStorageDeviceList  string
	StoragePoolList       string
	SDSStorageDeviceNames string
	PerfProfileForSDS     string
	IsSdc                 string
	PerfProfileForSDC     string
	SDCName               string
	IsRFCache             string
	RFcacheSSDDeviceList  string
	IsSdr                 string
	SDRName               string
	SDRPort               string
	SDRApplicationIPs     string
	SDRStorageIPs         string
	SDRExternalIPs        string
	SDRAllIPS             string
	PerfProfileForSDR     string
}

// StoragePoolCsvRow desfines the srtuct for the Storage Pool CSV Data
type StoragePoolCsvRow struct {
	ProtectionDomain                     string
	StoragePool                          string
	MediaType                            string
	ExternalAcceleration                 string
	DataLayout                           string
	ZeroPadding                          string
	CompressionMethod                    string
	ReplicationJournalCapacityPercentage int
}
