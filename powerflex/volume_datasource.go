package powerflex

import (
	"context"
	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &volumeDataSource{}
	_ datasource.DataSourceWithConfigure = &volumeDataSource{}
)

// VolumeDataSource returns the volume data source
func VolumeDataSource() datasource.DataSource {
	return &volumeDataSource{}
}

type volumeDataSource struct {
	client *goscaleio.Client
}

type volumeDataSourceModel struct {
	Volumes         []volumeModel `tfsdk:"volumes"`
	ID              types.String  `tfsdk:"id"`
	StoragePoolID   types.String  `tfsdk:"storage_pool_id"`
	StoragePoolName types.String  `tfsdk:"storage_pool_name"`
	Name            types.String  `tfsdk:"name"`
}

type volumeModel struct {
	ID                                 types.String         `tfsdk:"id"`
	Name                               types.String         `tfsdk:"name"`
	CreationTime                       types.Int64          `tfsdk:"creation_time"`
	SizeInKb                           types.Int64          `tfsdk:"size_in_kb"`
	AncestorVolumeID                   types.String         `tfsdk:"ancestor_volume_id"`
	VTreeID                            types.String         `tfsdk:"vtree_id"`
	ConsistencyGroupID                 types.String         `tfsdk:"consistency_group_id"`
	VolumeType                         types.String         `tfsdk:"volume_type"`
	UseRmCache                         types.Bool           `tfsdk:"use_rm_cache"`
	StoragePoolID                      types.String         `tfsdk:"storage_pool_id"`
	DataLayout                         types.String         `tfsdk:"data_layout"`
	NotGenuineSnapshot                 types.Bool           `tfsdk:"not_genuine_snapshot"`
	AccessModeLimit                    types.String         `tfsdk:"access_mode_limit"`
	SecureSnapshotExpTime              types.Int64          `tfsdk:"secure_snapshot_exp_time"`
	ManagedBy                          types.String         `tfsdk:"managed_by"`
	LockedAutoSnapshot                 types.Bool           `tfsdk:"locked_auto_snapshot"`
	LockedAutoSnapshotMarkedForRemoval types.Bool           `tfsdk:"locked_auto_snapshot_marked_for_removal"`
	CompressionMethod                  types.String         `tfsdk:"compression_method"`
	TimeStampIsAccurate                types.Bool           `tfsdk:"time_stamp_is_accurate"`
	OriginalExpiryTime                 types.Int64          `tfsdk:"original_expiry_time"`
	VolumeReplicationState             types.String         `tfsdk:"volume_replication_state"`
	ReplicationJournalVolume           types.Bool           `tfsdk:"replication_journal_volume"`
	ReplicationTimeStamp               types.Int64          `tfsdk:"replication_time_stamp"`
	Links                              []volumeLinkModel    `tfsdk:"links"`
	MappedSdcInfo                      []mappedSdcInfoModel `tfsdk:"mapped_sdc_info"`
}

type volumeLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

type mappedSdcInfoModel struct {
	SdcID                 types.String `tfsdk:"sdc_id"`
	SdcIP                 types.String `tfsdk:"sdc_ip"`
	LimitIops             types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps         types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName               types.String `tfsdk:"sdc_name"`
	AccessMode            types.String `tfsdk:"access_mode"`
	IsDirectBufferMapping types.Bool   `tfsdk:"is_direct_buffer_mapping"`
}

