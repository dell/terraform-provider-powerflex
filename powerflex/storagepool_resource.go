package powerflex

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &storagepoolResource{}
	_ resource.ResourceWithConfigure   = &storagepoolResource{}
	_ resource.ResourceWithImportState = &storagepoolResource{}
)

// StoragepoolResource - function to return resource interface
func StoragepoolResource() resource.Resource {
	return &storagepoolResource{}
}

type storagepoolResource struct {
	client *goscaleio.Client
}

type storagepoolResourceModel struct {
	ID                                                  types.String `tfsdk:"id"`
	ProtectionDomainID                                  types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName                                types.String `tfsdk:"protection_domain_name"`
	Name                                                types.String `tfsdk:"name"`
	MediaType                                           types.String `tfsdk:"media_type"`
	UseRmcache                                          types.Bool   `tfsdk:"use_rmcache"`
	UseRfcache                                          types.Bool   `tfsdk:"use_rfcache"`
	ZeroPaddingEnabled                                  types.Bool   `tfsdk:"zero_padding_enabled"`
	ReplicationJournalCapacity                          types.Int64  `tfsdk:"replication_journal_capacity"`
	CapacityAlertHighThreshold                          types.Int64  `tfsdk:"capacity_alert_high_threshold"`
	CapacityAlertCriticalThreshold                      types.Int64  `tfsdk:"capacity_alert_critical_threshold"`
	ProtectedMaintenanceModeIoPriorityPolicy            types.String `tfsdk:"protected_maintenance_mode_io_priority_policy"`
	ProtectedMaintenanceModeNumOfConcurrentIosPerDevice types.Int64  `tfsdk:"protected_maintenance_mode_num_of_concurrent_ios_per_device"`
	ProtectedMaintenanceModeBwLimitPerDeviceInKbps      types.Int64  `tfsdk:"protected_maintenance_mode_bw_limit_per_device_in_kbps"`
	RebalanceEnabled                                    types.Bool   `tfsdk:"rebalance_enabled"`
	RebalanceIoPriorityPolicy                           types.String `tfsdk:"rebalance_io_priority_policy"`
	RebalanceNumOfConcurrentIosPerDevice                types.Int64  `tfsdk:"rebalance_num_of_concurrent_ios_per_device"`
	RebalanceBwLimitPerDeviceInKbps                     types.Int64  `tfsdk:"rebalance_bw_limit_per_device_in_kbps"`
	VtreeMigrationIoPriorityPolicy                      types.String `tfsdk:"vtree_migration_io_priority_policy"`
	VtreeMigrationNumOfConcurrentIosPerDevice           types.Int64  `tfsdk:"vtree_migration_num_of_concurrent_ios_per_device"`
	VtreeMigrationBwLimitPerDeviceInKbps                types.Int64  `tfsdk:"vtree_migration_bw_limit_per_device_in_kbps"`
	SparePercentage                                     types.Int64  `tfsdk:"spare_percentage"`
	RmCacheWriteHandlingMode                            types.String `tfsdk:"rm_cache_write_handling_mode"`
	RebuildEnabled                                      types.Bool   `tfsdk:"rebuild_enabled"`
	RebuildRebalanceParallelism                         types.Int64  `tfsdk:"rebuild_rebalance_parallelism"`
	Fragmentation                                       types.Bool   `tfsdk:"fragmentation"`
}

func (r *storagepoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_pool"
}

func (r *storagepoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StoragepoolReourceSchema
}

func (r *storagepoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*goscaleio.Client)
}

