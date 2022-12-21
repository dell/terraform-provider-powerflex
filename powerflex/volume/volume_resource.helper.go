package volume

import (
	"errors"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updateVolumeState(vol *pftypes.Volume, plan volumeResourceModel) (response volumeResourceModel) {
	response = volumeResourceModel{
		ProtectionDomainID:                 plan.ProtectionDomainID,
		VolumeSizeInKb:                     plan.VolumeSizeInKb,
		StoragePoolID:                      types.StringValue(vol.StoragePoolID),
		UseRmCache:                         types.BoolValue(vol.UseRmCache),
		MappingToAllSdcsEnabled:            types.BoolValue(vol.MappingToAllSdcsEnabled),
		IsObfuscated:                       types.BoolValue(vol.IsObfuscated),
		VolumeType:                         types.StringValue(vol.VolumeType),
		ConsistencyGroupID:                 types.StringValue(vol.ConsistencyGroupID),
		VTreeID:                            types.StringValue(vol.VTreeID),
		AncestorVolumeID:                   types.StringValue(vol.AncestorVolumeID),
		MappedScsiInitiatorInfo:            types.StringValue(vol.MappedScsiInitiatorInfo),
		SizeInKb:                           types.Int64Value(int64(vol.SizeInKb)),
		CreationTime:                       types.Int64Value(int64(vol.CreationTime)),
		Name:                               types.StringValue(vol.Name),
		ID:                                 types.StringValue(vol.ID),
		DataLayout:                         types.StringValue(vol.DataLayout),
		NotGenuineSnapshot:                 types.BoolValue(vol.NotGenuineSnapshot),
		AccessModeLimit:                    types.StringValue(vol.AccessModeLimit),
		SecureSnapshotExpTime:              types.Int64Value(int64(vol.SecureSnapshotExpTime)),
		ManagedBy:                          types.StringValue(vol.ManagedBy),
		LockedAutoSnapshot:                 types.BoolValue(vol.LockedAutoSnapshot),
		LockedAutoSnapshotMarkedForRemoval: types.BoolValue(vol.LockedAutoSnapshotMarkedForRemoval),
		CompressionMethod:                  types.StringValue(vol.CompressionMethod),
		TimeStampIsAccurate:                types.BoolValue(vol.TimeStampIsAccurate),
		OriginalExpiryTime:                 types.Int64Value(int64(vol.OriginalExpiryTime)),
		VolumeReplicationState:             types.StringValue(vol.VolumeReplicationState),
		ReplicationJournalVolume:           types.BoolValue(vol.ReplicationJournalVolume),
		ReplicationTimeStamp:               types.Int64Value(int64(vol.ReplicationTimeStamp)),
	}
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
	response.Links = listVal
	response.MappedSdcInfo = mappedSdcInfoVal
	return response
}

func getStoragePoolInstance(c *goscaleio.Client, spId string, pdId string) (*goscaleio.StoragePool, error) {
	getSystems, _ := c.GetSystems()
	sr := goscaleio.NewSystem(c)
	sr.System = getSystems[0]
	getProtectionDomains, _ := sr.GetProtectionDomain("")
	pdr := goscaleio.NewProtectionDomain(c)
	for _, protectionDomain := range getProtectionDomains {
		pdr.ProtectionDomain = protectionDomain
		if pdr.ProtectionDomain.ID == pdId {
			getStoragePools, _ := pdr.GetStoragePool("")
			spr := goscaleio.NewStoragePool(c)
			for _, sp := range getStoragePools {
				spr.StoragePool = sp
				if spr.StoragePool.ID == spId {
					return spr, nil
				}
			}
		}
	}
	return nil, errors.New("couldn't find the storage pool")
}

