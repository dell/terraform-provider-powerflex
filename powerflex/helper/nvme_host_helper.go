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

// GetNvmeHostState sets the state for the NvmeHost datasource.
func GetNvmeHostState(host goscaleio_types.NvmeHost) models.NvmeHostDatasourceModel {
	model := models.NvmeHostDatasourceModel{
		ID:             types.StringValue(host.ID),
		Name:           types.StringValue(host.Name),
		Nqn:            types.StringValue(host.Nqn),
		SystemID:       types.StringValue(host.SystemID),
		MaxNumPaths:    types.Int64Value(int64(host.MaxNumPaths)),
		MaxNumSysPorts: types.Int64Value(int64(host.MaxNumSysPorts)),
	}
	for _, link := range host.Links {
		model.Links = append(model.Links, models.LinkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}
	return model
}
