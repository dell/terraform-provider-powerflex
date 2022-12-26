package powerflex

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// MiKB to convert size in megabytes
	MiKB = 1024
	// GiKB to convert size in gigabytes
	GiKB = 1024 * MiKB
	// TiKB to convert size in terabytes
	TiKB = 1024 * GiKB
)

// covertToKB fucntion to convert size into kb
func convertToKB(capacityUnit string, size int64) (int64, error) {
	var valInKiB int64
	switch capacityUnit {
	case "MB":
		valInKiB = size * MiKB
	case "TB":
		valInKiB = size * TiKB
	case "GB":
		valInKiB = size * GiKB
	default:
		return 0, errors.New("invalid capacity unit")
	}
	return int64(valInKiB), nil
}

// VolumeTerraformState function to convert goscaleio volume struct to terraform volume struct
func VolumeTerraformState(vol *pftypes.Volume, plan VolumeResourceModel) (state VolumeResourceModel) {
	state.ProtectionDomainID = plan.ProtectionDomainID
	state.Size = plan.Size
	state.CapacityUnit = plan.CapacityUnit
	state.MapSdcsID = plan.MapSdcsID
	state.LockedAutoSnapshot = plan.LockedAutoSnapshot
	VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	state.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
	state.StoragePoolID = types.StringValue(vol.StoragePoolID)
	state.UseRmCache = types.BoolValue(vol.UseRmCache)
	state.MappingToAllSdcsEnabled = types.BoolValue(vol.MappingToAllSdcsEnabled)
	state.IsObfuscated = types.BoolValue(vol.IsObfuscated)
	state.VolumeType = types.StringValue(vol.VolumeType)
	state.ConsistencyGroupID = types.StringValue(vol.ConsistencyGroupID)
	state.VTreeID = types.StringValue(vol.VTreeID)
	state.AncestorVolumeID = types.StringValue(vol.AncestorVolumeID)
	state.MappedScsiInitiatorInfo = types.StringValue(vol.MappedScsiInitiatorInfo)
	state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	state.CreationTime = types.Int64Value(int64(vol.CreationTime))
	state.Name = types.StringValue(vol.Name)
	state.ID = types.StringValue(vol.ID)
	state.DataLayout = types.StringValue(vol.DataLayout)
	state.NotGenuineSnapshot = types.BoolValue(vol.NotGenuineSnapshot)
	state.AccessModeLimit = types.StringValue(vol.AccessModeLimit)
	state.SecureSnapshotExpTime = types.Int64Value(int64(vol.SecureSnapshotExpTime))
	state.ManagedBy = types.StringValue(vol.ManagedBy)
	state.LockedAutoSnapshotMarkedForRemoval = types.BoolValue(vol.LockedAutoSnapshotMarkedForRemoval)
	state.CompressionMethod = types.StringValue(vol.CompressionMethod)
	state.TimeStampIsAccurate = types.BoolValue(vol.TimeStampIsAccurate)
	state.OriginalExpiryTime = types.Int64Value(int64(vol.OriginalExpiryTime))
	state.VolumeReplicationState = types.StringValue(vol.VolumeReplicationState)
	state.ReplicationJournalVolume = types.BoolValue(vol.ReplicationJournalVolume)
	state.ReplicationTimeStamp = types.Int64Value(int64(vol.ReplicationTimeStamp))

	linkAttrTypes := map[string]attr.Type{
		"rel":  types.StringType,
		"href": types.StringType,
	}
	mappedSdcInfoAttrTypes := map[string]attr.Type{
		"sdc_id":                   types.StringType,
		"sdc_ip":                   types.StringType,
		"limit_iops":               types.Int64Type,
		"limit_bw_in_mbps":         types.Int64Type,
		"sdc_name":                 types.StringType,
		"access_mode":              types.StringType,
		"is_direct_buffer_mapping": types.BoolType,
	}
	linkElemType := types.ObjectType{
		AttrTypes: linkAttrTypes,
	}
	mappedSdcInfoElemType := types.ObjectType{
		AttrTypes: mappedSdcInfoAttrTypes,
	}
	objectLinks := []attr.Value{}
	objectMappedSdcInfos := []attr.Value{}

	for _, link := range vol.Links {
		obj := map[string]attr.Value{
			"rel":  types.StringValue(link.Rel),
			"href": types.StringValue(link.HREF),
		}
		objVal, _ := types.ObjectValue(linkAttrTypes, obj)
		objectLinks = append(objectLinks, objVal)
	}
	listVal, _ := types.ListValue(linkElemType, objectLinks)

	for _, msi := range vol.MappedSdcInfo {
		obj := map[string]attr.Value{
			"sdc_id":                   types.StringValue(msi.SdcID),
			"sdc_ip":                   types.StringValue(msi.SdcIP),
			"limit_iops":               types.Int64Value(int64(msi.LimitIops)),
			"limit_bw_in_mbps":         types.Int64Value(int64(msi.LimitBwInMbps)),
			"sdc_name":                 types.StringValue(msi.SdcName),
			"access_mode":              types.StringValue(msi.AccessMode),
			"is_direct_buffer_mapping": types.BoolValue(msi.IsDirectBufferMapping),
		}
		objVal, _ := types.ObjectValue(mappedSdcInfoAttrTypes, obj)
		objectMappedSdcInfos = append(objectMappedSdcInfos, objVal)
	}
	mappedSdcInfoVal, _ := types.ListValue(mappedSdcInfoElemType, objectMappedSdcInfos)
	state.Links = listVal
	state.MappedSdcInfo = mappedSdcInfoVal
	return state
}

// getStoragePoolInstance function to get storage pool from storage pool id and protection domain id
func getStoragePoolInstance(c *goscaleio.Client, spID string, pdID string) (*goscaleio.StoragePool, error) {
	getSystems, _ := c.GetSystems()
	sr := goscaleio.NewSystem(c)
	sr.System = getSystems[0]
	pdr := goscaleio.NewProtectionDomain(c)
	protectionDomain, err := sr.FindProtectionDomain(pdID, "", "")
	if err != nil {
		return nil, err
	}
	pdr.ProtectionDomain = protectionDomain
	spr := goscaleio.NewStoragePool(c)
	storagePool, err := pdr.FindStoragePool(spID, "", "")
	spr.StoragePool = storagePool
	if err != nil {
		return nil, err
	}
	return spr, nil
}

// Difference function to find the state difference b/w sdcs
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// PlanModifyString runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown or known, do not set default value.
	// if !req.PlanValue.IsNull() {
	// 	return
	// }

	resp.PlanValue = types.StringValue(m.Default)
}

func stringDefault(defaultValue string) planmodifier.String {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}
