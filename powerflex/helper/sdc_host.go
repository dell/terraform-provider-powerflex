package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ client.Logger = &provisionerLogger{}
)

type provisionerLogger struct {
	ctx context.Context
}

func (l *provisionerLogger) Printf(format string, v ...any) {
	tflog.Info(l.ctx, fmt.Sprintf(format, v...))
}

func (l *provisionerLogger) Println(v ...any) {
	tflog.Info(l.ctx, fmt.Sprint(v...))
}

// SdcHostResource - helper for SDC host resource
type SdcHostResource struct {
	System *goscaleio.System
}

func (r *SdcHostResource) getSshProvisioner(ctx context.Context, plan models.SdcHostModel) (*client.SshProvisioner, string, error) {
	var remote models.SdcHostRemoteModel
	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	dir := ""
	if remote.Dir == nil {
		dir = "/tmp"
	} else {
		dir = *remote.Dir
	}
	prov, err := client.NewSshProvisioner(client.SshProvisionerConfig{
		IP:         plan.Host.ValueString(),
		Username:   remote.User,
		Password:   remote.Password,
		PrivateKey: remote.PrivateKey,
	}, &provisionerLogger{ctx: ctx})
	return prov, dir, err
}

func (r *SdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) ([]string, error) {
	var mdmIps []string
	if !plan.MdmIPs.IsNull() && len(plan.MdmIPs.Elements()) > 0 {
		diags := plan.MdmIPs.ElementsAs(ctx, &mdmIps, true)
		if diags.HasError() {
			return nil, fmt.Errorf("Error in getting MDM Details", diags.Errors())
		}
	} else {
		mdmDetails, err := r.System.GetMDMClusterDetails()

		mdmIps = GetMdmIPList(mdmDetails)
		if err != nil {
			return nil, fmt.Errorf("Error in getting MDM Details on the PowerFlex cluster", err.Error())

		}
	}

	return mdmIps, nil
}

func GetMdmIPList(mdmDetails *goscaleio_types.MdmCluster) []string {
	var ipmap []string

	for index := range mdmDetails.PrimaryMDM.IPs {
		ipmap = append(ipmap, mdmDetails.PrimaryMDM.IPs[index])
	}

	for _, mdm := range mdmDetails.SecondaryMDM {
		for index := range mdm.IPs {
			ipmap = append(ipmap, mdm.IPs[index])
		}
	}

	return ipmap
}

func (r *SdcHostResource) ReadSDCHost(ctx context.Context, state models.SdcHostModel) (models.SdcHostModel, error) {
	// get SDC by IP
	tflog.Info(ctx, "Finding SDC by IP")
	sdcData, err := r.System.FindSdc("SdcIP", state.Host.ValueString())
	if err != nil {
		return state, fmt.Errorf("error finding SDC by IP %s: %w", state.Host.ValueString(), err)
	}
	tflog.Info(ctx, "Found SDC by IP")
	state.ID = types.StringValue(sdcData.Sdc.ID)
	state.PerformanceProfile = types.StringValue(sdcData.Sdc.PerfProfile)
	state.Name = types.StringValue(sdcData.Sdc.Name)
	state.OS = types.StringValue(sdcData.Sdc.OSType)
	state.MdmConnectionState = types.StringValue(sdcData.Sdc.MdmConnectionState)
	state.IsApproved = types.BoolValue(sdcData.Sdc.SdcApproved)
	state.SystemID = types.StringValue(sdcData.Sdc.SystemID)
	state.OnVMWare = types.BoolValue(sdcData.Sdc.OnVMWare)
	state.GUID = types.StringValue(sdcData.Sdc.SdcGUID)
	return state, nil
}

// SetSDCParams - function to set SDC parameters
func (r *SdcHostResource) SetSDCParams(ctx context.Context, plan, state models.SdcHostModel) error {
	// set name
	if plan.Name.ValueString() != state.Name.ValueString() && !plan.Name.IsUnknown() {
		tflog.Info(ctx, "Setting SDC name")
		if _, err := r.System.ChangeSdcName(state.ID.ValueString(), plan.Name.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC name: %w", err)
		}
		tflog.Info(ctx, "SDC name set")
	}

	// set Performance Profile
	if plan.PerformanceProfile.ValueString() != state.PerformanceProfile.ValueString() && !plan.PerformanceProfile.IsUnknown() {
		tflog.Info(ctx, "Setting SDC performance profile")
		if _, err := r.System.ChangeSdcPerfProfile(state.ID.ValueString(), plan.PerformanceProfile.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC performance profile: %w", err)
		}
		tflog.Info(ctx, "SDC performance profile set")
	}

	return nil
}

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

	// defer winRMClient.Close() TODO

	mdmIPs, err := r.GetMdmIps(ctx, plan)

	if err != nil {
		respDiagnostics.AddError(
			"Error getting MDM IPs",
			err.Error(),
		)
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

			} else {

				respDiagnostics.AddError(
					"Error while installing command",
					winRMClient.Errors[0]["message"],
				)
				return respDiagnostics
			}
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

	// defer winRMClient.Close() TODO

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
