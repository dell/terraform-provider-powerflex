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

	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// contains checks if the value is in the items slice
func contains(items []types.String, val string) bool {
	for _, item := range items {
		if item.ValueString() == val {
			return true
		}
	}
	return false
}

// isMatchingFilter checks if the host matches the filter
func isMatchingFilter(filter *models.IDNameFilter, host goscaleio_types.NvmeHost) bool {
	// filter not specified, return all hosts
	if filter == nil {
		return true
	}

	// Check if filter.IDs and filter.Names are both empty
	if len(filter.IDs) == 0 && len(filter.Names) == 0 {
		return true
	}

	// Check if filter.IDs or filter.Names contain the corresponding host attributes
	if (len(filter.IDs) > 0 && contains(filter.IDs, host.ID)) || (len(filter.Names) > 0 && contains(filter.Names, host.Name)) {
		return true
	}

	return false
}

// GetNvmeHostState sets the state for the NvmeHost datasource.
func GetNvmeHostState(hosts []goscaleio_types.NvmeHost, filter *models.IDNameFilter) []models.NvmeHostModel {
	var response []models.NvmeHostModel
	for _, host := range hosts {
		if isMatchingFilter(filter, host) {
			hostState := models.NvmeHostModel{
				ID:             types.StringValue(host.ID),
				Name:           types.StringValue(host.Name),
				Nqn:            types.StringValue(host.Nqn),
				SystemID:       types.StringValue(host.SystemID),
				MaxNumPaths:    types.Int64Value(int64(host.MaxNumPaths)),
				MaxNumSysPorts: types.Int64Value(int64(host.MaxNumSysPorts)),
			}
			for _, link := range host.Links {
				hostState.Links = append(hostState.Links, models.LinkModel{
					Rel:  types.StringValue(link.Rel),
					HREF: types.StringValue(link.HREF),
				})
			}

			response = append(response, hostState)
		}
	}
	return response
}