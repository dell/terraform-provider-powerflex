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

	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// UpdateUploadPackageState updates the state
func UpdateUploadPackageState(packageDetails []*goscaleio_types.PackageDetails, plan models.PackageModel) (models.PackageModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	PackageAttrTypes := GetPackageType()
	PackageElemType := types.ObjectType{
		AttrTypes: PackageAttrTypes,
	}

	packages := []attr.Value{}
	for _, vol := range packageDetails {
		objVal, dgs := GetPackageValue(vol)
		diags = append(diags, dgs...)
		packages = append(packages, objVal)
	}
	setVal, dgs := types.SetValue(PackageElemType, packages)
	diags = append(diags, dgs...)
	state.PackageDetails = setVal
	state.ID = types.StringValue("placeholder")

	return state, diags
}

// GetPackageType returns the Package type required for mapping
func GetPackageType() map[string]attr.Type {
	return map[string]attr.Type{
		"file_name":        types.StringType,
		"operating_system": types.StringType,
		"linux_flavour":    types.StringType,
		"version":          types.StringType,
		"label":            types.StringType,
		"type":             types.StringType,
		"sio_patch_number": types.Int64Type,
		"size":             types.Int64Type,
		"latest":           types.BoolType,
	}
}

// GetPackageValue returns the Package object required for mapping
func GetPackageValue(packageDetails *goscaleio_types.PackageDetails) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetPackageType(), map[string]attr.Value{
		"file_name":        types.StringValue(packageDetails.Filename),
		"operating_system": types.StringValue(packageDetails.OperatingSystem),
		"linux_flavour":    types.StringValue(packageDetails.LinuxFlavour),
		"version":          types.StringValue(packageDetails.Version),
		"label":            types.StringValue(packageDetails.Label),
		"type":             types.StringValue(packageDetails.Type),
		"sio_patch_number": types.Int64Value(int64(packageDetails.SioPatchNumber)),
		"size":             types.Int64Value(int64(packageDetails.Size)),
		"latest":           types.BoolValue(bool(packageDetails.Latest)),
	})
}
