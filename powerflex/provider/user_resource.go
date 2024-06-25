/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"strings"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// UserResource - function to return resource interface
func UserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client  *goscaleio.Client
	system  *goscaleio.System
	version string
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage the User entity of the PowerFlex Array. We can Create, Update and Delete the user using this resource. We can also import an existing user from the PowerFlex array.",
		MarkdownDescription: "This resource is used to manage the User entity of the PowerFlex Array. We can Create, Update and Delete the user using this resource. We can also import an existing user from the PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the user.",
				Computed:            true,
				MarkdownDescription: "The ID of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system_id": schema.StringAttribute{
				Description:         "The ID of the system.",
				Computed:            true,
				MarkdownDescription: "The ID of the system.",
			},

			"name": schema.StringAttribute{
				Description: "The name of the user. For PowerFlex version 3.6, cannot be updated.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(31),
				},
				MarkdownDescription: "The name of the user. For PowerFlex version 3.6, cannot be updated.",
			},
			"role": schema.StringAttribute{
				Description: "The role of the user." +
					" Accepted values for PowerFlex version 3.6 are 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'." +
					" Accepted values for PowerFlex version 4.5 are 'Monitor', 'SuperUser', 'SystemAdmin', 'StorageAdmin', 'LifecycleAdmin', 'ReplicationManager', 'SnapshotManager', 'SecurityAdmin', 'DriveReplacer', 'Technician', 'Support'.",
				Required: true,
				MarkdownDescription: "The role of the user." +
					" Accepted values for PowerFlex version 3.6 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'." +
					" Accepted values for PowerFlex version 4.5 are 'Monitor', 'SuperUser', 'SystemAdmin', 'StorageAdmin', 'LifecycleAdmin', 'ReplicationManager', 'SnapshotManager', 'SecurityAdmin', 'DriveReplacer', 'Technician', 'Support'.",
				Validators: []validator.String{stringvalidator.OneOf(
					"Monitor",
					"Configure",
					"Administrator",
					"Security",
					"FrontendConfig",
					"BackendConfig",
					"SuperUser",
					"SystemAdmin",
					"StorageAdmin",
					"LifecycleAdmin",
					"ReplicationManager",
					"SnapshotManager",
					"SecurityAdmin",
					"DriveReplacer",
					"Technician",
					"Support",
				)},
			},
			"password": schema.StringAttribute{
				Description: "Password of the user. For PowerFlex version 3.6, cannot be updated.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(31),
				},
				MarkdownDescription: "Password of the user. For PowerFlex version 3.6, cannot be updated.",
			},
			"first_name": schema.StringAttribute{
				Description:         "First name of the user. PowerFlex version 3.6 does not support the first_name attribute. It is mandatory for PowerFlex version 4.6.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "First name of the user. PowerFlex version 3.6 does not support the first_name attribute. It is mandatory for PowerFlex version 4.6.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_name": schema.StringAttribute{
				Description:         "Last name of the user. PowerFlex version 3.6 does not support the last_name attribute. It is mandatory for PowerFlex version 4.6.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Last name of the user. PowerFlex version 3.6 does not support the last_name attribute. It is mandatory for PowerFlex version 4.6.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
	r.version = r.client.GetConfigConnect().Version
}

