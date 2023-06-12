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

package powerflex

import (
	"context"

	"github.com/dell/goscaleio"
	scaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sdsDataSource{}
	_ datasource.DataSourceWithConfigure = &sdsDataSource{}
)

// SDSDataSource is a helper function to simplify the provider implementation.
func SDSDataSource() datasource.DataSource {
	return &sdsDataSource{}
}

// sdsDataSource is the data source implementation.
type sdsDataSource struct {
	client *goscaleio.Client
}

type ipList struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

type sdsWindowType struct {
	Threshold            types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec      types.Int64 `tfsdk:"window_size_in_sec"`
	LastOscillationCount types.Int64 `tfsdk:"last_oscillation_count"`
	LastOscillationTime  types.Int64 `tfsdk:"last_oscillation_time"`
	MaxFailuresCount     types.Int64 `tfsdk:"max_failures_count"`
}

type sdsWindow struct {
	ShortWindow  sdsWindowType `tfsdk:"short_window"`
	MediumWindow sdsWindowType `tfsdk:"medium_window"`
	LongWindow   sdsWindowType `tfsdk:"long_window"`
}

type certificateInfoModel struct {
	Subject             types.String `tfsdk:"subject"`
	Issuer              types.String `tfsdk:"issuer"`
	ValidFrom           types.String `tfsdk:"valid_from"`
	ValidTo             types.String `tfsdk:"valid_to"`
	Thumbprint          types.String `tfsdk:"thumbprint"`
	ValidFromAsn1Format types.String `tfsdk:"valid_from_asn1_format"`
	ValidToAsn1Format   types.String `tfsdk:"valid_to_asn1_format"`
}

