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

// UpdateFaultSetState updates the State for Fault set Resource
func UpdateFaultSetState(faultset *scaleiotypes.FaultSet, plan models.FaultSetResourceModel) models.FaultSetResourceModel {
	state := plan
	state.ProtectionDomainID = types.StringValue(faultset.ProtectionDomainID)
	state.ID = types.StringValue(faultset.ID)
	state.Name = types.StringValue(faultset.Name)
	return state
}

// GetAllFaultSetState returns the state for fault set data source
func GetAllFaultSetState(faultSet scaleiotypes.FaultSet, sdsDetails []models.SdsDataModel) (response models.FaultSet) {
	response = models.FaultSet{
		ID:                 faultSet.ID,
		Name:               faultSet.Name,
		ProtectionDomainID: faultSet.ProtectionDomainID,
	}

	for _, link := range faultSet.Links {
		response.Links = append(response.Links, &models.LinkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}
	response.SdsDetails = sdsDetails
	return
}