func (r *storagepoolResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data storagepoolResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// validate that if the policy is unlimited then the user can't provide values to num of concurrent IOs per device and bandwidth limit per device in Kbps
	if data.ProtectedMaintenanceModeIoPriorityPolicy.ValueString() == "unlimited" {
		if !data.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsNull() || !data.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("protected_maintenance_mode_io_priority_policy"),
				"Attribute Error",
				"With policy as unlimited, it can't add values to num of concurrent IOs per device and bandwidth limit per device in Kbps",
			)
		}
	}

	// validate that if the policy is limitNumOfConcurrentIos then the user can't provide values to bandwidth limit per device in Kbps
	if data.ProtectedMaintenanceModeIoPriorityPolicy.ValueString() == "limitNumOfConcurrentIos" {
		if !data.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("protected_maintenance_mode_io_priority_policy"),
				"Attribute Error",
				"With policy as limitNumOfConcurrentIos, it can't add values to bandwidth limit per device in Kbps",
			)
		}
	}

	// Validate that the policy must be provided in the config in order to configure num of concurrent IOS or bandwidth limit.
	if !data.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsNull() && data.ProtectedMaintenanceModeIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("protected_maintenance_mode_bw_limit_per_device_in_kbps"),
			"Attribute Error",
			"protected_maintenance_mode_io_priority_policy must be provided with a valid value to configure bandwidth limit",
		)
	}

	if !data.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsNull() && data.ProtectedMaintenanceModeIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("protected_maintenance_mode_num_of_concurrent_ios_per_device"),
			"Attribute Error",
			"protected_maintenance_mode_io_priority_policy must be provided with a valid value to configure num of concurrent IOS per device",
		)
	}

	// validate that if the policy is unlimited then the user can't provide values to num of concurrent IOs per device and bandwidth limit per device in Kbps
	if data.RebalanceIoPriorityPolicy.ValueString() == "unlimited" {
		if !data.RebalanceNumOfConcurrentIosPerDevice.IsNull() || !data.RebalanceBwLimitPerDeviceInKbps.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rebalance_io_priority_policy"),
				"Attribute Error",
				"With policy as unlimited, it can't add values to num of concurrent IOs per device and bandwidth limit per device in Kbps",
			)
		}
	}

	// validate that if the policy is limitNumOfConcurrentIos then the user can't provide values to bandwidth limit per device in Kbps
	if data.RebalanceIoPriorityPolicy.ValueString() == "limitNumOfConcurrentIos" {
		if !data.RebalanceBwLimitPerDeviceInKbps.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rebalance_io_priority_policy"),
				"Attribute Error",
				"With policy as limitNumOfConcurrentIos, it can't add values to bandwidth limit per device in Kbps",
			)
		}
	}

	// Validate that the policy must be provided in the config in order to configure num of concurrent IOS or bandwidth limit.
	if !data.RebalanceNumOfConcurrentIosPerDevice.IsNull() && data.RebalanceIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rebalance_num_of_concurrent_ios_per_device"),
			"Attribute Error",
			"rebalance_io_priority_policy must be provided with a valid value to configure bandwidth limit",
		)
	}

	if !data.RebalanceBwLimitPerDeviceInKbps.IsNull() && data.RebalanceIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rebalance_bw_limit_per_device_in_kbps"),
			"Attribute Error",
			"rebalance_io_priority_policy must be provided with a valid value to configure num of concurrent IOS per device",
		)
	}

	// validate that if the policy is limitNumOfConcurrentIos then the user can't provide values to bandwidth limit per device in Kbps
	if data.VtreeMigrationIoPriorityPolicy.ValueString() == "limitNumOfConcurrentIos" {
		if !data.VtreeMigrationBwLimitPerDeviceInKbps.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vtree_migration_io_priority_policy"),
				"Attribute Error",
				"With policy as limitNumOfConcurrentIos, it can't add values to bandwidth limit per device in Kbps",
			)
		}
	}

	// Validate that the policy must be provided in the config in order to configure num of concurrent IOS or bandwidth limit.
	if !data.VtreeMigrationBwLimitPerDeviceInKbps.IsNull() && data.VtreeMigrationIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vtree_migration_bw_limit_per_device_in_kbps"),
			"Attribute Error",
			"vtree_migration_io_priority_policy must be provided with a valid value to configure bandwidth limit",
		)
	}

	if !data.VtreeMigrationNumOfConcurrentIosPerDevice.IsNull() && data.VtreeMigrationIoPriorityPolicy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vtree_migration_num_of_concurrent_ios_per_device"),
			"Attribute Error",
			"vtree_migration_io_priority_policy must be provided with a valid value to configure num of concurrent IOS per device",
		)
	}

	// The write handling mode selection comes into play if rm_cache is enabled
	if !(data.RmCacheWriteHandlingMode.IsNull() || data.RmCacheWriteHandlingMode.IsUnknown()) && !data.UseRmcache.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rm_cache_write_handling_mode"),
			"rm_cache_write_handling_mode cannot be specified while use_rmcache is not set to true",
			"Read Ram cache must be enabled in order to configure its write handling mode",
		)
	}
	// Do I need to add the validation that policy must be present
}