type raidControllersModel struct {
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

type sdsDataModel struct {
	ID                                          types.String           `tfsdk:"id"`
	Name                                        types.String           `tfsdk:"name"`
	IPList                                      []ipList               `tfsdk:"ip_list"`
	Port                                        types.Int64            `tfsdk:"port"`
	SdsState                                    types.String           `tfsdk:"sds_state"`
	MembershipState                             types.String           `tfsdk:"membership_state"`
	MdmConnectionState                          types.String           `tfsdk:"mdm_connection_state"`
	DrlMode                                     types.String           `tfsdk:"drl_mode"`
	RmcacheEnabled                              types.Bool             `tfsdk:"rmcache_enabled"`
	RmcacheSize                                 types.Int64            `tfsdk:"rmcache_size"`
	RmcacheFrozen                               types.Bool             `tfsdk:"rmcache_frozen"`
	OnVmware                                    types.Bool             `tfsdk:"on_vmware"`
	FaultsetID                                  types.String           `tfsdk:"faultset_id"`
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
	SdsDecoupled                                sdsWindow              `tfsdk:"sds_decoupled"`
	SdsConfigurationFailure                     sdsWindow              `tfsdk:"sds_configuration_failure"`
	SdsReceiveBufferAllocationFailures          sdsWindow              `tfsdk:"sds_receive_buffer_allocation_failures"`
	CertificateInfo                             certificateInfoModel   `tfsdk:"certificate_info"`
	RaidControllers                             []raidControllersModel `tfsdk:"raid_controllers"`
	Links                                       []linkModel            `tfsdk:"links"`
}

// sdsDataSourceModel maps the Sds data source schema data
type sdsDataSourceModel struct {
	SDSIDs               types.List     `tfsdk:"sds_ids"`
	SDSNames             types.List     `tfsdk:"sds_names"`
	ProtectionDomainID   types.String   `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String   `tfsdk:"protection_domain_name"`
	SDSDetails           []sdsDataModel `tfsdk:"sds_details"`
	ID                   types.String   `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *sdsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sds"
}

// Schema defines the schema for the data source.
func (d *sdsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SdsDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *sdsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

// sdsCounterModelValue processes the different types of windows information
func sdsCounterModelValue(s1 scaleio_types.SdsWindow) sdsWindow {
	return sdsWindow{
		ShortWindow: sdsWindowType{
			Threshold:            types.Int64Value(int64(s1.ShortWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
		},
		MediumWindow: sdsWindowType{
			Threshold:            types.Int64Value(int64(s1.MediumWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
		},
		LongWindow: sdsWindowType{
			Threshold:            types.Int64Value(int64(s1.LongWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
		},
	}
}

// sdsCertificateInfo process SDS certificate information and maps to certificateInfoModel struct
func sdsCertificateInfo(s1 scaleio_types.CertificateInfo) certificateInfoModel {
	certicateInfo := certificateInfoModel{}

	if v := s1.Subject; v != "" {
		certicateInfo.Subject = types.StringValue(v)
	} else {
		certicateInfo.Subject = types.StringNull()
	}
	if v := s1.Issuer; v != "" {
		certicateInfo.Issuer = types.StringValue(v)
	} else {
		certicateInfo.Issuer = types.StringNull()
	}
	if v := s1.ValidFrom; v != "" {
		certicateInfo.ValidFrom = types.StringValue(v)
	} else {
		certicateInfo.ValidFrom = types.StringNull()
	}
	if v := s1.ValidTo; v != "" {
		certicateInfo.ValidTo = types.StringValue(v)
	} else {
		certicateInfo.ValidTo = types.StringNull()
	}
	if v := s1.ValidFromAsn1Format; v != "" {
		certicateInfo.ValidFromAsn1Format = types.StringValue(v)
	} else {
		certicateInfo.ValidFromAsn1Format = types.StringNull()
	}
	if v := s1.ValidToAsn1Format; v != "" {
		certicateInfo.ValidFromAsn1Format = types.StringValue(v)
	} else {
		certicateInfo.ValidFromAsn1Format = types.StringNull()
	}

	return certicateInfo
}

// Read refreshes the Terraform state with the latest data.
func (d *sdsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started SDS data source read method")
	var state sdsDataSourceModel
	var pd *scaleio_types.ProtectionDomain
	var err3 error

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the system on the PowerFlex cluster
	c2, err := getFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance",
			err.Error(),
		)
		return
	}

	// Check if protection domain ID or name is provided
	if state.ProtectionDomainID.ValueString() != "" {
		pd, err3 = c2.FindProtectionDomain(state.ProtectionDomainID.ValueString(), "", "")
	} else {
		pd, err3 = c2.FindProtectionDomain("", state.ProtectionDomainName.ValueString(), "")
	}

	if err3 != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err3.Error(),
		)
		return
	}

	p1 := goscaleio.NewProtectionDomainEx(d.client, pd)

	sdsID := []string{}
	// Check if SDS ID or name is provided
	if !state.SDSIDs.IsNull() {
		diags = state.SDSIDs.ElementsAs(ctx, &sdsID, true)
	} else if !state.SDSNames.IsNull() {
		diags = state.SDSNames.ElementsAs(ctx, &sdsID, true)
	} else {
		// Get all the SDS associated with protection domain
		sds, _ := p1.GetSds()
		for sp := range sds {
			sdsID = append(sdsID, sds[sp].Name)
		}
	}

	// Iterate though the SDS for sacing details into state file
	for _, sdsIdentifier := range sdsID {
		var s1 *scaleio_types.Sds

		if !state.SDSIDs.IsNull() {
			s1, err3 = p1.FindSds("ID", sdsIdentifier)
		} else {
			s1, err3 = p1.FindSds("Name", sdsIdentifier)
		}

		if err3 != nil {
			resp.Diagnostics.AddError(
				"Unable to read SDS",
				err3.Error(),
			)
			return
		}
		sdsDetail := getSdsState(s1)
		state.SDSDetails = append(state.SDSDetails, sdsDetail)
	}

	// this is required for acceptance testing
	state.ID = types.StringValue("dummyID")

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// getSdsState saves SDS response in SDS struct
func getSdsState(s1 *scaleio_types.Sds) (sdsDetail sdsDataModel) {
	sdsDetail = sdsDataModel{
		ID:                             types.StringValue(s1.ID),
		Name:                           types.StringValue(s1.Name),
		Port:                           types.Int64Value(int64(s1.Port)),
		SdsState:                       types.StringValue(s1.SdsState),
		MembershipState:                types.StringValue(s1.MembershipState),
		MdmConnectionState:             types.StringValue(s1.MdmConnectionState),
		DrlMode:                        types.StringValue(s1.DrlMode),
		RmcacheEnabled:                 types.BoolValue(s1.RmcacheEnabled),
		RmcacheSize:                    types.Int64Value(int64(s1.RmcacheSizeInKb)),
		RmcacheFrozen:                  types.BoolValue(s1.RmcacheFrozen),
		OnVmware:                       types.BoolValue(s1.IsOnVMware),
		NumIOBuffers:                   types.Int64Value(int64(s1.NumOfIoBuffers)),
		RmcacheMemoryAllocationState:   types.StringValue(s1.RmcacheMemoryAllocationState),
		PerformanceProfile:             types.StringValue(s1.PerformanceProfile),
		SoftwareVersionInfo:            types.StringValue(s1.SoftwareVersionInfo),
		ConfiguredDrlMode:              types.StringValue(s1.ConfiguredDrlMode),
		RfcacheEnabled:                 types.BoolValue(s1.RfcacheEnabled),
		MaintenanceState:               types.StringValue(s1.MaintenanceState),
		MaintenanceType:                types.StringValue(s1.MaintenanceType),
		RfcacheErrorLowResources:       types.BoolValue(s1.RfcacheErrorLowResources),
		RfcacheErrorAPIVersionMismatch: types.BoolValue(s1.RfcacheErrorAPIVersionMismatch),
		RfcacheErrorInconsistentCacheConfiguration:  types.BoolValue(s1.RfcacheErrorInconsistentCacheConfiguration),
		RfcacheErrorInconsistentSourceConfiguration: types.BoolValue(s1.RfcacheErrorInconsistentSourceConfiguration),
		RfcacheErrorInvalidDriverPath:               types.BoolValue(s1.RfcacheErrorInvalidDriverPath),
		RfcacheErrorDeviceDoesNotExist:              types.BoolValue(s1.RfcacheErrorDeviceDoesNotExist),
		AuthenticationError:                         types.StringValue(s1.AuthenticationError),
		FglNumConcurrentWrites:                      types.Int64Value(int64(s1.FglNumConcurrentWrites)),
		FglMetadataCacheSize:                        types.Int64Value(int64(s1.FglMetadataCacheSize)),
		FglMetadataCacheState:                       types.StringValue(s1.FglMetadataCacheState),
		NumRestarts:                                 types.Int64Value(int64(s1.NumRestarts)),
		LastUpgradeTime:                             types.Int64Value(int64(s1.LastUpgradeTime)),
		SdsDecoupled:                                sdsCounterModelValue(s1.SdsDecoupled),
		SdsConfigurationFailure:                     sdsCounterModelValue(s1.SdsConfigurationFailure),
		SdsReceiveBufferAllocationFailures:          sdsCounterModelValue(s1.SdsReceiveBufferAllocationFailures),
		CertificateInfo:                             sdsCertificateInfo(s1.CertificateInfo),
	}

	if v := s1.FaultSetID; v != "" {
		sdsDetail.FaultsetID = types.StringValue(v)
	} else {
		sdsDetail.FaultsetID = types.StringNull()
	}

	// Iterate through IP list
	for _, ip := range s1.IPList {
		sdsDetail.IPList = append(sdsDetail.IPList, ipList{
			IP:   types.StringValue(ip.IP),
			Role: types.StringValue(ip.Role),
		})
	}

	// Iterate through the Links
	for _, link := range s1.Links {
		sdsDetail.Links = append(sdsDetail.Links, linkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}

	// Iterate through the RAID controllers
	for _, raid := range s1.RaidControllers {
		sdsDetail.RaidControllers = append(sdsDetail.RaidControllers, raidControllersModel{
			SerialNumber:    types.StringValue(raid.SerialNumber),
			ModelName:       types.StringValue(raid.ModelName),
			VendorName:      types.StringValue(raid.VendorName),
			FirmwareVersion: types.StringValue(raid.FirmwareVersion),
			DriverVersion:   types.StringValue(raid.DriverVersion),
			DriverName:      types.StringValue(raid.DriverName),
			PciAddress:      types.StringValue(raid.PciAddress),
			Status:          types.StringValue(raid.Status),
			BatteryStatus:   types.StringValue(raid.BatteryStatus),
		})
	}

	return
}
