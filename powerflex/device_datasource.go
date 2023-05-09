package powerflex

import (
	"context"
	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// deviceDataSource returns the volume data source
func DeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

type deviceDataSource struct {
	client *goscaleio.Client
	System *goscaleio.System
}

type deviceDataSourceModel struct {
	ID                   types.String      `tfsdk:"id"`
	Name                 types.String      `tfsdk:"name"`
	StoragePoolID        types.String      `tfsdk:"storage_pool_id"`
	StoragePoolName      types.String      `tfsdk:"storage_pool_name"`
	ProtectionDomainName types.String      `tfsdk:"protection_domain_name"`
	ProtectionDomainID   types.String      `tfsdk:"protection_domain_id"`
	DeviceModel          []deviceModelData `tfsdk:"device_model"`
}

type deviceModelData struct {
	FglNvdimmMetadataAmortizationX100 types.Int64            `tfsdk:"fgl_nvdimm_metadata_amortization_x100"`
	LogicalSectorSizeInBytes          types.Int64            `tfsdk:"logical_sector_size_in_bytes"`
	FglNvdimmWriteCacheSize           types.Int64            `tfsdk:"fgl_nvdimm_write_cache_size"`
	AccelerationPoolID                types.String           `tfsdk:"acceleration_pool_id"`
	RfcacheProps                      RfcachePropsModel      `tfsdk:"rfcache_props"`
	SdsID                             types.String           `tfsdk:"sds_id"`
	StoragePoolID                     types.String           `tfsdk:"storage_pool_id"`
	CapacityLimitInKb                 types.Int64            `tfsdk:"capacity_limit_in_kb"`
	ErrorState                        types.String           `tfsdk:"error_state"`
	Capacity                          types.Int64            `tfsdk:"capacity"`
	DeviceType                        types.String           `tfsdk:"device_type"`
	PersistentChecksumState           types.String           `tfsdk:"persistent_checksum_state"`
	DeviceState                       types.String           `tfsdk:"device_state"`
	LedSetting                        types.String           `tfsdk:"led_setting"`
	MaxCapacityInKb                   types.Int64            `tfsdk:"max_capacity_in_kb"`
	SpSdsID                           types.String           `tfsdk:"sp_sds_id"`
	LongSuccessfulIos                 LongSuccessfulIosModel `tfsdk:"long_successful_ios"`
	AggregatedState                   types.String           `tfsdk:"aggregated_state"`
	TemperatureState                  types.String           `tfsdk:"temperature_state"`
	SsdEndOfLifeState                 types.String           `tfsdk:"ssd_end_of_life_state"`
	ModelName                         types.String           `tfsdk:"model_name"`
	VendorName                        types.String           `tfsdk:"vendor_name"`
	RaidControllerSerialNumber        types.String           `tfsdk:"raid_controller_serial_number"`
	FirmwareVersion                   types.String           `tfsdk:"firmware_version"`
	CacheLookAheadActive              types.Bool             `tfsdk:"cache_look_ahead_active"`
	WriteCacheActive                  types.Bool             `tfsdk:"write_cache_active"`
	AtaSecurityActive                 types.Bool             `tfsdk:"ata_security_active"`
	PhysicalSectorSizeInBytes         types.Int64            `tfsdk:"physical_sector_size_in_bytes"`
	MediaFailing                      types.Bool             `tfsdk:"media_failing"`
	SlotNumber                        types.String           `tfsdk:"slot_number"`
	ExternalAccelerationType          types.String           `tfsdk:"external_acceleration_type"`
	AutoDetectMediaType               types.String           `tfsdk:"auto_detect_media_type"`
	StorageProps                      StoragePropsModel      `tfsdk:"storage_props"`
	AccelerationProps                 AccelerationPropsModel `tfsdk:"acceleration_props"`
	DeviceCurrentPathName             types.String           `tfsdk:"device_current_path_name"`
	DeviceOriginalPathName            types.String           `tfsdk:"device_original_path_name"`
	RfcacheErrorDeviceDoesNotExist    types.Bool             `tfsdk:"rfcache_error_device_does_not_exist"`
	MediaType                         types.String           `tfsdk:"media_type"`
	SerialNumber                      types.String           `tfsdk:"serial_number"`
	Name                              types.String           `tfsdk:"name"`
	ID                                types.String           `tfsdk:"id"`
	Links                             []deviceLinkModel      `tfsdk:"links"`
}

// RfcachePropsModel defines struct for Device
type RfcachePropsModel struct {
	DeviceUUID                     types.String `tfsdk:"device_uuid"`
	RfcacheErrorStuckIO            types.Bool   `tfsdk:"rfcache_error_stuck_io"`
	RfcacheErrorHeavyLoadCacheSkip types.Bool   `tfsdk:"rfcache_error_heavy_load_cache_skip"`
	RfcacheErrorCardIoError        types.Bool   `tfsdk:"rfcache_error_card_io_error"`
}

// LongSuccessfulIosModel defines struct for Device
type LongSuccessfulIosModel struct {
	ShortWindow  DeviceWindowTypeModel `tfsdk:"short_window"`
	MediumWindow DeviceWindowTypeModel `tfsdk:"medium_window"`
	LongWindow   DeviceWindowTypeModel `tfsdk:"long_window"`
}

// DeviceWindowType defines struct for LongSuccessfulIosModel
type DeviceWindowTypeModel struct {
	Threshold            types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec      types.Int64 `tfsdk:"window_size_in_sec"`
	LastOscillationCount types.Int64 `tfsdk:"last_oscillation_count"`
	LastOscillationTime  types.Int64 `tfsdk:"last_oscillation_time"`
	MaxFailuresCount     types.Int64 `tfsdk:"max_failures_count"`
}

// AccelerationPropsModel defines struct for Device
type AccelerationPropsModel struct {
	AccUsedCapacityInKb types.String `tfsdk:"acc_used_capacity_in_kb"`
}

// StoragePropsModel defines struct for Device
type StoragePropsModel struct {
	FglAccDeviceID                   types.String `tfsdk:"fgl_acc_device_id"`
	FglNvdimmSizeMb                  types.Int64  `tfsdk:"fgl_nvdimm_size_mb"`
	DestFglNvdimmSizeMb              types.Int64  `tfsdk:"dest_fgl_nvdimm_size_mb"`
	DestFglAccDeviceID               types.String `tfsdk:"dest_fgl_acc_device_id"`
	ChecksumMode                     types.String `tfsdk:"checksum_mode"`
	DestChecksumMode                 types.String `tfsdk:"dest_checksum_mode"`
	ChecksumAccDeviceID              types.String `tfsdk:"checksum_acc_device_id"`
	DestChecksumAccDeviceID          types.String `tfsdk:"dest_checksum_acc_device_id"`
	ChecksumSizeMb                   types.Int64  `tfsdk:"checksum_size_mb"`
	IsChecksumFullyCalculated        types.Bool   `tfsdk:"is_checksum_fully_calculated"`
	ChecksumChangelogAccDeviceID     types.String `tfsdk:"checksum_changelog_acc_device_id"`
	DestChecksumChangelogAccDeviceID types.String `tfsdk:"dest_checksum_changelog_acc_device_id"`
	ChecksumChangelogSizeMb          types.Int64  `tfsdk:"checksum_changelog_size_mb"`
	DestChecksumChangelogSizeMb      types.Int64  `tfsdk:"dest_checksum_changelog_size_mb"`
}

type deviceLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

func getStoragePool(d *deviceDataSource, storagePoolID string) (*goscaleio.StoragePool, error) {

	system, err := getFirstSystem(d.client)
	if err != nil {
		return nil, err
	}

	sp, err := system.GetStoragePoolByID(storagePoolID)
	if err != nil {
		return nil, err
	}

	sp1 := goscaleio.NewStoragePoolEx(d.client, sp)
	return sp1, nil
}

func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DeviceDataSourceSchema
}

