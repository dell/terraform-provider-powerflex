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
	"fmt"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewNvmeHostResource is a helper function to simplify the provider implementation.
func NewNvmeHostResource() resource.Resource {
	return &NvmeHostResource{}
}

// NvmeHostResource is the resource implementation.
type NvmeHostResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

// Metadata describes the resource arguments.
func (r *NvmeHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvme_host"
}

// Schema describes the resource arguments.
func (r *NvmeHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manager NVMe host from the PowerFlex array. We can Create, Update and Delete the PowerFlex NVMe host using this resource. We can also import an existing NVMe host from PowerFlex array.",
		MarkdownDescription: "This resource is used to manager NVMe host from the PowerFlex array. We can Create, Update and Delete the PowerFlex NVMe host using this resource. We can also import an existing NVMe host from PowerFlex array.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the NVMe host",
				MarkdownDescription: "ID of the NVMe host",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the NVMe host",
				MarkdownDescription: "Name of the NVMe host",
				Optional:            true,
				Computed:            true,
			},
			"system_id": schema.StringAttribute{
				Description:         "The ID of the system.",
				MarkdownDescription: "The ID of the system.",
				Computed:            true,
			},
			"nqn": schema.StringAttribute{
				Description:         "NQN of the NVMe host. This attribute must be set during creation and cannot be modified afterwards.",
				MarkdownDescription: "NQN of the NVMe host. This attribute must be set during creation and cannot be modified afterwards.",
				Required:            true,
			},
			"max_num_paths": schema.Int64Attribute{
				Description:         "Number of Paths Per Volume.",
				MarkdownDescription: "Number of Paths Per Volume.",
				Optional:            true,
				Computed:            true,
			},
			"max_num_sys_ports": schema.Int64Attribute{
				Description:         "Number of System Ports per Protection Domain.",
				MarkdownDescription: "Number of System Ports per Protection Domain.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NvmeHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	r.system = system
}

// Create used to Create NVMe host Resource
func (r *NvmeHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create NVMe host")

	var plan models.NvmeHostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := scaleiotypes.NvmeHostParam{}
	err := helper.ReadFromState(ctx, plan, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing plan model",
			fmt.Sprintf("Could not read NVMe host param with error: %s", err.Error()),
		)
		return
	}

	if !plan.Name.IsUnknown() && !plan.Name.IsNull() && plan.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Name cannot be empty",
			"Name cannot be empty",
		)
		return
	}

	hostResp, err := r.system.CreateNvmeHost(params)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create NVMe host with name %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	host, err := helper.GetNvmeHostByID(r.system, hostResp.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not read NVMe host with ID %s", hostResp.ID),
			err.Error(),
		)
		return
	}

	err = helper.CopyFields(ctx, host, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating NVMe host",
			fmt.Sprintf("Could not read NVMe host %s with error: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read NVMe host Resource
func (r *NvmeHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read NVMe host")
	// Get current state
	var state models.NvmeHostResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := helper.GetNvmeHostByID(r.system, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the NVMe host Details",
			err.Error(),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, host, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading NVMe host",
			fmt.Sprintf("Could not read NVMe host struct %s with error: %s", state.ID, err.Error()),
		)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update used to NVMe host Resource
func (r *NvmeHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.NvmeHostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.NvmeHostResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := goscaleio.CheckPfmpVersion(r.client, "4.6")

	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking PowerFlex version",
			err.Error(),
		)
		return
	}
	if result < 0 {
		resp.Diagnostics.AddError(
			"Updating NVMe host is not supported",
			"Updating the NVMe host is not supported in PowerFlex versions earlier than 4.6",
		)
		return
	}

	if !plan.Nqn.Equal(state.Nqn) {
		resp.Diagnostics.AddError(
			"nqn cannot be modified after creation",
			"nqn cannot be modified after creation",
		)
		return
	}

	if !plan.Name.IsNull() && plan.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Name cannot be empty",
			"Name cannot be empty",
		)
		return
	}

	if !plan.Name.IsNull() && plan.Name.ValueString() != state.Name.ValueString() {
		err := r.system.ChangeNvmeHostName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update NVMe host name",
				err.Error(),
			)
		}
	}

	if !plan.MaxNumPaths.IsNull() && plan.MaxNumPaths.ValueInt64() != state.MaxNumPaths.ValueInt64() {
		err := r.system.ChangeNvmeHostMaxNumPaths(state.ID.ValueString(), int(plan.MaxNumPaths.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update max_num_paths",
				err.Error(),
			)
		}
	}

	if !plan.MaxNumSysPorts.IsNull() && plan.MaxNumSysPorts.ValueInt64() != state.MaxNumSysPorts.ValueInt64() {
		err := r.system.ChangeNvmeHostMaxNumSysPorts(state.ID.ValueString(), int(plan.MaxNumSysPorts.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update max_num_sys_ports",
				err.Error(),
			)
		}
	}

	host, err := r.system.GetNvmeHostByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting NVMe host after update",
			err.Error(),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, host, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading NVMe host",
			fmt.Sprintf("Could not read NVMe host struct %s with error: %s", state.ID, err.Error()),
		)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *NvmeHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.NvmeHostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.system.DeleteNvmeHost(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete NVMe host",
			err.Error(),
		)

		return
	}
	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource state.
func (r *NvmeHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
