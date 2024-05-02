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
	err = scpProv.Upload(plan.DrvCfg.ValueString(), drvCfgTarget, "0755")
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

// // CreateWindows creates an windows SDC host
// func (r *SdcHostResource) CreateWindows(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
// 	var respDiagnostics diag.Diagnostics

// 	var remote models.SdcHostRemoteModel
// 	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

// 	winRMClient := &winrmclient.WinRMClient{}

// 	var contexts map[string]string

// 	_ = json.Unmarshal(byteData, &contexts)

// 	context := make(map[string]string)

// 	context["username"] = remote.User

// 	context["password"] = *remote.Password

// 	context["host"] = plan.Host.ValueString()

// 	winRMClient.GetConnection(context, false)

// 	// defer winRMClient.Close() TODO

// 	mdmIPs, err := r.GetMdmIps(ctx, plan)

// 	if err != nil {
// 		respDiagnostics.AddError(
// 			"Error getting MDM IPs",
// 			err.Error(),
// 		)
// 		return respDiagnostics
// 	}

// 	if winRMClient.Init() {
// 		winRMClient.Upload("C:\\EMC-ScaleIO-sdc.msi", plan.Pkg.ValueString())

// 		ouptut := winRMClient.ExecuteCommand("msiexec.exe /i \"C:\\EMC-ScaleIO-sdc.msi\" MDM_IP=\"" + strings.Join(mdmIPs, ",") + "\" /q")

// 		if ouptut == "SUCCESS" {

// 			time.Sleep(30 * time.Second)

// 			tflog.Info(ctx, "Installed SDC Package")

// 		} else {

// 			respDiagnostics.AddError(
// 				"Error while installing command",
// 				winRMClient.Errors[0]["message"],
// 			)
// 			return respDiagnostics
// 		}
// 	}

// 	return respDiagnostics
// }

// // DeleteWindows - function to uninstall SDC package in Windows host
// func (r *SdcHostResource) DeleteWindows(ctx context.Context, state models.SdcHostModel) diag.Diagnostics {
// 	var respDiagnostics diag.Diagnostics
// 	// Disconnect from PowerFlex
// 	tflog.Info(ctx, "Logging into host...")

// 	var remote models.SdcHostRemoteModel
// 	state.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

// 	winRMClient := &winrmclient.WinRMClient{}

// 	var contexts map[string]string

// 	_ = json.Unmarshal(byteData, &contexts)

// 	context := make(map[string]string)

// 	context["username"] = remote.User

// 	context["password"] = *remote.Password

// 	context["host"] = state.Host.ValueString()

// 	winRMClient.GetConnection(context, false)

// 	// defer winRMClient.Close() TODO

// 	if winRMClient.Init() {
// 		ouptut := winRMClient.ExecuteCommand("msiexec.exe /x \"C:\\EMC-ScaleIO-sdc.msi\" /q")

// 		if ouptut == "SUCCESS" {

// 			time.Sleep(10 * time.Second)

// 			tflog.Info(ctx, "Uninstalled SDC Package")

// 		} else {

// 			respDiagnostics.AddError(
// 				"Error while installing command",
// 				winRMClient.Errors[0]["message"],
// 			)
// 			return respDiagnostics
// 		}
// 	}

// 	return respDiagnostics
// }
