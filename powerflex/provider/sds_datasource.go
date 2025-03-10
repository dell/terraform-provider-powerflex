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

package provider

import (
	"context"
	"fmt"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

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

// Metadata returns the data source type name.
func (d *sdsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sds"
}

// Schema defines the schema for the data source.
func (d *sdsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SdsDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *sdsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

// sdsCounterModelValue processes the different types of windows information
func sdsCounterModelValue(s1 scaleio_types.SdsWindow) models.SdsWindow {
	return models.SdsWindow{
		ShortWindow: models.SdsWindowType{
			Threshold:            types.Int64Value(int64(s1.ShortWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.ShortWindow.WindowSizeInSec)),
		},
		MediumWindow: models.SdsWindowType{
			Threshold:            types.Int64Value(int64(s1.MediumWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.MediumWindow.WindowSizeInSec)),
		},
		LongWindow: models.SdsWindowType{
			Threshold:            types.Int64Value(int64(s1.LongWindow.Threshold)),
			WindowSizeInSec:      types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			LastOscillationCount: types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			LastOscillationTime:  types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
			MaxFailuresCount:     types.Int64Value(int64(s1.LongWindow.WindowSizeInSec)),
		},
	}
}

// sdsCertificateInfo process SDS certificate information and maps to certificateInfoModel struct
func sdsCertificateInfo(s1 scaleio_types.CertificateInfo) models.CertificateInfoModel {
	certicateInfo := models.CertificateInfoModel{}

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
	var (
		state   models.SdsDataSourceModel
		sdsList []models.SdsDataModel
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the system on the PowerFlex cluster
	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance",
			err.Error(),
		)
		return
	}

	//Gather all Sdses
	sdsUnTranslated, err := system.GetAllSds()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Sdses",
			err.Error(),
		)
		return
	}

	// Check if SdsDataFilter is provided
	if state.SdsDataFilter != nil {
		sdsfiltered, err := helper.GetDataSourceByValue(*state.SdsDataFilter, sdsUnTranslated)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering SDS: %v please validate the filter", state.SdsDataFilter), err.Error(),
			)
			return
		}
		var filteredSds []scaleio_types.Sds
		for _, val := range sdsfiltered {
			filteredSds = append(filteredSds, val.(scaleio_types.Sds))
		}
		sdsUnTranslated = filteredSds
	}

	// translate goscaleio model to provider model
	for _, sds := range sdsUnTranslated {
		sdsTranslated := getSdsState(sds)
		sdsList = append(sdsList, sdsTranslated)
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.SDSDetails = sdsList

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
func getSdsState(s1 scaleio_types.Sds) (sdsDetail models.SdsDataModel) {
	sdsDetail = models.SdsDataModel{
		ID:                             types.StringValue(s1.ID),
		Name:                           types.StringValue(s1.Name),
		ProtectionDomainID:             types.StringValue(s1.ProtectionDomainID),
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
		sdsDetail.IPList = append(sdsDetail.IPList, models.IPList{
			IP:   types.StringValue(ip.IP),
			Role: types.StringValue(ip.Role),
		})
	}

	// Iterate through the Links
	for _, link := range s1.Links {
		sdsDetail.Links = append(sdsDetail.Links, models.LinkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}

	// Iterate through the RAID controllers
	for _, raid := range s1.RaidControllers {
		sdsDetail.RaidControllers = append(sdsDetail.RaidControllers, models.RaidControllersModel{
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
