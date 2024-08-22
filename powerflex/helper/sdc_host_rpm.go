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
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateRhel creates a RHEL SDC host
func (r *SdcHostResource) CreateRhel(ctx context.Context, plan models.SdcHostModel, sshP *client.SSHProvisioner, dir string) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics

	if !plan.UseRemotePath.ValueBool() {
		// upload sw
		scpProv := client.NewScpProvisioner(sshP)
		pkgTarget := strings.TrimSuffix(dir, "/") + "/" + "emc-sdc-package.rpm"
		err := scpProv.Upload(plan.Pkg.ValueString(), pkgTarget, "")
		if err != nil {
			respDiagnostics.AddError(
				"Error uploading package",
				err.Error(),
			)
			return respDiagnostics
		}
	}

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	// install sw
	debName := "emc-sdc-package.rpm"
	op, err := sshP.RunWithDir(dir, fmt.Sprintf("MDM_IP=%s rpm -i %s", strings.Join(mdmIPs, ","), debName))
	if err != nil {
		respDiagnostics.AddError(
			"Error installing sdc package",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, op)

	// wait 30 seconds for scini to configure itself
	tflog.Info(ctx, "Waiting 30 seconds for scini to configure itself")
	time.Sleep(30 * time.Second)
	// check that scini status has the log SUCCESS
	op, err = sshP.Run("systemctl status scini")
	if err != nil {
		respDiagnostics.AddError(
			"Error checking scini status after installing sdc package",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	if !strings.Contains(op, "SUCCESS") {
		respDiagnostics.AddError(
			"scini service did not start successfully",
			op,
		)
		return respDiagnostics
	}

	return respDiagnostics
}

// DeleteRhel - function to uninstall SDC package in RHEL host
func (r *SdcHostResource) DeleteRhel(ctx context.Context, state models.SdcHostModel, sshP *client.SSHProvisioner) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")

	// list dpkg packages
	op, err := sshP.Run("rpm -qa | grep EMC")
	if err != nil {
		respDiagnostics.AddError(
			"Error listing installed packages",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	pkgList := client.GetLinesUnix(op)
	// get the package with sdc in the name
	var sdcPkg string
	for _, pkg := range pkgList {
		if strings.Contains(pkg, "sdc") {
			sdcPkg = pkg
			break
		}
	}
	if sdcPkg == "" {
		tflog.Info(ctx, "No sdc package installed")
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("Found sdc package %s", sdcPkg))

	// remove sdc package
	tflog.Info(ctx, "Removing installed sdc package")
	op, err = sshP.Run(fmt.Sprintf("rpm -e %s", sdcPkg))
	if err != nil {
		respDiagnostics.AddError(
			"Error uninstalling package",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, op)

	return respDiagnostics
}
