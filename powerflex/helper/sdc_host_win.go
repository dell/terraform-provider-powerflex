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
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	byteData []byte
)

// CreateWindows creates an windows SDC host
func (r *SdcHostResource) CreateWindows(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics

	var remote models.SdcHostRemoteModel
	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	winRMClient := &client.WinRMClient{}

	var contexts map[string]string

	_ = json.Unmarshal(byteData, &contexts)

	context := make(map[string]string)

	context["username"] = remote.User

	context["password"] = *remote.Password

	context["host"] = plan.Host.ValueString()

	winRMClient.GetConnection(context, false)

	defer winRMClient.Destroy()

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	if winRMClient.Init() {

		ouptut := winRMClient.ExecuteCommand("Get-Package -name \"EMC-scaleio-sdc\" -ErrorAction SilentlyContinue")

		if ouptut == "SUCCESS" {
			winRMClient.Upload("C:\\EMC-ScaleIO-sdc.msi", plan.Pkg.ValueString())

			ouptut := winRMClient.ExecuteCommand("msiexec.exe /i \"C:\\EMC-ScaleIO-sdc.msi\" MDM_IP=\"" + strings.Join(mdmIPs, ",") + "\" /q")

			if ouptut == "SUCCESS" {

				time.Sleep(30 * time.Second)

				tflog.Info(ctx, "Installed SDC Package")

				return respDiagnostics

			}

			respDiagnostics.AddError(
				"Error while installing command",
				winRMClient.Errors[0]["message"],
			)
			return respDiagnostics

		}

		respDiagnostics.AddError(
			"SDC Package is alredy installed",
			"SDC Package is alredy installed",
		)
		return respDiagnostics

	}

	return respDiagnostics
}

// DeleteWindows - function to uninstall SDC package in Windows host
func (r *SdcHostResource) DeleteWindows(ctx context.Context, state models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")

	var remote models.SdcHostRemoteModel
	state.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	winRMClient := &client.WinRMClient{}

	var contexts map[string]string

	_ = json.Unmarshal(byteData, &contexts)

	context := make(map[string]string)

	context["username"] = remote.User

	context["password"] = *remote.Password

	context["host"] = state.Host.ValueString()

	winRMClient.GetConnection(context, false)

	defer winRMClient.Destroy()

	if winRMClient.Init() {
		ouptut := winRMClient.ExecuteCommand("msiexec.exe /x \"C:\\EMC-ScaleIO-sdc.msi\" /q")

		if ouptut == "SUCCESS" {

			time.Sleep(10 * time.Second)

			tflog.Info(ctx, "Uninstalled SDC Package")

		} else {

			respDiagnostics.AddError(
				"Error while installing command",
				winRMClient.Errors[0]["message"],
			)
			return respDiagnostics
		}
	}

	return respDiagnostics
}
