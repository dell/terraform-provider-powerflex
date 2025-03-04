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
		Port:       remote.Port,
		IP:         plan.Host.ValueString(),
		Username:   remote.User,
		Password:   remote.Password,
		PrivateKey: remote.PrivateKey,
		HostKey:    remote.HostKey,
		CaCert:     remote.CaCert,
	}, &provisionerLogger{ctx: ctx})
	return prov, dir, err
}

// GetMdmIps - get mdm ips from PowerFlex
func (r *SdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) ([]string, diag.Diagnostics) {
	mdmDetails, err := r.System.GetMDMClusterDetails()
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Error in getting MDM Details on the PowerFlex cluster", err.Error())
		return nil, diags
	}

	return GetMdmIPList(mdmDetails), nil
}

// GetMdmIPList - get mdm ips from PowerFlex
func GetMdmIPList(mdmDetails *goscaleio_types.MdmCluster) []string {

	var ipmap []string

	if len(mdmDetails.VirtualIPs) > 0 {
		ipmap = mdmDetails.VirtualIPs
	} else {
		ipmap = mdmDetails.PrimaryMDM.IPs
		for _, mdm := range mdmDetails.SecondaryMDM {
			ipmap = append(ipmap, mdm.IPs...)
		}
	}

	return ipmap
}

