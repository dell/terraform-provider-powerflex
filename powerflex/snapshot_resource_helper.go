package powerflex

import (
	"strconv"

	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// SdcList struct for sdc info response mapping to terrafrom
type SdcList struct {
	SdcID                 string `tfsdk:"sdc_id"`
	LimitIops             int    `tfsdk:"limit_iops"`
	LimitBwInMbps         int    `tfsdk:"limit_bw_in_mbps"`
	SdcName               string `tfsdk:"sdc_name"`
	AccessMode            string `tfsdk:"access_mode"`
}

// SnapshotTerraformState function to convert goscaleio snapshot struct to terraform snapshot struct
func SnapshotTerraformState(vol *pftypes.Volume, plan SnapshotResourceModel, sdcListState []SdcList) (state SnapshotResourceModel) {

	state.AccessMode = plan.AccessMode
	state.Size = plan.Size
	state.RemoveMode = plan.RemoveMode
	state.Name = types.StringValue(vol.Name)
	state.ID = types.StringValue(vol.ID)
	state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	state.LockAutoSnapshot = types.BoolValue(vol.LockedAutoSnapshot)
	state.CapacityUnit = plan.CapacityUnit
	state.RetentionUnit = plan.RetentionUnit
	if !plan.DesiredRetention.IsNull() {
		state.DesiredRetention = plan.DesiredRetention
	}
	if plan.VolumeID.ValueString() != "" {
		state.VolumeID = plan.VolumeID
	}
	if plan.VolumeName.ValueString() != "" {
		state.VolumeName = plan.VolumeName
	}
	if plan.Size.IsUnknown() {
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(int64(vol.SizeInKb), 10))
	} else {
		VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
	}
	// state.MapSdcIds = plan.MapSdcIds
	state.SdcList = sdcMapState(vol.MappedSdcInfo, sdcListState)
	return state
}

func sdcMapState(sdcInfos []*pftypes.MappedSdcInfo, sdcListState []SdcList) basetypes.ListValue {
	sdcInfoAttrTypes := map[string]attr.Type{
		"sdc_id":                   types.StringType,
		"limit_iops":               types.Int64Type,
		"limit_bw_in_mbps":         types.Int64Type,
		"sdc_name":                 types.StringType,
		"access_mode":              types.StringType,
	}
	sdcInfoElemType := types.ObjectType{
		AttrTypes: sdcInfoAttrTypes,
	}
	objectSdcInfos := []attr.Value{}
	for _, sls := range sdcListState {
		for _, msi := range sdcInfos {
			if sls.SdcID == msi.SdcID {
				obj := map[string]attr.Value{
					"sdc_id":                   types.StringValue(msi.SdcID),
					"limit_iops":               types.Int64Value(int64(msi.LimitIops)),
					"limit_bw_in_mbps":         types.Int64Value(int64(msi.LimitBwInMbps)),
					"sdc_name":                 types.StringValue(msi.SdcName),
					"access_mode":              types.StringValue(msi.AccessMode),
				}
				objVal, _ := types.ObjectValue(sdcInfoAttrTypes, obj)
				objectSdcInfos = append(objectSdcInfos, objVal)
			}
		}
	}
	mappedSdcInfoVal, _ := types.ListValue(sdcInfoElemType, objectSdcInfos)
	return mappedSdcInfoVal
}

func convertToMin(desireRetention int64, retentionUnit string) string {
	retentionMin := ""
	if retentionUnit == "days" {
		retentionMin = strconv.FormatInt(desireRetention*24*60, 10)
	} else {
		retentionMin = strconv.FormatInt(desireRetention*60, 10)
	}
	return retentionMin
}
