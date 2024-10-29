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

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetPeerMdms gets all peer mdm
func GetPeerMdms(client *goscaleio.Client) ([]scaleiotypes.PeerMDM, error) {
	peerMdmList := []scaleiotypes.PeerMDM{}

	// Get All Peer Mdm
	peerMdms, err := client.GetPeerMDMs()
	if err != nil {
		return nil, err
	}
	for _, val := range peerMdms {
		peerMdmList = append(peerMdmList, *val)
	}

	return peerMdmList, nil
}

// MapPeerMdmState maps peer mdm state
func MapPeerMdmState(pairs []scaleiotypes.PeerMDM, state models.PeerMdmDataSourceModel) models.PeerMdmDataSourceModel {
	mappedPeerMdm := []models.PeerMdmModel{}
	for _, val := range pairs {

		temp := models.PeerMdmModel{
			ID:                  types.StringValue(val.ID),
			Name:                types.StringValue(val.Name),
			Port:                types.Int64Value(int64(val.Port)),
			PeerSystemID:        types.StringValue(val.PeerSystemID),
			SystemID:            types.StringValue(val.SystemID),
			SoftwareVersionInfo: types.StringValue(val.SoftwareVersionInfo),
			MembershipState:     types.StringValue(val.MembershipState),
			PerfProfile:         types.StringValue(val.PerfProfile),
			NetworkType:         types.StringValue(val.NetworkType),
			CouplingRC:          types.StringValue(val.CouplingRC),
			IPList:              []*models.IPListNoRole{},
		}

		for _, IP := range val.IPList {
			temp.IPList = append(temp.IPList, &models.IPListNoRole{
				IP: types.StringValue(IP.IP),
			})
		}

		mappedPeerMdm = append(mappedPeerMdm, temp)
	}
	return models.PeerMdmDataSourceModel{
		ID:             types.StringValue("peer_system_id"),
		PeerMdmFilter:  state.PeerMdmFilter,
		PeerMdmDetails: mappedPeerMdm,
	}
}
