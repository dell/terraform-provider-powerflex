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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &serviceResource{}
	_ resource.ResourceWithConfigure   = &serviceResource{}
	_ resource.ResourceWithImportState = &serviceResource{}
)

// ServiceResource - function to return resource interface
func ServiceResource() resource.Resource {
	return &serviceResource{}
}

// serviceResource - struct to define sdc resource
type serviceResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for SDC resource.
func (r *serviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Schema - function to return Schema for SDC resource.
func (r *serviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ServiceReourceSchema
}

// Configure - function to return Configuration for SDC resource.
func (r *serviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create - function to Create for SDC resource.
func (r *serviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")

	var plan models.ServiceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentResponse, err := r.gatewayClient.DeployService(plan.DeploymentName.ValueString(), plan.DeploymentDescription.ValueString(), plan.TemplateId.ValueString(), plan.FirmwareId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in deploying service",
			err.Error(),
		)
		return
	}

	deploymentID := deploymentResponse.ID

	couterForStopExecution := 0

	deadLineCount := 60 / 5 //TODO

	for couterForStopExecution <= deadLineCount {

		time.Sleep(5 * time.Minute)

		tflog.Info(ctx, "Service Details updated to state file successfully")

		deploymentResponse, err = r.gatewayClient.GetServiceDetailsByID(deploymentID, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting service details",
				err.Error(),
			)
			return
		}

		if deploymentResponse.Status == "complete" {

			data, dgs := helper.UpdateServiceState(deploymentResponse, plan)
			resp.Diagnostics.Append(dgs...)

			diags = resp.State.Set(ctx, data)
			resp.Diagnostics.Append(diags...)

			tflog.Info(ctx, "Service Details updated to state file successfully")

			return
		}

		couterForStopExecution++

	}

	resp.Diagnostics.AddError(
		"Timed Out For Getting the Deployemnt Status",
		err.Error(),
	)

}

// Read - function to Read for SDC resource.
func (r *serviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.ServiceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//For handling the import case
	if state.ID.ValueString() != "" && state.ID.ValueString() != "placeholder" {
		deploymentResponse, err := r.gatewayClient.GetServiceDetailsByID(state.ID.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting service details",
				err.Error(),
			)
			return
		}

		if deploymentResponse.Status == "complete" {

			data, dgs := helper.UpdateServiceState(deploymentResponse, state)
			resp.Diagnostics.Append(dgs...)

			diags = resp.State.Set(ctx, data)
			resp.Diagnostics.Append(diags...)

			tflog.Info(ctx, "Service Details updated to state file successfully")

			return
		}

	} else {
		resp.Diagnostics.AddError("[Read] Please provide valid Service ID", "Please provide valid Service ID")

		return
	}
}

// Update - function to Update for Service resource.
func (r *serviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	// Retrieve values from plan
	var plan models.ServiceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.ServiceResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete - function to Delete for SDC resource.
func (r *serviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.ServiceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.gatewayClient.DeleteService(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in Deleting service details",
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)

}

// ImportState - function to ImportState for SDC resource.
func (r *serviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
