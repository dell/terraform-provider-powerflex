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

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (r *SdcHostResource) getSSHProvisioner(ctx context.Context, plan models.SdcHostModel) (*client.SSHProvisioner, string, error) {
	var remote models.SdcHostRemoteModel
	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	dir := ""
	if remote.Dir == nil {
		dir = "/tmp"
	} else {
		dir = *remote.Dir
	}
	prov, err := client.NewSSHProvisioner(client.SSHProvisionerConfig{
		IP:         plan.Host.ValueString(),
		Username:   remote.User,
		Password:   remote.Password,
		PrivateKey: remote.PrivateKey,
		HostKey:    remote.HostKey,
		CaCert:     remote.CaCert,
	}, &provisionerLogger{ctx: ctx})
	return prov, dir, err
}

// GetMdmIps - get mdm ips from plan or from pflex
func (r *SdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) ([]string, diag.Diagnostics) {
	var mdmIps []string
	if !plan.MdmIPs.IsNull() && len(plan.MdmIPs.Elements()) > 0 {
		diags := plan.MdmIPs.ElementsAs(ctx, &mdmIps, true)
		if diags.HasError() {
			return nil, diags
		}
	} else {
		mdmDetails, err := r.System.GetMDMClusterDetails()
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError("Error in getting MDM Details on the PowerFlex cluster", err.Error())
			return nil, diags
		}
		mdmIps = GetMdmIPList(mdmDetails)
	}

	return mdmIps, nil
}

// GetMdmIPList - get mdm ips from pflex
func GetMdmIPList(mdmDetails *goscaleio_types.MdmCluster) []string {
	ipmap := mdmDetails.PrimaryMDM.IPs
	for _, mdm := range mdmDetails.SecondaryMDM {
		ipmap = append(ipmap, mdm.IPs...)
	}
	return ipmap
}

// ReadSDCHost - read SDC host and set state
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
	os := strings.ToLower(sdcData.Sdc.OSType)
	if strings.HasPrefix(os, "esx") {
		// both esxi is return as esx from API
		os = "esxi"
	}
	state.OS = types.StringValue(os)
	state.MdmConnectionState = types.StringValue(sdcData.Sdc.MdmConnectionState)
	state.IsApproved = types.BoolValue(sdcData.Sdc.SdcApproved)
	state.SystemID = types.StringValue(sdcData.Sdc.SystemID)
	state.OnVMWare = types.BoolValue(sdcData.Sdc.OnVMWare)
	state.GUID = types.StringValue(sdcData.Sdc.SdcGUID)

	mdmIPs, diags := r.GetMdmIps(ctx, state)
	if diags.HasError() {
		var errStr string
		for _, diag := range diags {
			errStr += diag.Summary() + "\n"
		}
		return state, fmt.Errorf(errStr)
	}

	objectMDMs := make([]attr.Value, len(mdmIPs))
	for i, link := range mdmIPs {
		obj := types.StringValue(link)
		objectMDMs[i] = obj
	}

	listVal, _ := types.ListValue(types.StringType, objectMDMs)
	state.MdmIPs = listVal
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

// LinuxOp creates or deletes a linux SDC host
func (r *SdcHostResource) LinuxOp(ctx context.Context, plan models.SdcHostModel, add bool) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	sshP, dir, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	op, err := sshP.Run("cat /etc/os-release")
	if err != nil {
		respDiagnostics.AddError(
			"Error checking /etc/os-release after restart",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}

	// Parse the output of "cat /etc/os-release" to determine the Linux distribution type
	linuxType := "unknown"
	lines := client.GetLinesUnix(op)
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			linuxType = strings.Trim(id, `"`) // Remove leading and trailing quotes
			break
		}
	}

	// You may need to customize this mapping based on your specific requirements
	linuxTypeMap := map[string]string{
		"ubuntu": "ubuntu",
		"sles":   "sles",
		"centos": "centos",
		"rhel":   "rhel",
		// Add more mappings as needed
	}

	// Check if the Linux distribution ID is mapped to a standard name
	standardLinuxType, ok := linuxTypeMap[linuxType]
	if ok {
		linuxType = standardLinuxType
	}

	tflog.Info(ctx, fmt.Sprintf("Linux distribution detected: %s", linuxType))

	switch linuxType {
	case "rhel", "sles", "centos":
		respDiagnostics.Append(r.CreateRhel(ctx, plan, sshP, dir)...)
	case "ubuntu":
		if add {
			respDiagnostics.Append(r.CreateUbuntu(ctx, plan, sshP, dir)...)
		} else {
			respDiagnostics.Append(r.DeleteUbuntu(ctx, plan, sshP)...)
		}
	default:
		respDiagnostics.AddError(
			"Could not able to find supported linux distribution",
			linuxType,
		)
	}

	return respDiagnostics
}

// DeleteLinux - delete linux SDC packages
func (r *SdcHostResource) DeleteLinux(ctx context.Context, plan models.SdcHostModel, add bool) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	sshP, _, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	op, err := sshP.Run("cat /etc/os-release")
	if err != nil {
		respDiagnostics.AddError(
			"Error checking /etc/os-release after restart",
			op+"\n"+err.Error(),
		)
		return respDiagnostics
	}

	// Parse the output of "cat /etc/os-release" to determine the Linux distribution type
	linuxType := "unknown"
	lines := strings.Split(op, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			linuxType = strings.Trim(id, `"`) // Remove leading and trailing quotes
			break
		}
	}

	// You may need to customize this mapping based on your specific requirements
	linuxTypeMap := map[string]string{
		"ubuntu": "ubuntu",
		"sles":   "sles",
		"centos": "centos",
		"rhel":   "rhel",
		// Add more mappings as needed
	}

	// Check if the Linux distribution ID is mapped to a standard name
	standardLinuxType, ok := linuxTypeMap[linuxType]
	if ok {
		linuxType = standardLinuxType
	}

	tflog.Info(ctx, fmt.Sprintf("Linux distribution detected: %s", linuxType))

	switch linuxType {
	case "rhel", "sles", "centos":
		respDiagnostics.Append(r.DeleteRhel(ctx, plan, sshP)...)
	case "ubuntu":
		respDiagnostics.Append(r.DeleteUbuntu(ctx, plan, sshP)...)

	}

	return respDiagnostics
}
