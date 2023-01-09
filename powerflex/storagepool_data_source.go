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
	_ datasource.DataSource              = &storagepoolDataSource{}
	_ datasource.DataSourceWithConfigure = &storagepoolDataSource{}
)

// StoragePoolDataSource is a helper function to simplify the provider implementation.
func StoragePoolDataSource() datasource.DataSource {
	return &storagepoolDataSource{}
}

// storagepoolDataSource is the data source implementation.
type storagepoolDataSource struct {
	client *goscaleio.Client
}

// volume maps the volume schema data.
type volume struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// sdsData maps the SDS schema data
type sdsData struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// linkModel maps the link schema data
type linkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// storagePoolModel maps the storagepool schema data
type storagePoolModel struct {
	ID                                                            types.String `tfsdk:"id"`
	Name                                                          types.String `tfsdk:"name"`
	RebalanceioPriorityPolicy                                     types.String `tfsdk:"rebalance_io_priority_policy"`
	RebuildioPriorityPolicy                                       types.String `tfsdk:"rebuild_io_priority_policy"`
	RebuildioPriorityBwLimitPerDeviceInKbps                       types.Int64  `tfsdk:"rebuild_io_priority_bw_limit_per_device_in_kbps"`
	RebuildioPriorityNumOfConcurrentIosPerDevice                  types.Int64  `tfsdk:"rebuild_io_priority_num_of_concurrent_ios_per_device"`
	RebalanceioPriorityNumOfConcurrentIosPerDevice                types.Int64  `tfsdk:"rebalance_io_priority_num_of_concurrent_ios_per_device"`
	RebalanceioPriorityBwLimitPerDeviceInKbps                     types.Int64  `tfsdk:"rebalance_io_priority_bw_limit_per_device_kbps"`
	RebuildioPriorityAppIopsPerDeviceThreshold                    types.Int64  `tfsdk:"rebuild_io_priority_app_iops_per_device_threshold"`
	RebalanceioPriorityAppIopsPerDeviceThreshold                  types.Int64  `tfsdk:"rebalance_io_priority_app_iops_per_device_threshold"`
	RebuildioPriorityAppBwPerDeviceThresholdInKbps                types.Int64  `tfsdk:"rebuild_io_priority_app_bw_per_device_threshold_kbps"`
	RebalanceioPriorityAppBwPerDeviceThresholdInKbps              types.Int64  `tfsdk:"rebalance_io_priority_app_bw_per_device_threshold_kbps"`
	RebuildioPriorityQuietPeriodInMsec                            types.Int64  `tfsdk:"rebuild_io_priority_quiet_period_msec"`
	RebalanceioPriorityQuietPeriodInMsec                          types.Int64  `tfsdk:"rebalance_io_priority_quiet_period_msec"`
	ZeroPaddingEnabled                                            types.Bool   `tfsdk:"zero_padding_enabled"`
	UseRmcache                                                    types.Bool   `tfsdk:"use_rm_cache"`
	SparePercentage                                               types.Int64  `tfsdk:"spare_percentage"`
	RmCacheWriteHandlingMode                                      types.String `tfsdk:"rm_cache_write_handling_mode"`
	RebuildEnabled                                                types.Bool   `tfsdk:"rebuild_enabled"`
	RebalanceEnabled                                              types.Bool   `tfsdk:"rebalance_enabled"`
	NumofParallelRebuildRebalanceJobsPerDevice                    types.Int64  `tfsdk:"num_of_parallel_rebuild_rebalance_jobs_per_device"`
	BackgroundScannerBWLimitKBps                                  types.Int64  `tfsdk:"background_scanner_bw_limit_kbps"`
	ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice types.Int64  `tfsdk:"protected_maintenance_mode_io_priority_num_of_concurrent_ios_per_device"`
	DataLayout                                                    types.String `tfsdk:"data_layout"`
	VtreeMigrationIoPriorityBwLimitPerDeviceInKbps                types.Int64  `tfsdk:"vtree_migration_io_priority_bw_limit_per_device_kbps"`
	VtreeMigrationIoPriorityPolicy                                types.String `tfsdk:"vtree_migration_io_priority_policy"`
	AddressSpaceUsage                                             types.String `tfsdk:"address_space_usage"`
	ExternalAccelerationType                                      types.String `tfsdk:"external_acceleration_type"`
	PersistentChecksumState                                       types.String `tfsdk:"persistent_checksum_state"`
	UseRfcache                                                    types.Bool   `tfsdk:"use_rf_cache"`
	ChecksumEnabled                                               types.Bool   `tfsdk:"checksum_enabled"`
	CompressionMethod                                             types.String `tfsdk:"compression_method"`
	FragmentationEnabled                                          types.Bool   `tfsdk:"fragmentation_enabled"`
	CapacityUsageState                                            types.String `tfsdk:"capacity_usage_state"`
	CapacityUsageType                                             types.String `tfsdk:"capacity_usage_type"`
	AddressSpaceUsageType                                         types.String `tfsdk:"address_space_usage_type"`
	BgScannerCompareErrorAction                                   types.String `tfsdk:"bg_scanner_compare_error_action"`
	BgScannerReadErrorAction                                      types.String `tfsdk:"bg_scanner_read_error_action"`
	ReplicationCapacityMaxRatio                                   types.Int64  `tfsdk:"replication_capacity_max_ratio"`
	PersistentChecksumEnabled                                     types.Bool   `tfsdk:"persistent_checksum_enabled"`
	PersistentChecksumBuilderLimitKb                              types.Int64  `tfsdk:"persistent_checksum_builder_limit_kb"`
	PersistentChecksumValidateOnRead                              types.Bool   `tfsdk:"persistent_checksum_validate_on_read"`
	VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice           types.Int64  `tfsdk:"vtree_migration_io_priority_num_of_concurrent_ios_per_device"`
	ProtectedMaintenanceModeIoPriorityPolicy                      types.String `tfsdk:"protected_maintenance_mode_io_priority_policy"`
	BackgroundScannerMode                                         types.String `tfsdk:"background_scanner_mode"`
	MediaType                                                     types.String `tfsdk:"media_type"`
	Volumes                                                       []volume     `tfsdk:"volumes"`
	SDS                                                           []sdsData    `tfsdk:"sds"`
	Links                                                         []linkModel  `tfsdk:"links"`
}

