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
	"time"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &osRepositoryResource{}
	_ resource.ResourceWithConfigure = &osRepositoryResource{}
)

// NewOsRepositoryResource - function to return resource interface
func NewOsRepositoryResource() resource.Resource {
	return &osRepositoryResource{}
}

type osRepositoryResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *osRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_os_repository"
}

func (r *osRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OsRepositoryResourceSchema
}

func (r *osRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Function used to Create OS Repository Resource
func (r *osRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create OS Repository")
	// Retrieve values from plan

	var plan models.OSRepositoryResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Create OS Repository
	osRepoResp, err := helper.CreateOSRepository(r.system, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create OS Repository",
			err.Error(),
		)
		return
	}
	// Wait for resource id to be generated
	time.Sleep(30 * time.Second)
	repoID, errGet := helper.GetOsRepositoryID(r.system, osRepoResp.Name)
	if errGet != nil {
		resp.Diagnostics.AddError(
			"Could not get the OS Repository id",
			errGet.Error(),
		)
		return
	}
	startTime := time.Now()
	var endTime time.Time
	if !plan.Timeout.IsNull() {
		endTime = startTime.Add(time.Duration(plan.Timeout.ValueInt64()) * time.Minute)
	}
	var osRepo *scaleiotypes.OSRepository
	// will loop until the os repository is created successfully and status becomes available or until it times out
	for time.Now().Before(endTime) {
		time.Sleep(1 * time.Minute)
		osRepo, err = helper.GetOSRepositoryByID(r.system, repoID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not get the OS Repository Details",
				err.Error(),
			)
			return
		}
		if osRepo.State == "available" {
			break
		} else if osRepo.State == "copying" {
			continue
		} else {
			resp.Diagnostics.AddError(
				"Could not create OS Repository",
				fmt.Sprintf(" with name: %s", osRepo.Name),
			)
			return
		}
	}
	state := helper.UpdateOsRepositoryState(osRepo, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read OS Repository Resource
func (r *osRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read OS Repository")
	// Get current state
	var state models.OSRepositoryResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	osRepo, err := helper.GetOSRepositoryByID(r.system, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get the OS Repository Details",
			err.Error(),
		)
		return
	}

	osRepoState := helper.UpdateOsRepositoryState(osRepo, state)
	diags = resp.State.Set(ctx, osRepoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update OS repository Resource
func (r *osRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		constants.UpdatesNotSupportedErrorMsg,
		constants.UpdatesNotSupportedErrorMsg,
	)
}

// Function used to Delete OS Repository Resource
func (r *osRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.OSRepositoryResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.system.RemoveOSRepository(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting OS repository",
			"Couldn't Delete OS repository "+err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}
