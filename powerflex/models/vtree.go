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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VTreeDataSourceModel maps the struct to VTree data source schema
type VTreeDataSourceModel struct {
	VTreeIDs    types.Set    `tfsdk:"vtree_ids"`
	VolumeIDs   types.Set    `tfsdk:"volume_ids"`
	VolumeNames types.Set    `tfsdk:"volume_names"`
	VTrees      []VTree      `tfsdk:"vtree_details"`
	ID          types.String `tfsdk:"id"`
}

// VTree maps the struct to VTree schema
type VTree struct {
	StoragePoolID      types.String       `tfsdk:"storage_pool_id"`
	DataLayout         types.String       `tfsdk:"data_layout"`
	CompressionMethod  types.String       `tfsdk:"compression_method"`
	RootVolumes        []types.String     `tfsdk:"root_volumes"`
	VtreeMigrationInfo VtreeMigrationInfo `tfsdk:"vtree_migration_info"`
	InDeletion         types.Bool         `tfsdk:"in_deletion"`
	Name               types.String       `tfsdk:"name"`
	ID                 types.String       `tfsdk:"id"`
	Links              []VTreeLinks       `tfsdk:"links"`
}

// VtreeMigrationInfo maps the struct to VTree migration schema
type VtreeMigrationInfo struct {
	MigrationQueuePosition   types.Int64  `tfsdk:"migration_queue_position"`
	MigrationPauseReason     types.String `tfsdk:"migration_pause_reason"`
	MigrationStatus          types.String `tfsdk:"migration_status"`
	SourceStoragePoolID      types.String `tfsdk:"source_storage_pool_id"`
	DestinationStoragePoolID types.String `tfsdk:"destination_storage_pool_id"`
	ThicknessConversionType  types.String `tfsdk:"thickness_conversion_type"`
}

// VTreeLinks maps the struct to VTree links
type VTreeLinks struct {
	Rel  types.String `tfsdk:"rel"`
	Href types.String `tfsdk:"href"`
}
