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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
)

var (
	_ resource.Resource = &userResource{}
)

// UserResource - function to return resource interface
func UserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *goscaleio.Client
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to manage mapping of volumes to an SDC on a PowerFlex array.",
		MarkdownDescription: "This resource can be used to manage mapping of volumes to an SDC on a PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the user.",
				Computed:            true,
				MarkdownDescription: "The ID of the user.",
			},
			"name": schema.StringAttribute{
				Description:         "The name of the user.",
				Required:            true,
				MarkdownDescription: "The name of the user.",
			},
			"role": schema.StringAttribute{
				Description:         "The role of the user.",
				Required:            true,
				MarkdownDescription: "The role of the user.",
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
				Description:         "Password of the user.",
				Required:            true,
				MarkdownDescription: "Password of the user.",
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
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.UserModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting the system", "Could not get the system, unexpected err: "+err.Error(),
		)
		return
	}

	payload := &scaleiotypes.UserParam{
		Name:     plan.Name.ValueString(),
		UserRole: plan.Role.ValueString(),
		Password: plan.Password.ValueString(),
	}

	// create the user
	response, err2 := system.CreateUser(payload)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error creating the user", "Could not create user, unexpected error: "+err2.Error(),
		)
		return
	}
	user, err3 := system.GetUserByIDName(response.ID, "")
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error fetching the user after creation", "Could not fetch user, unexpected error: "+err3.Error(),
		)
		return
	}

	state := UpdateUserState(user, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
func UpdateUserState(user *scaleiotypes.User, plan models.UserModel) models.UserModel {
	state := plan
	state.Name = types.StringValue(user.Name)
	state.Role = types.StringValue(user.UserRole)
	state.Password = plan.Password
	state.ID = types.StringValue(user.ID)
	return state
}
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.UserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster", err.Error(),
		)
		return
	}

	user, err := system.GetUserByIDName(state.ID.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
		)
		return
	}
	response := UpdateUserState(user, state)
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.UserModel

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
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting the user", err.Error(),
		)
		return
	}

	if !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		resp.Diagnostics.AddError(
			"username cannot be updated once the user is created.", "unexpected error: username change is not supported",
		)
	}
	if !plan.Password.IsUnknown() && !plan.Password.Equal(state.Password) {
		resp.Diagnostics.AddError(
			"password cannot be updated after user creation.", "unexpected error: password change is not supported",
		)
	}
	if plan.Role.ValueString() != state.Role.ValueString() {
		payload := &scaleiotypes.UserRoleParam{
			UserRole: plan.Role.ValueString(),
		}
		err2 := system.SetUserRole(payload, state.ID.ValueString())
		if err2 != nil {
			resp.Diagnostics.AddError(
				"Error while updating role of the user", err.Error(),
			)
		}
	}
	user, err := system.GetUserByIDName(state.ID.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
		)
		return
	}
	response := UpdateUserState(user, state)
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.UserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while getting the user", err.Error(),
		)
		return
	}
	err2 := system.RemoveUser(state.ID.ValueString())
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error Deleting User", "Couldn't Delete User "+err2.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}
