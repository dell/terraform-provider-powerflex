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
	_ resource.Resource                = &peerSystemResource{}
	_ resource.ResourceWithConfigure   = &peerSystemResource{}
	_ resource.ResourceWithImportState = &peerSystemResource{}
)

// PeerSystemResource - function to return resource interface
func PeerSystemResource() resource.Resource {
	return &peerSystemResource{}
}

// peerSystemResource - struct to define PeerSystem resource
type peerSystemResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for PeerSystem resource.
func (r *peerSystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_peer_system"
}

// Schema - function to return Schema for PeerSystem resource.
func (r *peerSystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PeerSystemReourceSchema
}

// Configure - function to return Configuration for PeerSystem resource.
func (r *peerSystemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create - function to Create for PeerSystem resource.
func (r *peerSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")
	var plan models.PeerMdmResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If certificate add_certificate is set to true, get the root certificate and add it to the source mdm trust store
	if plan.AddCertificate.ValueBool() {
		err := helper.AddCertificate(ctx, r.client, plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding certificate to trust store from destination primary mdm to source primary mdm",
				err.Error(),
			)
			return
		}
	}

	// Create PeerSystem
	id, err := helper.CreatePeerSystem(r.client, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating peer system",
			err.Error(),
		)
		return
	}

	peerSystem, err := helper.GetPeerSystem(r.client, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading peer system",
			err.Error(),
		)
		return
	}
	// Map
	updatedState := helper.MapPeerSystemResourceState(*peerSystem, plan)
	diagsState := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsState...)
}

// Read - function to Read for PeerSystem resource.
func (r *peerSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.PeerMdmResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	peerSystem, err := helper.GetPeerSystem(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading peer system",
			err.Error(),
		)
		return
	}
	updatedState := helper.MapPeerSystemResourceState(*peerSystem, state)
	diagsState := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsState...)
}

// Update - function to Update for PeerSystem resource.
func (r *peerSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	var state models.PeerMdmResourceModel
	var plan models.PeerMdmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagState...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If any of the state values are equal to null then just read plan and set because it could be an import
	if state.Name.IsNull() || state.PeerSystemID.IsNull() || len(state.IPList) == 0 || state.Port.IsNull() || state.PerfProfile.IsNull() {
		peerSystem, err := helper.GetPeerSystem(r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading peer system",
				err.Error(),
			)
			return
		}
		updatedState := helper.MapPeerSystemResourceState(*peerSystem, plan)
		diagsState := resp.State.Set(ctx, updatedState)
		resp.Diagnostics.Append(diagsState...)
		return
	}

	if plan.PeerSystemID.ValueString() != state.PeerSystemID.ValueString() {
		resp.Diagnostics.AddError(
			"peer_system_id cannot be updated",
			"peer_system_id cannot be updated")
		return
	}

	// Do the update
	updateError := helper.PeerSystemUpdate(r.client, state, plan)

	if updateError != nil {
		resp.Diagnostics.AddError(
			"Error updating peer system",
			updateError.Error(),
		)
		return
	}

	peerSystem, err := helper.GetPeerSystem(r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading peer system",
			err.Error(),
		)
		return
	}
	updatedState := helper.MapPeerSystemResourceState(*peerSystem, plan)
	diagsState := resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diagsState...)
}

// Delete - function to Delete for PeerSystem resource.
func (r *peerSystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.PeerMdmResourceModel
	diagState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagState...)
	if resp.Diagnostics.HasError() {
		return
	}

	errRemove := r.client.RemovePeerMdm(state.ID.ValueString())
	if errRemove != nil {
		resp.Diagnostics.AddError(
			"Error removing peer system",
			errRemove.Error(),
		)
		return
	}

}

// ImportState - function to ImportState for PeerSystem resource.
func (r *peerSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
