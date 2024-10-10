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
	"time"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &replicationPairResource{}
	_ resource.ResourceWithConfigure   = &replicationPairResource{}
	_ resource.ResourceWithImportState = &replicationPairResource{}
)

// ReplicationPairResource - function to return resource interface
func ReplicationPairResource() resource.Resource {
	return &replicationPairResource{}
}

// replicationPairResource - struct to define ReplicationPair resource
type replicationPairResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for ReplicationPair resource.
func (r *replicationPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_pair"
}

// Schema - function to return Schema for ReplicationPair resource.
func (r *replicationPairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ReplicationPairReourceSchema
}

// Configure - function to return Configuration for ReplicationPair resource.
func (r *replicationPairResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create - function to Create for ReplicationPair resource.
func (r *replicationPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")
	var plan models.ReplicationPairResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createID, createErr := helper.CreateReplicationPair(r.client, plan)
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating replication pair",
			createErr.Error(),
		)
		return
	}

	// If Pause Copy is true, then pause after the create
	if plan.PauseCopy.ValueBool() {
		// Sleep for 4 seconds after create before pausing
		time.Sleep(4 * time.Second)
		// Pause replication pair
		_, err := helper.PauseReplicationPair(r.client, createID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error pausing replication pair after create",
				err.Error(),
			)
			return
		}
	}

	rp, err := helper.GetSpecificReplicationPair(r.client, createID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication pair",
			err.Error(),
		)
		return
	}
	updateState := helper.MapReplicationPairState(*rp, plan)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)

	tflog.Info(ctx, "replication pair details updated to state file successfully")

}

// Read - function to Read for ReplicationPair resource.
func (r *replicationPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ReplicationPairResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rp, err := helper.GetSpecificReplicationPair(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication pair",
			err.Error(),
		)
		return
	}

	updateState := helper.MapReplicationPairState(*rp, state)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)
	tflog.Debug(ctx, "[POWERFLEX] Read")
}

// Update - function to Update for ReplicationPair resource.
func (r *replicationPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	var state models.ReplicationPairResourceModel
	var plan models.ReplicationPairResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diagsPlan := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diagsPlan...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the state values are equal to null then just read plan and set because it could be an import
	if state.Name.IsNull() || state.SourceVolumeID.IsNull() || state.DestinationVolumeID.IsNull() || state.ReplicationConsistencyGroupID.IsNull() {
		rp, err := helper.GetSpecificReplicationPair(r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading replication pair after update",
				err.Error(),
			)
			return
		}
		updateState := helper.MapReplicationPairState(*rp, plan)
		diagsState := resp.State.Set(ctx, updateState)
		resp.Diagnostics.Append(diagsState...)
		return
	}

	if (state.Name.ValueString() != plan.Name.ValueString()) ||
		(state.SourceVolumeID.ValueString() != plan.SourceVolumeID.ValueString()) ||
		(state.DestinationVolumeID.ValueString() != plan.DestinationVolumeID.ValueString()) ||
		(state.ReplicationConsistencyGroupID.ValueString() != plan.ReplicationConsistencyGroupID.ValueString()) {
		resp.Diagnostics.AddError(
			"name, source_volume_id, replication_consistency_group_id, and destination_volume_id cannot be updated",
			"name, source_volume_id, source_volume_id, and destination_volume_id cannot be updated",
		)
		return
	}

	if plan.PauseCopy.ValueBool() {
		// Pause replication pair
		_, err := helper.PauseReplicationPair(r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error pausing replication pair, only avaiable during initial copy",
				err.Error(),
			)
			return
		}
	} else {
		// Resume replication pair
		_, err := helper.ResumeReplicationPair(r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error resuming replication pair, only avaiable during initial copy",
				err.Error(),
			)
			return
		}
	}

	rp, err := helper.GetSpecificReplicationPair(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication pair after update",
			err.Error(),
		)
		return
	}
	updateState := helper.MapReplicationPairState(*rp, plan)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)

}

// Delete - function to Delete for ReplicationPair resource.
func (r *replicationPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.ReplicationPairResourceModel
	diagState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagState...)
	if resp.Diagnostics.HasError() {
		return
	}

	pair := goscaleio.NewReplicationPair(r.client)
	pair.ReplicaitonPair.ID = state.ID.ValueString()

	_, err := pair.RemoveReplicationPair(false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting replication pair",
			err.Error(),
		)
		return
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState - function to ImportState for ReplicationPair resource.
func (r *replicationPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