func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

// func (r *deviceDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
// 	var config deviceDataSourceModel

// 	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	if !config.StoragePoolName.IsNull() {
// 		if config.ProtectionDomainID.IsNull() && config.ProtectionDomainName.IsNull() {
// 			resp.Diagnostics.AddAttributeError(
// 				path.Root("storage_pool_name"),
// 				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
// 				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
// 			)
// 		}
// 	}
// }

func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		state deviceDataSourceModel
		pd    *goscaleio.ProtectionDomain
		err   error
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !state.ProtectionDomainID.IsNull() || !state.ProtectionDomainName.IsNull() {
		pd, err = getNewProtectionDomainEx(d.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting protection domain details with ID: "+state.ProtectionDomainID.ValueString()+" name: "+state.ProtectionDomainName.ValueString(),
				err.Error(),
			)
			return
		}
	}
	if !state.StoragePoolID.IsNull() || state.StoragePoolID.ValueString() != "" {
		sp, err := d.System.GetStoragePoolByID(state.StoragePoolID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting storage pool details with ID: "+state.StoragePoolID.ValueString(),
				err.Error(),
			)
			return
		}
		state.StoragePoolName = types.StringValue(sp.Name)
	} else if !state.StoragePoolName.IsUnknown() {
		sp, err1 := pd.FindStoragePool("", state.StoragePoolName.ValueString(), "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Error in getting storage pool details with name: "+state.StoragePoolName.ValueString(),
				err1.Error(),
			)
			return
		}
		state.StoragePoolID = types.StringValue(sp.ID)
	}

	sp, err := getStoragePool(d, state.StoragePoolID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool instance with ID: "+state.StoragePoolID.ValueString(),
			"unexpected error: "+err.Error(),
		)
		return
	}

	devices, err1 := sp.GetDevice()
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting all devices within storage pool instance with ID: "+state.StoragePoolID.ValueString(),
			"unexpected error: "+err1.Error(),
		)
	}
	resp.Diagnostics.AddError(
		helper.PrettyJSON(devices), "dev-avengers",
	)
	// Set refreshed state
	state.DeviceModel = getAllDeviceState(devices)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func getAllDeviceState(devices []scaleiotypes.Device) (response []deviceModelData) {
	for _, device := range devices {
		deviceState := deviceModelData{
			FglNvdimmMetadataAmortizationX100: types.Int64Value(int64(device.FglNvdimmMetadataAmortizationX100)),
			LogicalSectorSizeInBytes:          types.Int64Value(int64(device.LogicalSectorSizeInBytes)),
			FglNvdimmWriteCacheSize:           types.Int64Value(int64(device.FglNvdimmWriteCacheSize)),
			AccelerationPoolID:                types.StringValue(device.AccelerationPoolID),
			SdsID:                             types.StringValue(device.SdsID),
			StoragePoolID:                     types.StringValue(device.StoragePoolID),
			CapacityLimitInKb:                 types.Int64Value(int64(device.CapacityLimitInKb)),
			ErrorState:                        types.StringValue(device.ErrorState),
			Capacity:                          types.Int64Value(int64(device.Capacity)),
			DeviceType:                        types.StringValue(device.DeviceType),
			PersistentChecksumState:           types.StringValue(device.PersistentChecksumState),
			DeviceState:                       types.StringValue(device.DeviceState),
			LedSetting:                        types.StringValue(device.LedSetting),
			MaxCapacityInKb:                   types.Int64Value(int64(device.MaxCapacityInKb)),
			SpSdsID:                           types.StringValue(device.SpSdsID),
			AggregatedState:                   types.StringValue(device.AggregatedState),
			TemperatureState:                  types.StringValue(device.TemperatureState),
			SsdEndOfLifeState:                 types.StringValue(device.SsdEndOfLifeState),
			ModelName:                         types.StringValue(device.ModelName),
			VendorName:                        types.StringValue(device.VendorName),
			RaidControllerSerialNumber:        types.StringValue(device.RaidControllerSerialNumber),
			FirmwareVersion:                   types.StringValue(device.FirmwareVersion),
			CacheLookAheadActive:              types.BoolValue(device.CacheLookAheadActive),
			WriteCacheActive:                  types.BoolValue(device.WriteCacheActive),
			AtaSecurityActive:                 types.BoolValue(device.AtaSecurityActive),
			PhysicalSectorSizeInBytes:         types.Int64Value(int64(device.PhysicalSectorSizeInBytes)),
			MediaFailing:                      types.BoolValue(device.MediaFailing),
			SlotNumber:                        types.StringValue(device.SlotNumber),
			ExternalAccelerationType:          types.StringValue(device.ExternalAccelerationType),
			AutoDetectMediaType:               types.StringValue(device.AutoDetectMediaType),
			DeviceCurrentPathName:             types.StringValue(device.DeviceCurrentPathName),
			DeviceOriginalPathName:            types.StringValue(device.DeviceOriginalPathName),
			RfcacheErrorDeviceDoesNotExist:    types.BoolValue(device.RfcacheErrorDeviceDoesNotExist),
			MediaType:                         types.StringValue(device.MediaType),
			SerialNumber:                      types.StringValue(device.SerialNumber),
			Name:                              types.StringValue(device.SerialNumber),
			ID:                                types.StringValue(device.ID),
		}
		deviceState.RfcacheProps = RfcachePropsModel{
			DeviceUUID:                     types.StringValue(device.RfcacheProps.DeviceUUID),
			RfcacheErrorStuckIO:            types.BoolValue(device.RfcacheProps.RfcacheErrorStuckIO),
			RfcacheErrorHeavyLoadCacheSkip: types.BoolValue(device.RfcacheProps.RfcacheErrorHeavyLoadCacheSkip),
			RfcacheErrorCardIoError:        types.BoolValue(device.RfcacheProps.RfcacheErrorCardIoError),
		}
		deviceState.LongSuccessfulIos = LongSuccessfulIosModel{
			ShortWindow: DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.ShortWindow.MaxFailuresCount)),
			},
			MediumWindow: DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.MediumWindow.MaxFailuresCount)),
			},
			LongWindow: DeviceWindowTypeModel{
				Threshold:            types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.Threshold)),
				WindowSizeInSec:      types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.WindowSizeInSec)),
				LastOscillationCount: types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.LastOscillationCount)),
				LastOscillationTime:  types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.LastOscillationTime)),
				MaxFailuresCount:     types.Int64Value(int64(device.LongSuccessfulIos.LongWindow.MaxFailuresCount)),
			},
		}
		deviceState.StorageProps = StoragePropsModel{
			FglAccDeviceID:                   types.StringValue(device.StorageProps.FglAccDeviceID),
			FglNvdimmSizeMb:                  types.Int64Value(int64(device.StorageProps.FglNvdimmSizeMb)),
			DestFglNvdimmSizeMb:              types.Int64Value(int64(device.StorageProps.DestFglNvdimmSizeMb)),
			DestFglAccDeviceID:               types.StringValue(device.StorageProps.DestFglAccDeviceID),
			ChecksumMode:                     types.StringValue(device.StorageProps.ChecksumMode),
			DestChecksumMode:                 types.StringValue(device.StorageProps.DestChecksumMode),
			ChecksumAccDeviceID:              types.StringValue(device.StorageProps.ChecksumAccDeviceID),
			DestChecksumAccDeviceID:          types.StringValue(device.StorageProps.DestChecksumAccDeviceID),
			ChecksumSizeMb:                   types.Int64Value(int64(device.StorageProps.ChecksumSizeMb)),
			IsChecksumFullyCalculated:        types.BoolValue(device.StorageProps.IsChecksumFullyCalculated),
			ChecksumChangelogAccDeviceID:     types.StringValue(device.StorageProps.ChecksumChangelogAccDeviceID),
			DestChecksumChangelogAccDeviceID: types.StringValue(device.StorageProps.DestChecksumChangelogAccDeviceID),
			ChecksumChangelogSizeMb:          types.Int64Value(int64(device.StorageProps.ChecksumChangelogSizeMb)),
			DestChecksumChangelogSizeMb:      types.Int64Value(int64(device.StorageProps.DestChecksumChangelogSizeMb)),
		}
		deviceState.AccelerationProps = AccelerationPropsModel{
			AccUsedCapacityInKb: types.StringValue(device.AccelerationProps.AccUsedCapacityInKb),
		}
		for _, link := range device.Links {
			deviceState.Links = append(deviceState.Links, deviceLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		response = append(response, deviceState)
	}
	return
}
