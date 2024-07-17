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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CompatibilityManagementDatasourceModel defines the schema for the compatibility management datasource
type CompatibilityManagementDatasourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Source                 types.String `tfsdk:"source"`
	RepositoryPath         types.String `tfsdk:"repository_path"`
	CurrentVersion         types.String `tfsdk:"current_version"`
	AvailabieVersion       types.String `tfsdk:"available_version"`
	CompatibilityData      types.String `tfsdk:"compatibility_data"`
	CompatibilityDataBytes types.String `tfsdk:"compatibility_data_bytes"`
}
