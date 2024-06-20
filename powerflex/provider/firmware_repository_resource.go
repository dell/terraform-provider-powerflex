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
	"time"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &firmwareRepositoryResource{}
	_ resource.ResourceWithConfigure   = &firmwareRepositoryResource{}
	_ resource.ResourceWithImportState = &firmwareRepositoryResource{}
)

// NewFirmwareRepositoryResource - function to return resource interface
func NewFirmwareRepositoryResource() resource.Resource {
	return &firmwareRepositoryResource{}
}

type firmwareRepositoryResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (r *firmwareRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firmware_repository"
}

func (r *firmwareRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FirmwareRepositoryResourceSchema
}

func (r *firmwareRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Function used to Create Firmware Repository Resource
func (r *firmwareRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Upload Firmware Repository")
	// Retrieve values from plan
	var endTime time.Time
	var plan models.FirmwareRepositoryResourceModel
	var frDetails *scaleiotypes.UploadComplianceTopologyDetails
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ucParam := &scaleiotypes.UploadComplianceParam{
		SourceLocation: plan.SourceLocation.ValueString(),
	}

	if !plan.Username.IsNull() && !plan.Password.IsNull() {
		ucParam.Username = plan.Username.ValueString()
		ucParam.Password = plan.Password.ValueString()
	}

	err := r.gatewayClient.TestConnection(ucParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Please provide valid credentials",
			err.Error(),
		)
		return
	}

	// uploading the compliance file
	fr, err := r.gatewayClient.UploadCompliance(ucParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Upload the compliance File",
			err.Error(),
		)
		return
	}
	startTime := time.Now()

	if !plan.Timeout.IsNull() {
		endTime = startTime.Add(time.Duration(plan.Timeout.ValueInt64()) * time.Minute)
	}
	// will loop until the file gets loaded successfully and status becomes available or until it times out
	for time.Now().Before(endTime) {
		time.Sleep(1 * time.Minute)
		frDetails, err = r.gatewayClient.GetUploadComplianceDetails(fr.ID, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not get the Firmware Repository Details",
				err.Error(),
			)
			return
		}
		if frDetails.State == "available" {
			break
		} else if frDetails.State == "needsApproval" {
			if plan.Approve.ValueBool() {
				err := r.gatewayClient.ApproveUnsignedFile(frDetails.ID)
				if err != nil {
					resp.Diagnostics.AddError(
						"Could not approve the compliance File",
						err.Error(),
					)
					return
				}
			} else {
				resp.Diagnostics.AddWarning(
					"The compliance file is unsigned",
					"The compliance file needs approval to proceed ahead.",
				)
				break
			}
		} else if frDetails.State == "errors" {
			resp.Diagnostics.AddError(
				"Could not Upload the compliance File",
				"error while uploading compliance file",
			)
			return
		} else if frDetails.State == "copying" {
			continue
		}

	}

	frDetails, err = r.gatewayClient.GetUploadComplianceDetails(fr.ID, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the Firmware Repository Details",
			err.Error(),
		)
		return
	}

	if frDetails.State == "copying" || frDetails.State == "pending" {
		resp.Diagnostics.AddError(
			"The Operation got timed Out",
			"The Operation got timed Out",
		)
		return
	}

	state := helper.UpdateFrimwareRepositoryState(frDetails, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Firmware Repository Resource
func (r *firmwareRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read Storagepool")
	// Get current state
	var state models.FirmwareRepositoryResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	frDetails, err := r.gatewayClient.GetUploadComplianceDetails(state.ID.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the Firmware Repository Details",
			err.Error(),
		)
		return
	}
	frState := helper.UpdateFrimwareRepositoryState(frDetails, state)
	diags = resp.State.Set(ctx, frState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update firmware repository Resource
func (r *firmwareRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.FirmwareRepositoryResourceModel
	var frDetails *scaleiotypes.UploadComplianceTopologyDetails
	var err error
	var endTime time.Time
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state models.FirmwareRepositoryResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.SourceLocation.ValueString() != state.SourceLocation.ValueString() {
		resp.Diagnostics.AddError(
			"Source Location cannot be updated",
			"Source Location cannot be updated")
		return
	}

	if plan.Username.ValueString() != state.Username.ValueString() {
		resp.Diagnostics.AddError(
			"Username cannot be updated",
			"Username cannot be updated. If the resource is imported then username is not required for approving the unsigned file.")
		return
	}

	if plan.Password.ValueString() != state.Password.ValueString() {
		resp.Diagnostics.AddError(
			"Password cannot be updated",
			"Password cannot be updated. If the resource is imported then password is not required for approving the unsigned file.")
		return
	}

	if !plan.Approve.ValueBool() && state.Approve.ValueBool() {
		resp.Diagnostics.AddError(
			"Approve cannot be set to false once it is set to true.",
			"Approve cannot be set to false once it is set to true.")
		return
	}

	startTime := time.Now()

	if plan.Timeout.ValueInt64() != state.Timeout.ValueInt64() {
		endTime = startTime.Add(time.Duration(plan.Timeout.ValueInt64()) * time.Minute)
	} else {
		endTime = startTime.Add(time.Duration(state.Timeout.ValueInt64()) * time.Minute)
	}
	if plan.Approve.ValueBool() != state.Approve.ValueBool() {
		for time.Now().Before(endTime) {
			frDetails, err = r.gatewayClient.GetUploadComplianceDetails(state.ID.ValueString(), true)
			if err != nil {
				resp.Diagnostics.AddError(
					"Could not get the Firmware Repository Details",
					err.Error(),
				)
				return
			}
			if frDetails.State == "available" {
				break
			} else if frDetails.State == "needsApproval" {
				if plan.Approve.ValueBool() {
					err2 := r.gatewayClient.ApproveUnsignedFile(frDetails.ID)
					if err2 != nil {
						resp.Diagnostics.AddError(
							"Could not Upload the compliance File",
							err2.Error(),
						)
						return
					}
				} else {
					break
				}
			} else if frDetails.State == "errors" {
				resp.Diagnostics.AddError(
					"Could not Upload the compliance File",
					"error while uploading compliance file",
				)
				return
			} else if frDetails.State == "copying" || frDetails.State == "pending" {
				continue
			}
			time.Sleep(1 * time.Minute)
		}
	} else {
		resp.Diagnostics.AddError(
			"Approve attribute needs to be updated",
			"Please modify the approve attribute if you want to approve your unsigned file and proceed with the download."+
				"if the resource is imported then please check the value of approve in the state file. If it is set to false and you want to approve your unsigned file then update the approve attribute to true and proceed with the download. If it is set to true then no need to update anything.",
		)
		return
	}

	if frDetails.State == "copying" || frDetails.State == "pending" {
		resp.Diagnostics.AddError(
			"The Operation got timed Out",
			"The Operation got timed Out",
		)
		return
	}
	state = helper.UpdateFrimwareRepositoryState(frDetails, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Firmware Repository Resource
func (r *firmwareRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.FirmwareRepositoryResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.gatewayClient.DeleteFirmwareRepository(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting firmware repository",
			"Couldn't Delete firmware repository "+err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for fault set Resource
func (r *firmwareRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Please provide valid firmware repository ID", "Please provide valid firmware repository ID")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
