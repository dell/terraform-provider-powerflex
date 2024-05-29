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

// FirmwareRepositoryDatasourceModel defines the tfsdk model for firmware repository datasource
type FirmwareRepositoryDatasourceModel struct {
	FirmwareRepositoryIDs     types.Set                   `tfsdk:"firmware_repository_ids"`
	FirmwareRepositoryNames   types.Set                   `tfsdk:"firmware_repository_names"`
	FirmwareRepositoryDetails []FirmwareRepositoryDetails `tfsdk:"firmware_repository_details"`
	ID                        types.String                `tfsdk:"id"`
}

// FirmwareRepositoryDetails defines the tfsdk model of firmware repository details
type FirmwareRepositoryDetails struct {
	ID                  types.String  `tfsdk:"id"`
	Name                types.String  `tfsdk:"name"`
	SourceLocation      types.String  `tfsdk:"source_location"`
	SourceType          types.String  `tfsdk:"source_type"`
	DiskLocation        types.String  `tfsdk:"disk_location"`
	Filename            types.String  `tfsdk:"filename"`
	Username            types.String  `tfsdk:"username"`
	Password            types.String  `tfsdk:"password"`
	DownloadStatus      types.String  `tfsdk:"download_status"`
	CreatedDate         types.String  `tfsdk:"created_date"`
	CreatedBy           types.String  `tfsdk:"created_by"`
	UpdatedDate         types.String  `tfsdk:"updated_date"`
	UpdatedBy           types.String  `tfsdk:"updated_by"`
	DefaultCatalog      types.Bool    `tfsdk:"default_catalog"`
	Embedded            types.Bool    `tfsdk:"embedded"`
	State               types.String  `tfsdk:"state"`
	SoftwareComponents  []Component   `tfsdk:"software_components"`
	SoftwareBundles     []Bundle      `tfsdk:"software_bundles"`
	BundleCount         types.Int64   `tfsdk:"bundle_count"`
	ComponentCount      types.Int64   `tfsdk:"component_count"`
	UserBundleCount     types.Int64   `tfsdk:"user_bundle_count"`
	Minimal             types.Bool    `tfsdk:"minimal"`
	DownloadProgress    types.Int64   `tfsdk:"download_progress"`
	ExtractProgress     types.Int64   `tfsdk:"extract_progress"`
	FileSizeInGigabytes types.Float64 `tfsdk:"file_size_in_gigabytes"`
	Signature           types.String  `tfsdk:"signature"`
	Custom              types.Bool    `tfsdk:"custom"`
	NeedsAttention      types.Bool    `tfsdk:"needs_attention"`
	JobID               types.String  `tfsdk:"job_id"`
	Rcmapproved         types.Bool    `tfsdk:"rcmapproved"`
}

// Component is the tfsdk model of Component
type Component struct {
	ID                  types.String   `tfsdk:"id"`
	PackageID           types.String   `tfsdk:"package_id"`
	DellVersion         types.String   `tfsdk:"dell_version"`
	VendorVersion       types.String   `tfsdk:"vendor_version"`
	ComponentID         types.String   `tfsdk:"component_id"`
	DeviceID            types.String   `tfsdk:"device_id"`
	SubDeviceID         types.String   `tfsdk:"sub_device_id"`
	VendorID            types.String   `tfsdk:"vendor_id"`
	SubVendorID         types.String   `tfsdk:"sub_vendor_id"`
	CreatedDate         types.String   `tfsdk:"created_date"`
	CreatedBy           types.String   `tfsdk:"created_by"`
	UpdatedDate         types.String   `tfsdk:"updated_date"`
	UpdatedBy           types.String   `tfsdk:"updated_by"`
	Path                types.String   `tfsdk:"path"`
	HashMd5             types.String   `tfsdk:"hash_md5"`
	Name                types.String   `tfsdk:"name"`
	Category            types.String   `tfsdk:"category"`
	ComponentType       types.String   `tfsdk:"component_type"`
	OperatingSystem     types.String   `tfsdk:"operating_system"`
	SystemIDs           []types.String `tfsdk:"system_ids"`
	Custom              types.Bool     `tfsdk:"custom"`
	NeedsAttention      types.Bool     `tfsdk:"needs_attention"`
	Ignore              types.Bool     `tfsdk:"ignore"`
	OriginalComponentID types.String   `tfsdk:"original_component_id"`
	FirmwareRepoName    types.String   `tfsdk:"firmware_repo_name"`
}

// Bundle is the tfsdk model of Bundle
type Bundle struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Version            types.String `tfsdk:"version"`
	BundleDate         types.String `tfsdk:"bundle_date"`
	CreatedDate        types.String `tfsdk:"created_date"`
	CreatedBy          types.String `tfsdk:"created_by"`
	UpdatedDate        types.String `tfsdk:"updated_date"`
	UpdatedBy          types.String `tfsdk:"updated_by"`
	Description        types.String `tfsdk:"description"`
	UserBundle         types.Bool   `tfsdk:"user_bundle"`
	UserBundlePath     types.String `tfsdk:"user_bundle_path"`
	DeviceType         types.String `tfsdk:"device_type"`
	DeviceModel        types.String `tfsdk:"device_model"`
	FwRepositoryID     types.String `tfsdk:"fw_repository_id"`
	BundleType         types.String `tfsdk:"bundle_type"`
	Custom             types.Bool   `tfsdk:"custom"`
	NeedsAttention     types.Bool   `tfsdk:"needs_attention"`
	SoftwareComponents []Component  `tfsdk:"software_components"`
}
