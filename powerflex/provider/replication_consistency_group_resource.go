/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &replicationConsistencyGroupResource{}
	_ resource.ResourceWithConfigure   = &replicationConsistencyGroupResource{}
	_ resource.ResourceWithImportState = &replicationConsistencyGroupResource{}
)

// ReplicationConsistencyGroupResource - function to return resource interface
func ReplicationConsistencyGroupResource() resource.Resource {
	return &replicationConsistencyGroupResource{}
}

// replicationConsistencyGroupResource - struct to define ReplicationConsistencyGroup resource
type replicationConsistencyGroupResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_consistency_group"
}

// Schema - function to return Schema for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ReplicationConsistencyGroupReourceSchema
}

// Configure - function to return Configuration for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {

		r.client = req.ProviderData.(*powerflexProvider).client
	}

	if req.ProviderData.(*powerflexProvider).gatewayClient != nil {

		r.gatewayClient = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// Create - function to Create for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")
	var plan models.ReplicationConsistancyGroupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Some of the fields can only be updated after create
	if plan.CurrConsistMode.ValueString() != "Consistent" ||
		plan.FreezeState.ValueString() != "Unfrozen" ||
		plan.PauseMode.ValueString() != "Resume" ||
		plan.TargetVolumeAccessMode.ValueString() != "NoAccess" ||
		plan.LocalActivityState.ValueString() != "Active" {
		resp.Diagnostics.AddError(
			"local_activity_state, target_volume_access_mode, pause_mode, freeze_state, and curr_consist_mode "+
				"are not able to be modified from default values while creating a new replication consistency group."+
				"These fields can be modifed on an imported or already created replication consistency group",
			"Please remove these fields and try again to create the replication consistency group",
		)
		return
	}

	// Create
	createID, createErr := helper.CreateReplicationConsistencyGroup(r.client, plan)
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating replication consistency group",
			createErr.Error(),
		)
		return
	}

	// Read
	rcg, err := helper.GetSpecificReplicationConsistencyGroup(r.client, createID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication consistency group",
			err.Error(),
		)
		return
	}

	// Map
	updatedState := helper.MapReplicationConsistancyGroupsResourceState(*rcg, plan)
	diagsState := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsState...)
}

// Read - function to Read for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.ReplicationConsistancyGroupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rcg, err := helper.GetSpecificReplicationConsistencyGroup(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication consistency group",
			err.Error(),
		)
		return
	}

	updatedState := helper.MapReplicationConsistancyGroupsResourceState(*rcg, state)
	diagsState := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsState...)
}

// Update - function to Update for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")

	var state models.ReplicationConsistancyGroupModel
	var plan models.ReplicationConsistancyGroupModel

	diagsState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagsState...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsPlan := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diagsPlan...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the state values are equal to null then just read plan and set because it could be an import
	if state.Name.IsNull() || state.ProtectionDomainID.IsNull() || state.RemoteProtectionDomainID.IsNull() || state.DestinationSystemID.IsNull() {
		rcg, err := helper.GetSpecificReplicationConsistencyGroup(r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading replication consistency group",
				err.Error(),
			)
			return
		}
		updatedState := helper.MapReplicationConsistancyGroupsResourceState(*rcg, plan)
		diagsUpdate := resp.State.Set(ctx, updatedState)
		resp.Diagnostics.Append(diagsUpdate...)
		return
	}

	// Check for values which are not able to be updated, if update is attempted return an error
	if (state.ProtectionDomainID.ValueString() != plan.ProtectionDomainID.ValueString()) ||
		(state.RemoteProtectionDomainID.ValueString() != plan.RemoteProtectionDomainID.ValueString()) ||
		(state.DestinationSystemID.ValueString() != plan.DestinationSystemID.ValueString()) {
		resp.Diagnostics.AddError(
			"protection_domain_id, remote_protection_domain_id, and destination_system_id cannot be updated",
			"protection_domain_id, remote_protection_domain_id, and destination_system_id cannot be updated")
		return
	}

	// Do the update
	updateError := helper.RCGUpdates(r.client, state, plan)

	if updateError != nil {
		resp.Diagnostics.AddError(
			"Error updating replication consistency group",
			updateError.Error(),
		)
		return
	}

	rcg, err := helper.GetSpecificReplicationConsistencyGroup(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication consistency group",
			err.Error(),
		)
		return
	}

	updatedState := helper.MapReplicationConsistancyGroupsResourceState(*rcg, plan)
	diagsUpdate := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsUpdate...)
}

// Delete - function to Delete for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.ReplicationConsistancyGroupModel
	diagState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagState...)
	if resp.Diagnostics.HasError() {
		return
	}

	rcg, err := helper.GetSpecificReplicationConsistencyGroup(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication consistency group",
			err.Error(),
		)
		return
	}

	rcgC := goscaleio.NewReplicationConsistencyGroup(r.client)
	rcgC.ReplicationConsistencyGroup = rcg

	errRemove := rcgC.RemoveReplicationConsistencyGroup(false)
	if errRemove != nil {
		resp.Diagnostics.AddError(
			"Error deleting replication consistency group",
			errRemove.Error(),
		)
		return
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// ImportState - function to ImportState for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
