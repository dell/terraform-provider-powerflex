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
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"terraform-provider-powerflex/powerflex/client/winrmclient"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	byteData []byte
)

func HostWindowsOperations(ctx context.Context, apiClient goscaleio.Client, plan models.HostResourceModel, mdmIP []string, credential models.CredentialModel, system *goscaleio.System) ([]models.HostDetailModel, error) {

	hostDetailList := []models.HostDetailModel{}

	winRMClient := &winrmclient.WinRMClient{}

	var contexts map[string]string

	_ = json.Unmarshal(byteData, &contexts)

	context := make(map[string]string)

	context["username"] = credential.UserName.ValueString()

	context["password"] = credential.Password.ValueString()

	context["host"] = plan.IP.ValueString()

	winRMClient.GetConnection(context, false)

	if winRMClient.Init() {

		winRMClient.Upload("C:\\EMC-ScaleIO-sdc.msi", plan.PackagePath.ValueString())

		ouptut := winRMClient.ExecuteCommand("msiexec.exe /i \"C:\\EMC-ScaleIO-sdc.msi\" MDM_IP=\"" + strings.Join(mdmIP, ",") + "\" /q")

		if ouptut == "SUCCESS" {

			time.Sleep(30 * time.Second)

			sdcData, _ := system.FindSdc("SdcIP", plan.IP.ValueString())

			hostDetailList = append(hostDetailList, setSDCDetails(*sdcData.Sdc))

			return hostDetailList, nil

		} else {
			return nil, fmt.Errorf("Error while installing command %s", winRMClient.Errors[0]["message"])
		}

	} else {
		return nil, fmt.Errorf("Error while establishing connection %s", winRMClient.Errors[0]["message"])
	}

}

func setSDCDetails(sdc goscaleio_types.Sdc) models.HostDetailModel {

	var model models.HostDetailModel

	model.OperatingSystem = types.StringValue(sdc.OSType)
	model.HostID = types.StringValue(sdc.ID)
	model.HostName = types.StringValue(sdc.Name)
	model.HostGUID = types.StringValue(sdc.SdcGUID)
	model.IsApproved = types.BoolValue(sdc.SdcApproved)
	model.OnVMWare = types.BoolValue(sdc.OnVMWare)
	model.SystemID = types.StringValue(sdc.SystemID)
	model.PerformanceProfile = types.StringValue(sdc.PerfProfile)
	model.IP = types.StringValue(sdc.SdcIP)
	model.MdmConnectionState = types.StringValue(sdc.MdmConnectionState)
	return model
}

func UpdateHostState(hosts []models.HostDetailModel, plan models.HostResourceModel) (models.HostResourceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	HostAttrTypes := GetHostStateDetailType()

	HostElemType := types.ObjectType{
		AttrTypes: HostAttrTypes,
	}

	objectHosts := []attr.Value{}
	for _, host := range hosts {
		objVal, dgs := GetHostStateDetailValue(host)
		diags = append(diags, dgs...)
		objectHosts = append(objectHosts, objVal)
	}
	setHosts, dgs := types.ListValue(HostElemType, objectHosts)
	diags = append(diags, dgs...)

	state.HostDetails = setHosts
	state.ID = types.StringValue("placeholder")

	return state, diags
}

func GetHostStateDetailType() map[string]attr.Type {
	return map[string]attr.Type{
		"host_id":              types.StringType,
		"ip":                   types.StringType,
		"operating_system":     types.StringType,
		"performance_profile":  types.StringType,
		"host_name":            types.StringType,
		"system_id":            types.StringType,
		"is_approved":          types.BoolType,
		"on_vmware":            types.BoolType,
		"host_guid":            types.StringType,
		"mdm_connection_state": types.StringType,
	}
}

func GetHostStateDetailValue(host models.HostDetailModel) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(GetHostStateDetailType(), map[string]attr.Value{
		"ip":                   types.StringValue(host.IP.ValueString()),
		"operating_system":     types.StringValue(host.OperatingSystem.ValueString()),
		"performance_profile":  types.StringValue(host.PerformanceProfile.ValueString()),
		"host_name":            types.StringValue(host.HostName.ValueString()),
		"host_id":              types.StringValue(host.HostID.ValueString()),
		"system_id":            types.StringValue(host.SystemID.ValueString()),
		"is_approved":          types.BoolValue(host.IsApproved.ValueBool()),
		"on_vmware":            types.BoolValue(host.OnVMWare.ValueBool()),
		"host_guid":            types.StringValue(host.HostGUID.ValueString()),
		"mdm_connection_state": types.StringValue(host.MdmConnectionState.ValueString()),
	})
}
