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

// UpdateWindowsMdms updates the mdms on an windows SDC host
func (r *SdcHostResource) UpdateWindowsMdms(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
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

	context["port"] = remote.Port

	winRMClient.GetConnection(context, false)

	defer winRMClient.Destroy()

	connectionStatus, err := winRMClient.Init()

	if err != nil {
		respDiagnostics.AddError(
			"Error while connecting sdc remote host",
			err.Error(),
		)
		return respDiagnostics
	}

	var mdms []string
	respDiagnostics = append(respDiagnostics, plan.MdmIPs.ElementsAs(ctx, &mdms, true)...)
	if respDiagnostics.HasError() {
		return respDiagnostics
	}

	if connectionStatus {
		ouptut, qErr := winRMClient.ExecuteCommand(fmt.Sprintf("cd '%s'; .\\drv_cfg.exe --query_mdms", plan.WindowsDrvCfg.ValueString()))
		if qErr != nil {
			respDiagnostics.AddError(
				"Error retrieving mdms",
				qErr.Error(),
			)
			return respDiagnostics
		}
		tflog.Info(ctx, "Checking if MDMS err: "+ouptut)
		for _, mdm := range mdms {
			splitMdms := strings.Split(mdm, ",")
			tflog.Info(ctx, "Checking if MDMS exists: "+mdm+" "+splitMdms[0])
			if strings.Contains(ouptut, splitMdms[0]) {
				tflog.Info(ctx, "Updating MDMS: "+mdm)
				_, err := winRMClient.ExecuteCommand(fmt.Sprintf("cd '%s'; .\\drv_cfg.exe --mod_mdm_ip --ip=%s --new_mdm_ip=%s", plan.WindowsDrvCfg.ValueString(), splitMdms[0], mdm))
				if err != nil {
					respDiagnostics.AddError(
						"Error updating mdms: "+mdm,
						err.Error(),
					)
					return respDiagnostics
				}
				tflog.Info(ctx, "MDMS updated")
				// Add new MDMs to the sdc
			} else {
				tflog.Info(ctx, "Adding MDMS: "+mdm)
				_, err := winRMClient.ExecuteCommand(fmt.Sprintf("cd '%s'; .\\drv_cfg.exe --add_mdm --ip=%s", plan.WindowsDrvCfg.ValueString(), mdm))
				if err != nil {
					respDiagnostics.AddError(
						"Error adding mdms?: "+mdm,
						err.Error(),
					)
					return respDiagnostics
				}
				tflog.Info(ctx, "MDMS Added")
			}
		}

		// Use the drv_cfg to validate the mdms are actually set
		for _, mdm := range mdms {
			splitMdms := strings.Split(mdm, ",")
			// Check to see if the mdms are actually set
			outputAfter, qErrAfter := winRMClient.ExecuteCommand(fmt.Sprintf("cd '%s'; (.\\drv_cfg.exe --query_mdms) -match '%s'", plan.WindowsDrvCfg.ValueString(), splitMdms[0]))
			if qErrAfter != nil {
				respDiagnostics.AddError(
					"Error validating mdms after they were set",
					qErrAfter.Error(),
				)
				return respDiagnostics
			}
			// If it contains 0000000000000000 that means the MDM was not valid
			if !strings.Contains(outputAfter, splitMdms[0]) || strings.Contains(outputAfter, "0000000000000000") {
				respDiagnostics.AddError(
					"Error validating mdms",
					"MDMS "+mdm+" were invalid please check the configuration and try again",
				)
				return respDiagnostics
			}
		}
		return respDiagnostics
	}

	respDiagnostics.AddError(
		"Error while connecting sdc remote host",
		"Error while connecting sdc remote host",
	)
	return respDiagnostics
}

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

	context["port"] = remote.Port

	winRMClient.GetConnection(context, false)

	defer winRMClient.Destroy()

	connectionStatus, err := winRMClient.Init()

	if err != nil {
		respDiagnostics.AddError(
			"Error while connecting sdc remote host",
			err.Error(),
		)
		return respDiagnostics
	}

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	if connectionStatus {

		ouptut, err := winRMClient.ExecuteCommand("Get-Package -name \"EMC-scaleio-sdc\" -ErrorAction SilentlyContinue")
		if err != nil {
			respDiagnostics.AddError(
				"Error while checking for installed sdc package",
				err.Error(),
			)
			return respDiagnostics
		}

		if ouptut == "SUCCESS" {
			if !plan.UseRemotePath.ValueBool() {
				err := winRMClient.Upload("C:\\EMC-ScaleIO-sdc.msi", plan.Pkg.ValueString())

				if err != nil {
					respDiagnostics.AddError(
						"Error while uploading package",
						err.Error(),
					)
					return respDiagnostics
				}
			}

			ouptut, err := winRMClient.ExecuteCommand("msiexec.exe /i \"C:\\EMC-ScaleIO-sdc.msi\" MDM_IP=\"" + strings.Join(mdmIPs, ",") + "\" /q")
			if err != nil {
				respDiagnostics.AddError(
					"Error while installing command",
					err.Error(),
				)
				return respDiagnostics
			}

			if ouptut == "SUCCESS" {

				time.Sleep(30 * time.Second)

				tflog.Info(ctx, "Installed SDC Package")
				// If more then one cluster is set, update the mdms to include all clusters
				if !plan.MdmIPs.IsUnknown() && len(plan.MdmIPs.Elements()) > 0 {
					tflog.Info(ctx, "Updating MDMs")
					r.UpdateWindowsMdms(ctx, plan)
				}

				return respDiagnostics

			}

		}

		respDiagnostics.AddError(
			"SDC Package is alredy installed",
			"SDC Package is alredy installed",
		)
		return respDiagnostics

	}

	respDiagnostics.AddError(
		"Error while connecting sdc remote host",
		"Error while connecting sdc remote host",
	)
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

	context["port"] = remote.Port

	winRMClient.GetConnection(context, false)

	defer winRMClient.Destroy()

	connectionStatus, err := winRMClient.Init()
	if err != nil {
		respDiagnostics.AddError(
			"Error while connecting sdc remote host",
			err.Error(),
		)
		return respDiagnostics
	}

	if connectionStatus {
		ouptut, err := winRMClient.ExecuteCommand("msiexec.exe /x \"C:\\EMC-ScaleIO-sdc.msi\" /q")
		if err != nil {
			respDiagnostics.AddError(
				"Error while uninstalling sdc package",
				err.Error(),
			)
			return respDiagnostics
		}

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
