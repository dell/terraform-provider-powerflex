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
	_ resource.Resource                = &resourceGroupResource{}
	_ resource.ResourceWithConfigure   = &resourceGroupResource{}
	_ resource.ResourceWithImportState = &resourceGroupResource{}
)

// ResourceGroupResource - function to return resource interface
func ResourceGroupResource() resource.Resource {
	return &resourceGroupResource{}
}

// ResourceGroupResource - struct to define ResourceGroup resource
type resourceGroupResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for ResourceGroup resource.
func (r *resourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

// Schema - function to return Schema for ResourceGroup resource.
func (r *resourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceGroupReourceSchema
}

// Configure - function to return Configuration for ResourceGroup resource.
func (r *resourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create - function to Create for ResourceGroup resource.
func (r *resourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")

	var plan models.ResourceGroupResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CloneFromHost.ValueString() != "" {
		resp.Diagnostics.AddError("No need to pass clone_from_host during the deployment of ResourceGroup", "please validate your inputs")
		return
	}

	deploymentResponse, err := r.gatewayClient.DeployService(plan.DeploymentName.ValueString(), plan.DeploymentDescription.ValueString(), plan.TemplateID.ValueString(), plan.FirmwareID.ValueString(), plan.Nodes.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in deploying ResourceGroup",
			err.Error(),
		)
		return
	}

	deploymentResponse, diag := helper.HandleResourceGroupDeployment(ctx, deploymentResponse, plan, r.gatewayClient)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	data, dgs := helper.UpdateResourceGroupState(deploymentResponse, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "ResourceGroup Details updated to state file successfully")

}

// Read - function to Read for ResourceGroup resource.
func (r *resourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.ResourceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//For handling the import case
	if state.ID.ValueString() != "" {
		deploymentResponse, err := r.gatewayClient.GetServiceDetailsByID(state.ID.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting ResourceGroup details",
				err.Error(),
			)
			return
		}

		if deploymentResponse.Status == "complete" {

			data, dgs := helper.UpdateResourceGroupState(deploymentResponse, state)
			resp.Diagnostics.Append(dgs...)

			diags = resp.State.Set(ctx, data)
			resp.Diagnostics.Append(diags...)

			tflog.Info(ctx, "ResourceGroup Details updated to state file successfully")

			return
		}

	} else {
		resp.Diagnostics.AddError("[Read] Please provide valid ResourceGroup ID", "Please provide valid ResourceGroup ID")

		return
	}
}

// Update - function to Update for ResourceGroup resource.
func (r *resourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	// Retrieve values from plan
	var plan models.ResourceGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.ResourceGroupResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Nodes.ValueInt64() < state.Nodes.ValueInt64() {
		resp.Diagnostics.AddError("Removing node(s) is not supported", "please validate your inputs")
		return
	}

	if plan.Nodes.String() != state.Nodes.String() && plan.CloneFromHost.ValueString() == "" {
		resp.Diagnostics.AddError("Please provide clone_from_host for adding the resource", "please validate your inputs")
		return
	}

	if plan.TemplateID.String() != state.TemplateID.String() {
		resp.Diagnostics.AddError("Changing of template_id is not supported", "please validate your inputs")
		return
	}

	if plan.FirmwareID.String() != state.FirmwareID.String() {
		resp.Diagnostics.AddError("Changing of firmware_id is not supported", "please validate your inputs")
		return
	}

	deploymentResponse, err := r.gatewayClient.UpdateService(state.ID.ValueString(), plan.DeploymentName.ValueString(), plan.DeploymentDescription.ValueString(), plan.Nodes.String(), plan.CloneFromHost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in deploying service",
			err.Error(),
		)
		return
	}

	if deploymentResponse.Status == "complete" {
		data, dgs := helper.UpdateResourceGroupState(deploymentResponse, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "ResourceGroup Details updated to state file successfully")

		return

	}

	deploymentResponse, diag := helper.HandleResourceGroupDeployment(ctx, deploymentResponse, plan, r.gatewayClient)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	data, dgs := helper.UpdateResourceGroupState(deploymentResponse, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "ResourceGroup Details updated to state file successfully")

}

// Delete - function to Delete for ResourceGroup resource.
func (r *resourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.ResourceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.gatewayClient.DeleteService(state.ID.ValueString(), state.ServersInInventory.ValueString(), state.ServersManagedState.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in Deleting ResourceGroup details",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)

}

// ImportState - function to ImportState for ResourceGroup resource.
func (r *resourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
