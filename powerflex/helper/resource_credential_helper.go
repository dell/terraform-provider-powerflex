/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"context"
	"fmt"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Enum for resource credentials types
type RCType string

const (
	RCNode               RCType = "Node"
	RCSwitch             RCType = "Switch"
	RCvCenter            RCType = "vCenter"
	RCElementManager     RCType = "ElementManager"
	RCPowerflexGateway   RCType = "PowerflexGateway"
	RCPresentationServer RCType = "PresentationServer"
	RCOSAdmin            RCType = "OSAdmin"
	RCOSUser             RCType = "OSUser"
)

func (c RCType) String() string {
	switch c {
	case RCNode:
		return "Node"
	case RCSwitch:
		return "Switch"
	case RCvCenter:
		return "vCenter"
	case RCElementManager:
		return "ElementManager"
	case RCPowerflexGateway:
		return "PowerflexGateway"
	case RCPresentationServer:
		return "PresentationServer"
	case RCOSAdmin:
		return "OSAdmin"
	case RCOSUser:
		return "OSUser"
	default:
		return "Unknown"
	}
}

// GetResourceCredentials returns all the resource credentials
func GetResourceCredentials(ctx context.Context, system *goscaleio.System) ([]scaleiotypes.CredObj, error) {
	var resourceCredentials []scaleiotypes.CredObj
	creds, err := system.GetResourceCredentials()
	if err != nil {
		return nil, err
	}
	for _, val := range creds.Credentials {
		resourceCredentials = append(resourceCredentials, val.Credential)
	}
	return resourceCredentials, nil
}

// MapResourceCredentials a terraform mapped resource
func MapResourceCredentials(rcs []scaleiotypes.CredObj, state models.ResourceCredentialDataSourceModel) models.ResourceCredentialDataSourceModel {
	var resourceCredentialDetails []models.ResourceCredentialModel
	for _, val := range rcs {
		resourceCredentialDetails = append(resourceCredentialDetails,
			models.ResourceCredentialModel{
				ID:          types.StringValue(val.ID),
				Type:        types.StringValue(val.Type),
				CreateDate:  types.StringValue(val.CreateDate),
				CreatedBy:   types.StringValue(val.CreatedBy),
				UpdatedBy:   types.StringValue(val.UpdatedBy),
				UpdatedDate: types.StringValue(val.UpdatedDate),
				Label:       types.StringValue(val.Label),
				Domain:      types.StringValue(val.Domain),
				Username:    types.StringValue(val.Username),
			})
	}
	state.ResourceCredentialDetails = resourceCredentialDetails
	return state
}

func ValidateResourceCredentialResource(plan models.ResourceCredentialResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	// Validate plan based on type
	switch RCType(plan.Type.ValueString()) {
	case RCNode:
		if !plan.Snmpv2CommunityString.IsNull() && !plan.Snmpv3SecurityName.IsNull() {
			diags.AddError(
				"Only one of SNMPv2 or SNMPv3 can be set",
				"Only one of SNMPv2 or SNMPv3 can be set",
			)
			return diags
		}
		if plan.Snmpv3SecurityLevel.ValueString() == "Moderate" && plan.Snmpv3Md5AuthenticationPassword.IsNull() {
			diags.AddError(
				"SNMPv3 security level Moderate requires SNMPv3 MD5 authentication password.",
				"SNMPv3 security level Moderate requires SNMPv3 MD5 authentication password.",
			)
			return diags
		}
		if plan.Snmpv3SecurityLevel.ValueString() == "Maximal" && (plan.Snmpv3Md5AuthenticationPassword.IsNull() || plan.Snmpv3DesAuthenticationPassword.IsNull()) {
			diags.AddError(
				"SNMPv3 security level Maximal requires SNMPv3 MD5 DES authentication password.",
				"SNMPv3 security level Maximal requires SNMPv3 MD5 and DES authentication password.",
			)
			return diags
		}
		if (plan.SSHPrivateKey.IsNull() && !plan.KeyPairName.IsNull()) || (!plan.SSHPrivateKey.IsNull() && plan.KeyPairName.IsNull()) {
			diags.AddError(
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
			)
			return diags
		}
	case RCSwitch:
		if (plan.SSHPrivateKey.IsNull() && !plan.KeyPairName.IsNull()) || (!plan.SSHPrivateKey.IsNull() && plan.KeyPairName.IsNull()) {
			diags.AddError(
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
			)
			return diags
		}
	case RCvCenter:
		if plan.Domain.IsNull() {
			diags.AddError(
				"domain must be set for type vcenter_credential",
				"domain must be set for type vcenter_credential",
			)
			return diags
		}
	case RCElementManager:
		// EM doesnt have any required fields besides username password
	case RCPowerflexGateway:
		if plan.OSUsername.IsNull() || plan.OSPassword.IsNull() {
			diags.AddError(
				"os_username and os_password must be set for type powerflex_gateway",
				"os_username and os_password must be set for type powerflex_gateway",
			)
			return diags
		}
	case RCPresentationServer:
		// PS doesnt have any required fields besides username password
	case RCOSAdmin:
		if (plan.SSHPrivateKey.IsNull() && !plan.KeyPairName.IsNull()) || (!plan.SSHPrivateKey.IsNull() && plan.KeyPairName.IsNull()) {
			diags.AddError(
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
			)
			return diags
		}
	case RCOSUser:
		if (plan.SSHPrivateKey.IsNull() && !plan.KeyPairName.IsNull()) || (!plan.SSHPrivateKey.IsNull() && plan.KeyPairName.IsNull()) {
			diags.AddError(
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
				"ssh_private_key and key_pair_name must either both be set or neither be set.",
			)
			return diags
		}
	}
	return diags
}

