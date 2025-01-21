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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ResourceCredentialDataSourceModel model for Resource Credential
type ResourceCredentialDataSourceModel struct {
	ResourceCredentialDetails []ResourceCredentialModel `tfsdk:"resource_credential_details"`
	ID                        types.String              `tfsdk:"id"`
	ResourceCredentialFilter  *ResourceCredentialFilter `tfsdk:"filter"`
}

// ResourceCredentialFilter defines the model for filters used for ResourceCredentialDataSourceModel
type ResourceCredentialFilter struct {
	ID          []types.String `tfsdk:"id"`
	Type        []types.String `tfsdk:"type"`
	CreateDate  []types.String `tfsdk:"created_date"`
	CreatedBy   []types.String `tfsdk:"created_by"`
	UpdatedBy   []types.String `tfsdk:"updated_by"`
	UpdatedDate []types.String `tfsdk:"updated_date"`
	Label       []types.String `tfsdk:"label"`
	Domain      []types.String `tfsdk:"domain"`
	Username    []types.String `tfsdk:"username"`
}

// ResourceCredentialModel model for Resource Credential
type ResourceCredentialModel struct {
	ID          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	CreateDate  types.String `tfsdk:"created_date"`
	CreatedBy   types.String `tfsdk:"created_by"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
	UpdatedDate types.String `tfsdk:"updated_date"`
	Label       types.String `tfsdk:"label"`
	Domain      types.String `tfsdk:"domain"`
	Username    types.String `tfsdk:"username"`
}
