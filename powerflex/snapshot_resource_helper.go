/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package powerflex

import (
	"strconv"

	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// READWRITE represents access mode limit of snapshot
	READWRITE = "ReadWrite"
	// READONLY represents access mode limit of snapshot
	READONLY = "ReadOnly"
	// SecondsThreshold represents platform epoch drift
	SecondsThreshold = 300
	// DayInMins represents day in min
	DayInMins = 24 * HourInMins
	// HourInMins represents hour in min
	HourInMins = 60
	// MinuteInSeconds represents min in sec.
	MinuteInSeconds = 60
)

// SnapshotResourceModel maps the resource schema data.
type SnapshotResourceModel struct {
	Name             types.String `tfsdk:"name"`
	VolumeID         types.String `tfsdk:"volume_id"`
	VolumeName       types.String `tfsdk:"volume_name"`
	AccessMode       types.String `tfsdk:"access_mode"`
	ID               types.String `tfsdk:"id"`
	Size             types.Int64  `tfsdk:"size"`
	CapacityUnit     types.String `tfsdk:"capacity_unit"`
	SizeInKb         types.Int64  `tfsdk:"size_in_kb"`
	LockAutoSnapshot types.Bool   `tfsdk:"lock_auto_snapshot"`
	RemoveMode       types.String `tfsdk:"remove_mode"`
	DesiredRetention types.Int64  `tfsdk:"desired_retention"`
	RetentionUnit    types.String `tfsdk:"retention_unit"`
	RetentionInMin   types.String `tfsdk:"retention_in_min"`
}

// SdcList struct for sdc info response mapping to terrafrom
type SdcList struct {
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName       types.String `tfsdk:"sdc_name"`
	AccessMode    types.String `tfsdk:"access_mode"`
}

func refreshState(snap *pftypes.Volume, prestate *SnapshotResourceModel) (diags diag.Diagnostics) {
	var drift int64
	prestate.ID = types.StringValue(snap.ID)
	prestate.Name = types.StringValue(snap.Name)
	prestate.AccessMode = types.StringValue(snap.AccessModeLimit)
	prestate.SizeInKb = types.Int64Value(int64(snap.SizeInKb))
	switch prestate.CapacityUnit.ValueString() {
	case "TB":
		prestate.Size = types.Int64Value(int64(snap.SizeInKb / TiKB))
	case "GB":
		prestate.Size = types.Int64Value(int64(snap.SizeInKb / GiKB))
	}
	prestate.LockAutoSnapshot = types.BoolValue(snap.LockedAutoSnapshot)
	diff1 := int64(snap.SecureSnapshotExpTime) - int64(snap.CreationTime)
	if prestate.RetentionUnit.ValueString() == "days" {
		drift = diff1 - prestate.DesiredRetention.ValueInt64()*DayInMins*MinuteInSeconds
	} else {
		drift = diff1 - prestate.DesiredRetention.ValueInt64()*HourInMins*MinuteInSeconds
	}
	if diff1 > 0 && drift > SecondsThreshold && drift < -SecondsThreshold {
		prestate.RetentionInMin = types.StringValue(strconv.FormatInt(diff1/60, 10))
	}

	return diags
}

func convertToMin(desireRetention int64, retentionUnit string) string {
	retentionMin := ""
	if retentionUnit == "days" {
		retentionMin = strconv.FormatInt(desireRetention*DayInMins, 10)
	} else {
		retentionMin = strconv.FormatInt(desireRetention*HourInMins, 10)
	}
	return retentionMin
}

// coverterKB fucntion to convert size into kb
func converterKB(capacityUnit string, size int64) int64 {
	var valInKiB int64
	switch capacityUnit {
	case "TB":
		valInKiB = size * TiKB
	case "GB":
		valInKiB = size * GiKB
	default:
		return 0
	}
	return int64(valInKiB)
}
