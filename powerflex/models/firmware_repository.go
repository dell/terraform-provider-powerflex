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

// FirmwareRepositoryResourceModel defines struct firmware repository resource
type FirmwareRepositoryResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SourceLocation types.String `tfsdk:"source_location"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Approve        types.Bool   `tfsdk:"approve"`
	Name           types.String `tfsdk:"name"`
	DiskLocation   types.String `tfsdk:"disk_location"`
	FileName       types.String `tfsdk:"file_name"`
	DefaultCatalog types.Bool   `tfsdk:"default_catalog"`
	Timeout        types.Int64  `tfsdk:"timeout"`
}
