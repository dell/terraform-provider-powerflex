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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateFrimwareRepositoryState sets the state for the firmware repository resource.
func UpdateFrimwareRepositoryState(frDetails *scaleiotypes.UploadComplianceTopologyDetails, plan models.FirmwareRepositoryResourceModel) models.FirmwareRepositoryResourceModel {
	state := plan
	state.ID = types.StringValue(frDetails.ID)
	state.DiskLocation = types.StringValue(frDetails.DiskLocation)
	state.FileName = types.StringValue(frDetails.Filename)
	state.Name = types.StringValue(frDetails.Name)
	state.DefaultCatalog = types.BoolValue(frDetails.DefaultCatalog)
	state.SourceLocation = types.StringValue(frDetails.SourceLocation)
	if plan.Approve.IsNull() && frDetails.State == "needsApproval" {
		state.Approve = types.BoolValue(false)
	} else if plan.Approve.IsNull() && frDetails.State == "available" {
		state.Approve = types.BoolValue(true)
	}
	return state
}

// GetAllFirmwareRepositoryState sets the state for the firmware repository datasource.
func GetAllFirmwareRepositoryState(input *scaleiotypes.FirmwareRepositoryDetails) models.FirmwareRepositoryDetails {
	return models.FirmwareRepositoryDetails{
		ID:                  types.StringValue(input.ID),
		Name:                types.StringValue(input.Name),
		SourceLocation:      types.StringValue(input.SourceLocation),
		SourceType:          types.StringValue(input.SourceType),
		DiskLocation:        types.StringValue(input.DiskLocation),
		Filename:            types.StringValue(input.Filename),
		Username:            types.StringValue(input.Username),
		Password:            types.StringValue(input.Password),
		DownloadStatus:      types.StringValue(input.DownloadStatus),
		CreatedDate:         types.StringValue(input.CreatedDate),
		CreatedBy:           types.StringValue(input.CreatedBy),
		UpdatedDate:         types.StringValue(input.UpdatedDate),
		UpdatedBy:           types.StringValue(input.UpdatedBy),
		DefaultCatalog:      types.BoolValue(input.DefaultCatalog),
		Embedded:            types.BoolValue(input.Embedded),
		State:               types.StringValue(input.State),
		SoftwareComponents:  newComponentList(input.SoftwareComponents),
		SoftwareBundles:     newBundleList(input.SoftwareBundles),
		BundleCount:         types.Int64Value(int64(input.BundleCount)),
		ComponentCount:      types.Int64Value(int64(input.ComponentCount)),
		UserBundleCount:     types.Int64Value(int64(input.UserBundleCount)),
		Minimal:             types.BoolValue(input.Minimal),
		DownloadProgress:    types.Int64Value(int64(input.DownloadProgress)),
		ExtractProgress:     types.Int64Value(int64(input.ExtractProgress)),
		FileSizeInGigabytes: types.Float64Value(input.FileSizeInGigabytes),
		Signature:           types.StringValue(input.Signature),
		Custom:              types.BoolValue(input.Custom),
		NeedsAttention:      types.BoolValue(input.NeedsAttention),
		JobID:               types.StringValue(input.JobID),
		Rcmapproved:         types.BoolValue(input.Rcmapproved),
	}
}

// newComponentList converts list of client.Component to list of models.Component
func newComponentList(inputs []scaleiotypes.Component) []models.Component {
	out := make([]models.Component, 0)
	for _, input := range inputs {
		out = append(out, newComponent(input))
	}
	return out
}

// newBundleList converts list of client.Bundle to list of models.Bundle
func newBundleList(inputs []scaleiotypes.Bundle) []models.Bundle {
	out := make([]models.Bundle, 0)
	for _, input := range inputs {
		out = append(out, newBundle(input))
	}
	return out
}

// newComponent converts client.Component to models.Component
func newComponent(input scaleiotypes.Component) models.Component {
	return models.Component{
		ID:                  types.StringValue(input.ID),
		PackageID:           types.StringValue(input.PackageID),
		DellVersion:         types.StringValue(input.DellVersion),
		VendorVersion:       types.StringValue(input.VendorVersion),
		ComponentID:         types.StringValue(input.ComponentID),
		DeviceID:            types.StringValue(input.DeviceID),
		SubDeviceID:         types.StringValue(input.SubDeviceID),
		VendorID:            types.StringValue(input.VendorID),
		SubVendorID:         types.StringValue(input.SubVendorID),
		CreatedDate:         types.StringValue(input.CreatedDate),
		CreatedBy:           types.StringValue(input.CreatedBy),
		UpdatedDate:         types.StringValue(input.UpdatedDate),
		UpdatedBy:           types.StringValue(input.UpdatedBy),
		Path:                types.StringValue(input.Path),
		HashMd5:             types.StringValue(input.HashMd5),
		Name:                types.StringValue(input.Name),
		Category:            types.StringValue(input.Category),
		ComponentType:       types.StringValue(input.ComponentType),
		OperatingSystem:     types.StringValue(input.OperatingSystem),
		Custom:              types.BoolValue(input.Custom),
		NeedsAttention:      types.BoolValue(input.NeedsAttention),
		Ignore:              types.BoolValue(input.Ignore),
		OriginalComponentID: types.StringValue(input.OriginalComponentID),
		FirmwareRepoName:    types.StringValue(input.FirmwareRepoName),
		SystemIDs:           newSystemIDS(input.SystemIDs),
	}
}

// newBundle converts client.Bundle to models.Bundle
func newBundle(input scaleiotypes.Bundle) models.Bundle {
	return models.Bundle{
		ID:                 types.StringValue(input.ID),
		Name:               types.StringValue(input.Name),
		Version:            types.StringValue(input.Version),
		BundleDate:         types.StringValue(input.BundleDate),
		CreatedDate:        types.StringValue(input.CreatedDate),
		CreatedBy:          types.StringValue(input.CreatedBy),
		UpdatedDate:        types.StringValue(input.UpdatedDate),
		UpdatedBy:          types.StringValue(input.UpdatedBy),
		Description:        types.StringValue(input.Description),
		UserBundle:         types.BoolValue(input.UserBundle),
		UserBundlePath:     types.StringValue(input.UserBundlePath),
		DeviceType:         types.StringValue(input.DeviceType),
		DeviceModel:        types.StringValue(input.DeviceModel),
		FwRepositoryID:     types.StringValue(input.FwRepositoryID),
		BundleType:         types.StringValue(input.BundleType),
		Custom:             types.BoolValue(input.Custom),
		NeedsAttention:     types.BoolValue(input.NeedsAttention),
		SoftwareComponents: newComponentList(input.SoftwareComponents),
	}
}

func newSystemIDS(input []string) []types.String {
	var out []types.String
	for _, rspl := range input {
		out = append(out, types.StringValue(rspl))
	}
	return out
}
