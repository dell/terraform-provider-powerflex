package powerflex

import (
	"context"
	"strconv"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snapshotResource{}
	_ resource.ResourceWithConfigure   = &snapshotResource{}
	_ resource.ResourceWithImportState = &snapshotResource{}
)

// NewSnapshotResource is a helper function to simplify the provider implementation.
func NewSnapshotResource() resource.Resource {
	return &snapshotResource{}
}

// snapshotResource is the resource implementation.
type snapshotResource struct {
	client *goscaleio.Client
}

// Metadata returns the resource type name.
func (r *snapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema defines the schema for the resource.
func (r *snapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotResourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *snapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	gs, _ := r.client.GetSystems()
	sr := goscaleio.NewSystem(r.client)
	sr.System = gs[0]
	snapshotReqs := make([]*pftypes.SnapshotDef, 0)
	snapReq := &pftypes.SnapshotDef{
		VolumeID:     plan.VolumeID.ValueString(),
		SnapshotName: plan.Name.ValueString(),
	}
	snapshotReqs = append(snapshotReqs, snapReq)
	snapParam := &pftypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotReqs,
	}
	snapResps, err := sr.CreateSnapshotConsistencyGroup(snapParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			"Could not create snapshot, unexpected err: "+err.Error(),
		)
		return
	}
	snapID := snapResps.VolumeIDList[0]
	snapResponse, _ := r.client.GetVolume("", snapID, "", "", false)
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap
	err = snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting Volume Access Mode",
			"Could not set snapshots, unexpected err: "+err.Error(),
		)
	}
	if !plan.Size.IsNull() {
		vikb, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		tflog.Info(ctx, "vikb"+strconv.FormatInt(vikb, 10))
		if int64(snapResource.Volume.SizeInKb) > vikb {
			resp.Diagnostics.AddError(
				"Error setting the snapshot size",
				"Could not set the size for snapshot below volume size",
			)
			snapResource.RemoveVolume("")
			return
		}
		err3 := snapResource.SetVolumeSize(strconv.FormatInt(plan.Size.ValueInt64(), 10))
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error setting snapshot size",
				"Could not set snapshot size, unexpected err: "+err3.Error(),
			)
			snapResource.RemoveVolume("")
			return
		}

	}
	if plan.LockedAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
			snapResource.RemoveVolume("")
			return
		}
	}
	successMapped := make([]string, 0)
	msids := []string{}
	diags = plan.MapSdcIds.ElementsAs(ctx, &msids, true)
	resp.Diagnostics.Append(diags...)
	for _, msid := range msids {
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msid,
			AllowMultipleMappings: "true",
		}
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Snapshot to SDCs",
				"Could not map Snapshot to scs with id: "+msid+", unexpected error: "+err3.Error(),
			)
			for _, usi := range successMapped {
				snapResource.UnmapVolumeSdc(
					&pftypes.UnmapVolumeSdcParam{
						SdcID: usi,
					},
				)
			}
			snapResource.RemoveVolume("")
			return
		}
		successMapped = append(successMapped, msid)
	}
	snapResponse, _ = r.client.GetVolume("", snapID, "", "", false)
	snap = snapResponse[0]
	state := SnapshotTerraformState(snap, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	snapResponse, _ := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	snap := snapResponse[0]
	state = SnapshotTerraformState(snap, state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	var state SnapshotResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	plan.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))
	snapResponse, _ := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap
	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		errRename := snapResource.SetVolumeName(plan.Name.ValueString())
		if errRename != nil {
			resp.Diagnostics.AddError(
				"Error renaming the snapshot -> "+plan.Name.ValueString()+" : "+state.Name.ValueString(),
				"Could not rename the snapshot, unexpected error:"+errRename.Error(),
			)
			return
		}
	}
	// updating the size of the volume if there is change in plan
	if plan.VolumeSizeInKb.ValueString() != state.VolumeSizeInKb.ValueString() {
		sizeInGb, _ := strconv.Atoi(strconv.FormatInt(VSIKB, 10))
		sizeInGb = sizeInGb / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		err3 := snapResource.SetVolumeSize(sizeInGB)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error setting snapshot size -> "+plan.VolumeSizeInKb.ValueString()+":"+state.VolumeSizeInKb.ValueString(),
				"Could not set new snapshot size -> "+sizeInGB+", unexpected err: "+err3.Error(),
			)
			return
		}
	}
	if plan.AccessMode.ValueString() != state.AccessMode.ValueString() {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Snapshot Access Mode",
				"Could not set the Snapshot Access Mode, unexpected err: "+err.Error(),
			)
			return
		}
	}

	if plan.LockedAutoSnapshot.ValueBool() && !state.LockedAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
		}
	}
	if !plan.LockedAutoSnapshot.ValueBool() && state.LockedAutoSnapshot.ValueBool() {
		err := snapResource.UnlockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unlocking Auto Snapshots",
				"Could not unlock auto snapshots, unexpected error: "+err.Error(),
			)
		}
	}
	planSdcIds := []string{}
	stateSdcIds := []string{}
	diags = plan.MapSdcIds.ElementsAs(ctx, &planSdcIds, true)
	resp.Diagnostics.Append(diags...)

	diags = state.MapSdcIds.ElementsAs(ctx, &stateSdcIds, true)
	resp.Diagnostics.Append(diags...)
	mapSdcIds := Difference(planSdcIds, stateSdcIds)
	unmapSdcIds := Difference(stateSdcIds, planSdcIds)

	for _, msi := range mapSdcIds {
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msi,
			AllowMultipleMappings: "true",
		}
		err3 := snapResource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Snapshots to SDCs",
				"Could not map snapshot to scs with id: "+msi+", unexpected error: "+err3.Error(),
			)
			return
		}
	}

	for _, usi := range unmapSdcIds {
		err4 := snapResource.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: usi,
			},
		)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Snapshot to SDCs",
				"Could not Unmap snapshot to scs with id: "+usi+", unexpected error: "+err4.Error(),
			)
			return
		}
	}
	snapResponse, _ = r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	snap = snapResponse[0]
	snapResource.Volume = snap
	state = SnapshotTerraformState(snap, plan)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	snapResponse, _ := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	snapshot := goscaleio.NewVolume(r.client)
	snapshot.Volume = snapResponse[0]
	sdcsToUnmap := []string{}
	diags = state.MapSdcIds.ElementsAs(ctx, &sdcsToUnmap, true)
	resp.Diagnostics.Append(diags...)
	for _, stu := range sdcsToUnmap {
		err := snapshot.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: stu,
			},
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Volume to SDCs",
				"Couldn't unmap volume to scs with id: "+stu+", unexpected error: "+err.Error(),
			)
			return
		}
	}
	err := snapshot.RemoveVolume("")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing Volume",
			"Couldn't remove volume "+err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
