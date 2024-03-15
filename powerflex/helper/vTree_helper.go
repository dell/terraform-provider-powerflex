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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllVTreeState saves state of vTree data source
func GetAllVTreeState(vTrees []scaleiotypes.VTreeDetails) (response []models.VTree) {
	for _, vTree := range vTrees {
		vTreeState := models.VTree{
			ID:                types.StringValue(vTree.ID),
			Name:              types.StringValue(vTree.Name),
			CompressionMethod: types.StringValue(vTree.CompressionMethod),
			StoragePoolID:     types.StringValue(vTree.StoragePoolID),
			DataLayout:        types.StringValue(vTree.DataLayout),
			InDeletion:        types.BoolValue(vTree.InDeletion),
		}

		for _, vol := range vTree.RootVolumes {
			vTreeState.RootVolumes = append(vTreeState.RootVolumes, types.StringValue(vol))
		}

		for _, link := range vTree.Links {
			vTreeState.Links = append(vTreeState.Links, models.VTreeLinks{
				Rel:  types.StringValue(link.Rel),
				Href: types.StringValue(link.HREF),
			})
		}

		vTreeState.VtreeMigrationInfo = models.VtreeMigrationInfo{
			MigrationQueuePosition:   types.Int64Value(vTree.VtreeMigrationInfo.MigrationQueuePosition),
			MigrationPauseReason:     types.StringValue(vTree.VtreeMigrationInfo.MigrationPauseReason),
			MigrationStatus:          types.StringValue(vTree.VtreeMigrationInfo.MigrationStatus),
			SourceStoragePoolID:      types.StringValue(vTree.VtreeMigrationInfo.SourceStoragePoolID),
			DestinationStoragePoolID: types.StringValue(vTree.VtreeMigrationInfo.DestinationStoragePoolID),
			ThicknessConversionType:  types.StringValue(vTree.VtreeMigrationInfo.ThicknessConversionType),
		}

		response = append(response, vTreeState)
	}
	return
}
