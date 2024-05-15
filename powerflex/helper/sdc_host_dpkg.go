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
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UbuntuSdcPackage struct
type UbuntuSdcPackage struct {
	Siob        string
	Sig         string
	SiobExtract string
}

// GetUbuntuSdcPackage returns the siob, sig and siob_extract files
func GetUbuntuSdcPackage(files []string) (*UbuntuSdcPackage, error) {
	if len(files) != 3 {
		return nil, fmt.Errorf("invalid number of files: %d, expecting 3 files", len(files))
	}
	var pkg UbuntuSdcPackage
	// one of the files must end in .siob
	for _, file := range files {
		if strings.HasSuffix(file, ".siob") {
			pkg.Siob = file
		}
	}
	// one of the files must end in .sig
	for _, file := range files {
		if strings.HasSuffix(file, ".sig") {
			pkg.Sig = file
		}
	}
	// one of the filename must be siob_extract
	for _, file := range files {
		if file == "siob_extract" {
			pkg.SiobExtract = file
		}
	}

	var err error
	if pkg.Siob == "" {
		err = errors.Join(err, fmt.Errorf("no .siob file found"))
	}
	if pkg.Sig == "" {
		err = errors.Join(err, fmt.Errorf("no .sig file found"))
	}
	if pkg.SiobExtract == "" {
		err = errors.Join(err, fmt.Errorf("no siob_extract file found"))
	}
	return &pkg, err

}

// CreateUbuntu creates an linux ubuntu SDC host
func (r *SdcHostResource) CreateUbuntu(ctx context.Context, plan models.SdcHostModel, sshP *client.SSHProvisioner, dir string) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics

	// upload sw
	scpProv := client.NewScpProvisioner(sshP)
	pkgTarget := filepath.Join(dir, "emc-sdc-package.tar")
	err := scpProv.Upload(plan.Pkg.ValueString(), pkgTarget, "")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
		return respDiagnostics
	}

	// extract software
	files, err := sshP.UntarUnix("emc-sdc-package.tar", dir)
	if err != nil {
		respDiagnostics.AddError(
			"Error extracting package",
			err.Error(),
		)
		return respDiagnostics
	}
	// verify that there are 3 files only - a siob file, a sig file and a file called siob_extract
	pkg, err := GetUbuntuSdcPackage(files)
	if err != nil {
		respDiagnostics.AddError(
			"Error extracting package",
			err.Error(),
		)
		return respDiagnostics
	}
	// run siob extract
	op, err := sshP.RunWithDir(dir, fmt.Sprintf("./%s %s", pkg.SiobExtract, pkg.Siob))
	if err != nil {
		respDiagnostics.AddError(
			"Error extracting siob file",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, op)

	mdmIPs, dgs := r.GetMdmIps(ctx, plan)

	if dgs.HasError() {
		respDiagnostics = append(respDiagnostics, dgs...)
		return respDiagnostics
	}

	// install sw
	// the software name is same as siob file, but with .deb extension instead of .siob
	debName := strings.ReplaceAll(pkg.Siob, ".siob", ".deb")
	op, err = sshP.RunWithDir(dir, fmt.Sprintf("MDM_IP=%s dpkg -i %s", strings.Join(mdmIPs, ","), debName))
	if err != nil {
		respDiagnostics.AddError(
			"Error installing sdc package",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, op)

	// check that scini status has the log SUCCESS
	op, err = sshP.Run("systemctl status scini")
	if err != nil {
		respDiagnostics.AddError(
			"Error checking scini status after restart",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	if !strings.Contains(op, "SUCCESS") {
		respDiagnostics.AddError(
			"scini service did not restart successfully",
			op,
		)
		return respDiagnostics
	}

	return respDiagnostics
}

// DeleteUbuntu - function to uninstall SDC package in Linux Ubuntu host
func (r *SdcHostResource) DeleteUbuntu(ctx context.Context, state models.SdcHostModel, sshP *client.SSHProvisioner) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")

	// list dpkg packages
	op, err := sshP.Run("dpkg -l")
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
	op, err = sshP.Run(fmt.Sprintf("dpkg -P %s", sdcPkg))
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
