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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetAttachedNvmeHostInfo sets the NVMe hosts attached to the NVMe targets
func SetAttachedNvmeHostInfo(system *goscaleio.System, nvmeTargets []models.NvmeTargetDatasourceModel) error {
	allNvmeHosts, err := GetAllNvmeHosts(system)
	if err != nil {
		return err
	}

	for _, host := range allNvmeHosts {
		controllers, err := system.GetHostNvmeControllers(host)
		if err != nil {
			continue
		}
		for _, controller := range controllers {
			if !controller.IsConnected {
				continue
			}

			for i := range nvmeTargets {
				if controller.SdtID == nvmeTargets[i].ID.ValueString() {
					h := models.NvmeHostList{
						HostIP:      types.StringValue(controller.HostIP),
						IsConnected: types.BoolValue(controller.IsConnected),
						HostName:    types.StringValue(host.Name),
						HostID:      types.StringValue(controller.HostID),
						SysPortIP:   types.StringValue(controller.SysPortIP),
					}
					nvmeTargets[i].HostList = append(nvmeTargets[i].HostList, h)
				}
			}
		}
	}
	return nil
}
