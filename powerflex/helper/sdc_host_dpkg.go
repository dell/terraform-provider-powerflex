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

const (
	sciniv = "5.15.0-105-generic"
)

type UbuntuSdcPackage struct {
	Siob        string
	Sig         string
	SiobExtract string
}

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

func verifySciniPkg(files []string) error {
	if l := len(files); l != 2 {
		return fmt.Errorf("invalid number of files: %d, expecting 2 files", l)
	}
	found := false
	for _, file := range files {
		if file == "scini.siob" {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("scini.siob file not found")
	}
	found = false
	for _, file := range files {
		if file == "scini.siob.sig" {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("scini.siob.sig file not found")
	}
	return nil
}

// CreateEsxi creates an esxi SDC host
func (r *SdcHostResource) CreateUbuntu(ctx context.Context, plan models.SdcHostModel, sshP *client.SshProvisioner, dir string) diag.Diagnostics {
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

	// upload kernel specific scini file
	sciniTarget := filepath.Join(dir, "scini.tar")
	err = scpProv.Upload(plan.DrvCfg.ValueString(), sciniTarget, "")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading scini package",
			err.Error(),
		)
		return respDiagnostics
	}
	// extract scini
	scini_files, err := sshP.UntarUnix("scini.tar", dir)
	if err != nil {
		respDiagnostics.AddError(
			"Error extracting package",
			err.Error(),
		)
		return respDiagnostics
	}
	err = verifySciniPkg(scini_files)
	if err != nil {
		respDiagnostics.AddError(
			"Extraction of scini tar file did not yield expected files",
			err.Error(),
		)
		return respDiagnostics
	}
	// run siob extract
	op, err = sshP.RunWithDir(dir, fmt.Sprintf("./%s %s", pkg.SiobExtract, "scini.siob"))
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

	// install scini
	// list directory
	pflex_versions, err := sshP.ListDirUnix("/bin/emc/scaleio/scini_sync/driver_cache/Ubuntu")
	if err != nil {
		respDiagnostics.AddError(
			"Error listing directory",
			err.Error(),
		)
		return respDiagnostics
	}
	if len(pflex_versions) != 1 {
		respDiagnostics.AddError(
			"Error finding PowerFlex version directory",
			fmt.Sprintf("Expected 1 PowerFlex version directory in scini driver cache, got %d", len(pflex_versions)),
		)
	}
	pflex_version := pflex_versions[0]
	// make kernel specific scini directory
	sciniDir := fmt.Sprintf("/bin/emc/scaleio/scini_sync/driver_cache/Ubuntu/%s/%s", pflex_version, sciniv)
	op, err = sshP.Run(fmt.Sprintf("mkdir -p %s", sciniDir))
	if err != nil {
		respDiagnostics.AddError(
			"Error creating directory",
			err.Error()+": "+op,
		)
		return respDiagnostics
	}
	// move scini.ko to scini directory
	op, err = sshP.RunWithDir(dir, fmt.Sprintf("mv %s %s", "scini.ko", sciniDir))
	if err != nil {
		respDiagnostics.AddError(
			"Error installing scini.ko",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	// restart scini service
	op, err = sshP.Run("systemctl restart scini")
	if err != nil {
		respDiagnostics.AddError(
			"Error restarting scini service",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
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

// DeleteEsxi - function to uninstall SDC package in ESXi host
func (r *SdcHostResource) DeleteUbuntu(ctx context.Context, state models.SdcHostModel, sshP *client.SshProvisioner) diag.Diagnostics {
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

	// remove scini driver sync
	tflog.Info(ctx, "Removing scini driver sync files")
	op, err = sshP.Run("rm -rf /bin/emc/scaleio/scini_sync/driver_cache/Ubuntu")
	if err != nil {
		respDiagnostics.AddError(
			"Error removing scini driver sync files",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, op)

	return respDiagnostics
}