// storagepoolDataSourceModel maps the storage pool data source schema data
type storagepoolDataSourceModel struct {
	StoragePoolIDs       types.List         `tfsdk:"storage_pool_ids"`
	StoragePoolNames     types.List         `tfsdk:"storage_pool_names"`
	ProtectionDomainID   types.String       `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String       `tfsdk:"protection_domain_name"`
	StoragePools         []storagePoolModel `tfsdk:"storage_pools"`
	ID                   types.String       `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *storagepoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagepool"
}

// Schema defines the schema for the data source.
func (d *storagepoolDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *storagepoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *storagepoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started storage pool data source read method")
	var state storagepoolDataSourceModel
	var pd *scaleio_types.ProtectionDomain
	var err3 error

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the systems on the PowerFlex cluster
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

	sp := goscaleio.NewStoragePool(d.client)

	spID := []string{}
	// Check if storage pool ID or name is provided
	if !state.StoragePoolIDs.IsNull() {
		diags = state.StoragePoolIDs.ElementsAs(ctx, &spID, true)
	} else if !state.StoragePoolNames.IsNull() {
		diags = state.StoragePoolNames.ElementsAs(ctx, &spID, true)
	} else {
		// Get all the storage pools associated with protection domain
		storagePools, _ := p1.GetStoragePool("")
		for sp := range storagePools {
			spID = append(spID, storagePools[sp].Name)
		}
	}

	if numSP := len(spID); numSP == 0 {
		resp.Diagnostics.AddError("No storage pools found for the specified protection domain", "")
		return
	}

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	for _, spIdentifier := range spID {
		var s1 *scaleio_types.StoragePool

		if !state.StoragePoolIDs.IsNull() {
			s1, err3 = p1.FindStoragePool(spIdentifier, "", "")
		} else {
			s1, err3 = p1.FindStoragePool("", spIdentifier, "")
		}

		if err3 != nil {
			resp.Diagnostics.AddError(
				"Unable to read storage pool",
				err3.Error(),
			)
			return
		}
		sp.StoragePool = s1

		volList, err4 := sp.GetVolume("", "", "", "", false)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Unable to get volumes associated with storage pool",
				err4.Error(),
			)
			return
		}

		sdsList, err5 := sp.GetSDSStoragePool()
		if err5 != nil {
			resp.Diagnostics.AddError(
				"Unable to get SDS associated with storage pool",
				err5.Error(),
			)
			return
		}

		storagePool := getStoragePoolState(volList, sdsList, s1)
		state.StoragePools = append(state.StoragePools, storagePool)
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

