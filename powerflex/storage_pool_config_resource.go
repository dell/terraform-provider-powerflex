package powerflex

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &storagepoolConfigResource{}
	_ resource.ResourceWithConfigure   = &storagepoolConfigResource{}
	_ resource.ResourceWithImportState = &storagepoolConfigResource{}
)

// StoragepoolResource - function to return resource interface
func StoragepoolConfigResource() resource.Resource {
	return &storagepoolConfigResource{}
}

type storagepoolConfigResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

type storagepoolConfigResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	ProtectionDomainID         types.String   `tfsdk:"protection_domain_id"`
	RebuildEnabled             types.Bool     `tfsdk:"rebuild_enabled"`
	ReplicationJournalCapacity types.Int64    `tfsdk:"replication_journal_capacity"`
	OriginalConfig             originalConfig `tfsdk:"original_config"`
}

type originalConfig struct {
	RebuildEnabled             types.Bool  `tfsdk:"rebuild_enabled"`
	ReplicationJournalCapacity types.Int64 `tfsdk:"replication_journal_capacity"`
}

func (r *storagepoolConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_pool_config"
}

func (r *storagepoolConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StoragepoolReourceSchema
}

func (r *storagepoolConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*goscaleio.Client)
	system, err := getFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	r.system = system
}

// Function used to Create Storagepool Resource
func (r *storagepoolConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create storagepool")
	// Retrieve values from plan
	var plan, state storagepoolConfigResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	sp, err := pd.FindStoragePool(plan.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting Storagepool", err.Error(),
		)
		return
	}

	state.OriginalConfig = originalConfig{
		ReplicationJournalCapacity: types.Int64Value(int64(sp.ReplicationCapacityMaxRatio)),
		RebuildEnabled:             types.BoolValue(sp.RebuildEnabled),
	}

	// set rebuild enabled
	if !plan.RebuildEnabled.IsUnknown() && !plan.RebuildEnabled.IsNull() {
		err := pd.SetRebuildEnabled(sp.ID, plan.RebuildEnabled.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set rebuild enabled to %s", plan.RebuildEnabled.String()),
				err.Error(),
			)
		}
		state.RebuildEnabled = plan.RebuildEnabled
	}

	// set the replication journal capacity
	if !plan.ReplicationJournalCapacity.IsUnknown() && !plan.ReplicationJournalCapacity.IsNull() {
		err := pd.SetReplicationJournalCapacity(sp.ID, plan.ReplicationJournalCapacity.String())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set replication Journal capacity to %s", plan.ReplicationJournalCapacity.String()),
				err.Error(),
			)
		}
		state.ReplicationJournalCapacity = plan.ReplicationJournalCapacity
	}
	state.ID = plan.ID
	state.ProtectionDomainID = plan.ProtectionDomainID
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Storagepool Resource
func (r *storagepoolConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read Storagepool")
	// Get current state
	var state storagepoolConfigResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spr, err := r.system.GetStoragePoolByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get storagepool by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}

	state.ReplicationJournalCapacity = types.Int64Value(int64(spr.ReplicationCapacityMaxRatio))
	state.RebuildEnabled = types.BoolValue(spr.RebuildEnabled)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update Storagepool Resource
func (r *storagepoolConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update Storagepool")
	// Retrieve values from plan
	var plan, state storagepoolConfigResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	sp, err := pd.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting Storagepool", err.Error(),
		)
		return
	}

	if !plan.ReplicationJournalCapacity.IsUnknown() &&
		!state.ReplicationJournalCapacity.Equal(plan.ReplicationJournalCapacity) {
		errReplicationJournalCapacity := pd.SetReplicationJournalCapacity(sp.ID, strconv.FormatInt(plan.ReplicationJournalCapacity.ValueInt64(), 10))
		if errReplicationJournalCapacity != nil {
			resp.Diagnostics.AddError(
				"Error while updating ReplicationJournalCapacity of Storagepool", errReplicationJournalCapacity.Error(),
			)
		}
		state.ReplicationJournalCapacity = plan.ReplicationJournalCapacity
	}

	if !plan.RebuildEnabled.IsUnknown() &&
		!state.RebuildEnabled.Equal(plan.RebuildEnabled) {
		errRebuildEnabled := pd.SetRebuildEnabled(sp.ID, plan.RebuildEnabled.String())
		if errRebuildEnabled != nil {
			resp.Diagnostics.AddError(
				"Error while updating RebuildEnabled of Storagepool", errRebuildEnabled.Error(),
			)
		}
		state.RebuildEnabled = plan.RebuildEnabled
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Storagepool Resource
func (r *storagepoolConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete Storagepool")
	// Retrieve values from state
	var state storagepoolConfigResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := getNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(),"", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}
	sp, err := pd.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting Storagepool", err.Error(),
		)
		return
	}

	if !state.OriginalConfig.ReplicationJournalCapacity.IsUnknown() &&
		!state.ReplicationJournalCapacity.Equal(state.OriginalConfig.ReplicationJournalCapacity) {
		errReplicationJournalCapacity := pd.SetReplicationJournalCapacity(sp.ID, strconv.FormatInt(state.OriginalConfig.ReplicationJournalCapacity.ValueInt64(), 10))
		if errReplicationJournalCapacity != nil {
			resp.Diagnostics.AddError(
				"Error while updating ReplicationJournalCapacity of Storagepool", errReplicationJournalCapacity.Error(),
			)
		}
		state.ReplicationJournalCapacity = state.OriginalConfig.ReplicationJournalCapacity
	}

	if !state.OriginalConfig.RebuildEnabled.IsUnknown() &&
		!state.RebuildEnabled.Equal(state.OriginalConfig.RebuildEnabled) {
		errRebuildEnabled := pd.SetRebuildEnabled(sp.ID, state.OriginalConfig.RebuildEnabled.String())
		if errRebuildEnabled != nil {
			resp.Diagnostics.AddError(
				"Error while updating RebuildEnabled of Storagepool", errRebuildEnabled.Error(),
			)
		}
		state.RebuildEnabled = state.OriginalConfig.RebuildEnabled
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for Storagepool Resource
func (r *storagepoolConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
