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

func UpdateFrimwareRepositoryState(frDetails *scaleiotypes.UploadComplianceTopologyDetails, plan models.FirmwareRepositoryResourceModel) models.FirmwareRepositoryResourceModel {
	state := plan
	state.ID = types.StringValue(frDetails.ID)
	state.DiskLocation = types.StringValue(frDetails.DiskLocation)
	state.FileName = types.StringValue(frDetails.Filename)
	state.Name = types.StringValue(frDetails.Name)
	state.DefaultCatalog = types.BoolValue(frDetails.DefaultCatalog)
	return state
}
