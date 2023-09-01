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
	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"strings"
)

var (
	_ resource.Resource = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// UserResource - function to return resource interface
func UserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to manage user on a PowerFlex array.",
		MarkdownDescription: "This resource can be used to manage user on a PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the user.",
				Computed:            true,
				MarkdownDescription: "The ID of the user.",
			},
			"system_id": schema.StringAttribute{
				Description:         "The ID of the system.",
				Computed:            true,
				MarkdownDescription: "The ID of the system.",
			},

			"name": schema.StringAttribute{
				Description: "The name of the user." +
					" Cannot be updated.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(31),
				},
				MarkdownDescription: "The name of the user." +
					" Cannot be updated.",
			},
			"role": schema.StringAttribute{
				Description: "The role of the user." +
					" Accepted values are 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'.",
				Required: true,
				MarkdownDescription: "The role of the user." +
					" Accepted values are 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'.",
				Validators: []validator.String{stringvalidator.OneOf(
					"Monitor",
					"Configure",
					"Administrator",
					"Security",
					"FrontendConfig",
					"BackendConfig",
				)},
			},
			"password": schema.StringAttribute{
				Description: "Password of the user." +
					" Cannot be updated.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(31),
				},
				MarkdownDescription: "Password of the user." +
					" Cannot be updated.",
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.UserModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the user payload by setting the values from plan
	payload := &scaleiotypes.UserParam{
		Name:     plan.Name.ValueString(),
		UserRole: plan.Role.ValueString(),
		Password: plan.Password.ValueString(),
	}

	// create the user
	response, err2 := r.system.CreateUser(payload)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error creating the user", "Could not create user, unexpected error: "+err2.Error(),
		)
		return
	}

	//fetch the user
	user, err3 := r.system.GetUserByIDName(response.ID, "")
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error fetching the user after creation", "Could not fetch user, unexpected error: "+err3.Error(),
		)
		return
	}

	// update the state as per the values fetched
	state := helper.UpdateUserState(user, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.UserModel
	var user *scaleiotypes.User
	var err error
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//fetch the user
	if !state.ID.IsNull() {
		user, err = r.system.GetUserByIDName( state.ID.ValueString(),"")
		if err != nil {
			resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
			)
			return
		}

	} else {
		user, err = r.system.GetUserByIDName("", state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get user by name %s", state.Name.ValueString()), err.Error(),
			)
			return
		}	
	}
	// update the state as per the values fetched
	response := helper.UpdateUserState(user, state)
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}		

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.UserModel

	// Get the plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state models.UserModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check if name is updated and if it's updated then throw the error
	if !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		resp.Diagnostics.AddError(
			"username cannot be updated once the user is created.", "unexpected error: username change is not supported",
		)
	}
	// check if password is updated and if it's updated then throw the error
	if !plan.Password.IsUnknown() && !plan.Password.Equal(state.Password) {
		resp.Diagnostics.AddError(
			"password cannot be updated after user creation.", "unexpected error: password change is not supported",
		)
	}

	// check if role is updated and if it's updated then set the role as per the plan
	if plan.Role.ValueString() != state.Role.ValueString() {
		payload := &scaleiotypes.UserRoleParam{
			UserRole: plan.Role.ValueString(),
		}
		err2 := r.system.SetUserRole(payload, state.ID.ValueString())
		if err2 != nil {
			resp.Diagnostics.AddError(
				"Error while updating role of the user", err2.Error(),
			)
		}
	}

	// fetch the user
	user, err := r.system.GetUserByIDName(state.ID.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
		)
		return
	}

	// set the state
	response := helper.UpdateUserState(user, state)
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.UserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete the user
	err2 := r.system.RemoveUser(state.ID.ValueString())
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error Deleting User", "Couldn't Delete User "+err2.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// remove from state
	resp.State.RemoveResource(ctx)
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	var userName string
    if len(parts) == 2 {
        userName = parts[1]
    } else {
		resp.Diagnostics.AddError(
            "Unexpected Import Identifier",
            fmt.Sprintf("Expected import identifier with format: name:username. Got: %q", req.ID),
        )
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), userName)...)
}