func (d *volumeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (d *volumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VolumeDataSourceSchema
}

func (d *volumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *volumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state volumeDataSourceModel
	var volumes []*scaleiotypes.Volume
	var err error

	diags := req.Config.Get(ctx, &state)

	//Read the volumes based on volume id/name or storage pool id/name and if nothing
	//is mentioned , then return all volumes
	if state.Name.ValueString() != "" {
		volumes, err = d.client.GetVolume("", "", "", state.Name.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
	} else if state.ID.ValueString() != "" {
		volumes, err = d.client.GetVolume("", state.ID.ValueString(), "", "", false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
	} else if state.StoragePoolID.ValueString() != "" {
		sps, err1 := d.client.FindStoragePool(state.StoragePoolID.ValueString(), "", "", "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
		sp := goscaleio.NewStoragePool(d.client)
		sp.StoragePool = sps
		volumes, err = sp.GetVolume("", "", "", "", false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
	} else if state.StoragePoolName.ValueString() != "" {
		sps, err1 := d.client.FindStoragePool("", state.StoragePoolName.ValueString(), "", "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
		sp := goscaleio.NewStoragePool(d.client)
		sp.StoragePool = sps
		volumes, err = sp.GetVolume("", "", "", "", false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
	} else {
		volumes, err = d.client.GetVolume("", "", "", "", false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err.Error(),
			)
			return
		}
	}

	state.Volumes = updateVolumeState(volumes)
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// updateVolumeState iterates over the volume list and update the state
func updateVolumeState(volumes []*scaleiotypes.Volume) (response []volumeModel) {
	for _, volumeValue := range volumes {
		volumeState := volumeModel{
			ID:                                 types.StringValue(volumeValue.ID),
			Name:                               types.StringValue(volumeValue.Name),
			CreationTime:                       types.Int64Value((int64)(volumeValue.CreationTime)),
			SizeInKb:                           types.Int64Value((int64)(volumeValue.SizeInKb)),
			AncestorVolumeID:                   types.StringValue(volumeValue.AncestorVolumeID),
			VTreeID:                            types.StringValue(volumeValue.VTreeID),
			ConsistencyGroupID:                 types.StringValue(volumeValue.ConsistencyGroupID),
			VolumeType:                         types.StringValue(volumeValue.VolumeType),
			UseRmCache:                         types.BoolValue(volumeValue.UseRmCache),
			StoragePoolID:                      types.StringValue(volumeValue.StoragePoolID),
			DataLayout:                         types.StringValue(volumeValue.DataLayout),
			NotGenuineSnapshot:                 types.BoolValue(volumeValue.NotGenuineSnapshot),
			AccessModeLimit:                    types.StringValue(volumeValue.AccessModeLimit),
			SecureSnapshotExpTime:              types.Int64Value((int64)(volumeValue.SecureSnapshotExpTime)),
			ManagedBy:                          types.StringValue(volumeValue.ManagedBy),
			LockedAutoSnapshot:                 types.BoolValue(volumeValue.LockedAutoSnapshot),
			LockedAutoSnapshotMarkedForRemoval: types.BoolValue(volumeValue.LockedAutoSnapshotMarkedForRemoval),
			CompressionMethod:                  types.StringValue(volumeValue.CompressionMethod),
			TimeStampIsAccurate:                types.BoolValue(volumeValue.TimeStampIsAccurate),
			OriginalExpiryTime:                 types.Int64Value((int64)(volumeValue.OriginalExpiryTime)),
			VolumeReplicationState:             types.StringValue(volumeValue.VolumeReplicationState),
			ReplicationJournalVolume:           types.BoolValue(volumeValue.ReplicationJournalVolume),
			ReplicationTimeStamp:               types.Int64Value((int64)(volumeValue.ReplicationTimeStamp)),
		}

		for _, link := range volumeValue.Links {
			volumeState.Links = append(volumeState.Links, volumeLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		for _, sdc := range volumeValue.MappedSdcInfo {
			volumeState.MappedSdcInfo = append(volumeState.MappedSdcInfo, mappedSdcInfoModel{
				SdcID:                 types.StringValue(sdc.SdcID),
				SdcIP:                 types.StringValue(sdc.SdcIP),
				LimitIops:             types.Int64Value((int64)(sdc.LimitIops)),
				LimitBwInMbps:         types.Int64Value((int64)(sdc.LimitBwInMbps)),
				SdcName:               types.StringValue(sdc.SdcName),
				AccessMode:            types.StringValue(sdc.AccessMode),
				IsDirectBufferMapping: types.BoolValue(sdc.IsDirectBufferMapping),
			})
		}
		response = append(response, volumeState)
	}
	return
}
