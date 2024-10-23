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
	"fmt"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
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

// setNvmeHostName sets the name of the NVMe host.
// If the name is empty, it will be set to the host type + ID
func setNvmeHostName(host *goscaleio_types.NvmeHost) {
	if len(host.Name) == 0 {
		host.Name = fmt.Sprintf("%s:%s", host.HostType, host.ID)
	}
}

// GetNvmeHostByID returns an NVMe host by id
func GetNvmeHostByID(system *goscaleio.System, id string) (*goscaleio_types.NvmeHost, error) {
	host, err := system.GetNvmeHostByID(id)
	if err != nil {
		return host, err
	}
	setNvmeHostName(host)
	return host, nil
}

// GetAllNvmeHosts returns all NvmeHost list
func GetAllNvmeHosts(system *goscaleio.System) ([]goscaleio_types.NvmeHost, error) {
	hosts, err := system.GetAllNvmeHosts()
	if err != nil {
		return hosts, err
	}

	for i := range hosts {
		setNvmeHostName(&hosts[i])
	}

	return hosts, nil
}