// Function used to Create Storagepool Resource
func (r *storagepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create storagepool")
	// Retrieve values from plan
	var plan storagepoolResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	plan.ProtectionDomainName = types.StringValue(pd.ProtectionDomain.Name)
	payload := &scaleiotypes.StoragePoolParam{
		Name:      plan.Name.ValueString(),
		MediaType: plan.MediaType.ValueString(),
	}

	// enable/disable RmCache
	if plan.UseRmcache.ValueBool() {
		payload.UseRmcache = "true"
	} else {
		payload.UseRmcache = "false"
	}

	// enable/disable RfCache
	if plan.UseRfcache.ValueBool() {
		payload.UseRfcache = "true"
	} else {
		payload.UseRfcache = "false"
	}

	// enable/disable zero padding
	if !plan.ZeroPaddingEnabled.ValueBool() {
		payload.ZeroPaddingEnabled = "false"
	} else {
		payload.ZeroPaddingEnabled = "true"
	}

	// set the spare percentage
	if !plan.SparePercentage.IsNull() && !plan.SparePercentage.IsUnknown() {
		payload.SparePercentage = strconv.FormatInt(plan.SparePercentage.ValueInt64(), 10)
	}

	// set the Rmcache write handling mode when rm_cache is enabled
	if !plan.RmCacheWriteHandlingMode.IsUnknown() && !plan.RmCacheWriteHandlingMode.IsNull() {
		payload.RmcacheWriteHandlingMode = plan.RmCacheWriteHandlingMode.ValueString()
	}

	// create the storage pool
	sp, err := pd.CreateStoragePool(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Storage Pool",
			"Could not create Storage Pool, unexpected error: "+err.Error(),
		)
		return
	}

	// set the replication journal capacity
	if !plan.ReplicationJournalCapacity.IsUnknown() && !plan.ReplicationJournalCapacity.IsNull() {
		err := pd.SetReplicationJournalCapacity(sp, plan.ReplicationJournalCapacity.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set replication Journal capacity to %s", plan.ReplicationJournalCapacity.String()),
				err.Error(),
			)
		}
	}

	// set the capacity alert threshold - high or critical
	if (!plan.CapacityAlertHighThreshold.IsUnknown() && !plan.CapacityAlertHighThreshold.IsNull()) ||
		(!plan.CapacityAlertCriticalThreshold.IsUnknown() && !plan.CapacityAlertCriticalThreshold.IsNull()) {
		capacityAlertThreshold := &scaleiotypes.CapacityAlertThresholdParam{
			CapacityAlertHighThresholdPercent:     strconv.FormatInt(plan.CapacityAlertHighThreshold.ValueInt64(), 10),
			CapacityAlertCriticalThresholdPercent: strconv.FormatInt(plan.CapacityAlertCriticalThreshold.ValueInt64(), 10),
		}
		err := pd.SetCapacityAlertThreshold(sp, capacityAlertThreshold)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set capacity alert high threshold to %s", plan.CapacityAlertHighThreshold.String()),
				err.Error(),
			)
		}
	}

	// set the protected maintenance mode IO priority policy
	if (!plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsNull()) ||
		(!plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && !plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsNull()) ||
		(!plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown() && !plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsNull()) {
		protectedMaintenance := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy: plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString(),
		}
		// when the value is null or unknown the below line will return 0 and if it's 0 then we don't need to add it to the payload
		if strconv.FormatInt(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.ValueInt64(), 10) != "0" {
			protectedMaintenance.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.ValueInt64(), 10)
		}
		if strconv.FormatInt(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.ValueInt64(), 10) != "0" {
			protectedMaintenance.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.ValueInt64(), 10)
		}
		err := pd.SetProtectedMaintenanceModeIoPriorityPolicy(sp, protectedMaintenance)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set protected maintenance mode Io priority policy  to %s", plan.ProtectedMaintenanceModeIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	}

	// set rebalance enabled
	if !plan.RebalanceEnabled.IsUnknown() && !plan.RebalanceEnabled.IsNull() {
		err := pd.SetRebalanceEnabled(sp, plan.RebalanceEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalance enabled to %s", plan.RebalanceEnabled.String()),
				err.Error(),
			)
		}
	}

	// set the rebalance IO priority policy
	if (!plan.RebalanceIoPriorityPolicy.IsUnknown() && !plan.RebalanceIoPriorityPolicy.IsNull()) ||
		(!plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && !plan.RebalanceNumOfConcurrentIosPerDevice.IsNull()) ||
		(!plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown() && !plan.RebalanceBwLimitPerDeviceInKbps.IsNull()) {
		rebalanceIoPriorityPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy: plan.RebalanceIoPriorityPolicy.ValueString(),
		}
		if strconv.FormatInt(plan.RebalanceNumOfConcurrentIosPerDevice.ValueInt64(), 10) != "0" {
			rebalanceIoPriorityPolicy.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.RebalanceNumOfConcurrentIosPerDevice.ValueInt64(), 10)
		}
		if strconv.FormatInt(plan.RebalanceBwLimitPerDeviceInKbps.ValueInt64(), 10) != "0" {
			rebalanceIoPriorityPolicy.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.RebalanceBwLimitPerDeviceInKbps.ValueInt64(), 10)
		}
		err := pd.SetRebalanceIoPriorityPolicy(sp, rebalanceIoPriorityPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalance Io priority policy  to %s", plan.RebalanceIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	}

	// set vtree migration IO priority policy
	if (!plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && !plan.VtreeMigrationIoPriorityPolicy.IsNull()) ||
		(!plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && !plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsNull()) ||
		(!plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown() && !plan.VtreeMigrationBwLimitPerDeviceInKbps.IsNull()) {
		vtreeMigrationPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy: plan.VtreeMigrationIoPriorityPolicy.ValueString(),
		}
		if strconv.FormatInt(plan.VtreeMigrationNumOfConcurrentIosPerDevice.ValueInt64(), 10) != "0" {
			vtreeMigrationPolicy.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.VtreeMigrationNumOfConcurrentIosPerDevice.ValueInt64(), 10)
		}
		if strconv.FormatInt(plan.VtreeMigrationBwLimitPerDeviceInKbps.ValueInt64(), 10) != "0" {
			vtreeMigrationPolicy.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.VtreeMigrationBwLimitPerDeviceInKbps.ValueInt64(), 10)
		}
		err := pd.SetVTreeMigrationIOPriorityPolicy(sp, vtreeMigrationPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set Vtree migration Io priority policy  to %s", plan.VtreeMigrationIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	}

	// set rebuild enabled
	if !plan.RebuildEnabled.IsUnknown() && !plan.RebuildEnabled.IsNull() {
		err := pd.SetRebuildEnabled(sp, plan.RebuildEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild enabled to %s", plan.RebalanceEnabled.String()),
				err.Error(),
			)
		}
	}

	// set rebuild rebalance parallelism
	if !plan.RebuildRebalanceParallelism.IsUnknown() && !plan.RebuildRebalanceParallelism.IsNull() {
		err := pd.SetRebuildRebalanceParallelismParam(sp, plan.RebuildRebalanceParallelism.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild rebalance parallelism to %s", plan.RebuildRebalanceParallelism.String()),
				err.Error(),
			)
		}
	}

	// set the fragmentation
	if !plan.Fragmentation.IsUnknown() {
		err := pd.Fragmentation(sp, plan.Fragmentation.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set fragmentation to %s", plan.Fragmentation.String()),
				err.Error(),
			)
		}
	}

	spResponse, err := pd.FindStoragePool(sp, "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Storagepool after creation",
			"Could not get Storagepool, unexpected error: "+err.Error(),
		)
		return
	}

	state := updateStoragepoolState(spResponse, plan)
	tflog.Debug(ctx, "Create Storagepool :-- "+helper.PrettyJSON(sp))
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Storagepool Resource
func (r *storagepoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read Storagepool")
	// Get current state
	var state storagepoolResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster", err.Error(),
		)
		return
	}

	spr, err := system.GetStoragePoolByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get storagepool by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}
	spResponse := updateStoragepoolState(spr, state)

	if state.ProtectionDomainName.IsNull() {
		protectionDomain, err := system.FindProtectionDomain(spr.ProtectionDomainID, "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to read name of protection domain of ID %s for Storagepool %s", spr.ProtectionDomainID, spr.Name),
				err.Error(),
			)
		} else {
			spResponse.ProtectionDomainName = types.StringValue(protectionDomain.Name)
		}
	}

	tflog.Debug(ctx, "Read Storagepool :-- "+helper.PrettyJSON(spr))
	diags = resp.State.Set(ctx, spResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update Storagepool Resource
func (r *storagepoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update Storagepool")
	// Retrieve values from plan
	var plan storagepoolResourceModel
	var err1 error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state storagepoolResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	spResponse, err := pd.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting Storagepool", err.Error(),
		)
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		_, err := pd.ModifyStoragePoolName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating name of Storagepool", err.Error(),
			)
		}
	}

	if plan.MediaType.ValueString() != state.MediaType.ValueString() {
		_, err := pd.ModifyStoragePoolMedia(state.ID.ValueString(), plan.MediaType.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating media type of Storagepool", err.Error(),
			)
		}
	}

	rm := goscaleio.NewStoragePoolEx(r.client, spResponse)

	if !state.UseRmcache.Equal(plan.UseRmcache) {
		err := rm.ModifyRMCache(plan.UseRmcache.String())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating rm_cache of Storagepool", err.Error(),
			)
		}
	}

	if !state.UseRfcache.Equal(plan.UseRfcache) {
		if plan.UseRfcache.String() == "true" {
			_, err1 = pd.EnableRFCache(spResponse.ID)

		} else {
			_, err1 = pd.DisableRFCache(spResponse.ID)
		}
	}

	if !plan.ZeroPaddingEnabled.IsUnknown() &&
		!state.ZeroPaddingEnabled.Equal(plan.ZeroPaddingEnabled) {
		errZeroPaddingEnabled := pd.EnableOrDisableZeroPadding(spResponse.ID, plan.ZeroPaddingEnabled.String())
		if errZeroPaddingEnabled != nil {
			resp.Diagnostics.AddError(
				"Error while updating ZeroPadding settings of Storagepool", errZeroPaddingEnabled.Error(),
			)
		}
	}

	if !plan.ReplicationJournalCapacity.IsUnknown() &&
		!state.ReplicationJournalCapacity.Equal(plan.ReplicationJournalCapacity) {
		errReplicationJournalCapacity := pd.SetReplicationJournalCapacity(spResponse.ID, strconv.FormatInt(plan.ReplicationJournalCapacity.ValueInt64(), 10))
		if errReplicationJournalCapacity != nil {
			resp.Diagnostics.AddError(
				"Error while updating ReplicationJournalCapacity of Storagepool", errReplicationJournalCapacity.Error(),
			)
		}
	}

	if capacityAlertThresholdParam, ok := isCritcalAlert(plan, state); !ok {
		errSetCapacityAlertThreshold := pd.SetCapacityAlertThreshold(spResponse.ID, capacityAlertThresholdParam)
		if errSetCapacityAlertThreshold != nil {
			resp.Diagnostics.AddError(
				"Error while updating Capacity Alert Thresholds of Storagepool", errSetCapacityAlertThreshold.Error(),
			)
		}
	}

	if protectedMaintenanceModeParam, ok := isProtectedMaintenance(plan, state); !ok {
		errProtectedMaintenanceModeIoPriorityPolicy := pd.SetProtectedMaintenanceModeIoPriorityPolicy(spResponse.ID, protectedMaintenanceModeParam)
		if errProtectedMaintenanceModeIoPriorityPolicy != nil {
			resp.Diagnostics.AddError(
				"Error while updating Protect Maintance Policy/NumOfConcurrentIosPerDevice/BwLimitPerDeviceInKbps of Storagepool", errProtectedMaintenanceModeIoPriorityPolicy.Error(),
			)
		}
	}

	if !plan.RebalanceEnabled.IsUnknown() &&
		!state.RebalanceEnabled.Equal(plan.RebalanceEnabled) {
		err := pd.SetRebalanceEnabled(spResponse.ID, plan.RebalanceEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error while updating rebalance enabled to %s", plan.RebalanceEnabled.String()),
				err.Error(),
			)
		}
	}

	if rebalanceIoPriorityPolicy, ok := isRebalance(plan, state); !ok {
		errRebalanceIoPrioritypolicy := pd.SetRebalanceIoPriorityPolicy(spResponse.ID, rebalanceIoPriorityPolicy)
		if errRebalanceIoPrioritypolicy != nil {
			resp.Diagnostics.AddError(
				"Error while updating Rebalance Policy/NumOfConcurrentIosPerDevice/BwLimitPerDeviceInKbps of Storagepool", errRebalanceIoPrioritypolicy.Error(),
			)
		}
	}

	if vtreeMigrationPolicy, ok := isVtreeMigration(plan, state); !ok {
		errVtreeMigrationIoPriorityPolicy := pd.SetVTreeMigrationIOPriorityPolicy(spResponse.ID, vtreeMigrationPolicy)
		if errVtreeMigrationIoPriorityPolicy != nil {
			resp.Diagnostics.AddError(
				"Error while updating Vtree Migration Policy/NumOfConcurrentIosPerDevice/BwLimitPerDeviceInKbps of Storagepool", errVtreeMigrationIoPriorityPolicy.Error(),
			)
		}
	}

	if !plan.SparePercentage.IsUnknown() &&
		!state.SparePercentage.Equal(plan.SparePercentage) {
		errSparePercentage := pd.SetSparePercentage(spResponse.ID, strconv.FormatInt(plan.SparePercentage.ValueInt64(), 10))
		if errSparePercentage != nil {
			resp.Diagnostics.AddError(
				"Error while updating SparePercentage of Storagepool", errSparePercentage.Error(),
			)
		}
	}

	if !plan.RmCacheWriteHandlingMode.IsUnknown() &&
		!state.RmCacheWriteHandlingMode.Equal(plan.RmCacheWriteHandlingMode) {
		errRmCacheWriteHandlingMode := pd.SetRMcacheWriteHandlingMode(spResponse.ID, plan.RmCacheWriteHandlingMode.ValueString())
		if errRmCacheWriteHandlingMode != nil {
			resp.Diagnostics.AddError(
				"Error while updating RmCacheWriteHandlingMode of Storagepool", errRmCacheWriteHandlingMode.Error(),
			)
		}
	}

	if !plan.RebuildEnabled.IsUnknown() &&
		!state.RebuildEnabled.Equal(plan.RebuildEnabled) {
		errRebuildEnabled := pd.SetRebuildEnabled(spResponse.ID, plan.RebuildEnabled.String())
		if errRebuildEnabled != nil {
			resp.Diagnostics.AddError(
				"Error while updating RebuildEnabled of Storagepool", errRebuildEnabled.Error(),
			)
		}
	}

	if !plan.RebuildRebalanceParallelism.IsUnknown() &&
		!state.RebuildRebalanceParallelism.Equal(plan.RebuildRebalanceParallelism) {
		errRebuildRebalanceParallelism := pd.SetRebuildRebalanceParallelismParam(spResponse.ID, strconv.FormatInt(plan.RebuildRebalanceParallelism.ValueInt64(), 10))
		if errRebuildRebalanceParallelism != nil {
			resp.Diagnostics.AddError(
				"Error updating RebuildRebalanceParallelism settings of Storagepool", errRebuildRebalanceParallelism.Error(),
			)
		}
	}

	if !plan.Fragmentation.IsUnknown() &&
		!state.Fragmentation.Equal(plan.Fragmentation) {
		errFragmentation := pd.Fragmentation(spResponse.ID, plan.Fragmentation.ValueBool())
		if errFragmentation != nil {
			resp.Diagnostics.AddError(
				"Error updating Fragmentation settings of Storagepool", errFragmentation.Error(),
			)
		}
	}

	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error while updating rf_cache of Storagepool", err.Error(),
		)
	}

	spResponse, err = pd.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting Storagepool", err.Error(),
		)
		return
	}

	state1 := updateStoragepoolState(spResponse, state)
	tflog.Debug(ctx, "Update Storagepool :-- "+helper.PrettyJSON(spResponse))
	diags = resp.State.Set(ctx, state1)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Storagepool Resource