// ReadSDCHost - read SDC host and set state
func (r *SdcHostResource) ReadSDCHost(ctx context.Context, client *goscaleio.Client, state models.SdcHostModel) (models.SdcHostModel, error) {
	// get SDC if ID is set use that, if GUID is set use that, last try using IP
	tflog.Info(ctx, "Finding SDC")
	var sdcData *goscaleio.Sdc
	var err error
	// Check by ID
	if !state.ID.IsNull() && !state.ID.IsUnknown() {
		sdcData, err = r.System.GetSdcByID(state.ID.ValueString())
		if err != nil {
			return state, fmt.Errorf("error finding SDC by ID %s: %w", state.ID.ValueString(), err)
		}
		// Check by UUID
	} else if !state.GUID.IsNull() && !state.GUID.IsUnknown() {
		sdcss, err := r.System.GetSdc()

		if err != nil {
			return state, fmt.Errorf("error grabbing SDCs: %w", err)
		}
		for _, val := range sdcss {
			tflog.Debug(ctx, fmt.Sprintf("checking GUIDs IP: %s GUID: %s", val.SdcIP, val.SdcGUID))
			if val.SdcGUID == state.GUID.ValueString() {
				tflog.Debug(ctx, fmt.Sprintf("Found GUID: %s", val.SdcGUID))
				sdcData = goscaleio.NewSdc(client, &val)
				break
			}
		}
		if sdcData == nil {
			return state, fmt.Errorf("error finding SDC by GUID %s", state.GUID.ValueString())
		}
		tflog.Info(ctx, "Found SDC by GUID")
		// If all else fails try by IP
	} else {
		sdcData, err = r.System.FindSdc("SdcIP", state.Host.ValueString())
		if err != nil {
			return state, fmt.Errorf("error finding SDC by IP %s: %w", state.Host.ValueString(), err)
		}
		tflog.Info(ctx, "Found SDC by IP")
	}

	state.ID = types.StringValue(sdcData.Sdc.ID)
	state.PerformanceProfile = types.StringValue(sdcData.Sdc.PerfProfile)
	state.Name = types.StringValue(sdcData.Sdc.Name)
	state.Host = types.StringValue(sdcData.Sdc.SdcIP)
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
	if state.MdmIPs.IsUnknown() {
		// Just make an empty list if not set by user
		state.MdmIPs, _ = types.ListValue(types.StringType, []attr.Value{})
	}
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

// UpdateLinuxMdms - update linux SDC host
func (r *SdcHostResource) UpdateLinuxMdms(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Create the ssh provisioner
	sshP, _, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()
	// Check for existing mdms
	qMdms, qErr := sshP.RunWithDir(plan.LinuxDrvCfg.ValueString(), "./drv_cfg --query_mdms")
	if qErr != nil {
		respDiagnostics.AddError(
			"Error retrieving mdms with",
			qErr.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Existing MDMS: "+qMdms)
	var mdms []string
	respDiagnostics = append(respDiagnostics, plan.MdmIPs.ElementsAs(ctx, &mdms, true)...)
	if respDiagnostics.HasError() {
		return respDiagnostics
	}
	// If they already exist, attempt to update them othewise add them
	for _, mdm := range mdms {
		splitMdms := strings.Split(mdm, ",")
		if strings.Contains(qMdms, splitMdms[0]) {
			tflog.Info(ctx, "Updating MDMS: "+mdm)
			_, err := sshP.RunWithDir(plan.LinuxDrvCfg.ValueString(), fmt.Sprintf("./drv_cfg --mod_mdm_ip --ip=%s --new_mdm_ip=%s", splitMdms[0], mdm))
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
			_, err := sshP.RunWithDir(plan.LinuxDrvCfg.ValueString(), fmt.Sprintf("./drv_cfg --add_mdm --ip=%s", mdm))
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

	// Check to see if all mdms were set properly by the drv_cfg command
	for _, mdm := range mdms {
		splitMdms := strings.Split(mdm, ",")
		qMdmsAfter, qErrAfter := sshP.RunWithDir(plan.LinuxDrvCfg.ValueString(), "./drv_cfg --query_mdms | grep "+splitMdms[0])
		if qErrAfter != nil {
			respDiagnostics.AddError(
				"Error validating mdms",
				qErrAfter.Error(),
			)
			return respDiagnostics
		}
		// If it contains 0000000000000000 that means the MDM was not valid
		if !strings.Contains(qMdmsAfter, splitMdms[0]) || strings.Contains(qMdmsAfter, "0000000000000000") {
			respDiagnostics.AddError(
				"Error validating mdms",
				"MDMS "+mdm+" were invalid please check the configuration and try again",
			)
			return respDiagnostics
		}
	}

	// After setting the new config let it sleep to make sure it is set
	time.Sleep(5 * time.Second)
	return respDiagnostics
}

// LinuxOp creates or deletes a linux SDC host
func (r *SdcHostResource) LinuxOp(ctx context.Context, plan models.SdcHostModel, add bool) (models.SdcHostModel, diag.Diagnostics) {
	var respDiagnostics diag.Diagnostics
	sshP, dir, err := r.getSSHProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return plan, respDiagnostics
	}
	defer sshP.Close()

	op, err := sshP.Run("cat /etc/os-release")
	if err != nil {
		respDiagnostics.AddError(
			"Error retrieving contents of /etc/os-release",
			op+"\n"+err.Error(),
		)
		return plan, respDiagnostics
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
	if standardLinuxType, ok := linuxTypeMap[linuxType]; ok {
		linuxType = standardLinuxType
	}

	tflog.Info(ctx, fmt.Sprintf("Linux distribution detected: %s", linuxType))

	switch linuxType {
	case "rhel", "sles", "centos":
		if add {
			planCreate, daigsRhelCreate := r.CreateRhel(ctx, plan, sshP, dir)
			plan = planCreate
			respDiagnostics.Append(daigsRhelCreate...)
			// If the MdmIPs are set in the plan then use the drv_conf to update them
			if !plan.MdmIPs.IsUnknown() && len(plan.MdmIPs.Elements()) > 0 {
				respDiagnostics.Append(r.UpdateLinuxMdms(ctx, plan)...)
			}
		} else {
			respDiagnostics.Append(r.DeleteRhel(ctx, plan, sshP)...)
		}
	case "ubuntu":
		if add {
			planCreate, daigsUbuntuCreate := r.CreateUbuntu(ctx, plan, sshP, dir)
			plan = planCreate
			respDiagnostics.Append(daigsUbuntuCreate...)
			// If the MdmIPs are set in the plan then use the drv_conf to update them
			if !plan.MdmIPs.IsUnknown() && len(plan.MdmIPs.Elements()) > 0 {
				respDiagnostics.Append(r.UpdateLinuxMdms(ctx, plan)...)
			}
		} else {
			respDiagnostics.Append(r.DeleteUbuntu(ctx, plan, sshP)...)
		}
	default:
		respDiagnostics.AddError(
			"Could not find supported linux distribution",
			linuxType,
		)
	}

	return plan, respDiagnostics
}
