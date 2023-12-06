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

package provider

import (
	"context"
	"fmt"
	"terraform-provider-powerflex/powerflex/helper"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &faultSetResource{}
	_ resource.ResourceWithConfigure   = &faultSetResource{}
	_ resource.ResourceWithImportState = &faultSetResource{}
)

// NewFaultSetResource - function to return resource interface
func NewFaultSetResource() resource.Resource {
	return &faultSetResource{}
}

type faultSetResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *faultSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fault_set"
}

func (r *faultSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FaultSetResourceSchema
}

func (r *faultSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	r.client = req.ProviderData.(*powerflexProvider).client
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}
	r.system = system
}

// Function used to Create fault set Resource
func (r *faultSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create fault set")
	// Retrieve values from plan
	var plan models.FaultSetResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := helper.GetNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	payload := &scaleiotypes.FaultSetParam{
		Name:               plan.Name.ValueString(),
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
	}

	// create the fault set
	fsID, err := pd.CreateFaultSet(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating fault set",
			"Could not create fault set, unexpected error: "+err.Error(),
		)
		return
	}
	response, err := r.system.GetFaultSetByID(fsID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting fault set after creation",
			"Could not get fault set, unexpected error: "+err.Error(),
		)
		return
	}

	state := helper.UpdateFaultSetState(response, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Fault set Resource
func (r *faultSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read Storagepool")
	// Get current state
	var state models.FaultSetResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fs, err := r.system.GetFaultSetByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get fault set by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}
	faultsetState := helper.UpdateFaultSetState(fs, state)
	diags = resp.State.Set(ctx, faultsetState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update fault set Resource
func (r *faultSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.FaultSetResourceModel
	//var err1 error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state models.FaultSetResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ProtectionDomainID.ValueString() != state.ProtectionDomainID.ValueString() {
		resp.Diagnostics.AddError(
			"Protection Domain ID cannot be updated",
			"Protection Domain ID cannot be updated")
		return
	}

	pd, err := helper.GetNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		err := pd.ModifyFaultSetName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating name of fault set", err.Error(),
			)
		}
	}
	response, err := r.system.GetFaultSetByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting fault set", err.Error(),
		)
		return
	}

	state1 := helper.UpdateFaultSetState(response, state)
	diags = resp.State.Set(ctx, state1)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Fault Set Resource
func (r *faultSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.FaultSetResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := helper.GetNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	err = pd.DeleteFaultSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting fault set",
			"Couldn't Delete fault set "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for fault set Resource
func (r *faultSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