func (r *storagepoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete Storagepool")
	// Retrieve values from state
	var state storagepoolResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	err = pd.DeleteStoragePool(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Storagepool",
			"Couldn't Delete Storagepool "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for Storagepool Resource
func (r *storagepoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to update the State for Storagepool Resource
func updateStoragepoolState(storagepool *scaleiotypes.StoragePool, plan storagepoolResourceModel) storagepoolResourceModel {
	state := plan
	state.ProtectionDomainID = types.StringValue(storagepool.ProtectionDomainID)
	state.ID = types.StringValue(storagepool.ID)
	state.Name = types.StringValue(storagepool.Name)
	state.MediaType = types.StringValue(storagepool.MediaType)
	state.UseRmcache = types.BoolValue(storagepool.UseRmcache)
	state.UseRfcache = types.BoolValue(storagepool.UseRfcache)
	state.ZeroPaddingEnabled = types.BoolValue(storagepool.ZeroPaddingEnabled)
	state.ReplicationJournalCapacity = types.Int64Value(int64(storagepool.ReplicationCapacityMaxRatio))
	state.CapacityAlertHighThreshold = types.Int64Value(int64(storagepool.CapacityAlertHighThreshold))
	state.CapacityAlertCriticalThreshold = types.Int64Value(int64(storagepool.CapacityAlertCriticalThreshold))
	state.ProtectedMaintenanceModeIoPriorityPolicy = types.StringValue(storagepool.ProtectedMaintenanceModeIoPriorityPolicy)
	state.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice))
	state.ProtectedMaintenanceModeBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps))
	state.RebalanceEnabled = types.BoolValue(storagepool.RebalanceEnabled)
	state.RebalanceIoPriorityPolicy = types.StringValue(storagepool.RebalanceioPriorityPolicy)
	state.RebalanceNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.RebalanceioPriorityNumOfConcurrentIosPerDevice))
	state.RebalanceBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.RebalanceioPriorityBwLimitPerDeviceInKbps))
	state.VtreeMigrationIoPriorityPolicy = types.StringValue(storagepool.VtreeMigrationIoPriorityPolicy)
	state.VtreeMigrationNumOfConcurrentIosPerDevice = types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice))
	state.VtreeMigrationBwLimitPerDeviceInKbps = types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps))
	state.SparePercentage = types.Int64Value(int64(storagepool.SparePercentage))
	state.RmCacheWriteHandlingMode = types.StringValue(storagepool.RmCacheWriteHandlingMode)
	state.RebuildEnabled = types.BoolValue(storagepool.RebuildEnabled)
	state.RebuildRebalanceParallelism = types.Int64Value(int64(storagepool.NumofParallelRebuildRebalanceJobsPerDevice))
	state.Fragmentation = types.BoolValue(storagepool.FragmentationEnabled)
	return state
}