// ValidateConfig validates the configuration for the user resource.
func (r *userResource) ValidateConfig(plan models.UserModel) diag.Diagnostics {
	roleMap := make(map[string][]string)
	var diags diag.Diagnostics
	var flag bool

	roleMap["3.5"] = []string{"Monitor", "Configure", "Administrator", "Security", "FrontendConfig", "BackendConfig"}
	roleMap["4.0"] = []string{"Monitor", "SuperUser", "SystemAdmin", "StorageAdmin", "LifecycleAdmin", "ReplicationManager", "SnapshotManager", "SecurityAdmin", "DriveReplacer", "Technician", "Support"}

	if r.version == models.Version3X {
		roles := roleMap[models.Version3X]

		if !plan.FirstName.IsUnknown() || !plan.LastName.IsUnknown() {
			diags.AddError(
				"PowerFlex version 3.6 does not support the first_name and last_name attributes.",
				"PowerFlex version 3.6 does not support the first_name and last_name attributes.",
			)
		}

		for _, value := range roles {
			if value == plan.Role.ValueString() {
				flag = true
				break
			}
		}

		if !flag {
			diags.AddAttributeError(
				path.Root("role"),
				"Invalid user role",
				"Supported values for role in PowerFlex version 3.6 are 'Administrator', 'Monitor', 'Configure', 'Security', 'FrontendConfig', 'BackendConfig'",
			)
		}
	} else {
		roles := roleMap["4.0"]

		for _, value := range roles {
			if value == plan.Role.ValueString() {
				flag = true
				break
			}
		}

		if !flag {
			diags.AddAttributeError(
				path.Root("role"),
				"Invalid user role",
				"Supported values for role in PowerFlex version 4.5 are 'Monitor', 'SuperUser', 'SystemAdmin', 'StorageAdmin', 'LifecycleAdmin', 'ReplicationManager', 'SnapshotManager', 'SecurityAdmin', 'DriveReplacer', 'Technician', 'Support'",
			)
		}
	}
	return diags
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan    models.UserModel
		user    *scaleiotypes.User
		ssoUser *scaleiotypes.SSOUserDetails
		err     error
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.ValidateConfig(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the user payload by setting the values from plan
	if r.version == models.Version3X {
		payload := &scaleiotypes.UserParam{
			Name:     plan.Name.ValueString(),
			UserRole: plan.Role.ValueString(),
			Password: plan.Password.ValueString(),
		}
		// create the user
		response, err := r.system.CreateUser(payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating the user", "Could not create user, unexpected error: "+err.Error(),
			)
			return
		}

		//fetch the user
		user, err = r.system.GetUserByIDName(response.ID, "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error fetching the user after creation", "Could not fetch user, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		payload := &scaleiotypes.SSOUserCreateParam{
			UserName:  plan.Name.ValueString(),
			Role:      plan.Role.ValueString(),
			Password:  plan.Password.ValueString(),
			FirstName: plan.FirstName.ValueString(),
			LastName:  plan.LastName.ValueString(),
			Type:      "Local",
		}
		// create the user
		ssoUser, err = r.client.CreateSSOUser(payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating the user", "Could not create user, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// update the state as per the values fetched
	state := helper.UpdateUserState(user, plan, ssoUser)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		state    models.UserModel
		user     *scaleiotypes.User
		ssoUser  *scaleiotypes.SSOUserDetails
		ssoUsers *scaleiotypes.SSOUserList
		err      error
	)

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//fetch the user
	if !state.ID.IsNull() {
		if r.version == models.Version3X {
			user, err = r.system.GetUserByIDName(state.ID.ValueString(), "")
		} else {
			ssoUser, err = r.client.GetSSOUser(state.ID.ValueString())
		}

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
			)
			return
		}

	} else {
		if r.version == models.Version3X {
			user, err = r.system.GetUserByIDName("", state.Name.ValueString())
		} else {
			ssoUsers, err = r.client.GetSSOUserByFilters("username", state.Name.ValueString())
			if len(ssoUsers.SSOUsers) > 0 {
				ssoUser = &ssoUsers.SSOUsers[0]
			}
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not get user by name %s", state.Name.ValueString()), err.Error(),
			)
			return
		}
	}
	// update the state as per the values fetched
	response := helper.UpdateUserState(user, state, ssoUser)
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan     models.UserModel
		state    models.UserModel
		response models.UserModel
		user     *scaleiotypes.User
		ssoUser  *scaleiotypes.SSOUserDetails
		err      error
	)

	// Get the plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.ValidateConfig(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check if name is updated and if it's updated then throw the error
	if r.version == models.Version3X && !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		resp.Diagnostics.AddError(
			"username cannot be updated once the user is created.", "unexpected error: username change is not supported",
		)
	}
	// check if password is updated and if it's updated then throw the error
	if r.version == models.Version3X && !plan.Password.IsUnknown() && !plan.Password.Equal(state.Password) {
		resp.Diagnostics.AddError(
			"password cannot be updated after user creation.", "unexpected error: password change is not supported",
		)
	}

	// check if role is updated and if it's updated then set the role as per the plan
	if r.version == models.Version3X {
		if plan.Role.ValueString() != state.Role.ValueString() {
			payload := &scaleiotypes.UserRoleParam{
				UserRole: plan.Role.ValueString(),
			}
			err = r.system.SetUserRole(payload, state.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error while updating role of the user", err.Error(),
				)
			}
		}

		// fetch the user
		user, err = r.system.GetUserByIDName(state.ID.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
			)
			return
		}
	}

	if r.version == "4.0" {
		if plan.Role.ValueString() != state.Role.ValueString() || plan.Name.ValueString() != state.Name.ValueString() || plan.FirstName.ValueString() != state.FirstName.ValueString() || plan.LastName.ValueString() != state.LastName.ValueString() {
			payload := &scaleiotypes.SSOUserModifyParam{
				UserName:  plan.Name.ValueString(),
				Role:      plan.Role.ValueString(),
				FirstName: plan.FirstName.ValueString(),
				LastName:  plan.LastName.ValueString(),
			}

			_, err = r.client.ModifySSOUser(state.ID.ValueString(), payload)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error while updating role/username of the user", err.Error(),
				)
			}
		}

		if plan.Password.ValueString() != state.Password.ValueString() {
			payload := &scaleiotypes.SSOUserModifyParam{
				Password: plan.Password.ValueString(),
			}

			err = r.client.ResetSSOUserPassword(state.ID.ValueString(), payload)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error while updating password of the user", err.Error(),
				)
			}
		}

		// fetch the user
		ssoUser, err = r.client.GetSSOUser(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not get user by ID %s", state.ID.ValueString()), err.Error(),
			)
			return
		}
	}

	// set the state
	if r.version == models.Version3X {
		response = helper.UpdateUserState(user, state, ssoUser)
	} else {
		response = helper.UpdateUserState(user, plan, ssoUser)
	}
	diags = resp.State.Set(ctx, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		state models.UserModel
		err   error
	)

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete the user
	if r.version == models.Version3X {
		err = r.system.RemoveUser(state.ID.ValueString())
	} else {
		err = r.client.DeleteSSOUser(state.ID.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting User", "Couldn't Delete User "+err.Error(),
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

	if strings.Contains(req.ID, ":") {
		parts := strings.Split(req.ID, ":")
		if len(parts) == 2 && parts[0] == "id" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
		} else if len(parts) == 2 && parts[0] == "name" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
		} else {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier format: id:userId or name:userName. Got: %q", req.ID),
			)
		}
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	}
}
