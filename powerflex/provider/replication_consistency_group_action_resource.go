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
	"terraform-provider-powerflex/powerflex/constants"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &replicationConsistencyGroupResourceAction{}
	_ resource.ResourceWithConfigure = &replicationConsistencyGroupResourceAction{}
)

// ReplicationConsistencyGroupActionResource - function to return resource interface
func ReplicationConsistencyGroupActionResource() resource.Resource {
	return &replicationConsistencyGroupResourceAction{}
}

// replicationConsistencyGroupResourceAction - struct to define ReplicationConsistencyGroup resource
type replicationConsistencyGroupResourceAction struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResourceAction) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_consistency_group_action"
}

// Schema - function to return Schema for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResourceAction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ReplicationConsistencyGroupActionReourceSchema
}

// Configure - function to return Configuration for ReplicationConsistencyGroup resource.
func (r *replicationConsistencyGroupResourceAction) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create - function to Create for ReplicationConsistencyGroupAction resource.
func (r *replicationConsistencyGroupResourceAction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")
	var plan models.ReplicationConsistencyGroupAction

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.RCGDoAction(r.client, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error doing action %s on replication consistency group", plan.Action.ValueString()),
			err.Error(),
		)
		return
	}

	diagsState := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diagsState...)
}

// Read - function to Read for ReplicationConsistencyGroupAction resource.
func (r *replicationConsistencyGroupResourceAction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.ReplicationConsistencyGroupAction

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsState := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diagsState...)
}

// Update - function to Update for ReplicationConsistencyGroupAction resource.
func (r *replicationConsistencyGroupResourceAction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		constants.UpdatesNotSupportedErrorMsg+", to create update and delete replication consistency groups use the powerflex_replication_consistency_group resource",
		constants.UpdatesNotSupportedErrorMsg+", to create update and delete replication consistency groups use the powerflex_replication_consistency_group resource",
	)

	var plan models.ReplicationConsistencyGroupAction

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsState := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diagsState...)

}

// Delete - function to Delete for ReplicationConsistencyGroupAction resource.
func (r *replicationConsistencyGroupResourceAction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		constants.DeleteIsNotSupportedErrorMsg+", to create update and delete replication consistency groups use the powerflex_replication_consistency_group resource",
		constants.DeleteIsNotSupportedErrorMsg+", to create update and delete replication consistency groups use the powerflex_replication_consistency_group resource",
	)
}