func isCritcalAlert(plan, state storagepoolResourceModel) (*scaleiotypes.CapacityAlertThresholdParam, bool) {
	payload, ok := scaleiotypes.CapacityAlertThresholdParam{}, true
	if !plan.CapacityAlertHighThreshold.IsUnknown() && !state.CapacityAlertHighThreshold.Equal(plan.CapacityAlertHighThreshold) {
		ok = false
		payload.CapacityAlertHighThresholdPercent = strconv.FormatInt(plan.CapacityAlertHighThreshold.ValueInt64(), 10)
	}
	if !plan.CapacityAlertCriticalThreshold.IsUnknown() && !state.CapacityAlertCriticalThreshold.Equal(plan.CapacityAlertCriticalThreshold) {
		ok = false
		payload.CapacityAlertCriticalThresholdPercent = strconv.FormatInt(plan.CapacityAlertCriticalThreshold.ValueInt64(), 10)
	}
	return &payload, ok
}

func isProtectedMaintenance(plan, state storagepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && !state.ProtectedMaintenanceModeIoPriorityPolicy.Equal(plan.ProtectedMaintenanceModeIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.ProtectedMaintenanceModeIoPriorityPolicy.ValueString()
	}
	if !plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && !state.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.Equal(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown() && !state.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.Equal(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}

func isRebalance(plan, state storagepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.RebalanceIoPriorityPolicy.IsUnknown() && !state.RebalanceIoPriorityPolicy.Equal(plan.RebalanceIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.RebalanceIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.RebalanceIoPriorityPolicy.ValueString()
	}
	if !plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && !state.RebalanceNumOfConcurrentIosPerDevice.Equal(plan.RebalanceNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.RebalanceNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown() && !state.RebalanceBwLimitPerDeviceInKbps.Equal(plan.RebalanceBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.RebalanceBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}

func isVtreeMigration(plan, state storagepoolResourceModel) (*scaleiotypes.ProtectedMaintenanceModeParam, bool) {
	payload, ok := scaleiotypes.ProtectedMaintenanceModeParam{}, true
	if !plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && !state.VtreeMigrationIoPriorityPolicy.Equal(plan.VtreeMigrationIoPriorityPolicy) {
		ok = false
		payload.Policy = plan.VtreeMigrationIoPriorityPolicy.ValueString()
	} else {
		payload.Policy = state.VtreeMigrationIoPriorityPolicy.ValueString()
	}
	if !plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && !state.VtreeMigrationNumOfConcurrentIosPerDevice.Equal(plan.VtreeMigrationNumOfConcurrentIosPerDevice) {
		ok = false
		payload.NumOfConcurrentIosPerDevice = strconv.FormatInt(plan.VtreeMigrationNumOfConcurrentIosPerDevice.ValueInt64(), 10)
	}
	if !plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown() && !state.VtreeMigrationBwLimitPerDeviceInKbps.Equal(plan.VtreeMigrationBwLimitPerDeviceInKbps) {
		ok = false
		payload.BwLimitPerDeviceInKbps = strconv.FormatInt(plan.VtreeMigrationBwLimitPerDeviceInKbps.ValueInt64(), 10)
	}
	return &payload, ok
}
