package volume

import (
	"context"

	// "time"

	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"

	// "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewvolumeResource is a helper function to simplify the provider implementation.
func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

// volumeResource is the resource implementation.
type volumeResource struct {
	client *goscaleio.Client
}

// volumeResourceModel maps the resource schema data.
type volumeResourceModel struct {
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	StoragePoolID      types.String `tfsdk:"storage_pool_id"`
	VolumeType         types.String `tfsdk:"volume_type"`
	// UseRmCache                         types.Bool   `tfsdk:"use_rm_cache"`

	VolumeSizeInKb types.String `tfsdk:"volume_size_in_kb"`
	Name           types.String `tfsdk:"name"`
	// MappingToAllSdcsEnabled            types.Bool   `tfsdk:"mapping_to_all_sdcs_enabled"`
	// MappedSdcInfo                      types.List   `tfsdk:"mapped_sdc_info"`
	// IsObfuscated                       types.Bool   `tfsdk:"is_obfuscated"`
	ConsistencyGroupID types.String `tfsdk:"consistency_group_id"`
	VTreeID            types.String `tfsdk:"vtree_id"`
	AncestorVolumeID   types.String `tfsdk:"ancestor_volume_id"`
	// MappedScsiInitiatorInfo            types.String `tfsdk:"mapped_scsi_initiator_info"`
	SizeInKb                           types.Int64  `tfsdk:"size_in_kb"`
	CreationTime                       types.Int64  `tfsdk:"creation_time"`
	ID                                 types.String `tfsdk:"id"`
	DataLayout                         types.String `tfsdk:"data_layout"`
	NotGenuineSnapshot                 types.Bool   `tfsdk:"not_genuine_snapshot"`
	AccessModeLimit                    types.String `tfsdk:"access_mode_limit"`
	SecureSnapshotExpTime              types.Int64  `tfsdk:"secure_snapshot_exp_time"`
	ManagedBy                          types.String `tfsdk:"managed_by"`
	LockedAutoSnapshot                 types.Bool   `tfsdk:"locked_auto_snapshot"`
	LockedAutoSnapshotMarkedForRemoval types.Bool   `tfsdk:"locked_auto_snapshot_marked_for_removal"`
	CompressionMethod                  types.String `tfsdk:"compression_method"`
	TimeStampIsAccurate                types.Bool   `tfsdk:"time_stamp_is_accurate"`
	OriginalExpiryTime                 types.Int64  `tfsdk:"original_expiry_time"`
	VolumeReplicationState             types.String `tfsdk:"volume_replication_state"`
	ReplicationJournalVolume           types.Bool   `tfsdk:"replication_journal_volume"`
	ReplicationTimeStamp               types.Int64  `tfsdk:"replication_time_stamp"`
	// Links                              types.List   `tfsdk:"links"`
}

type MappedSdcInfo struct {
	SdcID                 types.String `tfsdk:"sdc_id"`
	SdcIP                 types.String `tfsdk:"sdc_ip"`
	LimitIops             types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps         types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName               types.String `tfsdk:"sdc_name"`
	AccessMode            types.String `tfsdk:"access_mode"`
	IsDirectBufferMapping types.Bool   `tfsdk:"is_direct_buffer_mapping"`
}

// Link defines struct of Link
type Link struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// Metadata returns the data source type name.
func (r *volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema defines the schema for the data source.
func (r *volumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an volume.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "name",
				Required:    true,
			},
			"storage_pool_id": schema.StringAttribute{
				Description: "storage pool id",
				Required:    true,
			},
			"protection_domain_id": schema.StringAttribute{
				Description: "protection domain id",
				Required:    true,
			},
			"volume_size_in_kb": schema.StringAttribute{
				Description: "volume size in kb",
				Required:    true,
			},
			"volume_type": schema.StringAttribute{
				Description: "volume type",
				Computed:    true,
			},
			// "use_rm_cache": schema.StringAttribute{
			// 	Description: "use rm cache",
			// 	Optional:    true,
			// },
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"creation_time": schema.Int64Attribute{
				Description: "Creation Time",
				Computed:    true,
			},
			"size_in_kb": schema.Int64Attribute{
				Description: "Size in KB",
				Computed:    true,
			},
			"ancestor_volume_id": schema.StringAttribute{
				Description: "ancestor volume id",
				Computed:    true,
			},
			"vtree_id": schema.StringAttribute{
				Description: "v tree id",
				Computed:    true,
			},
			"consistency_group_id": schema.StringAttribute{
				Description: "consistency group id",
				Computed:    true,
			},
			"data_layout": schema.StringAttribute{
				Description: "data layout",
				Computed:    true,
			},
			"not_genuine_snapshot": schema.BoolAttribute{
				Description: "not genuine snapshot",
				Computed:    true,
			},
			"access_mode_limit": schema.StringAttribute{
				Description: "access mode limit",
				Computed:    true,
			},
			"secure_snapshot_exp_time": schema.Int64Attribute{
				Description: "secure snapshot exp time",
				Computed:    true,
			},
			"managed_by": schema.StringAttribute{
				Description: "manged by",
				Computed:    true,
			},
			"locked_auto_snapshot": schema.BoolAttribute{
				Description: "locked auto snapshot",
				Computed:    true,
			},
			"locked_auto_snapshot_marked_for_removal": schema.BoolAttribute{
				Description: "locked auto snapshot marked for removal",
				Computed:    true,
			},
			"compression_method": schema.StringAttribute{
				Description: "compression method",
				Computed:    true,
			},
			"time_stamp_is_accurate": schema.BoolAttribute{
				Description: "time stamp is accurate",
				Computed:    true,
			},
			"original_expiry_time": schema.Int64Attribute{
				Description: "original expiry time",
				Computed:    true,
			},
			"volume_replication_state": schema.StringAttribute{
				Description: "volume replication state",
				Computed:    true,
			},
			"replication_journal_volume": schema.BoolAttribute{
				Description: "replication journal volume",
				Computed:    true,
			},
			"replication_time_stamp": schema.Int64Attribute{
				Description: "replication time stamp",
				Computed:    true,
			},
			// "links": schema.ListAttribute{
			// 	Computed: true,
			// 	ElementType: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"rel":  types.StringType,
			// 			"href": types.StringType,
			// 		},
			// 	},
			// },
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *volumeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *volumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	volumeCreate := &pftypes.VolumeParam{
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
		StoragePoolID:      plan.StoragePoolID.ValueString(),
		// UseRmCache:         strconv.FormatBool(plan.UseRmCache.ValueBool()),
		VolumeType:     plan.VolumeType.ValueString(),
		VolumeSizeInKb: plan.VolumeSizeInKb.ValueString(),
		Name:           plan.Name.ValueString(),
	}
	getSystems, _ := r.client.GetSystems()
	sr := goscaleio.NewSystem(r.client)
	sr.System = getSystems[0]
	getProtectionDomains, _ := sr.GetProtectionDomain("")
	tflog.Info(ctx, "2. [POWERFLEX] volume Resource State"+helper.PrettyJSON((getSystems[0])))
	for _, protectionDomain := range getProtectionDomains {
		pdr := goscaleio.NewProtectionDomain(r.client)
		pdr.ProtectionDomain = protectionDomain
		tflog.Info(ctx, "hello"+pdr.ProtectionDomain.ID+" "+plan.ProtectionDomainID.ValueString())
		if pdr.ProtectionDomain.ID == plan.ProtectionDomainID.ValueString() {
			getStoragePool, _ := pdr.GetStoragePool("")
			tflog.Info(ctx, "selected"+pdr.ProtectionDomain.ID+" "+plan.ProtectionDomainID.ValueString())
			for _, sp := range getStoragePool {
				spr := goscaleio.NewStoragePool(r.client)
				spr.StoragePool = sp
				tflog.Info(ctx, spr.StoragePool.ID+" "+plan.StoragePoolID.ValueString())
				if spr.StoragePool.ID == plan.StoragePoolID.ValueString() {
					tflog.Info(ctx, "selected : "+spr.StoragePool.ID+" "+plan.StoragePoolID.ValueString())
					volCreateResponse, err1 := spr.CreateVolume(volumeCreate)
					if err1 != nil {
						resp.Diagnostics.AddError(
							"Error creating volume",
							"Could not create volume, unexpected error: "+err1.Error(),
						)
						return
					}
					// plan.ID = types.StringValue(volCreateResponse.ID)
					volsResponse, err2 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
					if err2 != nil {
						resp.Diagnostics.AddError(
							"Error getting volume after creation",
							"Could not get volume, unexpected error: "+err2.Error(),
						)
						return
					}
					tflog.Info(ctx, "[Anshuman] volume Resource State"+helper.PrettyJSON((volsResponse[0])))
					vol := volsResponse[0]
					spi := types.StringValue(vol.StoragePoolID)
					tflog.Info(ctx, "[Anshuman-SPI] volume Resource State"+spi.ValueString())
					res := volumeResourceModel{
						ProtectionDomainID: plan.ProtectionDomainID,
						VolumeSizeInKb:     plan.VolumeSizeInKb,
						// VolumeType: plan.VolumeType,
						StoragePoolID: types.StringValue(vol.StoragePoolID),
						// UseRmCache:                         types.BoolValue(vol.UseRmCache),
						// MappingToAllSdcsEnabled:            types.BoolValue(vol.MappingToAllSdcsEnabled),
						// IsObfuscated:                       types.BoolValue(vol.IsObfuscated),
						VolumeType:         types.StringValue(vol.VolumeType),
						ConsistencyGroupID: types.StringValue(vol.ConsistencyGroupID),
						VTreeID:            types.StringValue(vol.VTreeID),
						AncestorVolumeID:   types.StringValue(vol.AncestorVolumeID),
						// MappedScsiInitiatorInfo:            types.StringValue(vol.MappedScsiInitiatorInfo),
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
					tflog.Info(ctx, "[Anshuman-Plan] volume Resource State"+helper.PrettyJSON((res)))
					plan = res
					diags = resp.State.Set(ctx, plan)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *volumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *volumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
