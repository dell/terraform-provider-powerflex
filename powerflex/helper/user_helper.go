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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserState is the helper function that marshals API response to UserModel
func UpdateUserState(user *scaleiotypes.User, plan models.UserModel, ssoUser *scaleiotypes.SSOUserDetails) models.UserModel {
	state := plan
	if user != nil {
		state.Name = types.StringValue(user.Name)
		state.Role = types.StringValue(user.UserRole)
		state.Password = plan.Password
		state.ID = types.StringValue(user.ID)
		state.SystemID = types.StringValue(user.SystemID)
	} else {
		state.Name = types.StringValue(ssoUser.Username)
		state.Role = types.StringValue(ssoUser.Permission.Role)
		state.Password = plan.Password
		state.ID = types.StringValue(ssoUser.ID)
		state.SystemID = types.StringNull()
	}

	return state
}
