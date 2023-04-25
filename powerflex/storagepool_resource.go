package powerflex

import (
	"context"
	"fmt"
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
	ID                   types.String `tfsdk:"id"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	Name                 types.String `tfsdk:"name"`
	MediaType            types.String `tfsdk:"media_type"`
	UseRmcache           types.Bool   `tfsdk:"use_rmcache"`
	UseRfcache           types.Bool   `tfsdk:"use_rfcache"`
	ZeroPaddingEnabled	 types.Bool	  `tfsdk:"zero_padding_enabled"`
	ReplicationJournalCapacity types.Int64	`tfsdk:"replication_journal_capacity"`
	CapacityAlertHighThreshold types.Int64 `tfsdk:"capacity_alert_high_threshold"`
	CapacityAlertCriticalThreshold types.Int64 `tfsdk:"capacity_alert_critical_threshold"`
	ProtectedMaintenanceModeIoPriorityPolicy  types.String `tfsdk:"protected_maintenance_mode_io_priority_policy"`
	ProtectedMaintenanceModeNumOfConcurrentIosPerDevice  types.Int64 `tfsdk:"protected_maintenance_mode_num_of_concurrent_ios_per_device"`
	ProtectedMaintenanceModeBwLimitPerDeviceInKbps  types.Int64 `tfsdk:"protected_maintenance_mode_bw_limit_per_device_in_kbps"`
	RebalanceEnabled  types.Bool `tfsdk:"rebalance_enabled"`
	RebalanceIoPriorityPolicy  types.String `tfsdk:"rebalance_io_priority_policy"`
	RebalanceNumOfConcurrentIosPerDevice  types.Int64 `tfsdk:"rebalance_num_of_concurrent_ios_per_device"`
	RebalanceBwLimitPerDeviceInKbps  types.Int64 `tfsdk:"rebalance_bw_limit_per_device_in_kbps"`
	VtreeMigrationIoPriorityPolicy  types.String `tfsdk:"vtree_migration_io_priority_policy"`
	VtreeMigrationNumOfConcurrentIosPerDevice  types.Int64 `tfsdk:"vtree_migration_num_of_concurrent_ios_per_device"`
	VtreeMigrationBwLimitPerDeviceInKbps  types.Int64 `tfsdk:"vtree_migration_bw_limit_per_device_in_kbps"`
	SparePercentage types.Int64 `tfsdk:"spare_percentage"`
	RmCacheWriteHandlingMode types.String `tfsdk:"rm_cache_write_handling_mode"`
	RebuildEnabled types.Bool `tfsdk:"rebuild_enabled"`
	RebuildRebalanceParallelism types.Int64 `tfsdk:"rebuild_rebalance_parallelism"`
	Fragmentation types.Bool `tfsdk:"fragmentation"`
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

	if plan.UseRmcache.String() == "true" {
		payload.UseRmcache = "true"
	} else {
		payload.UseRmcache = "false"
	}

	if plan.UseRfcache.String() == "true" {
		payload.UseRfcache = "true"
	} else {
		payload.UseRfcache = "false"
	}

	if plan.ZeroPaddingEnabled.String() == "false" {
		payload.ZeroPaddingEnabled = "false"
	} else {
		payload.ZeroPaddingEnabled = "true"
	}

	if !plan.SparePercentage.IsUnknown() && !plan.SparePercentage.IsNull() {
		payload.SparePercentage = plan.SparePercentage.String()
	}

	if !plan.RmCacheWriteHandlingMode.IsUnknown() && !plan.RmCacheWriteHandlingMode.IsNull() {
		payload.RmcacheWriteHandlingMode = plan.RmCacheWriteHandlingMode.ValueString()
	}

	sp, err := pd.CreateStoragePool(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Storage Pool",
			"Could not create Storage Pool, unexpected error: "+err.Error(),
		)
		return
	}

	if !plan.ReplicationJournalCapacity.IsUnknown() && !plan.ReplicationJournalCapacity.IsNull(){
		err := pd.SetReplicationJournalCapacity(sp,plan.ReplicationJournalCapacity.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set replication Journal capacity to %s", plan.ReplicationJournalCapacity.String()),
				err.Error(),
			)
		}
	}

	if !plan.CapacityAlertHighThreshold.IsUnknown() && plan.CapacityAlertCriticalThreshold.IsUnknown(){
		capacityAlertThreshold := &scaleiotypes.CapacityAlertThresholdParam{
			CapacityAlertHighThresholdPercent:      plan.CapacityAlertHighThreshold.String(),
		}
		err := pd.SetCapacityAlertThreshold(sp,capacityAlertThreshold)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set capacity alert high threshold to %s", plan.CapacityAlertHighThreshold.String()),
				err.Error(),
			)
		}
	} else if !plan.CapacityAlertCriticalThreshold.IsUnknown() && plan.CapacityAlertHighThreshold.IsUnknown() {
		capacityAlertThreshold := &scaleiotypes.CapacityAlertThresholdParam{
			CapacityAlertCriticalThresholdPercent:      plan.CapacityAlertCriticalThreshold.String(),
		}
		err := pd.SetCapacityAlertThreshold(sp,capacityAlertThreshold)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set capacity alert critical threshold to %s", plan.CapacityAlertCriticalThreshold.String()),
				err.Error(),
			)
		}
	} else if !plan.CapacityAlertCriticalThreshold.IsUnknown() && !plan.CapacityAlertCriticalThreshold.IsNull() && !plan.CapacityAlertHighThreshold.IsUnknown() && !plan.CapacityAlertHighThreshold.IsNull(){
		capacityAlertThreshold := &scaleiotypes.CapacityAlertThresholdParam{
			CapacityAlertCriticalThresholdPercent:      plan.CapacityAlertCriticalThreshold.String(),
			CapacityAlertHighThresholdPercent:      plan.CapacityAlertHighThreshold.String(),
		}
		err := pd.SetCapacityAlertThreshold(sp,capacityAlertThreshold)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set capacity alert high threshold to %s and critical threshold to %s", plan.CapacityAlertHighThreshold.String(),plan.CapacityAlertCriticalThreshold.String()),
				err.Error(),
			)
		}
	}

	if !plan.CapacityAlertHighThreshold.IsUnknown() && plan.CapacityAlertCriticalThreshold.IsUnknown(){
		capacityAlertThreshold := &scaleiotypes.CapacityAlertThresholdParam{
			CapacityAlertHighThresholdPercent:      plan.CapacityAlertHighThreshold.String(),
		}
		err := pd.SetCapacityAlertThreshold(sp,capacityAlertThreshold)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set capacity alert high threshold to %s", plan.CapacityAlertHighThreshold.String()),
				err.Error(),
			)
		}
	}

	if !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown(){
		protectedMaintenance := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString(),
		}
		err := pd.SetProtectedMaintenanceModeIoPriorityPolicy(sp,protectedMaintenance)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set protected maintenance mode Io priority policy  to %s", plan.ProtectedMaintenanceModeIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	} else if !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && !plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown(){
		protectedMaintenance := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.String(),
		}
		err := pd.SetProtectedMaintenanceModeIoPriorityPolicy(sp,protectedMaintenance)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set protected maintenance mode Io priority policy  to %s and num of concurrent IOs oer device to %s", plan.ProtectedMaintenanceModeIoPriorityPolicy.String(),plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.String()),
				err.Error(),
			)
		}
	} else if !plan.ProtectedMaintenanceModeIoPriorityPolicy.IsUnknown() && !plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.IsUnknown() && !plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.IsUnknown(){
		protectedMaintenance := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.ProtectedMaintenanceModeIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.String(),
			BwLimitPerDeviceInKbps: plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.String(),
		}
		err := pd.SetProtectedMaintenanceModeIoPriorityPolicy(sp,protectedMaintenance)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set protected maintenance mode Io priority policy  to %s , num of concurrent IOs oer device to %s and bandwidth limit to %s", plan.ProtectedMaintenanceModeIoPriorityPolicy.String(),plan.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice.String(),plan.ProtectedMaintenanceModeBwLimitPerDeviceInKbps.String()),
				err.Error(),
			)
		}
	}

	if !plan.RebalanceEnabled.IsUnknown() && !plan.RebalanceEnabled.IsNull(){
		err := pd.SetRebalanceEnabled(sp,plan.RebalanceEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalnace enabled to %s", plan.RebalanceEnabled.String()),
				err.Error(),
			)
		}
	}

	if !plan.RebalanceIoPriorityPolicy.IsUnknown() && plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown(){
		rebalanceIoPriorityPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.RebalanceIoPriorityPolicy.ValueString(),
		}
		err := pd.SetRebalanceIoPriorityPolicy(sp,rebalanceIoPriorityPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalance Io priority policy  to %s", plan.RebalanceIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	} else if !plan.RebalanceIoPriorityPolicy.IsUnknown() && !plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown(){
		rebalanceIoPriorityPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.RebalanceIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.RebalanceNumOfConcurrentIosPerDevice.String(),
		}
		err := pd.SetRebalanceIoPriorityPolicy(sp,rebalanceIoPriorityPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalance Io priority policy  to %s and num of concurrent IOs oer device to %s", plan.RebalanceIoPriorityPolicy.String(),plan.RebalanceNumOfConcurrentIosPerDevice.String()),
				err.Error(),
			)
		}
	} else if !plan.RebalanceIoPriorityPolicy.IsUnknown() && !plan.RebalanceNumOfConcurrentIosPerDevice.IsUnknown() && !plan.RebalanceBwLimitPerDeviceInKbps.IsUnknown(){
		rebalanceIoPriorityPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.RebalanceIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.RebalanceNumOfConcurrentIosPerDevice.String(),
			BwLimitPerDeviceInKbps: plan.RebalanceBwLimitPerDeviceInKbps.String(),
		}
		err := pd.SetRebalanceIoPriorityPolicy(sp,rebalanceIoPriorityPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebalance Io priority policy  to %s , num of concurrent IOs oer device to %s and bandwidth limit to %s", plan.RebalanceIoPriorityPolicy.String(),plan.RebalanceNumOfConcurrentIosPerDevice.String(),plan.RebalanceBwLimitPerDeviceInKbps.String()),
				err.Error(),
			)
		}
	}

	if !plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown(){
		vtreeMigrationPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.VtreeMigrationIoPriorityPolicy.ValueString(),
		}
		err := pd.SetVTreeMigrationIOPriorityPolicy(sp,vtreeMigrationPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set Vtree migration Io priority policy  to %s", plan.VtreeMigrationIoPriorityPolicy.String()),
				err.Error(),
			)
		}
	} else if !plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && !plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown(){
		vtreeMigrationPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.VtreeMigrationIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.VtreeMigrationNumOfConcurrentIosPerDevice.String(),
		}
		err := pd.SetVTreeMigrationIOPriorityPolicy(sp,vtreeMigrationPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set Vtree migration Io priority policy  to %s and num of concurrent IOs oer device to %s", plan.VtreeMigrationIoPriorityPolicy.String(),plan.VtreeMigrationNumOfConcurrentIosPerDevice.String()),
				err.Error(),
			)
		}
	} else if !plan.VtreeMigrationIoPriorityPolicy.IsUnknown() && !plan.VtreeMigrationNumOfConcurrentIosPerDevice.IsUnknown() && !plan.VtreeMigrationBwLimitPerDeviceInKbps.IsUnknown(){
		vtreeMigrationPolicy := &scaleiotypes.ProtectedMaintenanceModeParam{
			Policy:      plan.VtreeMigrationIoPriorityPolicy.ValueString(),
			NumOfConcurrentIosPerDevice : plan.VtreeMigrationNumOfConcurrentIosPerDevice.String(),
			BwLimitPerDeviceInKbps: plan.VtreeMigrationBwLimitPerDeviceInKbps.String(),
		}
		err := pd.SetVTreeMigrationIOPriorityPolicy(sp,vtreeMigrationPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set Vtree migration Io priority policy  to %s , num of concurrent IOs oer device to %s and bandwidth limit to %s", plan.VtreeMigrationIoPriorityPolicy.String(),plan.VtreeMigrationNumOfConcurrentIosPerDevice.String(),plan.VtreeMigrationBwLimitPerDeviceInKbps.String()),
				err.Error(),
			)
		}
	}

	if !plan.RebuildEnabled.IsUnknown() && !plan.RebuildEnabled.IsNull(){
		err := pd.SetRebuildEnabled(sp,plan.RebuildEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild enabled to %s", plan.RebalanceEnabled.String()),
				err.Error(),
			)
		}
	}

	if !plan.RebuildRebalanceParallelism.IsUnknown() && !plan.RebuildRebalanceParallelism.IsNull(){
		err := pd.SetRebuildRebalanceParallelismParam(sp,plan.RebuildRebalanceParallelism.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild rebalance parallelism to %s", plan.RebuildRebalanceParallelism.String()),
				err.Error(),
			)
		}
	}

	if !plan.RebuildRebalanceParallelism.IsUnknown() && !plan.RebuildRebalanceParallelism.IsNull(){
		err := pd.SetRebuildRebalanceParallelismParam(sp,plan.RebuildRebalanceParallelism.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild rebalance parallelism to %s", plan.RebuildRebalanceParallelism.String()),
				err.Error(),
			)
		}
	}

	if plan.Fragmentation.String() == "true" {
		err := pd.EnableFragmentation(sp)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set fragmentation to %s", plan.Fragmentation.String()),
				err.Error(),
			)
		}
	} else if plan.Fragmentation.String() == "false" {
		err := pd.DisableFragmentation(sp)
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
	state.ZeroPaddingEnabled= types.BoolValue(storagepool.ZeroPaddingEnabled)
	state.ReplicationJournalCapacity= types.Int64Value(int64(storagepool.ReplicationCapacityMaxRatio))
	state.CapacityAlertHighThreshold= types.Int64Value(int64(storagepool.CapacityAlertHighThreshold))
	state.CapacityAlertCriticalThreshold= types.Int64Value(int64(storagepool.CapacityAlertCriticalThreshold))
	state.ProtectedMaintenanceModeIoPriorityPolicy=types.StringValue(storagepool.ProtectedMaintenanceModeIoPriorityPolicy)
	state.ProtectedMaintenanceModeNumOfConcurrentIosPerDevice=types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice))
	state.ProtectedMaintenanceModeBwLimitPerDeviceInKbps=types.Int64Value(int64(storagepool.ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps))
	state.RebalanceEnabled=types.BoolValue(storagepool.RebalanceEnabled)
	state.RebalanceIoPriorityPolicy=types.StringValue(storagepool.RebalanceioPriorityPolicy)
	state.RebalanceNumOfConcurrentIosPerDevice=types.Int64Value(int64(storagepool.RebalanceioPriorityNumOfConcurrentIosPerDevice))
	state.RebalanceBwLimitPerDeviceInKbps=types.Int64Value(int64(storagepool.RebalanceioPriorityBwLimitPerDeviceInKbps))
	state.VtreeMigrationIoPriorityPolicy=types.StringValue(storagepool.VtreeMigrationIoPriorityPolicy)
	state.VtreeMigrationNumOfConcurrentIosPerDevice=types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice))
	state.VtreeMigrationBwLimitPerDeviceInKbps=types.Int64Value(int64(storagepool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps))
	state.SparePercentage=types.Int64Value(int64(storagepool.SparePercentage))
	state.RmCacheWriteHandlingMode=types.StringValue(storagepool.RmCacheWriteHandlingMode)
	state.RebuildEnabled=types.BoolValue(storagepool.RebuildEnabled)
	state.RebuildRebalanceParallelism=types.Int64Value(int64(storagepool.NumofParallelRebuildRebalanceJobsPerDevice))
	state.Fragmentation=types.BoolValue(storagepool.FragmentationEnabled)
	return state
}