func getStoragePoolState(volList []*scaleio_types.Volume, sdsList []scaleio_types.Sds, s1 *scaleio_types.StoragePool) (storagePool storagePoolModel) {
	storagePool = storagePoolModel{
		ID:   types.StringValue(s1.ID),
		Name: types.StringValue(s1.Name),
	}

	// Iterate through volume list
	for _, vol := range volList {
		storagePool.Volumes = append(storagePool.Volumes, volume{
			ID:   types.StringValue(vol.ID),
			Name: types.StringValue(vol.Name),
		})
	}

	// Iterate through SDS list
	for _, sds := range sdsList {
		storagePool.SDS = append(storagePool.SDS, sdsData{
			ID:   types.StringValue(sds.ID),
			Name: types.StringValue(sds.Name),
		})
	}

	// Iterate through the Links
	for _, link := range s1.Links {
		storagePool.Links = append(storagePool.Links, linkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}

	storagePool.RebalanceioPriorityPolicy = types.StringValue(s1.RebalanceioPriorityPolicy)
	storagePool.RebalanceioPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.RebalanceioPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.RebalanceioPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.RebalanceioPriorityAppIopsPerDeviceThreshold))
	storagePool.RebalanceioPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.RebalanceioPriorityBwLimitPerDeviceInKbps))
	storagePool.RebalanceioPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.RebalanceioPriorityQuietPeriodInMsec))
	storagePool.RebalanceioPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.RebalanceioPriorityNumOfConcurrentIosPerDevice))
	storagePool.RebuildioPriorityPolicy = types.StringValue(s1.RebuildioPriorityPolicy)
	storagePool.RebuildioPriorityAppBwPerDeviceThresholdInKbps = types.Int64Value(int64(s1.RebuildioPriorityAppBwPerDeviceThresholdInKbps))
	storagePool.RebuildioPriorityAppIopsPerDeviceThreshold = types.Int64Value(int64(s1.RebuildioPriorityAppIopsPerDeviceThreshold))
	storagePool.RebuildioPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.RebalanceioPriorityBwLimitPerDeviceInKbps))
	storagePool.RebuildioPriorityQuietPeriodInMsec = types.Int64Value(int64(s1.RebuildioPriorityQuietPeriodInMsec))
	storagePool.RebuildioPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.RebuildioPriorityNumOfConcurrentIosPerDevice))
	storagePool.ZeroPaddingEnabled = types.BoolValue(s1.ZeroPaddingEnabled)
	storagePool.UseRmcache = types.BoolValue(s1.UseRmcache)
	storagePool.SparePercentage = types.Int64Value(int64(s1.SparePercentage))
	storagePool.RmCacheWriteHandlingMode = types.StringValue(s1.RmCacheWriteHandlingMode)
	storagePool.RebalanceEnabled = types.BoolValue(s1.RebalanceEnabled)
	storagePool.RebuildEnabled = types.BoolValue(s1.RebuildEnabled)
	storagePool.NumofParallelRebuildRebalanceJobsPerDevice = types.Int64Value(int64(s1.NumofParallelRebuildRebalanceJobsPerDevice))
	storagePool.BackgroundScannerBWLimitKBps = types.Int64Value(int64(s1.BackgroundScannerBWLimitKBps))
	storagePool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice))
	storagePool.DataLayout = types.StringValue(s1.DataLayout)
	storagePool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps = types.Int64Value(int64(s1.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps))
	storagePool.VtreeMigrationIoPriorityPolicy = types.StringValue(s1.VtreeMigrationIoPriorityPolicy)
	storagePool.AddressSpaceUsage = types.StringValue(s1.AddressSpaceUsage)
	storagePool.ExternalAccelerationType = types.StringValue(s1.ExternalAccelerationType)
	storagePool.PersistentChecksumState = types.StringValue(s1.PersistentChecksumState)
	storagePool.UseRfcache = types.BoolValue(s1.UseRfcache)
	storagePool.ChecksumEnabled = types.BoolValue(s1.ChecksumEnabled)
	storagePool.CompressionMethod = types.StringValue(s1.CompressionMethod)
	storagePool.FragmentationEnabled = types.BoolValue(s1.FragmentationEnabled)
	storagePool.CapacityUsageState = types.StringValue(s1.CapacityUsageState)
	storagePool.CapacityUsageType = types.StringValue(s1.CapacityUsageType)
	storagePool.AddressSpaceUsageType = types.StringValue(s1.AddressSpaceUsageType)
	storagePool.BgScannerCompareErrorAction = types.StringValue(s1.BgScannerCompareErrorAction)
	storagePool.BgScannerReadErrorAction = types.StringValue(s1.BgScannerReadErrorAction)
	storagePool.ReplicationCapacityMaxRatio = types.Int64Value(int64(s1.ReplicationCapacityMaxRatio))
	storagePool.PersistentChecksumEnabled = types.BoolValue(s1.PersistentChecksumEnabled)
	storagePool.PersistentChecksumBuilderLimitKb = types.Int64Value(int64(s1.PersistentChecksumBuilderLimitKb))
	storagePool.PersistentChecksumValidateOnRead = types.BoolValue(s1.PersistentChecksumValidateOnRead)
	storagePool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice = types.Int64Value(int64(s1.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice))
	storagePool.ProtectedMaintenanceModeIoPriorityPolicy = types.StringValue(s1.ProtectedMaintenanceModeIoPriorityPolicy)
	storagePool.BackgroundScannerMode = types.StringValue(s1.BackgroundScannerMode)
	storagePool.MediaType = types.StringValue(s1.MediaType)
	return
}
