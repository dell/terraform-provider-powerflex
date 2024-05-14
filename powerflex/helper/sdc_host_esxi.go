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
	"path/filepath"
	"regexp"
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateEsxi creates an esxi SDC host
func (r *SdcHostResource) CreateEsxi(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	sshP, dir, err := r.getSshProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	// upload sw
	scpProv := client.NewScpProvisioner(sshP)
	pkgTarget := filepath.Join(dir, "emc-sdc-package.zip")
	err = scpProv.Upload(plan.Pkg.ValueString(), pkgTarget, "")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
		return respDiagnostics
	}

	// install sw
	esxi := client.NewEsxCli(sshP)
	pkgInstallCmd := client.VibInstallCommand{
		ZipFile:  pkgTarget,
		SigCheck: true,
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

	tflog.Info(ctx, "Setting scini module parameters")
	var esxiInput models.SdcHostEsxiModel
	respDiagnostics.Append(plan.Esxi.As(ctx, &esxiInput, basetypes.ObjectAsOptions{})...)
	if respDiagnostics.HasError() {
		return respDiagnostics
	}

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	params := map[string]string{
		"IoctlIniGuidStr":        esxiInput.Guid.ValueString(),
		"IoctlMdmIPStr":          strings.Join(mdmIPs, ","),
		"bBlkDevIsPdlActive":     "1",
		"blkDevPdlTimeoutMillis": "60000",
	}
	if op, err := esxi.SetModuleParameters("scini", params); err != nil {
		respDiagnostics.AddError(
			"Error setting module parameters",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Scini module parameters set")

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return respDiagnostics
	}

	// load esxi kernel modules
	tflog.Info(ctx, "Loading vmk modules")
	op, err = sshP.Run("vmkload_mod -l")
	if err != nil {
		respDiagnostics.AddError(
			"Error loading vmk modules",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Finished loading vmk modules")
	tflog.Debug(ctx, op)

	// upload driver config
	// recreate scpProvisioner
	scpProv = client.NewScpProvisioner(sshP)
	tflog.Info(ctx, "Uploading driver config")
	drvCfgTarget := filepath.Join(dir, "drv_cfg")
	err = scpProv.Upload(esxiInput.DrvCfg.ValueString(), drvCfgTarget, "0755")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
		return respDiagnostics
	}
	// query mdms via drv cfg
	tflog.Info(ctx, "Querying mdm ips via drv cfg")
	op, err = sshP.Run(drvCfgTarget + " --query_mdm")
	if err != nil {
		respDiagnostics.AddError(
			"Error querying mdm ips via drv cfg",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}

	return respDiagnostics
}

// DeleteEsxi - function to uninstall SDC package in ESXi host
func (r *SdcHostResource) DeleteEsxi(ctx context.Context, state models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")
	sshP, _, err := r.getSshProvisioner(ctx, state)
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
