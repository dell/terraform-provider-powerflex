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

package models

import (
	"context"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SdsResourceModel maps the resource schema data.
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

// SdsIPModel IP object
type SdsIPModel struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

// GetIPList converts list of IPs from tf model to go type
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

// IPList defines struct for SDS IP
type IPList struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

// SdsWindowType defines struct for SDS windows type
type SdsWindowType struct {
	Threshold            types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec      types.Int64 `tfsdk:"window_size_in_sec"`
	LastOscillationCount types.Int64 `tfsdk:"last_oscillation_count"`
	LastOscillationTime  types.Int64 `tfsdk:"last_oscillation_time"`
	MaxFailuresCount     types.Int64 `tfsdk:"max_failures_count"`
}

// SdsWindow defines struct for SDS window
type SdsWindow struct {
	ShortWindow  SdsWindowType `tfsdk:"short_window"`
	MediumWindow SdsWindowType `tfsdk:"medium_window"`
	LongWindow   SdsWindowType `tfsdk:"long_window"`
}

// CertificateInfoModel defines struct for certificate information
type CertificateInfoModel struct {
	Subject             types.String `tfsdk:"subject"`
	Issuer              types.String `tfsdk:"issuer"`
	ValidFrom           types.String `tfsdk:"valid_from" json:"validFrom"`
	ValidTo             types.String `tfsdk:"valid_to" json:"validTo"`
	Thumbprint          types.String `tfsdk:"thumbprint"`
	ValidFromAsn1Format types.String `tfsdk:"valid_from_asn1_format" json:"validFromAsn1Format"`
	ValidToAsn1Format   types.String `tfsdk:"valid_to_asn1_format" json:"validToAsn1Format"`
}

// RaidControllersModel defines struct for RAID controller
type RaidControllersModel struct {
	SerialNumber    types.String `tfsdk:"serial_number"`
	ModelName       types.String `tfsdk:"model_name"`
	VendorName      types.String `tfsdk:"vendor_name"`
	FirmwareVersion types.String `tfsdk:"firmware_version"`
	DriverVersion   types.String `tfsdk:"driver_version"`
	DriverName      types.String `tfsdk:"driver_name"`
	PciAddress      types.String `tfsdk:"pci_address"`
	Status          types.String `tfsdk:"status"`
	BatteryStatus   types.String `tfsdk:"battery_status"`
}

// SdsDataModel defines struct for SDS data model
type SdsDataModel struct {
	ID                                          types.String           `tfsdk:"id"`
	Name                                        types.String           `tfsdk:"name"`
	IPList                                      []IPList               `tfsdk:"ip_list"`
	Port                                        types.Int64            `tfsdk:"port"`
	SdsState                                    types.String           `tfsdk:"sds_state"`
	MembershipState                             types.String           `tfsdk:"membership_state"`
	MdmConnectionState                          types.String           `tfsdk:"mdm_connection_state"`
	DrlMode                                     types.String           `tfsdk:"drl_mode"`
	RmcacheEnabled                              types.Bool             `tfsdk:"rmcache_enabled"`
	RmcacheSize                                 types.Int64            `tfsdk:"rmcache_size"`
	RmcacheFrozen                               types.Bool             `tfsdk:"rmcache_frozen"`
	OnVmware                                    types.Bool             `tfsdk:"on_vmware"`
	FaultsetID                                  types.String           `tfsdk:"fault_set_id"`
	NumIOBuffers                                types.Int64            `tfsdk:"num_io_buffers"`
	RmcacheMemoryAllocationState                types.String           `tfsdk:"rmcache_memory_allocation_state"`
	PerformanceProfile                          types.String           `tfsdk:"performance_profile"`
	SoftwareVersionInfo                         types.String           `tfsdk:"software_version_info"`
	ConfiguredDrlMode                           types.String           `tfsdk:"configured_drl_mode"`
	RfcacheEnabled                              types.Bool             `tfsdk:"rfcache_enabled"`
	MaintenanceState                            types.String           `tfsdk:"maintenance_state"`
	MaintenanceType                             types.String           `tfsdk:"maintenance_type"`
	RfcacheErrorLowResources                    types.Bool             `tfsdk:"rfcache_error_low_resources"`
	RfcacheErrorAPIVersionMismatch              types.Bool             `tfsdk:"rfcache_error_api_version_mismatch"`
	RfcacheErrorInconsistentCacheConfiguration  types.Bool             `tfsdk:"rfcache_error_inconsistent_cache_configuration"`
	RfcacheErrorInconsistentSourceConfiguration types.Bool             `tfsdk:"rfcache_error_inconsistent_source_configuration"`
	RfcacheErrorInvalidDriverPath               types.Bool             `tfsdk:"rfcache_error_invalid_driver_path"`
	RfcacheErrorDeviceDoesNotExist              types.Bool             `tfsdk:"rfcache_error_device_does_not_exist"`
	AuthenticationError                         types.String           `tfsdk:"authentication_error"`
	FglNumConcurrentWrites                      types.Int64            `tfsdk:"fgl_num_concurrent_writes"`
	FglMetadataCacheState                       types.String           `tfsdk:"fgl_metadata_cache_state"`
	FglMetadataCacheSize                        types.Int64            `tfsdk:"fgl_metadata_cache_size"`
	NumRestarts                                 types.Int64            `tfsdk:"num_restarts"`
	LastUpgradeTime                             types.Int64            `tfsdk:"last_upgrade_time"`
	SdsDecoupled                                SdsWindow              `tfsdk:"sds_decoupled"`
	SdsConfigurationFailure                     SdsWindow              `tfsdk:"sds_configuration_failure"`
	SdsReceiveBufferAllocationFailures          SdsWindow              `tfsdk:"sds_receive_buffer_allocation_failures"`
	CertificateInfo                             CertificateInfoModel   `tfsdk:"certificate_info"`
	RaidControllers                             []RaidControllersModel `tfsdk:"raid_controllers"`
	Links                                       []LinkModel            `tfsdk:"links"`
}

// SdsDataSourceModel maps the Sds data source schema data
type SdsDataSourceModel struct {
	SDSIDs               types.List     `tfsdk:"sds_ids"`
	SDSNames             types.List     `tfsdk:"sds_names"`
	ProtectionDomainID   types.String   `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String   `tfsdk:"protection_domain_name"`
	SDSDetails           []SdsDataModel `tfsdk:"sds_details"`
	ID                   types.String   `tfsdk:"id"`
}
