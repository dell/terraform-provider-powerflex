package powerflex

import (
	"strconv"

	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type attrState struct {
	name, volumeID, volumeName, desiredRetention, retentionUnit,
	accessMode, capacityUnit, size, sizeInKb, lockAutoSnapshot bool
	sdcList map[string]*sdcAttr
	errMsg  map[string]string
}

type sdcAttr struct {
	isSdcMapped, isSdcLimitSet, isAccessModeSet bool
}

// SdcInfoAttrTypes for defining sdc list stuct into terraform type
var SdcInfoAttrTypes = map[string]attr.Type{
	"sdc_id":           types.StringType,
	"limit_iops":       types.Int64Type,
	"limit_bw_in_mbps": types.Int64Type,
	"sdc_name":         types.StringType,
	"access_mode":      types.StringType,
}

// SdcList struct for sdc info response mapping to terrafrom
type SdcList struct {
	SdcID         string `tfsdk:"sdc_id"`
	LimitIops     int    `tfsdk:"limit_iops"`
	LimitBwInMbps int    `tfsdk:"limit_bw_in_mbps"`
	SdcName       string `tfsdk:"sdc_name"`
	AccessMode    string `tfsdk:"access_mode"`
}

// SnapshotTerraformState for taking latest snapshot on infrastructure and storing latest snapshot data into state file
func SnapshotTerraformState(vol *pftypes.Volume, plan SnapshotResourceModel, success *attrState) (state SnapshotResourceModel) {
	// update name and id of snapshot in state
	if success.name && vol.Name != "" && vol.ID != "" {
		state.Name = types.StringValue(vol.Name)
		state.ID = types.StringValue(vol.ID)
	}
	if success.volumeID && !plan.VolumeID.IsNull() {
		state.VolumeID = plan.VolumeID
	}
	if success.volumeName && !plan.VolumeName.IsNull() {
		state.VolumeName = plan.VolumeName
	}
	if success.accessMode && !plan.AccessMode.IsNull() {
		state.AccessMode = plan.AccessMode
	}
	if success.capacityUnit && !plan.CapacityUnit.IsNull() {
		state.CapacityUnit = plan.CapacityUnit
	}
	if plan.Size.IsUnknown() {
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(int64(vol.SizeInKb), 10))
	} else {
		VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
		state.Size = plan.Size
	}
	if success.sizeInKb {
		state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	}
	state.LockAutoSnapshot = types.BoolValue(vol.LockedAutoSnapshot)
	state.RemoveMode = plan.RemoveMode
	if success.desiredRetention && !plan.DesiredRetention.IsNull() {
		state.DesiredRetention = plan.DesiredRetention
	}
	if success.retentionUnit && !plan.RetentionUnit.IsNull() {
		state.RetentionUnit = plan.RetentionUnit
	}
	state.SdcList = sdcMapState(vol.MappedSdcInfo, success.sdcList)
	return state
}

func sdcMapState(sdcInfos []*pftypes.MappedSdcInfo, sdcListState map[string]*sdcAttr) basetypes.SetValue {
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SdcInfoAttrTypes,
	}
	objectSdcInfos := []attr.Value{}
	for key := range sdcListState {
		for _, msi := range sdcInfos {
			if key == msi.SdcID {
				obj := map[string]attr.Value{
					"sdc_id":           types.StringValue(msi.SdcID),
					"limit_iops":       types.Int64Value(int64(msi.LimitIops)),
					"limit_bw_in_mbps": types.Int64Value(int64(msi.LimitBwInMbps)),
					"sdc_name":         types.StringValue(msi.SdcName),
					"access_mode":      types.StringValue(msi.AccessMode),
				}
				objVal, _ := types.ObjectValue(SdcInfoAttrTypes, obj)
				objectSdcInfos = append(objectSdcInfos, objVal)
			}
		}
	}
	mappedSdcInfoVal, _ := types.SetValue(sdcInfoElemType, objectSdcInfos)
	return mappedSdcInfoVal
}

func refreshState(snap *pftypes.Volume, prestate *SnapshotResourceModel) {
	prestate.ID = types.StringValue(snap.ID)
	prestate.Name = types.StringValue(snap.Name)
	prestate.AccessMode = types.StringValue(snap.AccessModeLimit)
	prestate.SizeInKb = types.Int64Value(int64(snap.SizeInKb))
	prestate.VolumeSizeInKb = types.StringValue(strconv.FormatInt(int64(snap.SizeInKb), 10))
	switch prestate.CapacityUnit.ValueString() {
	case "":
		prestate.Size = types.Int64Value(int64(snap.SizeInKb / TiKB))
	case "GB":
		prestate.Size = types.Int64Value(int64(snap.SizeInKb / GiKB))
	}
	prestate.LockAutoSnapshot = types.BoolValue(snap.LockedAutoSnapshot)
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SdcInfoAttrTypes,
	}
	objectSdcInfos := []attr.Value{}
	for _, msi := range snap.MappedSdcInfo {
		// refreshing state for drift outside terraform
		obj := map[string]attr.Value{
			"sdc_id":           types.StringValue(msi.SdcID),
			"limit_iops":       types.Int64Value(int64(msi.LimitIops)),
			"limit_bw_in_mbps": types.Int64Value(int64(msi.LimitBwInMbps)),
			"sdc_name":         types.StringValue(msi.SdcName),
			"access_mode":      types.StringValue(msi.AccessMode),
		}
		objVal, _ := types.ObjectValue(SdcInfoAttrTypes, obj)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	mappedSdcInfoVal, _ := types.SetValue(sdcInfoElemType, objectSdcInfos)
	prestate.SdcList = mappedSdcInfoVal
}

func convertToMin(desireRetention int64, retentionUnit string) string {
	retentionMin := ""
	if retentionUnit == "days" {
		// retentionMin = strconv.FormatInt(desireRetention*24*60, 10)
		retentionMin = strconv.FormatInt(desireRetention*4, 10)
	} else {
		// retentionMin = strconv.FormatInt(desireRetention*60, 10)
		retentionMin = strconv.FormatInt(desireRetention*2, 10)
	}
	return retentionMin
}
