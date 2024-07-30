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
	"terraform-provider-powerflex/powerflex/constants"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource              = &compatibilityManagmentResource{}
	_ resource.ResourceWithConfigure = &compatibilityManagmentResource{}
)

// NewCompatibilityManagementResource - function to return resource interface
func NewCompatibilityManagementResource() resource.Resource {
	return &compatibilityManagmentResource{}
}

type compatibilityManagmentResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *compatibilityManagmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compatibility_management"
}

func (r *compatibilityManagmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CompatibilityManagementResourceSchema
}

func (r *compatibilityManagmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Function used to set Compatibility Managment Resource
func (r *compatibilityManagmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.CompatibilityManagementDatasourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, errSet := helper.SetCompatibilityManagement(ctx, r.system, plan)
	if errSet != nil {
		resp.Diagnostics.AddError(
			"Error setting compatibility management",
			errSet.Error(),
		)
		return
	}

	// Set state
	state := helper.MapCompatibilityManagementState(ctx, response)
	diagsState := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diagsState...)

}

// Function used to Read Compatibility Managment Resource
func (r *compatibilityManagmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.CompatibilityManagementDatasourceModel
	// Get the Compatibility Management details
	cm, err := helper.GetCompatibilityManagement(ctx, r.system)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting compatibility management details",
			err.Error(),
		)
		return
	}

	// Set state
	state = helper.MapCompatibilityManagementState(ctx, cm)
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Function used to Update Compatibility Managment Resource
func (r *compatibilityManagmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		constants.UpdatesNotSupportedErrorMsg,
		constants.UpdatesNotSupportedErrorMsg,
	)
}

// Function used to Delete Compatibility Managment Resource
func (r *compatibilityManagmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		constants.DeleteIsNotSupportedErrorMsg,
		constants.DeleteIsNotSupportedErrorMsg,
	)
}
