/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	_ resource.Resource                = &resourceCredentialResource{}
	_ resource.ResourceWithConfigure   = &resourceCredentialResource{}
	_ resource.ResourceWithImportState = &resourceCredentialResource{}
)

// ResourceCredentialResource - function to return resource interface
func ResourceCredentialResource() resource.Resource {
	return &resourceCredentialResource{}
}

// resourceCredentialResource - struct to define ResourceCredential resource
type resourceCredentialResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

// Metadata - function to return metadata for ResourceCredential resource.
func (r *resourceCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_credential"
}

// Schema - function to return Schema for ResourceCredential resource.
func (r *resourceCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceCredentialResourceSchema
}

// Configure - function to return Configuration for ResourceCredential resource.
func (r *resourceCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {

		r.client = req.ProviderData.(*powerflexProvider).client
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

// Create - function to Create for ResourceCredential resource.
func (r *resourceCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan models.ResourceCredentialResourceModel
	// Get the plan
	diagsPlan := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diagsPlan...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Validate inputs
	resp.Diagnostics.Append(helper.ValidateResourceCredentialResource(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Creates the resource credential
	rc, errCreate := helper.CreateResourceCredential(ctx, r.system, plan)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error Creating Resource Credential",
			errCreate.Error(),
		)
		return
	}

	updateState := helper.MapResourceCredentialResource(rc.Credential, plan)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)

	tflog.Debug(ctx, "[POWERFLEX] Create")
}

// Read - function to Read for ResourceCredential resource.
func (r *resourceCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ResourceCredentialResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rc, err := r.system.GetResourceCredential(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Credential",
			err.Error(),
		)
		return
	}

	// Map RC and return
	updateState := helper.MapResourceCredentialResource(rc.Credential, state)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)
	tflog.Debug(ctx, "[POWERFLEX] Read")
}

// Update - function to Update for ResourceCredential resource.
func (r *resourceCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.ResourceCredentialResourceModel
	var plan models.ResourceCredentialResourceModel
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

	// If type state is null it means we are doing an import and can continue.
	// Otherwise type is not a modifiable field
	if !state.Type.IsNull() && state.Type.ValueString() != plan.Type.ValueString() {
		resp.Diagnostics.AddError(
			"Error Modifing Resource Credential",
			"Type cannot be modified",
		)
		return
	}

	// Validate inputs
	resp.Diagnostics.Append(helper.ValidateResourceCredentialResource(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Updates the resource credential
	rc, errModify := helper.ModifyResourceCredential(ctx, r.system, plan, state.ID.ValueString())
	if errModify != nil {
		resp.Diagnostics.AddError(
			"Error Modifing Resource Credential",
			errModify.Error(),
		)
		return
	}

	// Map RC and return
	updateState := helper.MapResourceCredentialResource(rc.Credential, plan)
	diagsState := resp.State.Set(ctx, updateState)
	resp.Diagnostics.Append(diagsState...)
	tflog.Debug(ctx, "[POWERFLEX] Update")
}

// Delete - function to Delete for ResourceCredential resource.
func (r *resourceCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ResourceCredentialResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.system.DeleteResourceCredential(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Resource Credential",
			err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "[POWERFLEX] Delete")
}

// ImportState - function to ImportState for ResourceCredential resource.
func (r *resourceCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
