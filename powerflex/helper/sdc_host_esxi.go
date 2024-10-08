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
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateEsxiMdms updates an esxi SDC host
func (r *SdcHostResource) UpdateEsxiMdms(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	var esxiInput models.SdcHostEsxiModel

	respDiagnostics.Append(plan.Esxi.As(ctx, &esxiInput, basetypes.ObjectAsOptions{})...)
	if respDiagnostics.HasError() {
		return respDiagnostics
	}

	sshP, _, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	esxi := client.NewEsxCli(sshP)
	// update mdms
	respDiagnostics = r.updateMdms(ctx, plan, esxiInput, esxi, sshP)

	return respDiagnostics
}

// CreateEsxi creates an esxi SDC host
func (r *SdcHostResource) CreateEsxi(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics

	var esxiInput models.SdcHostEsxiModel
	respDiagnostics.Append(plan.Esxi.As(ctx, &esxiInput, basetypes.ObjectAsOptions{})...)
	if respDiagnostics.HasError() {
		return respDiagnostics
	}

	sshP, dir, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	pkgTarget := strings.TrimSuffix(dir, "/") + "/" + "emc-sdc-package.zip"
	if !plan.UseRemotePath.ValueBool() {
		// upload sw
		scpProv := client.NewScpProvisioner(sshP)
		err = scpProv.Upload(plan.Pkg.ValueString(), pkgTarget, "")
		if err != nil {
			respDiagnostics.AddError(
				"Error uploading package",
				err.Error(),
			)
			return respDiagnostics
		}
	}

	// install sw
	esxi := client.NewEsxCli(sshP)
	pkgInstallCmd := client.VibInstallCommand{
		ZipFile:  pkgTarget,
		SigCheck: esxiInput.VerifyVibSign.ValueBool(),
	}
	op, err := esxi.SoftwareInstall(pkgInstallCmd)
	if err != nil {
		respDiagnostics.AddError(
			"Error installing package",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return respDiagnostics
	}

	// check sw
	tflog.Info(ctx, "Checking for installed sdc package")
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		respDiagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %s", sdc))

	// update mdms
	respDiagnostics = r.updateMdms(ctx, plan, esxiInput, esxi, sshP)

	return respDiagnostics
}

// updateMdms - function to update MDMs
func (r *SdcHostResource) updateMdms(ctx context.Context, plan models.SdcHostModel, esxiInput models.SdcHostEsxiModel, esxi *client.EsxCli, sshP *client.SSHProvisioner) diag.Diagnostics {
	tflog.Info(ctx, "Setting scini module parameters")
	var respDiagnostics diag.Diagnostics

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	// Multi cluster
	if !plan.MdmIPs.IsUnknown() && len(plan.MdmIPs.Elements()) > 0 {
		var mdms []string
		respDiagnostics = append(respDiagnostics, plan.MdmIPs.ElementsAs(ctx, &mdms, true)...)
		params := map[string]string{
			"IoctlIniGuidStr": esxiInput.GUID.ValueString(),
			// for multi cluster this is a comma seperated list with a `+` inbetween clusters
			"IoctlMdmIPStr":          strings.Join(mdms, "+"),
			"bBlkDevIsPdlActive":     "1",
			"blkDevPdlTimeoutMillis": "60000",
		}
		if op, err := esxi.SetModuleParameters("scini", params); err != nil {
			tflog.Debug(ctx, op)
			respDiagnostics.AddError(
				"Error setting module parameters",
				err.Error()+"\n"+op,
			)
			return respDiagnostics
		}
		tflog.Info(ctx, "Scini module parameters set")

		// Single cluster
	} else {
		params := map[string]string{
			"IoctlIniGuidStr":        esxiInput.GUID.ValueString(),
			"IoctlMdmIPStr":          strings.Join(mdmIPs, ","),
			"bBlkDevIsPdlActive":     "1",
			"blkDevPdlTimeoutMillis": "60000",
		}
		if op, err := esxi.SetModuleParameters("scini", params); err != nil {
			tflog.Debug(ctx, op)
			respDiagnostics.AddError(
				"Error setting module parameters",
				err.Error()+"\n"+op,
			)
			return respDiagnostics
		}
		tflog.Info(ctx, "Scini module parameters set")
	}

	// reboot
	errReboot := sshP.RebootUnix()
	if errReboot != nil {
		respDiagnostics.AddError(
			"Error rebooting",
			errReboot.Error(),
		)
		return respDiagnostics
	}

	// list esxi kernel modules
	tflog.Info(ctx, "Checking for installed scini module")
	op, err := esxi.GetModuleByName("scini")
	if err != nil {
		tflog.Info(ctx, "Scini module not found: "+op)
		respDiagnostics.AddError(
			"Scini module not found after second reboot",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Scini module found: "+op)

	// wait 30 seconds for SDC to connect properly
	time.Sleep(30 * time.Second)

	return respDiagnostics
}

// DeleteEsxi - function to uninstall SDC package in ESXi host
func (r *SdcHostResource) DeleteEsxi(ctx context.Context, state models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")
	sshP, _, err := r.getSSHProvisioner(ctx, state)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	tflog.Info(ctx, "Checking for installed sdc package")
	esxi := client.NewEsxCli(sshP)
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		respDiagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %v+", sdc))

	op, err := esxi.SoftwareRmv(sdc.Name)
	if err != nil {
		respDiagnostics.AddError(
			"Error removing sdc package",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("sdc package removed: %s", op))

	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting host",
			err.Error(),
		)
		return respDiagnostics
	}

	return respDiagnostics
}