// MapResourceCredentialResource a terraform mapped resource
func MapResourceCredentialResource(rc scaleiotypes.CredObj, state models.ResourceCredentialResourceModel) models.ResourceCredentialResourceModel {
	return models.ResourceCredentialResourceModel{
		Name:                            state.Name,
		Username:                        state.Username,
		Password:                        state.Password,
		Domain:                          state.Domain,
		Type:                            state.Type,
		Snmpv2CommunityString:           state.Snmpv2CommunityString,
		Snmpv3SecurityName:              state.Snmpv3SecurityName,
		Snmpv3SecurityLevel:             state.Snmpv3SecurityLevel,
		Snmpv3Md5AuthenticationPassword: state.Snmpv3Md5AuthenticationPassword,
		Snmpv3DesAuthenticationPassword: state.Snmpv3DesAuthenticationPassword,
		SSHPrivateKey:                   state.SSHPrivateKey,
		KeyPairName:                     state.KeyPairName,
		OSUsername:                      state.OSUsername,
		OSPassword:                      state.OSPassword,
		ID:                              types.StringValue(rc.ID),
		CreateDate:                      types.StringValue(rc.CreateDate),
		CreatedBy:                       types.StringValue(rc.CreatedBy),
		UpdatedBy:                       types.StringValue(rc.UpdatedBy),
		UpdatedDate:                     types.StringValue(rc.UpdatedDate),
	}
}

// CreateResourceCredential create a new resource credential
func CreateResourceCredential(ctx context.Context, sys *goscaleio.System, plan models.ResourceCredentialResourceModel) (*scaleiotypes.ResourceCredential, error) {
	// Create resource credential
	switch RCType(plan.Type.ValueString()) {
	case RCNode:
		return sys.CreateNodeResourceCredential(scaleiotypes.ServerCredential{
			Username:                        plan.Username.ValueString(),
			Password:                        plan.Password.ValueString(),
			Label:                           plan.Name.ValueString(),
			SNMPv2CommunityString:           plan.Snmpv2CommunityString.ValueString(),
			SNMPv3SecurityName:              plan.Snmpv3SecurityName.ValueString(),
			SNMPv3SecurityLevel:             translateSnmpv3SecurityLevel(plan.Snmpv3SecurityLevel.ValueString()),
			SNMPv3MD5AuthenticationPassword: plan.Snmpv3Md5AuthenticationPassword.ValueString(),
			SNMPv3DesPrivatePassword:        plan.Snmpv3DesAuthenticationPassword.ValueString(),
			SSHPrivateKey:                   plan.SSHPrivateKey.ValueString(),
			KeyPairName:                     plan.KeyPairName.ValueString(),
		})
	case RCSwitch:
		return sys.CreateSwitchResourceCredential(scaleiotypes.IomCredential{
			Username:              plan.Username.ValueString(),
			Password:              plan.Password.ValueString(),
			Label:                 plan.Name.ValueString(),
			SNMPv2CommunityString: plan.Snmpv2CommunityString.ValueString(),
			SSHPrivateKey:         plan.SSHPrivateKey.ValueString(),
			KeyPairName:           plan.KeyPairName.ValueString(),
		})
	case RCvCenter:
		return sys.CreateVCenterResourceCredential(scaleiotypes.VCenterCredential{
			Username: plan.Username.ValueString(),
			Password: plan.Password.ValueString(),
			Label:    plan.Name.ValueString(),
			Domain:   plan.Domain.ValueString(),
		})
	case RCElementManager:
		return sys.CreateElementManagerResourceCredential(scaleiotypes.EMCredential{
			Username:              plan.Username.ValueString(),
			Password:              plan.Password.ValueString(),
			Label:                 plan.Name.ValueString(),
			Domain:                plan.Domain.ValueString(),
			SNMPv2CommunityString: plan.Snmpv2CommunityString.ValueString(),
		})
	case RCPowerflexGateway:
		return sys.CreateScaleIOResourceCredential(scaleiotypes.ScaleIOCredential{
			AdminUsername: plan.Username.ValueString(),
			AdminPassword: plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			OSUsername:    plan.OSUsername.ValueString(),
			OSPassword:    plan.OSPassword.ValueString(),
		})
	case RCPresentationServer:
		return sys.CreatePresentationServerResourceCredential(scaleiotypes.PSCredential{
			Username: plan.Username.ValueString(),
			Password: plan.Password.ValueString(),
			Label:    plan.Name.ValueString(),
		})
	case RCOSAdmin:
		return sys.CreateOsAdminResourceCredential(scaleiotypes.OSAdminCredential{
			Username:      plan.Username.ValueString(),
			Password:      plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			SSHPrivateKey: plan.SSHPrivateKey.ValueString(),
			KeyPairName:   plan.KeyPairName.ValueString(),
		})
	case RCOSUser:
		return sys.CreateOsUserResourceCredential(scaleiotypes.OSUserCredential{
			Username:      plan.Username.ValueString(),
			Password:      plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			Domain:        plan.Domain.ValueString(),
			SSHPrivateKey: plan.SSHPrivateKey.ValueString(),
			KeyPairName:   plan.KeyPairName.ValueString(),
		})
	}

	return nil, fmt.Errorf("invalid type unable to create resource credential.")
}

