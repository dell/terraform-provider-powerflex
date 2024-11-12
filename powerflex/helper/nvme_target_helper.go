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

// GetNvmeTargetByID returns an nvme target searched by id
func GetNvmeTargetByID(system *goscaleio.System, id string) (*goscaleio_types.Sdt, error) {
	sdt, err := system.GetSdtByID(id)
	if err != nil {
		return sdt, err
	}

	if sdt.MaintenanceState == "NoMaintenance" || sdt.MaintenanceState == "ExitMaintenanceInProgress" {
		sdt.MaintenanceState = "Inactive"
	} else if sdt.MaintenanceState == "InMaintenance" || sdt.MaintenanceState == "SetMaintenanceInProgress" {
		sdt.MaintenanceState = "Active"
	}

	return sdt, nil
}

// NvmeTargetUpdate updates the NVMe target
func NvmeTargetUpdate(system *goscaleio.System, state, plan models.NvmeTargetResourceModel) error {
	if state.Name.ValueString() != plan.Name.ValueString() {
		err := system.RenameSdt(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			return fmt.Errorf("could not rename the NVMe target: %w", err)
		}
	}

	if !plan.DiscoveryPort.IsUnknown() && !plan.DiscoveryPort.IsNull() &&
		state.DiscoveryPort.ValueInt64() != plan.DiscoveryPort.ValueInt64() {
		err := system.SetSdtDiscoveryPort(state.ID.ValueString(), int(plan.DiscoveryPort.ValueInt64()))
		if err != nil {
			return fmt.Errorf("could not update discovery port: %w", err)
		}
	}

	if !plan.NvmePort.IsUnknown() && !plan.NvmePort.IsNull() &&
		state.NvmePort.ValueInt64() != plan.NvmePort.ValueInt64() {
		err := system.SetSdtNvmePort(state.ID.ValueString(), int(plan.NvmePort.ValueInt64()))
		if err != nil {
			return fmt.Errorf("could not update NVMe port: %w", err)
		}
	}

	if !plan.StoragePort.IsUnknown() && !plan.StoragePort.IsNull() &&
		state.StoragePort.ValueInt64() != plan.StoragePort.ValueInt64() {
		err := system.SetSdtStoragePort(state.ID.ValueString(), int(plan.StoragePort.ValueInt64()))
		if err != nil {
			return fmt.Errorf("could not update storage port: %w", err)
		}
	}

	added, removed, changed := CompareIPLists(state.IPList, plan.IPList)

	for _, ip := range added {
		err := system.AddSdtTargetIP(state.ID.ValueString(), ip.IP.ValueString(), ip.Role.ValueString())
		if err != nil {
			return fmt.Errorf("could not add target IP %s: %w", ip.IP.ValueString(), err)
		}
	}

	for _, ip := range removed {
		err := system.RemoveSdtTargetIP(state.ID.ValueString(), ip.IP.ValueString())
		if err != nil {
			return fmt.Errorf("could not remove target IP %s: %w", ip.IP.ValueString(), err)
		}
	}

	for _, ip := range changed {
		err := system.ModifySdtIPRole(state.ID.ValueString(), ip.IP.ValueString(), ip.Role.ValueString())
		if err != nil {
			return fmt.Errorf("could not update the role of IP %s to %s: %w", ip.IP.ValueString(), ip.Role.ValueString(), err)
		}
	}

	return nil
}

// CompareIPLists takes state and plan IPLists, and returns added, removed and changed entries
func CompareIPLists(stateIPs, planIPs []models.IPList) (added, removed, changed []models.IPList) {
	stateMap := make(map[string]models.IPList)
	planMap := make(map[string]models.IPList)

	// Fill stateMap and planMap with current state and expected plan
	for _, ip := range stateIPs {
		stateMap[ip.IP.ValueString()] = ip
	}
	for _, ip := range planIPs {
		planMap[ip.IP.ValueString()] = ip
	}

	// Compare state with plan to find added, removed, and changed entries
	for _, ip := range stateIPs {
		if _, exists := planMap[ip.IP.ValueString()]; !exists {
			removed = append(removed, ip)
		}
	}
	for _, ip := range planIPs {
		if existing, exists := stateMap[ip.IP.ValueString()]; !exists {
			added = append(added, ip)
		} else if existing.Role != ip.Role {
			changed = append(changed, ip)
		}
	}
	return
}

// ToggleSdtMaintenanceMode toggles sdt maintenance mode
func ToggleSdtMaintenanceMode(system *goscaleio.System, id, state string) error {
	var err error
	if state == "Active" {
		err = system.EnterSdtMaintenanceMode(id)
	} else {
		err = system.ExitSdtMaintenanceMode(id)
	}
	return err
}
