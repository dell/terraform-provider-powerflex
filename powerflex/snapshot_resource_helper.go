package powerflex

import (
	"strconv"

	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnapshotTerraformState function to convert goscaleio snapshot struct to terraform snapshot struct
func SnapshotTerraformState(vol *pftypes.Volume, plan SnapshotResourceModel) (state SnapshotResourceModel) {
	state.Name = types.StringValue(vol.Name)
	state.VolumeID = plan.VolumeID
	state.AccessMode = plan.AccessMode
	state.ID = types.StringValue(vol.ID)
	state.Size = plan.Size
	state.CapacityUnit = plan.CapacityUnit
	if plan.Size.IsUnknown() {
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(int64(vol.SizeInKb), 10))
	} else {
		VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
	}
	state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	state.LockedAutoSnapshot = types.BoolValue(vol.LockedAutoSnapshot)
	state.MapSdcIds = plan.MapSdcIds
	return state
}