// ModifyResourceCredential Modify an existing resource credential
func ModifyResourceCredential(ctx context.Context, sys *goscaleio.System, plan models.ResourceCredentialResourceModel, id string) (*scaleiotypes.ResourceCredential, error) {
	// Create resource credential
	switch RCType(plan.Type.ValueString()) {
	case RCNode:
		return sys.ModifyNodeResourceCredential(scaleiotypes.ServerCredential{
			Username:                        plan.Username.ValueString(),
			Password:                        plan.Password.ValueString(),
			Label:                           plan.Name.ValueString(),
			SNMPv2CommunityString:           plan.Snmpv2CommunityString.ValueString(),
			SNMPv3SecurityName:              plan.Snmpv3SecurityName.ValueString(),
			SNMPv3SecurityLevel:             translateSnmpv3SecurityLevel(plan.Snmpv3SecurityLevel.ValueString()),
			SNMPv3MD5AuthenticationPassword: plan.Snmpv3Md5AuthenticationPassword.ValueString(),
			SNMPv3DesPrivatePassword:        plan.Snmpv3DesAuthenticationPassword.ValueString(),
			SSHPrivateKey:                   plan.SSHPrivateKey.ValueString(),
			KeyPairName:                     plan.KeyPairName.ValueString(),
		}, id)
	case RCSwitch:
		return sys.ModifySwitchResourceCredential(scaleiotypes.IomCredential{
			Username:              plan.Username.ValueString(),
			Password:              plan.Password.ValueString(),
			Label:                 plan.Name.ValueString(),
			SNMPv2CommunityString: plan.Snmpv2CommunityString.ValueString(),
			SSHPrivateKey:         plan.SSHPrivateKey.ValueString(),
			KeyPairName:           plan.KeyPairName.ValueString(),
		}, id)
	case RCvCenter:
		return sys.ModifyVCenterResourceCredential(scaleiotypes.VCenterCredential{
			Username: plan.Username.ValueString(),
			Password: plan.Password.ValueString(),
			Label:    plan.Name.ValueString(),
			Domain:   plan.Domain.ValueString(),
		}, id)
	case RCElementManager:
		return sys.ModifyElementManagerResourceCredential(scaleiotypes.EMCredential{
			Username:              plan.Username.ValueString(),
			Password:              plan.Password.ValueString(),
			Label:                 plan.Name.ValueString(),
			Domain:                plan.Domain.ValueString(),
			SNMPv2CommunityString: plan.Snmpv2CommunityString.ValueString(),
		}, id)
	case RCPowerflexGateway:
		return sys.ModifyScaleIOResourceCredential(scaleiotypes.ScaleIOCredential{
			AdminUsername: plan.Username.ValueString(),
			AdminPassword: plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			OSUsername:    plan.OSUsername.ValueString(),
			OSPassword:    plan.OSPassword.ValueString(),
		}, id)
	case RCPresentationServer:
		return sys.ModifyPresentationServerResourceCredential(scaleiotypes.PSCredential{
			Username: plan.Username.ValueString(),
			Password: plan.Password.ValueString(),
			Label:    plan.Name.ValueString(),
		}, id)
	case RCOSAdmin:
		return sys.ModifyOsAdminResourceCredential(scaleiotypes.OSAdminCredential{
			Username:      plan.Username.ValueString(),
			Password:      plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			SSHPrivateKey: plan.SSHPrivateKey.ValueString(),
			KeyPairName:   plan.KeyPairName.ValueString(),
		}, id)
	case RCOSUser:
		return sys.ModifyOsUserResourceCredential(scaleiotypes.OSUserCredential{
			Username:      plan.Username.ValueString(),
			Password:      plan.Password.ValueString(),
			Label:         plan.Name.ValueString(),
			Domain:        plan.Domain.ValueString(),
			SSHPrivateKey: plan.SSHPrivateKey.ValueString(),
			KeyPairName:   plan.KeyPairName.ValueString(),
		}, id)
	}

	return nil, fmt.Errorf("invalid type unable to create resource credential.")
}

func translateSnmpv3SecurityLevel(level string) string {
	switch level {
	case "Minimal":
		return "1"
	case "Moderate":
		return "2"
	case "Maximal":
		return "3"
	}
	// Return empty string if not found, this means they are using snmpv2
	return ""
}
